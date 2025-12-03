// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package lexers

//go:generate stringer --type Kind

import "bytes"

type Token struct {
	Kind           Kind
	LeadingTrivia  []*Span
	Value          []*Span
	TrailingTrivia []*Span
}

// Position returns the line and column number of the token.
func (t *Token) Position() (int, int) {
	if t == nil {
		return 0, 0
	} else if len(t.Value) != 0 {
		return t.Value[0].Line, t.Value[0].Col
	} else if len(t.LeadingTrivia) != 0 {
		return t.LeadingTrivia[0].Line, t.LeadingTrivia[0].Col
	} else if len(t.TrailingTrivia) != 0 {
		return t.TrailingTrivia[0].Line, t.TrailingTrivia[0].Col
	}
	return 0, 0
}

// Length returns the length of the token.
func (t *Token) Length() int {
	if t == nil {
		return 0
	}
	length := 0
	for _, span := range t.Value {
		length += span.Length()
	}
	return length
}

// Bytes is a helper for diagnostics / debugging.
// Returns an empty slice if there is no value.
func (t *Token) Bytes() []byte {
	b := &bytes.Buffer{}
	if t != nil {
		for _, span := range t.Value {
			b.Write(span.Bytes())
		}
	}
	return b.Bytes()
}

// Source is a helper function to rebuild the input from the token stream.
func (t *Token) Source() []byte {
	b := &bytes.Buffer{}
	for _, span := range t.LeadingTrivia {
		b.Write(span.Bytes())
	}
	for _, span := range t.Value {
		b.Write(span.Bytes())
	}
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

// LineNo returns the line number of the start of the span
func (s *Span) LineNo() int {
	if s == nil {
		return 0
	}
	return s.Line
}

// ColNo returns the column number of the start of the span
func (s *Span) ColNo() int {
	if s == nil {
		return 0
	}
	return s.Col
}

// Bytes is a helper for diagnostics / debugging.
func (s *Span) Bytes() []byte {
	if s == nil {
		return nil
	}
	return s.Value
}

// Merge is a helper function to merge tokens tokens
func (t *Token) Merge(tokens ...*Token) {
	for _, tok := range tokens {
		if tok.LeadingTrivia != nil {
			t.Value = append(t.Value, tok.LeadingTrivia...)
		}
		if tok.Value != nil {
			t.Value = append(t.Value, tok.Value...)
		}
		if tok.TrailingTrivia != nil {
			t.Value = append(t.Value, tok.TrailingTrivia...)
		}
	}
}

// Merge is a helper function to merge tokens tokens
func Merge(kind Kind, tokens ...*Token) *Token {
	t := &Token{Kind: kind}
	for _, tok := range tokens {
		t.Merge(tok)
	}
	return t
}

// ToSource is a helper function to rebuild the source from the token stream.
func ToSource(tokens ...*Token) []byte {
	b := &bytes.Buffer{}
	for _, t := range tokens {
		for _, span := range t.LeadingTrivia {
			b.Write(span.Bytes())
		}
		for _, span := range t.Value {
			b.Write(span.Bytes())
		}
		for _, span := range t.TrailingTrivia {
			b.Write(span.Bytes())
		}
	}
	return b.Bytes()
}

type Kind int

const (
	Text Kind = iota

	Backslash
	Colon
	Comma
	Dash
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
