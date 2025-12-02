// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package lexers

import (
	"bytes"
	"errors"
	"regexp"
)

// Lexer_t implements a lexer that returns tokens for the CST.
type Lexer_t struct {
	input []byte
	pos   int // current position in input
	line  int // current line (1-based)
	col   int // current column (1-based, eol is last column on line)
}

// New returns an initialized lexer.
func New(input []byte) (*Lexer_t, error) {
	if input == nil {
		return nil, errors.New("input cannot be nil")
	}
	return &Lexer_t{
		input: input,
		pos:   0,
		line:  1,
		col:   1,
	}, nil
}

var (
	reGrid   = regexp.MustCompile(`^[A-Z]{2}$`)
	reNumber = regexp.MustCompile(`^\d+$`)
	keywords = map[string]TokenKind_t{
		"Courier":  TokenCourier,
		"Current":  TokenCurrent,
		"Element":  TokenElement,
		"FINE":     TokenWeather,
		"Fall":     TokenSeason,
		"Fleet":    TokenFleet,
		"Garrison": TokenGarrison,
		"Hex":      TokenHex,
		"Next":     TokenNext,
		"Previous": TokenPrevious,
		"Scout":    TokenScout,
		"Spring":   TokenSeason,
		"Status":   TokenStatus,
		"Summer":   TokenSeason,
		"Tribe":    TokenTribe,
		"Turn":     TokenTurn,
		"Winter":   TokenSeason,
	}
)

// Next returns the next token from the input or an EOF token on end of input.
func (l *Lexer_t) Next() *Token_t {
	if l.iseof() {
		return &Token_t{
			Kind: TokenEOF,
			Span: Span_t{
				Start: l.pos,
				End:   l.pos,
				Line:  l.line,
				Col:   l.col,
			},
		}
	}

	token := Token_t{}

	// collect leading trivia (spaces and tabs, not newlines)
	if isTrivia(l.curr()) {
		token.LeadingTrivia = []Trivia_t{{
			Kind: Whitespace,
			Span: Span_t{
				Start: l.pos,
				Line:  l.line,
				Col:   l.col,
			},
		}}
		for isTrivia(l.curr()) {
			l.next()
		}
		token.LeadingTrivia[0].Span.End = l.pos
	}

	token.Span.Start = l.pos
	token.Span.Line = l.line
	token.Span.Col = l.col

	// check for EOF
	if l.iseof() {
		token.Kind = TokenEOF
		token.Span.End = l.pos
		return &token
	}

	// check for delimiters
	if ch := l.curr(); isDelimiter(ch) {
		switch ch {
		case ' ':
			// this should not happen - the leading trivia code should consume these
			token.Kind = TokenError
		case '\\':
			token.Kind = TokenBackslash
		case ':':
			token.Kind = TokenColon
		case ',':
			token.Kind = TokenComma
		case '-':
			token.Kind = TokenDash
		case '$':
			token.Kind = TokenDollar
		case '.':
			token.Kind = TokenDot
		case '"':
			token.Kind = TokenDouble
		case '\n':
			token.Kind = TokenEOL
		case '=':
			token.Kind = TokenEquals
		case '#':
			token.Kind = TokenHash
		case '(':
			token.Kind = TokenParenL
		case ')':
			token.Kind = TokenParenR
		case '+':
			token.Kind = TokenPlus
		case '\'':
			token.Kind = TokenSingle
		case '*':
			token.Kind = TokenStar
		case '/':
			token.Kind = TokenSlash
		default:
			// this catches undefined delimiters
			token.Kind = TokenDelimiter
		}
		l.next() // consume the delimiter
		token.Span.End = l.pos
		return &token
	}

	// collect text up to the first delimiter
	token.Kind = TokenText
	for !isDelimiter(l.curr()) {
		l.next()
	}
	token.Span.End = l.pos

	// collect trailing trivia (spaces only, not tabs or newlines)
	if isTrivia(l.curr()) {
		token.TrailingTrivia = []Trivia_t{{
			Kind: Whitespace,
			Span: Span_t{
				Start: l.pos,
				Line:  l.line,
				Col:   l.col,
			},
		}}
		for isTrivia(l.curr()) {
			l.next()
		}
		token.TrailingTrivia[0].Span.End = l.pos
	}

	tokenText := token.Span.Text(l.input)

	if kind, ok := keywords[string(tokenText)]; ok {
		token.Kind = kind
	} else if len(tokenText) == 2 && reGrid.Match(tokenText) {
		token.Kind = TokenGrid
	} else if reNumber.Match(tokenText) {
		token.Kind = TokenNumber
	}

	return &token
}

func (l *Lexer_t) curr() byte {
	if l.iseof() {
		return 0
	}
	return l.input[l.pos]
}

func (l *Lexer_t) iseof() bool {
	return l.pos >= len(l.input)
}

func (l *Lexer_t) next() {
	if l.iseof() {
		return
	}
	if l.input[l.pos] == '\n' {
		l.line, l.col = l.line+1, 1
	} else {
		l.col = l.col + 1
	}
	l.pos++
}

// isDelimiter returns true if the character should terminate a word
func isDelimiter(ch byte) bool {
	return bytes.IndexByte([]byte{' ', '\t', '\n', '\'', '"', '.', ',', '(', ')', '#', '+', '-', '*', '/', '=', '\\', '$', ':', 0}, ch) != -1
}

// isTrivia returns true if the character is whitespace that counts as trivia (spaces only)
func isTrivia(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\r'
}
