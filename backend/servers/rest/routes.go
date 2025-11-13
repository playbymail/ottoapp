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

	// Public routes (no authentication required)
	mux.HandleFunc("GET /api/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok","msg":"pong"}`))
	})
	mux.HandleFunc("POST /api/login", s.services.sessionsSvc.HandlePostLogin)
	mux.HandleFunc("GET /api/session", s.services.sessionsSvc.HandleGetSession)
	mux.HandleFunc("GET /api/timezones", s.handleGetTimezones)
	mux.HandleFunc("GET /api/version", s.getVersion)
	mux.HandleFunc("POST /api/shutdown", s.handlePostShutdown(s.debug.shutdownKey))

	// Protected routes (authentication required)
	protected := http.NewServeMux()
	protected.HandleFunc("GET /api/cookies/delete", s.services.sessionsSvc.DeleteCookie)
	protected.Handle("GET /api/documents", s.handleGetDocuments())
	protected.HandleFunc("POST /api/logout", s.services.sessionsSvc.HandlePostLogout)
	protected.HandleFunc("GET /api/my/profile", s.services.usersSvc.HandleGetMyProfile)
	protected.HandleFunc("GET /api/profile", s.handleGetProfile)
	protected.HandleFunc("POST /api/profile", s.handlePostProfile)
	protected.HandleFunc("GET /api/users", s.services.usersSvc.HandleGetUsers)
	protected.HandleFunc("POST /api/users", s.services.usersSvc.HandlePostUser)
	protected.HandleFunc("GET /api/users/me", s.services.usersSvc.HandleGetMe)
	protected.HandleFunc("GET /api/users/{id}", s.services.usersSvc.HandleGetUser)
	protected.HandleFunc("PATCH /api/users/{id}", s.services.usersSvc.HandlePatchUser)
	protected.HandleFunc("PATCH /api/users/{id}/password", s.services.usersSvc.HandlePatchPassword)
	protected.HandleFunc("POST /api/users/{id}/reset-password", s.services.usersSvc.HandlePostResetPassword)
	protected.HandleFunc("PATCH /api/users/{id}/role", s.services.usersSvc.HandlePatchUserRole)

	// Apply auth middleware to protected routes
	mux.Handle("/api/", authOnly(protected))

	// convert mux to handler before we add any global middlewares
	var h http.Handler = mux

	// Protect all state-changing routes with CSRF:
	if s.csrfGuard {
		h = csrfOnly(h)
	}

	// Add session middleware (runs on all routes)
	h = s.sessionMiddleware(h)

	// Add logging middleware
	h = s.loggingMiddleware(h)

	return h
}
