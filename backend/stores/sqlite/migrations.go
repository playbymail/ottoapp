// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package sqlite

import (
	"context"
	"database/sql"
	"embed"
	_ "embed"
	"errors"
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

type MigrationStatus struct {
	Id          int
	IsCurrent   bool
	MigrationId string
	AppliedAt   time.Time
	FileName    string
}

func (db *DB) GetDatabaseMigrationStatus() ([]*MigrationStatus, error) {
	status := map[string]*MigrationStatus{}
	// fetch the migrations applied and save their status
	migrationsApplied, err := db.q.GetDatabaseMigrationsApplied(db.ctx)
	if err != nil {
		return nil, err
	}
	for _, row := range migrationsApplied {
		migration := &MigrationStatus{
			Id:          int(row.ID),
			MigrationId: row.MigrationID,
			AppliedAt:   time.Unix(row.AppliedAt, 0).UTC(),
		}
		status[migration.MigrationId] = migration
	}
	// fetch the migration files and save their status
	migrationFiles, err := listMigrationFiles(migrationsFS, false)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("list migrations"), err)
	}
	for _, fileName := range migrationFiles {
		migrationId := fileName[:len("YYYYMMDD_HHMM")]
		migration, ok := status[migrationId]
		if !ok {
			// migration file is missing!
			migration = &MigrationStatus{
				MigrationId: migrationId,
			}
			status[migrationId] = migration
		}
		migration.FileName = fileName
	}
	// fetch the current migration and save its status
	currentId, err := db.q.GetDatabaseVersion(db.ctx)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("get version"), err)
	}
	migration, ok := status[currentId]
	if !ok {
		// current migration doesn't exist!
		migration = &MigrationStatus{
			MigrationId: currentId,
		}
		status[currentId] = migration
	}
	migration.IsCurrent = true
	// convert the map to a list
	var list []*MigrationStatus
	for _, elem := range status {
		list = append(list, elem)
	}
	// sort the list by applied then by id
	sort.Slice(list, func(i, j int) bool {
		if list[i].AppliedAt.Before(list[j].AppliedAt) {
			return true
		} else if list[i].AppliedAt.After(list[j].AppliedAt) {
			return false
		} else if list[i].Id < list[j].Id {
			return true
		} else if list[i].Id > list[j].Id {
			return false
		}
		return list[i].MigrationId < list[j].MigrationId
	})
	return list, nil
}

func (db *DB) GetDatabaseVersion() (string, error) {
	return db.q.GetDatabaseVersion(db.ctx)
}

func MigrateUp(ctx context.Context, path string, isInitializing, quiet, verbose, debug bool) error {
	started := time.Now()
	wdb, err := Open(ctx, path, false, quiet, verbose, debug)
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
	candidates, err := listMigrationFiles(migrationsFS, false)
	if err != nil {
		return 0, fmt.Errorf("list migrations: %w", err)
	}

	appliedCount := 0
	for _, fname := range candidates {
		// the file name looks like "20251029_1540_init.sql".
		// the migration ID is the timestamp part of that.
		migrationId := fname[:len("YYYYMMDD_HHMM")]
		log.Printf("[sqldb] migrations: %q %q\n", migrationId, fname)

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

		// Record the migration.
		appliedAt := time.Now().UTC().Unix()
		log.Printf("[sqldb] migrations: %q %22d\n", migrationId, appliedAt)

		if _, err := tx.ExecContext(ctx,
			`INSERT INTO schema_migrations (migration_id, file_name, created_at, applied_at, updated_at) VALUES (?, ?, ?, ?, ?)`, migrationId, fname, appliedAt, appliedAt, appliedAt,
		); err != nil {
			_ = tx.Rollback()
			log.Printf("[sqldb] migrations: error %q: %v\n", fname, err)
			return appliedCount, fmt.Errorf("record %s: %w", fname, err)
		}
		if _, err := tx.ExecContext(ctx,
			`UPDATE config SET value = ?, updated_at = ? WHERE key = 'schema.version'`, migrationId, appliedAt,
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
	rows, err := db.QueryContext(ctx, `SELECT file_name FROM schema_migrations`)
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
func listMigrationFiles(migrationsFS fs.FS, debug bool) ([]string, error) {
	// timestamp looks like YYYYMMDD_HHMM
	var migRe = regexp.MustCompile(`^\d{8}_\d{4}_.+\.sql$`)

	entries, err := fs.ReadDir(migrationsFS, "migrations")
	if err != nil {
		return nil, err
	}
	var files []string
	for _, e := range entries {
		if debug {
			log.Printf("[sqldb] lmf %-45s: %8v %q\n", e.Name(), e.IsDir(), path.Ext(e.Name()))
		}
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
