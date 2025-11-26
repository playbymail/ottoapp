// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package domains

import (
	"time"
)

type SessionId string

type Session struct {
	Id   SessionId
	Csrf string
	User User_t
	//Data           map[string]any
	CreatedAt      time.Time // always UTC
	ExpiresAt      time.Time // always UTC
	LastActivityAt time.Time // always UTC
}

const (
	ErrInvalidGcInterval = Error("invalid gc interval")
	ErrSessionExpired    = Error("expired session")
	ErrSessionInvalid    = Error("invalid session")
	ErrSessionIdInvalid  = Error("invalid session id")
	ErrInvalidTtl        = Error("invalid ttl")
)

type SessionView struct {
	UserID ID     `json:"userId"`
	Handle string `json:"handle"`
	CSRF   string `json:"csrf"`
	Roles  struct {
		AccessAdminRoutes bool `json:"accessAdminRoutes"`
		AccessGMRoutes    bool `json:"accessGMRoutes"`
		AccessUserRoutes  bool `json:"accessUserRoutes"`
		EditHandle        bool `json:"editHandle"`
	} `json:"roles"`
}
