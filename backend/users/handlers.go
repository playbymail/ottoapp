// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package users

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/jsonapi"
	"github.com/mdhender/phrases/v2"
	"github.com/playbymail/ottoapp/backend/domains"
	"github.com/playbymail/ottoapp/backend/restapi"
)

// HandleGetMe returns the current user's profile.
// For completeness; use 401 when no session, 404 if record missing.
func (s *Service) HandleGetMe(w http.ResponseWriter, r *http.Request) {
	actor, err := s.authSvc.GetActor(r)
	if err != nil || !actor.IsValid() {
		restapi.WriteJsonApiError(w, http.StatusUnauthorized, "not_authenticated", "Unauthorized", "Sign in to access this resource.")
		return
	}

	user, err := s.GetUserByID(actor.ID)
	if err != nil {
		log.Printf("GET /api/users/me: get me: %v", actor.ID, err)
		restapi.WriteJsonApiError(w, http.StatusInternalServerError, "server_error", "Internal Server Error", "")
		return
	}
	if user == nil {
		log.Printf("GET /api/users/me: get me: not found", actor.ID)
		restapi.WriteJsonApiError(w, http.StatusInternalServerError, "server_error", "Internal Server Error", "")
		return
	}

	view := s.buildUserView(user, actor, actor)

	restapi.WriteJsonApiData(w, http.StatusOK, view)
}

// HandleGetMyProfile returns the current user's profile.
// GET /api/my/profile
func (s *Service) HandleGetMyProfile(w http.ResponseWriter, r *http.Request) {
	actor, err := s.authSvc.GetActor(r)
	if err != nil || !actor.IsValid() {
		restapi.WriteJsonApiError(w, http.StatusUnauthorized, "not_authenticated", "Unauthorized", "Sign in to access this resource.")
		return
	}

	user, err := s.GetUserByID(actor.ID)
	if err != nil {
		log.Printf("GET /api/my/profile: get me: %v", actor.ID, err)
		restapi.WriteJsonApiError(w, http.StatusInternalServerError, "server_error", "Internal Server Error", "")
		return
	}
	if user == nil {
		log.Printf("GET /api/my/profile: get me: not found", actor.ID)
		restapi.WriteJsonApiError(w, http.StatusInternalServerError, "server_error", "Internal Server Error", "")
		return
	}

	view := s.buildUserView(user, actor, actor)

	restapi.WriteJsonApiData(w, http.StatusOK, view)
}

// HandleGetUsers returns a list of all non-admin, non-sysop users (admin only)
// GET /api/users
func (s *Service) HandleGetUsers(w http.ResponseWriter, r *http.Request) {
	actor, err := s.authSvc.GetActor(r)
	if err != nil || !actor.IsValid() {
		restapi.WriteJsonApiError(w, http.StatusUnauthorized, "not_authenticated", "Unauthorized", "Sign in to access this resource.")
		return
	}

	if !s.authSvc.CanListUsers(actor) {
		restapi.WriteJsonApiError(w, http.StatusForbidden, "forbidden", "Forbidden", "You do not have permission to view user lists.")
		return
	}

	// Get all users from database (we'll filter them)
	// TODO: Filtering in Go after GetAllUsers is going to hurt later.
	// Prefer a store method that already enforces visibility/roles.
	// When you add pagination: include links.self/next/prev and meta.total/page/pageSize/returned.
	rows, err := s.db.Queries().GetAllUsers(s.db.Context())
	if err != nil {
		log.Printf("GET /api/users: query: %v", err)
		restapi.WriteJsonApiError(w, http.StatusInternalServerError, "server_error", "Internal Server Error", "")
		return
	}

	var view []*UserView
	for _, row := range rows {
		target, err := s.authSvc.GetActorById(domains.ID(row.UserID))
		if err != nil || !target.IsValid() {
			continue
		}
		if !s.authSvc.CanViewTarget(actor, target) {
			// skip targets the actor is not allowed to view
			continue
		}

		user, err := s.GetUserByID(target.ID)
		if err != nil {
			log.Printf("GET /api/users/%d: get target: %v", target.ID, err)
			restapi.WriteJsonApiError(w, http.StatusInternalServerError, "server_error", "Internal Server Error", "")
			return
		}
		if user == nil {
			log.Printf("GET /api/users/%d: get target: not found", target.ID)
			restapi.WriteJsonApiError(w, http.StatusInternalServerError, "server_error", "Internal Server Error", "")
			return
		}

		view = append(view, s.buildUserView(user, actor, target))
	}

	restapi.WriteJsonApiData(w, http.StatusOK, view)
}

// HandleGetUsersWithPagination returns a user list with optional pagination support
// GET /api/users?page[number]=1&page[size]=25
func (s *Service) HandleGetUsersWithPagination(w http.ResponseWriter, r *http.Request) {
	actor, err := s.authSvc.GetActor(r)
	if err != nil || !actor.IsValid() {
		restapi.WriteJsonApiError(w, http.StatusUnauthorized, "not_authenticated", "Unauthorized", "Sign in to access this resource.")
		return
	}

	if !s.authSvc.CanListUsers(actor) {
		restapi.WriteJsonApiError(w, http.StatusForbidden, "forbidden", "Forbidden", "You do not have permission to view user lists.")
		return
	}

	// Parse pagination (JSON:API convention)
	pageNum := parsePositiveInt(r.URL.Query().Get("page[number]"), 1)
	pageSize := parsePositiveInt(r.URL.Query().Get("page[size]"), 25)
	if pageSize > 200 {
		pageSize = 200
	}

	// Prefer the DB to do RBAC filtering and pagination:
	// total := s.store.CountVisibleUsers(ctx, actorID)
	// rows  := s.store.ListVisibleUsers(ctx, actorID, limit=pageSize, offset=(pageNum-1)*pageSize)
	total := 12 // this is wrong
	users, err := s.ListUsersVisibleToActor(actor, pageNum, pageSize)
	if err != nil {
		log.Printf("GET /api/users?page[number]=%d&page[size]=%d: query: %v", pageNum, pageSize, err)
		restapi.WriteJsonApiError(w, http.StatusInternalServerError, "server_error", "Internal Server Error", "")
		return
	}

	views := make([]*UserView, 0, len(users))
	for _, user := range users {
		views = append(views, user)
	}

	// Marshal JSON:API data array to a buffer
	var dataBuf bytes.Buffer
	if err := jsonapi.MarshalPayload(&dataBuf, views); err != nil {
		restapi.WriteJsonApiError(w, http.StatusInternalServerError, "encode_failed", "Encoding error", err.Error())
		return
	}

	// Unwrap to add links/meta at top-level
	var payload map[string]any
	if err := json.Unmarshal(dataBuf.Bytes(), &payload); err != nil {
		restapi.WriteJsonApiError(w, http.StatusInternalServerError, "encode_failed", "Encoding error", err.Error())
		return
	}

	self := restapi.PaginateURL(r, pageNum, pageSize)
	var next, prev string
	if (pageNum-1)*pageSize+len(views) < int(total) {
		next = restapi.PaginateURL(r, pageNum+1, pageSize)
	}
	if pageNum > 1 {
		prev = restapi.PaginateURL(r, pageNum-1, pageSize)
	}
	payload["links"] = map[string]string{
		"self": self,
		"next": next,
		"prev": prev,
	}
	payload["meta"] = map[string]any{
		"total":    total,
		"page":     pageNum,
		"pageSize": pageSize,
		"returned": len(views),
	}

	out, err := json.Marshal(payload)
	if err != nil {
		log.Printf("[users] GetUsersWithPagination %v\n", err)
		restapi.WriteJsonApiError(w, http.StatusInternalServerError, "encode_failed", "Encoding error", err.Error())
		return
	}

	restapi.WriteJsonApiData(w, http.StatusOK, out)
}

func parsePositiveInt(s string, def int) int {
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil || n <= 0 {
		return def
	}
	return n
}

// HandleGetUser returns a specific user's profile (with RBAC).
func (s *Service) HandleGetUser(w http.ResponseWriter, r *http.Request) {
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

	actor, err := s.authSvc.GetActor(r)
	if err != nil || !actor.IsValid() {
		restapi.WriteJsonApiError(w, http.StatusUnauthorized, "not_authenticated", "Unauthorized", "Sign in to access this resource.")
		return
	}
	target, err := s.authSvc.GetActorById(domains.ID(targetID))
	if err != nil || !actor.IsValid() {
		// don't leak valid or invalid id
		restapi.WriteJsonApiError(w, http.StatusForbidden, "forbidden", "Forbidden", "You do not have permission to view this user.")
		return
	}

	if !s.authSvc.CanViewTarget(actor, target) {
		restapi.WriteJsonApiError(w, http.StatusForbidden, "forbidden", "Forbidden", "You do not have permission to view this user.")
		return
	}

	user, err := s.GetUserByID(target.ID)
	if err != nil {
		log.Printf("GET /api/users/%d: get target: %v", targetID, err)
		restapi.WriteJsonApiError(w, http.StatusInternalServerError, "server_error", "Internal Server Error", "")
		return
	} else if user == nil {
		restapi.WriteJsonApiError(w, http.StatusNotFound, "user_not_found", "Not Found", "User not found.")
		return
	}

	view := s.buildUserView(user, actor, target)

	restapi.WriteJsonApiData(w, http.StatusOK, view)
}

// HandlePatchUser updates a user's profile (with RBAC).
// Ember Data lies and sends the entire record up, not just the change set.
// PATCH /api/users/:id
func (s *Service) HandlePatchUser(w http.ResponseWriter, r *http.Request) {
	targetID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || domains.ID(targetID) == domains.InvalidID {
		restapi.WriteJsonApiMalformedPathParameter(w, "user_id", "User ID", r.PathValue("id"))
		return
	}

	actor, err := s.authSvc.GetActor(r)
	if err != nil || !actor.IsValid() {
		restapi.WriteJsonApiError(w, http.StatusUnauthorized, "not_authenticated", "Unauthorized", "Sign in to access this resource.")
		return
	}
	target, err := s.authSvc.GetActorById(domains.ID(targetID))
	if err != nil || !target.IsValid() {
		// don't leak valid or invalid id
		restapi.WriteJsonApiError(w, http.StatusForbidden, "forbidden", "Forbidden", "You do not have permission to view this user.")
		return
	}

	if !s.authSvc.CanEditTarget(actor, target) {
		restapi.WriteJsonApiError(w, http.StatusForbidden, "forbidden", "Forbidden", "You do not have permission to update this user.")
		return
	}

	var p UserPatchRequest
	if err := jsonapi.UnmarshalPayload(r.Body, &p); err != nil {
		log.Printf("PATCH /api/users/%d/role: update: %v", targetID, err)
		restapi.WriteJsonApiError(w, http.StatusBadRequest, "bad_request", "Invalid Request Body", err.Error())
		return
	}
	log.Printf("%s %s: unmarshal %+v\n", r.Method, r.URL.Path, p)

	user, err := s.GetUserByID(target.ID)
	if err != nil {
		log.Printf("GET /api/users/%d: get target: %v", targetID, err)
		restapi.WriteJsonApiError(w, http.StatusInternalServerError, "server_error", "Internal Server Error", "")
		return
	}
	if user == nil {
		log.Printf("GET /api/users/%d: get target: not found", targetID)
		restapi.WriteJsonApiError(w, http.StatusInternalServerError, "server_error", "Internal Server Error", "")
		return
	}

	// Calculate the change set and capture any validation errors
	var validationErrors []*jsonapi.ErrorObject
	changedUsername := false
	if p.Username != user.Username {
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
		}
	}
	changedTimezone := false
	var newTimezone *time.Location
	if p.Timezone != user.Locale.Timezone.Location.String() {
		if newTimezone, err = time.LoadLocation(p.Timezone); err != nil {
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
		}
	}
	if validationErrors != nil {
		restapi.WriteJsonApiErrorObjects(w, http.StatusUnprocessableEntity, validationErrors...)
		return
	}
	if !(changedUsername || changedEmail || changedTimezone) {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Check if trying to edit username
	if changedUsername && !s.authSvc.CanEditTargetUsername(actor, target) {
		restapi.WriteJsonApiError(w, http.StatusForbidden, "forbidden", "Forbidden", "You do not have permission to update user names.")
		return
	}

	// Update user
	updatedUser, err := s.UpsertUser(user.Handle, p.Email, p.Username, newTimezone)
	log.Printf("PATCH /api/users/%d: update user: %v", targetID, err)
	if err != nil {
		log.Printf("PATCH /api/users/%d: update: %v", targetID, err)
		restapi.WriteJsonApiError(w, http.StatusInternalServerError, "server_error", "Internal Server Error", "")
		return
	}

	// Return updated user
	log.Printf("PATCH /api/users/%d: get user by id: %v", targetID, err)
	view := s.buildUserView(updatedUser, actor, target)

	restapi.WriteJsonApiData(w, http.StatusOK, view)
}

/*
| Scenario                                                                                 | Meaning                                | Status                       |
| ---------------------------------------------------------------------------------------- | -------------------------------------- | ---------------------------- |
| **Unauthenticated** (no valid session/JWT)                                               | “You must sign in first.”              | **401 Unauthorized**         |
| **Authenticated but not allowed** (role too low, wrong resource, bad old password, etc.) | “You’re signed in, but can’t do this.” | **403 Forbidden**            |
| **Wrong resource identifier** (the IDs don’t match, or you’re editing someone else)      | “You can’t edit this user.”            | **403 Forbidden**            |
| **User doesn’t exist / soft-deleted**                                                    | “Resource missing.”                    | **404 Not Found**            |
| **Validation or business rule fails** (password too short, missing field)                | Client error but not authorization     | **422 Unprocessable Entity** |
*/

// HandlePatchPassword updates the target's credentials
// PATCH /api/users/:id/password
func (s *Service) HandlePatchPassword(w http.ResponseWriter, r *http.Request) {
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
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		restapi.WriteJsonApiError(w, http.StatusBadRequest, "bad_request", "Invalid Request Body", "")
		return
	}

	actor, err := s.authSvc.GetActor(r)
	if err != nil || !actor.IsValid() {
		restapi.WriteJsonApiError(w, http.StatusUnauthorized, "not_authenticated", "Unauthorized", "Sign in to access this resource.")
		return
	}
	target, err := s.authSvc.GetActorById(domains.ID(targetID))
	if err != nil || !target.IsValid() {
		restapi.WriteJsonApiError(w, http.StatusForbidden, "forbidden", "Forbidden", "You do not have permission to update user passwords.")
		return
	}

	err = s.authSvc.UpdateCredentials(actor, target, req.CurrentPassword, req.NewPassword)
	if err != nil {
		restapi.WriteJsonApiError(w, http.StatusForbidden, "forbidden", "Forbidden", "You do not have permission to update user passwords.")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HandlePostResetPassword resets a user's password (admin only)
// POST /api/users/:id/reset-password
func (s *Service) HandlePostResetPassword(w http.ResponseWriter, r *http.Request) {
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

	actor, err := s.authSvc.GetActor(r)
	if err != nil || !actor.IsValid() {
		restapi.WriteJsonApiError(w, http.StatusUnauthorized, "not_authenticated", "Unauthorized", "Sign in to access this resource.")
		return
	}
	target, err := s.authSvc.GetActorById(domains.ID(targetID))
	if err != nil || !target.IsValid() {
		restapi.WriteJsonApiError(w, http.StatusForbidden, "forbidden", "Forbidden", "You do not have permission to reset the user's password.")
		return
	}

	// Generate temporary password reset link
	magicLink := phrases.Generate(6)
	err = s.authSvc.UpdateCredentials(actor, target, "", magicLink)
	if err != nil {
		log.Printf("POST /api/users/%d/reset-password: %d: %d: %v", actor.ID, target.ID, err)
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

// HandlePostUser creates a new user (admin only)
// POST /api/users
func (s *Service) HandlePostUser(w http.ResponseWriter, r *http.Request) {
	actor, err := s.authSvc.GetActor(r)
	if err != nil || !actor.IsValid() {
		restapi.WriteJsonApiError(w, http.StatusUnauthorized, "not_authenticated", "Unauthorized", "Sign in to access this resource.")
		return
	}

	if !s.authSvc.CanCreateTarget(actor) {
		restapi.WriteJsonApiError(w, http.StatusForbidden, "forbidden", "Forbidden", "You do not have permission to create users.")
		return
	}

	var p UserCreateRequest
	if err := jsonapi.UnmarshalPayload(r.Body, &p); err != nil {
		log.Printf("POST /api/users: update: %v", err)
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
	user, err := s.UpsertUser(p.Handle, p.Email, p.Username, loc)
	if err != nil {
		log.Printf("POST /api/users: create: %v", err)
		restapi.WriteJsonApiError(w, http.StatusInternalServerError, "server_error", "Internal Server Error", "")
		return
	}
	target, err := s.authSvc.GetActorById(user.ID)
	if err != nil {
		log.Printf("POST /api/users: %d: %d: create: %v", actor.ID, user.ID, err)
		restapi.WriteJsonApiError(w, http.StatusInternalServerError, "server_error", "Internal Server Error", "")
		return
	}
	err = s.authSvc.UpdateCredentials(actor, target, "", password)
	if err != nil {
		log.Printf("POST /api/users: %d: %d: create: %v", actor.ID, user.ID, err)
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
			s.authSvc.RemoveRole(user.ID, "user")
			s.authSvc.AssignRole(user.ID, "guest")
		}

		// Assign additional roles
		for _, role := range p.Roles {
			if role != "active" && role != "user" { // active and user already assigned
				err = s.authSvc.AssignRole(user.ID, role)
				if err != nil {
					log.Printf("POST /api/users: assign role %q: %v", role, err)
				}
			}
		}
	}

	view := s.buildUserView(user, actor, target)

	// set a Location key in the header to let the client know how to find the new user
	w.Header().Set("Location", restapi.AbsURL(r, "/api/users/"+view.ID))

	restapi.WriteJsonApiData(w, http.StatusCreated, view)
}

// HandlePatchUserRelationshipRoles is a route to make JSON:API happier
// PATCH /api/users/:id/relationships/roles
//
//	{
//	 "data": [ { "type": "roles", "id": "admin" }, ... ]
//	}
func (s *Service) HandlePatchUserRelationshipRoles(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

// HandlePatchUserRole updates a user's role (admin only)
// PATCH /api/users/:id/role
func (s *Service) HandlePatchUserRole(w http.ResponseWriter, r *http.Request) {
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

	actor, err := s.authSvc.GetActor(r)
	if err != nil || !actor.IsValid() {
		restapi.WriteJsonApiError(w, http.StatusUnauthorized, "not_authenticated", "Unauthorized", "Sign in to access this resource.")
		return
	}
	target, err := s.authSvc.GetActorById(domains.ID(targetID))
	if err != nil || !target.IsValid() {
		restapi.WriteJsonApiError(w, http.StatusForbidden, "forbidden", "Forbidden", "You do not have permission to manage this user.")
		return
	}

	if !s.authSvc.CanManageTargetRoles(actor, target) {
		restapi.WriteJsonApiError(w, http.StatusForbidden, "forbidden", "Forbidden", "You do not have permission to manage this user.")
		return
	}

	// Add roles
	for _, roleID := range req.Add {
		err = s.authSvc.AssignRole(target.ID, roleID)
		if err != nil {
			log.Printf("PATCH /api/users/%d/role: add %q: %v", target.ID, roleID, err)
			restapi.WriteJsonApiError(w, http.StatusInternalServerError, "server_error", "Internal Server Error", "")
			return
		}
	}

	// Remove roles
	for _, roleID := range req.Remove {
		err = s.authSvc.RemoveRole(target.ID, roleID)
		if err != nil {
			log.Printf("PATCH /api/users/%d/role: remove %q: %v", target.ID, roleID, err)
			restapi.WriteJsonApiError(w, http.StatusInternalServerError, "server_error", "Internal Server Error", "")
			return
		}
	}

	// Return updated user
	user, err := s.GetUserByID(target.ID)
	if err != nil {
		restapi.WriteJsonApiError(w, http.StatusInternalServerError, "server_error", "Internal Server Error", "")
		return
	}

	view := s.buildUserView(user, actor, target)

	restapi.WriteJsonApiData(w, http.StatusOK, view)
}

// buildUserView constructs a UserView with permissions based on actor's privileges
func (s *Service) buildUserView(user *domains.User_t, actor, target *domains.Actor) *UserView {
	aa := s.authSvc.BuildActorAuth(actor, target)
	return &UserView{
		ID:          fmt.Sprintf("%d", user.ID),
		Username:    user.Username,
		Email:       user.Email,
		Handle:      user.Handle,
		Timezone:    user.Locale.Timezone.Location.String(),
		Roles:       aa.Roles,
		Permissions: aa.Permissions,
		CreatedAt:   user.Created.UTC(),
		UpdatedAt:   user.Updated.UTC(),
	}
}
