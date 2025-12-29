// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package rest

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/jsonapi"
	"github.com/playbymail/ottoapp/backend/domains"
	"github.com/playbymail/ottoapp/backend/iana"
	"github.com/playbymail/ottoapp/backend/restapi"
	"github.com/playbymail/ottoapp/backend/services/authz"
	"github.com/playbymail/ottoapp/backend/services/users"
)

// handleGetProfile returns the current user's profile information
// using the users service to fetch fresh data from the database.
func handleGetProfile(authzSvc *authz.Service, usersSvc *users.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type response struct {
			Username string   `json:"username,omitempty"`
			Email    string   `json:"email,omitempty"`
			Timezone string   `json:"timezone,omitempty"`
			Errors   []string `json:"errors,omitempty"`
		}

		actor, err := authzSvc.GetActor(r)
		if err != nil || !actor.IsValid() {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		// Fetch user data from the users service
		user, err := usersSvc.GetUserByID(actor.ID)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// Create a response payload with the fields the frontend needs
		profileData := response{
			Username: user.Username,
			Email:    user.Email,
			Timezone: user.Locale.Timezone.Location.String(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(profileData)
	}
}

// handleGetMyProfile returns the current user's profile.
// GET /api/my/profile
func handleGetMyProfile(authzSvc *authz.Service, usersSvc *users.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		actor, err := authzSvc.GetActor(r)
		if err != nil || !actor.IsValid() {
			restapi.WriteJsonApiError(w, http.StatusUnauthorized, "not_authenticated", "Unauthorized", "Sign in to access this resource.")
			return
		}

		user, err := usersSvc.GetUserByID(actor.ID)
		if err != nil {
			log.Printf("%s %s: GetUserByID(%d) %v", r.Method, r.URL, actor.ID, err)
			restapi.WriteJsonApiError(w, http.StatusInternalServerError, "server_error", "Internal Server Error", "")
			return
		}
		if user == nil {
			log.Printf("%s %s: GetUserByID(%d) %v", r.Method, r.URL, actor.ID, err)
			restapi.WriteJsonApiError(w, http.StatusInternalServerError, "server_error", "Internal Server Error", "")
			return
		}

		restapi.WriteJsonApiData(w, http.StatusOK, usersSvc.UserView(user, actor, actor))
	}
}

// handlePostProfile updates the current user's profile information.
// Users can update their email and timezone, but not their username.
func handlePostProfile(authzSvc *authz.Service, ianaSvc *iana.Service, usersSvc *users.Service) http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Timezone string `json:"timezone"`
	}
	type response struct {
		Handle   string   `json:"handle,omitempty"`
		Username string   `json:"username,omitempty"`
		Email    string   `json:"email,omitempty"`
		Timezone string   `json:"timezone,omitempty"`
		Errors   []string `json:"errors,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		actor, err := authzSvc.GetActor(r)
		if err != nil || !actor.IsValid() {
			restapi.WriteJsonApiError(w, http.StatusUnauthorized, "not_authenticated", "Unauthorized", "Sign in to access this resource.")
			return
		}

		target, err := authzSvc.GetActorById(domains.ID(actor.ID))
		if err != nil || !target.IsValid() {
			// don't leak valid or invalid id
			restapi.WriteJsonApiError(w, http.StatusForbidden, "forbidden", "Forbidden", "You do not have permission to update this user.")
			return
		}
		if !authzSvc.CanEditTarget(actor, target) {
			restapi.WriteJsonApiError(w, http.StatusForbidden, "forbidden", "Forbidden", "You do not have permission to update this user.")
			return
		}

		user, err := usersSvc.GetUserByID(target.ID)
		if err != nil {
			log.Printf("%s %s: get target: %v", r.Method, r.URL.Path, err)
			restapi.WriteJsonApiError(w, http.StatusInternalServerError, "server_error", "Internal Server Error", "")
			return
		}
		if user == nil {
			log.Printf("%s %s: get target: not found", r.Method, r.URL.Path)
			restapi.WriteJsonApiError(w, http.StatusInternalServerError, "server_error", "Internal Server Error", "")
			return
		}
		updatedUser := &domains.User_t{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Handle:   user.Handle,
			Locale: domains.UserLocale_t{
				DateFormat: user.Locale.DateFormat,
				Timezone: domains.UserTimezone_t{
					Location: user.Locale.Timezone.Location,
				},
			},
			Roles:   user.Roles,
			Created: user.Created,
			Updated: time.Now().UTC(),
		}

		// Parse the request body and update our user
		var req request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(response{
				Errors: []string{"Invalid request body"},
			})
			return
		}

		var validationErrors []*jsonapi.ErrorObject
		changedEmail := false
		if req.Email != user.Email {
			if err := domains.ValidateEmail(req.Email); err != nil {
				validationErrors = append(validationErrors, &jsonapi.ErrorObject{
					Status: strconv.Itoa(http.StatusUnprocessableEntity),
					Code:   "invalid_email",
					Title:  "Invalid email",
					Detail: "Provide a valid email.",
					Source: &jsonapi.ErrorSource{
						Pointer: "/data/attributes/email",
					},
				})
			} else {
				changedEmail = true
				updatedUser.Email = req.Email
			}
		}
		changedTimezone := false
		if req.Timezone != user.Locale.Timezone.Location.String() {
			if newTimezone, err := time.LoadLocation(req.Timezone); err != nil {
				validationErrors = append(validationErrors, &jsonapi.ErrorObject{
					Status: strconv.Itoa(http.StatusUnprocessableEntity),
					Code:   "invalid_timezone",
					Title:  "Invalid timezone",
					Detail: "Use an IANA timezone (e.g., America/Chicago).",
					Source: &jsonapi.ErrorSource{
						Pointer: "/data/attributes/timezone",
					},
				})
			} else {
				changedTimezone = true
				updatedUser.Locale.Timezone.Location = newTimezone
			}
		}
		// return if there are any validation errors
		if validationErrors != nil {
			restapi.WriteJsonApiErrorObjects(w, http.StatusUnprocessableEntity, validationErrors...)
			return
		}
		// return if there are no changes
		if !(changedEmail || changedTimezone) {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		err = usersSvc.UpdateUser(updatedUser)
		if err != nil {
			log.Printf("%s %s: update: %v", r.Method, r.URL.Path, err)
			restapi.WriteJsonApiError(w, http.StatusInternalServerError, "server_error", "Internal Server Error", "")
			return
		}

		// Return updated user
		log.Printf("%s %s: get user by id: %v", r.Method, r.URL.Path, err)
		restapi.WriteJsonApiData(w, http.StatusOK, usersSvc.UserView(updatedUser, actor, target))
	}
}
