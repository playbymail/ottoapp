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
	q := s.db.Queries()
	ctx := s.db.Context()
	updatedAt := time.Now().UTC().Unix()
	switch roleID {
	case "active":
		return q.UpdateUserActiveRole(ctx, sqlc.UpdateUserActiveRoleParams{UserID: int64(userID), HasRole: true, UpdatedAt: updatedAt})
	case "admin":
		return q.UpdateUserAdminRole(ctx, sqlc.UpdateUserAdminRoleParams{UserID: int64(userID), HasRole: true, UpdatedAt: updatedAt})
	case "gm":
		return q.UpdateUserGMRole(ctx, sqlc.UpdateUserGMRoleParams{UserID: int64(userID), HasRole: true, UpdatedAt: updatedAt})
	case "guest":
		return q.UpdateUserGuestRole(ctx, sqlc.UpdateUserGuestRoleParams{UserID: int64(userID), HasRole: true, UpdatedAt: updatedAt})
	case "player":
		return q.UpdateUserPlayerRole(ctx, sqlc.UpdateUserPlayerRoleParams{UserID: int64(userID), HasRole: true, UpdatedAt: updatedAt})
	case "service":
		return domains.ErrNotAuthorized
	case "sysop":
		return domains.ErrNotAuthorized
	case "user":
		return q.UpdateUserUserRole(ctx, sqlc.UpdateUserUserRoleParams{UserID: int64(userID), HasRole: true, UpdatedAt: updatedAt})
	}
	return domains.ErrInvalidRole
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
	log.Printf("[auth] getActorByEmail(%q)", email)
	user, err := s.db.Queries().ReadUserByEmail(s.db.Context(), email)
	if err != nil {
		log.Printf("[auth] getActorByEmail(%q): %v", email, err)
		return nil, err
	}
	return s.GetActorById(domains.ID(user.UserID))
}

// GetActorByHandle returns a domain Actor or an error.
func (s *Service) GetActorByHandle(handle string) (*domains.Actor, error) {
	//log.Printf("[auth] getActorByHandle(%q)", handle)
	user, err := s.db.Queries().ReadUserByHandle(s.db.Context(), handle)
	if err != nil {
		//log.Printf("[auth] getActorByHandle(%q): %v", handle, err)
		return nil, err
	}
	return s.GetActorById(domains.ID(user.UserID))
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
		return &domains.Actor{ID: SysopId, Roles: domains.Roles{Sysop: true}}, nil
	}
	user, err := s.db.Queries().ReadUserByUserId(s.db.Context(), int64(actorId))
	if err != nil {
		return nil, err
	}
	return &domains.Actor{
		ID: actorId,
		Roles: domains.Roles{
			Active:  user.IsActive,
			Admin:   user.IsAdmin,
			Gm:      user.IsGm,
			Guest:   user.IsGuest,
			Player:  user.IsPlayer,
			Service: user.IsService,
			Sysop:   user.IsSysop,
			User:    user.IsUser,
		},
	}, nil
}

func (s *Service) GetSessionData(r *http.Request) (*domains.SessionView, error) {
	//log.Printf("[authz] GetSessionData\n")
	sessionCookie, err := r.Cookie("sid")
	if err != nil || sessionCookie.Value == "" {
		return &domains.SessionView{}, ErrMissingSessionCookie
	}
	row, err := s.db.Queries().ReadSessionData(s.db.Context(), sqlc.ReadSessionDataParams{
		SessionID: sessionCookie.Value,
		ExpiresAt: time.Now().UTC().Unix(),
	})
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Printf("[authz] GetSessionData: %v\n", err)
			return &domains.SessionView{}, errors.Join(domains.ErrDatabaseError, err)
		}
		log.Printf("[authz] GetSessionData: no data found\n")
		return &domains.SessionView{}, ErrNoSessionData
	}
	return &domains.SessionView{
		CSRF:   row.Csrf,
		UserID: domains.ID(row.UserID),
		Handle: row.Handle,
		Roles: struct {
			AccessAdminRoutes bool `json:"accessAdminRoutes"`
			AccessGMRoutes    bool `json:"accessGMRoutes"`
			AccessUserRoutes  bool `json:"accessUserRoutes"`
			EditHandle        bool `json:"editHandle"`
		}{
			AccessAdminRoutes: s.CanAccessAdminRoutes(row.IsAdmin),
			AccessGMRoutes:    s.CanAccessGMRoutes(row.IsAdmin, row.IsGm),
			AccessUserRoutes:  s.CanAccessUserRoutes(row.IsAdmin, row.IsUser),
			EditHandle:        row.IsAdmin,
		},
	}, nil
}

func (s *Service) GetUserRoles(userID domains.ID) (domains.Roles, error) {
	user, err := s.db.Queries().ReadUserByUserId(s.db.Context(), int64(userID))
	if err != nil {
		return domains.Roles{}, err
	}
	return domains.Roles{
		Active:  user.IsActive,
		Admin:   user.IsAdmin,
		Gm:      user.IsGm,
		Guest:   user.IsGuest,
		Player:  user.IsPlayer,
		Service: user.IsService,
		Sysop:   user.IsSysop,
		User:    user.IsUser,
	}, nil
}

// RemoveRole removes a role from a user.
func (s *Service) RemoveRole(userID domains.ID, roleID string) error {
	q := s.db.Queries()
	ctx := s.db.Context()
	updatedAt := time.Now().UTC().Unix()
	switch roleID {
	case "active":
		return q.UpdateUserActiveRole(ctx, sqlc.UpdateUserActiveRoleParams{UserID: int64(userID), HasRole: false, UpdatedAt: updatedAt})
	case "admin":
		return q.UpdateUserAdminRole(ctx, sqlc.UpdateUserAdminRoleParams{UserID: int64(userID), HasRole: false, UpdatedAt: updatedAt})
	case "gm":
		return q.UpdateUserGMRole(ctx, sqlc.UpdateUserGMRoleParams{UserID: int64(userID), HasRole: false, UpdatedAt: updatedAt})
	case "guest":
		return q.UpdateUserGuestRole(ctx, sqlc.UpdateUserGuestRoleParams{UserID: int64(userID), HasRole: false, UpdatedAt: updatedAt})
	case "player":
		return q.UpdateUserPlayerRole(ctx, sqlc.UpdateUserPlayerRoleParams{UserID: int64(userID), HasRole: false, UpdatedAt: updatedAt})
	case "service":
		return domains.ErrNotAuthorized
	case "sysop":
		return domains.ErrNotAuthorized
	case "user":
		return q.UpdateUserUserRole(ctx, sqlc.UpdateUserUserRoleParams{UserID: int64(userID), HasRole: false, UpdatedAt: updatedAt})
	}
	return domains.ErrInvalidRole
}
