// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package domains

import (
	"time"
)

type SessionId string

type Session struct {
	Id             SessionId
	Csrf           string
	User           User_t
	Data           map[string]any
	CreatedAt      time.Time // always UTC
	ExpiresAt      time.Time // always UTC
	LastActivityAt time.Time // always UTC
}

const (
	ErrSessionExpired   = Error("expired session")
	ErrSessionInvalid   = Error("invalid session")
	ErrSessionIdInvalid = Error("invalid session id")
	ErrTtlInvalid       = Error("invalid ttl")
)
