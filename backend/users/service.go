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
func (s *Service) UpsertUser(handle, email, userName string, timezone *time.Location) (*domains.User_t, error) {
	handle = strings.ToLower(handle)
	if err := domains.ValidateHandle(handle); err != nil {
		return nil, errors.Join(domains.ErrInvalidHandle, err)
	}
	if err := domains.ValidateUsername(userName); err != nil {
		return nil, errors.Join(domains.ErrInvalidUsername, err)
	}
	email = strings.ToLower(email)
	if err := domains.ValidateEmail(email); err != nil {
		return nil, errors.Join(domains.ErrInvalidEmail, err)
	}
	if timezone == nil {
		return nil, errors.Join(domains.ErrInvalidTimezone, fmt.Errorf("timezone is required"))
	}
	timeZone, ok := iana.CanonicalName(timezone.String())
	if !ok {
		return nil, errors.Join(domains.ErrInvalidTimezone, domains.ErrBadInput, fmt.Errorf("%q: invalid timezone", timezone.String()))
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
	userId, err := qtx.UpsertUser(s.db.Context(), sqlc.UpsertUserParams{
		Handle:    handle,
		Username:  userName,
		Email:     email,
		Timezone:  timeZone,
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

func (s *Service) GetUser(actor *domains.Actor) (*domains.User_t, error) {
	if !actor.IsValid() {
		return nil, domains.ErrBadInput
	}
	return s.getUserByID(actor.ID)
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

// GetUserByHandle returns user data associated with the handle.
// Warning: callers expect this to return the same data that would be returned from GetUserByID!
func (s *Service) GetUserByHandle(handle string) (*domains.User_t, error) {
	q := s.db.Queries()
	ctx := s.db.Context()

	row, err := q.GetUserByHandle(ctx, handle)
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
	return s.getUserByID(userID)
}

func (s *Service) ListUsersVisibleToActor(actor *domains.Actor, pageNum, pageSize int) ([]*UserView, error) {
	rows, err := s.db.Queries().ListUsersVisibleToActor(s.db.Context(), sqlc.ListUsersVisibleToActorParams{
		ActorID:  actor.ID,
		PageSize: pageSize,
		PageNum:  (pageNum - 1) * pageSize,
	})
	if err != nil {
		return nil, err
	}
	var view []*UserView
	for _, row := range rows {
		view = append(view, &UserView{
			ID:       fmt.Sprintf("%d", row.UserID),
			Username: row.Username,
			Email:    row.Email,
			Timezone: row.Timezone,
			//Roles:       row.Roles,       // if you prejoin/aggregate; otherwise compute once for the list
			//Permissions: row.Permissions, // same note as above
			CreatedAt: time.Unix(row.CreatedAt, 0).UTC(),
			UpdatedAt: time.Unix(row.UpdatedAt, 0).UTC(),
		})
	}
	// todo: paginate response
	return view, nil
}

// getUserByID returns the user data associated with the given ID.
// If the user does not exist, it returns an error.
func (s *Service) getUserByID(userID domains.ID) (*domains.User_t, error) {
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
