// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package restapi

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/hashicorp/jsonapi"
)

// JSON:API helpers

func AbsURL(r *http.Request, path string) string {
	scheme := "https"
	if r.Header.Get("X-Forwarded-Proto") != "" {
		scheme = r.Header.Get("X-Forwarded-Proto")
	} else if r.TLS == nil {
		scheme = "http"
	}
	host := r.Header.Get("X-Forwarded-Host")
	if host == "" {
		host = r.Host
	}
	return scheme + "://" + host + path
}

// GetPaginationParameters query params:
//   - page[number]
//   - page[size]
func GetPaginationParameters(r *http.Request) (pageNumber, pageSize int, err error) {
	pageNumberParm := r.URL.Query().Get("page[number]")
	pageSizeParm := r.URL.Query().Get("page[size]")
	if pageNumberParm == "" && pageSizeParm == "" {
		return 0, 0, nil
	}
	if pageNumberParm != "" {
		pageNumber, err = strconv.Atoi(pageNumberParm)
		if err != nil {
			return 0, 0, errors.Join(ErrInvalidPageNumber, err)
		}
	} else {
		pageNumber = 1
	}
	if pageSizeParm != "" {
		pageSize, err = strconv.Atoi(pageSizeParm)
		if err != nil {
			return 0, 0, errors.Join(ErrInvalidPageSize, err)
		}
	} else {
		pageSize = 25
	}
	return pageNumber, pageSize, nil
}

func PaginateURL(r *http.Request, page, size int) string {
	u := *r.URL // copy
	q := u.Query()
	q.Set("page[number]", strconv.Itoa(page))
	q.Set("page[size]", strconv.Itoa(size))
	u.RawQuery = q.Encode()
	return AbsURL(r, u.Path+"?"+u.RawQuery)
}

func WriteJsonApiData(w http.ResponseWriter, status int, view any) {
	buf := &bytes.Buffer{}
	if err := jsonapi.MarshalPayload(buf, view); err != nil { // As an absolute last resort, write a minimal JSON:API error payload.
		// todo: log this!
		status = http.StatusInternalServerError // reset the status since we're not able to report the right error
		buf = &bytes.Buffer{}
		buf.WriteString(fmt.Sprintf(`{"errors":[{"status":"500","code":"encode_failed","title":"Encoding error","detail":%q}]}`, err.Error()))
	}
	WriteJsonApiResponse(w, status, buf.Bytes())
}

// WriteJsonApiError
//
//	Quick status-code guide (JSON:API)
//
//	  200/201: success with a data document
//
//	  204: success without a body (DELETE, or rare no-content updates)
//
//	  401/403/404/422/500: JSON:API error document ({ "errors": [ … ] })
//
// You should Map errors → statuses (classify the failure) and return the right status + error object(s):
//   - Bad payload / shape → 400 Bad Request
//   - AuthN/AuthZ → 401/403
//   - Target missing → 404 Not Found
//   - Validation failed (e.g., bad timezone, username too short) → 422 Unprocessable Entity
//   - Uniqueness/foreign-key constraint (e.g., email already taken) → usually 409 Conflict (some teams prefer 422; pick one and be consistent)
//   - Optimistic concurrency (ETag/If-Match mismatch) → 412 Precondition Failed
//   - DB unavailable / timeout / deadlock → 503 Service Unavailable (optionally Retry-After)
//   - Unknown / internal → 500 Internal Server Error
func WriteJsonApiError(w http.ResponseWriter, status int, code, title, detail string) {
	WriteJsonApiErrorObjects(w, status, &jsonapi.ErrorObject{
		Status: strconv.Itoa(status),
		Code:   code,
		Title:  title,
		Detail: detail,
	})
}

// WriteJsonApiErrorObjects writes a list of errors to the response.
func WriteJsonApiErrorObjects(w http.ResponseWriter, status int, errs ...*jsonapi.ErrorObject) {
	buf := &bytes.Buffer{}
	err := jsonapi.MarshalErrors(buf, errs)
	if err != nil { // As an absolute last resort, write a minimal JSON:API error payload.
		buf = &bytes.Buffer{}
		status = http.StatusInternalServerError
		buf.WriteString(fmt.Sprintf(`{"errors":[{"status":"500","code":"encode_failed","title":"Encoding error","detail":%q}]}`, err.Error()))
	}
	WriteJsonApiResponse(w, status, buf.Bytes())
}

func WriteJsonApiDatabaseError(w http.ResponseWriter) {
	WriteJsonApiError(w, http.StatusInternalServerError, "database_error", "Internal Server Error", "Could not process request due to a database error.")
}

// WriteJsonApiInvalidQueryParameter uses the error helper to quickly build
// a JSON:API error for an invalid URL query parameter.
func WriteJsonApiInvalidQueryParameter(w http.ResponseWriter, field, title string) {
	WriteJsonApiErrorObjects(w, http.StatusBadRequest, &jsonapi.ErrorObject{
		Status: strconv.Itoa(http.StatusBadRequest),
		Code:   fmt.Sprintf("invalid_query_%s", field),
		Title:  fmt.Sprintf("Invalid Parameter: %s", title),
		Detail: fmt.Sprintf("The value provided for the %q query parameter is invalid.", title),
		Source: &jsonapi.ErrorSource{
			Parameter: title,
		},
	})
}

func WriteJsonApiInternalServerError(w http.ResponseWriter, details ...string) {
	var detail string
	for _, s := range details {
		if detail != "" {
			detail += "\n"
		}
		detail += s
	}
	WriteJsonApiError(w, http.StatusInternalServerError, "server_error", "Internal Server Error", detail)
}

// WriteJsonApiMalformedPathParameter usage is
//
//	WriteJsonApiMalformedPathParameter(w, "document_id", "Document Identifier", value)
func WriteJsonApiMalformedPathParameter(w http.ResponseWriter, field, title, value string) {
	WriteJsonApiError(w, http.StatusBadRequest,
		fmt.Sprintf("invalid_%s", field),
		fmt.Sprintf("Invalid %s", title),
		fmt.Sprintf("The %s provided in the path ('%s') is malformed.", title, value))
}

// WriteJsonApiResponse adds the right content type and then writes the response.
func WriteJsonApiResponse(w http.ResponseWriter, status int, buf []byte) {
	w.Header().Set("Content-Type", jsonapi.MediaType)
	w.WriteHeader(status)
	_, _ = w.Write(buf)
}

// ValidationDetail is a 422 Unprocessable Entity Helper
type ValidationDetail struct {
	// FieldName is the human-readable name (e.g., "Document Name")
	FieldName string
	// JsonPointer is the JSON Pointer to the offending field (e.g., "/data/attributes/document-name")
	JsonPointer string
	// Detail is the specific reason the field failed (e.g., "must be at least 5 characters")
	Detail string
}

// WriteJsonApiValidationErrors is for unprocessable entities.
func WriteJsonApiValidationErrors(w http.ResponseWriter, details ...ValidationDetail) {
	const status = http.StatusUnprocessableEntity
	var list []*jsonapi.ErrorObject
	statusStr := strconv.Itoa(status)

	for _, detail := range details {
		list = append(list, &jsonapi.ErrorObject{
			Status: statusStr,
			Code:   "validation_failed",
			Title:  fmt.Sprintf("Invalid Field: %s", detail.FieldName),
			Detail: detail.Detail,
			Source: &jsonapi.ErrorSource{
				Pointer: detail.JsonPointer, // KEY: Use Source.Pointer for request body errors
			},
		})
	}
	if list == nil {
		panic("assert(details != nil)")
	}

	// Delegate final writing to the core helper
	WriteJsonApiErrorObjects(w, status, list...)
}
