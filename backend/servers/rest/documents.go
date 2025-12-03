// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/playbymail/ottoapp/backend/domains"
	"github.com/playbymail/ottoapp/backend/restapi"
	"github.com/playbymail/ottoapp/backend/services/authz"
	"github.com/playbymail/ottoapp/backend/services/documents"
)

// GetDocument returns a document for the current actor.
//
// Route: GET /api/documents/{id}
//
// Response type: documents.DocumentView
func GetDocument(authzSvc *authz.Service, documentsSvc *documents.Service, quiet, verbose, debug bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		actor, err := authzSvc.GetActor(r)
		if err != nil {
			log.Printf("%s %s: restapi: GetActor: %v\n", r.Method, r.URL.Path, err)
			restapi.WriteJsonApiError(w, http.StatusUnauthorized, "not_authenticated", "Unauthenticated", "Sign in to access this resource.")
			return
		} else if !actor.IsValid() {
			restapi.WriteJsonApiError(w, http.StatusUnauthorized, "not_authenticated", "Unauthenticated", "Sign in to access this resource.")
			return
		}

		var docId domains.ID = domains.InvalidID
		if value, err := strconv.Atoi(r.PathValue("id")); err != nil {
			restapi.WriteJsonApiMalformedPathParameter(w, "document_id", "Document ID", r.PathValue("id"))
			return
		} else {
			docId = domains.ID(value)
		}

		clan, err := documentsSvc.ReadDocumentOwner(docId, quiet, verbose, debug)
		if err != nil {
			if errors.Is(err, domains.ErrNotExists) {
				// not found, return a 404 response structured as a JSON:API error object
				restapi.WriteJsonApiError(w, http.StatusNotFound, "document_not_found",
					"Resource Not Found",
					fmt.Sprintf("Document with ID %d could not be found.", docId))
				return
			}
			restapi.WriteJsonApiDatabaseError(w)
			return
		} else if clan.UserID != actor.ID {
			restapi.WriteJsonApiError(w, http.StatusForbidden, "forbidden", "Forbidden", "You are not allowed access to this document.")
			return
		}

		view, err := documentsSvc.ReadDocument(actor, clan, docId, quiet, verbose, debug)
		if err != nil {
			log.Printf("%s %s: restapi: GetDocument: %v\n", r.Method, r.URL.Path, err)
			restapi.WriteJsonApiDatabaseError(w)
			return
		} else if view == nil {
			// not found, return a 404 response structured as a JSON:API error object
			restapi.WriteJsonApiError(w, http.StatusNotFound, "document_not_found",
				"Resource Not Found",
				fmt.Sprintf("Document with ID %d could not be found.", docId))
			return
		}
		if debug {
			b, _ := json.MarshalIndent(view, "", "  ")
			log.Printf("json %s\n", string(b))
		}
		restapi.WriteJsonApiData(w, http.StatusOK, view)
	}
}

// GetDocumentContents returns the document contents for the current actor.
//
// Route: GET /api/documents/{id}/contents
//
// Response type: depends on the content type
func GetDocumentContents(authzSvc *authz.Service, documentsSvc *documents.Service, quiet, verbose, debug bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		actor, err := authzSvc.GetActor(r)
		if err != nil {
			log.Printf("%s %s: restapi: GetActor: %v\n", r.Method, r.URL.Path, err)
			restapi.WriteJsonApiError(w, http.StatusUnauthorized, "not_authenticated", "Unauthenticated", "Sign in to access this resource.")
			return
		} else if !actor.IsValid() {
			restapi.WriteJsonApiError(w, http.StatusUnauthorized, "not_authenticated", "Unauthenticated", "Sign in to access this resource.")
			return
		}

		var docId domains.ID = domains.InvalidID
		if value, err := strconv.Atoi(r.PathValue("id")); err != nil {
			restapi.WriteJsonApiMalformedPathParameter(w, "document_id", "Document ID", r.PathValue("id"))
			return
		} else {
			docId = domains.ID(value)
		}

		clan, err := documentsSvc.ReadDocumentOwner(docId, quiet, verbose, debug)
		if err != nil {
			if errors.Is(err, domains.ErrNotExists) {
				// not found, return a 404 response structured as a JSON:API error object
				restapi.WriteJsonApiError(w, http.StatusNotFound, "document_not_found",
					"Resource Not Found",
					fmt.Sprintf("Document with ID %d could not be found.", docId))
				return
			}
			restapi.WriteJsonApiDatabaseError(w)
			return
		} else if clan.UserID != actor.ID {
			restapi.WriteJsonApiError(w, http.StatusForbidden, "forbidden", "Forbidden", "You are not allowed access to this document.")
			return
		}

		doc, err := documentsSvc.ReadDocumentContents(actor, clan, docId, quiet, verbose, debug)
		if err != nil {
			log.Printf("%s %s: restapi: GetDocumentContents: %v\n", r.Method, r.URL.Path, err)
			restapi.WriteJsonApiDatabaseError(w)
			return
		} else if doc == nil {
			// not found, return a 404 response structured as a JSON:API error object
			restapi.WriteJsonApiError(w, http.StatusNotFound, "document_not_found",
				"Resource Not Found",
				fmt.Sprintf("Document with ID %d could not be found.", docId))
			return
		}
		w.Header().Set("Content-Type", doc.ContentType)
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", doc.Path))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(doc.Contents)
	}
}

// GetDocumentList returns a list of documents for the current actor.
//
// Route: GET /api/documents
// Query params:
//   - filter[kind]=worldographer-map – only maps
//   - page[number], page[size] – standard pagination if needed
//
// Response type: []documents.DocumentView
func GetDocumentList(authzSvc *authz.Service, documentsSvc *documents.Service, quiet, verbose, debug bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filterKind := domains.DocumentType(r.URL.Query().Get("filter[kind]"))
		//log.Printf("%s %s: filterKind %q: %v\n", r.Method, r.URL.Path, filterKind, filterKind.IsValid())
		if filterKind != "" && !filterKind.IsValid() {
			restapi.WriteJsonApiInvalidQueryParameter(w, "filter_kind", "filter[kind]")
			return
		}
		pageNumber, pageSize, err := restapi.GetPaginationParameters(r)
		if err != nil {
			if errors.Is(err, restapi.ErrInvalidPageNumber) {
				restapi.WriteJsonApiInvalidQueryParameter(w, "page_number", "page[number]")
			} else if errors.Is(err, restapi.ErrInvalidPageSize) {
				restapi.WriteJsonApiInvalidQueryParameter(w, "page_size", "page[size]")
			} else {
				log.Printf("%s %s: restapi: GetPaginationParameters: %v\n", r.Method, r.URL.Path, err)
				restapi.WriteJsonApiInternalServerError(w)
			}
			return
		}

		actor, err := authzSvc.GetActor(r)
		if err != nil {
			log.Printf("%s %s: restapi: GetActor: %v\n", r.Method, r.URL.Path, err)
			restapi.WriteJsonApiError(w, http.StatusUnauthorized, "not_authenticated", "Unauthenticated", "Sign in to access this resource.")
			return
		} else if !actor.IsValid() {
			restapi.WriteJsonApiError(w, http.StatusUnauthorized, "not_authenticated", "Unauthenticated", "Sign in to access this resource.")
			return
		}
		//log.Printf("%s %s: actor %d\n", r.Method, r.URL.Path, actor.ID)
		userId := actor.ID

		view, err := documentsSvc.ReadDocumentsByUser(actor, userId, filterKind, pageNumber, pageSize, quiet, verbose, debug)
		if err != nil {
			log.Printf("%s %s: restapi: GetDocuments: %v\n", r.Method, r.URL.Path, err)
			restapi.WriteJsonApiInternalServerError(w)
			return
		}
		restapi.WriteJsonApiData(w, http.StatusOK, view)
	}
}
