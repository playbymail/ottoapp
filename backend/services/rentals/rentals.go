// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package rentals implements an API to serve data to the
// EmberJS Super Rentals tutorial. Whee!
package rentals

import (
	"encoding/json"
	"net/http"

	rs "github.com/playbymail/ottoapp/backend/stores/rentals"
)

type Service struct {
	store *rs.Store
}

func New() (*Service, error) {
	store, err := rs.New()
	return &Service{
		store: store,
	}, err
}

func (s *Service) IndexHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Data []rs.Rental `json:"data"`
	}
	payload.Data, _ = s.store.FetchRentals()
	// Set the Content-Type header to indicate JSON data
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(payload)
	return
}

func (s *Service) IdHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var payload struct {
		Data rs.Rental `json:"data"`
	}
	var ok bool
	payload.Data, ok = s.store.FetchRental(id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	// Set the Content-Type header to indicate JSON data
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(payload)
	return
}
