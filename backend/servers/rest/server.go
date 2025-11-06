// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package rest

import (
	"context"
	"fmt"
	"log"
	"net"
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
	csrfGuard bool
	logRoutes bool
	channels  struct {
		graceTimer    time.Duration
		shutdown      chan os.Signal
		shutdownTimer time.Duration
	}
	network struct {
		scheme string
		host   string
		port   string
	}
	debug struct {
		debug       bool
		shutdownKey []byte
	}
}

func New(sm ssi.SessionManager, options ...Option) (*Server, error) {
	s := &Server{
		Server: http.Server{
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
		},
	}
	s.network.scheme, s.network.host, s.network.port = "http", "localhost", "8181"
	s.services.sessionManager = sm
	s.debug.debug = true

	for _, opt := range options {
		err := opt(s)
		if err != nil {
			return nil, err
		}
	}

	log.Printf("[server] host %q: port %q\n", s.network.host, s.network.port)

	s.Addr = net.JoinHostPort(s.network.host, s.network.port)
	log.Printf("[server] address %q\n", s.Addr)

	s.Handler = Routes(s)

	return s, nil
}

func (s *Server) Run() error {
	log.Printf("[rest] serving API on %q\n", fmt.Sprintf("%s://%s", s.network.scheme, s.Addr))
	if s.channels.shutdownTimer != 0 {
		log.Printf("[rest] server timeout %v\n", s.channels.shutdownTimer)
	}
	if s.channels.graceTimer != 0 {
		log.Printf("[rest] shutdown delay %v\n", s.channels.graceTimer)
	}

	serverErrors := make(chan error, 1)
	go func() {
		serverErrors <- s.ListenAndServe()
	}()

	s.channels.shutdown = make(chan os.Signal, 1)
	signal.Notify(s.channels.shutdown, syscall.SIGINT, syscall.SIGTERM)

	if s.channels.shutdownTimer > 0 {
		go func() {
			time.Sleep(s.channels.shutdownTimer)
			log.Printf("[rest] timeout reached (%v), initiating shutdown\n", s.channels.shutdownTimer)
			s.channels.shutdown <- syscall.SIGTERM
		}()
	}

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case sig := <-s.channels.shutdown:
		log.Printf("[rest] received signal %v, starting graceful shutdown\n", sig)

		ctx, cancel := context.WithTimeout(context.Background(), s.channels.graceTimer)
		defer cancel()

		if err := s.Shutdown(ctx); err != nil {
			log.Printf("[rest] error during shutdown: %v\n", err)
			return s.Close()
		}

		log.Println("[rest] server stopped gracefully")
	}

	return nil
}

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

	mux.HandleFunc("GET /api/cookies/delete", s.services.sessionManager.DeleteCookie)
	mux.HandleFunc("POST /api/login", s.services.sessionManager.PostLoginHandler)
	mux.HandleFunc("POST /api/logout", s.services.sessionManager.PostLogoutHandler)
	mux.HandleFunc("GET /api/me", s.services.sessionManager.GetMeHandler)
	mux.HandleFunc("GET /api/session", s.services.sessionManager.GetSessionHandler) // returns CSRF
	mux.HandleFunc("POST /api/shutdown", s.postShutdown(s.debug.shutdownKey))

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
