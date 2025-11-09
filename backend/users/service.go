// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package users

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/playbymail/ottoapp/backend/auth"
	"github.com/playbymail/ottoapp/backend/domains"
	"github.com/playbymail/ottoapp/backend/iana"
	"github.com/playbymail/ottoapp/backend/stores/sqlite"
	"github.com/playbymail/ottoapp/backend/stores/sqlite/sqlc"
)

// Service provides user operations.
type Service struct {
	db      *sqlite.DB
	authSvc *auth.Service
}

func New(db *sqlite.DB, authSvc *auth.Service) *Service {
	return &Service{
		db:      db,
		authSvc: authSvc,
	}
}

// ChangePassword verifies the user's current credentials and, if valid,
// updates the password via the auth service.
func (s *Service) ChangePassword(email, oldPassword, newPassword string) error {
	userID, err := s.authSvc.AuthenticateWithEmailSecret(email, oldPassword)
	if err != nil {
		// could log as suspicious if this happens often
		return domains.ErrInvalidCredentials
	}

	// let auth enforce password policy
	return s.authSvc.UpdateUserSecret(userID, newPassword)
}

func (s *Service) CreateUser(handle, email, plainTextSecret string, timezone *time.Location) (*domains.User_t, error) {
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

	// start transaction
	tx, err := s.db.Stdlib().BeginTx(s.db.Context(), nil)
	if err != nil {
		return nil, err
	}
	// rollback if we return early; harmless after commit
	defer tx.Rollback()

	qtx := s.db.Queries().WithTx(tx)

	now := time.Now().UTC()
	userId, err := qtx.CreateUser(s.db.Context(), sqlc.CreateUserParams{
		Handle:    handle,
		Email:     email,
		Timezone:  timeZone,
		CreatedAt: now.Unix(),
		UpdatedAt: now.Unix(),
	})
	if err != nil {
		return nil, err
	}

	err = s.authSvc.CreateUserSecret(s.db.Context(), qtx, domains.ID(userId), plainTextSecret, now)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// safe to use the normal queries again after commit
	return s.GetUserByID(domains.ID(userId))
}

// GetUserByEmail returns user data associated with the email.
// Warning: callers expect this to return the same data that would be returned from GetUserByID!
func (s *Service) GetUserByEmail(email string) (*domains.User_t, error) {
	q := s.db.Queries()
	ctx := s.db.Context()

	row, err := q.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	loc, ok := iana.NormalizeTimeZone(row.Timezone)
	if !ok {
		return nil, domains.ErrInvalidTimezone
	}

	user := &domains.User_t{
		ID:       domains.ID(row.UserID),
		Username: row.Handle,
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

	return user, nil
}

// GetUserByHandle returns user data associated with the handle.
// Warning: callers expect this to return the same data that would be returned from GetUserByID!
func (s *Service) GetUserByHandle(handle string) (*domains.User_t, error) {
	q := s.db.Queries()
	ctx := s.db.Context()

	row, err := q.GetUserByHandle(ctx, handle)
	if err != nil {
		return nil, err
	}
	loc, ok := iana.NormalizeTimeZone(row.Timezone)
	if !ok {
		return nil, domains.ErrInvalidTimezone
	}

	user := &domains.User_t{
		ID:       domains.ID(row.UserID),
		Username: row.Handle,
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

	return user, nil
}

func (s *Service) GetUserIDByEmail(email string) (domains.ID, error) {
	id, err := s.db.Queries().GetUserIDByEmail(s.db.Context(), email)
	if err != nil {
		return domains.InvalidID, err
	}
	return domains.ID(id), nil
}

func (s *Service) GetUserIDByHandle(handle string) (domains.ID, error) {
	id, err := s.db.Queries().GetUserIDByHandle(s.db.Context(), handle)
	if err != nil {
		return domains.InvalidID, err
	}
	return domains.ID(id), nil
}

// GetUserByID returns the user data associated with the given ID.
// If the user does not exist, it returns an error.
func (s *Service) GetUserByID(userID domains.ID) (*domains.User_t, error) {
	q := s.db.Queries()
	ctx := s.db.Context()

	row, err := q.GetUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}
	loc, ok := iana.NormalizeTimeZone(row.Timezone)
	if !ok {
		return nil, domains.ErrInvalidTimezone
	}

	user := &domains.User_t{
		ID:       userID,
		Username: row.Handle,
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

	return user, nil
}

func (s *Service) UpdateUser(userId domains.ID, newHandle *string, newEmail *string, newTimeZone *time.Location) error {
	q := s.db.Queries()
	ctx := s.db.Context()

	// short circuit
	if userId == domains.InvalidID {
		return domains.ErrInvalidCredentials
	} else if newHandle == nil && newEmail == nil && newTimeZone == nil {
		return nil
	}

	// fetch the current user values
	user, err := q.GetUser(ctx, int64(userId))
	if err != nil {
		return domains.ErrInvalidCredentials
	}

	// merge the updated values into the current values
	handle := user.Handle
	if newHandle != nil {
		handle = strings.ToLower(*newHandle)
		if !ValidateHandle(handle) {
			return errors.Join(fmt.Errorf("%q: invalid handle", handle), domains.ErrInvalidHandle)
		}
	}
	email := user.Email
	if newEmail != nil {
		email = strings.ToLower(*newEmail)
		if !ValidateEmail(email) {
			return errors.Join(fmt.Errorf("%q: invalid email", email), domains.ErrInvalidEmail)
		}
	}
	timeZone := user.Timezone
	if newTimeZone != nil {
		timeZone = newTimeZone.String()
	}

	err = q.UpdateUser(ctx, sqlc.UpdateUserParams{
		UserID:    user.UserID,
		Handle:    handle,
		Email:     email,
		Timezone:  timeZone,
		UpdatedAt: time.Now().UTC().Unix(),
	})
	if err != nil {
		return errors.Join(fmt.Errorf("user %d: update failed", userId), err)
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
	} else if !(3 <= len(handle) && len(handle) < 14) {
		return false
	} else if strings.Contains(handle, "@") {
		return false
	}
	return true
}
