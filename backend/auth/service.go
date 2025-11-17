// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package auth implements an authentication / authorization service.
package auth

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/playbymail/ottoapp/backend/domains"
	"github.com/playbymail/ottoapp/backend/stores/sqlite"
	"github.com/playbymail/ottoapp/backend/stores/sqlite/sqlc"
	"golang.org/x/crypto/bcrypt"
)

// Service provides authentication and authorization operations.
type Service struct {
	db *sqlite.DB
}

func New(db *sqlite.DB) *Service {
	return &Service{db: db}
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

// verifyPasswordByID fetches the user's hash and compares it.
func (s *Service) verifyPasswordByID(userID int64, plain string) (domains.ID, error) {
	panic("!consolidated")
}

// hashPassword hashes the password with a reasonable cost.
func hashPassword(plainTextPassword string) (string, error) {
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPasswordBytes), nil
}

// AssignRole assigns a role to a user.
func (s *Service) AssignRole(userID domains.ID, roleID string) error {
	if err := domains.ValidateRole(roleID); err != nil {
		return errors.Join(domains.ErrInvalidRole, err)
	}

	q := s.db.Queries()
	ctx := s.db.Context()
	now := time.Now().UTC()

	return q.UpsertUserRole(ctx, sqlc.UpsertUserRoleParams{
		UserID:    int64(userID),
		RoleID:    roleID,
		CreatedAt: now.Unix(),
		UpdatedAt: now.Unix(),
	})
}

// RemoveRole removes a role from a user.
func (s *Service) RemoveRole(userID domains.ID, roleID string) error {
	if err := domains.ValidateRole(roleID); err != nil {
		return errors.Join(domains.ErrInvalidRole, err)
	}

	q := s.db.Queries()
	ctx := s.db.Context()

	return q.RemoveUserRole(ctx, sqlc.RemoveUserRoleParams{
		UserID: int64(userID),
		RoleID: roleID,
	})
}

// BuildActorAuth returns the actor's authorizations for the target user.
// Includes roles and permissions that the frontend needs for UI decisions.
func (s *Service) BuildActorAuth(actor, target *domains.Actor) *domains.ActorAuthorizations {
	aa := &domains.ActorAuthorizations{Permissions: map[string]bool{}}
	if s.CanEditTarget(actor, target) {
		aa.Roles = append(aa.Roles, "canEditProfile")
		aa.Permissions["canEditProfile"] = true
	} else {
		aa.Permissions["canEditProfile"] = false
	}
	if s.CanEditTargetUsername(actor, target) {
		aa.Roles = append(aa.Roles, "canEditUsername")
		aa.Permissions["canEditUsername"] = true
	} else {
		aa.Permissions["canEditUsername"] = false
	}
	if s.CanResetTargetCredentials(actor, target) {
		aa.Roles = append(aa.Roles, "canResetPassword")
		aa.Permissions["canResetPassword"] = true
	} else {
		aa.Permissions["canResetPassword"] = false
	}
	if s.CanUpdateTargetCredentials(actor, target) {
		aa.Roles = append(aa.Roles, "canChangePassword")
		aa.Permissions["canChangePassword"] = true
	} else {
		aa.Permissions["canChangePassword"] = false
	}
	return aa
}

// updateUserSecret updates a secret without checking authorization.
// It forwards the update to upserUserSecret using a generic transaction
// and context.
func (s *Service) updateUserSecret(userID domains.ID, newPlainTextSecret string) error {
	return s.upsertUserSecret(s.db.Context(), s.db.Queries(), userID, newPlainTextSecret, time.Now().UTC())
}

// upsertUserSecret does the actual insert or update of a user secret.
// It requires a sqlc.Queries parameter because it expects that we
// will want to call it within transactions sometimes. Call updateUserSecret
// if you want to use an immediate statement and generic context.
func (s *Service) upsertUserSecret(ctx context.Context, q *sqlc.Queries, userId domains.ID, plainTextSecret string, now time.Time) error {
	// log.Printf("user %d: password %q\n", userId, plainTextSecret)
	if err := domains.ValidatePassword(plainTextSecret); err != nil {
		return errors.Join(domains.ErrInvalidCredentials, err)
	}
	hashedPassword, err := hashPassword(plainTextSecret)
	if err != nil {
		return err
	}
	return q.UpsertUserSecret(ctx, sqlc.UpsertUserSecretParams{
		UserID:            int64(userId),
		HashedPassword:    hashedPassword,
		PlaintextPassword: sql.NullString{Valid: true, String: plainTextSecret},
		CreatedAt:         now.UTC().Unix(),
		UpdatedAt:         now.UTC().Unix(),
	})
}
