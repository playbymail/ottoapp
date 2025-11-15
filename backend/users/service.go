// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package users

//go:generate go run ../../cmd/godel -input handlers.go -struct UserView -output ../../frontend/app/models/user.js

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/playbymail/ottoapp/backend/auth"
	"github.com/playbymail/ottoapp/backend/domains"
	"github.com/playbymail/ottoapp/backend/iana"
	"github.com/playbymail/ottoapp/backend/stores/sqlite"
	"github.com/playbymail/ottoapp/backend/stores/sqlite/sqlc"
)

// SessionService defines the interface for session operations needed by handlers
type SessionService interface {
	GetCurrentUserID(r *http.Request) (domains.ID, error)
}

// Service provides user operations.
type Service struct {
	db          *sqlite.DB
	authSvc     *auth.Service
	tzSvc       *iana.Service
	sessionsSvc SessionService
}

func New(db *sqlite.DB, authSvc *auth.Service, tzSvc *iana.Service) *Service {
	return &Service{
		db:      db,
		authSvc: authSvc,
		tzSvc:   tzSvc,
	}
}

// SetSessionService injects the session service (to avoid circular dependency)
func (s *Service) SetSessionService(sessionsSvc SessionService) {
	s.sessionsSvc = sessionsSvc
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

func (s *Service) CreateUser(userName, email, handle, plainTextSecret string, timezone *time.Location) (*domains.User_t, error) {
	if !s.ValidateUsername(userName) {
		return nil, domains.ErrInvalidUsername
	}
	email = strings.ToLower(email)
	if !s.ValidateEmail(email) {
		return nil, domains.ErrInvalidEmail
	}
	handle = strings.ToLower(handle)
	if !s.ValidateHandle(handle) {
		return nil, domains.ErrInvalidHandle
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
		Username:  userName,
		Email:     email,
		Handle:    handle,
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

	// Assign default roles for CLI-created users: "active" and "user"
	err = qtx.AssignUserRole(s.db.Context(), sqlc.AssignUserRoleParams{
		UserID:    userId,
		RoleID:    "active",
		CreatedAt: now.Unix(),
		UpdatedAt: now.Unix(),
	})
	if err != nil {
		return nil, err
	}

	err = qtx.AssignUserRole(s.db.Context(), sqlc.AssignUserRoleParams{
		UserID:    userId,
		RoleID:    "user",
		CreatedAt: now.Unix(),
		UpdatedAt: now.Unix(),
	})
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
	loc, err := s.tzSvc.Location(row.Timezone)
	if err != nil {
		return nil, domains.ErrInvalidTimezone
	}

	user := &domains.User_t{
		ID:       domains.ID(row.UserID),
		Username: row.Username,
		Email:    row.Email,
		Handle:   row.Handle,
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

// GetUserByUsername returns user data associated with the username.
// Warning: callers expect this to return the same data that would be returned from GetUserByID!
func (s *Service) GetUserByUsername(userName string) (*domains.User_t, error) {
	q := s.db.Queries()
	ctx := s.db.Context()

	row, err := q.GetUserByUsername(ctx, userName)
	if err != nil {
		return nil, err
	}
	loc, err := s.tzSvc.Location(row.Timezone)
	if err != nil {
		return nil, domains.ErrInvalidTimezone
	}

	user := &domains.User_t{
		ID:       domains.ID(row.UserID),
		Username: row.Username,
		Email:    row.Email,
		Handle:   row.Handle,
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

func (s *Service) GetUserIDByUsername(userName string) (domains.ID, error) {
	id, err := s.db.Queries().GetUserIDByUsername(s.db.Context(), userName)
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

	row, err := q.GetUserByID(ctx, int64(userID))
	if err != nil {
		return nil, err
	}
	loc, err := s.tzSvc.Location(row.Timezone)
	if err != nil {
		return nil, domains.ErrInvalidTimezone
	}

	user := &domains.User_t{
		ID:       userID,
		Username: row.Username,
		Email:    row.Email,
		Handle:   row.Handle,
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

func (s *Service) UpdateUser(userId domains.ID, newUserName *string, newEmail *string, newTimeZone *time.Location) error {
	q := s.db.Queries()
	ctx := s.db.Context()

	// short circuit
	if userId == domains.InvalidID {
		return domains.ErrInvalidCredentials
	} else if newUserName == nil && newEmail == nil && newTimeZone == nil {
		return nil
	}

	// fetch the current user values
	user, err := q.GetUserByID(ctx, int64(userId))
	if err != nil {
		return domains.ErrInvalidCredentials
	}

	// merge the updated values into the current values
	userName := user.Username
	if newUserName != nil {
		userName = strings.ToLower(*newUserName)
		if !s.ValidateUsername(userName) {
			return errors.Join(fmt.Errorf("%q: invalid userName", userName), domains.ErrInvalidUsername)
		}
	}
	email := user.Email
	if newEmail != nil {
		email = strings.ToLower(*newEmail)
		if !s.ValidateEmail(email) {
			return errors.Join(fmt.Errorf("%q: invalid email", email), domains.ErrInvalidEmail)
		}
	}
	timeZone := user.Timezone
	if newTimeZone != nil {
		timeZone = newTimeZone.String()
	}

	err = q.UpdateUser(ctx, sqlc.UpdateUserParams{
		UserID:    user.UserID,
		Username:  userName,
		Email:     email,
		Handle:    user.Handle,
		Timezone:  timeZone,
		UpdatedAt: time.Now().UTC().Unix(),
	})
	if err != nil {
		return errors.Join(fmt.Errorf("user %d: update failed", userId), err)
	}
	return nil
}

func (s *Service) ValidateEmail(email string) bool {
	return ValidateEmail(email)
}

func (s *Service) ValidateHandle(handle string) bool {
	return ValidateHandle(handle)
}

func (s *Service) ValidateUsername(userName string) bool {
	return ValidateUsername(userName)
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
	// Must not be empty
	if len(handle) == 0 || len(handle) > 14 {
		return false
	}
	// Must be lowercase (no uppercase allowed)
	if handle != strings.ToLower(handle) {
		return false
	}
	// Must start with a letter
	if handle[0] < 'a' || handle[0] > 'z' {
		return false
	}
	// All characters must be lowercase letters, digits, or underscores
	for _, ch := range handle {
		if !((ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') || ch == '_') {
			return false
		}
	}
	return true
}

func ValidateUsername(userName string) bool {
	if userName != strings.TrimSpace(userName) {
		return false
	} else if !(3 <= len(userName) && len(userName) < 35) {
		return false
	} else if strings.Contains(userName, "@") {
		return false
	}
	return true
}
