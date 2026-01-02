// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package rest

import "net/http"

func Routes(s *Server, quiet, verbose, debug bool) http.Handler {
	mux := http.NewServeMux()

	// Public routes (no authentication required)
	mux.HandleFunc("GET /api/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok","msg":"pong"}`))
	})
	mux.HandleFunc("POST /api/login", handlePostLogin(s.services.authnSvc, s.services.authzSvc, s.services.sessionsSvc, s.services.usersSvc))
	mux.HandleFunc("GET /api/session", handleGetSession(s.services.authzSvc, s.services.sessionsSvc))
	mux.HandleFunc("POST /api/shutdown", s.handlePostShutdown(s.debug.shutdownKey))
	mux.HandleFunc("GET /api/timezones", s.services.ianaSvc.HandleGetTimezones(true, false, false))
	mux.HandleFunc("GET /api/versions", s.getAllVersions())
	mux.HandleFunc("GET /api/versions/{id}", s.getVersions())

	// Protected routes (authentication required)
	protected := http.NewServeMux()
	protected.HandleFunc("GET /api/cookies/delete", s.services.sessionsSvc.DeleteCookie)
	protected.Handle("GET /api/documents", GetDocumentList(s.services.authzSvc, s.services.documentsSvc, quiet, verbose, debug))
	protected.Handle("GET /api/documents/{id}", GetDocument(s.services.authzSvc, s.services.documentsSvc, quiet, verbose, debug))
	protected.Handle("GET /api/documents/{id}/contents", GetDocumentContents(s.services.authzSvc, s.services.documentsSvc, quiet, verbose, debug))
	protected.Handle("POST /api/games/{id}/turn-report-files", PostGamesTurnReportFiles(s.services.authzSvc, s.services.documentsSvc, s.services.gamesSvc, quiet, verbose, debug))
	protected.HandleFunc("POST /api/logout", s.services.sessionsSvc.HandlePostLogout)
	protected.HandleFunc("GET /api/my/profile", handleGetMyProfile(s.services.authzSvc, s.services.usersSvc))
	protected.HandleFunc("GET /api/profile", handleGetProfile(s.services.authzSvc, s.services.usersSvc))
	protected.HandleFunc("POST /api/profile", handlePostProfile(s.services.authzSvc, s.services.ianaSvc, s.services.usersSvc))
	protected.HandleFunc("GET /api/users", handleGetUsers(s.services.authzSvc, s.services.usersSvc))
	protected.HandleFunc("POST /api/users", handlePostUser(s.services.authnSvc, s.services.authzSvc, s.services.usersSvc))
	protected.HandleFunc("GET /api/users/me", handleGetMe(s.services.authzSvc, s.services.usersSvc))
	protected.HandleFunc("GET /api/users/{id}", handleGetUser(s.services.authzSvc, s.services.usersSvc))
	protected.HandleFunc("PATCH /api/users/{id}", handlePatchUser(s.services.authzSvc, s.services.usersSvc))
	protected.HandleFunc("PATCH /api/users/{id}/password", handlePatchPassword(s.services.authnSvc, s.services.authzSvc))
	protected.HandleFunc("POST /api/users/{id}/reset-password", handlePostResetPassword(s.services.authnSvc, s.services.authzSvc))
	protected.HandleFunc("PATCH /api/users/{id}/role", handlePatchUserRole(s.services.authzSvc, s.services.usersSvc))

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
