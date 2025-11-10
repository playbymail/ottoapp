// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/playbymail/ottoapp/backend/auth"
	"github.com/playbymail/ottoapp/backend/domains"
	"github.com/playbymail/ottoapp/backend/iana"
	"github.com/playbymail/ottoapp/backend/sessions"
	"github.com/playbymail/ottoapp/backend/users"
)

type Server struct {
	http.Server
	services struct {
		authSvc     *auth.Service
		sessionsSvc *sessions.Service
		tzSvc       *iana.Service
		usersSvc    *users.Service
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

func New(authSvc *auth.Service, sessionsSvc *sessions.Service, tzSvc *iana.Service, usersSvc *users.Service, options ...Option) (*Server, error) {
	s := &Server{
		Server: http.Server{
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
		},
	}
	s.network.scheme, s.network.host, s.network.port = "http", "localhost", "8181"
	s.debug.debug = true
	s.services.authSvc = authSvc
	s.services.sessionsSvc = sessionsSvc
	s.services.tzSvc = tzSvc
	s.services.usersSvc = usersSvc

	for _, opt := range options {
		err := opt(s)
		if err != nil {
			return nil, err
		}
	}

	if s.services.authSvc == nil {
		log.Printf("[rest] authSvc not initialized")
		return nil, domains.ErrInvalidArgument
	}
	if s.services.sessionsSvc == nil {
		log.Printf("[rest] sessionsSvc not initialized")
		return nil, domains.ErrInvalidArgument
	}
	if s.services.tzSvc == nil {
		log.Printf("[rest] tzSvc not initialized")
		return nil, domains.ErrInvalidArgument
	}
	if s.services.usersSvc == nil {
		log.Printf("[rest] usersSvc not initialized")
		return nil, domains.ErrInvalidArgument
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

// helpers for encoding and decoding
// from https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years/

func encode[T any](w http.ResponseWriter, r *http.Request, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

func decode[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}
	return v, nil
}

func decodeValid[T Validator](r *http.Request) (T, map[string]string, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, nil, fmt.Errorf("decode json: %w", err)
	}
	if problems := v.Valid(r.Context()); len(problems) > 0 {
		return v, problems, fmt.Errorf("invalid %T: %d problems", v, len(problems))
	}
	return v, nil, nil
}

// Validator is an object that can be validated.
type Validator interface {
	// Valid checks the object and returns any
	// problems. If len(problems) == 0 then
	// the object is valid.
	Valid(ctx context.Context) (problems map[string]string)
}
