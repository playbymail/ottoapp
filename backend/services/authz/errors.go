// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package authz

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrMissingSessionCookie = Error("missing session cookie")
	ErrNoSessionData        = Error("no session data")
)
