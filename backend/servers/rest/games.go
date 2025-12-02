// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package rest

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/playbymail/ottoapp/backend/domains"
	"github.com/playbymail/ottoapp/backend/restapi"
	"github.com/playbymail/ottoapp/backend/services/authz"
	"github.com/playbymail/ottoapp/backend/services/documents"
	"github.com/playbymail/ottoapp/backend/services/games"
	parsers "github.com/playbymail/ottoapp/backend/services/reports/docx"
	"github.com/playbymail/ottoapp/backend/services/reports/office"
)

// PostGamesTurnReportFiles creates a new turn report document for the clan specified
// in the document. If the document already exists, overwrite it.
//
// The GameID is extracted from the route path and the ClanNo from the document. These
// are used to create the document - don't use this handler if you want to upload a
// document for a different user!
//
// Route: POST /api/games/:game_id/turn-report-files
//
// Response type: TBD
func PostGamesTurnReportFiles(authzSvc *authz.Service, documentsSvc *documents.Service, gamesSvc *games.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		actor, err := authzSvc.GetActor(r)
		if err != nil {
			log.Printf("%s %s: GetActor: %v\n", r.Method, r.URL.Path, err)
			restapi.WriteJsonApiError(w, http.StatusUnauthorized, "not_authenticated", "Unauthenticated", "Sign in to access this resource.")
			return
		} else if !actor.IsValid() {
			restapi.WriteJsonApiError(w, http.StatusUnauthorized, "not_authenticated", "Unauthenticated", "Sign in to access this resource.")
			return
		}
		if !authzSvc.CanUploadTurnReports(actor) {
			log.Printf("%s %s: CanUploadTurnReports(%d): %v\n", r.Method, r.URL.Path, actor.ID, false)
			restapi.WriteJsonApiError(w, http.StatusForbidden, "forbidden", "Forbidden", "You must have the gm role to upload documents.")
			return
		}
		log.Printf("%s %s: actor %d\n", r.Method, r.URL.Path, actor.ID)

		gameId := domains.GameID(r.PathValue("id"))
		if gameId == domains.InvalidGameID {
			restapi.WriteJsonApiMalformedPathParameter(w, "game_id", "Game ID", r.PathValue("id"))
			return
		}
		log.Printf("%s %s: game %q\n", r.Method, r.URL.Path, gameId)

		r.Body = http.MaxBytesReader(w, r.Body, 150*1024)
		data, err := io.ReadAll(r.Body)
		if err != nil {
			if err.Error() == "http: request body too large" {
				restapi.WriteJsonApiError(w, http.StatusRequestEntityTooLarge, "too_large", "File Too Large", "File size exceeds 150KB limit.")
				return
			}
			restapi.WriteJsonApiError(w, http.StatusBadRequest, "bad_request", "Bad Request", "Error reading request body.")
			return
		}
		log.Printf("%s %s: game %q: data %d\n", r.Method, r.URL.Path, gameId, len(data))

		docx, err := parsers.ParseDocx(bytes.NewReader(data), true, true)
		log.Printf("%s %s: docx %v\n", r.Method, r.URL.Path, err)
		if err != nil {
			log.Printf("%s %s: game %q: ParseDocx %v\n", r.Method, r.URL.Path, gameId, err)
			if errors.Is(err, office.ErrNotAWordDocument) {
				restapi.WriteJsonApiError(w, http.StatusUnsupportedMediaType, "unsupported_media_type", "Unsupported file type", "Only Word documents (DOCX) are accepted.")
				return
			}
			restapi.WriteJsonApiError(w, http.StatusUnprocessableEntity, "invalid_file", "Invalid File", "Could not parse file: "+err.Error())
			return
		}
		header, err := parsers.ParseClanHeading(docx)
		log.Printf("%s %s: header %v\n", r.Method, r.URL.Path, err)
		if err != nil {
			// todo: add meta?
			// "meta": {"first-line": "TN3.1 Clan 9999", "expected-format": "GAME CLAN TURN"}
			restapi.WriteJsonApiError(w, http.StatusUnprocessableEntity, "invalid_clan_heading", "Invalid clan heading", "Could not parse a valid clan heading from the first two lines: "+err.Error())
			return
		}
		clanNo, err := strconv.Atoi(header.Id)
		log.Printf("%s %s: clanNo %d: %v\n", r.Method, r.URL.Path, clanNo, err)
		if err != nil {
			// todo: add meta?
			// "meta": {"first-line": "TN3.1 Clan 9999", "expected-format": "GAME CLAN TURN"}
			restapi.WriteJsonApiError(w, http.StatusUnprocessableEntity, "invalid_clan_heading", "Invalid Clan", "Clan ID is not a number: "+header.Id)
			return
		}

		// we must have a valid clan id to upload a document, so use the gameId and the clanNo from the header to find the clan
		clan, err := gamesSvc.GetClan(gameId, clanNo)
		log.Printf("%s %s: getClan(%q, %d) %v\n", r.Method, r.URL.Path, gameId, clanNo, err)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				restapi.WriteJsonApiError(w, http.StatusUnprocessableEntity, "clan_not_found", "Clan Not Found", fmt.Sprintf("game %q: clan %q: not found", gameId, header.Id))
				return
			}
			log.Printf("%s %s: FindClan: %v\n", r.Method, r.URL.Path, err)
			restapi.WriteJsonApiInternalServerError(w)
			return
		}

		// the documents service will calculate the size and hash for us, but
		// we are responsible for assigning the correct type and MIME type.
		doc := &domains.Document{
			Path:      fmt.Sprintf("%04d-%02d.%s.report.docx", header.Turn.Year, header.Turn.Month, header.Id),
			Type:      domains.TurnReportFile,
			MimeType:  domains.DOCXMimeType,
			Contents:  data,
			CanRead:   true, // Owner can read
			CanWrite:  false,
			CanDelete: false,
			CanShare:  true,
		}
		log.Printf("%s %s: doc %q\n", r.Method, r.URL.Path, doc.Path)

		docID, err := documentsSvc.CreateDocument(actor, clan, doc)
		log.Printf("%s %s: doc %d %v\n", r.Method, r.URL.Path, docID, err)
		if err != nil {
			// todo: handle unique constraint violation (idempotency)
			log.Printf("%s %s: CreateDocument: %v\n", r.Method, r.URL.Path, err)
			restapi.WriteJsonApiInternalServerError(w)
			return
		}

		view := struct {
			ID           string `jsonapi:"primary,document"`   // singular when sending a payload
			GameId       string `jsonapi:"attr,game-id"`       // game for this document
			ClanNo       string `jsonapi:"attr,clan"`          // clan for this document
			DocumentName string `jsonapi:"attr,document-name"` // untainted name of document
		}{
			ID:           fmt.Sprintf("%s", docID),
			GameId:       string(gameId),
			ClanNo:       fmt.Sprintf("%d", clan.ClanNo),
			DocumentName: doc.Path,
		}

		restapi.WriteJsonApiData(w, http.StatusCreated, &view)
	}
}
