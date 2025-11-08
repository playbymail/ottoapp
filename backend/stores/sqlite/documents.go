// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package sqlite

import (
	"time"

	"github.com/playbymail/ottoapp/backend/domains"
	"github.com/playbymail/ottoapp/backend/stores/sqlite/sqlc"
)

func (db *DB) InsertDocument(doc *domains.Document) (domains.ID, error) {
	// start transaction
	tx, err := db.db.BeginTx(db.ctx, nil)
	if err != nil {
		return 0, err
	}
	// rollback if we return early; harmless after commit
	defer tx.Rollback()

	userId, createdById := int64(1), int64(1) // sysop

	qtx := db.q.WithTx(tx)

	now := time.Now().UTC().Unix()

	id, err := qtx.CreateDocument(db.ctx, sqlc.CreateDocumentParams{
		MimeType:      doc.MimeType,
		ContentsHash:  doc.ContentsHash,
		ContentLength: doc.ContentLength,
		CreatedAt:     now,
		UpdatedAt:     now,
	})
	if err != nil {
		return 0, err
	}
	err = qtx.CreateDocumentContent(db.ctx, sqlc.CreateDocumentContentParams{
		DocumentID: id,
		Contents:   doc.Contents,
		CreatedAt:  now,
		UpdatedAt:  now,
	})
	if err != nil {
		return 0, err
	}
	err = qtx.CreateDocumentAcl(db.ctx, sqlc.CreateDocumentAclParams{
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
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return domains.ID(id), err
	}
	return domains.ID(id), nil
}
