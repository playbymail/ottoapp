// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Backup creates a compact, consistent backup into the database path.
// It uses VACUUM INTO, so the result is a single .db file with no sidecars.
//
// It's safest if the server is not running during the backup.
func Backup(ctx context.Context, path string) (string, error) {
	started := time.Now()
	wdb, err := Open(ctx, path, false)
	if err != nil {
		return "", err
	} else if wdb == nil || wdb.db == nil {
		panic("assert(wdb && wdb.db)")
	}
	defer func() {
		_ = wdb.Close()
	}()
	db := wdb.db

	log.Printf("[sqldb] backup: verifying WAL mode\n")
	var jm string
	if err := db.QueryRowContext(ctx, `PRAGMA journal_mode;`).Scan(&jm); err != nil {
		return "", fmt.Errorf("backup: check journal_mode: %w", err)
	} else {
		// Not strictly required, but nice to know
		if jm != "wal" && jm != "WAL" {
			/* OK even if not WAL */
		}
	}

	ts := started.UTC().Format("20060102-150405")
	bkup := filepath.Join(path, fmt.Sprintf("backup-%s.db", ts))
	tmp := bkup + ".part"
	_ = os.Remove(tmp) // best-effort cleanup

	log.Printf("[sqldb] backup: backing up to temp file...\n")
	if _, err := db.ExecContext(ctx, `VACUUM INTO ?;`, tmp); err != nil {
		return "", fmt.Errorf("backup: VACUUM INTO: %w", err)
	}

	log.Printf("[sqldb] backup: moving temp to backup\n")
	if err := os.Rename(tmp, bkup); err != nil {
		_ = os.Remove(tmp)
		return "", fmt.Errorf("backup: atomic rename: %w", err)
	}

	// Optional: quick integrity check on the new file
	log.Printf("[sqldb] backup: checking backupe...\n")
	if err := quickIntegrityCheck(ctx, bkup); err != nil {
		return "", fmt.Errorf("backup: integrity check failed: %w", err)
	}
	log.Printf("[sqldb] backup: completed in %v\n", time.Since(started))

	return bkup, nil
}

// quickIntegrityCheck opens the given DB file read-only and runs PRAGMA integrity_check.
func quickIntegrityCheck(ctx context.Context, dbPath string) error {
	dsn := fmt.Sprintf("file:%s?mode=ro&_pragma=busy_timeout(5000)", dbPath)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return err
	} else if db == nil {
		panic("assert(db)")
	}
	defer func() {
		_ = db.Close()
	}()

	var res string
	if err := db.QueryRowContext(ctx, `PRAGMA integrity_check;`).Scan(&res); err != nil {
		return err
	}
	if res != "ok" {
		return fmt.Errorf("integrity_check returned %q", res)
	}
	return nil
}
