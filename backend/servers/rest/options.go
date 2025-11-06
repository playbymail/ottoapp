// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package rest

import (
	"crypto/sha256"
	"fmt"
	"log"
	"time"
)

type Option func(*Server) error

func WithCsrfGuard(csrfGuard bool) Option {
	return func(s *Server) error {
		s.csrfGuard = csrfGuard
		return nil
	}
}

func WithGrace(d time.Duration) Option {
	return func(s *Server) error {
		if d < 0 {
			return fmt.Errorf("invalid grace timer")
		}
		s.channels.graceTimer = d
		return nil
	}
}

func WithHost(host string) Option {
	return func(s *Server) error {
		s.network.host = host
		return nil
	}
}

func WithPort(port string) Option {
	return func(s *Server) error {
		s.network.port = port
		return nil
	}
}

func WithRouteLogging(logRoutes bool) Option {
	return func(s *Server) error {
		s.logRoutes = logRoutes
		return nil
	}
}

func WithShutdownKey(key string) Option {
	log.Printf("option: withShutdownKey(%q)\n", key)
	return func(s *Server) error {
		b := sha256.Sum256([]byte(key))
		s.debug.shutdownKey = append([]byte{}, b[:]...)
		return nil
	}
}

func WithTimer(d time.Duration) Option {
	return func(s *Server) error {
		if d < 0 {
			return fmt.Errorf("invalid shutdown timer")
		}
		s.channels.shutdownTimer = d
		return nil
	}
}
