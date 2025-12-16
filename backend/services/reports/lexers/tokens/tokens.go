// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package tokens

//go:generate stringer --type Kind

import "fmt"

type Token struct {
	Kind  Kind
	Line  int
	Col   int
	Value string
}

func (tok Token) String() string {
	return fmt.Sprintf("<%d:%d %s(%q)>", tok.Line, tok.Col, tok.Kind, tok.Value)
}

type Kind int

const (
	Text Kind = iota

	Backslash
	Colon
	Comma
	Dash
	DayMonthYear
	Delimiter
	Direction
	Equals
	EOL
	EOF
	Grid
	Hash
	NA
	Note
	Number
	LeftParen
	RightParen
	ScoutId
	Slash
	Spaces
	TerrainCode
	TurnYearMonth
	UnitId

	// unit line keywords
	Tribe
	Courier
	Element
	Fleet
	Garrison
	Current
	Hex
	Previous
	Turn

	// turn line keywords
	Season
	Next
	Weather

	// unit movement keywords
	Goes
	To
	Movement
	Move

	// scout movement keywords
	Scout

	// status line keywords
	Status
)
