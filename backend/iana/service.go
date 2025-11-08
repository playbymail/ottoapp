// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package iana

import (
	"log"
	"time"
)

func NormalizeTimeZone(tz string) (loc *time.Location, ok bool) {
	tz, ok = Normalize(tz)
	if !ok {
		return nil, false
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		log.Printf("[iana] internal error: tz %q: %v\n", tz, err)
	}
	return loc, true
}
