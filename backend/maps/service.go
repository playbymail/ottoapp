// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package maps

import (
	"errors"
	"os"
	"time"

	"github.com/playbymail/ottoapp/backend/domains"
	"github.com/playbymail/ottoapp/backend/services/authn"
)

// Service provides authentication and authorization operations.
type Service struct {
	authnSvc *authn.Service
	path     string // path to root of data
}

func New(authnSvc *authn.Service, path string) (*Service, error) {
	// verify that path is a valid location
	if sb, err := os.Stat(path); err != nil {
		return nil, errors.Join(domains.ErrInvalidPath, err)
	} else if !sb.IsDir() {
		return nil, errors.Join(domains.ErrInvalidPath, domains.ErrNotDirectory)
	}
	return &Service{authnSvc: authnSvc, path: path}, nil
}

// MapView is the JSON:API view for a map
type MapView struct {
	ID        string         `jsonapi:"primary,map"` // singular when sending a payload
	Game      string         `jsonapi:"attr,game"`
	Clan      string         `jsonapi:"attr,clan"`
	Turn      domains.TurnID `jsonapi:"attr,turn"`
	CreatedAt time.Time      `jsonapi:"attr,created-at,iso8601"`
	UpdatedAt time.Time      `jsonapi:"attr,updated-at,iso8601"`
}

func (s *Service) userMaps(userID domains.ID, game string) (string, error) {
	// todo: fetch user clan, etc
	panic("!implemented")

	//return filepath.Join(s.path, user.Clan, "data", "output"), nil
}

// ListMaps returns a slice containing all the maps in the user's data directory
func (s *Service) ListMaps(userID domains.ID, game string) ([]*MapView, error) {
	panic("!implemented")
}
