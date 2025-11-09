// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package rest

import (
	"encoding/json"
	"net/http"
)

// handleGetTimezones returns the list of available timezones.
// If the "active" query parameter is set to "true", it returns only
// timezones that are actively used in the database.
func (s *Server) handleGetTimezones(w http.ResponseWriter, r *http.Request) {
	var timezones []string
	var err error

	// Check if the "active" query parameter is set
	if r.URL.Query().Get("active") == "true" {
		timezones, err = s.services.ianaSvc.Active()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	} else {
		timezones = s.services.ianaSvc.Names()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(timezones)
}
