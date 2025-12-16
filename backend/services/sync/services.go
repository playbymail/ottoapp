// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package sync

import (
	"log"

	"github.com/playbymail/ottoapp/backend/domains"
	"github.com/playbymail/ottoapp/backend/services/authn"
	"github.com/playbymail/ottoapp/backend/services/config"
	"github.com/playbymail/ottoapp/backend/services/users"
	"github.com/playbymail/ottoapp/backend/stores/sqlite"
)

type Service struct {
	db        *sqlite.DB
	authnSvc  *authn.Service
	configSvc *config.Service
	usersSvc  *users.Service
}

func New(db *sqlite.DB, authnSvc *authn.Service, configSvc *config.Service, usersSvc *users.Service) (*Service, error) {
	if authnSvc == nil {
		log.Printf("sync: authnSvc is required\n")
		return nil, domains.ErrBadInput
	} else if configSvc == nil {
		log.Printf("sync: configSvc is required\n")
		return nil, domains.ErrBadInput
	} else if usersSvc == nil {
		log.Printf("sync: usersSvc is required\n")
		return nil, domains.ErrBadInput
	}
	return &Service{
		db:        db,
		authnSvc:  authnSvc,
		configSvc: configSvc,
		usersSvc:  usersSvc,
	}, nil
}
