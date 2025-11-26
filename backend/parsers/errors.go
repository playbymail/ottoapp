// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package parsers

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrBadInput         = Error("bad input")
	ErrNotAWordDocument = Error("not a word document")
	ErrNotATurnReport   = Error("not a turn report")
)
