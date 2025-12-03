// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package games

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/playbymail/ottoapp/backend/domains"
	"github.com/playbymail/ottoapp/backend/iana"
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

func New(db *sqlite.DB, authnSvc *authn.Service, authzSvc *authz.Service, usersSvc *users.Service) (*Service, error) {
	if authzSvc == nil {
		authzSvc = authz.New(db)
	}
	if authnSvc == nil {
		authnSvc = authn.New(db, authzSvc)
	}
	if usersSvc == nil {
		ianaSvc, err := iana.New(db)
		if err != nil {
			return nil, errors.Join(fmt.Errorf("new iana service"), err)
		}
		usersSvc = users.New(db, authnSvc, authzSvc, ianaSvc)
	}
	return &Service{db: db, authnSvc: authnSvc, authzSvc: authzSvc, usersSvc: usersSvc}, nil
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

func (s *Service) GetGamesList() ([]*domains.Game, error) {
	rows, err := s.db.Queries().GetGamesList(s.db.Context())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []*domains.Game{}, nil
		}
		return nil, err
	}
	var games []*domains.Game
	for _, row := range rows {
		games = append(games, &domains.Game{
			ID:          domains.GameID(row.GameID),
			Description: row.Description,
			IsActive:    row.IsActive,
		})
	}
	if len(games) == 0 {
		return []*domains.Game{}, nil
	}
	return games, nil
}

func (s *Service) ReadClansByGame(gameId domains.GameID, quiet, verbose, debug bool) ([]*domains.Clan, error) {
	rows, err := s.db.Queries().ReadClansByGame(s.db.Context(), string(gameId))
	if err != nil {
		return nil, err
	}
	var clans []*domains.Clan
	for _, row := range rows {
		clans = append(clans, &domains.Clan{
			GameID: row.GameID,
			UserID: domains.ID(row.UserID),
			ClanID: domains.ID(row.ClanID),
			ClanNo: int(row.Clan),
		})
	}
	if clans == nil {
		clans = []*domains.Clan{}
	}
	return clans, nil
}
