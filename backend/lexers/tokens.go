// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package lexers

import (
	"bytes"
	"fmt"
)

// Token_t represents a lexical token with type, value, and trivia
type Token_t struct {
	Kind           TokenKind_t
	Span           Span_t
	LeadingTrivia  []Trivia_t
	TrailingTrivia []Trivia_t
}

// Text is a helper for diagnostics / debugging.
func (t Token_t) Text(src []byte) []byte {
	return t.Span.Text(src)
}

// LeadingTriviaText is a helper for diagnostics / debugging.
func (t Token_t) LeadingTriviaText(src []byte) []byte {
	if len(t.LeadingTrivia) == 0 {
		return nil
	}
	b := &bytes.Buffer{}
	for _, trivia := range t.LeadingTrivia {
		b.Write(trivia.Text(src))
	}
	return b.Bytes()
}

// TrailingTriviaText is a helper for diagnostics / debugging.
func (t Token_t) TrailingTriviaText(src []byte) []byte {
	if len(t.TrailingTrivia) == 0 {
		return nil
	}
	b := &bytes.Buffer{}
	for _, trivia := range t.TrailingTrivia {
		b.Write(trivia.Text(src))
	}
	return b.Bytes()
}

// TokenKind_t represents the kind of token
type TokenKind_t int // ClanId, TribeId, etc

const (
	TokenError TokenKind_t = iota
	TokenBackslash
	TokenColon
	TokenComma
	TokenDash
	TokenDollar
	TokenDot
	TokenEOF
	TokenEOL
	TokenEquals
	TokenGrid
	TokenHash
	TokenNumber
	TokenParenL
	TokenParenR
	TokenPlus
	TokenSingle
	TokenDouble
	TokenSlash
	TokenStar
	TokenText
	TokenDelimiter
	// Keywords
	TokenCourier
	TokenCurrent
	TokenElement
	TokenFleet
	TokenGarrison
	TokenHex
	TokenNext
	TokenNote
	TokenPrevious
	TokenScout
	TokenSeason
	TokenStatus
	TokenTribe
	TokenTurn
	TokenWeather
)

func (t TokenKind_t) String() string {
	switch t {
	case TokenError:
		return "error"
	case TokenBackslash:
		return "\\"
	case TokenColon:
		return ":"
	case TokenComma:
		return ","
	case TokenDash:
		return "-"
	case TokenDollar:
		return "$"
	case TokenDot:
		return "."
	case TokenEOF:
		return "eof"
	case TokenEOL:
		return "eol"
	case TokenEquals:
		return "="
	case TokenHash:
		return "#"
	case TokenNumber:
		return "number"
	case TokenParenL:
		return "("
	case TokenParenR:
		return ")"
	case TokenPlus:
		return "+"
	case TokenSingle:
		return "'"
	case TokenDouble:
		return "\""
	case TokenSlash:
		return "/"
	case TokenStar:
		return "*"
	case TokenText:
		return "text"
	case TokenDelimiter:
		return "delimiter"
	case TokenGrid:
		return "grid"
	case TokenTribe:
		return "tribe"
	case TokenCurrent:
		return "current"
	case TokenPrevious:
		return "previous"
	case TokenTurn:
		return "turn"
	case TokenNext:
		return "next"
	case TokenScout:
		return "scout"
	case TokenSeason:
		return "season"
	case TokenStatus:
		return "status"
	case TokenHex:
		return "hex"
	case TokenNote:
		return "note"
	case TokenWeather:
		return "weather"
	}
	panic(fmt.Sprintf("assert(kind != %d)", t))
}

type Span_t struct {
	Start int // byte offset (inclusive)
	End   int // byte offset (exclusive)
	Line  int // 1-based
	Col   int // 1-based, in UTF-8 code points
}

// Length returns the length of the span.
func (s Span_t) Length() int {
	return s.End - s.Start
}

// Text is a helper for diagnostics / debugging.
func (s Span_t) Text(src []byte) []byte {
	return src[s.Start:s.End]
}

type Trivia_t struct {
	Kind TriviaKind_t
	Span Span_t
}

// Text is a helper for diagnostics / debugging.
func (t Trivia_t) Text(src []byte) []byte {
	return t.Span.Text(src)
}

// TriviaKind_t represents the kind of trivia.
type TriviaKind_t int

const (
	Whitespace TriviaKind_t = iota
)

func (t TriviaKind_t) String() string {
	switch t {
	case Whitespace:
		return "whitespace"
	}
	panic(fmt.Sprintf("assert(kind != %d)", t))
}
