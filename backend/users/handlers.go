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
	"github.com/playbymail/ottoapp/backend/stores/sqlite/sqlc"
)

// UserView is the JSON:API view for a user
type UserView struct {
	ID          string          `jsonapi:"primary,user"` // singular when sending a payload
	Username    string          `jsonapi:"attr,username"`
	Email       string          `jsonapi:"attr,email"`
	Timezone    string          `jsonapi:"attr,timezone"`
	Roles       []string        `jsonapi:"attr,roles,omitempty"`
	Permissions map[string]bool `jsonapi:"attr,permissions,omitempty"`
	CreatedAt   time.Time       `jsonapi:"attr,created-at,iso8601"`
	UpdatedAt   time.Time       `jsonapi:"attr,updated-at,iso8601"`
}

// HandleGetMe returns the current user's profile.
// For completeness; use 401 when no session, 404 if record missing.
func (s *Service) HandleGetMe(w http.ResponseWriter, r *http.Request) {
	actorID, err := s.authSvc.GetActor(r)
	if err != nil || actorID == domains.InvalidID {
		WriteJsonApiError(w, http.StatusUnauthorized, "not_authenticated", "Unauthorized", "Sign in to access this resource.")
		return
	}

	actor, err := s.GetUserByID(actorID)
	if err != nil || actor == nil {
		log.Printf("GET /api/users/me: user %d not found: %v", actorID, err)
		WriteJsonApiError(w, http.StatusNotFound, "user_not_found", "Not Found", "User not found.")
		return
	}

	view := s.buildUserView(actor, actorID)

	WriteJsonApiData(w, http.StatusOK, view)
}

// HandleGetMyProfile returns the current user's profile.
// GET /api/my/profile
func (s *Service) HandleGetMyProfile(w http.ResponseWriter, r *http.Request) {
	actorID, err := s.authSvc.GetActor(r)
	if err != nil || actorID == domains.InvalidID {
		WriteJsonApiError(w, http.StatusUnauthorized, "not_authenticated", "Unauthorized", "Sign in to access this resource.")
		return
	}

	target, err := s.GetUserByID(actorID)
	if err != nil || target == nil {
		log.Printf("GET /api/my/profile: user %d not found: %v", actorID, err)
		WriteJsonApiError(w, http.StatusNotFound, "user_not_found", "Not Found", "User not found.")
		return
	}

	view := s.buildUserView(target, actorID)

	WriteJsonApiData(w, http.StatusOK, view)
}

// HandleGetUsers returns a list of all non-admin, non-sysop users (admin only)
// GET /api/users
func (s *Service) HandleGetUsers(w http.ResponseWriter, r *http.Request) {
	actorID, err := s.authSvc.GetActor(r)
	if err != nil || actorID == domains.InvalidID {
		WriteJsonApiError(w, http.StatusUnauthorized, "not_authenticated", "Unauthorized", "Sign in to access this resource.")
		return
	}

	canList, err := s.authSvc.CanListUsers(actorID)
	if err != nil {
		log.Printf("GET /api/users: check access: %v", err)
		WriteJsonApiError(w, http.StatusInternalServerError, "server_error", "Internal Server Error", "")
		return
	}
	if !canList {
		WriteJsonApiError(w, http.StatusForbidden, "forbidden", "Forbidden", "You do not have permission to view user lists.")
		return
	}

	// Get all users from database (we'll filter them)
	// TODO: Filtering in Go after GetAllUsers is going to hurt later.
	// Prefer a store method that already enforces visibility/roles.
	// When you add pagination: include links.self/next/prev and meta.total/page/pageSize/returned.
	rows, err := s.db.Queries().GetAllUsers(s.db.Context())
	if err != nil {
		log.Printf("GET /api/users: query: %v", err)
		WriteJsonApiError(w, http.StatusInternalServerError, "server_error", "Internal Server Error", "")
		return
	}

	var view []*UserView
	for _, row := range rows {
		targetID := domains.ID(row.UserID)

		// Check authorization
		canView, err := s.authSvc.CanViewUser(actorID, targetID)
		if err != nil {
			log.Printf("GET /api/users/%d: check access: %v", targetID, err)
			WriteJsonApiError(w, http.StatusInternalServerError, "server_error", "Internal Server Error", "")
			return
		}
		if !canView { // skip targets the actor is not allowed to view
			continue
		}

		target, err := s.GetUserByID(targetID)
		if err != nil {
			log.Printf("GET /api/users/%d: get target: %v", targetID, err)
			WriteJsonApiError(w, http.StatusInternalServerError, "server_error", "Internal Server Error", "")
			return
		}
		if target == nil {
			log.Printf("GET /api/users/%d: get target: not found", targetID)
			WriteJsonApiError(w, http.StatusInternalServerError, "server_error", "Internal Server Error", "")
			return
		}

		view = append(view, s.buildUserView(target, actorID))
	}

	WriteJsonApiData(w, http.StatusOK, view)
}

// HandleGetUsersWithPagination returns a user list with optional pagination support
// GET /api/users?page[number]=1&page[size]=25
func (s *Service) HandleGetUsersWithPagination(w http.ResponseWriter, r *http.Request) {
	actorID, err := s.authSvc.GetActor(r)
	if err != nil || actorID == domains.InvalidID {
		WriteJsonApiError(w, http.StatusUnauthorized, "not_authenticated", "Unauthorized", "Sign in to access this resource.")
		return
	}

	// Only admins/sysops may list users
	canList, err := s.authSvc.CanListUsers(actorID)
	if err != nil {
		log.Printf("GET /api/users: check access: %v", err)
		WriteJsonApiError(w, http.StatusInternalServerError, "server_error", "Internal Server Error", "")
		return
	}
	if !canList {
		WriteJsonApiError(w, http.StatusForbidden, "forbidden", "Forbidden", "You do not have permission to view user lists.")
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
	rows, err := s.db.Queries().ListUsersVisibleToActor(s.db.Context(), sqlc.ListUsersVisibleToActorParams{
		ActorID:  actorID,
		PageSize: pageSize,
		PageNum:  (pageNum - 1) * pageSize,
	})
	if err != nil {
		log.Printf("GET /api/users?page[number]=%d&page[size]=%d: query: %v", pageNum, pageSize, err)
		WriteJsonApiError(w, http.StatusInternalServerError, "server_error", "Internal Server Error", "")
		return
	}

	// Build views
	views := make([]*UserView, 0, len(rows))
	for _, row := range rows {
		views = append(views, &UserView{
			ID:       fmt.Sprintf("%d", row.UserID),
			Username: row.Username,
			Email:    row.Email,
			Timezone: row.Timezone,
			//Roles:       row.Roles,       // if you prejoin/aggregate; otherwise compute once for the list
			//Permissions: row.Permissions, // same note as above
			CreatedAt: time.Unix(row.CreatedAt, 0).UTC(),
			UpdatedAt: time.Unix(row.UpdatedAt, 0).UTC(),
		})
	}

	// Marshal JSON:API data array to a buffer
	var dataBuf bytes.Buffer
	if err := jsonapi.MarshalPayload(&dataBuf, views); err != nil {
		WriteJsonApiError(w, http.StatusInternalServerError, "encode_failed", "Encoding error", err.Error())
		return
	}

	// Unwrap to add links/meta at top-level
	var payload map[string]any
	if err := json.Unmarshal(dataBuf.Bytes(), &payload); err != nil {
		WriteJsonApiError(w, http.StatusInternalServerError, "encode_failed", "Encoding error", err.Error())
		return
	}

	self := paginateURL(r, pageNum, pageSize)
	var next, prev string
	if (pageNum-1)*pageSize+len(views) < int(total) {
		next = paginateURL(r, pageNum+1, pageSize)
	}
	if pageNum > 1 {
		prev = paginateURL(r, pageNum-1, pageSize)
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

	out, _ := json.Marshal(payload)
	WriteJsonApiResponse(w, http.StatusOK, out)
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

func paginateURL(r *http.Request, page, size int) string {
	u := *r.URL // copy
	q := u.Query()
	q.Set("page[number]", strconv.Itoa(page))
	q.Set("page[size]", strconv.Itoa(size))
	u.RawQuery = q.Encode()
	return AbsURL(r, u.Path+"?"+u.RawQuery)
}

// HandleGetUser returns a specific user's profile (with RBAC).
func (s *Service) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	// Parse target user ID from path
	targetIDInt, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		WriteJsonApiErrorObjects(w, http.StatusBadRequest, &jsonapi.ErrorObject{
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
	targetID := domains.ID(targetIDInt)

	actorID, err := s.authSvc.GetActor(r)
	if err != nil || actorID == domains.InvalidID {
		WriteJsonApiError(w, http.StatusUnauthorized, "not_authenticated", "Unauthorized", "Sign in to access this resource.")
		return
	}

	// Check authorization
	canView, err := s.authSvc.CanViewUser(actorID, targetID)
	if err != nil {
		log.Printf("GET /api/users/%d: check access: %v", targetID, err)
		WriteJsonApiError(w, http.StatusInternalServerError, "server_error", "Internal Server Error", "")
		return
	}
	if !canView {
		WriteJsonApiError(w, http.StatusForbidden, "forbidden", "Forbidden", "You do not have permission to view this user.")
		return
	}

	target, err := s.GetUserByID(targetID)
	if err != nil {
		log.Printf("GET /api/users/%d: get target: %v", targetID, err)
		WriteJsonApiError(w, http.StatusInternalServerError, "server_error", "Internal Server Error", "")
		return
	} else if target == nil {
		WriteJsonApiError(w, http.StatusNotFound, "user_not_found", "Not Found", "User not found.")
		return
	}

	view := s.buildUserView(target, actorID)

	WriteJsonApiData(w, http.StatusOK, view)
}

// HandlePatchUser updates a user's profile (with RBAC).
// Ember Data lies and sends the entire record up, not just the change set.
// PATCH /api/users/:id
func (s *Service) HandlePatchUser(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s: entered\n", r.Method, r.URL.Path)
	// Parse target user ID from path
	targetIDInt, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	targetID := domains.ID(targetIDInt)
	if err != nil || targetID == domains.InvalidID {
		WriteJsonApiErrorObjects(w, http.StatusBadRequest, &jsonapi.ErrorObject{
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
	log.Printf("%s %s: targetID %d\n", r.Method, r.URL.Path, targetID)

	actorID, err := s.authSvc.GetActor(r)
	log.Printf("%s %s: actorID %d, err %v\n", r.Method, r.URL.Path, actorID, err)
	if err != nil || actorID == domains.InvalidID {
		WriteJsonApiError(w, http.StatusUnauthorized, "not_authenticated", "Unauthorized", "Sign in to access this resource.")
		return
	}

	// Check authorization
	canEdit, err := s.authSvc.CanEditUser(actorID, targetID)
	log.Printf("%s %s: canEdit %v, err %v\n", r.Method, r.URL.Path, canEdit, err)
	if err != nil {
		log.Printf("PATCH /api/users/%d: check access: %v", targetID, err)
		WriteJsonApiError(w, http.StatusInternalServerError, "server_error", "Internal Server Error", "")
		return
	}
	if !canEdit {
		log.Printf("%s %s: canEdit false, returning 403\n", r.Method, r.URL.Path)
		WriteJsonApiError(w, http.StatusForbidden, "forbidden", "Forbidden", "You do not have permission to update this user.")
		return
	}

	type UserPatchPayload struct {
		ID       string `jsonapi:"primary,users"` // plural when receiving from Ember Data
		Username string `jsonapi:"attr,username,omitempty"`
		Email    string `jsonapi:"attr,email,omitempty"`
		Timezone string `jsonapi:"attr,timezone,omitempty"`
	}
	var p UserPatchPayload
	if err := jsonapi.UnmarshalPayload(r.Body, &p); err != nil {
		log.Printf("PATCH /api/users/%d/role: update: %v", targetID, err)
		WriteJsonApiError(w, http.StatusBadRequest, "bad_request", "Invalid Request Body", err.Error())
		return
	}
	log.Printf("%s %s: unmarshal %+v\n", r.Method, r.URL.Path, p)

	target, err := s.GetUserByID(targetID)
	if err != nil {
		log.Printf("GET /api/users/%d: get target: %v", targetID, err)
		WriteJsonApiError(w, http.StatusInternalServerError, "server_error", "Internal Server Error", "")
		return
	}
	if target == nil {
		log.Printf("GET /api/users/%d: get target: not found", targetID)
		WriteJsonApiError(w, http.StatusInternalServerError, "server_error", "Internal Server Error", "")
		return
	}

	// Calculate the change set and capture any validation errors
	var validationErrors []*jsonapi.ErrorObject
	changedUsername := false
	if p.Username != target.Username {
		if !ValidateUsername(p.Username) {
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
	if p.Email != target.Email {
		if !ValidateEmail(p.Email) {
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
	if p.Timezone != target.Locale.Timezone.Location.String() {
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
		WriteJsonApiErrorObjects(w, http.StatusUnprocessableEntity, validationErrors...)
		return
	}
	if !(changedUsername || changedEmail || changedTimezone) {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Check if trying to edit username
	if changedUsername {
		canEditUsername, err := s.authSvc.CanEditUsername(actorID, targetID)
		if err != nil {
			log.Printf("PATCH /api/users/%d: check username edit: %v", targetID, err)
			WriteJsonApiError(w, http.StatusInternalServerError, "server_error", "Internal Server Error", "")
			return
		}
		if !canEditUsername {
			WriteJsonApiError(w, http.StatusForbidden, "forbidden", "Forbidden", "You do not have permission to update user names.")
			return
		}
	}

	// Update user
	var newUsername, newEmail *string
	if changedUsername {
		newUsername = &p.Username
	}
	if changedEmail {
		newEmail = &p.Email
	}
	err = s.UpdateUser(targetID, newUsername, newEmail, newTimezone)
	log.Printf("PATCH /api/users/%d: update user: %v", targetID, err)
	if err != nil {
		log.Printf("PATCH /api/users/%d: update: %v", targetID, err)
		WriteJsonApiError(w, http.StatusInternalServerError, "server_error", "Internal Server Error", "")
		return
	}

	// Return updated user
	user, err := s.GetUserByID(targetID)
	log.Printf("PATCH /api/users/%d: get user by id: %v", targetID, err)
	if err != nil {
		WriteJsonApiError(w, http.StatusInternalServerError, "server_error", "Internal Server Error", "")
		return
	}

	view := s.buildUserView(user, actorID)

	WriteJsonApiData(w, http.StatusOK, view)
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

// HandlePatchPassword updates the current user's password
// PATCH /api/users/:id/password
func (s *Service) HandlePatchPassword(w http.ResponseWriter, r *http.Request) {
	// Parse target user ID from path
	targetIDInt, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	targetID := domains.ID(targetIDInt)
	if err != nil || targetID == domains.InvalidID {
		WriteJsonApiErrorObjects(w, http.StatusBadRequest, &jsonapi.ErrorObject{
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
		WriteJsonApiError(w, http.StatusBadRequest, "bad_request", "Invalid Request Body", "")
		return
	}

	actorID, err := s.authSvc.GetActor(r)
	if err != nil || actorID == domains.InvalidID {
		WriteJsonApiError(w, http.StatusUnauthorized, "not_authenticated", "Unauthorized", "Sign in to access this resource.")
		return
	}
	if actorID != targetID {
		WriteJsonApiError(w, http.StatusForbidden, "forbidden", "Forbidden", "You do not have permission to update user passwords.")
		return
	}

	canChange, _ := s.authSvc.CanChangeOwnPassword(actorID, targetID)
	if !canChange {
		WriteJsonApiError(w, http.StatusForbidden, "forbidden", "Forbidden", "You do not have permission to update user passwords.")
		return
	}

	// Get email for authentication
	target, err := s.GetUserByID(targetID)
	if err != nil {
		WriteJsonApiError(w, http.StatusNotFound, "user_not_found", "Not Found", "User not found.")
		return
	}

	// Verify current password
	emailID, err := s.authSvc.AuthenticateWithEmailSecret(target.Email, req.CurrentPassword)
	if err != nil {
		WriteJsonApiError(w, http.StatusForbidden,
			"invalid_current_password",
			"Forbidden",
			"You do not have permission to update the password with the provided credentials.")
		return
	} else if emailID != targetID {
		log.Printf("PUT /api/users/%d/password: update: insanity %q -> %d", targetID, target.Email, emailID)
		WriteJsonApiError(w, http.StatusInternalServerError, "server_error", "Internal Server Error", "")
	}

	// Update password
	err = s.authSvc.UpdateUserSecret(targetID, req.NewPassword)
	if err != nil {
		log.Printf("PUT /api/users/%d/password: update: %v", targetID, err)
		WriteJsonApiError(w, http.StatusInternalServerError, "server_error", "Internal Server Error", "")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HandlePostResetPassword resets a user's password (admin only)
// POST /api/users/:id/reset-password
func (s *Service) HandlePostResetPassword(w http.ResponseWriter, r *http.Request) {
	// Parse target user ID from path
	targetIDInt, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	targetID := domains.ID(targetIDInt)
	if err != nil || targetID == domains.InvalidID {
		WriteJsonApiErrorObjects(w, http.StatusBadRequest, &jsonapi.ErrorObject{
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

	actorID, err := s.authSvc.GetActor(r)
	if err != nil || actorID == domains.InvalidID {
		WriteJsonApiError(w, http.StatusUnauthorized, "not_authenticated", "Unauthorized", "Sign in to access this resource.")
		return
	}

	// Check authorization
	canReset, err := s.authSvc.CanResetPassword(actorID, targetID)
	if err != nil {
		log.Printf("POST /api/users/%d/reset-password: check access: %v", targetID, err)
		WriteJsonApiError(w, http.StatusInternalServerError, "server_error", "Internal Server Error", "")
		return
	}
	if !canReset {
		WriteJsonApiError(w, http.StatusForbidden, "forbidden", "Forbidden", "You do not have permission to reset the user's password.")
		return
	}

	// Generate temporary password reset link
	magicLink := phrases.Generate(6)

	// Update password
	err = s.authSvc.UpdateUserSecret(targetID, magicLink)
	if err != nil {
		log.Printf("POST /api/users/%d/reset-password: update: %v", targetID, err)
		WriteJsonApiError(w, http.StatusInternalServerError, "server_error", "Internal Server Error", "")
		return
	}

	WriteJsonApiData(w, http.StatusOK, &struct {
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
	actorID, err := s.authSvc.GetActor(r)
	if err != nil || actorID == domains.InvalidID {
		WriteJsonApiError(w, http.StatusUnauthorized, "not_authenticated", "Unauthorized", "Sign in to access this resource.")
		return
	}

	canCreate, err := s.authSvc.CanCreateUser(actorID)
	if err != nil {
		log.Printf("POST /api/users: check access: %v", err)
		WriteJsonApiError(w, http.StatusInternalServerError, "server_error", "Internal Server Error", "")
		return
	}
	if !canCreate {
		WriteJsonApiError(w, http.StatusForbidden, "forbidden", "Forbidden", "You do not have permission to create users.")
		return
	}

	type UserCreatePayload struct {
		ID       string   `jsonapi:"primary,users"` // plural when receiving from Ember Data
		Username string   `jsonapi:"attr,username,omitempty"`
		Email    string   `jsonapi:"attr,email,omitempty"`
		Password string   `jsonapi:"attr,password,omitempty"`
		Timezone string   `jsonapi:"attr,timezone,omitempty"`
		Roles    []string `json:"attr,roles,omitempty"`
	}
	var p UserCreatePayload
	if err := jsonapi.UnmarshalPayload(r.Body, &p); err != nil {
		log.Printf("POST /api/users: update: %v", err)
		WriteJsonApiError(w, http.StatusBadRequest, "bad_request", "Invalid Request Body", err.Error())
		return
	}
	var validationErrors []*jsonapi.ErrorObject
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
	} else if !ValidateUsername(p.Username) {
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
	} else if !ValidateEmail(p.Email) {
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
		WriteJsonApiErrorObjects(w, http.StatusUnprocessableEntity, validationErrors...)
		return
	}

	// Generate password if not provided
	password := p.Password
	if password == "" {
		password = phrases.Generate(6)
	}

	// Create user (this will assign "active" and "user" roles by default)
	user, err := s.CreateUser(p.Username, p.Email, password, loc)
	if err != nil {
		log.Printf("POST /api/users: create: %v", err)
		WriteJsonApiError(w, http.StatusInternalServerError, "server_error", "Internal Server Error", "")
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

	view := s.buildUserView(user, actorID)

	// set a Location key in the header to let the client know how to find the new user
	w.Header().Set("Location", AbsURL(r, "/api/users/"+view.ID))

	WriteJsonApiData(w, http.StatusCreated, view)
}

// Consider adding a route to make JSON:API happier
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
	targetIDInt, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	targetID := domains.ID(targetIDInt)
	if err != nil || targetID == domains.InvalidID {
		WriteJsonApiErrorObjects(w, http.StatusBadRequest, &jsonapi.ErrorObject{
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
		WriteJsonApiError(w, http.StatusBadRequest, "bad_request", "Invalid Request Body", "")
		return
	}

	actorID, err := s.authSvc.GetActor(r)
	if err != nil || actorID == domains.InvalidID {
		WriteJsonApiError(w, http.StatusUnauthorized, "not_authenticated", "Unauthorized", "Sign in to access this resource.")
		return
	}

	canManage, err := s.authSvc.CanManageRoles(actorID, targetID)
	if err != nil {
		log.Printf("PATCH /api/users/%d/role: check access: %v", targetID, err)
		WriteJsonApiError(w, http.StatusInternalServerError, "server_error", "Internal Server Error", "")
		return
	}
	if !canManage {
		WriteJsonApiError(w, http.StatusForbidden, "forbidden", "Forbidden", "You do not have permission to manage this user.")
		return
	}

	// Add roles
	for _, roleID := range req.Add {
		err = s.authSvc.AssignRole(targetID, roleID)
		if err != nil {
			log.Printf("PATCH /api/users/%d/role: add %q: %v", targetID, roleID, err)
			WriteJsonApiError(w, http.StatusInternalServerError, "server_error", "Internal Server Error", "")
			return
		}
	}

	// Remove roles
	for _, roleID := range req.Remove {
		err = s.authSvc.RemoveRole(targetID, roleID)
		if err != nil {
			log.Printf("PATCH /api/users/%d/role: remove %q: %v", targetID, roleID, err)
			WriteJsonApiError(w, http.StatusInternalServerError, "server_error", "Internal Server Error", "")
			return
		}
	}

	// Return updated user
	user, err := s.GetUserByID(targetID)
	if err != nil {
		WriteJsonApiError(w, http.StatusInternalServerError, "server_error", "Internal Server Error", "")
		return
	}

	view := s.buildUserView(user, actorID)

	WriteJsonApiData(w, http.StatusOK, view)
}

// buildUserView constructs a UserView with permissions based on actor's privileges
func (s *Service) buildUserView(user *domains.User_t, actorID domains.ID) *UserView {
	auth, _ := s.authSvc.BuildActorAuth(actorID, user.ID)
	return &UserView{
		ID:          fmt.Sprintf("%d", user.ID),
		Username:    user.Username,
		Email:       user.Email,
		Timezone:    user.Locale.Timezone.Location.String(),
		Roles:       auth.Roles,
		Permissions: auth.Permissions,
		CreatedAt:   user.Created.UTC(),
		UpdatedAt:   user.Updated.UTC(),
	}
}

// JSON:API helpers

func AbsURL(r *http.Request, path string) string {
	scheme := "https"
	if r.Header.Get("X-Forwarded-Proto") != "" {
		scheme = r.Header.Get("X-Forwarded-Proto")
	} else if r.TLS == nil {
		scheme = "http"
	}
	host := r.Header.Get("X-Forwarded-Host")
	if host == "" {
		host = r.Host
	}
	return scheme + "://" + host + path
}

func WriteJsonApiData(w http.ResponseWriter, status int, view any) {
	buf := &bytes.Buffer{}
	err := jsonapi.MarshalPayload(buf, view)
	if err != nil { // As an absolute last resort, write a minimal JSON:API error payload.
		buf = &bytes.Buffer{}
		status = http.StatusInternalServerError
		buf.WriteString(fmt.Sprintf(`{"errors":[{"status":"500","code":"encode_failed","title":"Encoding error","detail":%q}]}`, err.Error()))
	}
	WriteJsonApiResponse(w, status, buf.Bytes())
}

// WriteJsonApiError
//
//	Quick status-code guide (JSON:API)
//
//	  200/201: success with a data document
//
//	  204: success without a body (DELETE, or rare no-content updates)
//
//	  401/403/404/422/500: JSON:API error document ({ "errors": [ … ] })
//
// You should Map errors → statuses (classify the failure) and return the right status + error object(s):
//   - Bad payload / shape → 400 Bad Request
//   - AuthN/AuthZ → 401/403
//   - Target missing → 404 Not Found
//   - Validation failed (e.g., bad timezone, username too short) → 422 Unprocessable Entity
//   - Uniqueness/foreign-key constraint (e.g., email already taken) → usually 409 Conflict (some teams prefer 422; pick one and be consistent)
//   - Optimistic concurrency (ETag/If-Match mismatch) → 412 Precondition Failed
//   - DB unavailable / timeout / deadlock → 503 Service Unavailable (optionally Retry-After)
//   - Unknown / internal → 500 Internal Server Error
func WriteJsonApiError(w http.ResponseWriter, status int, code, title, detail string) {
	WriteJsonApiErrorObjects(w, status, &jsonapi.ErrorObject{
		Status: strconv.Itoa(status),
		Code:   code,
		Title:  title,
		Detail: detail,
	})
}

func WriteJsonApiErrorObjects(w http.ResponseWriter, status int, errs ...*jsonapi.ErrorObject) {
	var list []*jsonapi.ErrorObject
	for _, err := range errs {
		list = append(list, err)
	}
	buf := &bytes.Buffer{}
	err := jsonapi.MarshalErrors(buf, list)
	if err != nil { // As an absolute last resort, write a minimal JSON:API error payload.
		buf = &bytes.Buffer{}
		status = http.StatusInternalServerError
		buf.WriteString(fmt.Sprintf(`{"errors":[{"status":"500","code":"encode_failed","title":"Encoding error","detail":%q}]}`, err.Error()))
	}
	WriteJsonApiResponse(w, status, buf.Bytes())
}

func WriteJsonApiResponse(w http.ResponseWriter, status int, buf []byte) {
	w.Header().Set("Content-Type", jsonapi.MediaType)
	w.WriteHeader(status)
	_, _ = w.Write(buf)
}
