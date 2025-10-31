// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package sqlite

import (
	"errors"
	"fmt"
	"strings"

	"github.com/playbymail/ottoapp/backend/domains"
	"github.com/playbymail/ottoapp/backend/stores/sqlite/sqlc"
	"golang.org/x/crypto/bcrypt"
)

// AuthenticateUser  verifies the user's credentials.
// Callers should implement a delay returning until a timer has expired to prevent timing attacks.
//
//	timingAttackGuard := time.NewTimer(500 * time.Microsecond)
//	// call this function
//	// delay until the timer expires to prevent timing attacks
//	<-timingAttackGuard.C
func (db *DB) AuthenticateUser(handle, plainTextSecret string) (domains.ID, error) {
	var hashedPassword string
	id, err := db.q.GetUserIDByHandle(db.ctx, handle)
	if err != nil {
		id, err = 0, domains.ErrInvalidCredentials
	} else if hashedPassword, err = db.q.GetUserSecrets(db.ctx, id); err != nil {
		id, err = 0, domains.ErrInvalidCredentials
	} else if err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainTextSecret)); err != nil {
		id, err = 0, domains.ErrInvalidCredentials
	}
	return domains.ID(id), err
}

// HashPassword uses the cheapest bcrypt cost to hash the password because we are not going to use
// it for anything other than authentication in non-production environments.
func HashPassword(plainTextPassword string) (string, error) {
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hashedPasswordBytes), err
}

func (db *DB) UpdateUserPassword(handle string, newPlainTextSecret string) error {
	if strings.TrimSpace(newPlainTextSecret) != newPlainTextSecret || len(newPlainTextSecret) < 4 {
		return domains.ErrInvalidPassword
	}

	userID, err := db.q.GetUserIDByHandle(db.ctx, handle)
	if err != nil {
		return errors.Join(fmt.Errorf("%q: user not found", handle), err)
	}

	hashedPassword, err := HashPassword(newPlainTextSecret)
	if err != nil {
		return err
	}
	err = db.q.UpdateUserPassword(db.ctx, sqlc.UpdateUserPasswordParams{
		UserID:         userID,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		return errors.Join(fmt.Errorf("%q: update failed", handle), err)
	}
	return nil
}

func ValidatePassword(plainTextPassword string) bool {
	if strings.TrimSpace(plainTextPassword) != plainTextPassword {
		return false
	} else if !(4 < len(plainTextPassword) && len(plainTextPassword) < 128) {
		return false
	}
	return true
}

func ValidateRole(role string) bool {
	return containsWord(role, "guest", "chief", "admin")
}

func containsWord(word string, list ...string) bool {
	for _, elem := range list {
		if elem == word {
			return true
		}
	}
	return false
}
