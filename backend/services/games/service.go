// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package games

import (
	"strconv"
	"time"

	"github.com/playbymail/ottoapp/backend/auth"
	"github.com/playbymail/ottoapp/backend/domains"
	"github.com/playbymail/ottoapp/backend/restapi"
	"github.com/playbymail/ottoapp/backend/stores/sqlite"
	"github.com/playbymail/ottoapp/backend/stores/sqlite/sqlc"
	"github.com/playbymail/ottoapp/backend/users"
)

type Service struct {
	db       *sqlite.DB
	authSvc  *auth.Service
	usersSvc *users.Service
}

func New(db *sqlite.DB, authSvc *auth.Service, usersSvc *users.Service) *Service {
	return &Service{db: db, authSvc: authSvc, usersSvc: usersSvc}
}

func (s *Service) GetClan(game string, clanNo int) (*domains.Clan, error) {
	clan, err := s.db.Queries().GetClanByGameClanNo(s.db.Context(), sqlc.GetClanByGameClanNoParams{
		GameID: game,
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

func (s *Service) GetClanForUser(game string, userId domains.ID) (*domains.Clan, error) {
	clan, err := s.db.Queries().GetClanByGameUser(s.db.Context(), sqlc.GetClanByGameUserParams{
		GameID: game,
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

func (s *Service) NewClanDocumentResource(
	clanDocID int64,
	gameID string,
	clanID int64,
	turnNo *int64,
	kind string,
	documentID int64,
) ClanDocumentResource {
	createdAt, updatedAt := time.Now().UTC(), time.Now().UTC()
	return ClanDocumentResource{
		Type: "clan-documents",
		ID:   strconv.FormatInt(clanDocID, 10),
		Attributes: ClanDocumentAttributes{
			GameID:    gameID,
			ClanID:    clanID,
			TurnNo:    turnNo,
			Kind:      kind,
			CreatedAt: createdAt.Format(time.RFC3339),
			UpdatedAt: updatedAt.Format(time.RFC3339),
		},
		Relationships: ClanDocumentRelationships{
			Document: &restapi.ToOneRelationship{
				Data: &restapi.ResourceIdentifier{
					Type: "documents",
					ID:   strconv.FormatInt(documentID, 10),
				},
			},
		},
	}
}
