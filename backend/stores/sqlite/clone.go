// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package sqlite

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Clone creates a working copy of the database for testing.
// It uses VACUUM INTO to create a clean, single-file copy named ottoapp.db.
//
// The output directory must exist and must not already contain ottoapp.db.
// This is a safety measure to prevent accidentally overwriting existing instances.
//
// It's safest if the server is not running during the clone.
func Clone(ctx context.Context, path string, outputPath string, quiet, verbose, debug bool) (string, error) {
	started := time.Now()

	// Validate output path
	if outputPath == "" {
		return "", fmt.Errorf("clone: output path is required")
	}

	// Verify output directory exists
	if info, err := os.Stat(outputPath); err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("clone: output directory does not exist: %s", outputPath)
		}
		return "", fmt.Errorf("clone: cannot access output directory: %w", err)
	} else if !info.IsDir() {
		return "", fmt.Errorf("clone: output path is not a directory: %s", outputPath)
	}

	// Ensure source and destination are different
	sourcePath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("clone: cannot resolve source path: %w", err)
	}
	destPath, err := filepath.Abs(outputPath)
	if err != nil {
		return "", fmt.Errorf("clone: cannot resolve destination path: %w", err)
	}
	if sourcePath == destPath {
		return "", fmt.Errorf("clone: source and destination paths must be different")
	}

	// Check if ottoapp.db already exists in output directory
	clonePath := filepath.Join(outputPath, "ottoapp.db")
	if _, err := os.Stat(clonePath); err == nil {
		return "", fmt.Errorf("clone: destination file already exists: %s (refusing to overwrite)", clonePath)
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("clone: cannot check destination file: %w", err)
	}

	// Open source database
	wdb, err := Open(ctx, path, false, quiet, verbose, debug)
	if err != nil {
		return "", err
	} else if wdb == nil || wdb.db == nil {
		panic("assert(wdb && wdb.db)")
	}
	defer func() {
		_ = wdb.Close()
	}()
	db := wdb.db

	log.Printf("[sqldb] clone: verifying WAL mode\n")
	var jm string
	if err := db.QueryRowContext(ctx, `PRAGMA journal_mode;`).Scan(&jm); err != nil {
		return "", fmt.Errorf("clone: check journal_mode: %w", err)
	}

	// Create temporary file
	tmp := clonePath + ".part"
	_ = os.Remove(tmp) // best-effort cleanup

	log.Printf("[sqldb] clone: creating clone...\n")
	if _, err := db.ExecContext(ctx, `VACUUM INTO ?;`, tmp); err != nil {
		return "", fmt.Errorf("clone: VACUUM INTO: %w", err)
	}

	log.Printf("[sqldb] clone: moving temp to final location\n")
	if err := os.Rename(tmp, clonePath); err != nil {
		_ = os.Remove(tmp)
		return "", fmt.Errorf("clone: atomic rename: %w", err)
	}

	// Quick integrity check
	log.Printf("[sqldb] clone: checking integrity...\n")
	if err := quickIntegrityCheck(ctx, clonePath); err != nil {
		return "", fmt.Errorf("clone: integrity check failed: %w", err)
	}

	log.Printf("[sqldb] clone: completed in %v\n", time.Since(started))
	return clonePath, nil
}
