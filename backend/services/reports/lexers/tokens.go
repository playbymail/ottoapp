// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package lexers

//go:generate stringer --type Kind

import "bytes"

type Token struct {
	LeadingTrivia  []*Span
	Value          *Span
	TrailingTrivia []*Span
}

// Bytes is a helper for diagnostics / debugging.
func (t *Token) Bytes(src []byte) []byte {
	return t.Value.Bytes()
}

// Source is a helper function to rebuild the input from the token stream.
func (t *Token) Source() []byte {
	b := &bytes.Buffer{}
	for _, span := range t.LeadingTrivia {
		b.Write(span.Bytes())
	}
	b.Write(t.Value.Bytes())
	for _, span := range t.TrailingTrivia {
		b.Write(span.Bytes())
	}
	return b.Bytes()
}

type Span struct {
	Line  int // 1-based
	Col   int // 1-based, in UTF-8 code points
	Kind  Kind
	Value []byte
}

// Length returns the length of the span.
func (s *Span) Length() int {
	if s == nil {
		return 0
	}
	return len(s.Value)
}

// Bytes is a helper for diagnostics / debugging.
func (s *Span) Bytes() []byte {
	if s == nil {
		return nil
	}
	return s.Value
}

// ToSource is helper function to rebuild the source from the token stream.
func ToSource(tokens ...*Token) []byte {
	b := &bytes.Buffer{}
	for _, t := range tokens {
		for _, span := range t.LeadingTrivia {
			b.Write(span.Bytes())
		}
		b.Write(t.Value.Bytes())
		for _, span := range t.TrailingTrivia {
			b.Write(span.Bytes())
		}
	}
	return b.Bytes()
}

type Kind int

const (
	Text Kind = iota

	Comma
	Dash
	Delimiter
	Equals
	EOL
	EOF
	Grid
	Hash
	Note
	Number
	LeftParen
	RightParen
	Slash
	Spaces

	// keywords
	Tribe
	Courier
	Element
	Fleet
	Garrison
	Current
	Hex
	Previous
	Turn

	Season
	Next
	Weather
)
