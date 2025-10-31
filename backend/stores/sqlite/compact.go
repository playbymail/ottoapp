// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package sqlite

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
)

// Compact runs a checkpoint and truncates the WAL files. This reduces
// the file size and sets up the database file to be copied.
func Compact(ctx context.Context, path string) error {
	started := time.Now()

	dbw, err := Open(ctx, path, false)
	if err != nil {
		return err
	} else if dbw == nil || dbw.db == nil {
		panic("assert(dbw && dbw.db)")
	}
	db := dbw.db
	defer func() {
		_ = dbw.Close()
	}()

	log.Printf("[sqldb] compact: verifying WAL mode\n")
	var jm string
	if err := db.QueryRow(`PRAGMA journal_mode;`).Scan(&jm); err != nil {
		return errors.Join(fmt.Errorf("check journal_mode failed"), err)
	} else if jm != "wal" && jm != "WAL" {
		return fmt.Errorf("expected journal_mode=WAL, got %q", jm)
	}

	log.Printf("[sqldb] compact: forcing checkpoint and truncating WAL files\n")
	if _, err := db.Exec(`PRAGMA wal_checkpoint(TRUNCATE);`); err != nil {
		return errors.Join(fmt.Errorf("wal_checkpoint(TRUNCATE) failed"), err)
	}

	// (Optional) Inspect status
	// var busy, logg, checkpointed int
	// _ = db.QueryRow(`PRAGMA wal_checkpoint;`).Scan(&busy, &logg, &checkpointed)

	log.Printf("[sqldb] compact: compacting database\n")
	if _, err := db.Exec(`VACUUM;`); err != nil {
		return errors.Join(fmt.Errorf("VACUUM failed"), err)
	}

	log.Printf("[sqldb] compact: optimizing internal storage\n")
	_, _ = db.Exec(`PRAGMA optimize;`)

	log.Printf("[sqldb] compact: completed in %v\n", time.Since(started))
	return nil
}
