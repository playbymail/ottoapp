// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package documents

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"

	"github.com/playbymail/ottoapp/backend/domains"
)

// LoadDocx loads the file, hashes it, and returns a Document.
func LoadDocx(path string) (*domains.Document, error) {
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

	// Hash the same bytes we're returning in Contents so there's no chance
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

	doc := &domains.Document{
		Path:          path, // not tainted since we trust our admin
		MimeType:      "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		ContentLength: sb.Size(),
		Contents:      data,
		ContentsHash:  hex.EncodeToString(h.Sum(nil)), // stable 64-char hex
		CreatedAt:     sb.ModTime().UTC(),             // all timestamps must be normalized to UTC
		UpdatedAt:     sb.ModTime().UTC(),             // all timestamps must be normalized to UTC
	}

	return doc, nil
}
