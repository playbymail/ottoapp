// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package documents

import (
	"net/http"

	"github.com/playbymail/ottoapp/backend/domains"
)

func (s *Service) LoadMapFromFS(actor *domains.Actor, clan *domains.Clan, path, name string, canDelete, canRead, canShare, canWrite bool) (domains.ID, error) {
	//log.Printf("[documents] LoadMapFromFS(%d, %d, %q, %q) %v", actor.ID, clan.ClanNo, path, name, canRead)
	return s.loadFromFS(actor, clan, path, name, canDelete, canRead, canShare, canWrite, domains.WXXMimeType, domains.WorldographerMap)
}

func (s *Service) LoadMapFromRequest(r *http.Request) (domains.ID, error) {
	return s.loadFromRequest(r)
}
