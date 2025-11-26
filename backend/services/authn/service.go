// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package authn implements an authentication service.
package authn

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/playbymail/ottoapp/backend/domains"
	"github.com/playbymail/ottoapp/backend/services/authz"
	"github.com/playbymail/ottoapp/backend/stores/sqlite"
	"github.com/playbymail/ottoapp/backend/stores/sqlite/sqlc"
	"golang.org/x/crypto/bcrypt"
)

// Service provides authentication operations.
type Service struct {
	db       *sqlite.DB
	authzSvc *authz.Service
}

func New(db *sqlite.DB, authzSvc *authz.Service) *Service {
	return &Service{db: db, authzSvc: authzSvc}
}

// AuthenticateActorCredentials verifies credentials
func (s *Service) AuthenticateActorCredentials(actor *domains.Actor, secret string) (domains.ID, error) {
	return s.authenticateCredentials(actor, secret)
}

// AuthenticateEmailCredentials verifies the user's credentials (email + secret).
func (s *Service) AuthenticateEmailCredentials(email, secret string) (domains.ID, error) {
	actor, err := s.authzSvc.GetActorByEmail(email)
	if err != nil {
		// we need to plow through to prevent timing attacks
		actor = &domains.Actor{ID: domains.InvalidID}
	}
	return s.authenticateCredentials(actor, secret)
}

// UpdateCredentials updates the target's credentials after verifying the
// actor is authorized.
//
// The oldSecret parameter is only used when actor == target.
// It's okay to pass in an empty string for all other cases.
//
// Note: database errors are returned as authorization errors.
func (s *Service) UpdateCredentials(actor, target *domains.Actor, oldSecret, newSecret string) (domains.ID, error) {
	if actor == nil || target == nil {
		return domains.InvalidID, domains.ErrNotAuthorized
	}
	if !s.authzSvc.CanUpdateTargetCredentials(actor, target) {
		return domains.InvalidID, domains.ErrNotAuthorized
	}
	if actor.ID == target.ID {
		// users must re-authenticate to change their own secret
		_, err := s.AuthenticateActorCredentials(actor, oldSecret)
		if err != nil {
			// failed authentication
			log.Printf("[auth] actor %d: target %d: authentication %v", actor.ID, target.ID, err)
			return domains.InvalidID, domains.ErrNotAuthorized
		}
	}
	return s.updateUserSecret(target.ID, newSecret)
}

var dummyHash = []byte("$2a$10$uG3ThGlwW4vB0hUHd8OQ8u4JXkiS2EeMaZDD8f5U2J1cG6r3G5gW6") // random valid bcrypt

// authenticateCredentials verifies credentials.
//
// It uses dummy bcrypt comparisons on error paths to provide
// some resistance to timing or user-enumeration attacks.
// Not production grade but good enough for our environment.
//
// If authentication succeeds, the actor's last login time
// is updated.
func (s *Service) authenticateCredentials(actor *domains.Actor, secret string) (domains.ID, error) {
	if actor == nil || !s.authzSvc.CanAuthenticate(actor) {
		// burn time to mitigate user-enum: compare against dummy
		_ = bcrypt.CompareHashAndPassword(dummyHash, []byte(secret))
		return domains.InvalidID, domains.ErrInvalidCredentials
	}
	hashedPassword, err := s.db.Queries().GetUserSecret(s.db.Context(), int64(actor.ID))
	if err != nil {
		// burn time to mitigate user-enum: compare against dummy
		_ = bcrypt.CompareHashAndPassword(dummyHash, []byte(secret))
		return domains.InvalidID, domains.ErrInvalidCredentials
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(secret))
	if err != nil {
		return domains.InvalidID, domains.ErrInvalidCredentials
	}
	_ = s.db.Queries().UpdateUserLastLogin(s.db.Context(), sqlc.UpdateUserLastLoginParams{
		UserID:    int64(actor.ID),
		LastLogin: time.Now().UTC().Unix(),
	})

	return actor.ID, nil
}

// updateUserSecret updates a secret without checking authorization.
// It forwards the update to upserUserSecret using a generic transaction and context.
func (s *Service) updateUserSecret(userID domains.ID, newPlainTextSecret string) (domains.ID, error) {
	return s.upsertUserSecret(s.db.Context(), s.db.Queries(), userID, newPlainTextSecret, time.Now().UTC())
}

// upsertUserSecret does the actual insert or update of a user secret.
// It requires a sqlc.Queries parameter because it expects that we will want
// to call it within transactions sometimes. Call updateUserSecret if you
// want to use an immediate statement and generic context.
func (s *Service) upsertUserSecret(ctx context.Context, q *sqlc.Queries, userId domains.ID, plainTextSecret string, now time.Time) (domains.ID, error) {
	// log.Printf("user %d: password %q\n", userId, plainTextSecret)
	if err := domains.ValidatePassword(plainTextSecret); err != nil {
		return domains.InvalidID, errors.Join(domains.ErrInvalidCredentials, err)
	}
	hashedPassword, err := hashPassword(plainTextSecret)
	if err != nil {
		return domains.InvalidID, err
	}
	err = q.UpsertUserSecret(ctx, sqlc.UpsertUserSecretParams{
		UserID:            int64(userId),
		HashedPassword:    hashedPassword,
		PlaintextPassword: sql.NullString{Valid: true, String: plainTextSecret},
		CreatedAt:         now.UTC().Unix(),
		UpdatedAt:         now.UTC().Unix(),
	})
	if err != nil {
		return domains.InvalidID, err
	}
	return userId, nil
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
