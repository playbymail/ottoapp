// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package games

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/playbymail/ottoapp/backend/domains"
	"github.com/playbymail/ottoapp/backend/iana"
	"github.com/playbymail/ottoapp/backend/services/authn"
	"github.com/playbymail/ottoapp/backend/services/authz"
	"github.com/playbymail/ottoapp/backend/services/users"
	"github.com/playbymail/ottoapp/backend/stores/sqlite"
	"github.com/playbymail/ottoapp/backend/stores/sqlite/sqlc"
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

func (s *Service) ReadClanByGameIdAndClanNo(gameId domains.GameID, clanNo int, quiet, verbose, debug bool) (*domains.Clan, error) {
	clan, err := s.db.Queries().GetClanByGameClanNo(s.db.Context(), sqlc.GetClanByGameClanNoParams{
		GameID: int64(gameId),
		ClanNo: int64(clanNo),
	})
	if err != nil {
		return nil, err
	}
	return &domains.Clan{
		GameID:   domains.GameID(clan.GameID),
		UserID:   domains.ID(clan.UserID),
		ClanID:   domains.ID(clan.ClanID),
		ClanNo:   int(clan.Clan),
		IsActive: clan.IsActive,
	}, nil
}

func (s *Service) ReadClanByGameIdAndUserId(gameId domains.GameID, userId domains.ID) (*domains.Clan, error) {
	clan, err := s.db.Queries().GetClanByGameUser(s.db.Context(), sqlc.GetClanByGameUserParams{
		GameID: int64(gameId),
		UserID: int64(userId),
	})
	if err != nil {
		return nil, err
	}
	return &domains.Clan{
		GameID:   domains.GameID(clan.GameID),
		UserID:   domains.ID(clan.UserID),
		ClanID:   domains.ID(clan.ClanID),
		ClanNo:   int(clan.Clan),
		IsActive: clan.IsActive,
	}, nil
}

func (s *Service) ReadGames() ([]*domains.Game, error) {
	rows, err := s.db.Queries().ReadGames(s.db.Context())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []*domains.Game{}, nil
		}
		return nil, err
	}
	var games []*domains.Game
	for _, row := range rows {
		game := &domains.Game{
			ID:          domains.GameID(row.GameID),
			Code:        row.Code,
			Description: row.Description,
			IsActive:    row.IsActive,
			ActiveTurn: &domains.Turn{
				ID:        domains.TurnID(row.ActiveTurnID),
				Year:      int(row.ActiveTurnYear),
				Month:     int(row.ActiveTurnMonth),
				No:        int(row.ActiveTurnNo),
				OrdersDue: time.Time{},
			},
		}
		games = append(games, game)
	}
	if len(games) == 0 {
		return []*domains.Game{}, nil
	}
	return games, nil
}

func (s *Service) ReadClansByGame(gameId domains.GameID, quiet, verbose, debug bool) ([]*domains.Clan, error) {
	rows, err := s.db.Queries().ReadClansByGame(s.db.Context(), int64(gameId))
	if err != nil {
		return nil, err
	}
	var clans []*domains.Clan
	for _, row := range rows {
		clans = append(clans, &domains.Clan{
			GameID: domains.GameID(row.GameID),
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

func (s *Service) GameIdClanNoToClan(id domains.GameID, clanNo int) (*domains.Clan, error) {
	return nil, domains.ErrNotImplemented
}

func (s *Service) GameCodeYearMonthToGameTurnId(code, yearMonth string) (domains.GameID, domains.TurnID, error) {
	return domains.InvalidGameID, domains.InvalidTurnID, domains.ErrNotImplemented
}
