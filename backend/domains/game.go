// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package domains

import (
	"fmt"
	"time"
)

type GameID int // database key?

const (
	InvalidGameID GameID = 0
)

type Game struct {
	ID          GameID
	Code        string // something like 0300, 0301
	Description string
	IsActive    bool
	ActiveTurn  *Turn
	SetupTurn   *Turn
	OrdersDue   time.Time
}

type Turn struct {
	ID    string // YYYY-MM
	Year  int    // 899...9999
	Month int    // 1...12
	No    int    // 0...9_999_999
}

func (t Turn) String() string {
	return fmt.Sprintf("%04d-%02d", t.Year, t.Month)
}

const (
	InvalidTurnID = "0000-00"
)

type Clan struct {
	GameID    GameID
	UserID    ID
	ClanID    ID
	ClanNo    int
	SetupTurn *Turn
	IsActive  bool
}
