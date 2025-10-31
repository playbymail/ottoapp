// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package domains

import "time"

const (
	ErrInvalidEmail  = Error("invalid email")
	ErrInvalidHandle = Error("invalid handle")
)

// User_t is the type for a user.
type User_t struct {
	ID ID // unique identifier

	Handle string // must be unique and lowercase
	Email  string // must be unique and lowercase
	Locale UserLocale_t

	Created time.Time // always UTC
	Updated time.Time // always UTC
}

// 	LastLogin time.Time // always UTC, time.Zero if never logged in

type UserLocale_t struct {
	DateFormat string // date format, e.g. "2006-01-02"
	Timezone   UserTimezone_t
}

type UserRoles_t struct {
	IsActive        bool // true if the user is active
	IsAdministrator bool // true if the user is an administrator
	IsAuthenticated bool // true if the user is authenticated
	IsOperator      bool // true if the user is an operator
	IsUser          bool // true if the user is a user
}

type UserTimezone_t struct {
	Location *time.Location // timezone, e.g. "America/New_York", should default to UTC
}
