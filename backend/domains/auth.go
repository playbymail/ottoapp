// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package domains

// ID is the type for identity. It is unique and immutable.
//
// It is used to identify a user, organization, or other entity.
//
// We assume that the ID is never deleted or reused.
type ID int64

const (
	InvalidID ID = 0
)

// contextKey is the type for context keys to avoid collisions
type contextKey string

const (
	// ContextKeyUserID is the context key for the authenticated user ID
	ContextKeyUserID contextKey = "userID"
)

const (
	ErrInvalidCredentials = Error("invalid credentials")
	ErrInvalidPassword    = Error("invalid password")
	ErrInvalidRole        = Error("invalid role")
	ErrNoRolesAssigned    = Error("no roles assigned")
	ErrNotAuthenticated   = Error("not authenticated")
)

type Role string
type Roles map[Role]bool
