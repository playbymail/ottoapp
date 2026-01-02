// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package status

import "bytes"

type Lexer struct {
	input []byte
}

type Token struct {
	Lexeme string
	Trivia string
}

type position struct{}

// Next returns the next token.
func Next(b []byte) (Token, []byte) {
	if len(b) == 0 {
		return Token{}, nil
	}
	lexeme, rest := Split(b)
	t := Token{Lexeme: string(lexeme)}
	if len(rest) > 0 && isspace(rest[0]) {
		var trivia []byte
		trivia, rest = Split(rest)
		t.Trivia = string(trivia)
	}
	return t, rest
}

// Split splits the input on the next delimiter or space.
func Split(b []byte) ([]byte, []byte) {
	if len(b) == 0 {
		return nil, nil
	}
	n := 0
	if isdelim(b[n]) {
		// split at the delimiter
		n++
	} else if isspace(b[n]) {
		for n < len(b) && isspace(b[n]) {
			n++
		}
	} else {
		// split at the first delimiter
		n = bytes.IndexAny(b, ",: \t")
		if n == -1 {
			// did not find a delimiter
			n = len(b)
		}
	}
	return b[:n], b[n:]
}

func isdelim(ch byte) bool {
	return ch == ',' || ch == ':'
}

func isspace(ch byte) bool {
	return ch == ' ' || ch == '\t'
}
