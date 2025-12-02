// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package cst implements a parser for scrubbed TribeNet reports. It returns
// a concrete syntax tree (CST). It tries hard to parse as much of the input
// as it can, with error reporting included in the nodes.
package cst

import (
	"bytes"
	"fmt"

	"github.com/playbymail/ottoapp/backend/services/reports/lexers"
)

// Parse parses the token stream and returns a CST.
func Parse(tokens []*lexers.Token) *TurnReportNode {
	p := &parser{
		tokens: tokens,
		pos:    0,
	}
	return p.parseTurnReport()
}

// ParseTurnLine parses a single turn line from the token stream.
func ParseTurnLine(tokens []*lexers.Token) *TurnLineNode {
	p := &parser{
		tokens: tokens,
		pos:    0,
	}
	return p.parseTurnLine()
}

// parser holds the state for parsing.
type parser struct {
	tokens []*lexers.Token
	pos    int // current position in tokens
}

// peek returns the current token without advancing.
// Returns nil if at end of tokens.
func (p *parser) peek() *lexers.Token {
	if p.pos >= len(p.tokens) {
		return nil
	}
	return p.tokens[p.pos]
}

// peekKind returns the Kind of the current token.
// Returns EOF if at end of tokens.
func (p *parser) peekKind() lexers.Kind {
	if tok := p.peek(); tok != nil && tok.Value != nil {
		return tok.Kind
	}
	return lexers.EOF
}

// peekMatch returns true if the Kind of the current token matches.
func (p *parser) peekMatch(kinds ...lexers.Kind) (lexers.Kind, bool) {
	if tok := p.peek(); tok != nil && tok.Value != nil {
		for _, kind := range kinds {
			if tok.Kind == kind {
				return tok.Kind, true
			}
		}
	}
	return lexers.EOF, false
}

// advance consumes and returns the current token.
// Returns nil if at end of tokens.
func (p *parser) advance() *lexers.Token {
	if p.pos >= len(p.tokens) {
		return nil
	}
	tok := p.tokens[p.pos]
	p.pos++
	return tok
}

// match checks if the current token matches any of the given kinds.
// If it matches, it advances and returns the token.
// Otherwise, it returns nil without advancing.
func (p *parser) match(kinds ...lexers.Kind) *lexers.Token {
	kind := p.peekKind()
	for _, k := range kinds {
		if kind == k {
			return p.advance()
		}
	}
	return nil
}

// expect consumes a token of the given kind.
// If the current token doesn't match, it returns nil and an error.
func (p *parser) expect(kind lexers.Kind) (*lexers.Token, error) {
	if tok := p.match(kind); tok != nil {
		return tok, nil
	}
	got := p.peekKind()
	return nil, fmt.Errorf("expected %s, got %s", kind, got)
}

// isAtEnd returns true if all tokens have been consumed.
func (p *parser) isAtEnd() bool {
	return p.pos >= len(p.tokens) || p.peekKind() == lexers.EOF
}

// isUnitKeyword returns true if the given kind is a unit keyword.
func isUnitKeyword(k lexers.Kind) bool {
	switch k {
	case lexers.Courier, lexers.Element, lexers.Fleet, lexers.Garrison, lexers.Tribe:
		return true
	}
	return false
}

// syncToNextLine advances past the next EOL token.
// Used for error recovery in line-oriented parsing.
func (p *parser) syncToNextLine() []*lexers.Token {
	var skipped []*lexers.Token
	for !p.isAtEnd() {
		tok := p.advance()
		skipped = append(skipped, tok)
		if tok.Value != nil && tok.Kind == lexers.EOL {
			break
		}
	}
	return skipped
}

// syncToUnitKeyword advances until a unit keyword or EOF is found.
// The unit keyword is NOT consumed.
// Used for error recovery between unit sections.
func (p *parser) syncToUnitKeyword() []*lexers.Token {
	var skipped []*lexers.Token
	for !p.isAtEnd() {
		if isUnitKeyword(p.peekKind()) {
			break
		}
		skipped = append(skipped, p.advance())
	}
	return skipped
}

// parseTurnReport parses the entire token stream into a TurnReportNode.
// turn_report = unit_section, { unit_section }, EOF ;
func (p *parser) parseTurnReport() *TurnReportNode {
	node := &TurnReportNode{}

	// Parse one or more unit sections
	for !p.isAtEnd() {
		// Check if we're at a unit keyword
		if isUnitKeyword(p.peekKind()) {
			section := p.parseUnitSection()
			node.Sections = append(node.Sections, section)
			node.tokens = append(node.tokens, section.tokens...)
			node.errors = append(node.errors, section.errors...)
		} else {
			// Not a unit keyword - skip to next unit or EOF
			skipped := p.syncToUnitKeyword()
			if len(skipped) > 0 {
				node.tokens = append(node.tokens, skipped...)
				node.errors = append(node.errors, fmt.Errorf("unexpected tokens before unit section"))
			}
		}
	}

	return node
}

// parseUnitSection parses a unit section.
// unit_section = unit_line ;
func (p *parser) parseUnitSection() *UnitSectionNode {
	node := &UnitSectionNode{}

	unitLine := p.parseUnitLine()
	node.UnitLine = unitLine
	node.tokens = unitLine.tokens
	node.errors = unitLine.errors

	return node
}

// parseUnitLine parses a unit line.
// unit_line = unit_keyword, unit_id, Comma, [ Note ], Comma,
//
//	Current, Hex, Equals, coords, Comma,
//	LeftParen, Previous, Hex, Equals, coords, RightParen, EOL ;
func (p *parser) parseUnitLine() *UnitLineNode {
	node := &UnitLineNode{}

	// unit_keyword
	if tok := p.match(lexers.Courier, lexers.Element, lexers.Fleet, lexers.Garrison, lexers.Tribe); tok != nil {
		node.Keyword = tok
		node.tokens = append(node.tokens, tok)
	} else {
		node.errors = append(node.errors, fmt.Errorf("expected unit keyword"))
		skipped := p.syncToNextLine()
		node.tokens = append(node.tokens, skipped...)
		return node
	}

	// unit_id (Number or UnitId)
	if tok := p.match(lexers.Number, lexers.UnitId); tok != nil {
		node.UnitID = tok
		node.tokens = append(node.tokens, tok)
	} else {
		node.errors = append(node.errors, fmt.Errorf("expected unit id"))
		skipped := p.syncToNextLine()
		node.tokens = append(node.tokens, skipped...)
		return node
	}

	// Comma
	if tok, err := p.expect(lexers.Comma); err != nil {
		node.errors = append(node.errors, err)
		skipped := p.syncToNextLine()
		node.tokens = append(node.tokens, skipped...)
		return node
	} else {
		node.Comma1 = tok
		node.tokens = append(node.tokens, tok)
	}

	// Optional Note - all tokens to first comma or end of line
	var noteTokens []*lexers.Token
	for _, ok := p.peekMatch(lexers.Comma, lexers.EOL); !ok; _, ok = p.peekMatch() {
		tok := p.advance()
		node.tokens = append(node.tokens, tok)
	}
	if len(noteTokens) != 0 {
		node.Note = lexers.Merge(lexers.Note, noteTokens...)
	}

	// Comma
	if tok, err := p.expect(lexers.Comma); err != nil {
		node.errors = append(node.errors, err)
		skipped := p.syncToNextLine()
		node.tokens = append(node.tokens, skipped...)
		return node
	} else {
		node.Comma2 = tok
		node.tokens = append(node.tokens, tok)
	}

	// Current
	if tok, err := p.expect(lexers.Current); err != nil {
		node.errors = append(node.errors, err)
		skipped := p.syncToNextLine()
		node.tokens = append(node.tokens, skipped...)
		return node
	} else {
		node.Current = tok
		node.tokens = append(node.tokens, tok)
	}

	// Hex
	if tok, err := p.expect(lexers.Hex); err != nil {
		node.errors = append(node.errors, err)
		skipped := p.syncToNextLine()
		node.tokens = append(node.tokens, skipped...)
		return node
	} else {
		node.Hex1 = tok
		node.tokens = append(node.tokens, tok)
	}

	// Equals
	if tok, err := p.expect(lexers.Equals); err != nil {
		node.errors = append(node.errors, err)
		skipped := p.syncToNextLine()
		node.tokens = append(node.tokens, skipped...)
		return node
	} else {
		node.Equals1 = tok
		node.tokens = append(node.tokens, tok)
	}

	// coords (current hex)
	node.CurrentHex = p.parseCoords()
	node.tokens = append(node.tokens, node.CurrentHex.Tokens()...)
	node.errors = append(node.errors, node.CurrentHex.Errors()...)

	// Comma
	if tok, err := p.expect(lexers.Comma); err != nil {
		node.errors = append(node.errors, err)
		skipped := p.syncToNextLine()
		node.tokens = append(node.tokens, skipped...)
		return node
	} else {
		node.Comma3 = tok
		node.tokens = append(node.tokens, tok)
	}

	// LeftParen
	if tok, err := p.expect(lexers.LeftParen); err != nil {
		node.errors = append(node.errors, err)
		skipped := p.syncToNextLine()
		node.tokens = append(node.tokens, skipped...)
		return node
	} else {
		node.LeftParen = tok
		node.tokens = append(node.tokens, tok)
	}

	// Previous
	if tok, err := p.expect(lexers.Previous); err != nil {
		node.errors = append(node.errors, err)
		skipped := p.syncToNextLine()
		node.tokens = append(node.tokens, skipped...)
		return node
	} else {
		node.Previous = tok
		node.tokens = append(node.tokens, tok)
	}

	// Hex
	if tok, err := p.expect(lexers.Hex); err != nil {
		node.errors = append(node.errors, err)
		skipped := p.syncToNextLine()
		node.tokens = append(node.tokens, skipped...)
		return node
	} else {
		node.Hex2 = tok
		node.tokens = append(node.tokens, tok)
	}

	// Equals
	if tok, err := p.expect(lexers.Equals); err != nil {
		node.errors = append(node.errors, err)
		skipped := p.syncToNextLine()
		node.tokens = append(node.tokens, skipped...)
		return node
	} else {
		node.Equals2 = tok
		node.tokens = append(node.tokens, tok)
	}

	// coords (previous hex)
	node.PreviousHex = p.parseCoords()
	node.tokens = append(node.tokens, node.PreviousHex.Tokens()...)
	node.errors = append(node.errors, node.PreviousHex.Errors()...)

	// RightParen
	if tok, err := p.expect(lexers.RightParen); err != nil {
		node.errors = append(node.errors, err)
		skipped := p.syncToNextLine()
		node.tokens = append(node.tokens, skipped...)
		return node
	} else {
		node.RightParen = tok
		node.tokens = append(node.tokens, tok)
	}

	// EOL
	if tok, err := p.expect(lexers.EOL); err != nil {
		node.errors = append(node.errors, err)
		skipped := p.syncToNextLine()
		node.tokens = append(node.tokens, skipped...)
		return node
	} else {
		node.EOL = tok
		node.tokens = append(node.tokens, tok)
	}

	return node
}

// parseTurnLine parses a turn line.
// turn_line = Current, Turn, TurnYearMonth, turn_number, Comma, Season, Comma, Weather,
//
//	[ Next, Turn, TurnYearMonth, turn_number, Comma, report_date ],
//	EOL ;
func (p *parser) parseTurnLine() *TurnLineNode {
	node := &TurnLineNode{}

	// Current
	if tok, err := p.expect(lexers.Current); err != nil {
		node.errors = append(node.errors, err)
		skipped := p.syncToNextLine()
		node.tokens = append(node.tokens, skipped...)
		return node
	} else {
		node.Current1 = tok
		node.tokens = append(node.tokens, tok)
	}

	// Turn
	if tok, err := p.expect(lexers.Turn); err != nil {
		node.errors = append(node.errors, err)
		skipped := p.syncToNextLine()
		node.tokens = append(node.tokens, skipped...)
		return node
	} else {
		node.Turn1 = tok
		node.tokens = append(node.tokens, tok)
	}

	// TurnYearMonth
	if tok, err := p.expect(lexers.TurnYearMonth); err != nil {
		node.errors = append(node.errors, err)
		skipped := p.syncToNextLine()
		node.tokens = append(node.tokens, skipped...)
		return node
	} else {
		node.TurnYearMonth1 = tok
		node.tokens = append(node.tokens, tok)
	}

	// turn_number
	turnNum := p.parseTurnNumber()
	node.TurnNumber1 = turnNum
	node.tokens = append(node.tokens, turnNum.tokens...)
	node.errors = append(node.errors, turnNum.errors...)

	// Comma
	if tok, err := p.expect(lexers.Comma); err != nil {
		node.errors = append(node.errors, err)
		skipped := p.syncToNextLine()
		node.tokens = append(node.tokens, skipped...)
		return node
	} else {
		node.Comma1 = tok
		node.tokens = append(node.tokens, tok)
	}

	// Season
	if tok, err := p.expect(lexers.Season); err != nil {
		node.errors = append(node.errors, err)
		skipped := p.syncToNextLine()
		node.tokens = append(node.tokens, skipped...)
		return node
	} else {
		node.Season = tok
		node.tokens = append(node.tokens, tok)
	}

	// Comma
	if tok, err := p.expect(lexers.Comma); err != nil {
		node.errors = append(node.errors, err)
		skipped := p.syncToNextLine()
		node.tokens = append(node.tokens, skipped...)
		return node
	} else {
		node.Comma2 = tok
		node.tokens = append(node.tokens, tok)
	}

	// Weather
	if tok, err := p.expect(lexers.Weather); err != nil {
		node.errors = append(node.errors, err)
		skipped := p.syncToNextLine()
		node.tokens = append(node.tokens, skipped...)
		return node
	} else {
		node.Weather = tok
		node.tokens = append(node.tokens, tok)
	}

	// Optional: Next Turn TurnYearMonth turn_number Comma report_date
	if p.peekKind() == lexers.Next {
		node.Next = p.advance()
		node.tokens = append(node.tokens, node.Next)

		// Turn
		if tok, err := p.expect(lexers.Turn); err != nil {
			node.errors = append(node.errors, err)
			skipped := p.syncToNextLine()
			node.tokens = append(node.tokens, skipped...)
			return node
		} else {
			node.Turn2 = tok
			node.tokens = append(node.tokens, tok)
		}

		// TurnYearMonth
		if tok, err := p.expect(lexers.TurnYearMonth); err != nil {
			node.errors = append(node.errors, err)
			skipped := p.syncToNextLine()
			node.tokens = append(node.tokens, skipped...)
			return node
		} else {
			node.TurnYearMonth2 = tok
			node.tokens = append(node.tokens, tok)
		}

		// turn_number
		turnNum2 := p.parseTurnNumber()
		node.TurnNumber2 = turnNum2
		node.tokens = append(node.tokens, turnNum2.tokens...)
		node.errors = append(node.errors, turnNum2.errors...)

		// Comma
		if tok, err := p.expect(lexers.Comma); err != nil {
			node.errors = append(node.errors, err)
			skipped := p.syncToNextLine()
			node.tokens = append(node.tokens, skipped...)
			return node
		} else {
			node.Comma3 = tok
			node.tokens = append(node.tokens, tok)
		}

		// report_date
		reportDate := p.parseReportDate()
		node.ReportDate = reportDate
		node.tokens = append(node.tokens, reportDate.tokens...)
		node.errors = append(node.errors, reportDate.errors...)
	}

	// EOL
	if tok, err := p.expect(lexers.EOL); err != nil {
		node.errors = append(node.errors, err)
		skipped := p.syncToNextLine()
		node.tokens = append(node.tokens, skipped...)
		return node
	} else {
		node.EOL = tok
		node.tokens = append(node.tokens, tok)
	}

	return node
}

// parseTurnNumber parses a turn number.
// turn_number = LeftParen, Hash, Number, RightParen ;
func (p *parser) parseTurnNumber() *TurnNumberNode {
	node := &TurnNumberNode{}

	// LeftParen
	if tok, err := p.expect(lexers.LeftParen); err != nil {
		node.errors = append(node.errors, err)
		return node
	} else {
		node.LeftParen = tok
		node.tokens = append(node.tokens, tok)
	}

	// Hash
	if tok, err := p.expect(lexers.Hash); err != nil {
		node.errors = append(node.errors, err)
		return node
	} else {
		node.Hash = tok
		node.tokens = append(node.tokens, tok)
	}

	// Number
	if tok, err := p.expect(lexers.Number); err != nil {
		node.errors = append(node.errors, err)
		return node
	} else {
		node.Number = tok
		node.tokens = append(node.tokens, tok)
	}

	// RightParen
	if tok, err := p.expect(lexers.RightParen); err != nil {
		node.errors = append(node.errors, err)
		return node
	} else {
		node.RightParen = tok
		node.tokens = append(node.tokens, tok)
	}

	return node
}

// parseReportDate parses a report date.
// report_date = Number, Slash, Number, Slash, Number ;
func (p *parser) parseReportDate() *ReportDateNode {
	node := &ReportDateNode{}

	// Day (Number)
	if tok, err := p.expect(lexers.Number); err != nil {
		node.errors = append(node.errors, err)
		return node
	} else {
		node.Day = tok
		node.tokens = append(node.tokens, tok)
	}

	// Slash
	if tok, err := p.expect(lexers.Slash); err != nil {
		node.errors = append(node.errors, err)
		return node
	} else {
		node.Slash1 = tok
		node.tokens = append(node.tokens, tok)
	}

	// Month (Number)
	if tok, err := p.expect(lexers.Number); err != nil {
		node.errors = append(node.errors, err)
		return node
	} else {
		node.Month = tok
		node.tokens = append(node.tokens, tok)
	}

	// Slash
	if tok, err := p.expect(lexers.Slash); err != nil {
		node.errors = append(node.errors, err)
		return node
	} else {
		node.Slash2 = tok
		node.tokens = append(node.tokens, tok)
	}

	// Year (Number)
	if tok, err := p.expect(lexers.Number); err != nil {
		node.errors = append(node.errors, err)
		return node
	} else {
		node.Year = tok
		node.tokens = append(node.tokens, tok)
	}

	return node
}

// parseCoords parses coordinates.
// coords = grid_coords | na_coords | obscured_coords ;
// grid_coords = Grid, Number ;
// na_coords = Text, Slash, Text ;
// obscured_coords = Hash, Hash, Number ;
func (p *parser) parseCoords() CoordsNode {
	switch p.peekKind() {
	case lexers.Grid:
		// grid_coords = Grid, Number, but lexer returns Grid for Hash Hash
		if bytes.Equal(p.peek().Bytes(), []byte{'#', '#'}) {
			// obscured_coords = Hash, Hash, Number (e.g., ## 1315)
			node := &ObscuredCoordsNode{}
			node.Grid = p.advance()
			node.tokens = append(node.tokens, node.Grid)

			if tok, err := p.expect(lexers.Number); err != nil {
				node.errors = append(node.errors, err)
			} else {
				node.Number = tok
				node.tokens = append(node.tokens, tok)
			}
			return node
		}

		node := &GridCoordsNode{}
		node.Grid = p.advance()
		node.tokens = append(node.tokens, node.Grid)

		if tok, err := p.expect(lexers.Number); err != nil {
			node.errors = append(node.errors, err)
		} else {
			node.Number = tok
			node.tokens = append(node.tokens, tok)
		}
		return node

	case lexers.NA:
		// na_coords = Text, Slash, Text (e.g., N/A)
		node := &NACoordsNode{}
		node.Text = p.advance()
		node.tokens = append(node.tokens, node.Text)
		return node

	default:
		// Error case - unexpected token
		node := &ErrorCoordsNode{
			Message: fmt.Sprintf("expected coordinates, got %s", p.peekKind()),
		}
		node.errors = append(node.errors, fmt.Errorf("%s", node.Message))
		return node
	}
}
