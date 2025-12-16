// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package documents implements a service for managing documents.
package documents

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/playbymail/ottoapp/backend/domains"
	"github.com/playbymail/ottoapp/backend/iana"
	"github.com/playbymail/ottoapp/backend/services/authn"
	"github.com/playbymail/ottoapp/backend/services/authz"
	"github.com/playbymail/ottoapp/backend/services/users"
	"github.com/playbymail/ottoapp/backend/stores/sqlite"
	"github.com/playbymail/ottoapp/backend/stores/sqlite/sqlc"
)

// Service provides document management operations.
type Service struct {
	db       *sqlite.DB
	authzSvc *authz.Service
	usersSvc *users.Service
}

func New(db *sqlite.DB, authzSvc *authz.Service, usersSvc *users.Service) (*Service, error) {
	if authzSvc == nil {
		authzSvc = authz.New(db)
	}
	if usersSvc == nil {
		authnSvc := authn.New(db, authzSvc)
		ianaSvc, err := iana.New(db)
		if err != nil {
			return nil, errors.Join(fmt.Errorf("new iana service"), err)
		}
		usersSvc = users.New(db, authnSvc, authzSvc, ianaSvc)
	}
	return &Service{db: db, authzSvc: authzSvc, usersSvc: usersSvc}, nil
}

// CreateDocument creates a document.
// Returns ErrExists if the document name already exists for the clan.
//
// Actor is the user/service creating the document.
// Owner is the clan will own the new document.
func (s *Service) CreateDocument(actor *domains.Actor, owner *domains.Clan, doc *domains.Document, quiet, verbose, debug bool) (domains.ID, error) {
	if debug {
		log.Printf("[documents] CreateDocument(%d, (%q, %d), %q, %q)\n", actor.ID, owner.GameID, owner.UserID, doc.Path, doc.Type)
	}
	if doc.Path != html.EscapeString(doc.Path) {
		return domains.InvalidID, ErrInvalidPath
	}
	if !s.authzSvc.CanCreateDocuments(actor) {
		return domains.InvalidID, domains.ErrNotAuthorized
	}

	// don't trust the caller on important metadata
	contentLength, contentsHash, err := Hash(doc.Contents)
	if err != nil {
		return domains.InvalidID, errors.Join(domains.ErrHashFailed, err)
	}
	var documentType string
	switch doc.Type {
	case domains.TurnReportFile:
		documentType = "docx"
	case domains.TurnReportExtract:
		documentType = "txt"
	case domains.WorldographerMap:
		documentType = "wxx"
	default:
		return domains.InvalidID, fmt.Errorf("%q: unknown type", doc.Type)
	}

	// start transaction
	ctx := s.db.Context()
	tx, err := s.db.Stdlib().BeginTx(ctx, nil)
	if err != nil {
		return domains.InvalidID, err
	}
	defer tx.Rollback() // rollback if we return early; harmless after commit
	qtx := s.db.Queries().WithTx(tx)
	now := time.Now().UTC()
	createdAt, updatedAt := now.Unix(), now.Unix()

	// return ErrExists if the clan already has a document with this name
	_, err = qtx.ReadDocumentByClanAndName(ctx, sqlc.ReadDocumentByClanAndNameParams{
		ClanID:       int64(owner.ClanID),
		DocumentName: doc.Path,
	})
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Printf("[documents] CreateDocument(%d, (%q, %d), %q) %v", actor.ID, owner.GameID, owner.UserID, doc.Path, err)
			return domains.InvalidID, errors.Join(domains.ErrDatabaseError, err)
		}
		// document does not exist; okay to create
	} else {
		// document exists
		return domains.InvalidID, ErrExists
	}

	// create the document
	documentId, err := qtx.CreateDocument(ctx, sqlc.CreateDocumentParams{
		ClanID:       int64(owner.ClanID),
		DocumentName: doc.Path,
		DocumentType: documentType,
		ModifiedAt:   doc.ModifiedAt.Unix(),
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	})
	if err != nil {
		log.Printf("[documents] CreateDocument(%d, (%q, %d), %q) %v", actor.ID, owner.GameID, owner.UserID, doc.Path, err)
		return domains.InvalidID, errors.Join(domains.ErrDatabaseError, err)
	}

	// upload the contents
	err = qtx.CreateDocumentContents(ctx, sqlc.CreateDocumentContentsParams{
		DocumentID:    documentId,
		ContentLength: int64(contentLength),
		ContentsHash:  contentsHash,
		Contents:      doc.Contents,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	})
	if err != nil {
		log.Printf("[documents] CreateDocument(%d, (%q, %d), (%q, %d), %q) %v", actor.ID, owner.GameID, owner.UserID, owner.GameID, owner.UserID, doc.Path, err)
		return domains.InvalidID, errors.Join(domains.ErrDatabaseError, err)
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("[documents] CreateDocument(%d, (%q, %d), %q) %v", actor.ID, owner.GameID, owner.UserID, doc.Path, err)
		return domains.InvalidID, errors.Join(domains.ErrDatabaseError, err)
	}

	if debug {
		log.Printf("[documents] CreateDocument(%d, (%q, %d), %q, %q) %d\n", actor.ID, owner.GameID, owner.UserID, doc.Path, doc.Type, documentId)
	}

	return domains.ID(documentId), nil
}

// ReadDocument will
func (s *Service) ReadDocument(actor *domains.Actor, owner *domains.Clan, documentId domains.ID, quiet, verbose, debug bool) (*DocumentView, error) {
	// start transaction
	if debug {
		log.Printf("[documents] ReadDocument(%d, (%d, %d) %q) %d\n", actor.ID, owner.GameID, owner.ClanID, documentId)
	}
	ctx := s.db.Context()
	tx, err := s.db.Stdlib().BeginTx(ctx, nil)
	if err != nil {
		log.Printf("[documents] ReadDocument(%d, (%q, %d), %d) %v\n", actor.ID, owner.GameID, owner.ClanID, documentId, err)
		return nil, errors.Join(domains.ErrDatabaseError, err)
	}
	defer tx.Rollback() // rollback if we return early; harmless after commit
	qtx := s.db.Queries().WithTx(tx)

	// cache handles by user_id to avoid repeated database calls
	handles := map[domains.ID]string{}
	for _, id := range []domains.ID{actor.ID, owner.UserID} {
		if _, ok := handles[id]; ok {
			continue
		}
		handle, err := qtx.ReadHandleByUserId(ctx, int64(id))
		if err != nil {
			log.Printf("[documents] ReadDocument(%d, (%q, %d), %d) %v\n", actor.ID, owner.GameID, owner.ClanID, documentId, err)
			return nil, errors.Join(fmt.Errorf("GetUserHandle(%d)", id), err)
		}
		handles[id] = handle
	}

	d, err := qtx.ReadDocumentContentsByIdAuthorized(s.db.Context(), sqlc.ReadDocumentContentsByIdAuthorizedParams{
		DocumentID: int64(documentId),
		ClanID:     int64(owner.ClanID),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domains.ErrNotAuthorized
		}
		log.Printf("[documents] ReadDocument(%d, (%q, %d), %d) %v\n", actor.ID, owner.GameID, owner.ClanID, documentId, err)
		return nil, errors.Join(domains.ErrDatabaseError, err)
	}

	view := &DocumentView{
		ID:           fmt.Sprintf("%d", d.DocumentID),
		OwnerHandle:  handles[owner.UserID],
		UserHandle:   handles[owner.UserID],
		GameId:       fmt.Sprintf("%s", d.GameID),
		ClanNo:       fmt.Sprintf("%04d", owner.ClanNo),
		DocumentName: d.DocumentName,
		DocumentType: d.DocumentType,
		ModifiedAt:   time.Unix(d.ModifiedAt, 0).UTC(),
		CreatedAt:    time.Unix(d.CreatedAt, 0).UTC(),
		UpdatedAt:    time.Unix(d.UpdatedAt, 0).UTC(),
	}

	return view, nil
}

// ReadDocumentContents will
func (s *Service) ReadDocumentContents(actor *domains.Actor, owner *domains.Clan, documentId domains.ID, quiet, verbose, debug bool) (*domains.Document, error) {
	// start transaction
	if debug {
		log.Printf("[documents] ReadDocumentContents(%d, (%d, %d) %q) %d\n", actor.ID, owner.GameID, owner.ClanID, documentId)
	}
	ctx := s.db.Context()
	tx, err := s.db.Stdlib().BeginTx(ctx, nil)
	if err != nil {
		log.Printf("[documents] ReadDocumentContents(%d, (%q, %d), %d) %v\n", actor.ID, owner.GameID, owner.ClanID, documentId, err)
		return nil, errors.Join(domains.ErrDatabaseError, err)
	}
	defer tx.Rollback() // rollback if we return early; harmless after commit
	qtx := s.db.Queries().WithTx(tx)

	// cache handles by user_id to avoid repeated database calls
	handles := map[domains.ID]string{}
	for _, id := range []domains.ID{actor.ID, owner.UserID} {
		if _, ok := handles[id]; ok {
			continue
		}
		handle, err := qtx.ReadHandleByUserId(ctx, int64(id))
		if err != nil {
			log.Printf("[documents] ReadDocumentContents(%d, (%q, %d), %d) %v\n", actor.ID, owner.GameID, owner.ClanID, documentId, err)
			return nil, errors.Join(domains.ErrDatabaseError, err)
		}
		handles[id] = handle
	}

	d, err := qtx.ReadDocumentContentsByIdAuthorized(s.db.Context(), sqlc.ReadDocumentContentsByIdAuthorizedParams{
		DocumentID: int64(documentId),
		ClanID:     int64(owner.ClanID),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domains.ErrNotAuthorized
		}
		log.Printf("[documents] ReadDocumentContents(%d, (%q, %d), %d) %v\n", actor.ID, owner.GameID, owner.ClanID, documentId, err)
		return nil, errors.Join(domains.ErrDatabaseError, err)
	}
	var documentType domains.DocumentType
	switch d.DocumentType {
	case "turn-report-file":
		documentType = domains.TurnReportFile
	case "turn-report-extract":
		documentType = domains.TurnReportExtract
	case "worldographer-map":
		documentType = domains.WorldographerMap
	default:
		return nil, fmt.Errorf("%q: unknown type", d.DocumentType)
	}

	return &domains.Document{
		ID:             domains.ID(d.DocumentID),
		GameID:         domains.GameID(d.GameID),
		ClanId:         domains.ID(d.ClanID),
		ClanNo:         int(d.Clan),
		Path:           d.DocumentName,
		Type:           documentType,
		ContentsLength: len(d.Contents),
		ContentType:    d.ContentType,
		Contents:       d.Contents,
		ModifiedAt:     time.Unix(d.ModifiedAt, 0).UTC(),
		CreatedAt:      time.Unix(d.CreatedAt, 0).UTC(),
		UpdatedAt:      time.Unix(d.UpdatedAt, 0).UTC(),
	}, nil
}

// ReadDocumentsByUser returns an unsorted list of documents that the actor has permissions to view.
// Returns an empty list (not a nil list) if there are no documents.
func (s *Service) ReadDocumentsByUser(actor *domains.Actor, userId domains.ID, docType domains.DocumentType, pageNumber, pageSize int, quiet, verbose, debug bool) ([]*DocumentView, error) {
	var documentType string
	if docType == "" {
		documentType = "*"
	} else {
		switch docType {
		case domains.TurnReportFile:
			documentType = string(domains.TurnReportFile)
		case domains.TurnReportExtract:
			documentType = string(domains.TurnReportExtract)
		case domains.WorldographerMap:
			documentType = string(domains.WorldographerMap)
		default:
			return nil, fmt.Errorf("%q: unknown type", docType)
		}
	}

	// start transaction
	if debug {
		log.Printf("[documents] ReadDocumentsByUser(%d, %d, %q) %q\n", actor.ID, userId, docType, documentType)
	}
	ctx := s.db.Context()
	tx, err := s.db.Stdlib().BeginTx(ctx, nil)
	if err != nil {
		log.Printf("[documents] ReadDocumentsByUser(%d, %d, %q) %v\n", actor.ID, userId, docType, err)
		return nil, errors.Join(domains.ErrDatabaseError, err)
	}
	defer tx.Rollback() // rollback if we return early; harmless after commit
	qtx := s.db.Queries().WithTx(tx)

	docs, err := qtx.ReadDocumentsByUser(ctx, int64(userId))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []*DocumentView{}, nil
		}
		log.Printf("[documents] ReadDocumentsByUser(%d, %d, %q) %v\n", actor.ID, userId, docType, err)
		return nil, errors.Join(domains.ErrDatabaseError, err)
	}
	if debug {
		log.Printf("[documents] ReadDocumentsByUser(%d, %d, %q) %d docs\n", actor.ID, userId, docType, len(docs))
	}

	// cache handles by user_id to avoid repeated database calls
	handles := map[domains.ID]string{}
	actorHandle, err := qtx.ReadHandleByUserId(ctx, int64(actor.ID))
	if err != nil {
		return nil, errors.Join(fmt.Errorf("GetUserHandle(%d)", actor.ID), err)
	}
	//log.Printf("[documents] GetAllDocumentsForUserAcrossGames(%d): actor %q\n", actor.ID, actorHandle)
	handles[actor.ID] = actorHandle

	// cache clans by clan_id to avoid repeated database calls
	clans := map[domains.ID]*domains.Clan{}

	var list []*DocumentView
	for _, doc := range docs {
		if debug {
			log.Printf("[documents] ReadDocumentsByUser(%d, %d, %q) %q\n", actor.ID, userId, docType, doc.DocumentName)
		}
		matchedType := documentType == "*" || documentType == doc.DocumentType
		if !matchedType {
			continue
		}

		ownerClan, ok := clans[domains.ID(doc.ClanID)]
		if !ok {
			clan, err := qtx.GetClan(ctx, doc.ClanID)
			if err != nil {
				return nil, errors.Join(fmt.Errorf("ReadClan(%d)", doc.ClanID), err)
			}
			ownerClan = &domains.Clan{GameID: domains.GameID(clan.GameID), UserID: domains.ID(clan.UserID), ClanID: domains.ID(clan.ClanID), ClanNo: int(clan.Clan)}
		}
		var ownerHandle string
		if ownerHandle, ok = handles[ownerClan.UserID]; !ok {
			ownerHandle, err = qtx.ReadHandleByUserId(ctx, int64(ownerClan.UserID))
			if err != nil {
				return nil, errors.Join(fmt.Errorf("GetUserHandle(%d)", ownerClan.UserID), err)
			}
			handles[ownerClan.UserID] = ownerHandle
		}
		view := &DocumentView{
			ID:           fmt.Sprintf("%d", doc.DocumentID),
			OwnerHandle:  ownerHandle,
			UserHandle:   ownerHandle,
			GameId:       strconv.FormatInt(doc.GameID, 10),
			ClanNo:       fmt.Sprintf("%04d", ownerClan.ClanNo),
			DocumentName: doc.DocumentName,
			DocumentType: doc.DocumentType,
			ModifiedAt:   time.Unix(doc.ModifiedAt, 0).UTC(),
			CreatedAt:    time.Unix(doc.CreatedAt, 0).UTC(),
			UpdatedAt:    time.Unix(doc.UpdatedAt, 0).UTC(),
		}
		list = append(list, view)
	}
	if list == nil {
		list = []*DocumentView{}
	}
	return list, nil
}

// ReadDocumentOwner returns the clan that owns the document.
func (s *Service) ReadDocumentOwner(documentId domains.ID, quiet, verbose, debug bool) (*domains.Clan, error) {
	clan, err := s.db.Queries().ReadDocumentOwner(s.db.Context(), int64(documentId))
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Printf("[documents] ReadDocumentOwner(%d) %v\n", documentId, err)
			return nil, errors.Join(domains.ErrDatabaseError, err)
		}
		return nil, domains.ErrNotExists
	}
	return &domains.Clan{
		GameID:   domains.GameID(clan.GameID),
		UserID:   domains.ID(clan.UserID),
		ClanID:   domains.ID(clan.ClanID),
		ClanNo:   int(clan.Clan),
		IsActive: clan.IsActive,
	}, nil
}

// ReplaceDocument overwrites a document by deleting (if it exists) and creating a new one
func (s *Service) ReplaceDocument(actor *domains.Actor, owner *domains.Clan, doc *domains.Document, quiet, verbose, debug bool) (domains.ID, error) {
	if doc.Path != html.EscapeString(doc.Path) {
		return domains.InvalidID, ErrInvalidPath
	}
	if !s.authzSvc.CanCreateDocuments(actor) {
		return domains.InvalidID, domains.ErrNotAuthorized
	}

	// don't trust the caller on important metadata
	contentLength, contentsHash, err := Hash(doc.Contents)
	if err != nil {
		return domains.InvalidID, errors.Join(domains.ErrHashFailed, err)
	}
	var documentType string
	switch doc.Type {
	case domains.TurnReportFile:
		documentType = string(domains.TurnReportFile)
	case domains.TurnReportExtract:
		documentType = string(domains.TurnReportExtract)
	case domains.WorldographerMap:
		documentType = string(domains.WorldographerMap)
	default:
		return domains.InvalidID, fmt.Errorf("%q: unknown type", doc.Type)
	}

	// start transaction
	ctx := s.db.Context()
	tx, err := s.db.Stdlib().BeginTx(ctx, nil)
	if err != nil {
		log.Printf("[documents] ReplaceDocument(%d, (%q, %d), %q) %v\n", actor.ID, owner.GameID, owner.ClanID, doc.Path, err)
		return domains.InvalidID, errors.Join(domains.ErrDatabaseError, err)
	}
	defer tx.Rollback() // rollback if we return early; harmless after commit
	qtx := s.db.Queries().WithTx(tx)
	now := time.Now().UTC()
	createdAt, updatedAt := now.Unix(), now.Unix()

	err = qtx.DeleteDocumentByClanAndNameAuthorized(ctx, sqlc.DeleteDocumentByClanAndNameAuthorizedParams{
		ClanID:       int64(owner.ClanID),
		DocumentName: doc.Path,
	})
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Printf("[documents] ReplaceDocument(%d, (%q, %d), %q) %v\n", actor.ID, owner.GameID, owner.ClanID, doc.Path, err)
			return domains.InvalidID, errors.Join(domains.ErrDatabaseError, err)
		}
	}

	// create the document
	documentId, err := qtx.CreateDocument(ctx, sqlc.CreateDocumentParams{
		ClanID:       int64(owner.ClanID),
		DocumentName: doc.Path,
		DocumentType: documentType,
		ModifiedAt:   doc.ModifiedAt.Unix(),
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	})
	if err != nil {
		log.Printf("[documents] ReplaceDocument(%d, (%q, %d), %q) %v\n", actor.ID, owner.GameID, owner.ClanID, doc.Path, err)
		return domains.InvalidID, errors.Join(domains.ErrDatabaseError, err)
	}

	// upload the contents
	err = qtx.CreateDocumentContents(ctx, sqlc.CreateDocumentContentsParams{
		DocumentID:    documentId,
		ContentLength: int64(contentLength),
		ContentsHash:  contentsHash,
		Contents:      doc.Contents,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	})
	if err != nil {
		log.Printf("[documents] ReplaceDocument(%d, (%q, %d), %q) %v\n", actor.ID, owner.GameID, owner.ClanID, doc.Path, err)
		return domains.InvalidID, errors.Join(domains.ErrDatabaseError, err)
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("[documents] ReplaceDocument(%d, (%q, %d), %q) %v\n", actor.ID, owner.GameID, owner.ClanID, doc.Path, err)
		return domains.InvalidID, errors.Join(domains.ErrDatabaseError, err)
	}

	if debug {
		log.Printf("[documents] ReplaceDocument(%d, (%q, %d), %q) %d\n", actor.ID, owner.GameID, owner.ClanID, doc.Path, documentId)
	}

	return domains.ID(documentId), nil
}

// SyncDocument will create the document if it does not exist; otherwise it will
// update it if it has changed.
func (s *Service) SyncDocument(actor *domains.Actor, owner *domains.Clan, doc *domains.Document, quiet, verbose, debug bool) (domains.ID, error) {
	if doc.Path != html.EscapeString(doc.Path) {
		return domains.InvalidID, ErrInvalidPath
	}
	if !s.authzSvc.CanCreateDocuments(actor) {
		return domains.InvalidID, domains.ErrNotAuthorized
	}

	// don't trust the caller on important metadata
	contentLength, contentsHash, err := Hash(doc.Contents)
	if err != nil {
		return domains.InvalidID, errors.Join(domains.ErrHashFailed, err)
	}
	var documentType string
	switch doc.Type {
	case domains.TurnReportFile:
		documentType = string(domains.TurnReportFile)
	case domains.TurnReportExtract:
		documentType = string(domains.TurnReportExtract)
	case domains.WorldographerMap:
		documentType = string(domains.WorldographerMap)
	default:
		return domains.InvalidID, fmt.Errorf("%q: unknown type", doc.Type)
	}

	// start transaction
	ctx := s.db.Context()
	tx, err := s.db.Stdlib().BeginTx(ctx, nil)
	if err != nil {
		log.Printf("[documents] SyncDocument(%d, (%q, %d), %q) %v\n", actor.ID, owner.GameID, owner.ClanID, doc.Path, err)
		return domains.InvalidID, errors.Join(domains.ErrDatabaseError, err)
	}
	defer tx.Rollback() // rollback if we return early; harmless after commit
	qtx := s.db.Queries().WithTx(tx)
	now := time.Now().UTC()
	createdAt, updatedAt := now.Unix(), now.Unix()

	d, err := qtx.ReadDocumentByClanAndName(ctx, sqlc.ReadDocumentByClanAndNameParams{
		ClanID:       int64(owner.ClanID),
		DocumentName: doc.Path,
	})
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Printf("[documents] SyncDocument(%d, (%q, %d), %q) %v\n", actor.ID, owner.GameID, owner.ClanID, doc.Path, err)
			return domains.InvalidID, errors.Join(domains.ErrDatabaseError, err)
		}

		// create the document
		documentId, err := qtx.CreateDocument(ctx, sqlc.CreateDocumentParams{
			ClanID:       int64(owner.ClanID),
			DocumentName: doc.Path,
			DocumentType: documentType,
			ModifiedAt:   doc.ModifiedAt.Unix(),
			CreatedAt:    createdAt,
			UpdatedAt:    updatedAt,
		})
		if err != nil {
			log.Printf("[documents] SyncDocument(%d, (%q, %d), %q) %v\n", actor.ID, owner.GameID, owner.ClanID, doc.Path, err)
			return domains.InvalidID, errors.Join(domains.ErrDatabaseError, err)
		}
		doc.ID = domains.ID(documentId)

		// upload the contents
		err = qtx.CreateDocumentContents(ctx, sqlc.CreateDocumentContentsParams{
			DocumentID:    documentId,
			ContentLength: int64(contentLength),
			ContentsHash:  contentsHash,
			Contents:      doc.Contents,
			CreatedAt:     createdAt,
			UpdatedAt:     updatedAt,
		})
		if err != nil {
			log.Printf("[documents] SyncDocument(%d, (%q, %d), %q) %v\n", actor.ID, owner.GameID, owner.ClanID, doc.Path, err)
			return domains.InvalidID, errors.Join(domains.ErrDatabaseError, err)
		}

		err = tx.Commit()
		if err != nil {
			log.Printf("[documents] SyncDocument(%d, (%q, %d), %q) %v\n", actor.ID, owner.GameID, owner.ClanID, doc.Path, err)
			return domains.InvalidID, errors.Join(domains.ErrDatabaseError, err)
		}
		return domains.ID(documentId), nil
	}
	doc.ID = domains.ID(d.DocumentID)

	// do we need to update the document contents and meta-data?
	if doc.ContentsHash != d.ContentsHash {
		err = qtx.UpdateDocumentContentsById(ctx, sqlc.UpdateDocumentContentsByIdParams{
			DocumentID:    d.DocumentID,
			ContentLength: int64(contentLength),
			ContentsHash:  contentsHash,
			Contents:      doc.Contents,
			UpdatedAt:     updatedAt,
		})
		if err != nil {
			log.Printf("[documents] SyncDocument(%d, (%q, %d), %q) %v\n", actor.ID, owner.GameID, owner.ClanID, doc.Path, err)
			return domains.InvalidID, errors.Join(domains.ErrDatabaseError, err)
		}

		err = qtx.UpdateDocumentById(ctx, sqlc.UpdateDocumentByIdParams{
			ClanID:       d.ClanID,
			DocumentID:   d.DocumentID,
			DocumentName: d.DocumentName,
			DocumentType: d.DocumentType,
			ModifiedAt:   doc.ModifiedAt.Unix(),
			UpdatedAt:    updatedAt,
		})
		if err != nil {
			log.Printf("[documents] SyncDocument(%d, (%q, %d), %q) %v\n", actor.ID, owner.GameID, owner.ClanID, doc.Path, err)
			return domains.InvalidID, errors.Join(domains.ErrDatabaseError, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("[documents] SyncDocument(%d, (%q, %d), %q) %v\n", actor.ID, owner.GameID, owner.ClanID, doc.Path, err)
		return domains.InvalidID, errors.Join(domains.ErrDatabaseError, err)
	}

	if debug {
		log.Printf("[documents] SyncDocument(%d, (%q, %d), %q) %d\n", actor.ID, owner.GameID, owner.ClanID, doc.Path, doc.ID)
	}

	return doc.ID, nil
}

// UpdateDocument will update the document if it exists and the contents have changed.
func (s *Service) UpdateDocument(actor *domains.Actor, owner *domains.Clan, doc *domains.Document, quiet, verbose, debug bool) error {
	if doc.Path != html.EscapeString(doc.Path) {
		return ErrInvalidPath
	}
	if !s.authzSvc.CanCreateDocuments(actor) {
		return domains.ErrNotAuthorized
	}

	// don't trust the caller on important metadata
	contentLength, contentsHash, err := Hash(doc.Contents)
	if err != nil {
		return errors.Join(domains.ErrHashFailed, err)
	}
	var documentType string
	switch doc.Type {
	case domains.TurnReportFile:
		documentType = string(domains.TurnReportFile)
	case domains.TurnReportExtract:
		documentType = string(domains.TurnReportExtract)
	case domains.WorldographerMap:
		documentType = string(domains.WorldographerMap)
	default:
		return fmt.Errorf("%q: unknown type", doc.Type)
	}

	// start transaction
	ctx := s.db.Context()
	tx, err := s.db.Stdlib().BeginTx(ctx, nil)
	if err != nil {
		log.Printf("[documents] UpdateDocument(%d, (%q, %d), %q) %v\n", actor.ID, owner.GameID, owner.ClanID, doc.Path, err)
		return errors.Join(domains.ErrDatabaseError, err)
	}
	defer tx.Rollback() // rollback if we return early; harmless after commit
	qtx := s.db.Queries().WithTx(tx)
	now := time.Now().UTC()
	updatedAt := now.Unix()

	d, err := qtx.ReadDocumentByClanAndName(ctx, sqlc.ReadDocumentByClanAndNameParams{
		ClanID:       int64(owner.ClanID),
		DocumentName: doc.Path,
	})
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Printf("[documents] UpdateDocument(%d, (%q, %d), %q) %v\n", actor.ID, owner.GameID, owner.ClanID, doc.Path, err)
			return errors.Join(domains.ErrDatabaseError, err)
		}
		return domains.ErrNotExists
	}
	doc.ID = domains.ID(d.DocumentID)

	// do we need to update the document contents?
	updatedContents := false
	if doc.ContentsHash != d.ContentsHash {
		err = qtx.UpdateDocumentContentsById(ctx, sqlc.UpdateDocumentContentsByIdParams{
			DocumentID:    d.DocumentID,
			ContentLength: int64(contentLength),
			ContentsHash:  contentsHash,
			Contents:      doc.Contents,
			UpdatedAt:     updatedAt,
		})
		if err != nil {
			log.Printf("[documents] UpdateDocument(%d, (%q, %d), %q) %v\n", actor.ID, owner.GameID, owner.ClanID, doc.Path, err)
			return errors.Join(domains.ErrDatabaseError, err)
		}
		updatedContents = true
	}

	// do we need to update the document meta-data?
	if updatedContents || doc.Path != d.DocumentName || documentType != d.DocumentType {
		err = qtx.UpdateDocumentById(ctx, sqlc.UpdateDocumentByIdParams{
			ClanID:       d.ClanID,
			DocumentID:   d.DocumentID,
			DocumentName: doc.Path,
			DocumentType: documentType,
			ModifiedAt:   doc.ModifiedAt.Unix(),
			UpdatedAt:    updatedAt,
		})
		if err != nil {
			log.Printf("[documents] UpdateDocument(%d, (%q, %d), %q) %v\n", actor.ID, owner.GameID, owner.ClanID, doc.Path, err)
			return errors.Join(domains.ErrDatabaseError, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("[documents] UpdateDocument(%d, (%q, %d), %q) %v\n", actor.ID, owner.GameID, owner.ClanID, doc.Path, err)
		return errors.Join(domains.ErrDatabaseError, err)
	}

	if debug {
		log.Printf("[documents] UpdateDocument(%d, (%q, %d), %q) %d\n", actor.ID, owner.GameID, owner.ClanID, doc.Path, doc.ID)
	}

	return nil
}

// DeleteDocument will return nil if the document doesn't exist.
// Returns an error if the actor is not authorized to delete it.
func (s *Service) DeleteDocument(actor *domains.Actor, owner *domains.Clan, documentId domains.ID) error {
	err := s.db.Queries().DeleteDocumentByIdAuthorized(s.db.Context(), sqlc.DeleteDocumentByIdAuthorizedParams{
		DocumentID: int64(documentId),
		ClanID:     int64(owner.ClanID),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// no document, so not an error
			return nil
		}
		log.Printf("[documents] DeleteDocument(%d, (%q, %d), %d) %v\n", actor.ID, owner.GameID, owner.ClanID, documentId, err)
		return errors.Join(domains.ErrDatabaseError, err)
	}
	return nil
}

func (s *Service) ReadReportExtractContents(documentId domains.ID) ([]byte, error) {
	return s.db.Queries().ReadDocumentContents(s.db.Context(), int64(documentId))
}

func (s *Service) ReadReportExtractMeta() ([]*domains.Document, error) {
	rows, err := s.db.Queries().ReadReportExtracts(s.db.Context())
	if err != nil {
		log.Printf("[documents] ReadReportExtractMeta:\n", err)
		return nil, errors.Join(domains.ErrDatabaseError, err)
	}

	var docs []*domains.Document
	for _, d := range rows {
		docs = append(docs, &domains.Document{
			ID:         domains.ID(d.DocumentID),
			GameID:     domains.GameID(d.GameID),
			ClanId:     domains.InvalidID,
			ClanNo:     int(d.Clan),
			Path:       d.DocumentName,
			Type:       domains.TurnReportExtract,
			ModifiedAt: time.Unix(d.ModifiedAt, 0).UTC(),
			CreatedAt:  time.Unix(d.CreatedAt, 0).UTC(),
			UpdatedAt:  time.Unix(d.UpdatedAt, 0).UTC(),
		})
	}
	if docs == nil {
		docs = []*domains.Document{}
	}
	return docs, nil
}

// loadFromFS loads a file, creates a Document, and returns the document ID.
func (s *Service) loadFromFS(actor *domains.Actor, clan *domains.Clan, path, name string, docType domains.DocumentType, quiet, verbose, debug bool) (domains.ID, error) {
	//log.Printf("[documents] loadFromFS(%d, (%q, %d, %d), %q, %q, _, %q) read %v", actor.ID, clan.GameID, clan.UserID, clan.ClanNo, path, name, docType, canRead)
	if name == "" {
		return domains.InvalidID, errors.Join(domains.ErrBadInput, fmt.Errorf("missing name"))
	}
	if !actor.IsValid() || clan == nil {
		return domains.InvalidID, domains.ErrNotAuthorized
	}

	// Stat so we can validate the path and get size/timestamps.
	sb, err := os.Stat(path)
	if err != nil {
		return domains.InvalidID, errors.Join(domains.ErrInvalidPath, fmt.Errorf("%q: stat failed", path), err)
	}
	if sb.IsDir() || !sb.Mode().IsRegular() {
		// keep ErrInvalidPath but add context so logs make sense
		return domains.InvalidID, errors.Join(domains.ErrInvalidPath, fmt.Errorf("%q: not a regular file", path))
	}
	createdAt, updatedAt := sb.ModTime().UTC(), sb.ModTime().UTC()

	data, err := os.ReadFile(path)
	if err != nil {
		return domains.InvalidID, errors.Join(domains.ErrReadFailed, fmt.Errorf("%q: read failed", path), err)
	}

	doc := &domains.Document{
		Path:      name,
		ClanId:    clan.ClanID,
		Type:      docType,
		Contents:  data,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	return s.CreateDocument(actor, clan, doc, quiet, verbose, debug)
}

// loadFromRequest loads a document from an http.Request.
// It's used by the file upload handlers.
func (s *Service) loadFromRequest(r *http.Request) (domains.ID, error) {
	return domains.InvalidID, domains.ErrNotImplemented
}

func (s *Service) FindClanByGameAndNumber(gameId domains.GameID, clanNo int) (*domains.Clan, error) {
	q := s.db.Queries()
	ctx := s.db.Context()
	row, err := q.GetClanByGameClanNo(ctx, sqlc.GetClanByGameClanNoParams{
		GameID: int64(gameId),
		ClanNo: int64(clanNo),
	})
	if err != nil {
		return nil, err
	}
	return &domains.Clan{
		GameID: domains.GameID(row.GameID),
		UserID: domains.ID(row.UserID),
		ClanID: domains.ID(row.ClanID),
		ClanNo: int(row.Clan),
	}, nil
}

func Hash(contents []byte) (length int, hash string, err error) {
	h := sha256.New()
	if _, err := h.Write(contents); err != nil {
		return 0, "", errors.Join(domains.ErrHashFailed, err)
	}
	return len(contents), hex.EncodeToString(h.Sum(nil)), nil
}
