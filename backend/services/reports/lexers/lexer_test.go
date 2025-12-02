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
	//		//break
	//	}
	//	line, col := tok.Position()
	//	fmt.Printf("{id: %3d, line: %3d, col: %3d, length: %3d, leading: %3d, trailing: %3d, kind: %-16q, text: %q},\n",
	//		i, line, col, tok.Length(), len(tok.LeadingTrivia), len(tok.TrailingTrivia), tok.Kind.String(), string(tok.Bytes()))
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
		{id: 1, line: 1, col: 7, length: 4, leading: 0, trailing: 0, kind: "UnitId", text: "0987"},
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
		{id: 20, line: 2, col: 14, length: 6, leading: 0, trailing: 1, kind: "TurnYearMonth", text: "900-01"},
		{id: 21, line: 2, col: 21, length: 1, leading: 0, trailing: 0, kind: "LeftParen", text: "("},
		{id: 22, line: 2, col: 22, length: 1, leading: 0, trailing: 0, kind: "Hash", text: "#"},
		{id: 23, line: 2, col: 23, length: 1, leading: 0, trailing: 0, kind: "Number", text: "1"},
		{id: 24, line: 2, col: 24, length: 1, leading: 0, trailing: 0, kind: "RightParen", text: ")"},
		{id: 25, line: 2, col: 25, length: 1, leading: 0, trailing: 1, kind: "Comma", text: ","},
		{id: 26, line: 2, col: 27, length: 6, leading: 0, trailing: 0, kind: "Season", text: "Spring"},
		{id: 27, line: 2, col: 33, length: 1, leading: 0, trailing: 1, kind: "Comma", text: ","},
		{id: 28, line: 2, col: 35, length: 4, leading: 0, trailing: 1, kind: "Weather", text: "FINE"},
		{id: 29, line: 2, col: 40, length: 4, leading: 0, trailing: 1, kind: "Next", text: "Next"},
		{id: 30, line: 2, col: 45, length: 4, leading: 0, trailing: 1, kind: "Turn", text: "Turn"},
		{id: 31, line: 2, col: 50, length: 6, leading: 0, trailing: 1, kind: "TurnYearMonth", text: "900-02"},
		{id: 32, line: 2, col: 57, length: 1, leading: 0, trailing: 0, kind: "LeftParen", text: "("},
		{id: 33, line: 2, col: 58, length: 1, leading: 0, trailing: 0, kind: "Hash", text: "#"},
		{id: 34, line: 2, col: 59, length: 1, leading: 0, trailing: 0, kind: "Number", text: "2"},
		{id: 35, line: 2, col: 60, length: 1, leading: 0, trailing: 0, kind: "RightParen", text: ")"},
		{id: 36, line: 2, col: 61, length: 1, leading: 0, trailing: 1, kind: "Comma", text: ","},
		{id: 37, line: 2, col: 63, length: 2, leading: 0, trailing: 0, kind: "Number", text: "12"},
		{id: 38, line: 2, col: 65, length: 1, leading: 0, trailing: 0, kind: "Slash", text: "/"},
		{id: 39, line: 2, col: 66, length: 2, leading: 0, trailing: 0, kind: "Number", text: "12"},
		{id: 40, line: 2, col: 68, length: 1, leading: 0, trailing: 0, kind: "Slash", text: "/"},
		{id: 41, line: 2, col: 69, length: 4, leading: 0, trailing: 0, kind: "Number", text: "2025"},
		{id: 42, line: 2, col: 73, length: 1, leading: 0, trailing: 0, kind: "EOL", text: "\n"},
		{id: 43, line: 3, col: 1, length: 5, leading: 0, trailing: 1, kind: "Tribe", text: "Tribe"},
		{id: 44, line: 3, col: 7, length: 8, leading: 0, trailing: 0, kind: "Text", text: "Movement"},
		{id: 45, line: 3, col: 15, length: 1, leading: 0, trailing: 1, kind: "Colon", text: ":"},
		{id: 46, line: 3, col: 17, length: 4, leading: 0, trailing: 1, kind: "Text", text: "Move"},
		{id: 47, line: 3, col: 22, length: 1, leading: 0, trailing: 0, kind: "Text", text: "N"},
		{id: 48, line: 3, col: 23, length: 1, leading: 0, trailing: 0, kind: "Dash", text: "-"},
		{id: 49, line: 3, col: 24, length: 2, leading: 0, trailing: 0, kind: "Grid", text: "PR"},
		{id: 50, line: 3, col: 26, length: 1, leading: 0, trailing: 1, kind: "Comma", text: ","},
		{id: 51, line: 3, col: 28, length: 1, leading: 0, trailing: 0, kind: "Backslash", text: "\\"},
		{id: 52, line: 3, col: 29, length: 2, leading: 0, trailing: 0, kind: "Grid", text: "NE"},
		{id: 53, line: 3, col: 31, length: 1, leading: 0, trailing: 0, kind: "Dash", text: "-"},
		{id: 54, line: 3, col: 32, length: 2, leading: 0, trailing: 0, kind: "Grid", text: "GH"},
		{id: 55, line: 3, col: 34, length: 1, leading: 0, trailing: 1, kind: "Comma", text: ","},
		{id: 56, line: 3, col: 36, length: 1, leading: 0, trailing: 1, kind: "Comma", text: ","},
		{id: 57, line: 3, col: 38, length: 1, leading: 0, trailing: 1, kind: "Text", text: "L"},
		{id: 58, line: 3, col: 40, length: 2, leading: 0, trailing: 0, kind: "Grid", text: "NE"},
		{id: 59, line: 3, col: 42, length: 1, leading: 0, trailing: 1, kind: "Comma", text: ","},
		{id: 60, line: 3, col: 44, length: 1, leading: 0, trailing: 0, kind: "Text", text: "N"},
		{id: 61, line: 3, col: 45, length: 1, leading: 0, trailing: 0, kind: "Backslash", text: "\\"},
		{id: 62, line: 3, col: 46, length: 1, leading: 0, trailing: 0, kind: "EOL", text: "\n"},
	} {
		tok := tokens[i]
		if line, col := tok.Position(); expected.line != line {
			t.Errorf("%d: token %d: line: want %d, got %d\n", expected.id, i, expected.line, line)
		} else if expected.col != col {
			t.Errorf("%d: token %d: col: want %d, got %d\n", expected.id, i, expected.col, col)
		} else if expected.length != tok.Length() {
			t.Errorf("%d: token %d: length: want %d, got %d\n", expected.id, i, expected.length, tok.Length())
		} else if expected.kind != tok.Kind.String() {
			t.Errorf("%d: token %d: kind: want %q, got %q\n", expected.id, i, expected.kind, tok.Kind.String())
		} else if text := string(tok.Bytes()); expected.text != text {
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
		line, _ := tok.Position()
		if line == prevLine && len(tok.LeadingTrivia) != 0 {
			t.Error("only first tokens should have leading trivia")
		}
		prevLine = line
	}

	// verify that EOL tokens never have trailing trivia
	for _, tok := range tokens {
		if tok.Kind == lexers.EOL && len(tok.TrailingTrivia) != 0 {
			t.Error("EOL token should not have trailing trivia")
		}
	}

	// verify the round trip
	if output := lexers.ToSource(tokens...); !bytes.Equal(input, output) {
		t.Errorf("round trip: failed: input %d, output %d\n", len(input), len(output))
		t.Errorf("round trip: mismatch:\n%s", cmp.Diff(string(input), string(output)))
	}
}
