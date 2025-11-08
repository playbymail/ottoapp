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

	"github.com/playbymail/ottoapp/backend/sessions"
)

type Server struct {
	http.Server
	services struct {
		sessionsSvc *sessions.Service
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

func New(sessionsSvc *sessions.Service, options ...Option) (*Server, error) {
	s := &Server{
		Server: http.Server{
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
		},
	}
	s.network.scheme, s.network.host, s.network.port = "http", "localhost", "8181"
	s.services.sessionsSvc = sessionsSvc
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
