// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package cst_test

import (
	"os"
	"testing"

	"github.com/playbymail/ottoapp/backend/services/reports/cst"
	"github.com/playbymail/ottoapp/backend/services/reports/lexers"
)

func TestParse_SingleUnitLine(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantUnits   int
		wantKeyword string
		wantUnitID  string
		wantCurrHex string // "Grid Number" or "Text/Text" or "## Number"
		wantPrevHex string
		wantErrors  int
	}{
		{
			name:        "tribe with grid coords",
			input:       "Tribe 0987, , Current Hex = QQ 1509, (Previous Hex = QQ 1410)\n",
			wantUnits:   1,
			wantKeyword: "Tribe",
			wantUnitID:  "0987",
			wantCurrHex: "QQ 1509",
			wantPrevHex: "QQ 1410",
			wantErrors:  0,
		},
		{
			name:        "element with mixed grids",
			input:       "Element 0987e1, , Current Hex = QQ 1407, (Previous Hex = FF 1410)\n",
			wantUnits:   1,
			wantKeyword: "Element",
			wantUnitID:  "0987e1",
			wantCurrHex: "QQ 1407",
			wantPrevHex: "FF 1410",
			wantErrors:  0,
		},
		{
			name:        "garrison same hex",
			input:       "Garrison 0987g1, , Current Hex = QQ 1408, (Previous Hex = QQ 1408)\n",
			wantUnits:   1,
			wantKeyword: "Garrison",
			wantUnitID:  "0987g1",
			wantCurrHex: "QQ 1408",
			wantPrevHex: "QQ 1408",
			wantErrors:  0,
		},
		{
			name:        "tribe with N/A previous",
			input:       "Tribe 1234, , Current Hex = AA 0101, (Previous Hex = N/A)\n",
			wantUnits:   1,
			wantKeyword: "Tribe",
			wantUnitID:  "1234",
			wantCurrHex: "AA 0101",
			wantPrevHex: "N/A",
			wantErrors:  0,
		},
		{
			name:        "fleet with obscured previous",
			input:       "Fleet 5678f1, , Current Hex = BB 0202, (Previous Hex = ## 9999)\n",
			wantUnits:   1,
			wantKeyword: "Fleet",
			wantUnitID:  "5678f1",
			wantCurrHex: "BB 0202",
			wantPrevHex: "## 9999",
			wantErrors:  0,
		},
		{
			name:        "courier unit",
			input:       "Courier 9999c1, , Current Hex = ZZ 0505, (Previous Hex = ZZ 0404)\n",
			wantUnits:   1,
			wantKeyword: "Courier",
			wantUnitID:  "9999c1",
			wantCurrHex: "ZZ 0505",
			wantPrevHex: "ZZ 0404",
			wantErrors:  0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tokens := lexers.Scan([]byte(tc.input))
			result := cst.Parse(tokens)

			if got := len(result.Sections); got != tc.wantUnits {
				t.Errorf("units: want %d, got %d", tc.wantUnits, got)
				return
			}

			if got := len(result.Errors()); got != tc.wantErrors {
				t.Errorf("errors: want %d, got %d", tc.wantErrors, got)
				for _, err := range result.Errors() {
					t.Logf("  error: %v", err)
				}
			}

			if tc.wantUnits == 0 {
				return
			}

			section := result.Sections[0]
			unitLine := section.UnitLine

			if unitLine.Keyword == nil {
				t.Fatal("keyword is nil")
			}
			if got := string(unitLine.Keyword.Bytes()); got != tc.wantKeyword {
				t.Errorf("keyword: want %q, got %q", tc.wantKeyword, got)
			}

			if unitLine.UnitID == nil {
				t.Fatal("unit id is nil")
			}
			if got := string(unitLine.UnitID.Bytes()); got != tc.wantUnitID {
				t.Errorf("unit id: want %q, got %q", tc.wantUnitID, got)
			}

			currHex := formatCoords(unitLine.CurrentHex)
			if currHex != tc.wantCurrHex {
				t.Errorf("current hex: want %q, got %q", tc.wantCurrHex, currHex)
			}

			prevHex := formatCoords(unitLine.PreviousHex)
			if prevHex != tc.wantPrevHex {
				t.Errorf("previous hex: want %q, got %q", tc.wantPrevHex, prevHex)
			}
		})
	}
}

func TestParse_MultipleUnits(t *testing.T) {
	input := `Tribe 0987, , Current Hex = QQ 1509, (Previous Hex = QQ 1410)
Element 0987e1, , Current Hex = QQ 1407, (Previous Hex = FF 1410)
Garrison 0987g1, , Current Hex = QQ 1408, (Previous Hex = QQ 1408)
`
	tokens := lexers.Scan([]byte(input))
	result := cst.Parse(tokens)

	if got := len(result.Sections); got != 3 {
		t.Fatalf("sections: want 3, got %d", got)
	}

	expected := []struct {
		keyword string
		unitID  string
	}{
		{"Tribe", "0987"},
		{"Element", "0987e1"},
		{"Garrison", "0987g1"},
	}

	for i, exp := range expected {
		section := result.Sections[i]
		if section.UnitLine.Keyword == nil {
			t.Errorf("section %d: keyword is nil", i)
			continue
		}
		if got := string(section.UnitLine.Keyword.Bytes()); got != exp.keyword {
			t.Errorf("section %d: keyword: want %q, got %q", i, exp.keyword, got)
		}
		if got := string(section.UnitLine.UnitID.Bytes()); got != exp.unitID {
			t.Errorf("section %d: unit id: want %q, got %q", i, exp.unitID, got)
		}
	}

	if len(result.Errors()) != 0 {
		t.Errorf("unexpected errors: %v", result.Errors())
	}
}

func TestParse_ReportFile(t *testing.T) {
	input, err := os.ReadFile("../lexers/testdata/0900-01.0987.scrubbed.txt")
	if err != nil {
		t.Fatalf("failed to read testdata file: %v", err)
	}

	tokens := lexers.Scan(input)
	result := cst.Parse(tokens)

	// The file has lines starting with unit keywords on lines 1, 3, 8, 9, 15, 17, 18
	// Lines 1, 8, 15, 17 are valid unit lines
	// Lines 3, 9, 18 are "Tribe Movement:" lines that start with "Tribe" but fail to parse
	expectedSections := []struct {
		keyword   string
		unitID    string
		wantError bool
	}{
		{"Tribe", "0987", false},      // line 1: valid unit line
		{"Tribe", "Movement", true},   // line 3: "Tribe Movement:" - fails after Movement
		{"Element", "0987e1", false},  // line 8: valid unit line
		{"Tribe", "Movement", true},   // line 9: "Tribe Movement:" - fails after Movement
		{"Garrison", "0987g1", false}, // line 15: valid unit line
		{"Tribe", "1987", false},      // line 17: valid unit line
		{"Tribe", "Movement", true},   // line 18: "Tribe Movement:" - fails after Movement
	}

	if got := len(result.Sections); got != len(expectedSections) {
		t.Fatalf("sections: want %d, got %d", len(expectedSections), got)
	}

	for i, exp := range expectedSections {
		section := result.Sections[i]
		if section.UnitLine.Keyword == nil {
			t.Errorf("section %d: keyword is nil", i)
			continue
		}
		if got := string(section.UnitLine.Keyword.Bytes()); got != exp.keyword {
			t.Errorf("section %d: keyword: want %q, got %q", i, exp.keyword, got)
		}
		if section.UnitLine.UnitID != nil {
			if got := string(section.UnitLine.UnitID.Bytes()); got != exp.unitID {
				t.Errorf("section %d: unit id: want %q, got %q", i, exp.unitID, got)
			}
		}
		hasErrors := len(section.Errors()) > 0
		if hasErrors != exp.wantError {
			t.Errorf("section %d: wantError=%v, got errors=%v", i, exp.wantError, section.Errors())
		}
	}

	// Expect errors for non-unit lines and failed "Tribe Movement" parses
	if len(result.Errors()) == 0 {
		t.Error("expected errors for non-unit content and failed parses")
	} else {
		t.Logf("errors (expected): %d", len(result.Errors()))
	}
}

func TestParse_ErrorRecovery(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantUnits  int
		wantErrors int
	}{
		{
			name:       "missing unit id recovers to next line",
			input:      "Tribe , , Current Hex = QQ 1509, (Previous Hex = QQ 1410)\nTribe 0987, , Current Hex = QQ 1509, (Previous Hex = QQ 1410)\n",
			wantUnits:  2,
			wantErrors: 1,
		},
		{
			name:       "garbage before unit line",
			input:      "some garbage here\nTribe 0987, , Current Hex = QQ 1509, (Previous Hex = QQ 1410)\n",
			wantUnits:  1,
			wantErrors: 1,
		},
		{
			name:       "empty input",
			input:      "",
			wantUnits:  0,
			wantErrors: 0,
		},
		{
			name:       "only garbage",
			input:      "no units here at all\n",
			wantUnits:  0,
			wantErrors: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tokens := lexers.Scan([]byte(tc.input))
			result := cst.Parse(tokens)

			if got := len(result.Sections); got != tc.wantUnits {
				t.Errorf("units: want %d, got %d", tc.wantUnits, got)
			}

			if got := len(result.Errors()); got != tc.wantErrors {
				t.Errorf("errors: want %d, got %d", tc.wantErrors, got)
				for _, err := range result.Errors() {
					t.Logf("  error: %v", err)
				}
			}
		})
	}
}

func TestParseTurnLine(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		wantTurnYM1     string // first TurnYearMonth
		wantTurnNum1    string // first turn number (just the number)
		wantSeason      string
		wantWeather     string
		wantHasNextTurn bool
		wantTurnYM2     string // second TurnYearMonth (if present)
		wantTurnNum2    string // second turn number (if present)
		wantReportDate  string // DD/MM/YYYY (if present)
		wantErrors      int
	}{
		{
			name:            "full turn line with next turn",
			input:           "Current Turn 899-12 (#0), Winter, FINE Next Turn 900-01 (#1), 28/11/2025\n",
			wantTurnYM1:     "899-12",
			wantTurnNum1:    "0",
			wantSeason:      "Winter",
			wantWeather:     "FINE",
			wantHasNextTurn: true,
			wantTurnYM2:     "900-01",
			wantTurnNum2:    "1",
			wantReportDate:  "28/11/2025",
			wantErrors:      0,
		},
		{
			name:            "turn line without next turn",
			input:           "Current Turn 900-01 (#1), Spring, FINE\n",
			wantTurnYM1:     "900-01",
			wantTurnNum1:    "1",
			wantSeason:      "Spring",
			wantWeather:     "FINE",
			wantHasNextTurn: false,
			wantErrors:      0,
		},
		{
			name:            "turn line from test file",
			input:           "Current Turn 900-01 (#1), Spring, FINE Next Turn 900-02 (#2), 12/12/2025\n",
			wantTurnYM1:     "900-01",
			wantTurnNum1:    "1",
			wantSeason:      "Spring",
			wantWeather:     "FINE",
			wantHasNextTurn: true,
			wantTurnYM2:     "900-02",
			wantTurnNum2:    "2",
			wantReportDate:  "12/12/2025",
			wantErrors:      0,
		},
		{
			name:            "summer season",
			input:           "Current Turn 900-06 (#6), Summer, FINE\n",
			wantTurnYM1:     "900-06",
			wantTurnNum1:    "6",
			wantSeason:      "Summer",
			wantWeather:     "FINE",
			wantHasNextTurn: false,
			wantErrors:      0,
		},
		{
			name:            "fall/autumn season",
			input:           "Current Turn 900-09 (#9), Fall, FINE\n",
			wantTurnYM1:     "900-09",
			wantTurnNum1:    "9",
			wantSeason:      "Fall",
			wantWeather:     "FINE",
			wantHasNextTurn: false,
			wantErrors:      0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tokens := lexers.Scan([]byte(tc.input))
			result := cst.ParseTurnLine(tokens)

			if got := len(result.Errors()); got != tc.wantErrors {
				t.Errorf("errors: want %d, got %d", tc.wantErrors, got)
				for _, err := range result.Errors() {
					t.Logf("  error: %v", err)
				}
				if tc.wantErrors == 0 {
					return
				}
			}

			if result.TurnYearMonth1 == nil {
				t.Fatal("TurnYearMonth1 is nil")
			}
			if got := string(result.TurnYearMonth1.Bytes()); got != tc.wantTurnYM1 {
				t.Errorf("TurnYearMonth1: want %q, got %q", tc.wantTurnYM1, got)
			}

			if result.TurnNumber1 == nil || result.TurnNumber1.Number == nil {
				t.Fatal("TurnNumber1 or its Number is nil")
			}
			if got := string(result.TurnNumber1.Number.Bytes()); got != tc.wantTurnNum1 {
				t.Errorf("TurnNumber1: want %q, got %q", tc.wantTurnNum1, got)
			}

			if result.Season == nil {
				t.Fatal("Season is nil")
			}
			if got := string(result.Season.Bytes()); got != tc.wantSeason {
				t.Errorf("Season: want %q, got %q", tc.wantSeason, got)
			}

			if result.Weather == nil {
				t.Fatal("Weather is nil")
			}
			if got := string(result.Weather.Bytes()); got != tc.wantWeather {
				t.Errorf("Weather: want %q, got %q", tc.wantWeather, got)
			}

			hasNextTurn := result.Next != nil
			if hasNextTurn != tc.wantHasNextTurn {
				t.Errorf("hasNextTurn: want %v, got %v", tc.wantHasNextTurn, hasNextTurn)
			}

			if tc.wantHasNextTurn {
				if result.TurnYearMonth2 == nil {
					t.Fatal("TurnYearMonth2 is nil")
				}
				if got := string(result.TurnYearMonth2.Bytes()); got != tc.wantTurnYM2 {
					t.Errorf("TurnYearMonth2: want %q, got %q", tc.wantTurnYM2, got)
				}

				if result.TurnNumber2 == nil || result.TurnNumber2.Number == nil {
					t.Fatal("TurnNumber2 or its Number is nil")
				}
				if got := string(result.TurnNumber2.Number.Bytes()); got != tc.wantTurnNum2 {
					t.Errorf("TurnNumber2: want %q, got %q", tc.wantTurnNum2, got)
				}

				if result.ReportDate == nil {
					t.Fatal("ReportDate is nil")
				}
				gotDate := string(result.ReportDate.Day.Bytes()) + "/" +
					string(result.ReportDate.Month.Bytes()) + "/" +
					string(result.ReportDate.Year.Bytes())
				if gotDate != tc.wantReportDate {
					t.Errorf("ReportDate: want %q, got %q", tc.wantReportDate, gotDate)
				}
			}
		})
	}
}

// formatCoords returns a string representation of coordinates for testing.
func formatCoords(coords cst.CoordsNode) string {
	if coords == nil {
		return "<nil>"
	}
	switch c := coords.(type) {
	case *cst.GridCoordsNode:
		if c.Grid == nil || c.Number == nil {
			return "<incomplete grid>"
		}
		return string(c.Grid.Bytes()) + " " + string(c.Number.Bytes())
	case *cst.NACoordsNode:
		if c.Text == nil {
			return "<incomplete na>"
		}
		return string(c.Text.Bytes())
	case *cst.ObscuredCoordsNode:
		if c.Grid == nil || c.Number == nil {
			return "<incomplete obscured>"
		}
		return string(c.Grid.Bytes()) + " " + string(c.Number.Bytes())
	case *cst.ErrorCoordsNode:
		return "<error: " + c.Message + ">"
	default:
		return "<unknown>"
	}
}
