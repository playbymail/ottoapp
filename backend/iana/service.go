// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package iana implements a service for validating and normalizing timezone names.
//
// The `gentz` command should be used to update the list of timezones in the `normalize.go` file.
package iana

import (
	"log"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/playbymail/ottoapp/backend/domains"
	"github.com/playbymail/ottoapp/backend/stores/sqlite"
)

// Service provides timezone operations.
type Service struct {
	sync.Mutex
	db    *sqlite.DB
	names []string
}

func New(db *sqlite.DB) *Service {
	// be kind and sort the names
	names := make([]string, 0, len(canonicalNames))
	for _, name := range canonicalNames {
		if strings.HasPrefix(name, "Etc/") {
			continue
		}
		names = append(names, name)
	}
	slices.Sort(names)
	return &Service{
		db:    db,
		names: names,
	}
}

func (s *Service) Active() ([]string, error) {
	names, err := s.db.Queries().GetActiveTimezones(s.db.Context())
	if err != nil {
		return nil, err
	}
	return names, nil
}

// Location returns the time.Location associated with the timezone or an error.
func (s *Service) Location(tz string) (*time.Location, error) {
	ntz, ok := Normalize(tz)
	if !ok {
		return nil, domains.ErrInvalidTimezone
	}
	loc, err := time.LoadLocation(ntz)
	if err != nil {
		// this happens when the normalized data is out of sync with
		// the Go distribution that created the binary. we could panic.
		log.Printf("[iana] tz %q: out of sync: %v\n", tz, err)
		return nil, domains.ErrInvalidTimezone
	}
	return loc, nil
}

func (s *Service) Names() []string {
	return s.names
}

func nNormalizeTimeZone(tz string) (loc *time.Location, ok bool) {
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
