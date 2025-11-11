// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package users

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/mdhender/phrases/v2"
	"github.com/playbymail/ottoapp/backend/domains"
)

// UserResponse represents a user in API responses
type UserResponse struct {
	ID          int64           `json:"id"`
	Username    string          `json:"username"`
	Email       string          `json:"email"`
	Timezone    string          `json:"timezone"`
	Roles       []string        `json:"roles"`
	Permissions map[string]bool `json:"permissions"`
	Created     string          `json:"created"`
	Updated     string          `json:"updated"`
}

// HandleGetMe returns the current user's profile
// GET /api/users/me
func (s *Service) HandleGetMe(w http.ResponseWriter, r *http.Request) {
	// Get current user ID from session (stored by session middleware)
	userID, err := s.sessionsSvc.GetCurrentUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := s.GetUserByID(userID)
	if err != nil {
		log.Printf("GET /api/users/me: user %d: %v", userID, err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	roles, err := s.authSvc.GetUserRoles(userID)
	if err != nil {
		log.Printf("GET /api/users/me: roles: user %d: %v", userID, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := s.buildUserResponse(user, roles, userID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleGetUsers returns a list of all non-admin, non-sysop users (admin only)
// GET /api/users
func (s *Service) HandleGetUsers(w http.ResponseWriter, r *http.Request) {
	// Get current user ID from session
	actorID, ok := r.Context().Value("userID").(domains.ID)
	if !ok || actorID == domains.InvalidID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check if actor is admin
	isAdmin, err := s.authSvc.IsAdmin(actorID)
	if err != nil {
		log.Printf("GET /api/users: check admin: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !isAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Get all users from database (we'll filter them)
	rows, err := s.db.Queries().GetAllUsers(s.db.Context())
	if err != nil {
		log.Printf("GET /api/users: query: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var users []UserResponse
	for _, row := range rows {
		userID := domains.ID(row.UserID)

		// Get user roles
		roles, err := s.authSvc.GetUserRoles(userID)
		if err != nil {
			log.Printf("GET /api/users: roles: user %d: %v", userID, err)
			continue
		}

		// Skip sysop users
		if roles[domains.Role("sysop")] {
			continue
		}

		// Skip admin users
		if roles[domains.Role("admin")] {
			continue
		}

		loc, err := s.tzSvc.Location(row.Timezone)
		if err != nil {
			log.Printf("GET /api/users: timezone: user %d: %v", userID, err)
			continue
		}

		user := &domains.User_t{
			ID:       userID,
			Username: row.Username,
			Email:    row.Email,
			Locale: domains.UserLocale_t{
				DateFormat: "2006-01-02",
				Timezone: domains.UserTimezone_t{
					Location: loc,
				},
			},
			Created: time.Unix(row.CreatedAt, 0).UTC(),
			Updated: time.Unix(row.UpdatedAt, 0).UTC(),
		}

		users = append(users, s.buildUserResponse(user, roles, actorID))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// HandleGetUser returns a specific user's profile (with RBAC)
// GET /api/users/:id
func (s *Service) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	// Get current user ID from session
	actorID, ok := r.Context().Value("userID").(domains.ID)
	if !ok || actorID == domains.InvalidID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse target user ID from path
	idStr := r.PathValue("id")
	targetIDInt, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	targetID := domains.ID(targetIDInt)

	// Check authorization
	canView, err := s.authSvc.CanEditUser(actorID, targetID)
	if err != nil {
		log.Printf("GET /api/users/%d: check access: %v", targetID, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !canView {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	user, err := s.GetUserByID(targetID)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	roles, err := s.authSvc.GetUserRoles(targetID)
	if err != nil {
		log.Printf("GET /api/users/%d: roles: %v", targetID, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := s.buildUserResponse(user, roles, actorID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandlePatchUser updates a user's profile (with RBAC)
// PATCH /api/users/:id
func (s *Service) HandlePatchUser(w http.ResponseWriter, r *http.Request) {
	// Get current user ID from session
	actorID, ok := r.Context().Value("userID").(domains.ID)
	if !ok || actorID == domains.InvalidID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse target user ID from path
	idStr := r.PathValue("id")
	targetIDInt, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	targetID := domains.ID(targetIDInt)

	// Check authorization
	canEdit, err := s.authSvc.CanEditUser(actorID, targetID)
	if err != nil {
		log.Printf("PATCH /api/users/%d: check access: %v", targetID, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !canEdit {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Parse request body
	var req struct {
		Username *string `json:"username,omitempty"`
		Email    *string `json:"email,omitempty"`
		Timezone *string `json:"timezone,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check if trying to edit username
	if req.Username != nil {
		canEditUsername, err := s.authSvc.CanEditUsername(actorID, targetID)
		if err != nil {
			log.Printf("PATCH /api/users/%d: check username edit: %v", targetID, err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if !canEditUsername {
			http.Error(w, "Forbidden: cannot edit username", http.StatusForbidden)
			return
		}
	}

	// Parse timezone if provided
	var newTimezone *time.Location
	if req.Timezone != nil {
		loc, err := time.LoadLocation(*req.Timezone)
		if err != nil {
			http.Error(w, "Invalid timezone", http.StatusBadRequest)
			return
		}
		newTimezone = loc
	}

	// Update user
	err = s.UpdateUser(targetID, req.Username, req.Email, newTimezone)
	if err != nil {
		log.Printf("PATCH /api/users/%d: update: %v", targetID, err)
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	// Return updated user
	user, err := s.GetUserByID(targetID)
	if err != nil {
		http.Error(w, "Failed to fetch updated user", http.StatusInternalServerError)
		return
	}

	roles, err := s.authSvc.GetUserRoles(targetID)
	if err != nil {
		log.Printf("PATCH /api/users/%d: roles: %v", targetID, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := s.buildUserResponse(user, roles, actorID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandlePutPassword updates the current user's password
// PUT /api/users/:id/password
func (s *Service) HandlePutPassword(w http.ResponseWriter, r *http.Request) {
	// Get current user ID from session
	actorID, ok := r.Context().Value("userID").(domains.ID)
	if !ok || actorID == domains.InvalidID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse target user ID from path
	idStr := r.PathValue("id")
	targetIDInt, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	targetID := domains.ID(targetIDInt)

	// Users can only change their own password via this endpoint
	if actorID != targetID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Parse request body
	var req struct {
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get user email for authentication
	user, err := s.GetUserByID(actorID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Verify current password
	_, err = s.authSvc.AuthenticateWithEmailSecret(user.Email, req.CurrentPassword)
	if err != nil {
		http.Error(w, "Invalid current password", http.StatusUnauthorized)
		return
	}

	// Update password
	err = s.authSvc.UpdateUserSecret(actorID, req.NewPassword)
	if err != nil {
		log.Printf("PUT /api/users/%d/password: update: %v", actorID, err)
		http.Error(w, "Failed to update password", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "message": "Password updated"})
}

// HandlePostResetPassword resets a user's password (admin only)
// POST /api/users/:id/reset-password
func (s *Service) HandlePostResetPassword(w http.ResponseWriter, r *http.Request) {
	// Get current user ID from session
	actorID, ok := r.Context().Value("userID").(domains.ID)
	if !ok || actorID == domains.InvalidID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse target user ID from path
	idStr := r.PathValue("id")
	targetIDInt, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	targetID := domains.ID(targetIDInt)

	// Check authorization
	canReset, err := s.authSvc.CanResetPassword(actorID, targetID)
	if err != nil {
		log.Printf("POST /api/users/%d/reset-password: check access: %v", targetID, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !canReset {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Generate temporary password
	tempPassword := phrases.Generate(6)

	// Update password
	err = s.authSvc.UpdateUserSecret(targetID, tempPassword)
	if err != nil {
		log.Printf("POST /api/users/%d/reset-password: update: %v", targetID, err)
		http.Error(w, "Failed to reset password", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":       "ok",
		"message":      "Password reset",
		"tempPassword": tempPassword,
	})
}

// HandlePostUser creates a new user (admin only)
// POST /api/users
func (s *Service) HandlePostUser(w http.ResponseWriter, r *http.Request) {
	// Get current user ID from session
	actorID, ok := r.Context().Value("userID").(domains.ID)
	if !ok || actorID == domains.InvalidID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check if actor is admin
	isAdmin, err := s.authSvc.IsAdmin(actorID)
	if err != nil {
		log.Printf("POST /api/users: check admin: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !isAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Parse request body
	var req struct {
		Username string   `json:"username"`
		Email    string   `json:"email"`
		Password string   `json:"password,omitempty"`
		Timezone string   `json:"timezone"`
		Roles    []string `json:"roles,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Username == "" || req.Email == "" || req.Timezone == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Generate password if not provided
	password := req.Password
	if password == "" {
		password = phrases.Generate(6)
	}

	// Parse timezone
	loc, err := time.LoadLocation(req.Timezone)
	if err != nil {
		http.Error(w, "Invalid timezone", http.StatusBadRequest)
		return
	}

	// Create user (this will assign "active" and "user" roles by default)
	user, err := s.CreateUser(req.Username, req.Email, password, loc)
	if err != nil {
		log.Printf("POST /api/users: create: %v", err)
		http.Error(w, fmt.Sprintf("Failed to create user: %v", err), http.StatusInternalServerError)
		return
	}

	// If roles were specified and don't include "user" or "admin", assign "guest" instead
	if len(req.Roles) > 0 {
		hasUserOrAdmin := false
		for _, role := range req.Roles {
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
		for _, role := range req.Roles {
			if role != "active" && role != "user" { // active and user already assigned
				err = s.authSvc.AssignRole(user.ID, role)
				if err != nil {
					log.Printf("POST /api/users: assign role %q: %v", role, err)
				}
			}
		}
	}

	roles, err := s.authSvc.GetUserRoles(user.ID)
	if err != nil {
		log.Printf("POST /api/users: roles: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := s.buildUserResponse(user, roles, actorID)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandlePatchUserRole updates a user's role (admin only)
// PATCH /api/users/:id/role
func (s *Service) HandlePatchUserRole(w http.ResponseWriter, r *http.Request) {
	// Get current user ID from session
	actorID, ok := r.Context().Value("userID").(domains.ID)
	if !ok || actorID == domains.InvalidID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check if actor is admin
	isAdmin, err := s.authSvc.IsAdmin(actorID)
	if err != nil {
		log.Printf("PATCH /api/users/:id/role: check admin: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !isAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Parse target user ID from path
	idStr := r.PathValue("id")
	targetIDInt, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	targetID := domains.ID(targetIDInt)

	// Cannot modify sysop or admin roles
	isSysop, _ := s.authSvc.IsSysop(targetID)
	isTargetAdmin, _ := s.authSvc.IsAdmin(targetID)
	if isSysop || isTargetAdmin {
		http.Error(w, "Forbidden: cannot modify sysop or admin roles", http.StatusForbidden)
		return
	}

	// Parse request body
	var req struct {
		Add    []string `json:"add,omitempty"`
		Remove []string `json:"remove,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Add roles
	for _, roleID := range req.Add {
		err = s.authSvc.AssignRole(targetID, roleID)
		if err != nil {
			log.Printf("PATCH /api/users/%d/role: add %q: %v", targetID, roleID, err)
			http.Error(w, fmt.Sprintf("Failed to add role: %v", err), http.StatusInternalServerError)
			return
		}
	}

	// Remove roles
	for _, roleID := range req.Remove {
		err = s.authSvc.RemoveRole(targetID, roleID)
		if err != nil {
			log.Printf("PATCH /api/users/%d/role: remove %q: %v", targetID, roleID, err)
			http.Error(w, fmt.Sprintf("Failed to remove role: %v", err), http.StatusInternalServerError)
			return
		}
	}

	// Return updated user
	user, err := s.GetUserByID(targetID)
	if err != nil {
		http.Error(w, "Failed to fetch updated user", http.StatusInternalServerError)
		return
	}

	roles, err := s.authSvc.GetUserRoles(targetID)
	if err != nil {
		log.Printf("PATCH /api/users/%d/role: roles: %v", targetID, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := s.buildUserResponse(user, roles, actorID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// buildUserResponse constructs a UserResponse with permissions based on actor's privileges
func (s *Service) buildUserResponse(user *domains.User_t, roles domains.Roles, actorID domains.ID) UserResponse {
	// Convert roles map to slice
	roleSlice := make([]string, 0, len(roles))
	for role := range roles {
		roleSlice = append(roleSlice, string(role))
	}

	// Determine permissions based on actor's relationship to this user
	canEditProfile, _ := s.authSvc.CanEditUser(actorID, user.ID)
	canEditUsername, _ := s.authSvc.CanEditUsername(actorID, user.ID)
	canResetPassword, _ := s.authSvc.CanResetPassword(actorID, user.ID)

	return UserResponse{
		ID:       int64(user.ID),
		Username: user.Username,
		Email:    user.Email,
		Timezone: user.Locale.Timezone.Location.String(),
		Roles:    roleSlice,
		Permissions: map[string]bool{
			"canEditProfile":    canEditProfile,
			"canEditUsername":   canEditUsername,
			"canResetPassword":  canResetPassword,
			"canChangePassword": actorID == user.ID, // Only own password
		},
		Created: user.Created.Format(time.RFC3339),
		Updated: user.Updated.Format(time.RFC3339),
	}
}
