// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package auth

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/playbymail/ottoapp/backend/domains"
	"github.com/playbymail/ottoapp/backend/stores/sqlite/sqlc"
	"golang.org/x/crypto/bcrypt"
)

const (
	// SysopId is asserted in the database initialization scripts
	// and initial database connection.
	SysopId = domains.ID(1)
)

// AuthenticateCredentials verifies credentials
func (s *Service) AuthenticateCredentials(actor *domains.Actor, secret string) error {
	return s.authenticateCredentials(actor, secret)
}

// UpdateCredentials updates the target's credentials after verifying the
// actor is authorized.
//
// The oldSecret parameter is only used when actor == target.
// It's okay to pass in an empty string for all other cases.
//
// Note: database errors are returned as authorization errors.
func (s *Service) UpdateCredentials(actor, target *domains.Actor, oldSecret, newSecret string) error {
	if actor == nil || target == nil {
		return domains.ErrNotAuthorized
	}
	if !s.CanUpdateTargetCredentials(actor, target) {
		return domains.ErrNotAuthorized
	}
	if actor.ID == target.ID {
		// users must re-authenticate to change their own secret
		err := s.AuthenticateCredentials(actor, oldSecret)
		if err != nil {
			// failed authentication
			log.Printf("[auth] actor %d: target %d: authentication %v", actor.ID, target.ID, err)
			return domains.ErrNotAuthorized
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
func (s *Service) authenticateCredentials(actor *domains.Actor, secret string) error {
	if actor == nil || !s.CanAuthenticate(actor) {
		// burn time
		_ = bcrypt.CompareHashAndPassword(dummyHash, []byte(secret))
		return domains.ErrInvalidCredentials
	}

	q := s.db.Queries()
	ctx := s.db.Context()
	if hashedPassword, err := q.GetUserSecret(ctx, int64(actor.ID)); err != nil {
		// mitigate user-enum: compare against dummy
		_ = bcrypt.CompareHashAndPassword(dummyHash, []byte(secret))
		return domains.ErrInvalidCredentials
	} else if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(secret)); err != nil {
		return domains.ErrInvalidCredentials
	}
	_ = q.UpdateUserLastLogin(ctx, sqlc.UpdateUserLastLoginParams{
		UserID:    int64(actor.ID),
		LastLogin: time.Now().UTC().Unix(),
	})

	return nil
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
