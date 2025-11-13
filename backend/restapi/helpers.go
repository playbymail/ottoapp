// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package restapi

import (
	"bytes"
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
	err := jsonapi.MarshalPayload(buf, view)
	if err != nil { // As an absolute last resort, write a minimal JSON:API error payload.
		buf = &bytes.Buffer{}
		status = http.StatusInternalServerError
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

func WriteJsonApiErrorObjects(w http.ResponseWriter, status int, errs ...*jsonapi.ErrorObject) {
	var list []*jsonapi.ErrorObject
	for _, err := range errs {
		list = append(list, err)
	}
	buf := &bytes.Buffer{}
	err := jsonapi.MarshalErrors(buf, list)
	if err != nil { // As an absolute last resort, write a minimal JSON:API error payload.
		buf = &bytes.Buffer{}
		status = http.StatusInternalServerError
		buf.WriteString(fmt.Sprintf(`{"errors":[{"status":"500","code":"encode_failed","title":"Encoding error","detail":%q}]}`, err.Error()))
	}
	WriteJsonApiResponse(w, status, buf.Bytes())
}

func WriteJsonApiResponse(w http.ResponseWriter, status int, buf []byte) {
	w.Header().Set("Content-Type", jsonapi.MediaType)
	w.WriteHeader(status)
	_, _ = w.Write(buf)
}
