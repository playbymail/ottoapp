// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package lemon

import (
	"regexp"

	"github.com/playbymail/ottoapp/backend/services/reports/lexers/tokens"
)

// an experiment to try lemon style parsing

func New(input []byte) *Lexer {
	if len(input) == 0 {
		return &Lexer{input: []byte{}, line: 0, col: 0, pos: 0}
	}
	return &Lexer{input: input, line: 1, col: 1, pos: 0}
}

type Lexer struct {
	input          []byte
	line, col, pos int
}

func (l *Lexer) Next() tokens.Token {
	if l.pos >= len(l.input) {
		return tokens.Token{Line: l.line, Col: l.col, Kind: tokens.EOF}
	}

	start, line, col := l.pos, l.line, l.col

	if isSpace[l.input[l.pos]] {
		for l.pos < len(l.input) && isSpace[l.input[l.pos]] {
			l.pos, l.col = l.pos+1, l.col+1
		}
		return tokens.Token{Line: line, Col: col, Kind: tokens.Spaces, Value: string(l.input[start:l.pos])}
	} else if l.input[l.pos] == '\n' {
		l.pos, l.line, l.col = l.pos+1, l.line+1, 1
		return tokens.Token{Line: line, Col: col, Kind: tokens.EOL, Value: string(l.input[start:l.pos])}
	} else if isDelimiter[l.input[l.pos]] {
		l.pos, l.col = l.pos+1, l.col+1
		return tokens.Token{Line: line, Col: col, Kind: tokens.Delimiter, Value: string(l.input[start:l.pos])}
	}

	// sniff out tokens that include delimiters. can be too greedy.
	if text := reDayMonthYear.Find(l.input[l.pos:]); text != nil {
		l.pos, l.col = l.pos+len(text), l.col+len(text)
		return tokens.Token{Line: line, Col: col, Kind: tokens.DayMonthYear, Value: string(l.input[start:l.pos])}
	} else if text = reNA.Find(l.input[l.pos:]); text != nil {
		l.pos, l.col = l.pos+len(text), l.col+len(text)
		return tokens.Token{Line: line, Col: col, Kind: tokens.NA, Value: string(l.input[start:l.pos])}
	} else if text = reTurnYearMonth.Find(l.input[l.pos:]); text != nil {
		l.pos, l.col = l.pos+len(text), l.col+len(text)
		return tokens.Token{Line: line, Col: col, Kind: tokens.TurnYearMonth, Value: string(l.input[start:l.pos])}
	}

	for l.pos < len(l.input) && !isDelimiter[l.input[l.pos]] {
		l.pos, l.col = l.pos+1, l.col+1
	}
	tok := tokens.Token{Line: line, Col: col, Kind: tokens.Text, Value: string(l.input[start:l.pos])}
	if kw, ok := keywords[tok.Value]; ok {
		tok.Kind = kw
	} else if reGrid.MatchString(tok.Value) {
		tok.Kind = tokens.Grid
	} else if reNumber.MatchString(tok.Value) {
		tok.Kind = tokens.Number
	} else if reUnitId.MatchString(tok.Value) {
		// warning - must be after the check for numbers since
		// the unit id for tribes is a valid four-digit number.
		tok.Kind = tokens.UnitId
	}
	return tok
}

func init() {
	for _, ch := range []byte{' ', '\t', '\r'} {
		isSpace[ch] = true
		isDelimiter[ch] = true
	}
	for _, ch := range []byte{0, '\n', '\'', '"', '.', ',', '(', ')', '#', '+', '-', '*', '/', '=', '\\', '$', ':'} {
		isDelimiter[ch] = true
	}
}

var (
	isDelimiter = [256]bool{}
	isSpace     = [256]bool{}

	reDayMonthYear  = regexp.MustCompile(`^\d{1,2}/\d{1,2}/\d{4}`)
	reGrid          = regexp.MustCompile(`^[A-Z]{2}$`)
	reNA            = regexp.MustCompile(`^N/A`)
	reNumber        = regexp.MustCompile(`^\d+$`)
	reTurnYearMonth = regexp.MustCompile(`^\d{3,4}-\d{1,2}`)
	reUnitId        = regexp.MustCompile(`^\d{4}([cefg][1-9])?$`)

	keywords = map[string]tokens.Kind{
		// unit line keywords
		"Tribe":    tokens.Tribe,
		"Courier":  tokens.Courier,
		"Element":  tokens.Element,
		"Garrison": tokens.Garrison,
		"Fleet":    tokens.Fleet,
		"Current":  tokens.Current,
		"Hex":      tokens.Hex,
		"Previous": tokens.Previous,

		// turn line keywords
		"Turn":   tokens.Turn,
		"Spring": tokens.Season,
		"Summer": tokens.Season,
		"Fall":   tokens.Season,
		"Winter": tokens.Season,
		"FINE":   tokens.Weather,
		"Next":   tokens.Next,

		// movement line keywords
		"Goes":     tokens.Goes,
		"to":       tokens.To,
		"Movement": tokens.Movement,
		"Move":     tokens.Move,

		// directions
		"N":  tokens.Direction,
		"NE": tokens.Direction, // conflicts with Grid
		"SE": tokens.Direction, // conflicts with Grid
		"S":  tokens.Direction,
		"SW": tokens.Direction, // conflicts with Grid and TerrainCode for Swamp
		"NW": tokens.Direction, // conflicts with Grid

		// terrain codes
		"ALPS": tokens.TerrainCode,
		"AH":   tokens.TerrainCode, // conflicts with Grid
		"AR":   tokens.TerrainCode, // conflicts with Grid
		"BF":   tokens.TerrainCode, // conflicts with Grid
		"BH":   tokens.TerrainCode, // conflicts with Grid
		"CH":   tokens.TerrainCode, // conflicts with Grid
		"D":    tokens.TerrainCode,
		"DH":   tokens.TerrainCode, // conflicts with Grid
		"DE":   tokens.TerrainCode, // conflicts with Grid
		"GH":   tokens.TerrainCode, // conflicts with Grid
		"GHP":  tokens.TerrainCode,
		"Hsm":  tokens.TerrainCode,
		"JG":   tokens.TerrainCode, // conflicts with Grid
		"JH":   tokens.TerrainCode, // conflicts with Grid
		"L":    tokens.TerrainCode,
		"Lam":  tokens.TerrainCode,
		"Lcm":  tokens.TerrainCode,
		"Ljm":  tokens.TerrainCode,
		"Lsm":  tokens.TerrainCode,
		"Lvm":  tokens.TerrainCode,
		"O":    tokens.TerrainCode,
		"PI":   tokens.TerrainCode, // conflicts with Grid
		"PPR":  tokens.TerrainCode,
		"PR":   tokens.TerrainCode, // conflicts with Grid
		"RH":   tokens.TerrainCode, // conflicts with Grid
		"SH":   tokens.TerrainCode, // conflicts with Grid
		// SW conflicts with Direction and Grid
		"TU": tokens.TerrainCode, // conflicts with Grid

		// scout line keywords
		"Scout": tokens.Scout,

		// status line keywords
		"Status": tokens.Status,
	}
)
