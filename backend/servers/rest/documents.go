// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/playbymail/ottoapp/backend/auth"
	"github.com/playbymail/ottoapp/backend/domains"
	"github.com/playbymail/ottoapp/backend/restapi"
	"github.com/playbymail/ottoapp/backend/services/documents"
)

// CreateClanDocument creates a new clan-document association.
//
// Route: POST /api/clan-documents
func CreateClanDocument() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
	}
}

// DeleteClanDocument deletes a clan-document association.
//
// Route: DELETE /api/clan-documents/{id}
// Response type: http.StatusNoContent
func DeleteClanDocument() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
	}
}

// GetClanDocumentList returns a list of clan documents.
//
// Route: GET /api/clan-documents
// Response type: ???
func GetClanDocumentList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		restapi.WriteJsonApiResponse(w, http.StatusOK, []byte(`{
    "data": [
        {
            "type": "clan-document",
            "id": "11",
            "attributes": {
                "game": "0301",
                "clan-no": "0987",
                "document-name": "0301.0901-04.0987.report.docx",
                "document-type": "turn-report-file",
                "processing-status": "success",
                "created-at": "2025-11-19T21:33:57Z",
                "updated-at": "2025-11-19T21:33:57Z"
            },
            "links": {
                "self": "/api/documents/11",
                "contents": {
                    "href": "/api/documents/11/contents"
                },
                "log": {
                    "href": "/api/documents/12/contents"
                },
                "extract": {
                    "href": "/api/documents/13"
                }
            }
        },{
            "type": "clan-document",
            "id": "13",
            "attributes": {
                "game": "0301",
                "clan-no": "0987",
                "document-name": "0301.0901-04.0987.json",
                "document-type": "turn-report-extract",
                "processing-status": "success",
                "created-at": "2025-11-19T21:33:58Z",
                "updated-at": "2025-11-19T21:33:58Z"
            },
            "links": {
                "self": "/api/documents/13",
                "contents": {
                    "href": "/api/documents/13/contents"
                },
                "input": {
                    "href": "/api/documents/11"
                },
                "log": {
                    "href": "/api/documents/14/contents"
                },
                "output": {
                    "href": "/api/documents/15"
                }
            }
        },{
            "type": "clan-document",
            "id": "15",
            "attributes": {
                "game": "0301",
                "clan-no": "0987",
                "document-name": "0301.0901-04.0987.map",
                "document-type": "worldographer-map",
                "processing-status": "success",
                "created-at": "2025-11-19T21:33:59Z",
                "updated-at": "2025-11-19T21:33:59Z"
            },
            "links": {
                "self": "/api/documents/15",
                "contents": {
                    "href": "/api/documents/15/contents"
                },
                "input": {
                    "href": "/api/documents/13"
                }
            }
        },{
            "type": "clan-document",
            "id": "51",
            "attributes": {
                "game": "0301",
                "clan-no": "0987",
                "document-name": "0301.0901-05.0987.report.docx",
                "document-type": "turn-report-file",
                "processing-status": "success",
                "created-at": "2025-11-19T21:34:57Z",
                "updated-at": "2025-11-19T21:34:57Z"
            },
            "links": {
                "self": "/api/documents/51",
                "contents": {
                    "href": "/api/documents/51/contents"
                },
                "log": {
                    "href": "/api/documents/52"
                },
                "extract": {
                    "href": "/api/documents/53"
                }
            }
        },{
            "type": "clan-document",
            "id": "53",
            "attributes": {
                "game": "0301",
                "clan-no": "0987",
                "document-name": "0301.0901-04.0987.json",
                "document-type": "turn-report-extract",
                "processing-status": "failed",
                "created-at": "2025-11-19T21:35:58Z",
                "updated-at": "2025-11-19T21:35:58Z"
            },
            "links": {
                "self": "/api/documents/53",
                "input": {
                    "href": "/api/documents/51"
                },
                "log": {
                    "href": "/api/documents/54"
                }
            }
        }
    ]
}`))
	}
}

// GetClanDocument returns a single document.
//
// Route: GET /api/clan-documents/{id}
// Response type: ???
func GetClanDocument() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
	}
}

// GetDocument returns a document for the current actor.
//
// Route: GET /api/documents/{id}
//
// Response type: documents.DocumentView
func GetDocument(authSvc *auth.Service, documentsSvc *documents.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		actor, err := authSvc.GetActor(r)
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

		view, err := documentsSvc.GetDocument(actor, docId)
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
		b, _ := json.MarshalIndent(view, "", "  ")
		log.Printf("json %s\n", string(b))
		restapi.WriteJsonApiData(w, http.StatusOK, view)
	}
}

// GetDocumentContents returns the document contents for the current actor.
//
// Route: GET /api/documents/{id}/contents
//
// Response type: depends on the content type
func GetDocumentContents(authSvc *auth.Service, documentsSvc *documents.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		actor, err := authSvc.GetActor(r)
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

		doc, err := documentsSvc.GetDocumentContents(actor, docId)
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
		w.Header().Set("Content-Type", string(doc.MimeType))
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
func GetDocumentList(authSvc *auth.Service, documentsSvc *documents.Service) http.HandlerFunc {
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

		actor, err := authSvc.GetActor(r)
		if err != nil {
			log.Printf("%s %s: restapi: GetActor: %v\n", r.Method, r.URL.Path, err)
			restapi.WriteJsonApiError(w, http.StatusUnauthorized, "not_authenticated", "Unauthenticated", "Sign in to access this resource.")
			return
		} else if !actor.IsValid() {
			restapi.WriteJsonApiError(w, http.StatusUnauthorized, "not_authenticated", "Unauthenticated", "Sign in to access this resource.")
			return
		}
		//log.Printf("%s %s: actor %d\n", r.Method, r.URL.Path, actor.ID)

		view, err := documentsSvc.GetAllDocumentsForUserAcrossGames(actor, filterKind, pageNumber, pageSize)
		if err != nil {
			log.Printf("%s %s: restapi: GetDocuments: %v\n", r.Method, r.URL.Path, err)
			restapi.WriteJsonApiInternalServerError(w)
			return
		}
		restapi.WriteJsonApiData(w, http.StatusOK, view)
	}
}
