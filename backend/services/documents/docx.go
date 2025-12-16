// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package documents

import (
	"net/http"

	"github.com/playbymail/ottoapp/backend/domains"
)

// LoadDocxFromFS loads the file, creates a Document, and returns the document ID.
func (s *Service) LoadDocxFromFS(actor *domains.Actor, clan *domains.Clan, path, name string, quiet, verbose, debug bool) (domains.ID, error) {
	return s.loadFromFS(actor, clan, path, name, domains.TurnReportFile, quiet, verbose, debug)
}

// LoadDocxFromRequest loads a document from an http.Request.
func (s *Service) LoadDocxFromRequest(r *http.Request) (domains.ID, error) {
	return s.loadFromRequest(r)
}

func (s *Service) UpdateDocx(doc *domains.Document) error {
	return domains.ErrNotImplemented
}
