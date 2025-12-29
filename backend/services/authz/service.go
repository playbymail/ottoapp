// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package authz implements an authorization service.
package authz

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/playbymail/ottoapp/backend/domains"
	"github.com/playbymail/ottoapp/backend/stores/sqlite"
	"github.com/playbymail/ottoapp/backend/stores/sqlite/sqlc"
)

// Service provides authorization operations.
type Service struct {
	db *sqlite.DB
}

func New(db *sqlite.DB) *Service {
	return &Service{db: db}
}

const (
	// SysopId is asserted in the database initialization scripts
	// and initial database connection.
	SysopId = domains.ID(1)
)

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

// GetActor extracts the actor from the request context.
// This is a convenience helper that assumes session middleware
// has added the user to the context.
func (s *Service) GetActor(r *http.Request) (*domains.Actor, error) {
	actorId, ok := r.Context().Value(domains.ContextKeyUserID).(domains.ID)
	if !ok || actorId == domains.InvalidID {
		return nil, domains.ErrNotAuthenticated
	}
	return s.GetActorById(actorId)
}

// GetActorByEmail returns a domain Actor or an error.
func (s *Service) GetActorByEmail(email string) (*domains.Actor, error) {
	//log.Printf("[auth] getActorByEmail(%q)", email)
	userId, err := s.db.Queries().ReadUserIdByEmail(s.db.Context(), email)
	if err != nil {
		//log.Printf("[auth] getActorByEmail(%q): %v", email, err)
		return nil, err
	}
	return s.GetActorById(domains.ID(userId))
}

// GetActorByHandle returns a domain Actor or an error.
func (s *Service) GetActorByHandle(handle string) (*domains.Actor, error) {
	//log.Printf("[auth] getActorByHandle(%q)", handle)
	userId, err := s.db.Queries().ReadUserIdByHandle(s.db.Context(), handle)
	if err != nil {
		//log.Printf("[auth] getActorByHandle(%q): %v", handle, err)
		return nil, err
	}
	return s.GetActorById(domains.ID(userId))
}

// GetActorById returns a domain Actor or an error.
// TODO: security considerations from handing out a sysop actor.
//
// Background processes can construct an actor directly without DB:
//
//	var ServiceActor = domains.Actor{
//	   ID:    domains.InvalidID,              // not a user
//	   Service: true,
//	}
func (s *Service) GetActorById(actorId domains.ID) (*domains.Actor, error) {
	if actorId == SysopId {
		return &domains.Actor{ID: SysopId, Sysop: true}, nil
	}
	userRoles, err := s.db.Queries().GetUserRoles(s.db.Context(), int64(actorId))
	if err != nil {
		return nil, err
	}
	actor := domains.Actor{ID: actorId}
	for _, role := range userRoles {
		switch role {
		case "admin":
			actor.Admin = true
		case "service":
			actor.Service = true
		case "user":
			actor.User = true
		}
	}
	return &actor, nil
}

func (s *Service) GetSessionData(r *http.Request) (*domains.SessionView, error) {
	//log.Printf("[authz] GetSessionData\n")
	sessionCookie, err := r.Cookie("sid")
	if err != nil || sessionCookie.Value == "" {
		return &domains.SessionView{}, ErrMissingSessionCookie
	}
	var view domains.SessionView
	rows, err := s.db.Queries().GetSessionData(s.db.Context(), sqlc.GetSessionDataParams{
		SessionID: sessionCookie.Value,
		ExpiresAt: time.Now().UTC().Unix(),
	})
	//log.Printf("[authz] GetSessionData: %d: %v\n", len(rows), err)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("[authz] GetSessionData: no data found\n")
			return &domains.SessionView{}, ErrNoSessionData
		}
		log.Printf("[authz] GetSessionData: %v\n", err)
		return &domains.SessionView{}, errors.Join(domains.ErrDatabaseError, err)
	} else if len(rows) == 0 {
		log.Printf("[authz] GetSessionData: no data found\n")
		return &domains.SessionView{}, ErrNoSessionData
	}
	var isAdmin, isGM, isUser bool
	for _, row := range rows {
		view.CSRF = row.Csrf
		view.UserID = domains.ID(row.UserID)
		view.Handle = row.Handle
		switch row.RoleID {
		case "admin":
			isAdmin = true
		case "gm":
			isGM = true
		case "user":
			isUser = true
		}
	}
	view.Roles.AccessAdminRoutes = s.CanAccessAdminRoutes(isAdmin)
	view.Roles.AccessGMRoutes = s.CanAccessGMRoutes(isAdmin, isGM)
	view.Roles.AccessUserRoutes = s.CanAccessUserRoutes(isAdmin, isUser)
	view.Roles.EditHandle = isAdmin
	//log.Printf("[authz] GetSessionData: %+v\n", view)
	return &view, nil
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
