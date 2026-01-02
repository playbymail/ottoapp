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
func (s *Server) handleGetTimezones(quiet, verbose, debug bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var list []*iana.TimezoneItem
		var err error

		// Check if the "active" query parameter is set
		if r.URL.Query().Get("active") == "true" {
			list, err = s.services.ianaSvc.Active(quiet, verbose, debug)
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
}
