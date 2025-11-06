// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package documents implements a store for documents.
// It supports a testing mode using a temporary Afero file system.
package documents

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/playbymail/ottoapp/backend/domains"
	"github.com/spf13/afero"
)

type Store struct {
	fs afero.Fs
}

// NewStore creates a Store over any afero filesystem.
func NewStore(fs afero.Fs) *Store {
	return &Store{fs: fs}
}

// NewPersistentStore creates a Store rooted at `path` on the OS filesystem.
func NewPersistentStore(path string) (*Store, error) {
	sb, err := os.Stat(path)
	if err != nil {
		return nil, errors.Join(domains.ErrInvalidPath, err)
	}
	if !sb.IsDir() {
		return nil, errors.Join(domains.ErrInvalidPath, fmt.Errorf("not directory: %s", path))
	}

	base := afero.NewBasePathFs(afero.NewOsFs(), path)
	return NewStore(base), nil
}

// NewTemporaryStore creates a temporary in-memory store.
func NewTemporaryStore() (*Store, error) {
	mem := afero.NewMemMapFs()
	if err := mem.MkdirAll(".", 0o755); err != nil {
		return nil, err
	}
	return NewStore(mem), nil
}

// Document represents a document that has been written to the server.
// We use Data in pipeline stages (parsing, cleaning) and writing back
// to the file system.
type Document struct {
	ID         string    // SHA-256 hash of document contents (hex, 64 chars) – used as document ID / dedupe key
	Path       string    // Path on the file system
	SourcePath string    // Client's original name for the source (tainted!)
	Length     int64     // File size in bytes
	CreatedAt  time.Time // File mod time, stored in UTC
	Data       []byte    // Current file contents
}

// some loader functions

// LoadDocumentFromFS loads a document from the file system.
func LoadDocumentFromFS(fs afero.Fs, path string) (*Document, error) {
	// Stat first so we can validate the path and get size/timestamps.
	sb, err := fs.Stat(path)
	if err != nil {
		return nil, errors.Join(domains.ErrInvalidPath, err)
	}
	if sb.IsDir() || !sb.Mode().IsRegular() {
		// keep ErrInvalidPath but add context so logs make sense
		return nil, errors.Join(domains.ErrInvalidPath, fmt.Errorf("not a regular file: %s", path))
	}

	// Read the entire file – required for later parsing steps.
	// todo: future optimizations need to consider streaming the content instead of loading the entire file into memory
	fp, err := fs.Open(path)
	if err != nil {
		return nil, errors.Join(domains.ErrOpenFailed, err)
	}
	data, err := io.ReadAll(fp)
	if err != nil {
		return nil, errors.Join(domains.ErrReadFailed, err)
	}

	// The ID of a Document is the hash of the contents.
	id, err := hashContents(data)
	if err != nil {
		return nil, errors.Join(domains.ErrWriteFailed, err)
	}

	doc := &Document{
		ID:         id,
		Path:       path,
		SourcePath: path, // tainted?
		Length:     sb.Size(),
		CreatedAt:  sb.ModTime().UTC(), // all timestamps must be normalized to UTC
		Data:       data,
	}

	return doc, nil
}

// LoadDocumentFromRequest loads a document from an http.Request.
// It's used by the file upload handlers.
func LoadDocumentFromRequest(r *http.Request) (*Document, error) {
	// todo: can we untaint the source path?
	panic("!implemented")
}

// Save writes the document's data to the store's filesystem.
// Rules:
//   - if d.ID is empty, it is computed from d.Data
//   - if d.Path is empty, d.Path = d.ID
//   - directories are created as needed
//   - on success, d.Length and d.CreatedAt are updated
func (s *Store) Save(d *Document) error {
	if d == nil {
		return errors.Join(domains.ErrWriteFailed, fmt.Errorf("document is nil"))
	}
	if len(d.Data) == 0 {
		return errors.Join(domains.ErrWriteFailed, fmt.Errorf("document has no data"))
	}

	// ensure we have an ID
	if d.ID == "" {
		id, err := hashContents(d.Data)
		if err != nil {
			return errors.Join(domains.ErrHashFailed, err)
		}
		d.ID = id
	}

	// choose path
	destPath := d.Path
	if destPath == "" {
		destPath = d.ID
	}

	// ensure directory exists (handle nested paths later if you need them)
	if dir := filepath.Dir(destPath); dir != "." {
		if err := s.fs.MkdirAll(dir, 0o755); err != nil {
			return errors.Join(domains.ErrWriteFailed, err)
		}
	}

	// write file
	if err := afero.WriteFile(s.fs, destPath, d.Data, 0o644); err != nil {
		return errors.Join(domains.ErrWriteFailed, err)
	}

	// refresh metadata
	fi, err := s.fs.Stat(destPath)
	if err != nil {
		// fallback if stat fails for some reason
		d.Length = int64(len(d.Data))
		if d.CreatedAt.IsZero() {
			d.CreatedAt = time.Now().UTC()
		}
	} else {
		d.Length = fi.Size()
		// Only set CreatedAt if it wasn't set before
		if d.CreatedAt.IsZero() {
			d.CreatedAt = fi.ModTime().UTC()
		}
	}

	d.Path = destPath
	return nil
}

func hashContents(b []byte) (string, error) {
	h := sha256.New()
	if _, err := h.Write(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
