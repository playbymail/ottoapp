// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package documents

import (
	"net/http"

	"github.com/playbymail/ottoapp/backend/domains"
)

// LoadReportFromFS loads the file, creates a Document, and returns the document ID.
func (s *Service) LoadReportFromFS(actor *domains.Actor, clan *domains.Clan, path, name string, quiet, verbose, debug bool) (domains.ID, error) {
	return s.loadFromFS(actor, clan, path, name, domains.TurnReportExtract, quiet, verbose, debug)
}

// LoadReportFromRequest loads a document from an http.Request.
func (s *Service) LoadReportFromRequest(r *http.Request) (domains.ID, error) {
	return s.loadFromRequest(r)
}
