// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package rest

import (
	"log"
	"net/http"
	"time"
)

/*** middleware & helpers ***/

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !s.logRoutes {
			next.ServeHTTP(w, r)
			return
		}
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		log.Printf("%s %s - %d - %v", r.Method, r.URL.Path, rec.status, time.Since(start))
	})
}

func authOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s: authOnly: entered\n", r.Method, r.URL.Path)
		if _, ok := currentSession(r); !ok {
			log.Printf("%s %s: authOnly: currentSession: false\n", r.Method, r.URL.Path)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
