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

const (
	ErrInvalidCredentials = Error("invalid credentials")
	ErrInvalidPassword    = Error("invalid password")
	ErrNoRolesAssigned    = Error("no roles assigned")
)

type Role string
type Roles map[Role]bool
