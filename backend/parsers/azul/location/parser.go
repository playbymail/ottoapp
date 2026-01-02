// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package location

//go:generate pigeon -o grammar.go grammar.peg

import (
	"errors"
	"log"

	"github.com/playbymail/ottoapp/backend/domains"
)

type Location struct {
	UnitId         string
	Note           string
	CurrentCoords  string
	PreviousCoords string
}

func Expect(filename string, b []byte, opts ...Option) (Location, error) {
	p := newParser(filename, b, opts...)
	val, err := p.parse(g)
	if err != nil {
		return Location{}, err
	}
	l, ok := val.(Location)
	if !ok {
		log.Printf("%s: error parsing turn\n", filename)
		log.Printf("error: expected %T, got %T\n", l, val)
		log.Printf("please report this error\n")
		return Location{}, domains.ErrBadInput
	}
	return l, nil
}

// ParseError contains position information extracted from parser errors.
type ParseError struct {
	Error  error
	Prefix string
	Inner  string
	Pos    struct {
		Line   int
		Col    int
		Offset int
	}
	Expected []string
}

// ExtractParseError attempts to extract position information from a pigeon parser error.
// Returns nil if the error is not a parser error or position cannot be extracted.
func ExtractParseError(err error) *ParseError {
	if err == nil {
		return nil
	}

	// Try to unwrap to errList first
	var el errList
	if errors.As(err, &el) && len(el) > 0 {
		// Get the first error which should be a parserError
		var pe *parserError
		if errors.As(el[0], &pe) {
			return &ParseError{
				Error:  err,
				Prefix: pe.prefix,
				Inner:  pe.Inner.Error(),
				Pos: struct {
					Line   int
					Col    int
					Offset int
				}{
					Line:   pe.pos.line,
					Col:    pe.pos.col,
					Offset: pe.pos.offset,
				},
				Expected: pe.expected,
			}
		}
	}

	// Try direct parserError
	var pe *parserError
	if errors.As(err, &pe) {
		return &ParseError{
			Error:  err,
			Prefix: pe.prefix,
			Inner:  pe.Inner.Error(),
			Pos: struct {
				Line   int
				Col    int
				Offset int
			}{
				Line:   pe.pos.line,
				Col:    pe.pos.col,
				Offset: pe.pos.offset,
			},
			Expected: pe.expected,
		}
	}

	return nil
}
