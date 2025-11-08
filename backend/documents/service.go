// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package documents implements a service for managing documents.
package documents

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/playbymail/ottoapp/backend/domains"
	"github.com/playbymail/ottoapp/backend/stores/sqlite"
	"github.com/playbymail/ottoapp/backend/stores/sqlite/sqlc"
)

// Service provides document management operations.
type Service struct {
	db *sqlite.DB
}

func New(db *sqlite.DB) *Service {
	return &Service{db: db}
}

func (s *Service) CreateDocument(doc *domains.Document) (domains.ID, error) {
	ctx := s.db.Context()

	// don't trust the caller on important metadata
	contentLength := len(doc.Contents)
	contentsHash, err := hashContents(doc.Contents)
	if err != nil {
		return domains.InvalidID, errors.Join(domains.ErrHashFailed, err)
	}

	// start transaction
	tx, err := s.db.Stdlib().BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	// rollback if we return early; harmless after commit
	defer tx.Rollback()

	userId, createdById := int64(1), int64(1) // sysop

	qtx := s.db.Queries().WithTx(tx)

	now := time.Now().UTC().Unix()

	id, err := qtx.CreateDocument(s.db.Context(), sqlc.CreateDocumentParams{
		MimeType:      doc.MimeType,
		ContentsHash:  contentsHash,
		ContentLength: int64(contentLength),
		CreatedAt:     now,
		UpdatedAt:     now,
	})
	if err != nil {
		return domains.InvalidID, err
	}
	err = qtx.CreateDocumentContent(s.db.Context(), sqlc.CreateDocumentContentParams{
		DocumentID: id,
		Contents:   doc.Contents,
		CreatedAt:  now,
		UpdatedAt:  now,
	})
	if err != nil {
		return domains.InvalidID, err
	}
	err = qtx.CreateDocumentAcl(s.db.Context(), sqlc.CreateDocumentAclParams{
		DocumentID:   id,
		UserID:       userId,
		DocumentName: doc.Path,
		CreatedBy:    createdById,
		IsOwner:      true,
		CanRead:      true,
		CanWrite:     true,
		CanDelete:    true,
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	if err != nil {
		return domains.InvalidID, err
	}

	err = tx.Commit()
	if err != nil {
		return domains.InvalidID, err
	}
	return domains.ID(id), nil
}

// LoadDocxFromFS loads the file, creates a Document, and returns the document ID.
func (s *Service) LoadDocxFromFS(path string) (domains.ID, error) {
	// Stat first so we can validate the path and get size/timestamps.
	sb, err := os.Stat(path)
	if err != nil {
		return domains.InvalidID, errors.Join(domains.ErrInvalidPath, err)
	}
	if sb.IsDir() || !sb.Mode().IsRegular() {
		// keep ErrInvalidPath but add context so logs make sense
		return domains.InvalidID, errors.Join(domains.ErrInvalidPath, fmt.Errorf("%s: not a regular file", path))
	}

	// Read the entire file – required for later parsing steps.
	data, err := os.ReadFile(path)
	if err != nil {
		return domains.InvalidID, errors.Join(domains.ErrReadFailed, err)
	}

	doc := &domains.Document{
		Path:     path, // not tainted since we trust our admin
		MimeType: "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		Contents: data,
	}

	return s.CreateDocument(doc)
}

// LoadDocxFromRequest loads a document from an http.Request.
// It's used by the file upload handlers.
func (s *Service) LoadDocxFromRequest(r *http.Request) (*domains.Document, error) {
	// todo: can we untaint the source path?
	return nil, domains.ErrNotImplemented
}

func (s *Service) UpdateDocx(doc *domains.Document) error {
	return domains.ErrNotImplemented
}

func hashContents(b []byte) (string, error) {
	h := sha256.New()
	if _, err := h.Write(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
