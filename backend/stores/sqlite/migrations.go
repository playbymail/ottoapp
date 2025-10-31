// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package sqlite

import (
	"context"
	"database/sql"
	"embed"
	_ "embed"
	"fmt"
	"io/fs"
	"log"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

//go:embed migrations/**.sql
var migrationsFS embed.FS

func MigrateUp(ctx context.Context, path string, isInitializing bool) error {
	started := time.Now()
	wdb, err := Open(ctx, path, false)
	if err != nil {
		return err
	} else if wdb == nil || wdb.db == nil {
		panic("assert(wdb && wdb.db)")
	}
	defer func() {
		_ = wdb.Close()
	}()

	n, err := migrateUp(ctx, wdb.db, migrationsFS, isInitializing)
	log.Printf("[sqldb] migrations: completed %d migrations in %v\n", n, time.Since(started))
	if err != nil {
		return err
	}
	return nil
}

// migrateUp applies all missing migrations found in the root of migrationsFS.
// Returns how many were applied.
func migrateUp(ctx context.Context, db *sql.DB, migrationsFS fs.FS, isInitializing bool) (int, error) {
	// 1) Get already-applied migration names into a set.
	applied, err := fetchAppliedMigrations(ctx, db)
	if err != nil {
		if isInitializing && isNoSuchTable(err) {
			// this is expected during initialization
			applied = map[string]bool{}
		} else {
			return 0, fmt.Errorf("fetch applied: %w", err)
		}
	}

	// 2) List *.sql files that match TIMESTAMP_name.sql; sort by filename.
	candidates, err := listMigrationFiles(migrationsFS)
	if err != nil {
		return 0, fmt.Errorf("list migrations: %w", err)
	}

	appliedCount := 0
	for _, fname := range candidates {
		if applied[fname] {
			log.Printf("[sqldb] migrations: skipping %q\n", fname)
			continue // already applied
		}

		// Read SQL from FS.
		sqlBytes, err := fs.ReadFile(migrationsFS, filepath.Join("migrations", fname))
		if err != nil {
			log.Printf("[sqldb] migrations: error %q: %v\n", fname, err)
			return appliedCount, fmt.Errorf("read %s: %w", fname, err)
		}
		sqlText := string(sqlBytes)

		// 3) Run each migration in its own transaction.
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			log.Printf("[sqldb] migrations: error %q: %v\n", fname, err)
			return appliedCount, fmt.Errorf("begin tx for %s: %w", fname, err)
		}

		if _, err := tx.ExecContext(ctx, sqlText); err != nil {
			_ = tx.Rollback()
			log.Printf("[sqldb] migrations: error %q: %v\n", fname, err)
			return appliedCount, fmt.Errorf("exec %s: %w", fname, err)
		}

		// Record success.
		appliedAt := time.Now().UTC().Format(time.RFC3339Nano)
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO schema_migrations (name, applied_at) VALUES (?, ?)`, fname, appliedAt,
		); err != nil {
			_ = tx.Rollback()
			log.Printf("[sqldb] migrations: error %q: %v\n", fname, err)
			return appliedCount, fmt.Errorf("record %s: %w", fname, err)
		}

		if err := tx.Commit(); err != nil {
			log.Printf("[sqldb] migrations: error %q: %v\n", fname, err)
			return appliedCount, fmt.Errorf("commit %s: %w", fname, err)
		}

		log.Printf("[sqldb] migrations: applied %q\n", fname)
		appliedCount++
	}

	return appliedCount, nil
}

// fetchAppliedMigrations returns the list of migration scripts that have already
// been applied to the database.
func fetchAppliedMigrations(ctx context.Context, db *sql.DB) (map[string]bool, error) {
	rows, err := db.QueryContext(ctx, `SELECT name FROM schema_migrations`)
	if err != nil {
		return nil, err
	}
	applied := map[string]bool{}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			_ = rows.Close()
			return nil, err
		}
		applied[name] = true
	}
	_ = rows.Close()
	return applied, rows.Err()
}

// listMigrationFiles returns the migration scripts in the root of the
// filesystem in the order they should be applied.
func listMigrationFiles(migrationsFS fs.FS) ([]string, error) {
	// timestamp looks like YYYYMMDD_HHMM
	var migRe = regexp.MustCompile(`^\d{8}_\d{4}_.+\.sql$`)

	entries, err := fs.ReadDir(migrationsFS, "migrations")
	if err != nil {
		return nil, err
	}
	var files []string
	for _, e := range entries {
		log.Printf("[sqldb] lmf %-45s: %8v %q\n", e.Name(), e.IsDir(), path.Ext(e.Name()))
		if e.IsDir() {
			continue
		} else if name := e.Name(); path.Ext(name) != ".sql" {
			continue
		} else if !migRe.MatchString(name) {
			// Strict: reject unexpected names so problems are obvious.
			log.Printf("[sqldb] migrate: invalid filename %q\n", name)
			continue
		} else {
			files = append(files, name)
		}
	}
	sort.Strings(files) // lexicographic sort respects the leading timestamp
	return files, nil
}

func isNoSuchTable(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), ": no such table:")
}
