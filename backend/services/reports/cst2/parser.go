// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package cst2

import (
	"fmt"

	"github.com/playbymail/ottoapp/backend/services/reports/lexers"
)

type Option func(p *parser) error

func Parse(tokens []*lexers.Token, options ...Option) (Node, error) {
	p := &parser{
		tokens: tokens,
	}

	for _, option := range options {
		if err := option(p); err != nil {
			return nil, err
		}
	}

	turnReport := p.parseTurnReport()

	return turnReport, nil
}

type parser struct {
	tokens []*lexers.Token
	pos    int // offset into tokens buffer
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

// advanceTo returns the tokens up to (but not including) a delimiter token.
// Returns tokens, false if there were tokens but no delimiter.
// Returns nil, true if there were no tokens found before a delimiter.
// Returns nil, false when at EOF
func (p *parser) advanceTo(delimiters ...lexers.Kind) ([]*lexers.Token, bool) {
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

// syncToDelimiter advances the next delimiter token.
// Used for error recovery.
func (p *parser) syncToDelimiter(kind lexers.Kind) []*lexers.Token {
	tokens, _ := p.advanceTo(kind)
	return tokens
}

// wantPattern returns true if the next set of tokens matches the want sequence
func (p *parser) wantPattern(wants ...lexers.Kind) bool {
	for offset, want := range wants {
		if p.pos+offset >= len(p.tokens) {
			return false
		} else if p.tokens[p.pos+offset].Kind != want {
			return false
		}
	}
	return true
}

// parseTurnReport parses the entire token stream into a TurnReportNode.
// turn_report = [ prologue_section ], unit_section, { unit_section }, EOF ;
func (p *parser) parseTurnReport() *TurnReportNode {
	node := newTurnReportNode()

	// parse optional prologue section
	node.prologue = p.parsePrologueSection()

	// parse zero or more unit sections
	for section := p.parseUnitSection(); section != nil; section = p.parseUnitSection() {
		node.sections = append(node.sections, section)
	}

	// parse optional epilogue section
	node.epilogue = p.parseEpilogueSection()

	return node
}

// parsePrologueSection collects lines up to the first unit_section.
// prologue_section = { line } ;
func (p *parser) parsePrologueSection() *PrologueSectionNode {
	var node *PrologueSectionNode
	for !p.isAtEnd() {
		if p.wantPattern(lexers.Tribe, lexers.UnitId) {
			// we found the start of a unit section
			break
		}
		if node == nil {
			node = newPreambleSectionNode()
		}
		line := p.parseLine()
		if line == nil {
			break
		}
		node.lines = append(node.lines, line)
	}
	return node
}

// parseLine parses a single line - all the tokens up to (and including) the EOF.
func (p *parser) parseLine() *LineNode {
	if p.isAtEnd() {
		return nil
	}
	line := newLineNode()
	tokens, foundEol := p.advanceTo(lexers.EOL)
	if len(tokens) != 0 {
		line.line = append(line.line, tokens...)
		line.tokens = append(line.tokens, tokens...)
	}
	if foundEol {
		eol := p.advance()
		line.line = append(line.line, eol)
		line.tokens = append(line.tokens, eol)
	}
	return line
}

// parseUnitSection parses the current unit section.
// We expect this to be called at the start of a unit, so we will
// return either a UnitSectionNode or BadUnitSectionNode.
//
// unit_section = unit_line, turn_line, [ unit_movement_line ], { scout_line } status_line ;
func (p *parser) parseUnitSection() Node {
	if !p.wantPattern(lexers.Tribe, lexers.UnitId) {
		// bad unit section?
		return nil
	}

	node := newUnitSectionNode()

	unitLine := p.parseUnitLine()
	_ = unitLine

	heading := p.parseLine()
	node.nodes = append(node.nodes, heading)
	foundStatusLine := false
	for !(foundStatusLine || p.isAtEnd()) {
		// sync point - stop if we find another unit section
		if p.wantPattern(lexers.Tribe, lexers.UnitId) {
			break
		}
		foundStatusLine = p.wantPattern(lexers.UnitId, lexers.Status, lexers.Colon)
		node.nodes = append(node.nodes, p.parseLine())
	}
	if foundStatusLine {
		// consume other lines up to the next unit section
		for !(p.isAtEnd() || p.wantPattern(lexers.Tribe, lexers.UnitId)) {
			node.epilogue = append(node.epilogue, p.parseLine())
		}
	}

	return node
}

// parseUnitLine parses a unit line.
// unit_line = unit_keyword, unit_id, Comma, [ Note ], Comma, Current, Hex, Equals, coords, Comma, LeftParen, Previous, Hex, Equals, coords, RightParen, EOL ;
func (p *parser) parseUnitLine() *UnitLineNode {
	node := newUnitLineNode()
	node.UnitKeyword = p.match(lexers.Courier, lexers.Element, lexers.Fleet, lexers.Garrison, lexers.Tribe)
	node.UnitId = p.match(lexers.UnitId)
	node.Comma1 = p.match(lexers.Comma)

	return node
}

// parseEpilogueSection collects lines up to end of input.
// epilogue_section = { line } EOF ;
func (p *parser) parseEpilogueSection() *EpilogueSectionNode {
	var node *EpilogueSectionNode
	for line := p.parseLine(); line != nil; line = p.parseLine() {
		if node == nil {
			node = newEpilogueSectionNode()
		}
		node.lines = append(node.lines, line)
	}
	return node
}
