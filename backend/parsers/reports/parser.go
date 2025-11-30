// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package reports

//go:generate pigeon -o grammar.go grammar.peg
//go:generate pigeon -o currentturn/grammar.go currentturn/grammar.peg

type ReportFile_t struct {
	Clan *ClanSection_t
	Turn *TurnInfo_t
}

type ClanSection_t struct {
	Id          UnitId_t
	CurrentHex  Coords_t
	PreviousHex Coords_t
	CurrentTurn *TurnInfo_t
	ScoutLines  []*ScoutLine_t
}

type LocationLine_t struct {
	UnitID      UnitId_t
	CurrentHex  Coords_t
	PreviousHex Coords_t
}

type UnitId_t string

type Coords_t struct {
	Grid string
	Col  int
	Row  int
}

func (c Coords_t) IsNA() bool {
	return c.Grid == "N/A"
}
func (c Coords_t) IsObscured() bool {
	return c.Grid == "##"
}

type TurnInfo_t struct {
	Year  int
	Month int
	No    int
}

type ScoutLine_t struct {
	Id ScoutId_t
}

type ScoutId_t int

// these types are internal to the grammar

type gSection struct {
	location   *gLocation
	turn       *gTurnLine
	scoutLines []*gScoutLine
}

type gCoords struct {
	grid string
	col  int
	row  int
}

type gLocation struct {
	unitId   gUnitId
	note     string
	current  gCoords
	previous gCoords
}

type gScoutId int

type gScoutLine struct {
	Id    gScoutId      `json:"id"`
	Moves []*gScoutMove `json:"moves,omitempty"`
}

type gScoutMove struct {
	Move string `json:"move,omitempty"`
	// moved results
	Moved *gScoutMoved `json:"moved,omitempty"`
	// failed to move results
	// - Can't Move on
	// - Not enough M.P's to Move to
	Failed *gScoutFailed `json:"failed,omitempty"`
	// patrolled results
	Patrolled *gScoutPatrolled `json:"patrolled,omitempty"`
}

type gScoutMoved struct {
	Move      string `json:"move,omitempty"`
	Direction string `json:"direction"`
	Terrain   string `json:"terrain"`
	resources string
	edges     string
	neighbors string
}

type gScoutFailed struct {
	Move           string `json:"failed,omitempty"`
	Direction      string `json:"direction,omitempty"`
	Terrain        string `json:"terrain,omitempty"`
	BlockedByLake  bool   `json:"blockedByLake,omitempty"`
	BlockedByOcean bool   `json:"blockedByOcean,omitempty"`
	BlockedByRiver bool   `json:"blockedByRiver,omitempty"`
	NotEnoughMPs   bool   `json:"notEnoughMPs,omitempty"`
}

type gScoutPatrolled struct {
	Move              string `json:"move,omitempty"`
	NothingOfInterest bool   `json:"nothingOfInterest,omitempty"`
}

type gTurn struct {
	year  int
	month int
	no    int
}

type gTurnLine struct {
	current  *gTurn
	next     *gTurn
	season   string
	weather  string
	reportDt string
}

type gUnitId string

func toAnySlice(v any) []any {
	if v == nil {
		return nil
	}
	return v.([]any)
}
