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

func TestParse_UnitSectionWithTurnLine(t *testing.T) {
	input := `Tribe 0987, , Current Hex = QQ 1509, (Previous Hex = QQ 1410)
Current Turn 900-01 (#1), Spring, FINE Next Turn 900-02 (#2), 12/12/2025
Element 0987e1, , Current Hex = QQ 1407, (Previous Hex = FF 1410)
Current Turn 900-01 (#1), Spring, FINE
`
	tokens := lexers.Scan([]byte(input))
	result := cst.Parse(tokens)

	if got := len(result.Sections); got != 2 {
		t.Fatalf("sections: want 2, got %d", got)
	}

	// First section: Tribe with turn line including next turn
	section1 := result.Sections[0]
	if got := string(section1.UnitLine.Keyword.Bytes()); got != "Tribe" {
		t.Errorf("section 0: keyword: want %q, got %q", "Tribe", got)
	}
	if section1.TurnLine == nil {
		t.Fatal("section 0: TurnLine is nil")
	}
	if got := string(section1.TurnLine.TurnYearMonth1.Bytes()); got != "900-01" {
		t.Errorf("section 0: TurnYearMonth1: want %q, got %q", "900-01", got)
	}
	if got := string(section1.TurnLine.Season.Bytes()); got != "Spring" {
		t.Errorf("section 0: Season: want %q, got %q", "Spring", got)
	}
	if section1.TurnLine.Next == nil {
		t.Error("section 0: expected Next turn info")
	}
	if section1.TurnLine.ReportDate == nil {
		t.Error("section 0: expected ReportDate")
	}

	// Second section: Element with turn line without next turn
	section2 := result.Sections[1]
	if got := string(section2.UnitLine.Keyword.Bytes()); got != "Element" {
		t.Errorf("section 1: keyword: want %q, got %q", "Element", got)
	}
	if section2.TurnLine == nil {
		t.Fatal("section 1: TurnLine is nil")
	}
	if got := string(section2.TurnLine.TurnYearMonth1.Bytes()); got != "900-01" {
		t.Errorf("section 1: TurnYearMonth1: want %q, got %q", "900-01", got)
	}
	if section2.TurnLine.Next != nil {
		t.Error("section 1: expected no Next turn info")
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

	// The file has 4 unit sections:
	// - Tribe 0987 (line 1) with land movement (line 3)
	// - Element 0987e1 (line 8) with land movement (line 9)
	// - Garrison 0987g1 (line 15) - no movement
	// - Tribe 1987 (line 17) with land movement (line 18)
	expectedSections := []struct {
		keyword      string
		unitID       string
		hasMovement  bool
		movementType string // "land" or "goes_to"
	}{
		{"Tribe", "0987", true, "land"},
		{"Element", "0987e1", true, "land"},
		{"Garrison", "0987g1", false, ""},
		{"Tribe", "1987", true, "land"},
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

		hasMovement := section.UnitMovementLine != nil
		if hasMovement != exp.hasMovement {
			t.Errorf("section %d (%s): hasMovement: want %v, got %v", i, exp.unitID, exp.hasMovement, hasMovement)
		}

		if exp.hasMovement && exp.movementType == "land" {
			if _, ok := section.UnitMovementLine.(*cst.LandMovementLineNode); !ok {
				t.Errorf("section %d (%s): expected LandMovementLineNode, got %T", i, exp.unitID, section.UnitMovementLine)
			}
		}
	}

	// Expect errors for non-unit content (Scout lines, Status lines, etc.)
	if len(result.Errors()) == 0 {
		t.Error("expected errors for non-unit content")
	} else {
		t.Logf("errors (expected for Scout/Status lines): %d", len(result.Errors()))
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

func TestParse_UnitGoesToLine(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantCoords string // "Grid Number"
		wantErrors int
	}{
		{
			name:       "basic unit goes to",
			input:      "Tribe 0987, , Current Hex = QQ 1509, (Previous Hex = QQ 1410)\nCurrent Turn 900-01 (#1), Spring, FINE\nTribe Goes to QQ 1612\n",
			wantCoords: "QQ 1612",
			wantErrors: 0,
		},
		{
			name:       "unit goes to different grid",
			input:      "Tribe 1234, , Current Hex = AA 0101, (Previous Hex = N/A)\nCurrent Turn 900-01 (#1), Spring, FINE\nTribe Goes to BB 2222\n",
			wantCoords: "BB 2222",
			wantErrors: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tokens := lexers.Scan([]byte(tc.input))
			result := cst.Parse(tokens)

			if got := len(result.Errors()); got != tc.wantErrors {
				t.Errorf("errors: want %d, got %d", tc.wantErrors, got)
				for _, err := range result.Errors() {
					t.Logf("  error: %v", err)
				}
			}

			if len(result.Sections) == 0 {
				t.Fatal("expected at least one section")
			}

			section := result.Sections[0]
			if section.UnitMovementLine == nil {
				t.Fatal("UnitMovementLine is nil")
			}

			moveLine, ok := section.UnitMovementLine.(*cst.UnitGoesToLineNode)
			if !ok {
				t.Fatalf("expected UnitGoesToLineNode, got %T", section.UnitMovementLine)
			}

			if moveLine.Coords == nil {
				t.Fatal("Coords is nil")
			}

			gotCoords := string(moveLine.Coords.Grid.Bytes()) + " " + string(moveLine.Coords.Number.Bytes())
			if gotCoords != tc.wantCoords {
				t.Errorf("coords: want %q, got %q", tc.wantCoords, gotCoords)
			}
		})
	}
}

func TestParse_LandMovementLine(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantSteps  int
		wantMoves  []string // "DIR-TERRAIN" or "" for empty steps
		wantErrors int
	}{
		{
			name:       "single step",
			input:      "Tribe 0987, , Current Hex = QQ 1509, (Previous Hex = QQ 1410)\nCurrent Turn 900-01 (#1), Spring, FINE\nTribe Movement: Move N-PR,\n",
			wantSteps:  1,
			wantMoves:  []string{"N-PR"},
			wantErrors: 0,
		},
		{
			name:       "two steps",
			input:      "Tribe 0987, , Current Hex = QQ 1509, (Previous Hex = QQ 1410)\nCurrent Turn 900-01 (#1), Spring, FINE\nTribe Movement: Move N-PR, \\NE-GH,\n",
			wantSteps:  2,
			wantMoves:  []string{"N-PR", "NE-GH"},
			wantErrors: 0,
		},
		{
			name:       "empty first step",
			input:      "Tribe 0987, , Current Hex = QQ 1509, (Previous Hex = QQ 1410)\nCurrent Turn 900-01 (#1), Spring, FINE\nTribe Movement: Move \\NE-PR,\n",
			wantSteps:  2,
			wantMoves:  []string{"", "NE-PR"},
			wantErrors: 0,
		},
		{
			name:       "trailing empty step",
			input:      "Tribe 0987, , Current Hex = QQ 1509, (Previous Hex = QQ 1410)\nCurrent Turn 900-01 (#1), Spring, FINE\nTribe Movement: Move N-PR, \\\n",
			wantSteps:  2,
			wantMoves:  []string{"N-PR", ""},
			wantErrors: 0,
		},
		{
			name:       "multiple empty steps",
			input:      "Tribe 0987, , Current Hex = QQ 1509, (Previous Hex = QQ 1410)\nCurrent Turn 900-01 (#1), Spring, FINE\nTribe Movement: Move SE-PR, \\\\\n",
			wantSteps:  3,
			wantMoves:  []string{"SE-PR", "", ""},
			wantErrors: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tokens := lexers.Scan([]byte(tc.input))
			result := cst.Parse(tokens)

			if got := len(result.Errors()); got != tc.wantErrors {
				t.Errorf("errors: want %d, got %d", tc.wantErrors, got)
				for _, err := range result.Errors() {
					t.Logf("  error: %v", err)
				}
			}

			if len(result.Sections) == 0 {
				t.Fatal("expected at least one section")
			}

			section := result.Sections[0]
			if section.UnitMovementLine == nil {
				t.Fatal("UnitMovementLine is nil")
			}

			moveLine, ok := section.UnitMovementLine.(*cst.LandMovementLineNode)
			if !ok {
				t.Fatalf("expected LandMovementLineNode, got %T", section.UnitMovementLine)
			}

			if moveLine.LandMovement == nil {
				t.Fatal("LandMovement is nil")
			}

			if got := len(moveLine.LandMovement.Steps); got != tc.wantSteps {
				t.Errorf("steps: want %d, got %d", tc.wantSteps, got)
			}

			for i, wantMove := range tc.wantMoves {
				if i >= len(moveLine.LandMovement.Steps) {
					break
				}
				step := moveLine.LandMovement.Steps[i]
				gotMove := ""
				if step.Direction != nil && step.Terrain != nil {
					gotMove = string(step.Direction.Bytes()) + "-" + string(step.Terrain.Bytes())
				}
				if gotMove != wantMove {
					t.Errorf("step %d: want %q, got %q", i, wantMove, gotMove)
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
