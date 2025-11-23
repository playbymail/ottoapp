// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package games

import (
	"time"

	"github.com/playbymail/ottoapp/backend/restapi"
)

// ClanDocumentView is the JSON:API view for a clan document.
type ClanDocumentView struct {
	ID           string    `jsonapi:"primary,document"`   // singular when sending a payload
	Game         string    `jsonapi:"attr,game"`          // id of game the document is from
	ClanNo       string    `jsonapi:"attr,owner"`         // clan the document is for
	DocumentName string    `jsonapi:"attr,document-name"` // untainted name of document
	DocumentType string    `jsonapi:"attr,document-type"` // for categorizing on the dashboards
	CanRead      bool      `jsonapi:"attr,can-read"`
	CanWrite     bool      `jsonapi:"attr,can-write"`
	CanDelete    bool      `jsonapi:"attr,can-delete"`
	CanShare     bool      `jsonapi:"attr,can-share"`
	CreatedBy    string    `jsonapi:"attr,created-by"` // handle of user that created this document
	CreatedAt    time.Time `jsonapi:"attr,created-at,iso8601"`
	UpdatedAt    time.Time `jsonapi:"attr,updated-at,iso8601"`
	//Li
}

// ClanDocumentAttributes are the JSON:API "attributes" for clan-document.
// This is the game-domain metadata derived by the parsers / game engine.
type ClanDocumentAttributes struct {
	GameID string `json:"gameId"` // games.game_id
	ClanID int64  `json:"clanId"` // clans.clan_id

	// Optional but useful for maps / reports
	TurnNo *int64 `json:"turnNo,omitempty"`
	Kind   string `json:"kind,omitempty"` // e.g. "map", "report", "orders"

	CreatedAt string `json:"createdAt"` // RFC3339
	UpdatedAt string `json:"updatedAt"`
}

// ClanDocumentRelationships wires up links to other resources.
// Most important: the underlying document.
type ClanDocumentRelationships struct {
	Document *restapi.ToOneRelationship `json:"document,omitempty"`
	Clan     *restapi.ToOneRelationship `json:"clan,omitempty"`
	// You can add Game as a relationship later if you expose a game resource.
	// Game *ToOneRelationship `json:"game,omitempty"`
}

// ClanDocumentResource is the JSON:API resource object for clan-documents.
type ClanDocumentResource struct {
	Type          string                    `json:"type"` // always "clan-documents"
	ID            string                    `json:"id"`
	Attributes    ClanDocumentAttributes    `json:"attributes"`
	Relationships ClanDocumentRelationships `json:"relationships,omitempty"`
}

// SingleClanDocumentResponse wraps a single clan-document resource.
type SingleClanDocumentResponse struct {
	Data  ClanDocumentResource    `json:"data"`
	Meta  *restapi.PaginationMeta `json:"meta,omitempty"`
	Links *restapi.Links          `json:"links,omitempty"`
}

// ManyClanDocumentsResponse wraps multiple clan-document resources.
type ManyClanDocumentsResponse struct {
	Data  []ClanDocumentResource  `json:"data"`
	Meta  *restapi.PaginationMeta `json:"meta,omitempty"`
	Links *restapi.Links          `json:"links,omitempty"`
}
