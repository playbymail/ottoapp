// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package rest

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/jsonapi"
	"github.com/mdhender/phrases/v2"
	"github.com/playbymail/ottoapp/backend/domains"
	"github.com/playbymail/ottoapp/backend/restapi"
	"github.com/playbymail/ottoapp/backend/services/authn"
	"github.com/playbymail/ottoapp/backend/services/authz"
	"github.com/playbymail/ottoapp/backend/services/users"
)

// handleGetMe returns the current user's profile.
// For completeness; use 401 when no session, 404 if record missing.
func handleGetMe(authzSvc *authz.Service, usersSvc *users.Service) http.HandlerFunc {
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
			log.Printf("%s %s: GetUserByID(%d) not found", r.Method, r.URL, actor.ID)
			restapi.WriteJsonApiError(w, http.StatusInternalServerError, "server_error", "Internal Server Error", "")
			return
		}

		restapi.WriteJsonApiData(w, http.StatusOK, usersSvc.UserView(user, actor, actor))
	}
}

// handleGetUser returns a specific user's profile (with RBAC).
func handleGetUser(authzSvc *authz.Service, usersSvc *users.Service) http.HandlerFunc {
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
		if err != nil || !actor.IsValid() {
			// don't leak valid or invalid id
			restapi.WriteJsonApiError(w, http.StatusForbidden, "forbidden", "Forbidden", "You do not have permission to view this user.")
			return
		}

		view, err := usersSvc.ReadUser(actor, target)
		if err != nil {
			if errors.Is(err, domains.ErrNotExists) || errors.Is(err, domains.ErrNotFound) || errors.Is(err, domains.ErrNotAuthorized) {
				restapi.WriteJsonApiError(w, http.StatusForbidden, "forbidden", "Forbidden", "You do not have permission to view this user.")
				return
			}
			restapi.WriteJsonApiDatabaseError(w)
			return
		}

		restapi.WriteJsonApiData(w, http.StatusOK, view)
	}
}

// handleGetUsers returns a list of all non-admin, non-sysop users (admin only)
// GET /api/users
func handleGetUsers(authzSvc *authz.Service, usersSvc *users.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		actor, err := authzSvc.GetActor(r)
		if err != nil || !actor.IsValid() {
			restapi.WriteJsonApiError(w, http.StatusUnauthorized, "not_authenticated", "Unauthorized", "Sign in to access this resource.")
			return
		}

		views, err := usersSvc.ReadUsers(actor)
		if err != nil {
			if errors.Is(err, domains.ErrNotAuthorized) {
				restapi.WriteJsonApiError(w, http.StatusForbidden, "forbidden", "Forbidden", "You do not have permission to view user lists.")
				return
			}
			restapi.WriteJsonApiDatabaseError(w)
			return
		}

		restapi.WriteJsonApiData(w, http.StatusOK, views)
	}
}

// handlePatchUserRelationshipRoles is a route to make JSON:API happier
// PATCH /api/users/:id/relationships/roles
//
//	{
//	 "data": [ { "type": "roles", "id": "admin" }, ... ]
//	}
func handlePatchUserRelationshipRoles() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
	}
}

// handlePatchUser updates a user's profile (with RBAC).
// Ember Data lies and sends the entire record up, not just the change set.
// PATCH /api/users/:id
func handlePatchUser(authzSvc *authz.Service, usersSvc *users.Service) http.HandlerFunc {
	type userPatchRequest struct {
		ID       string `json:"id"`
		Username string `json:"username,omitempty"`
		Email    string `json:"email,omitempty"`
		Timezone string `json:"timezone,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		targetID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil || domains.ID(targetID) == domains.InvalidID {
			restapi.WriteJsonApiMalformedPathParameter(w, "user_id", "User ID", r.PathValue("id"))
			return
		}

		actor, err := authzSvc.GetActor(r)
		if err != nil || !actor.IsValid() {
			restapi.WriteJsonApiError(w, http.StatusUnauthorized, "not_authenticated", "Unauthorized", "Sign in to access this resource.")
			return
		}
		target, err := authzSvc.GetActorById(domains.ID(targetID))
		if err != nil || !target.IsValid() {
			// don't leak valid or invalid id
			restapi.WriteJsonApiError(w, http.StatusForbidden, "forbidden", "Forbidden", "You do not have permission to view this user.")
			return
		}

		if !authzSvc.CanEditTarget(actor, target) {
			restapi.WriteJsonApiError(w, http.StatusForbidden, "forbidden", "Forbidden", "You do not have permission to update this user.")
			return
		}

		var p userPatchRequest
		if err := jsonapi.UnmarshalPayload(r.Body, &p); err != nil {
			log.Printf("%s %s: update: %v", r.Method, r.URL.Path, err)
			restapi.WriteJsonApiError(w, http.StatusBadRequest, "bad_request", "Invalid Request Body", err.Error())
			return
		}
		log.Printf("%s %s: unmarshal %+v\n", r.Method, r.URL.Path, p)

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
		// UpdatePassword user
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

		// Calculate the change set and capture any validation errors
		var validationErrors []*jsonapi.ErrorObject
		changedUsername := false
		if p.Username != user.Username {
			if !authzSvc.CanEditTargetUsername(actor, target) {
				restapi.WriteJsonApiError(w, http.StatusForbidden, "forbidden", "Forbidden", "You do not have permission to update user names.")
				return
			}
			if err := domains.ValidateUsername(p.Username); err != nil {
				validationErrors = append(validationErrors, &jsonapi.ErrorObject{
					Status: strconv.Itoa(http.StatusUnprocessableEntity),
					Code:   "invalid_username",
					Title:  "Invalid username",
					Detail: "Provide a valid username.",
					Source: &jsonapi.ErrorSource{
						Pointer: "/data/attributes/username",
					},
				})
			} else {
				changedUsername = true
				user.Username = p.Username
			}
		}
		changedEmail := false
		if p.Email != user.Email {
			if err := domains.ValidateEmail(p.Email); err != nil {
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
				user.Email = p.Email
			}
		}
		changedTimezone := false
		if p.Timezone != user.Locale.Timezone.Location.String() {
			if newTimezone, err := time.LoadLocation(p.Timezone); err != nil {
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
		if !(changedUsername || changedEmail || changedTimezone) {
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

// handlePatchUserRole updates a user's role (admin only)
// PATCH /api/users/:id/role
func handlePatchUserRole(authzSvc *authz.Service, usersSvc *users.Service) http.HandlerFunc {
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
		var req struct {
			Add    []string `json:"add,omitempty"`
			Remove []string `json:"remove,omitempty"`
		}
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
			restapi.WriteJsonApiError(w, http.StatusForbidden, "forbidden", "Forbidden", "You do not have permission to manage this user.")
			return
		}

		if !authzSvc.CanManageTargetRoles(actor, target) {
			restapi.WriteJsonApiError(w, http.StatusForbidden, "forbidden", "Forbidden", "You do not have permission to manage this user.")
			return
		}

		// Add roles
		for _, roleID := range req.Add {
			err = authzSvc.AssignRole(target.ID, roleID)
			if err != nil {
				log.Printf("PATCH /api/users/%d/role: add %q: %v", target.ID, roleID, err)
				restapi.WriteJsonApiError(w, http.StatusInternalServerError, "server_error", "Internal Server Error", "")
				return
			}
		}

		// Remove roles
		for _, roleID := range req.Remove {
			err = authzSvc.RemoveRole(target.ID, roleID)
			if err != nil {
				log.Printf("PATCH /api/users/%d/role: remove %q: %v", target.ID, roleID, err)
				restapi.WriteJsonApiError(w, http.StatusInternalServerError, "server_error", "Internal Server Error", "")
				return
			}
		}

		// Return updated user
		user, err := usersSvc.GetUserByID(target.ID)
		if err != nil {
			restapi.WriteJsonApiError(w, http.StatusInternalServerError, "server_error", "Internal Server Error", "")
			return
		}

		restapi.WriteJsonApiData(w, http.StatusOK, usersSvc.UserView(user, actor, target))
	}
}

// handlePostUser creates a new user (admin only)
// POST /api/users
func handlePostUser(authnSvc *authn.Service, authzSvc *authz.Service, usersSvc *users.Service) http.HandlerFunc {
	// userCreateRequest is the request payload for creating a user
	type userCreateRequest struct {
		Handle   string   `json:"handle,omitempty"`
		Email    string   `json:"email,omitempty"`
		Username string   `json:"username,omitempty"`
		Password string   `json:"password,omitempty"`
		Timezone string   `json:"timezone,omitempty"`
		Roles    []string `json:"roles,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		actor, err := authzSvc.GetActor(r)
		if err != nil || !actor.IsValid() {
			restapi.WriteJsonApiError(w, http.StatusUnauthorized, "not_authenticated", "Unauthorized", "Sign in to access this resource.")
			return
		}

		if !authzSvc.CanCreateTarget(actor) {
			restapi.WriteJsonApiError(w, http.StatusForbidden, "forbidden", "Forbidden", "You do not have permission to create users.")
			return
		}

		var p userCreateRequest
		if err := jsonapi.UnmarshalPayload(r.Body, &p); err != nil {
			log.Printf("%s %s: update: %v", r.Method, r.URL.Path, err)
			restapi.WriteJsonApiError(w, http.StatusBadRequest, "bad_request", "Invalid Request Body", err.Error())
			return
		}
		var validationErrors []*jsonapi.ErrorObject
		if p.Handle == "" {
			validationErrors = append(validationErrors, &jsonapi.ErrorObject{
				Status: strconv.Itoa(http.StatusUnprocessableEntity),
				Code:   "missing_handle",
				Title:  "Missing handle",
				Detail: "Provide a valid handle.",
				Source: &jsonapi.ErrorSource{
					Pointer: "/data/attributes/handle",
				},
			})
		} else if err := domains.ValidateHandle(p.Handle); err != nil {
			validationErrors = append(validationErrors, &jsonapi.ErrorObject{
				Status: strconv.Itoa(http.StatusUnprocessableEntity),
				Code:   "missing_handle",
				Title:  "Missing handle",
				Detail: "Provide a valid handle.",
				Source: &jsonapi.ErrorSource{
					Pointer: "/data/attributes/handle",
				},
			})
		}
		if p.Username == "" {
			validationErrors = append(validationErrors, &jsonapi.ErrorObject{
				Status: strconv.Itoa(http.StatusUnprocessableEntity),
				Code:   "missing_username",
				Title:  "Missing username",
				Detail: "Provide a valid username.",
				Source: &jsonapi.ErrorSource{
					Pointer: "/data/attributes/username",
				},
			})
		} else if err := domains.ValidateUsername(p.Username); err != nil {
			validationErrors = append(validationErrors, &jsonapi.ErrorObject{
				Status: strconv.Itoa(http.StatusUnprocessableEntity),
				Code:   "invalid_username",
				Title:  "Invalid username",
				Detail: "Provide a valid username.",
				Source: &jsonapi.ErrorSource{
					Pointer: "/data/attributes/username",
				},
			})
		}
		if p.Email == "" {
			validationErrors = append(validationErrors, &jsonapi.ErrorObject{
				Status: strconv.Itoa(http.StatusUnprocessableEntity),
				Code:   "missing_email",
				Title:  "missing email",
				Detail: "Provide a valid email.",
				Source: &jsonapi.ErrorSource{
					Pointer: "/data/attributes/email",
				},
			})
		} else if err := domains.ValidateEmail(p.Email); err != nil {
			validationErrors = append(validationErrors, &jsonapi.ErrorObject{
				Status: strconv.Itoa(http.StatusUnprocessableEntity),
				Code:   "invalid_email",
				Title:  "Invalid email",
				Detail: "Provide a valid email.",
				Source: &jsonapi.ErrorSource{
					Pointer: "/data/attributes/email",
				},
			})
		}
		var loc *time.Location
		if p.Timezone == "" {
			validationErrors = append(validationErrors, &jsonapi.ErrorObject{
				Status: strconv.Itoa(http.StatusUnprocessableEntity),
				Code:   "missing_timezone",
				Title:  "Missing timezone",
				Detail: "Use an IANA timezone (e.g., America/Chicago).",
				Source: &jsonapi.ErrorSource{
					Pointer: "/data/attributes/timezone",
				},
			})
		} else if loc, err = time.LoadLocation(p.Timezone); err != nil {
			validationErrors = append(validationErrors, &jsonapi.ErrorObject{
				Status: strconv.Itoa(http.StatusUnprocessableEntity),
				Code:   "invalid_timezone",
				Title:  "Invalid timezone",
				Detail: "Use an IANA timezone (e.g., America/Chicago).",
				Source: &jsonapi.ErrorSource{
					Pointer: "/data/attributes/timezone",
				},
			})
		}
		if validationErrors != nil {
			restapi.WriteJsonApiErrorObjects(w, http.StatusUnprocessableEntity, validationErrors...)
			return
		}

		// Generate password if not provided
		password := p.Password
		if password == "" {
			password = phrases.Generate(6)
		}

		// Create user (this will assign "active" and "user" roles by default)
		// TODO: Add handle field to payload; for now use username as handle
		user, err := usersSvc.CreateUser(p.Handle, p.Email, p.Username, loc)
		if err != nil {
			log.Printf("%s %s: create: %v", r.Method, r.URL.Path, err)
			restapi.WriteJsonApiError(w, http.StatusInternalServerError, "server_error", "Internal Server Error", "")
			return
		}
		target, err := authzSvc.GetActorById(user.ID)
		if err != nil {
			log.Printf("%s %s: %d: %d: create: %v", r.Method, r.URL.Path, actor.ID, user.ID, err)
			restapi.WriteJsonApiError(w, http.StatusInternalServerError, "server_error", "Internal Server Error", "")
			return
		}
		_, err = authnSvc.UpdateCredentials(actor, target, "", password)
		if err != nil {
			log.Printf("%s %s: %d: %d: create: %v", r.Method, r.URL.Path, actor.ID, user.ID, err)
			restapi.WriteJsonApiError(w, http.StatusInternalServerError, "server_error", "Internal Server Error", "")
			return
		}

		// If roles were specified and don't include "user" or "admin", assign "guest" instead
		if len(p.Roles) > 0 {
			hasUserOrAdmin := false
			for _, role := range p.Roles {
				if role == "user" || role == "admin" {
					hasUserOrAdmin = true
					break
				}
			}
			if !hasUserOrAdmin {
				// Remove "user" role and add "guest"
				_ = authzSvc.RemoveRole(user.ID, "user")
				_ = authzSvc.AssignRole(user.ID, "guest")
			}

			// Assign additional roles
			for _, role := range p.Roles {
				if role != "active" && role != "user" { // active and user already assigned
					err = authzSvc.AssignRole(user.ID, role)
					if err != nil {
						log.Printf("%s %s: assign role %q: %v", r.Method, r.URL.Path, role, err)
					}
				}
			}
		}

		view := usersSvc.UserView(user, actor, target)

		// set a Location key in the header to let the client know how to find the new user
		w.Header().Set("Location", restapi.AbsURL(r, "/api/users/"+view.ID))

		restapi.WriteJsonApiData(w, http.StatusCreated, view)
	}
}
