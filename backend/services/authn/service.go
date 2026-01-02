// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package authn implements an authentication service.
package authn

import (
	"context"
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
	log.Printf("authn: AuthenticateEmailCredentials(%q, %q)\n", email, secret)
	actor, err := s.authzSvc.GetActorByEmail(email)
	log.Printf("authn: AuthenticateEmailCredentials(%q, %q) %v\n", email, secret, err)
	if err != nil {
		// we need to plow through to prevent timing attacks
		actor = &domains.Actor{ID: domains.InvalidID}
	}
	log.Printf("authn: AuthenticateEmailCredentials(%q, %q) %+v\n", email, secret, *actor)
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
	if actor == nil {
		log.Printf("authn: authenticateCredentials(nil, %q)\n", secret)
		// burn time to mitigate user-enum: compare against dummy
		_ = bcrypt.CompareHashAndPassword(dummyHash, []byte(secret))
		return domains.InvalidID, domains.ErrInvalidCredentials
	}
	log.Printf("authn: authenticateCredentials(%+v, %q)\n", *actor, secret)
	if !s.authzSvc.CanAuthenticate(actor) {
		log.Printf("authn: authenticateCredentials(%+v, %q) canAuth false\n", *actor, secret)
		// burn time to mitigate user-enum: compare against dummy
		_ = bcrypt.CompareHashAndPassword(dummyHash, []byte(secret))
		return domains.InvalidID, domains.ErrInvalidCredentials
	}
	hashedPassword, err := s.db.Queries().ReadUserSecret(s.db.Context(), int64(actor.ID))
	if err != nil {
		log.Printf("authn: authenticateCredentials(%+v, %q) getUserSecret(%d) %v\n", *actor, secret, actor.ID, err)
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
func (s *Service) updateUserSecret(userId domains.ID, plainTextSecret string) (domains.ID, error) {
	updatedAt := time.Now().UTC().Unix()
	if err := domains.ValidatePassword(plainTextSecret); err != nil {
		return domains.InvalidID, errors.Join(domains.ErrInvalidCredentials, err)
	}
	hashedPassword, err := hashPassword(plainTextSecret)
	if err != nil {
		return domains.InvalidID, err
	}
	err = s.db.Queries().UpdateUserSecret(s.db.Context(), sqlc.UpdateUserSecretParams{
		UserID:            int64(userId),
		HashedPassword:    hashedPassword,
		PlaintextPassword: plainTextSecret,
		UpdatedAt:         updatedAt,
	})
	if err != nil {
		return domains.InvalidID, err
	}
	return userId, nil
}

// updateUserSecretInTransaction accepts sqlc.Queries to allow it to work within
// a database transaction. Call updateUserSecret if you want to use an immediate
// statement and generic context.
func (s *Service) updateUserSecretInTransaction(ctx context.Context, q *sqlc.Queries, userId domains.ID, plainTextSecret string, updatedAt time.Time) (domains.ID, error) {
	// log.Printf("user %d: password %q\n", userId, plainTextSecret)
	if err := domains.ValidatePassword(plainTextSecret); err != nil {
		return domains.InvalidID, errors.Join(domains.ErrInvalidCredentials, err)
	}
	hashedPassword, err := hashPassword(plainTextSecret)
	if err != nil {
		return domains.InvalidID, err
	}
	err = q.UpdateUserSecret(ctx, sqlc.UpdateUserSecretParams{
		UserID:            int64(userId),
		HashedPassword:    hashedPassword,
		PlaintextPassword: plainTextSecret,
		UpdatedAt:         updatedAt.UTC().Unix(),
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
	log.Printf("hashPassword(%q)\n", plainTextPassword)
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPasswordBytes), nil
}
