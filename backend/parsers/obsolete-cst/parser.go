// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package obsolete_cst

func Parse(input []byte) *Node {
	var tokens []*Token

	pos, line, col := 0, 1, 1
	for pos < len(input) {
		var token Token

		// check for leading trivia
		if pos < len(input) && isTrivia(input[pos]) {
			token.LeadingTrivia.Line, token.LeadingTrivia.Col = line, col
			start := pos
			for pos < len(input) && isTrivia(input[pos]) {
				pos, line, col = pos+1, line, col+1
			}
			token.LeadingTrivia.Bytes = input[start:pos]
		}

		// check for text
		if pos < len(input) {
			token.Text.Line, token.Text.Col = line, col
			start := pos
			if input[pos] == '\n' {
				pos, line, col = pos+1, line+1, 1
			} else if isDelimiter(input[pos]) {
				pos, line, col = pos+1, line, col+1
			} else {
				for pos < len(input) && !(isTrivia(input[pos]) || isDelimiter(input[pos])) {
					pos, line, col = pos+1, line, col+1
				}
			}
			token.Text.Bytes = input[start:pos]
		}

		// check for trailing trivia
		if pos < len(input) && isTrivia(input[pos]) {
			token.TrailingTrivia.Line, token.TrailingTrivia.Col = line, col
			start := pos
			for pos < len(input) && isTrivia(input[pos]) {
				pos, line, col = pos+1, line, col+1
			}
			token.TrailingTrivia.Bytes = input[start:pos]
		}

		tokens = append(tokens, &token)
	}

	root := &Node{}

	return root
}

func isDelimiter(ch byte) bool {
	return ch == '\n'
}

func isTrivia(ch byte) bool {
	return ch == ' ' || ch == '\t'
}

type Node struct{}

type Token struct {
	LeadingTrivia  Span
	Text           Span
	TrailingTrivia Span
}

type Span struct {
	Start int // byte offset (inclusive)
	End   int // byte offset (exclusive)
	Line  int // 1-based
	Col   int // 1-based, in UTF-8 code points
}

// Length returns the length of the span.
func (s Span) Length() int {
	return s.End - s.Start
}

// Text is a helper for diagnostics / debugging.
func (s Span) Text(src []byte) []byte {
	return src[s.Start:s.End]
}
