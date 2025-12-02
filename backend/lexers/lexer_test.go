// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package lexers

import (
	"bytes"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLexer_ReportFile(t *testing.T) {
	// Read the test data file
	data, err := os.ReadFile("testdata/0899-12.0500.scrubbed.txt")
	if err != nil {
		t.Fatalf("failed to read testdata file: %v", err)
	}

	// Create lexer
	lexer, err := New(data)
	if err != nil {
		t.Fatalf("failed to create lexer: %v", err)
	}

	// Lex all tokens
	tokens := []*Token_t{}
	var prevToken *Token_t
	for {
		tok := lexer.Next()
		tokens = append(tokens, tok)
		if tok.Kind == TokenEOF {
			break
		}

		// Detect infinite loop: if position hasn't changed and we got the same token kind
		if prevToken != nil && tok.Span.Start == prevToken.Span.Start {
			t.Fatalf("infinite loop: position stuck at %d, token kind %s", tok.Span.Start, tok.Kind)
		}
		prevToken = tok
	}

	// Verify we got tokens (even on empty input we should get at least the EOF token)
	if len(tokens) == 0 {
		t.Fatal("expected tokens but got none")
	}

	//// log some tokens for creating golden test files
	//for i, tok := range tokens {
	//	if i > 46 {
	//		break
	//	}
	//	fmt.Printf("{id: %3d, line: %3d, col: %3d, start: %4d, end: %4d, leading: %2d, trailing: %2d, kind: %-12q, text: %q},\n",
	//		i, tok.Span.Line, tok.Span.Col, tok.Span.Start, tok.Span.End, len(tok.LeadingTrivia), len(tok.TrailingTrivia), tok.Kind, string(tok.Text(data)))
	//}

	type expectedToken_t struct {
		id       int
		line     int
		col      int
		start    int
		end      int
		kind     string
		text     string
		leading  int
		trailing int
	}
	for i, expected := range []expectedToken_t{
		{id: 0, line: 1, col: 1, start: 0, end: 5, leading: 0, trailing: 1, kind: "tribe", text: "Tribe"},
		{id: 1, line: 1, col: 7, start: 6, end: 10, leading: 0, trailing: 0, kind: "number", text: "0500"},
		{id: 2, line: 1, col: 11, start: 10, end: 11, leading: 0, trailing: 0, kind: ",", text: ","},
		{id: 3, line: 1, col: 13, start: 12, end: 13, leading: 1, trailing: 0, kind: ",", text: ","},
		{id: 4, line: 1, col: 15, start: 14, end: 21, leading: 1, trailing: 1, kind: "current", text: "Current"},
		{id: 5, line: 1, col: 23, start: 22, end: 25, leading: 0, trailing: 1, kind: "hex", text: "Hex"},
		{id: 6, line: 1, col: 27, start: 26, end: 27, leading: 0, trailing: 0, kind: "=", text: "="},
		{id: 7, line: 1, col: 29, start: 28, end: 30, leading: 1, trailing: 1, kind: "grid", text: "HK"},
		{id: 8, line: 1, col: 32, start: 31, end: 35, leading: 0, trailing: 0, kind: "number", text: "1315"},
		{id: 9, line: 1, col: 36, start: 35, end: 36, leading: 0, trailing: 0, kind: ",", text: ","},
		{id: 10, line: 1, col: 38, start: 37, end: 38, leading: 1, trailing: 0, kind: "(", text: "("},
		{id: 11, line: 1, col: 39, start: 38, end: 46, leading: 0, trailing: 1, kind: "previous", text: "Previous"},
		{id: 12, line: 1, col: 48, start: 47, end: 50, leading: 0, trailing: 1, kind: "hex", text: "Hex"},
		{id: 13, line: 1, col: 52, start: 51, end: 52, leading: 0, trailing: 0, kind: "=", text: "="},
		{id: 14, line: 1, col: 54, start: 53, end: 55, leading: 1, trailing: 1, kind: "grid", text: "HK"},
		{id: 15, line: 1, col: 57, start: 56, end: 60, leading: 0, trailing: 0, kind: "number", text: "1315"},
		{id: 16, line: 1, col: 61, start: 60, end: 61, leading: 0, trailing: 0, kind: ")", text: ")"},
		{id: 17, line: 1, col: 62, start: 61, end: 62, leading: 0, trailing: 0, kind: "eol", text: "\n"},
		{id: 18, line: 2, col: 1, start: 62, end: 69, leading: 0, trailing: 1, kind: "current", text: "Current"},
		{id: 19, line: 2, col: 9, start: 70, end: 74, leading: 0, trailing: 1, kind: "turn", text: "Turn"},
		{id: 20, line: 2, col: 14, start: 75, end: 78, leading: 0, trailing: 0, kind: "number", text: "899"},
		{id: 21, line: 2, col: 17, start: 78, end: 79, leading: 0, trailing: 0, kind: "-", text: "-"},
		{id: 22, line: 2, col: 18, start: 79, end: 81, leading: 0, trailing: 1, kind: "number", text: "12"},
		{id: 23, line: 2, col: 21, start: 82, end: 83, leading: 0, trailing: 0, kind: "(", text: "("},
		{id: 24, line: 2, col: 22, start: 83, end: 84, leading: 0, trailing: 0, kind: "#", text: "#"},
		{id: 25, line: 2, col: 23, start: 84, end: 85, leading: 0, trailing: 0, kind: "number", text: "0"},
		{id: 26, line: 2, col: 24, start: 85, end: 86, leading: 0, trailing: 0, kind: ")", text: ")"},
		{id: 27, line: 2, col: 25, start: 86, end: 87, leading: 0, trailing: 0, kind: ",", text: ","},
		{id: 28, line: 2, col: 27, start: 88, end: 94, leading: 1, trailing: 0, kind: "season", text: "Winter"},
		{id: 29, line: 2, col: 33, start: 94, end: 95, leading: 0, trailing: 0, kind: ",", text: ","},
		{id: 30, line: 2, col: 35, start: 96, end: 100, leading: 1, trailing: 1, kind: "weather", text: "FINE"},
		{id: 31, line: 2, col: 40, start: 101, end: 105, leading: 0, trailing: 1, kind: "next", text: "Next"},
		{id: 32, line: 2, col: 45, start: 106, end: 110, leading: 0, trailing: 1, kind: "turn", text: "Turn"},
		{id: 33, line: 2, col: 50, start: 111, end: 114, leading: 0, trailing: 0, kind: "number", text: "900"},
		{id: 34, line: 2, col: 53, start: 114, end: 115, leading: 0, trailing: 0, kind: "-", text: "-"},
		{id: 35, line: 2, col: 54, start: 115, end: 117, leading: 0, trailing: 1, kind: "number", text: "01"},
		{id: 36, line: 2, col: 57, start: 118, end: 119, leading: 0, trailing: 0, kind: "(", text: "("},
		{id: 37, line: 2, col: 58, start: 119, end: 120, leading: 0, trailing: 0, kind: "#", text: "#"},
		{id: 38, line: 2, col: 59, start: 120, end: 121, leading: 0, trailing: 0, kind: "number", text: "1"},
		{id: 39, line: 2, col: 60, start: 121, end: 122, leading: 0, trailing: 0, kind: ")", text: ")"},
		{id: 40, line: 2, col: 61, start: 122, end: 123, leading: 0, trailing: 0, kind: ",", text: ","},
		{id: 41, line: 2, col: 63, start: 124, end: 126, leading: 1, trailing: 0, kind: "number", text: "28"},
		{id: 42, line: 2, col: 65, start: 126, end: 127, leading: 0, trailing: 0, kind: "/", text: "/"},
		{id: 43, line: 2, col: 66, start: 127, end: 129, leading: 0, trailing: 0, kind: "number", text: "11"},
		{id: 44, line: 2, col: 68, start: 129, end: 130, leading: 0, trailing: 0, kind: "/", text: "/"},
		{id: 45, line: 2, col: 69, start: 130, end: 134, leading: 0, trailing: 0, kind: "number", text: "2025"},
		{id: 46, line: 2, col: 73, start: 134, end: 135, leading: 0, trailing: 0, kind: "eol", text: "\n"},
	} {
		tok := tokens[i]
		if expected.kind != tok.Kind.String() {
			t.Errorf("%d: token %d: kind: want %q, got %q\n", expected.id, i, expected.kind, tok.Kind)
		} else if expected.start != tok.Span.Start {
			t.Errorf("%d: token %d: start: want %d, got %d\n", expected.id, i, expected.start, tok.Span.Start)
		} else if expected.end != tok.Span.End {
			t.Errorf("%d: token %d: end: want %d, got %d\n", expected.id, i, expected.end, tok.Span.End)
		} else if expected.line != tok.Span.Line {
			t.Errorf("%d: token %d: line: want %d, got %d\n", expected.id, i, expected.line, tok.Span.Line)
		} else if expected.col != tok.Span.Col {
			t.Errorf("%d: token %d: col: want %d, got %d\n", expected.id, i, expected.col, tok.Span.Col)
		} else if text := string(tok.Text(data)); expected.text != text {
			t.Errorf("%d: token %d: text: want %q, got %q\n", expected.id, i, expected.text, text)
		} else if count := len(tok.LeadingTrivia); expected.leading != count {
			t.Errorf("%d: token %d: leading trivia: want %d, got %d\n", expected.id, i, expected.leading, count)
		} else if count := len(tok.TrailingTrivia); expected.trailing != count {
			t.Errorf("%d: token %d: trailing trivia: want %d, got %d\n", expected.id, i, expected.trailing, count)
		}
	}

	// Verify last token is EOF (Kind will be TokenEOF)
	lastToken := tokens[len(tokens)-1]
	if lastToken.Kind != TokenEOF {
		t.Errorf("expected last token to be EOF, got %s", lastToken.Kind)
	}

	// Verify we have both Text and EOL tokens
	hasText := false
	hasEOL := false
	for _, tok := range tokens {
		if tok.Kind == TokenText {
			hasText = true
		}
		if tok.Kind == TokenEOL {
			hasEOL = true
		}
	}

	if !hasText {
		t.Error("expected to find TokenText")
	}
	if !hasEOL {
		t.Error("expected to find TokenEOL")
	}

	// Verify EOL tokens never have trailing trivia
	for _, tok := range tokens {
		if tok.Kind == TokenEOL && len(tok.TrailingTrivia) > 0 {
			t.Error("EOL token should not have trailing trivia")
		}
	}

	// Verify the round trip
	if input, output := data, tokensToSource(data, tokens); !bytes.Equal(input, output) {
		t.Errorf("round trip: failed: input %d, output %d\n", len(input), len(output))
		t.Errorf("round trip: mismatch:\n%s", cmp.Diff(string(input), string(output)))
	}
}

// helper function to rebuild the source
func tokensToSource(src []byte, tokens []*Token_t) []byte {
	b := &bytes.Buffer{}
	for _, tok := range tokens {
		b.Write(tok.LeadingTriviaText(src))
		b.Write(tok.Text(src))
		b.Write(tok.TrailingTriviaText(src))
	}
	return b.Bytes()
}
