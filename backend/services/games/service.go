// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package games

import (
	"github.com/playbymail/ottoapp/backend/domains"
	"github.com/playbymail/ottoapp/backend/services/authn"
	"github.com/playbymail/ottoapp/backend/services/authz"
	"github.com/playbymail/ottoapp/backend/stores/sqlite"
	"github.com/playbymail/ottoapp/backend/stores/sqlite/sqlc"
	"github.com/playbymail/ottoapp/backend/users"
)

type Service struct {
	db       *sqlite.DB
	authnSvc *authn.Service
	authzSvc *authz.Service
	usersSvc *users.Service
}

func New(db *sqlite.DB, authnSvc *authn.Service, authzSvc *authz.Service, usersSvc *users.Service) *Service {
	return &Service{db: db, authnSvc: authnSvc, authzSvc: authzSvc, usersSvc: usersSvc}
}

func (s *Service) GetClan(gameId domains.GameID, clanNo int) (*domains.Clan, error) {
	clan, err := s.db.Queries().GetClanByGameClanNo(s.db.Context(), sqlc.GetClanByGameClanNoParams{
		GameID: string(gameId),
		ClanNo: int64(clanNo),
	})
	if err != nil {
		return nil, err
	}
	return &domains.Clan{
		GameID:   clan.GameID,
		UserID:   domains.ID(clan.UserID),
		ClanID:   domains.ID(clan.ClanID),
		ClanNo:   int(clan.Clan),
		IsActive: clan.IsActive,
	}, nil
}

func (s *Service) GetClanForUser(gameId domains.GameID, userId domains.ID) (*domains.Clan, error) {
	clan, err := s.db.Queries().GetClanByGameUser(s.db.Context(), sqlc.GetClanByGameUserParams{
		GameID: string(gameId),
		UserID: int64(userId),
	})
	if err != nil {
		return nil, err
	}
	return &domains.Clan{
		GameID:   clan.GameID,
		UserID:   domains.ID(clan.UserID),
		ClanID:   domains.ID(clan.ClanID),
		ClanNo:   int(clan.Clan),
		IsActive: clan.IsActive,
	}, nil
}
