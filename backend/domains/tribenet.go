// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package domains

// This file defines the domain for TribeNet turn reports. It describes the
// data from a report that is needed to create a map. It does not include
// all the data that is in a report.

// TurnReport defines the data from a turn report required to generate a map.
// There are many data elements (such as season and weather) that are ignored.
type TurnReport struct {
	Turn YearMonth // all sections in a report must be for the same turn

	UnitSections []*UnitSection
}

// YearMonth is the year and month for the turn.
// It is formatted as YYYY-MM.
type YearMonth string

// UnitSection is data for a single unit in the report.
type UnitSection struct {
	UnitId UnitId // unit that owns the data in the section

	StartingLocation *Location // unit's starting location on the map for the turn
	CurrentLocation  *Location // unit's current location on the map for the turn

	GoesTo *Location // optional, the location that a unit teleports to before moving

	MovementReport *MovementReport // optional, not every unit moves every turn

	ScoutingReport []*ScoutingReport // optional, scouting results, one record per scout

	Status *Status // details on the unit's location at the end of the turn
}

// UnitId is the unique identifier for a unit.
// It matches the pattern of type followed by an optional code and sequence number.
type UnitId string

// Location is a set of coordinates on the map
// It matches N/A, ## 0102, or QQ 0102.
type Location struct {
	Hex string
}

// MovementReport can be for land units or fleets
type MovementReport struct {
	UnitId   UnitId    // unit making the report
	Location *Location // the starting location for the steps, usually the unit's current location
	Steps    []*Step   // every "step" in the report
}

type Step struct {
	Location *Location // the location the step is reporting on
	Results  []*Result // the results of the step, usually information on the "to" location and visible neighbors
}

type Result struct {
	Observation Observation          // how the location was observed
	Terrain     Terrain              // the terrain code for the location
	Resources   []Resource           // a list of resources in the location
	Borders     map[Direction][]Edge // optional information about the borders of the location
	Things      []*Thing             // optional list of things found in the location
}

// Observation is an enum describing the type of observation
// (in the location, the location is a neighbor, the location is in the outer ring)
type Observation int

const (
	ScoutInLocation    Observation = iota + 1 // scout, direct, full report
	UnitInLocation                            // non-scout, direct, brief report
	UnitNextToLocation                        // any unit, indirect, next to
	FleetNextToRing                           // fleet reported on location two hexes away
)

// Thing describes something that was found in a location
type Thing struct {
	SpecialHex  *SpecialHex // the thing is a special hex
	Settlements *Settlement // the thing is a settlement
	Units       UnitId      // the thing is a unit in the hex
}

// SpecialHex describes a special hex
type SpecialHex struct {
	Label string // name of the special hex
}

// Settlement describes a settlement (or village) in a hex
type Settlement struct {
	Label string // name of the settlement
}

// ScoutingReport is the results of a scouting expedition
type ScoutingReport struct {
	ScoutId  ScoutId   // scout making the report
	Location *Location // the starting location for the steps, usually the unit's current location
	Steps    []*Step   // every "step" in the report
}

// ScoutId is the unique identifier for a scout.
// It matches a UnitId followed by 's' and a sequence number from 1..8.
type ScoutId string

// Status provides data for a unit's location at the end of the turn.
type Status struct {
	UnitId   UnitId    // unit reporting the status
	Location *Location // the starting location for the results, usually the unit's current location
	Results  []*Result // results for each "step" in the status line
}
