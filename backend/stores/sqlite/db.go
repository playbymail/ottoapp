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
	// don't forget to update the schema version to match the expected migrations
	schemaVersion = 4
)

type DB struct {
	path string
	name string
	db   *sql.DB
	ctx  context.Context
	q    *sqlc.Queries
}
