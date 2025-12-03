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

// runUpTo returns the tokens up to (but not including) a delimiter token.
// Returns tokens, false if there were tokens but no delimiter.
// Returns nil, true if there were no tokens found before a delimiter.
// Returns nil, false when at EOF
func (p *parser) runUpTo(delimiters ...lexers.Kind) ([]*lexers.Token, bool) {
	var tokens []*lexers.Token
	for tok := p.peek(); tok != nil; tok = p.peek() {
		for _, delimiter := range delimiters {
			if tok.Kind == delimiter {
				return tokens, true
			}
		}
		// consume the token
		tokens = append(tokens, p.advance())
	}
	return tokens, false
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
	node.LineNo, node.ColNo = p.peek().Position()

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
// unit_section = unit_line, turn_line, [ unit_movement_line ] ;
func (p *parser) parseUnitSection() *UnitSectionNode {
	node := &UnitSectionNode{}
	node.LineNo, node.ColNo = p.peek().Position()

	unitLine := p.parseUnitLine()
	node.UnitLine = unitLine
	node.tokens = append(node.tokens, unitLine.tokens...)
	node.errors = append(node.errors, unitLine.errors...)

	// Parse turn_line if the next token is Current (start of turn line)
	if p.peekKind() == lexers.Current {
		turnLine := p.parseTurnLine()
		node.TurnLine = turnLine
		node.tokens = append(node.tokens, turnLine.tokens...)
		node.errors = append(node.errors, turnLine.errors...)
	}

	// Parse optional unit_movement_line
	// unit_movement_line = unit_goes_to_line | land_movement_line ;
	// unit_goes_to_line starts with Tribe, Goes
	// land_movement_line starts with Tribe, Movement
	if p.peekKind() == lexers.Tribe && p.pos+1 < len(p.tokens) {
		nextKind := p.tokens[p.pos+1].Kind
		switch nextKind {
		case lexers.Goes:
			// unit_goes_to_line: Tribe Goes To ...
			movementLine := p.parseUnitGoesToLine()
			node.UnitMovementLine = movementLine
			node.tokens = append(node.tokens, movementLine.tokens...)
			node.errors = append(node.errors, movementLine.errors...)
		case lexers.Movement:
			// land_movement_line: Tribe Movement: Move ...
			movementLine := p.parseLandMovementLine()
			node.UnitMovementLine = movementLine
			node.tokens = append(node.tokens, movementLine.tokens...)
			node.errors = append(node.errors, movementLine.errors...)
		}
	}

	return node
}

// parseUnitLine parses a unit line.
// unit_line = unit_keyword, unit_id, Comma, [ Note ], Comma,
//
//	Current, Hex, Equals, coords, Comma,
//	LeftParen, Previous, Hex, Equals, coords, RightParen, EOL ;
func (p *parser) parseUnitLine() *UnitLineNode {
	node := &UnitLineNode{}
	node.LineNo, node.ColNo = p.peek().Position()

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
	if noteTokens, _ := p.runUpTo(lexers.Comma, lexers.EOL); noteTokens != nil {
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
	node.LineNo, node.ColNo = p.peek().Position()

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
	node.LineNo, node.ColNo = p.peek().Position()

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
	node.LineNo, node.ColNo = p.peek().Position()

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

// parseUnitGoesToLine parses a "Tribe Goes to" movement line.
// unit_goes_to_line = Tribe, Goes, To, grid_coords, EOL ;
func (p *parser) parseUnitGoesToLine() *UnitGoesToLineNode {
	node := &UnitGoesToLineNode{}
	node.LineNo, node.ColNo = p.peek().Position()

	// Tribe
	if tok, err := p.expect(lexers.Tribe); err != nil {
		node.errors = append(node.errors, err)
		skipped := p.syncToNextLine()
		node.tokens = append(node.tokens, skipped...)
		return node
	} else {
		node.Tribe = tok
		node.tokens = append(node.tokens, tok)
	}

	// Goes
	if tok, err := p.expect(lexers.Goes); err != nil {
		node.errors = append(node.errors, err)
		skipped := p.syncToNextLine()
		node.tokens = append(node.tokens, skipped...)
		return node
	} else {
		node.Goes = tok
		node.tokens = append(node.tokens, tok)
	}

	// To
	if tok, err := p.expect(lexers.To); err != nil {
		node.errors = append(node.errors, err)
		skipped := p.syncToNextLine()
		node.tokens = append(node.tokens, skipped...)
		return node
	} else {
		node.To = tok
		node.tokens = append(node.tokens, tok)
	}

	// grid_coords
	coords := p.parseCoords()
	node.tokens = append(node.tokens, coords.Tokens()...)
	node.errors = append(node.errors, coords.Errors()...)
	if gc, ok := coords.(*GridCoordsNode); ok {
		node.Coords = gc
	} else {
		node.errors = append(node.errors, fmt.Errorf("expected grid coordinates"))
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

// parseLandMovementLine parses a "Tribe Movement:" line.
// land_movement_line = Tribe, Movement, Colon, land_movement, EOL ;
func (p *parser) parseLandMovementLine() *LandMovementLineNode {
	node := &LandMovementLineNode{}
	node.LineNo, node.ColNo = p.peek().Position()

	// Tribe
	if tok, err := p.expect(lexers.Tribe); err != nil {
		node.errors = append(node.errors, err)
		skipped := p.syncToNextLine()
		node.tokens = append(node.tokens, skipped...)
		return node
	} else {
		node.Tribe = tok
		node.tokens = append(node.tokens, tok)
	}

	// Movement
	if tok, err := p.expect(lexers.Movement); err != nil {
		node.errors = append(node.errors, err)
		skipped := p.syncToNextLine()
		node.tokens = append(node.tokens, skipped...)
		return node
	} else {
		node.Movement = tok
		node.tokens = append(node.tokens, tok)
	}

	// Colon
	if tok, err := p.expect(lexers.Colon); err != nil {
		node.errors = append(node.errors, err)
		skipped := p.syncToNextLine()
		node.tokens = append(node.tokens, skipped...)
		return node
	} else {
		node.Colon = tok
		node.tokens = append(node.tokens, tok)
	}

	// land_movement
	landMovement := p.parseLandMovement()
	node.LandMovement = landMovement
	node.tokens = append(node.tokens, landMovement.tokens...)
	node.errors = append(node.errors, landMovement.errors...)

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

// parseLandMovement parses the movement portion.
// land_movement = Move, land_step, { Backslash, land_step } ;
func (p *parser) parseLandMovement() *LandMovementNode {
	node := &LandMovementNode{}
	node.LineNo, node.ColNo = p.peek().Position()

	// Move
	if tok, err := p.expect(lexers.Move); err != nil {
		node.errors = append(node.errors, err)
		return node
	} else {
		node.Move = tok
		node.tokens = append(node.tokens, tok)
	}

	// First land_step
	step := p.parseLandStep()
	node.Steps = append(node.Steps, step)
	node.tokens = append(node.tokens, step.tokens...)
	node.errors = append(node.errors, step.errors...)

	// { Backslash, land_step }
	for p.peekKind() == lexers.Backslash {
		// consume the backslash
		backslash := p.advance()
		node.tokens = append(node.tokens, backslash)

		// parse the next land_step
		step := p.parseLandStep()
		node.Steps = append(node.Steps, step)
		node.tokens = append(node.tokens, step.tokens...)
		node.errors = append(node.errors, step.errors...)
	}

	return node
}

// parseLandStep parses a single step.
// land_step = [ [ land_step_movement ], land_step_result ];
// land_step_movement = Direction, Dash, Terrain ;
// land_step_result = Comma ;
func (p *parser) parseLandStep() *LandStepNode {
	node := &LandStepNode{}
	node.LineNo, node.ColNo = p.peek().Position()

	// Check if we're at the end of the step sequence (EOL or next unit keyword)
	// If so, return an empty step
	if p.peekKind() == lexers.EOL || p.peekKind() == lexers.Backslash || p.isAtEnd() {
		return node
	}

	// Try to parse land_step_movement: Direction, Dash, Terrain
	if p.peekKind() == lexers.Direction {
		node.Direction = p.advance()
		node.tokens = append(node.tokens, node.Direction)

		// Dash
		if tok, err := p.expect(lexers.Dash); err != nil {
			node.errors = append(node.errors, err)
			return node
		} else {
			node.Dash = tok
			node.tokens = append(node.tokens, tok)
		}

		// Terrain
		if tok, err := p.expect(lexers.TerrainCode); err != nil {
			node.errors = append(node.errors, err)
			return node
		} else {
			node.Terrain = tok
			node.tokens = append(node.tokens, tok)
		}
	}

	// land_step_result = Comma (required if step content is present)
	if p.peekKind() == lexers.Comma {
		node.Comma = p.advance()
		node.tokens = append(node.tokens, node.Comma)
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
			node.LineNo, node.ColNo = p.peek().Position()

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
		node.LineNo, node.ColNo = p.peek().Position()

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
		node.LineNo, node.ColNo = p.peek().Position()

		node.Text = p.advance()
		node.tokens = append(node.tokens, node.Text)
		return node

	default:
		// Error case - unexpected token
		node := &ErrorCoordsNode{
			Message: fmt.Sprintf("expected coordinates, got %s", p.peekKind()),
		}
		node.LineNo, node.ColNo = p.peek().Position()

		node.errors = append(node.errors, fmt.Errorf("%s", node.Message))
		return node
	}
}
