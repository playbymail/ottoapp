// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package domains

import "time"

const (
	ErrInvalidEmail    = Error("invalid email")
	ErrInvalidHandle   = Error("invalid handle")
	ErrInvalidTimezone = Error("invalid timezone")
	ErrInvalidUserId   = Error("invalid user id")
)

// User_t is the type for a user.
type User_t struct {
	ID ID // unique identifier

	Username string // must be unique and lowercase
	Email    string // must be unique and lowercase
	Locale   UserLocale_t

	Roles map[string]bool

	Created time.Time // always UTC
	Updated time.Time // always UTC
}

// 	LastLogin time.Time // always UTC, time.Zero if never logged in

type UserLocale_t struct {
	DateFormat string // date format, e.g. "2006-01-02"
	Timezone   UserTimezone_t
}

type UserTimezone_t struct {
	Location *time.Location // timezone, e.g. "America/New_York", should default to UTC
}
