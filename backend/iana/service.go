// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package iana implements a service for validating and normalizing timezone names.
//
// The `gentz` command should be used to update the list of timezones in the `normalize.go` file.
package iana

//go:generate go run ../../cmd/godel -input handlers.go -struct TimezoneView -output ../../frontend/app/models/timezone.js

import (
	"log"
	"sort"
	"strings"
	"time"

	"github.com/playbymail/ottoapp/backend/domains"
	"github.com/playbymail/ottoapp/backend/stores/sqlite"
)

type TimezoneItem struct {
	ID       int
	Label    string
	Location *time.Location
}

// Service provides timezone operations.
type Service struct {
	db *sqlite.DB
	// canonicalNames is the map of canonical names that Go recognizes
	canonicalNames map[string]*TimezoneItem
	// tzList is the sorted list of canonical names that Go recognizes
	tzList []*TimezoneItem
}

func New(db *sqlite.DB, quiet, verbose, debug bool) (*Service, error) {
	s := &Service{
		db:             db,
		canonicalNames: map[string]*TimezoneItem{},
		tzList:         []*TimezoneItem{},
	}
	for _, cn := range canonicalNames {
		if strings.HasPrefix(cn, "Etc/") {
			continue
		}
		loc, err := time.LoadLocation(cn)
		if err != nil {
			if debug {
				log.Printf("[iana] tz %-55s: location not found\n", cn)
			}
			continue
		}
		item := &TimezoneItem{
			ID:       len(s.tzList) + 1,
			Label:    cn,
			Location: loc,
		}
		s.canonicalNames[cn] = item
		s.tzList = append(s.tzList, item)
	}
	sort.Slice(s.tzList, func(i, j int) bool {
		return s.tzList[i].Label < s.tzList[j].Label
	})
	for i, n := range s.tzList {
		n.ID = i + 1
	}
	return s, nil
}

func (s *Service) Active(quiet, verbose, debug bool) ([]*TimezoneItem, error) {
	names, err := s.db.Queries().GetActiveTimezones(s.db.Context())
	if err != nil {
		return nil, err
	}
	var list []*TimezoneItem
	for _, name := range names {
		cn, ok := Normalize(name)
		if !ok {
			if !quiet {
				log.Printf("[iana] tz %-55s: name not canonical\n", name)
			}
			continue
		}
		item, ok := s.canonicalNames[cn]
		if !ok {
			if !quiet {
				log.Printf("[iana] tz %-55s: location not canonical\n", cn)
			}
			continue
		}
		list = append(list, item)
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

func (s *Service) Names() []*TimezoneItem {
	return s.tzList
}
