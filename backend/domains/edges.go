// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package domains

import (
	"encoding/json"
	"fmt"
)

// Edge is an enum for the edge of a hex
type Edge int

const (
	Open Edge = iota
	Canal
	Ford
	Pass
	River
	StonyRoad
)

// MarshalJSON implements the json.Marshaler interface.
func (e Edge) MarshalJSON() ([]byte, error) {
	return json.Marshal(EdgeToString[e])
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (e *Edge) UnmarshalJSON(data []byte) error {
	var s string
	var ok bool
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	} else if *e, ok = StringToEdge[s]; !ok {
		return fmt.Errorf("invalid Edge %q", s)
	}
	return nil
}

// String implements the fmt.Stringer interface.
func (e Edge) String() string {
	if str, ok := EdgeToString[e]; ok {
		return str
	}
	return fmt.Sprintf("Edge(%d)", int(e))
}

var (
	// EdgeToString is a helper map for marshalling the enum
	EdgeToString = map[Edge]string{
		Open:      "Open",
		Canal:     "Canal",
		Ford:      "Ford",
		Pass:      "Pass",
		River:     "River",
		StonyRoad: "StonyRoad",
	}
	// StringToEdge is a helper map for unmarshalling the enum
	StringToEdge = map[string]Edge{
		"Open":      Open,
		"Canal":     Canal,
		"Ford":      Ford,
		"Pass":      Pass,
		"River":     River,
		"StonyRoad": StonyRoad,
	}
)
