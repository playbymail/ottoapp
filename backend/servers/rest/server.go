// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package rest

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	ssi "github.com/playbymail/ottoapp/backend/services/sessions"
)

type Server struct {
	http.Server
	services struct {
		sessionManager ssi.SessionManager
	}
	csrfGuard     bool
	graceTimer    time.Duration
	logRoutes     bool
	shutdownTimer time.Duration
	debug         struct {
		autoLogin bool
		debug     bool
	}
}

func New(sm ssi.SessionManager, options ...Option) (*Server, error) {
	s := &Server{
		Server: http.Server{
			Addr:         ":8181",
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
		},
	}
	s.services.sessionManager = sm
	s.debug.autoLogin = true
	s.debug.debug = true

	for _, opt := range options {
		err := opt(s)
		if err != nil {
			return nil, err
		}
	}

	s.Handler = Routes(s)

	return s, nil
}

func (s *Server) Run() error {
	log.Printf("[rest] serving API on %q\n", fmt.Sprintf("http://localhost%s", s.Addr))
	if s.shutdownTimer != 0 {
		log.Printf("[rest] server timeout %v\n", s.shutdownTimer)
	}
	if s.graceTimer != 0 {
		log.Printf("[rest] shutdown delay %v\n", s.graceTimer)
	}

	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("server listening on port %s\n", s.Addr)
		serverErrors <- s.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	if s.shutdownTimer > 0 {
		go func() {
			time.Sleep(s.shutdownTimer)
			log.Printf("[rest] timeout reached (%v), initiating shutdown\n", s.shutdownTimer)
			shutdown <- syscall.SIGTERM
		}()
	}

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case sig := <-shutdown:
		log.Printf("[rest] received signal %v, starting graceful shutdown\n", sig)

		ctx, cancel := context.WithTimeout(context.Background(), s.graceTimer)
		defer cancel()

		if err := s.Shutdown(ctx); err != nil {
			log.Printf("[rest] error during shutdown: %v\n", err)
			return s.Close()
		}

		log.Println("[rest] server stopped gracefully")
	}

	return nil
}

func Routes(s *Server) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok","msg":"pong"}`))
	})

	mux.HandleFunc("GET /api/cookies/delete", s.services.sessionManager.DeleteCookie)
	mux.HandleFunc("POST /api/login", s.services.sessionManager.PostLoginHandler)
	mux.HandleFunc("POST /api/logout", s.services.sessionManager.PostLogoutHandler)
	mux.HandleFunc("GET /api/me", s.services.sessionManager.GetMeHandler)
	mux.HandleFunc("GET /api/session", s.services.sessionManager.GetSessionHandler) // returns CSRF

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
