// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/playbymail/ottoapp/backend/domains"
	"github.com/playbymail/ottoapp/backend/stores/sqlite/sqlc"
)

// Init initializes a new store. It enables WAL (write-ahead logging) for concurrency
// and verifies that the sqlite library supports foreign keys.
//
// Returns an error if the path already exists or there are errors initializing it.
func Init(ctx context.Context, path string, overwrite bool) error {
	started := time.Now()

	if sb, err := os.Stat(path); err != nil || !sb.IsDir() {
		return errors.Join(fmt.Errorf("invalid db path"), err)
	}

	name := filepath.Join(path, "ottoapp.db")
	if _, err := os.Stat(name); err == nil {
		if !overwrite {
			return domains.ErrDatabaseExists
		}
		if err := os.Remove(name); err != nil {
			return errors.Join(fmt.Errorf("overwrite failed"), err)
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return errors.Join(fmt.Errorf("invalid db-name"), err)
	}

	log.Printf("[sqldb] %s: initializing...\n", name)

	// Apply PRAGMA's per-connection via DSN so the pool always has them.
	// modernc.org/sqlite supports repeated _pragma=... parameters.
	dsn := fmt.Sprintf(
		"file:%s?_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)&_pragma=foreign_keys(ON)&_pragma=busy_timeout(5000)",
		name,
	)

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return err
	}
	defer func() {
		_ = db.Close()
	}()

	// Optional: size your pool; WAL supports multiple readers + 1 writer.
	// db.SetMaxOpenConns(10)
	// db.SetMaxIdleConns(10)

	// Sanity checks: ensure WAL and FK actually stuck for this connection
	var jm string
	if err := db.QueryRow(`PRAGMA journal_mode;`).Scan(&jm); err != nil {
		return errors.Join(fmt.Errorf("check journal_mode failed"), err)
	} else if jm != "wal" && jm != "WAL" {
		return fmt.Errorf("expected WAL journal_mode, got %q", jm)
	}
	var fk int
	if err := db.QueryRow(`PRAGMA foreign_keys;`).Scan(&fk); err != nil {
		return errors.Join(fmt.Errorf("check foreign_keys failed"), err)
	} else if fk != 1 {
		return fmt.Errorf("foreign_keys pragma not enabled (got %d)", fk)
	}

	log.Printf("[sqldb] init: migrating up\n")
	n, err := migrateUp(ctx, db, migrationsFS, true)
	if err != nil {
		log.Printf("[sqldb] init: migration %v\n", err)
		return err
	}
	log.Printf("[sqldb] init: migration: applied %d\n", n)

	log.Printf("[sqldb] init: completed in %v\n", time.Since(started))

	return nil
}

// OpenTempDB returns an initialized in-memory store for testing.
// It verifies that the sqlite library supports foreign keys.
//
// Returns an error if the path already exists or there are errors initializing it.
func OpenTempDB(ctx context.Context) (*DB, error) {
	started := time.Now()

	log.Printf("[sqldb] initializing temp db...\n")

	// Apply PRAGMA's per-connection via DSN so the pool always has them.
	// modernc.org/sqlite supports repeated _pragma=... parameters.
	dsn := "file::memory:?&cache=shared&_pragma=foreign_keys(ON)"

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = db.Close()
		}
	}()

	// Sanity checks: ensure FK actually stuck for this connection
	var fk int
	if err := db.QueryRow(`PRAGMA foreign_keys;`).Scan(&fk); err != nil {
		return nil, errors.Join(fmt.Errorf("check foreign_keys failed"), err)
	} else if fk != 1 {
		return nil, fmt.Errorf("foreign_keys pragma not enabled (got %d)", fk)
	}

	log.Printf("[sqldb] init: migrating up\n")
	n, err := migrateUp(ctx, db, migrationsFS, true)
	if err != nil {
		log.Printf("[sqldb] init: migration %v\n", err)
		return nil, err
	}
	log.Printf("[sqldb] init: migration: applied %d\n", n)

	log.Printf("[sqldb] init: completed in %v\n", time.Since(started))

	// return the store.
	return &DB{path: ":memory:", name: ":memory:", db: db, ctx: ctx, q: sqlc.New(db)}, nil
}
