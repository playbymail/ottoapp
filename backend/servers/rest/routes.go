// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package rest

import "net/http"

func Routes(s *Server) http.Handler {
	mux := http.NewServeMux()

	// Public routes (no authentication required)
	mux.HandleFunc("GET /api/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok","msg":"pong"}`))
	})
	mux.HandleFunc("POST /api/login", HandlePostLogin(s.services.authnSvc, s.services.authzSvc, s.services.sessionsSvc, s.services.usersSvc))
	mux.HandleFunc("GET /api/session", HandleGetSession(s.services.authzSvc, s.services.sessionsSvc))
	mux.HandleFunc("POST /api/shutdown", s.handlePostShutdown(s.debug.shutdownKey))
	mux.HandleFunc("GET /api/timezones", s.services.tzSvc.HandleGetTimezones())
	mux.HandleFunc("GET /api/versions", s.getAllVersions())
	mux.HandleFunc("GET /api/versions/{id}", s.getVersions())

	// Protected routes (authentication required)
	protected := http.NewServeMux()
	protected.HandleFunc("GET /api/cookies/delete", s.services.sessionsSvc.DeleteCookie)
	protected.Handle("GET /api/documents", GetDocumentList(s.services.authzSvc, s.services.documentsSvc))
	protected.Handle("GET /api/documents/{id}", GetDocument(s.services.authzSvc, s.services.documentsSvc))
	protected.Handle("GET /api/documents/{id}/contents", GetDocumentContents(s.services.authzSvc, s.services.documentsSvc))
	protected.Handle("POST /api/games/{id}/turn-report-files", PostGamesTurnReportFiles(s.services.authzSvc, s.services.documentsSvc, s.services.gamesSvc))
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
	//h = s.sessionMiddleware(h)
	h = s.services.sessionsSvc.Middleware(h)

	//// Add logging middleware
	//if s.logRoutes {
	//	h = s.loggingMiddleware(h)
	//}

	return h
}
