// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package versions

//go:generate go run ../../cmd/godel -input service.go -struct VersionView -output ../../frontend/app/models/version.js

import "github.com/maloquacious/semver"

// Service provides authentication and authorization operations.
type Service struct {
	version semver.Version
	short   string
	full    string
}

func New(version semver.Version) *Service {
	return &Service{
		version: version,
		short:   version.Core(),
		full:    version.String(),
	}
}

func (s *Service) Version() VersionView {
	return VersionView{
		ID:    "1",
		Short: s.short,
		Full:  s.full,
	}
}

// VersionView is the JSON:API view for a version
type VersionView struct {
	ID    string `jsonapi:"primary,version"` // singular when sending a payload
	Short string `jsonapi:"attr,short"`
	Full  string `jsonapi:"attr,full"`
}
