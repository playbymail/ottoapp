// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package rest

import (
	"encoding/json"
	"net/http"

	"github.com/playbymail/ottoapp/backend/iana"
)

type timezoneItem struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
}

// handleGetTimezones returns the list of available timezones.
// If the "active" query parameter is set to "true", it returns only
// timezones that are actively used in the database.
// Returns array of {id, label} objects.
func (s *Server) handleGetTimezones(w http.ResponseWriter, r *http.Request) {
	var list []iana.TimezoneItem
	var err error

	// Check if the "active" query parameter is set
	if r.URL.Query().Get("active") == "true" {
		list, err = s.services.ianaSvc.Active()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	} else {
		list = s.services.ianaSvc.Names()
	}

	// Convert to structured format with id and label
	timezones := make([]timezoneItem, 0, len(list))
	for _, item := range list {
		timezones = append(timezones, timezoneItem{
			ID:    item.ID,
			Label: item.Label,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(timezones)
}

// handleGetTimezonesRegions returns the list of available timezones in a zone/sub-zone hierarchy.
func (s *Server) handleGetTimezonesRegions(w http.ResponseWriter, r *http.Request) {
	type Zone struct {
		Name     string  `json:"name"`
		Location string  `json:"tz,omitempty"`
		SubZones []*Zone `json:"children,omitempty"`
	}
	var root []*Zone
	for _, l1 := range s.services.ianaSvc.Regions() {
		z1 := &Zone{Name: l1.Name}
		root = append(root, z1)
		if l1.Location != nil {
			z1.Location = l1.Location.String()
		} else {
			for _, l2 := range l1.SubZones {
				z2 := &Zone{Name: l2.Name}
				z1.SubZones = append(z1.SubZones, z2)
				if l2.Location != nil {
					z2.Location = l2.Location.String()
				} else {
					for _, l3 := range l2.SubZones {
						z3 := &Zone{Name: l3.Name}
						z2.SubZones = append(z2.SubZones, z3)
						if l3.Location != nil {
							z3.Location = l3.Location.String()
						} else {
							panic("assert(iana.region.depth <= 3)")
						}
					}
				}
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(root)
}
