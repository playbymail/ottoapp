// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package documents

import (
	"fmt"
	"time"

	"github.com/hashicorp/jsonapi"
	"github.com/playbymail/ottoapp/backend/domains"
)

// DocumentView is the JSON:API view for a document.
type DocumentView struct {
	ID           string    `jsonapi:"primary,document"`   // singular when sending a payload
	OwnerHandle  string    `jsonapi:"attr,owner"`         // handle of user that owns this document
	UserHandle   string    `jsonapi:"attr,created-by"`    // handle of user for this document
	GameId       string    `jsonapi:"attr,game-id"`       // game for this document
	ClanNo       string    `jsonapi:"attr,clan"`          // clan for this document
	DocumentName string    `jsonapi:"attr,document-name"` // untainted name of document
	DocumentType string    `jsonapi:"attr,document-type"` // for categorizing on the dashboards
	ModifiedAt   time.Time `jsonapi:"attr,modified-at,iso8601"`
	CreatedAt    time.Time `jsonapi:"attr,created-at,iso8601"`
	UpdatedAt    time.Time `jsonapi:"attr,updated-at,iso8601"`
}

// JSONAPILinks implements the jsonapi.Linkable interface for document-links
func (d *DocumentView) JSONAPILinks() *jsonapi.Links {
	return &jsonapi.Links{
		"self": fmt.Sprintf("/api/documents/%s", d.ID),
		"contents": jsonapi.Link{
			Href: fmt.Sprintf("/api/documents/%s/contents", d.ID),
		},
	}
}

// UserDocumentView is the JSON:API for a user document.
// This view of the document includes the game information.
type UserDocumentView struct {
	ID               string    `jsonapi:"primary,user-document"` // singular when sending a payload
	Owner            string    `jsonapi:"attr,owner"`            // handle of user that owns this document
	Game             string    `jsonapi:"attr,game"`
	Clan             string    `jsonapi:"attr,clan"`          // four-digit clan number
	DocumentName     string    `jsonapi:"attr,document-name"` // untainted name of document
	DocumentType     string    `jsonapi:"attr,document-type"` // for categorizing on the dashboards
	ProcessingStatus string    `jsonapi:"attr,processing-status,omitempty"`
	UpdatedAt        time.Time `jsonapi:"attr,updated-at,iso8601"`
}

// JSONAPILinks implements the jsonapi.Linkable interface for user-document-links
func (d *UserDocumentView) JSONAPILinks() *jsonapi.Links {
	switch domains.DocumentType(d.DocumentType) {
	case domains.TurnReportFile:
		return &jsonapi.Links{
			"self": fmt.Sprintf("/api/documents/%s", d.ID),
			"show": jsonapi.Link{
				Href: fmt.Sprintf("/documents/%s", d.ID),
			},
		}
	case domains.TurnReportExtract:
		return &jsonapi.Links{
			"self": fmt.Sprintf("/api/documents/%s", d.ID),
			"show": jsonapi.Link{
				Href: fmt.Sprintf("/documents/%s", d.ID),
			},
		}
	case domains.WorldographerMap:
		return &jsonapi.Links{
			"self": fmt.Sprintf("/api/documents/%s", d.ID),
			"contents": jsonapi.Link{
				Href: fmt.Sprintf("/api/documents/%s/contents", d.ID),
			},
		}
	}
	panic(fmt.Sprintf("assert(docType != %q)", d.DocumentType))
}
