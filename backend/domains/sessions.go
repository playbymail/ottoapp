// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package domains

import "time"

type Session_t struct {
	Id             string
	UserId         ID // id of the user that owns the session
	Data           map[string]any
	CreatedAt      time.Time // always UTC
	LastActivityAt time.Time // always UTC
}
