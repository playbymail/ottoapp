// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package sqlite implements the Sqlite database store.
package sqlite

//go:generate sqlc generate

import (
	"context"
	"database/sql"
	_ "embed"

	"github.com/playbymail/ottoapp/backend/stores/sqlite/sqlc"
	_ "modernc.org/sqlite"
)

const (
	// the version of the database this application expects
	expectedSchemaVersion = "20251106_2132"
)

type DB struct {
	path string
	name string // :memory: for a temporary database
	db   *sql.DB
	ctx  context.Context
	q    *sqlc.Queries
}
