// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package lexers_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/playbymail/ottoapp/backend/services/reports/lexers"
)

func TestLexer_ReportFile(t *testing.T) {
	// Read the test data file
	input, err := os.ReadFile("testdata/0900-01.0987.scrubbed.txt")
	if err != nil {
		t.Fatalf("failed to read testdata file: %v", err)
	}

	// scan the input
	tokens := lexers.Scan(input)

	//// log some tokens for creating golden test files
	//for i, tok := range tokens {
	//	if i > 46 {
	//		break
	//	}
	//	fmt.Printf("{id: %3d, line: %3d, col: %3d, length: %3d, leading: %3d, trailing: %3d, kind: %-12q, text: %q},\n",
	//		i, tok.Value.Line, tok.Value.Col, tok.Value.Length(), len(tok.LeadingTrivia), len(tok.TrailingTrivia), tok.Value.Kind.String(), string(tok.Value.Bytes()))
	//}

	for i, expected := range []struct {
		id       int
		line     int
		col      int
		length   int
		leading  int
		trailing int
		kind     string
		text     string
	}{
		{id: 0, line: 1, col: 1, length: 5, leading: 0, trailing: 1, kind: "Tribe", text: "Tribe"},
		{id: 1, line: 1, col: 7, length: 4, leading: 0, trailing: 0, kind: "Number", text: "0987"},
		{id: 2, line: 1, col: 11, length: 1, leading: 0, trailing: 1, kind: "Comma", text: ","},
		{id: 3, line: 1, col: 13, length: 1, leading: 0, trailing: 1, kind: "Comma", text: ","},
		{id: 4, line: 1, col: 15, length: 7, leading: 0, trailing: 1, kind: "Current", text: "Current"},
		{id: 5, line: 1, col: 23, length: 3, leading: 0, trailing: 1, kind: "Hex", text: "Hex"},
		{id: 6, line: 1, col: 27, length: 1, leading: 0, trailing: 1, kind: "Equals", text: "="},
		{id: 7, line: 1, col: 29, length: 2, leading: 0, trailing: 1, kind: "Grid", text: "QQ"},
		{id: 8, line: 1, col: 32, length: 4, leading: 0, trailing: 0, kind: "Number", text: "1509"},
		{id: 9, line: 1, col: 36, length: 1, leading: 0, trailing: 1, kind: "Comma", text: ","},
		{id: 10, line: 1, col: 38, length: 1, leading: 0, trailing: 0, kind: "LeftParen", text: "("},
		{id: 11, line: 1, col: 39, length: 8, leading: 0, trailing: 1, kind: "Previous", text: "Previous"},
		{id: 12, line: 1, col: 48, length: 3, leading: 0, trailing: 1, kind: "Hex", text: "Hex"},
		{id: 13, line: 1, col: 52, length: 1, leading: 0, trailing: 1, kind: "Equals", text: "="},
		{id: 14, line: 1, col: 54, length: 2, leading: 0, trailing: 1, kind: "Grid", text: "QQ"},
		{id: 15, line: 1, col: 57, length: 4, leading: 0, trailing: 0, kind: "Number", text: "1410"},
		{id: 16, line: 1, col: 61, length: 1, leading: 0, trailing: 0, kind: "RightParen", text: ")"},
		{id: 17, line: 1, col: 62, length: 1, leading: 0, trailing: 0, kind: "EOL", text: "\n"},
		{id: 18, line: 2, col: 1, length: 7, leading: 0, trailing: 1, kind: "Current", text: "Current"},
		{id: 19, line: 2, col: 9, length: 4, leading: 0, trailing: 1, kind: "Turn", text: "Turn"},
		{id: 20, line: 2, col: 14, length: 3, leading: 0, trailing: 0, kind: "Number", text: "900"},
		{id: 21, line: 2, col: 17, length: 1, leading: 0, trailing: 0, kind: "Dash", text: "-"},
		{id: 22, line: 2, col: 18, length: 2, leading: 0, trailing: 1, kind: "Number", text: "01"},
		{id: 23, line: 2, col: 21, length: 1, leading: 0, trailing: 0, kind: "LeftParen", text: "("},
		{id: 24, line: 2, col: 22, length: 1, leading: 0, trailing: 0, kind: "Hash", text: "#"},
		{id: 25, line: 2, col: 23, length: 1, leading: 0, trailing: 0, kind: "Number", text: "1"},
		{id: 26, line: 2, col: 24, length: 1, leading: 0, trailing: 0, kind: "RightParen", text: ")"},
		{id: 27, line: 2, col: 25, length: 1, leading: 0, trailing: 1, kind: "Comma", text: ","},
		{id: 28, line: 2, col: 27, length: 6, leading: 0, trailing: 0, kind: "Season", text: "Spring"},
		{id: 29, line: 2, col: 33, length: 1, leading: 0, trailing: 1, kind: "Comma", text: ","},
		{id: 30, line: 2, col: 35, length: 4, leading: 0, trailing: 1, kind: "Weather", text: "FINE"},
		{id: 31, line: 2, col: 40, length: 4, leading: 0, trailing: 1, kind: "Next", text: "Next"},
		{id: 32, line: 2, col: 45, length: 4, leading: 0, trailing: 1, kind: "Turn", text: "Turn"},
		{id: 33, line: 2, col: 50, length: 3, leading: 0, trailing: 0, kind: "Number", text: "900"},
		{id: 34, line: 2, col: 53, length: 1, leading: 0, trailing: 0, kind: "Dash", text: "-"},
		{id: 35, line: 2, col: 54, length: 2, leading: 0, trailing: 1, kind: "Number", text: "02"},
		{id: 36, line: 2, col: 57, length: 1, leading: 0, trailing: 0, kind: "LeftParen", text: "("},
		{id: 37, line: 2, col: 58, length: 1, leading: 0, trailing: 0, kind: "Hash", text: "#"},
		{id: 38, line: 2, col: 59, length: 1, leading: 0, trailing: 0, kind: "Number", text: "2"},
		{id: 39, line: 2, col: 60, length: 1, leading: 0, trailing: 0, kind: "RightParen", text: ")"},
		{id: 40, line: 2, col: 61, length: 1, leading: 0, trailing: 1, kind: "Comma", text: ","},
		{id: 41, line: 2, col: 63, length: 2, leading: 0, trailing: 0, kind: "Number", text: "12"},
		{id: 42, line: 2, col: 65, length: 1, leading: 0, trailing: 0, kind: "Slash", text: "/"},
		{id: 43, line: 2, col: 66, length: 2, leading: 0, trailing: 0, kind: "Number", text: "12"},
		{id: 44, line: 2, col: 68, length: 1, leading: 0, trailing: 0, kind: "Slash", text: "/"},
		{id: 45, line: 2, col: 69, length: 4, leading: 0, trailing: 0, kind: "Number", text: "2025"},
		{id: 46, line: 2, col: 73, length: 1, leading: 0, trailing: 0, kind: "EOL", text: "\n"},
	} {
		tok := tokens[i]
		if expected.line != tok.Value.Line {
			t.Errorf("%d: token %d: line: want %d, got %d\n", expected.id, i, expected.line, tok.Value.Line)
		} else if expected.col != tok.Value.Col {
			t.Errorf("%d: token %d: col: want %d, got %d\n", expected.id, i, expected.col, tok.Value.Col)
		} else if expected.length != tok.Value.Length() {
			t.Errorf("%d: token %d: length: want %d, got %d\n", expected.id, i, expected.length, tok.Value.Length())
		} else if expected.kind != tok.Value.Kind.String() {
			t.Errorf("%d: token %d: kind: want %q, got %q\n", expected.id, i, expected.kind, tok.Value.Kind.String())
		} else if text := string(tok.Value.Bytes()); expected.text != text {
			t.Errorf("%d: token %d: text: want %q, got %q\n", expected.id, i, expected.text, text)
		} else if count := len(tok.LeadingTrivia); expected.leading != count {
			t.Errorf("%d: token %d: leading trivia: want %d, got %d\n", expected.id, i, expected.leading, count)
		} else if count := len(tok.TrailingTrivia); expected.trailing != count {
			t.Errorf("%d: token %d: trailing trivia: want %d, got %d\n", expected.id, i, expected.trailing, count)
		}
	}

	// verify that only first tokens on a line have leading trivia
	prevLine := 0
	for _, tok := range tokens {
		if tok.Value.Line == prevLine && len(tok.LeadingTrivia) != 0 {
			t.Error("only first tokens should have leading trivia")
		}
		prevLine = tok.Value.Line
	}

	// verify that EOL tokens never have trailing trivia
	for _, tok := range tokens {
		if tok.Value.Value[0] == '\n' && len(tok.TrailingTrivia) != 0 {
			t.Error("EOL token should not have trailing trivia")
		}
	}

	// verify the round trip
	if output := lexers.ToSource(tokens...); !bytes.Equal(input, output) {
		t.Errorf("round trip: failed: input %d, output %d\n", len(input), len(output))
		t.Errorf("round trip: mismatch:\n%s", cmp.Diff(string(input), string(output)))
	}
}
