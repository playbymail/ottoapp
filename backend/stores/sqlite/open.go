// Copyright (c) 2024 Michael D Henderson. All rights reserved.

// Package sqlite implements the database layer to a Sqlite database.
package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/playbymail/ottoapp/backend/domains"
	"github.com/playbymail/ottoapp/backend/stores/sqlite/sqlc"
)

// Close closes the database connection.
func (db *DB) Close() error {
	var err error
	if db != nil {
		if db.db != nil {
			err = db.db.Close()
			db.db = nil
		}
	}
	return err
}

// Open opens an existing store. It verifies that WAL is enabled and that foreign
// keys are supported.
//
// Returns an error if the path is not a directory, or if the database does not exist.
// Caller must call Close() when done.
func Open(ctx context.Context, path string, checkVersion, debug bool) (*DB, error) {
	// it is an error if the path does not already exist and is not a directory.
	if sb, err := os.Stat(path); err != nil {
		log.Printf("[sqldb] %q: %s\n", path, err)
		return nil, err
	} else if !sb.IsDir() {
		log.Printf("[sqldb] %q: %s\n", path, err)
		return nil, domains.ErrInvalidPath
	}

	name := filepath.Join(path, "ottoapp.db")

	// it is an error if the database does not already exist and is not a file.
	if sb, err := os.Stat(name); err != nil {
		log.Printf("[sqldb] %q: %s\n", name, err)
		return nil, err
	} else if sb.IsDir() || !sb.Mode().IsRegular() {
		return nil, domains.ErrInvalidPath
	}

	if debug {
		log.Printf("[sqldb] opening %s\n", path)
	}

	// Apply PRAGMA's per-connection via DSN so the pool always has them.
	// modernc.org/sqlite supports repeated _pragma=... parameters.
	dsn := fmt.Sprintf(
		"file:%s?_pragma=foreign_keys(ON)&_pragma=busy_timeout(5000)",
		name,
	)

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	// don't defer a close here since we're returning the connection!

	// Sanity checks: ensure FK actually stuck for this connection
	var fk int
	if err := db.QueryRow(`PRAGMA foreign_keys;`).Scan(&fk); err != nil {
		_ = db.Close()
		log.Printf("[sqldb] %q: check foreign keys %v\n", name, err)
		return nil, errors.Join(fmt.Errorf("check foreign_keys: failed"), err)
	} else if fk != 1 {
		_ = db.Close()
		log.Printf("[sqldb] %q: check foreign keys: disabled\n", name)
		return nil, fmt.Errorf("foreign_keys pragma not enabled (got %d)", fk)
	}
	if checkVersion {
		var actualVersion string
		if err := db.QueryRow(`SELECT value FROM config WHERE key = 'schema.version';`).Scan(&actualVersion); err != nil {
			_ = db.Close()
			log.Printf("[sqldb] %q: check version %v\n", name, err)
			return nil, errors.Join(fmt.Errorf("check version: failed"), err)
		} else if actualVersion != expectedSchemaVersion {
			_ = db.Close()
			log.Printf("[sqldb] %q: version mismatch: want %q: got %q\n", name, expectedSchemaVersion, actualVersion)
			return nil, domains.ErrSchemaVersionMismatch
		}
	}

	// return the store.
	return &DB{path: path, name: name, db: db, ctx: ctx, q: sqlc.New(db)}, nil
}

func executePragma(db *sql.DB, pragma string) error {
	if rslt, err := db.Exec("PRAGMA " + pragma); err != nil {
		return domains.ErrPragmaFailed
	} else if rslt == nil {
		return domains.ErrPragmaReturnedNil
	}
	return nil
}

func enableCheckpointFull(db *sql.DB) error {
	// Most common: wait for writers to pause, then copy frames into the db,
	// reset the WAL so it can be reused; size may remain.
	return executePragma(db, "wal_checkpoint(FULL)")
}

func enableCheckpointPassive(db *sql.DB) error {
	// Non-blocking: do as much as possible without interfering with writers.
	return executePragma(db, "wal_checkpoint(PASSIVE)")
}

func enableCheckpointRestart(db *sql.DB) error {
	// Like FULL, but also resets the WAL file to size 0.
	return executePragma(db, "wal_checkpoint(RESTART)")
}

func enableCheckpointTruncate(db *sql.DB) error {
	// Strongest clean-up: copy, reset *and* truncate the WAL to 0 bytes.
	return executePragma(db, "wal_checkpoint(TRUNCATE)")
}

func enableForeignKeys(db *sql.DB) error {
	return executePragma(db, "foreign_keys = ON")
}

func enableWriteAheadLogging(db *sql.DB) error {
	if err := executePragma(db, "journal_mode = WAL"); err != nil {
		return err
	} else if err = executePragma(db, "wal_autocheckpoint = 1000"); err != nil {
		return err
	} else if err = executePragma(db, "synchronous = NORMAL"); err != nil {
		return err
	}
	return nil
}
