// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package rest

import (
	"fmt"
	"time"
)

type Option func(*Server) error

func WithGrace(d time.Duration) Option {
	return func(s *Server) error {
		if d < 0 {
			return fmt.Errorf("invalid grace timer")
		}
		s.graceTimer = d
		return nil
	}
}

func WithTimer(d time.Duration) Option {
	return func(s *Server) error {
		if d < 0 {
			return fmt.Errorf("invalid shutdown timer")
		}
		s.shutdownTimer = d
		return nil
	}
}
