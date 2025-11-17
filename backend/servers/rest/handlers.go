// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package rest

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"syscall"
	"time"

	"github.com/hashicorp/jsonapi"
	"github.com/playbymail/ottoapp/backend/restapi"
	"github.com/playbymail/ottoapp/backend/versions"
)

// GET /api/versions
func (s *Server) getAllVersions() http.HandlerFunc {
	view := s.services.versionsSvc.Version()
	views := []*versions.VersionView{&view}
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		restapi.WriteJsonApiData(w, http.StatusOK, views)
	}
}

// GET /api/versions/{id}
func (s *Server) getVersions() http.HandlerFunc {
	view := s.services.versionsSvc.Version()
	views := map[string]*versions.VersionView{view.ID: &view}
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		versionId := r.PathValue("id")
		view, ok := views[versionId]
		if !ok {
			restapi.WriteJsonApiErrorObjects(w, http.StatusNotFound, &jsonapi.ErrorObject{
				Status: strconv.Itoa(http.StatusNotFound),
				Code:   "unknown_version_id",
				Title:  "Unknown VersionID",
				Detail: "Provide a valid VersionID.",
				Source: &jsonapi.ErrorSource{
					Parameter: "id",
				},
			})
		}

		restapi.WriteJsonApiData(w, http.StatusOK, view)
	}
}

func (s *Server) handlePostShutdown(key []byte) http.HandlerFunc {
	if len(key) == 0 {
		return func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
			return
		}
	}
	return func(w http.ResponseWriter, r *http.Request) {
		// Basic server preconditions
		if s.debug.shutdownKey == nil {
			// do not reveal "disabled" vs "wrong key" — keep response generic
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		if s.channels.shutdown == nil {
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		}

		// Content-Type check (allow "application/json; charset=...")
		ct := r.Header.Get("Content-Type")
		if ct == "" || len(ct) < len("application/json") || ct[:len("application/json")] != "application/json" {
			http.Error(w, "content-type must be application/json", http.StatusUnsupportedMediaType)
			return
		}

		// Decode payload
		var payload struct {
			Key string `json:"key"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			// bad JSON — return quickly (no need to delay)
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		if payload.Key == "" {
			http.Error(w, "missing key", http.StatusBadRequest)
			return
		}

		// compare keys in psuedo-constant time
		delayTimer := time.NewTimer(250 * time.Millisecond)
		defer func() {
			// ensure the timer is drained to avoid memory pressure.
			// not sure if this is required with Go 1.23+
			if !delayTimer.Stop() {
				<-delayTimer.C
			}
		}()
		// Compare keys in constant time using SHA-256 digests
		da := sha256.Sum256([]byte(payload.Key))
		ok := bytes.Equal(key[:], da[:])
		if !ok {
			<-delayTimer.C // delay the return to prevent a timing attack
			// log only that a failed attempt happened (no key, no length)
			log.Printf("%s %s: shutdown api rejected invalid key\n", r.Method, r.URL.Path)
			// respond with a generic forbidden status (don't reveal which part failed)
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		// key is valid — accept and initiate shutdown
		log.Printf("%s %s: shutdown api accepted; initiating shutdown\n", r.Method, r.URL.Path)

		// non-blocking send so handler can't hang if receiver isn't ready
		select {
		case s.channels.shutdown <- syscall.SIGTERM:
			// signalled successfully
		default:
			// receiver not ready — return service unavailable so caller can retry
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(struct {
			Status  string `json:"status"`
			Message string `json:"msg"`
		}{
			Status:  "ok",
			Message: "shutdown initiated",
		})
	}
}
