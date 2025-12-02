// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package lexers implements a lexer for turn reports.
// Returns tokens that contain copies from the input buffer.
package lexers

import (
	"regexp"
)

// Scan returns all the tokens in the input buffer.
func Scan(input []byte) []*Token {
	var tokens []*Token

	// todo: we implemented two tokens of look back for N/A. we updated
	// the text branch to look ahead instead. we might be able to reduce
	// this to a single token of look back
	prevKind := [2]Kind{EOF, EOF}

	pos, line, col := 0, 1, 1
	for pos < len(input) {
		var token Token

		// check for leading trivia
		if pos < len(input) && isTrivia[input[pos]] {
			start, span := pos, &Span{Line: line, Col: col, Kind: Spaces}
			for pos < len(input) && isTrivia[input[pos]] {
				pos, line, col = pos+1, line, col+1
			}
			span.Value = bdup(input[start:pos])
			token.LeadingTrivia = []*Span{span}
		}

		// check for text
		if pos < len(input) {
			start, span := pos, &Span{Line: line, Col: col}
			if input[pos] == '\n' {
				span.Kind = EOL
				pos, line, col = pos+1, line+1, 1
			} else if isDelimiter[input[pos]] {
				length := 1
				switch input[pos] {
				case '\\':
					span.Kind = Backslash
				case ':':
					span.Kind = Colon
				case ',':
					span.Kind = Comma
				case '-':
					span.Kind = Dash
				case '=':
					span.Kind = Equals
				case '#':
					span.Kind = Hash
					if pos+1 < len(input) && input[pos+1] == '#' {
						span.Kind, length = Grid, length+1
					}
				case '(':
					span.Kind = LeftParen
				case ')':
					span.Kind = RightParen
				case '/':
					span.Kind = Slash
				default:
					span.Kind = Delimiter
				}
				pos, line, col = pos+length, line, col+length
			} else {
				if text := reNA.Find(input[pos:]); text != nil {
					// treat N/A as a single token
					span.Kind = NA
					pos, line, col = pos+len(text), line, col+len(text)
				} else if text = reTurnYearMonth.Find(input[pos:]); text != nil {
					// treat yyyy-mm as a single token
					span.Kind = TurnYearMonth
					pos, line, col = pos+len(text), line, col+len(text)
				} else {
					span.Kind = Text
					for pos < len(input) && !(isTrivia[input[pos]] || isDelimiter[input[pos]]) {
						pos, line, col = pos+1, line, col+1
					}
				}
			}
			span.Value = bdup(input[start:pos])
			if span.Kind == Text {
				text := span.Bytes()
				if kw, ok := keywords[string(text)]; ok {
					span.Kind = kw
				} else if len(text) == 2 && reGrid.Match(text) {
					span.Kind = Grid
				} else if reNumber.Match(text) {
					span.Kind = Number
					if prevKind[0] == Tribe && len(text) == 4 {
						span.Kind = UnitId
					}
				} else if reUnitId.Match(text) {
					// warning - must be after the check for numbers since
					// the unit id for tribes is a valid four-digit number.
					span.Kind = UnitId
				}
			}
			token.Kind, token.Value = span.Kind, []*Span{span}
		}

		// check for trailing trivia
		if pos < len(input) && isTrivia[input[pos]] {
			start, span := pos, &Span{Line: line, Col: col, Kind: Spaces}
			for pos < len(input) && isTrivia[input[pos]] {
				pos, line, col = pos+1, line, col+1
			}
			span.Value = bdup(input[start:pos])
			token.TrailingTrivia = []*Span{span}
		}

		// if there's no value, then the kind is still text
		if len(token.Value) == 0 {
			token.Kind = Text
		}

		tokens = append(tokens, &token)
		prevKind[0], prevKind[1] = token.Kind, prevKind[0]
	}

	return tokens
}

func bdup(src []byte) []byte {
	dst := make([]byte, len(src))
	copy(dst, src)
	return dst
}

func init() {
	for _, ch := range []byte{' ', '\t', '\r'} {
		isTrivia[ch] = true
		isDelimiter[ch] = true
	}
	for _, ch := range []byte{0, '\n', '\'', '"', '.', ',', '(', ')', '#', '+', '-', '*', '/', '=', '\\', '$', ':'} {
		isDelimiter[ch] = true
	}
}

var (
	isDelimiter = [256]bool{}
	isTrivia    = [256]bool{}

	keywords = map[string]Kind{
		"Tribe":    Tribe,
		"Courier":  Courier,
		"Element":  Element,
		"Garrison": Garrison,
		"Fleet":    Fleet,
		"Current":  Current,
		"Hex":      Hex,
		"Previous": Previous,

		"Turn":   Turn,
		"Spring": Season,
		"Summer": Season,
		"Fall":   Season,
		"Winter": Season,
		"FINE":   Weather,
		"Next":   Next,
	}

	reGrid          = regexp.MustCompile(`^[A-Z]{2}$`)
	reNA            = regexp.MustCompile(`^N/A`)
	reNumber        = regexp.MustCompile(`^\d+$`)
	reTurnYearMonth = regexp.MustCompile(`^\d{3,4}-\d{1,2}`)
	reUnitId        = regexp.MustCompile(`^\d{4}([cefg][1-9])?$`)
)
