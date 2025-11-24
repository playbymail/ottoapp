// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package domains

type GameID string // something like 0300, 0301

const (
	InvalidGameID GameID = ""
)

type TurnNo string // YYYY-MM

type Clan struct {
	GameID      string
	UserID      ID
	ClanID      ID
	ClanNo      int
	SetupTurnNo TurnNo
	IsActive    bool
}
