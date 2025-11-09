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

type TimezoneItem struct {
	ID    int
	Label string
}

// Service provides timezone operations.
type Service struct {
	sync.Mutex
	db      *sqlite.DB
	names   []TimezoneItem
	regions Zones
}

func New(db *sqlite.DB) (*Service, error) {
	s := &Service{
		db: db,
	}
	var names []string
	// be kind and sort the names
	for _, name := range canonicalNames {
		if strings.HasPrefix(name, "Etc/") {
			continue
		}
		names = append(names, name)
	}
	slices.Sort(names)
	for n, name := range names {
		s.names = append(s.names, TimezoneItem{ID: n + 1, Label: name})
	}
	regions, err := loadZones(names)
	if err != nil {
		return nil, err
	}
	s.regions = regions
	return s, nil
}

func (s *Service) Active() ([]TimezoneItem, error) {
	// assumes that the query sorts the names
	names, err := s.db.Queries().GetActiveTimezones(s.db.Context())
	if err != nil {
		return nil, err
	}
	var list []TimezoneItem
	for n, name := range names {
		list = append(list, TimezoneItem{ID: n + 1, Label: name})
	}

	return list, nil
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

func (s *Service) Names() []TimezoneItem {
	return s.names
}

func (s *Service) Regions() Zones {
	return s.regions
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
