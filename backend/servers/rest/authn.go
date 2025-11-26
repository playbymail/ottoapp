// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/playbymail/ottoapp/backend/domains"
	"github.com/playbymail/ottoapp/backend/restapi"
	"github.com/playbymail/ottoapp/backend/services/authn"
	"github.com/playbymail/ottoapp/backend/services/authz"
	"github.com/playbymail/ottoapp/backend/sessions"
	"github.com/playbymail/ottoapp/backend/users"
)

// HandleGetSession returns the current session
func HandleGetSession(authzSvc *authz.Service, sessionsSvc *sessions.Service) http.HandlerFunc {
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

// HandlePostLogin creates a session and sets the cookie.
func HandlePostLogin(authnSvc *authn.Service, authzSvc *authz.Service, sessionsSvc *sessions.Service, usersSvc *users.Service) http.HandlerFunc {
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

		// authenticate
		actorId, err := authnSvc.AuthenticateEmailCredentials(body.Email, body.Password)
		if err != nil {
			log.Printf("%s %s: auth(%q, %q): %v\n", r.Method, r.URL.Path, body.Email, body.Password, err)
			<-timer.C
			restapi.WriteJsonApiError(w, http.StatusUnauthorized, "invalid_credentials", "Invalid Credentials", "Provide valid credentials to access this resource.")
			return
		}
		user, err := usersSvc.GetUserByID(actorId)
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
