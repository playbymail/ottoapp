// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/playbymail/ottoapp/backend/domains"
)

// Document_t represents an uploaded document *as it was first seen*.
// We keep the original bytes because later pipeline stages (parsing,
// cleaning, writing to server storage) need the exact input.
type Document_t struct {
	ID         string    // SHA-256 hash of document contents (hex, 64 chars) – used as document ID / dedupe key
	Path       string    // Path on the server filesystem (where we wrote it to)
	SourcePath string    // Path to the source (where we read it from, untrusted!)
	Length     int64     // Original file size in bytes
	CreatedAt  time.Time // Source file mod time, stored in UTC
	Data       []byte    // Original file contents
}

// NewDocument loads the file, hashes it, and returns a Document_t.
// It reads the whole file into memory on purpose, because later
// stages need the original bytes.
func NewDocument(path string) (*Document_t, error) {
	// Stat first so we can validate the path and get size/timestamps.
	sb, err := os.Stat(path)
	if err != nil {
		return nil, errors.Join(domains.ErrInvalidPath, err)
	}
	if sb.IsDir() || !sb.Mode().IsRegular() {
		// keep ErrInvalidPath but add context so logs make sense
		return nil, errors.Join(domains.ErrInvalidPath, fmt.Errorf("not a regular file: %s", path))
	}

	// Read the entire file – required for later parsing steps.
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Join(domains.ErrReadFailed, err)
	}

	// Hash the same bytes we're returning in Data so there's no chance
	// of mismatch between what we store and what we identified.
	h := sha256.New()
	n, err := h.Write(data)
	if err != nil {
		return nil, errors.Join(domains.ErrWriteFailed, err)
	}
	// Sanity check: number of bytes seen by the hasher should match file size.
	if int64(n) != sb.Size() {
		return nil, domains.ErrHashFailed
	}

	doc := &Document_t{
		ID:         hex.EncodeToString(h.Sum(nil)), // stable 64-char hex
		SourcePath: path,
		Length:     sb.Size(),
		CreatedAt:  sb.ModTime().UTC(), // all timestamps must be normalized to UTC
		Data:       data,
	}

	return doc, nil
}
