// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package users

//go:generate go run ../../cmd/godel -input handlers.go -struct UserView -output ../../../frontend/app/models/user.js

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/playbymail/ottoapp/backend/domains"
	"github.com/playbymail/ottoapp/backend/iana"
	"github.com/playbymail/ottoapp/backend/services/authn"
	"github.com/playbymail/ottoapp/backend/services/authz"
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
	authnSvc    *authn.Service
	authzSvc    *authz.Service
	ianaSvc     *iana.Service
	sessionsSvc SessionService
}

func New(db *sqlite.DB, authnSvc *authn.Service, authzSvc *authz.Service, ianaSvc *iana.Service) *Service {
	return &Service{
		db:       db,
		authnSvc: authnSvc,
		authzSvc: authzSvc,
		ianaSvc:  ianaSvc,
	}
}

// SetSessionService injects the session service (to avoid circular dependency)
func (s *Service) SetSessionService(sessionsSvc SessionService) {
	s.sessionsSvc = sessionsSvc
}

func (s *Service) CreateUser(handle, email, userName string, timezone *time.Location) (*domains.User_t, error) {
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

	now := time.Now().UTC()
	createdAt, updatedAt := now.Unix(), now.Unix()

	userId, err := s.db.Queries().CreateUser(s.db.Context(), sqlc.CreateUserParams{
		Handle:    handle,
		Username:  userName,
		Email:     email,
		Timezone:  timeZone,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	})
	if err != nil {
		return nil, err
	}

	return s.getUserByID(domains.ID(userId))
}

func (s *Service) ReadUser(actor, target *domains.Actor) (*UserView, error) {
	if !s.authzSvc.CanViewTarget(actor, target) {
		return nil, domains.ErrNotAuthorized
	}
	user, err := s.GetUserByID(target.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domains.ErrNotFound
		}
		return nil, err
	}
	if user == nil {
		return nil, domains.ErrNotFound
	}
	return s.UserView(user, actor, target), nil
}

func (s *Service) ReadUsers(actor *domains.Actor) ([]*UserView, error) {
	if !s.authzSvc.CanListUsers(actor) {
		return nil, domains.ErrNotAuthorized
	}

	// Get all users from database (we'll filter them)
	// TODO: Filtering in Go after GetAllUsers is going to hurt later.
	// Prefer a store method that already enforces visibility/roles.
	rows, err := s.db.Queries().ReadUsers(s.db.Context())
	if err != nil {
		return nil, err
	}
	var views []*UserView
	for _, row := range rows {
		target, err := s.authzSvc.GetActorById(domains.ID(row.UserID))
		if err != nil || !target.IsValid() {
			continue
		}
		if !s.authzSvc.CanViewTarget(actor, target) {
			// skip targets the actor is not allowed to view
			continue
		}
		user, err := s.GetUserByID(target.ID)
		if err != nil {
			return nil, err
		}
		if user == nil {
			return nil, domains.ErrNotExists
		}
		views = append(views, s.UserView(user, actor, target))
	}
	if views == nil {
		return []*UserView{}, nil
	}
	return views, nil
}

func (s *Service) UpdateUser(user *domains.User_t) error {
	if err := domains.ValidateHandle(user.Handle); err != nil {
		return errors.Join(domains.ErrInvalidHandle, domains.ErrBadInput, err)
	}
	if err := domains.ValidateUsername(user.Username); err != nil {
		return errors.Join(domains.ErrInvalidUsername, domains.ErrBadInput, err)
	}
	if err := domains.ValidateEmail(user.Email); err != nil {
		return errors.Join(domains.ErrInvalidEmail, domains.ErrBadInput, err)
	}
	if user.Locale.Timezone.Location == nil {
		return errors.Join(domains.ErrInvalidTimezone, domains.ErrBadInput, fmt.Errorf("timezone is required"))
	}
	timeZone, ok := iana.CanonicalName(user.Locale.Timezone.Location.String())
	if !ok {
		return errors.Join(domains.ErrInvalidTimezone, domains.ErrBadInput, fmt.Errorf("%q: invalid timezone", user.Locale.Timezone.Location.String()))
	}

	err := s.db.Queries().UpdateUserByUserId(s.db.Context(), sqlc.UpdateUserByUserIdParams{
		Email:     user.Email,
		Handle:    user.Handle,
		Timezone:  timeZone,
		Username:  user.Username,
		UpdatedAt: time.Now().UTC().Unix(),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.Join(domains.ErrNotExists, err)
		}
		return errors.Join(domains.ErrDatabaseError, err)
	}

	return nil
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
	userId, err := s.db.Queries().ReadUserIdByEmail(s.db.Context(), email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Join(domains.ErrNotExists, err)
		}
		return nil, errors.Join(domains.ErrDatabaseError, err)
	}
	return s.getUserByID(domains.ID(userId))
}

// GetUserByHandle returns user data associated with the handle.
// Warning: callers expect this to return the same data that would be returned from GetUserByID!
func (s *Service) GetUserByHandle(handle string) (*domains.User_t, error) {
	userId, err := s.db.Queries().ReadUserIdByHandle(s.db.Context(), handle)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Join(domains.ErrNotExists, err)
		}
		return nil, errors.Join(domains.ErrDatabaseError, err)
	}
	return s.getUserByID(domains.ID(userId))
}

func (s *Service) GetUserHandle(userId domains.ID) (string, error) {
	handle, err := s.db.Queries().ReadHandleByUserId(s.db.Context(), int64(userId))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.Join(domains.ErrNotExists, err)
		}
		return "", errors.Join(domains.ErrDatabaseError, err)
	}
	return handle, nil
}

func (s *Service) GetUserIDByEmail(email string) (domains.ID, error) {
	id, err := s.db.Queries().ReadUserIdByEmail(s.db.Context(), email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domains.InvalidID, errors.Join(domains.ErrNotExists, err)
		}
		return domains.InvalidID, errors.Join(domains.ErrDatabaseError, err)
	}
	return domains.ID(id), nil
}

func (s *Service) GetUserIDByHandle(handle string) (domains.ID, error) {
	id, err := s.db.Queries().ReadUserIdByHandle(s.db.Context(), handle)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domains.InvalidID, errors.Join(domains.ErrNotExists, err)
		}
		return domains.InvalidID, errors.Join(domains.ErrDatabaseError, err)
	}
	return domains.ID(id), nil
}

// GetUserByID returns the user data associated with the given ID.
// If the user does not exist, it returns an error.
func (s *Service) GetUserByID(userID domains.ID) (*domains.User_t, error) {
	return s.getUserByID(userID)
}

func (s *Service) ListUsersVisibleToActor(actor *domains.Actor, pageNum, pageSize int) ([]*UserView, error) {
	rows, err := s.db.Queries().ReadUsersVisibleToActor(s.db.Context(), sqlc.ReadUsersVisibleToActorParams{
		ActorID:  actor.ID,
		PageSize: pageSize,
		PageNum:  (pageNum - 1) * pageSize,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Join(domains.ErrNotExists, err)
		}
		return nil, errors.Join(domains.ErrDatabaseError, err)
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
// It does not populate user roles.
func (s *Service) getUserByID(userID domains.ID) (*domains.User_t, error) {
	user, err := s.db.Queries().ReadUserByUserId(s.db.Context(), int64(userID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domains.ErrNotExists
		}
		return nil, err
	}
	loc, err := s.ianaSvc.Location(user.Timezone)
	if err != nil {
		return nil, domains.ErrInvalidTimezone
	}

	return &domains.User_t{
		ID:       userID,
		Username: user.Username,
		Email:    user.Email,
		Handle:   user.Handle,
		Locale: domains.UserLocale_t{
			DateFormat: "2006-01-02",
			Timezone: domains.UserTimezone_t{
				Location: loc,
			},
		},
		Created: time.Unix(user.CreatedAt, 0).UTC(),
		Updated: time.Unix(user.UpdatedAt, 0).UTC(),
	}, nil
}
