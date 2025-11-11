// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package rest

import "net/http"

/*
Sessions:what are the 4 routes for?

  POST /api/login → create session + Set-Cookie

  POST /api/logout → delete session + clear cookie

  GET /api/session → “is this cookie valid? give me csrf + user”

  GET /api/me → (optional) “give me just user again”
*/

func Routes(s *Server) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok","msg":"pong"}`))
	})

	mux.HandleFunc("GET /api/cookies/delete", s.services.sessionsSvc.DeleteCookie)
	mux.Handle("GET /api/documents", s.handleGetDocuments())
	mux.HandleFunc("POST /api/login", s.services.sessionsSvc.HandlePostLogin)
	mux.HandleFunc("POST /api/logout", s.services.sessionsSvc.HandlePostLogout)
	mux.HandleFunc("GET /api/me", s.services.sessionsSvc.HandleGetMe)
	mux.HandleFunc("GET /api/profile", s.handleGetProfile)
	mux.HandleFunc("POST /api/profile", s.handlePostProfile)
	mux.HandleFunc("GET /api/session", s.services.sessionsSvc.HandleGetSession) // returns CSRF
	mux.HandleFunc("POST /api/shutdown", s.handlePostShutdown(s.debug.shutdownKey))
	mux.HandleFunc("GET /api/timezones", s.handleGetTimezones)

	// User management routes
	mux.HandleFunc("GET /api/users/me", s.services.usersSvc.HandleGetMe)
	mux.HandleFunc("GET /api/users", s.services.usersSvc.HandleGetUsers)
	mux.HandleFunc("POST /api/users", s.services.usersSvc.HandlePostUser)
	mux.HandleFunc("GET /api/users/{id}", s.services.usersSvc.HandleGetUser)
	mux.HandleFunc("PATCH /api/users/{id}", s.services.usersSvc.HandlePatchUser)
	mux.HandleFunc("PUT /api/users/{id}/password", s.services.usersSvc.HandlePutPassword)
	mux.HandleFunc("POST /api/users/{id}/reset-password", s.services.usersSvc.HandlePostResetPassword)
	mux.HandleFunc("PATCH /api/users/{id}/role", s.services.usersSvc.HandlePatchUserRole)

	mux.HandleFunc("GET /api/version", s.getVersion)

	// convert mux to handler before we add any global middlewares
	var h http.Handler = mux

	// Protect all state-changing routes with CSRF:
	if s.csrfGuard {
		h = csrfOnly(h)
	}

	// Add logging middleware
	h = s.loggingMiddleware(h)

	return h
}
