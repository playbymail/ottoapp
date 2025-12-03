// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package domains

import (
	"encoding/json"
	"fmt"
)

// Direction is an enum for the direction
type Direction int

const (
	North Direction = iota
	NorthEast
	SouthEast
	South
	SouthWest
	NorthWest
)

// Directions is a helper for iterating over the directions
var Directions = []Direction{
	North,
	NorthEast,
	SouthEast,
	South,
	SouthWest,
	NorthWest,
}

// MarshalJSON implements the json.Marshaler interface.
func (d Direction) MarshalJSON() ([]byte, error) {
	return json.Marshal(DirectionToString[d])
}

// MarshalText implements the encoding.TextMarshaler interface.
// This is needed for marshalling the enum as map keys.
//
// Note that this is called by the json package, unlike the UnmarshalText function.
func (d Direction) MarshalText() (text []byte, err error) {
	return []byte(DirectionToString[d]), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (d *Direction) UnmarshalJSON(data []byte) error {
	var s string
	var ok bool
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	} else if *d, ok = StringToDirection[s]; !ok {
		return fmt.Errorf("invalid Direction %q", s)
	}
	return nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// This is needed for unmarshalling the enum as map keys.
//
// Note that this is never called; it just changes the code path in UnmarshalJSON.
func (d Direction) UnmarshalText(text []byte) error {
	panic("!")
}

// String implements the fmt.Stringer interface.
func (d Direction) String() string {
	if str, ok := DirectionToString[d]; ok {
		return str
	}
	return fmt.Sprintf("Direction(%d)", int(d))
}

var (
	// DirectionToString is a helper map for marshalling the enum
	DirectionToString = map[Direction]string{
		North:     "N",
		NorthEast: "NE",
		SouthEast: "SE",
		South:     "S",
		SouthWest: "SW",
		NorthWest: "NW",
	}

	// StringToDirection is a helper map for unmarshalling the enum
	StringToDirection = map[string]Direction{
		"N":  North,
		"NE": NorthEast,
		"SE": SouthEast,
		"S":  South,
		"SW": SouthWest,
		"NW": NorthWest,
	}
)

// column direction vectors defines the vectors used to determine the coordinates
// of the neighboring column based on the direction and the odd/even column
// property of the starting hex.
//
// NB: grids and hexes in TribeNet start at 1, 1 so "odd" and "even" are based on the hex coordinates.
//
// NB: indexed by [odd/even][direction]
var columnDirectionVectors = [2][6]struct {
	column int
	row    int
}{
	{ // even columns
		North:     {column: +0, row: -1}, // ## 1206 -> ## 1205
		NorthEast: {column: +1, row: +0}, // ## 1206 -> ## 1306
		SouthEast: {column: +1, row: +1}, // ## 1206 -> ## 1307
		South:     {column: +0, row: +1}, // ## 1206 -> ## 1207
		SouthWest: {column: -1, row: +1}, // ## 1206 -> ## 1107
		NorthWest: {column: -1, row: +0}, // ## 1206 -> ## 1106
	},
	{ // odd columns
		North:     {column: +0, row: -1}, // ## 1306 -> ## 1305
		NorthEast: {column: +1, row: -1}, // ## 1306 -> ## 1405
		SouthEast: {column: +1, row: +0}, // ## 1306 -> ## 1406
		South:     {column: +0, row: +1}, // ## 1306 -> ## 1307
		SouthWest: {column: -1, row: +0}, // ## 1306 -> ## 1206
		NorthWest: {column: -1, row: -1}, // ## 1306 -> ## 1205
	},
}

// AddDirection moves in the given direction and returns the new row and column.
// It always moves a single hex and allows for moving between grids and wrapping around the big map.
// Will panic on a bad direction.
func AddDirection(row, col int, d Direction) (newRow, newColumn int) {
	return row + columnDirectionVectors[col%2][d].row, col + columnDirectionVectors[col%2][d].column
}
