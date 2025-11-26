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
	"time"

	"github.com/playbymail/ottoapp/backend/domains"
	"github.com/playbymail/ottoapp/backend/services/authz"
	"github.com/playbymail/ottoapp/backend/stores/sqlite"
	"github.com/playbymail/ottoapp/backend/stores/sqlite/sqlc"
	"github.com/playbymail/ottoapp/backend/users"
)

// Service provides document management operations.
type Service struct {
	db       *sqlite.DB
	authzSvc *authz.Service
	usersSvc *users.Service
}

func New(db *sqlite.DB, authzSvc *authz.Service, usersSvc *users.Service) *Service {
	return &Service{db: db, authzSvc: authzSvc, usersSvc: usersSvc}
}

// CreateDocument creates a document.
//
// Actor is the user/service requesting the creation.
// Owner is the user that will own the new document.
func (s *Service) CreateDocument(actor *domains.Actor, clan *domains.Clan, doc *domains.Document) (domains.ID, error) {
	//log.Printf("[documents] CreateDocument(%d, (%q, %d), %q, %q) d:%v r:%v s:%v w:%v", actor.ID, clan.GameID, clan.UserID, doc.Path, doc.Type, doc.CanDelete, doc.CanRead, doc.CanShare, doc.CanWrite)
	if doc.Path != html.EscapeString(doc.Path) {
		return domains.InvalidID, ErrInvalidPath
	}
	if !s.authzSvc.CanCreateDocuments(actor) {
		return domains.InvalidID, domains.ErrNotAuthorized
	}

	// don't trust the caller on important metadata
	contentLength := len(doc.Contents)
	contentsHash, err := hashContents(doc.Contents)
	if err != nil {
		return domains.InvalidID, errors.Join(domains.ErrHashFailed, err)
	}

	ctx := s.db.Context()

	// start transaction
	tx, err := s.db.Stdlib().BeginTx(ctx, nil)
	if err != nil {
		return domains.InvalidID, err
	}
	defer tx.Rollback() // rollback if we return early; harmless after commit

	qtx := s.db.Queries().WithTx(tx)
	now := time.Now().UTC().Unix()

	err = qtx.CreateDocumentContents(ctx, sqlc.CreateDocumentContentsParams{
		ContentsHash:  contentsHash,
		ContentLength: int64(contentLength),
		MimeType:      string(doc.MimeType),
		Contents:      doc.Contents,
		CreatedAt:     now,
		UpdatedAt:     now,
	})
	//log.Printf("[documents] CreateDocument(%d, (%q, %d), %q, %q) %v", actor.ID, clan.GameID, clan.UserID, doc.Path, doc.Type, err)
	if err != nil {
		return domains.InvalidID, errors.Join(domains.ErrDatabaseError, fmt.Errorf("CreateDocumentContents(%q)", contentsHash), err)
	}

	id, err := qtx.CreateDocument(ctx, sqlc.CreateDocumentParams{
		ClanID:       int64(clan.ClanID),
		CanRead:      doc.CanRead,
		CanWrite:     doc.CanWrite,
		CanDelete:    doc.CanDelete,
		CanShare:     doc.CanShare,
		DocumentName: doc.Path,
		DocumentType: string(doc.Type),
		ContentsHash: contentsHash,
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	//log.Printf("[documents] CreateDocument(%d, (%q, %d), %q, %q) %d %v", actor.ID, clan.GameID, clan.UserID, doc.Path, doc.Type, id, err)
	if err != nil {
		return domains.InvalidID, errors.Join(domains.ErrDatabaseError, fmt.Errorf("CreateDocument(%d)", clan.ClanID), err)
	}

	err = tx.Commit()
	if err != nil {
		return domains.InvalidID, err
	}

	return domains.ID(id), nil
}

// DeleteDocument will return nil if the document doesn't exist.
// Returns an error if the actor is not authorized to delete it.
func (s *Service) DeleteDocument(actor *domains.Actor, documentId domains.ID) error {
	// todo: should verify that actor is allowed to delete the target
	ctx := s.db.Context()

	// start transaction
	tx, err := s.db.Stdlib().BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() // rollback if we return early; harmless after commit

	qtx := s.db.Queries().WithTx(tx)

	// fetch the document
	doc, err := qtx.GetDocumentById(ctx, int64(documentId))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return errors.Join(domains.ErrDatabaseError, err)
	}
	// we have to look up the clan to determine if the actor owns it or if it is shared
	clan, err := qtx.GetClan(ctx, doc.ClanID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { // should never happen
			log.Printf("[documents] DeleteDocument(%d, %d): GetClan(%d): does not exist", actor.ID, documentId, doc.ClanID)
			return nil
		}
		return errors.Join(domains.ErrDatabaseError, err)
	}
	isOwner := actor.ID == domains.ID(clan.UserID)
	if isOwner {
		if !doc.CanDelete {
			return domains.ErrNotAuthorized
		}
		err = qtx.DeleteDocumentAuthorized(ctx, sqlc.DeleteDocumentAuthorizedParams{
			DocumentID: int64(documentId),
			ClanID:     clan.ClanID,
		})
		if err != nil {
			return errors.Join(domains.ErrDatabaseError, err)
		}
		return nil
	}

	err = qtx.DeleteSharedDocumentById(ctx, sqlc.DeleteSharedDocumentByIdParams{
		DocumentID: int64(documentId),
		ClanID:     int64(doc.ClanID),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { // should never happen
			log.Printf("[documents] DeleteDocument(%d, %d): DeleteSharedDocumentById(%d, %d): does not exist", actor.ID, documentId, doc.ClanID)
			return nil
		}
		return err
	}

	return tx.Commit()
}

// GetDocument returns nil, nil if no data found
func (s *Service) GetDocument(actor *domains.Actor, documentId domains.ID) (*DocumentView, error) {
	//log.Printf("[documents] GetDocument(%d, %d)\n", actor.ID, documentId)
	q := s.db.Queries()
	ctx := s.db.Context()

	doc, err := q.GetDocumentForUserAuthorized(ctx, sqlc.GetDocumentForUserAuthorizedParams{
		DocumentID: int64(documentId),
		UserID:     int64(actor.ID),
	})
	//log.Printf("[documents] GetDocumentForUserAuthorized(%d, %d) %v\n", actor.ID, documentId, err)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Join(domains.ErrDatabaseError, err)
	}

	actorHandle, err := q.GetUserHandle(ctx, int64(actor.ID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { // should never happen
			log.Printf("[documents] GetDocument(%d, %d): GetUserHandle(%d): does not exist", actor.ID, documentId, actor.ID)
			return nil, nil
		}
		return nil, errors.Join(domains.ErrDatabaseError, err)
	}
	var ownerHandle string
	if actor.ID == domains.ID(doc.OwnerID) {
		ownerHandle = actorHandle
	} else if ownerHandle, err = q.GetUserHandle(ctx, doc.OwnerID); err != nil {
		if errors.Is(err, sql.ErrNoRows) { // should never happen
			log.Printf("[documents] GetDocument(%d, %d): GetUserHandle(%d): does not exist", actor.ID, documentId, doc.OwnerID)
			return nil, nil
		}
		return nil, errors.Join(domains.ErrDatabaseError, err)
	}

	return &DocumentView{
		ID:           fmt.Sprintf("%d", doc.DocumentID),
		OwnerHandle:  ownerHandle,
		UserHandle:   actorHandle,
		GameId:       doc.GameID,
		ClanNo:       fmt.Sprintf("%04d", doc.ClanID),
		DocumentName: doc.DocumentName,
		DocumentType: doc.DocumentType,
		CanRead:      doc.CanRead,
		CanWrite:     doc.CanWrite,
		CanDelete:    doc.CanDelete,
		CanShare:     doc.CanShare,
		IsShared:     doc.IsShared,
		CreatedAt:    time.Unix(doc.CreatedAt, 0).UTC(),
		UpdatedAt:    time.Unix(doc.UpdatedAt, 0).UTC(),
	}, nil
}

// GetDocumentContents returns nil, nil if no data found
func (s *Service) GetDocumentContents(actor *domains.Actor, documentId domains.ID) (*domains.Document, error) {
	q := s.db.Queries()
	ctx := s.db.Context()

	doc, err := q.GetDocumentForUserAuthorized(ctx, sqlc.GetDocumentForUserAuthorizedParams{
		DocumentID: int64(documentId),
		UserID:     int64(actor.ID),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Join(domains.ErrDatabaseError, err)
	}
	if !doc.CanRead {
		return nil, domains.ErrNotAuthorized
	}

	content, err := s.db.Queries().GetDocumentContents(s.db.Context(), doc.ContentsHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { // should never happen
			log.Printf("[documents] GetDocumentContents(%d, %d): GetDocumentContents(%q): does not exist", actor.ID, documentId, doc.ContentsHash)
			return nil, nil
		}
		return nil, errors.Join(domains.ErrDatabaseError, err)
	}
	return &domains.Document{
		ID:            domains.ID(doc.DocumentID),
		Path:          doc.DocumentName,
		Type:          domains.DocumentType(doc.DocumentType),
		MimeType:      domains.MimeType(content.MimeType),
		ContentLength: content.ContentLength,
		Contents:      content.Contents,
		ContentsHash:  doc.ContentsHash,
		CreatedAt:     time.Unix(doc.CreatedAt, 0).UTC(),
		UpdatedAt:     time.Unix(doc.UpdatedAt, 0).UTC(),
	}, nil
}

// GetAllDocumentsForUserAcrossGames returns an unsorted list of documents that the actor has permissions to view.
// Returns an empty list (not a nil list) if there are no documents.
func (s *Service) GetAllDocumentsForUserAcrossGames(actor *domains.Actor, docType domains.DocumentType, pageNumber, pageSize int) ([]*DocumentView, error) {
	//log.Printf("[documents] GetAllDocumentsForUserAcrossGames(%d, %q)\n", actor.ID, docType)

	q := s.db.Queries()
	ctx := s.db.Context()

	docs, err := q.GetAllDocumentsForUserAcrossGames(ctx, int64(actor.ID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Join(domains.ErrDatabaseError, fmt.Errorf("GetAllDocumentsForUserAcrossGames(%d)", actor.ID), err)
	}
	//log.Printf("[documents] GetAllDocumentsForUserAcrossGames(%d): %d docs\n", actor.ID, len(docs))

	// cache handles by user_id to avoid repeated database calls
	handles := map[domains.ID]string{}
	actorHandle, err := q.GetUserHandle(ctx, int64(actor.ID))
	if err != nil {
		return nil, errors.Join(fmt.Errorf("GetUserHandle(%d)", actor.ID), err)
	}
	//log.Printf("[documents] GetAllDocumentsForUserAcrossGames(%d): actor %q\n", actor.ID, actorHandle)
	handles[actor.ID] = actorHandle

	// cache clans by clan_id to avoid repeated database calls
	clans := map[domains.ID]*domains.Clan{}

	var list []*DocumentView
	for _, doc := range docs {
		//log.Printf("[documents] GetAllDocumentsForUserAcrossGames(%d): doc %d: game_id %q\n", actor.ID, doc.DocumentID, doc.GameID)
		//log.Printf("[documents] GetAllDocumentsForUserAcrossGames(%d): doc %d: user_id %d\n", actor.ID, doc.DocumentID, doc.UserID)
		//log.Printf("[documents] GetAllDocumentsForUserAcrossGames(%d): doc %d: clan_id %d\n", actor.ID, doc.DocumentID, doc.ClanID)
		//log.Printf("[documents] GetAllDocumentsForUserAcrossGames(%d): doc %d: owner   %d\n", actor.ID, doc.DocumentID, doc.OwnerID)
		//log.Printf("doc %+v\n", doc)
		if docType != "" && !(string(docType) == doc.DocumentType) {
			//log.Printf("doc kind %q != %q\n", doc.DocumentType, docType)
			continue
		}

		ownerClan, ok := clans[domains.ID(doc.OwnerID)]
		if !ok {
			clan, err := q.GetClan(ctx, doc.OwnerID)
			if err != nil {
				return nil, errors.Join(fmt.Errorf("GetClan(%d)", doc.OwnerID), err)
			}
			ownerClan = &domains.Clan{GameID: clan.GameID, UserID: domains.ID(clan.UserID), ClanID: domains.ID(clan.ClanID), ClanNo: int(clan.Clan)}
		}
		var ownerHandle string
		if ownerHandle, ok = handles[ownerClan.UserID]; !ok {
			ownerHandle, err = q.GetUserHandle(ctx, int64(ownerClan.UserID))
			if err != nil {
				return nil, errors.Join(fmt.Errorf("GetUserHandle(%d)", ownerClan.UserID), err)
			}
			handles[ownerClan.UserID] = ownerHandle
		}
		view := &DocumentView{
			ID:           fmt.Sprintf("%d", doc.DocumentID),
			OwnerHandle:  ownerHandle,
			UserHandle:   actorHandle,
			GameId:       doc.GameID,
			ClanNo:       fmt.Sprintf("%04d", ownerClan.ClanNo),
			DocumentName: doc.DocumentName,
			DocumentType: doc.DocumentType,
			CanRead:      doc.CanRead,
			CanWrite:     doc.CanWrite,
			CanDelete:    doc.CanDelete,
			CanShare:     doc.CanShare,
			IsShared:     doc.IsShared,
			CreatedAt:    time.Unix(doc.CreatedAt, 0).UTC(),
			UpdatedAt:    time.Unix(doc.UpdatedAt, 0).UTC(),
		}
		//log.Printf("view %+v\n", *view)
		list = append(list, view)
	}
	if list == nil {
		list = []*DocumentView{}
	}
	return list, nil
}

// ShareDocumentById shares a document.
//
// Actor is the user requesting the share. The actor must own the
// document and have permission to read and share it.
//
// Ally is the user that will own the shared document.
func (s *Service) ShareDocumentById(actor, ally *domains.Actor, documentId domains.ID, canRead, canDelete bool) error {
	q := s.db.Queries()
	ctx := s.db.Context()

	doc, err := q.GetDocumentForUserAuthorized(ctx, sqlc.GetDocumentForUserAuthorizedParams{
		DocumentID: int64(documentId),
		UserID:     int64(actor.ID),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return errors.Join(domains.ErrDatabaseError, fmt.Errorf("GetDocumentForUserAuthorized(%d, %d)", documentId, actor.ID), err)
	}
	if actor.ID != domains.ID(doc.OwnerID) || doc.IsShared || !doc.CanShare {
		return domains.ErrNotAuthorized
	}

	allyClan, err := q.GetClanByGameUser(ctx, sqlc.GetClanByGameUserParams{
		GameID: doc.GameID,
		UserID: int64(ally.ID),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { // should never happen
			log.Printf("[documents] ShareDocumentById(%d, %d, %d): GetClanByGameUser(%q, %d): does not exist", actor.ID, ally.ID, documentId, doc.GameID, ally.ID)
			return nil
		}
	}

	now := time.Now().UTC().Unix()
	err = s.db.Queries().ShareDocumentById(s.db.Context(), sqlc.ShareDocumentByIdParams{
		DocumentID: int64(documentId),
		ClanID:     allyClan.ClanID,
		CanRead:    canRead,
		CanDelete:  canDelete,
		CreatedAt:  now,
		UpdatedAt:  now,
	})
	if err != nil {
		return errors.Join(domains.ErrDatabaseError, fmt.Errorf("ShareDocumentById(%d)", documentId), err)
	}
	return nil
}
func hashContents(b []byte) (string, error) {
	h := sha256.New()
	if _, err := h.Write(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// loadFromFS loads a file, creates a Document, and returns the document ID.
func (s *Service) loadFromFS(actor *domains.Actor, clan *domains.Clan, path, name string, canDelete, canRead, canShare, canWrite bool, mimeType domains.MimeType, docType domains.DocumentType) (domains.ID, error) {
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

	data, err := os.ReadFile(path)
	if err != nil {
		return domains.InvalidID, errors.Join(domains.ErrReadFailed, fmt.Errorf("%q: read failed", path), err)
	}

	doc := &domains.Document{
		Path:      name,
		ClanId:    clan.ClanID,
		Type:      docType,
		MimeType:  mimeType,
		Contents:  data,
		CanDelete: canDelete,
		CanRead:   canRead,
		CanShare:  canShare,
		CanWrite:  canWrite,
	}

	return s.CreateDocument(actor, clan, doc)
}

// loadFromRequest loads a document from an http.Request.
// It's used by the file upload handlers.
func (s *Service) loadFromRequest(r *http.Request) (domains.ID, error) {
	return domains.InvalidID, domains.ErrNotImplemented
}

func (s *Service) FindClanByGameAndNumber(gameID string, clanNo int) (*domains.Clan, error) {
	q := s.db.Queries()
	ctx := s.db.Context()
	row, err := q.GetClanByGameClanNo(ctx, sqlc.GetClanByGameClanNoParams{
		GameID: gameID,
		ClanNo: int64(clanNo),
	})
	if err != nil {
		return nil, err
	}
	return &domains.Clan{
		GameID: row.GameID,
		UserID: domains.ID(row.UserID),
		ClanID: domains.ID(row.ClanID),
		ClanNo: int(row.Clan),
	}, nil
}
