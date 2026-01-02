// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/jsonapi"
	"github.com/mdhender/phrases/v2"
	"github.com/playbymail/ottoapp/backend/domains"
	"github.com/playbymail/ottoapp/backend/restapi"
	"github.com/playbymail/ottoapp/backend/services/authn"
	"github.com/playbymail/ottoapp/backend/services/authz"
	"github.com/playbymail/ottoapp/backend/services/users"
	"github.com/playbymail/ottoapp/backend/sessions"
)

// handleGetSession returns the current session
func handleGetSession(authzSvc *authz.Service, sessionsSvc *sessions.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//log.Printf("%s %s: entered", r.Method, r.URL.Path)
		if r.Method != http.MethodGet {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		ss, err := authzSvc.GetSessionData(r)
		if err != nil {
			log.Printf("%s %s: gsd %v", r.Method, r.URL.Path, err)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		} else if ss.UserID == domains.InvalidID {
			log.Printf("%s %s: gsd 0: no data found", r.Method, r.URL.Path)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(ss)
	}
}

// handlePatchPassword updates the target's credentials
// PATCH /api/users/:id/password
func handlePatchPassword(authnSvc *authn.Service, authzSvc *authz.Service) http.HandlerFunc {
	type patchPasswordRequest struct {
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// Parse target user ID from path
		targetID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil || domains.ID(targetID) == domains.InvalidID {
			restapi.WriteJsonApiErrorObjects(w, http.StatusBadRequest, &jsonapi.ErrorObject{
				Status: strconv.Itoa(http.StatusBadRequest),
				Code:   "invalid_user_id",
				Title:  "Invalid UserID",
				Detail: "Provide a valid UserID.",
				Source: &jsonapi.ErrorSource{
					Parameter: "id",
				},
			})
			return
		}

		// Parse request body
		var req patchPasswordRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			restapi.WriteJsonApiError(w, http.StatusBadRequest, "bad_request", "Invalid Request Body", "")
			return
		}

		actor, err := authzSvc.GetActor(r)
		if err != nil || !actor.IsValid() {
			restapi.WriteJsonApiError(w, http.StatusUnauthorized, "not_authenticated", "Unauthorized", "Sign in to access this resource.")
			return
		}
		target, err := authzSvc.GetActorById(domains.ID(targetID))
		if err != nil || !target.IsValid() {
			restapi.WriteJsonApiError(w, http.StatusForbidden, "forbidden", "Forbidden", "You do not have permission to update user passwords.")
			return
		}

		_, err = authnSvc.UpdateCredentials(actor, target, req.CurrentPassword, req.NewPassword)
		if err != nil {
			restapi.WriteJsonApiError(w, http.StatusForbidden, "forbidden", "Forbidden", "You do not have permission to update user passwords.")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// handlePostLogin creates a session and sets the cookie.
func handlePostLogin(authnSvc *authn.Service, authzSvc *authz.Service, sessionsSvc *sessions.Service, usersSvc *users.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s: entered", r.Method, r.URL.Path)
		if r.Method != http.MethodPost {
			restapi.WriteJsonApiError(w, http.StatusMethodNotAllowed, "method_not_allowed", fmt.Sprintf("%s not allowed", r.Method), "Only POST is allowed.")
			return
		}

		// invalidate current session
		sessionsSvc.HandleInvalidateSession(w, r)

		// fetch credentials from the request body, limiting body to 1kb of data
		r.Body = http.MaxBytesReader(w, r.Body, 1024)
		data, err := io.ReadAll(r.Body)
		if err != nil {
			if err.Error() == "http: request body too large" {
				restapi.WriteJsonApiError(w, http.StatusRequestEntityTooLarge, "too_large", "File Too Large", "File size exceeds 150KB limit.")
				return
			}
			restapi.WriteJsonApiError(w, http.StatusBadRequest, "bad_request", "Bad Request", "Error reading request body.")
			return
		}
		log.Printf("%s %s: body %q\n", r.Method, r.URL.Path, string(data))

		var body struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		err = json.Unmarshal(data, &body)
		if err != nil {
			log.Printf("%s %s: body %v\n", r.Method, r.URL.Path, err)
			restapi.WriteJsonApiError(w, http.StatusUnprocessableEntity, "invalid_body", "Invalid Body", "Could not parse body: "+err.Error())
			return
		}
		body.Email = strings.ToLower(body.Email)
		log.Printf("%s %s: email %q: password %q\n", r.Method, r.URL.Path, body.Email, body.Password)

		// ensure that we spend at least 500ms if the authentication path fails
		timer := time.NewTimer(500 * time.Millisecond)
		defer timer.Stop()

		// authenticate
		actorId, err := authnSvc.AuthenticateEmailCredentials(body.Email, body.Password)
		log.Printf("%s %s: auth(%q, %q): %v\n", r.Method, r.URL.Path, body.Email, body.Password, err)
		if err != nil {
			log.Printf("%s %s: auth(%q, %q): %v\n", r.Method, r.URL.Path, body.Email, body.Password, err)
			<-timer.C
			restapi.WriteJsonApiError(w, http.StatusUnauthorized, "invalid_credentials", "Invalid Credentials", "Provide valid credentials to access this resource.")
			return
		}
		user, err := usersSvc.GetUserByID(actorId)
		log.Printf("%s %s: getUser(%d): %v\n", r.Method, r.URL.Path, actorId, err)
		if err != nil {
			log.Printf("%s %s: getUser(%d): %v\n", r.Method, r.URL.Path, actorId, err)
			<-timer.C
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		sessionsSvc.HandleCreateSession(w, r, user)
		<-timer.C
	}
}

// handlePostResetPassword resets a user's password (admin only)
// POST /api/users/:id/reset-password
func handlePostResetPassword(authnSvc *authn.Service, authzSvc *authz.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse target user ID from path
		targetID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil || domains.ID(targetID) == domains.InvalidID {
			restapi.WriteJsonApiErrorObjects(w, http.StatusBadRequest, &jsonapi.ErrorObject{
				Status: strconv.Itoa(http.StatusBadRequest),
				Code:   "invalid_user_id",
				Title:  "Invalid UserID",
				Detail: "Provide a valid UserID.",
				Source: &jsonapi.ErrorSource{
					Parameter: "id",
				},
			})
			return
		}

		actor, err := authzSvc.GetActor(r)
		if err != nil || !actor.IsValid() {
			restapi.WriteJsonApiError(w, http.StatusUnauthorized, "not_authenticated", "Unauthorized", "Sign in to access this resource.")
			return
		}
		target, err := authzSvc.GetActorById(domains.ID(targetID))
		if err != nil || !target.IsValid() {
			restapi.WriteJsonApiError(w, http.StatusForbidden, "forbidden", "Forbidden", "You do not have permission to reset the user's password.")
			return
		}

		// Generate temporary password reset link
		magicLink := phrases.Generate(6)
		_, err = authnSvc.UpdateCredentials(actor, target, "", magicLink)
		if err != nil {
			log.Printf("%s %s: %d: %d: %v", r.Method, r.URL.Path, actor.ID, target.ID, err)
			restapi.WriteJsonApiError(w, http.StatusForbidden, "forbidden", "Forbidden", "You do not have permission to reset the user's password.")
			return
		}

		restapi.WriteJsonApiData(w, http.StatusOK, &struct {
			Message   string `jsonapi:"attr,message"`
			MagicLink string `jsonapi:"attr,link"`
		}{
			Message:   "Magic link generated",
			MagicLink: magicLink,
		})
	}
}
