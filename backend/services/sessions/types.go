// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package sessions

import "time"

type User_t struct {
	ID       string          `json:"id"`
	Username string          `json:"username"`
	Roles    map[string]bool `json:"roles,omitempty"`
}

type Session_t struct {
	User   User_t
	CSRF   string
	Expiry time.Time
}
