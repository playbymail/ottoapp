// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package sync

import (
	"log"

	"github.com/playbymail/ottoapp/backend/domains"
	"github.com/playbymail/ottoapp/backend/services/authn"
	"github.com/playbymail/ottoapp/backend/services/authz"
	"github.com/playbymail/ottoapp/backend/services/config"
	"github.com/playbymail/ottoapp/backend/services/documents"
	"github.com/playbymail/ottoapp/backend/services/games"
	"github.com/playbymail/ottoapp/backend/services/users"
	"github.com/playbymail/ottoapp/backend/stores/sqlite"
)

type Service struct {
	db           *sqlite.DB
	authnSvc     *authn.Service
	authzSvc     *authz.Service
	configSvc    *config.Service
	documentsSvc *documents.Service
	gameSvc      *games.Service
	usersSvc     *users.Service
}

func New(db *sqlite.DB, authnSvc *authn.Service, authzSvc *authz.Service, configSvc *config.Service, documentsSvc *documents.Service, gameSvc *games.Service, usersSvc *users.Service) (*Service, error) {
	if authnSvc == nil {
		log.Printf("sync: authnSvc is required\n")
		return nil, domains.ErrBadInput
	} else if authzSvc == nil {
		log.Printf("sync: authzSvc is required\n")
		return nil, domains.ErrBadInput
	} else if configSvc == nil {
		log.Printf("sync: configSvc is required\n")
		return nil, domains.ErrBadInput
	} else if documentsSvc == nil {
		log.Printf("sync: documentsSvc is required\n")
		return nil, domains.ErrBadInput
	} else if gameSvc == nil {
		log.Printf("sync: gameSvc is required\n")
		return nil, domains.ErrBadInput
	} else if usersSvc == nil {
		log.Printf("sync: usersSvc is required\n")
		return nil, domains.ErrBadInput
	}
	return &Service{
		db:           db,
		authnSvc:     authnSvc,
		authzSvc:     authzSvc,
		configSvc:    configSvc,
		documentsSvc: documentsSvc,
		gameSvc:      gameSvc,
		usersSvc:     usersSvc,
	}, nil
}
