// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package iana

import (
	"net/http"
	"strings"
	"time"

	"github.com/playbymail/ottoapp/backend/restapi"
)

// TimezoneView is the JSON:API view for a timezone
type TimezoneView struct {
	ID            string    `jsonapi:"primary,timezone"`        // singular when sending a payload
	CanonicalName string    `jsonapi:"attr,canonical-name"`     // canonical IANA name for the timezone
	CreatedAt     time.Time `jsonapi:"attr,created-at,iso8601"` // is this needed for a read-only service?
	UpdatedAt     time.Time `jsonapi:"attr,updated-at,iso8601"` // is this needed for a read-only service?
}

// HandleGetTimezones returns the list of available timezones.
// If the "active" query parameter is set to "true", it returns only
// timezones that are actively used in the database.
// Returns array of {id, label} objects.
//
// GET /api/timezones
func (s *Service) HandleGetTimezones(quiet, verbose, debug bool) http.HandlerFunc {
	//var list []*TimezoneItem
	//var err error

	var allTimezonesView []*TimezoneView
	for _, v := range s.Names() {
		allTimezonesView = append(allTimezonesView, &TimezoneView{
			ID:            strings.ToLower(v.Label),
			CanonicalName: v.Location.String(),
			CreatedAt:     time.Now().UTC(),
			UpdatedAt:     time.Now().UTC(),
		})
	}

	return func(w http.ResponseWriter, r *http.Request) {
		timezonesView := allTimezonesView

		// Check if the "active" query parameter is set
		if r.URL.Query().Get("active") == "true" {
			activeList, err := s.Active(quiet, verbose, debug)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			var activeTimezonesView []*TimezoneView
			for _, v := range activeList {
				activeTimezonesView = append(activeTimezonesView, &TimezoneView{
					ID:            strings.ToLower(v.Label),
					CanonicalName: v.Location.String(),
					CreatedAt:     time.Now().UTC(),
					UpdatedAt:     time.Now().UTC(),
				})
			}
			timezonesView = activeTimezonesView
		}

		restapi.WriteJsonApiData(w, http.StatusCreated, timezonesView)
		return
	}
}
