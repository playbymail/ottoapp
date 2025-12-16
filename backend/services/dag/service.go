// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package dag implements services to manage dependencies.
package dag

import (
	"github.com/playbymail/ottoapp/backend/stores/sqlite"
)

// Service provides dependency management operations.
type Service struct {
	db *sqlite.DB
}

func New(db *sqlite.DB) (*Service, error) {
	return &Service{db: db}, nil
}
