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
				switch input[pos] {
				case ',':
					span.Kind = Comma
				case '-':
					span.Kind = Dash
				case '=':
					span.Kind = Equals
				case '#':
					span.Kind = Hash
				case '(':
					span.Kind = LeftParen
				case ')':
					span.Kind = RightParen
				case '/':
					span.Kind = Slash
				default:
					span.Kind = Delimiter
				}
				pos, line, col = pos+1, line, col+1
			} else {
				span.Kind = Text
				for pos < len(input) && !(isTrivia[input[pos]] || isDelimiter[input[pos]]) {
					pos, line, col = pos+1, line, col+1
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
				}
			}
			token.Value = span
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

		tokens = append(tokens, &token)
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

	reGrid   = regexp.MustCompile(`^[A-Z]{2}$`)
	reNumber = regexp.MustCompile(`^\d+$`)
)
