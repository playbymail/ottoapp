// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package sqlite

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/playbymail/ottoapp/backend/domains"
	"github.com/playbymail/ottoapp/backend/iana"
	"github.com/playbymail/ottoapp/backend/stores/sqlite/sqlc"
)

func (db *DB) CreateUser(handle, email, plainTextSecret string, timezone *time.Location) (*domains.User_t, error) {
	handle = strings.ToLower(handle)
	if !ValidateHandle(handle) {
		return nil, domains.ErrInvalidHandle
	}
	email = strings.ToLower(email)
	if !ValidateEmail(email) {
		return nil, domains.ErrInvalidEmail
	}
	if timezone == nil {
		return nil, fmt.Errorf("timezone is required")
	}
	timeZone, ok := iana.CanonicalName(timezone.String())
	if !ok {
		return nil, fmt.Errorf("%q: invalid timezone", timezone.String())
	}

	// hash the password. can fail if the password is too long.
	hashedPassword, err := HashPassword(plainTextSecret)
	if err != nil {
		return nil, err
	}

	// todo: wrap this in a transaction

	id, err := db.q.CreateUser(db.ctx, sqlc.CreateUserParams{
		Handle:   handle,
		Email:    email,
		Timezone: timeZone,
	})
	if err != nil {
		return nil, err
	}

	// note: we let LastLogin be the zero-value for time.Time, which means never logged in.
	err = db.q.CreateUserSecrets(db.ctx, sqlc.CreateUserSecretsParams{
		UserID:         id,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		log.Printf("[sqldb] need to rollback create user\n")
		return nil, err
	}

	return db.GetUserByID(domains.ID(id))
}

// GetUserByEmail returns the user with the given email.
// If the user does not exist, it returns an error.
func (db *DB) GetUserByEmail(email string) (*domains.User_t, error) {
	id, err := db.q.GetUserIDByEmail(db.ctx, email)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("%q: user not found", email), err)
	}
	return db.GetUserByID(domains.ID(id))
}

// GetUserByID returns the user with the given ID.
// If the user does not exist, it returns an error.
func (db *DB) GetUserByID(userID domains.ID) (*domains.User_t, error) {
	row, err := db.q.GetUser(db.ctx, int64(userID))
	if err != nil {
		return nil, err
	}
	// convert row.Timezone to a time.Location
	loc, err := time.LoadLocation(row.Timezone)
	if err != nil {
		return nil, err
	}

	user := &domains.User_t{
		ID:    userID,
		Email: row.Email,
		Locale: domains.UserLocale_t{
			DateFormat: "2006-01-02",
			Timezone: domains.UserTimezone_t{
				Location: loc,
			},
		},
		Created: row.CreatedAt,
		Updated: row.UpdatedAt,
	}

	return user, nil
}

func (db *DB) UpdateUser(handle string, newEmail *string, newTimeZone *time.Location) error {
	userID, err := db.q.GetUserIDByHandle(db.ctx, handle)
	if err != nil {
		return errors.Join(fmt.Errorf("%q: user not found", handle), err)
	}
	if newEmail == nil && newTimeZone == nil {
		return nil
	}
	// fetch the current user values
	row, err := db.q.GetUser(db.ctx, userID)
	if err != nil {
		return errors.Join(fmt.Errorf("%d: user not found", userID), err)
	}
	// merge the updated values into the current values
	email := row.Email
	if newEmail != nil {
		email = strings.ToLower(*newEmail)
		if !ValidateEmail(email) {
			return errors.Join(fmt.Errorf("%q: invalid email", email), domains.ErrInvalidEmail)
		}
	}
	timeZone := row.Timezone
	if newTimeZone != nil {
		timeZone = newTimeZone.String()
	}

	err = db.q.UpdateUser(db.ctx, sqlc.UpdateUserParams{
		UserID:   userID,
		Handle:   handle,
		Email:    email,
		Timezone: timeZone,
	})
	if err != nil {
		return errors.Join(fmt.Errorf("%q: update failed", email), err)
	}
	return nil
}

func ValidateEmail(email string) bool {
	if email != strings.TrimSpace(email) {
		return false
	} else if !(4 < len(email) && len(email) < 48) {
		return false
	} else if !strings.Contains(email, "@") {
		return false
	}
	return true
}

func ValidateHandle(handle string) bool {
	if handle != strings.TrimSpace(handle) {
		return false
	} else if !(4 < len(handle) && len(handle) < 14) {
		return false
	} else if strings.Contains(handle, "@") {
		return false
	}
	return true
}
