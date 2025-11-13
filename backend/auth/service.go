// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package auth implements an authentication / authorization service.
package auth

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/playbymail/ottoapp/backend/domains"
	"github.com/playbymail/ottoapp/backend/stores/sqlite"
	"github.com/playbymail/ottoapp/backend/stores/sqlite/sqlc"
	"golang.org/x/crypto/bcrypt"
)

var dummyHash = []byte("$2a$10$uG3ThGlwW4vB0hUHd8OQ8u4JXkiS2EeMaZDD8f5U2J1cG6r3G5gW6") // random valid bcrypt

// Service provides authentication and authorization operations.
type Service struct {
	db *sqlite.DB
}

func New(db *sqlite.DB) *Service {
	return &Service{db: db}
}

// GetActor extracts the authenticated user ID from the request context.
// This is a convenience helper that assumes sessionMiddleware has run.
func (s *Service) GetActor(r *http.Request) (domains.ID, error) {
	userID, ok := r.Context().Value(domains.ContextKeyUserID).(domains.ID)
	if !ok || userID == domains.InvalidID {
		return domains.InvalidID, domains.ErrNotAuthenticated
	}
	return userID, nil
}

// AuthenticateUser verifies the user's credentials (username + password).
func (s *Service) AuthenticateUser(userName, plainTextSecret string) (domains.ID, error) {
	q := s.db.Queries()
	ctx := s.db.Context()

	user, err := q.GetUserByUsername(ctx, userName)
	if err != nil {
		// mitigate user-enum: compare against dummy
		_ = bcrypt.CompareHashAndPassword(dummyHash, []byte(plainTextSecret))
		return domains.InvalidID, domains.ErrInvalidCredentials
	}

	return s.verifyPasswordByID(user.UserID, plainTextSecret)
}

// AuthenticateWithEmailSecret verifies the user's credentials (email + password).
func (s *Service) AuthenticateWithEmailSecret(email, plainTextSecret string) (domains.ID, error) {
	q := s.db.Queries()
	ctx := s.db.Context()

	user, err := q.GetUserByEmail(ctx, email)
	if err != nil {
		// mitigate user-enum: compare against dummy
		_ = bcrypt.CompareHashAndPassword(dummyHash, []byte(plainTextSecret))
		return domains.InvalidID, domains.ErrInvalidCredentials
	}

	return s.verifyPasswordByID(user.UserID, plainTextSecret)
}

// CreateUserSecret requires a sqlc.Queries parameter because it expects that we
// will want to call it within transactions sometimes.
func (s *Service) CreateUserSecret(ctx context.Context, q *sqlc.Queries, userId domains.ID, plainTextSecret string, now time.Time) error {
	if !ValidatePassword(plainTextSecret) {
		return domains.ErrInvalidCredentials
	}
	hashedPassword, err := hashPassword(plainTextSecret)
	if err != nil {
		return err
	}
	return q.CreateUserSecret(ctx, sqlc.CreateUserSecretParams{
		UserID:         int64(userId),
		HashedPassword: hashedPassword,
		CreatedAt:      now.UTC().Unix(),
		UpdatedAt:      now.UTC().Unix(),
	})
}

func (s *Service) GetUserRoles(userID domains.ID) (domains.Roles, error) {
	q := s.db.Queries()
	ctx := s.db.Context()

	userRoles, err := q.GetUserRoles(ctx, int64(userID))
	if err != nil {
		return nil, err
	}
	roles := map[domains.Role]bool{}
	for _, role := range userRoles {
		roles[domains.Role(role)] = true
	}
	return roles, nil
}

func (s *Service) UpdateUserSecret(userID domains.ID, newPlainTextSecret string) error {
	if !ValidatePassword(newPlainTextSecret) {
		return domains.ErrInvalidPassword
	}

	q := s.db.Queries()
	ctx := s.db.Context()

	hashedPassword, err := hashPassword(newPlainTextSecret)
	if err != nil {
		return err
	}

	return q.UpdateUserSecret(ctx, sqlc.UpdateUserSecretParams{
		UserID:         int64(userID),
		HashedPassword: hashedPassword,
		UpdatedAt:      time.Now().UTC().Unix(),
	})
}

// verifyPasswordByID fetches the user's hash and compares it.
func (s *Service) verifyPasswordByID(userID int64, plain string) (domains.ID, error) {
	q := s.db.Queries()
	ctx := s.db.Context()

	hashedPassword, err := q.GetUserSecret(ctx, userID)
	if err != nil {
		return domains.InvalidID, domains.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plain)); err != nil {
		return domains.InvalidID, domains.ErrInvalidCredentials
	}

	return domains.ID(userID), nil
}

func containsWord(word string, list ...string) bool {
	for _, elem := range list {
		if elem == word {
			return true
		}
	}
	return false
}

// hashPassword hashes the password with a reasonable cost.
func hashPassword(plainTextPassword string) (string, error) {
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPasswordBytes), nil
}

func ValidatePassword(plainTextPassword string) bool {
	// no leading/trailing whitespace
	if strings.TrimSpace(plainTextPassword) != plainTextPassword {
		return false
	}
	// length bounds 8...128 bytes
	l := len(plainTextPassword)
	if !(8 <= l && l <= 128) {
		return false
	}
	return true
}

func ValidateRole(role string) bool {
	return containsWord(role, "active", "sysop", "admin", "player", "guest", "user", "tn3", "tn3.1")
}

// AssignRole assigns a role to a user.
func (s *Service) AssignRole(userID domains.ID, roleID string) error {
	if !ValidateRole(roleID) {
		return domains.ErrInvalidRole
	}

	q := s.db.Queries()
	ctx := s.db.Context()
	now := time.Now().UTC()

	return q.AssignUserRole(ctx, sqlc.AssignUserRoleParams{
		UserID:    int64(userID),
		RoleID:    roleID,
		CreatedAt: now.Unix(),
		UpdatedAt: now.Unix(),
	})
}

// RemoveRole removes a role from a user.
func (s *Service) RemoveRole(userID domains.ID, roleID string) error {
	if !ValidateRole(roleID) {
		return domains.ErrInvalidRole
	}

	q := s.db.Queries()
	ctx := s.db.Context()

	return q.RemoveUserRole(ctx, sqlc.RemoveUserRoleParams{
		UserID: int64(userID),
		RoleID: roleID,
	})
}

// HasRole checks if a user has a specific role.
func (s *Service) HasRole(userID domains.ID, roleID string) (bool, error) {
	roles, err := s.GetUserRoles(userID)
	if err != nil {
		return false, err
	}
	return roles[domains.Role(roleID)], nil
}

// IsAdmin checks if a user has the admin role.
func (s *Service) IsAdmin(userID domains.ID) (bool, error) {
	return s.HasRole(userID, "admin")
}

// IsSysop checks if a user has the sysop role.
func (s *Service) IsSysop(userID domains.ID) (bool, error) {
	return s.HasRole(userID, "sysop")
}

// CanViewUser returns true if actor can view target user's profile.
// Rules: user can view self, admin can view all (excluding sysop)
func (s *Service) CanViewUser(actorID, targetID domains.ID) (bool, error) {
	// User can view themselves
	if actorID == targetID {
		return true, nil
	}

	// Check if actor is admin
	isAdmin, err := s.IsAdmin(actorID)
	if err != nil {
		return false, err
	}
	if !isAdmin {
		return false, nil
	}

	// Admin cannot view sysop
	targetIsSysop, err := s.IsSysop(targetID)
	if err != nil {
		return false, err
	}
	if targetIsSysop {
		return false, nil
	}

	// Admin can all other users, including other admins
	return true, nil
}

// CanEditUser checks if actor can edit target user's profile.
// Rules: user can edit self, admin can edit non-admins (excluding sysop).
func (s *Service) CanEditUser(actorID, targetID domains.ID) (bool, error) {
	// User can edit themselves
	if actorID == targetID {
		return true, nil
	}

	// Check if actor is admin
	isAdmin, err := s.IsAdmin(actorID)
	if err != nil {
		return false, err
	}
	if !isAdmin {
		return false, nil
	}

	// Admin cannot edit sysop
	targetIsSysop, err := s.IsSysop(targetID)
	if err != nil {
		return false, err
	}
	if targetIsSysop {
		return false, nil
	}

	// Admin cannot edit other admins
	targetIsAdmin, err := s.IsAdmin(targetID)
	if err != nil {
		return false, err
	}
	if targetIsAdmin {
		return false, nil
	}

	return true, nil
}

// CanEditUsername checks if actor can edit target user's username.
// Only admins can edit usernames.
func (s *Service) CanEditUsername(actorID, targetID domains.ID) (bool, error) {
	isAdmin, err := s.IsAdmin(actorID)
	if err != nil {
		return false, err
	}
	if !isAdmin {
		return false, nil
	}

	// Admin cannot edit sysop
	targetIsSysop, err := s.IsSysop(targetID)
	if err != nil {
		return false, err
	}
	if targetIsSysop {
		return false, nil
	}

	// Admin cannot edit other admins
	targetIsAdmin, err := s.IsAdmin(targetID)
	if err != nil {
		return false, err
	}
	if targetIsAdmin {
		return false, nil
	}

	return true, nil
}

// CanResetPassword checks if actor can reset target user's password.
// Only admins can reset passwords for non-admins.
func (s *Service) CanResetPassword(actorID, targetID domains.ID) (bool, error) {
	isAdmin, err := s.IsAdmin(actorID)
	if err != nil {
		return false, err
	}
	if !isAdmin {
		return false, nil
	}

	// Admin cannot reset sysop password (sysop has no password)
	targetIsSysop, err := s.IsSysop(targetID)
	if err != nil {
		return false, err
	}
	if targetIsSysop {
		return false, nil
	}

	// Admin cannot reset other admin passwords
	targetIsAdmin, err := s.IsAdmin(targetID)
	if err != nil {
		return false, err
	}
	if targetIsAdmin {
		return false, nil
	}

	return true, nil
}

// CanCreateUser checks if actor can create new users.
// Only admins can create users.
func (s *Service) CanCreateUser(actorID domains.ID) (bool, error) {
	return s.IsAdmin(actorID)
}

// CanListUsers checks if actor can list all users.
// Only admins can list users.
func (s *Service) CanListUsers(actorID domains.ID) (bool, error) {
	return s.IsAdmin(actorID)
}

// CanManageRoles checks if actor can modify target user's roles.
// Only admins can manage roles, and they cannot modify sysop or admin roles.
func (s *Service) CanManageRoles(actorID, targetID domains.ID) (bool, error) {
	isAdmin, err := s.IsAdmin(actorID)
	if err != nil {
		return false, err
	}
	if !isAdmin {
		return false, nil
	}

	// Cannot modify sysop roles
	targetIsSysop, err := s.IsSysop(targetID)
	if err != nil {
		return false, err
	}
	if targetIsSysop {
		return false, nil
	}

	// Cannot modify admin roles
	targetIsAdmin, err := s.IsAdmin(targetID)
	if err != nil {
		return false, err
	}
	if targetIsAdmin {
		return false, nil
	}

	return true, nil
}

// CanChangeOwnPassword checks if actor can change their own password.
// Users can only change their own password (not others').
func (s *Service) CanChangeOwnPassword(actorID, targetID domains.ID) (bool, error) {
	return actorID == targetID, nil
}

// BuildActorAuth returns the actor's authorizations for the target user.
// Includes roles and permissions that the frontend needs for UI decisions.
func (s *Service) BuildActorAuth(actorID, targetID domains.ID) (*domains.ActorAuthorizations, error) {
	// Get target user's roles
	roles, err := s.GetUserRoles(targetID)
	if err != nil {
		return nil, err
	}

	// Convert roles map to slice
	roleSlice := make([]string, 0, len(roles))
	for role := range roles {
		roleSlice = append(roleSlice, string(role))
	}

	// Determine what actor can do to target
	canEditProfile, _ := s.CanEditUser(actorID, targetID)
	canEditUsername, _ := s.CanEditUsername(actorID, targetID)
	canResetPassword, _ := s.CanResetPassword(actorID, targetID)
	canChangePassword, _ := s.CanChangeOwnPassword(actorID, targetID)

	return &domains.ActorAuthorizations{
		Roles: roleSlice,
		Permissions: map[string]bool{
			"canEditProfile":    canEditProfile,
			"canEditUsername":   canEditUsername,
			"canResetPassword":  canResetPassword,
			"canChangePassword": canChangePassword,
		},
	}, nil
}
