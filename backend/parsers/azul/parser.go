// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package azul implements a parser that splits a turn report into
// sections. Sections contain only the lines needed to create maps.
package azul

import (
	"bytes"
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/playbymail/ottoapp/backend/parsers/azul/follows"
	"github.com/playbymail/ottoapp/backend/parsers/azul/goes"
	"github.com/playbymail/ottoapp/backend/parsers/azul/location"
	"github.com/playbymail/ottoapp/backend/parsers/azul/scout"
	"github.com/playbymail/ottoapp/backend/parsers/azul/turn"
	"github.com/playbymail/ottoapp/backend/parsers/docx"
)

// Parse parses a turn report file (the original Word document) and returns
// the sections for every unit in it.
func Parse(source string, input []byte, quiet, verbose, debug bool) (Report, error) {
	const trimLeading, trimTrailing = false, false

	d, err := docx.ParseReader(bytes.NewReader(input), trimLeading, trimTrailing, quiet, verbose, debug)
	if err != nil {
		if debug {
			log.Printf("azul: docx.Parse %v\n", err)
		}
		return Report{}, fmt.Errorf("azul: parse: %w", err)
	} else if d == nil {
		if debug {
			log.Printf("azul: docx.Parse returned nil\n")
		}
		return Report{}, fmt.Errorf("azul: parse: %w", err)
	}

	r := Report{
		Path: filepath.Dir(source),
		Name: filepath.Base(source),
	}

	var section *Section

	for _, line := range bytes.Split(d.Text, []byte{LF}) {
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		if reClanSection.Match(line) {
			// found a unit location line, so close out the prior section
			section = nil

			if debug {
				log.Printf("input %q\n", string(line))
			}

			l, err := location.Expect(r.Name, line)
			if err != nil {
				if pe := turn.ExtractParseError(err); pe == nil {
					log.Printf("%s: section: location %v", r.Name, err)
				} else {
					log.Printf("%s: section: location: parsing error\n%s\n%s^^^\n%v\n",
						r.Name, string(line), strings.Repeat(" ", pe.Pos.Offset), err)
				}
				continue
			}
			section = &Section{
				UnitId:         l.UnitId,
				Kind:           "clan",
				PreviousCoords: l.PreviousCoords,
				CurrentCoords:  l.CurrentCoords,
			}
			r.Sections = append(r.Sections, section)
			continue
		}

		if reCourierSection.Match(line) {
			// found a unit location line, so close out the prior section
			section = nil

			if debug {
				log.Printf("input %q\n", string(line))
			}

			l, err := location.Expect(r.Name, line)
			if err != nil {
				if pe := turn.ExtractParseError(err); pe == nil {
					log.Printf("%s: section: location %v", r.Name, err)
				} else {
					log.Printf("%s: section: location: parsing error\n%s\n%s^^^\n%v\n",
						r.Name, string(line), strings.Repeat(" ", pe.Pos.Offset), err)
				}
				continue
			}
			section = &Section{
				UnitId:         l.UnitId,
				Kind:           "courier",
				PreviousCoords: l.PreviousCoords,
				CurrentCoords:  l.CurrentCoords,
			}
			r.Sections = append(r.Sections, section)
			continue
		}

		if reElementSection.Match(line) {
			// found a unit location line, so close out the prior section
			section = nil

			if debug {
				log.Printf("input %q\n", string(line))
			}

			l, err := location.Expect(r.Name, line)
			if err != nil {
				if pe := turn.ExtractParseError(err); pe == nil {
					log.Printf("%s: section: location %v", r.Name, err)
				} else {
					log.Printf("%s: section: location: parsing error\n%s\n%s^^^\n%v\n",
						r.Name, string(line), strings.Repeat(" ", pe.Pos.Offset), err)
				}
				continue
			}
			section = &Section{
				UnitId:         l.UnitId,
				Kind:           "element",
				PreviousCoords: l.PreviousCoords,
				CurrentCoords:  l.CurrentCoords,
			}
			r.Sections = append(r.Sections, section)
			continue
		}

		if reFleetSection.Match(line) {
			// found a unit location line, so close out the prior section
			section = nil

			if debug {
				log.Printf("input %q\n", string(line))
			}

			l, err := location.Expect(r.Name, line)
			if err != nil {
				if pe := turn.ExtractParseError(err); pe == nil {
					log.Printf("%s: section: location %v", r.Name, err)
				} else {
					log.Printf("%s: section: location: parsing error\n%s\n%s^^^\n%v\n",
						r.Name, string(line), strings.Repeat(" ", pe.Pos.Offset), err)
				}
				continue
			}
			log.Printf("%s: location %+v\n", r.Name, l)
			section = &Section{
				UnitId:         l.UnitId,
				Kind:           "fleet",
				PreviousCoords: l.PreviousCoords,
				CurrentCoords:  l.CurrentCoords,
			}
			r.Sections = append(r.Sections, section)
			continue
		}

		if reGarrisonSection.Match(line) {
			// found a unit location line, so close out the prior section
			section = nil

			if debug {
				log.Printf("input %q\n", string(line))
			}

			l, err := location.Expect(r.Name, line)
			if err != nil {
				if pe := turn.ExtractParseError(err); pe == nil {
					log.Printf("%s: section: location %v", r.Name, err)
				} else {
					log.Printf("%s: section: location: parsing error\n%s\n%s^^^\n%v\n",
						r.Name, string(line), strings.Repeat(" ", pe.Pos.Offset), err)
				}
				continue
			}
			section = &Section{
				UnitId:         l.UnitId,
				Kind:           "clan",
				PreviousCoords: l.PreviousCoords,
				CurrentCoords:  l.CurrentCoords,
			}
			r.Sections = append(r.Sections, section)
			continue
		}

		if reTribeSection.Match(line) {
			// found a unit location line, so close out the prior section
			section = nil

			if debug {
				log.Printf("input %q\n", string(line))
			}

			l, err := location.Expect(r.Name, line)
			if err != nil {
				if pe := turn.ExtractParseError(err); pe == nil {
					log.Printf("%s: section: location %v", r.Name, err)
				} else {
					log.Printf("%s: section: location: parsing error\n%s\n%s^^^\n%v\n",
						r.Name, string(line), strings.Repeat(" ", pe.Pos.Offset), err)
				}
				continue
			}
			section = &Section{
				UnitId:         l.UnitId,
				Kind:           "tribe",
				PreviousCoords: l.PreviousCoords,
				CurrentCoords:  l.CurrentCoords,
			}
			r.Sections = append(r.Sections, section)
			continue
		}

		if section == nil {
			continue
		}

		if reCurrentTurn.Match(line) {
			if debug {
				log.Printf("input %q\n", string(line))
			}
			t, err := turn.Expect(r.Name, line)
			if err != nil {
				if pe := turn.ExtractParseError(err); pe == nil {
					log.Printf("%s: section: turn %v", r.Name, err)
				} else {
					log.Printf("%s: section: turn: parsing error\n%s\n%s^^^\n%v\n",
						r.Name, string(line), strings.Repeat(" ", pe.Pos.Offset), err)
				}
				continue
			}
			section.TurnNo = t.TurnNo
			continue
		}

		if reClanScry.Match(line) {
			if debug {
				log.Printf("input %q\n", string(line))
			}
			panic("!implemented")
		} else if reCourierScry.Match(line) {
			if debug {
				log.Printf("input %q\n", string(line))
			}
			panic("!implemented")
		} else if reElementScry.Match(line) {
			if debug {
				log.Printf("input %q\n", string(line))
			}
			panic("!implemented")
		} else if reFleetScry.Match(line) {
			if debug {
				log.Printf("input %q\n", string(line))
			}
			panic("!implemented")
		} else if reGarrisonScry.Match(line) {
			if debug {
				log.Printf("input %q\n", string(line))
			}
			panic("!implemented")
		} else if reTribeScry.Match(line) {
			if debug {
				log.Printf("input %q\n", string(line))
			}
			panic("!implemented")
		}

		if reFleetMovement.Match(line) {
			if debug {
				log.Printf("input %q\n", string(line))
			}
			panic("!implemented")
		} else if reTribeFollows.Match(line) {
			if debug {
				log.Printf("input %q\n", string(line))
			}
			f, err := follows.Expect(r.Name, line)
			if err != nil {
				if pe := follows.ExtractParseError(err); pe == nil {
					log.Printf("%s: section: follows %v", r.Name, err)
				} else {
					log.Printf("%s: section: follows: parsing error\n%s\n%s^^^\n%v\n",
						r.Name, string(line), strings.Repeat(" ", pe.Pos.Offset), err)
				}
				continue
			}
			section.Follows = &Follows{
				UnitId: f.UnitId,
			}
			continue
		} else if reTribeGoesTo.Match(line) {
			if debug {
				log.Printf("input %q\n", string(line))
			}
			g, err := goes.Expect(r.Name, line)
			if err != nil {
				if pe := goes.ExtractParseError(err); pe == nil {
					log.Printf("%s: section: goes %v", r.Name, err)
				} else {
					log.Printf("%s: section: goes: parsing error\n%s\n%s^^^\n%v\n",
						r.Name, string(line), strings.Repeat(" ", pe.Pos.Offset), err)
				}
				continue
			}
			section.GoesTo = &GoesTo{
				Coords: g.Coords,
			}
			continue
		} else if reTribeMovement.Match(line) {
			if debug {
				log.Printf("input %q\n", string(line))
			}
			panic("!implemented")
		}

		if reScout.Match(line) {
			if debug {
				log.Printf("input %q\n", string(line))
			}
			s, err := scout.Expect(r.Name, line)
			if err != nil {
				if pe := scout.ExtractParseError(err); pe == nil {
					log.Printf("%s: section: scout %v", r.Name, err)
				} else {
					log.Printf("%s: section: scout: parsing error\n%s\n%s^^^\n%v\n",
						r.Name, string(line), strings.Repeat(" ", pe.Pos.Offset), err)
				}
				continue
			}
			log.Printf("scout %+v\n", s)
			continue
		}

		//if reClanStatus.Match(line) {
		//	if debug {
		//		log.Printf("input %q\n", string(line))
		//	}
		//	s, err := status.Expect(r.Name, line)
		//	if err != nil {
		//		if pe := status.ExtractParseError(err); pe == nil {
		//			log.Printf("%s: section: status %v", r.Name, err)
		//		} else {
		//			log.Printf("%s: section: status: parsing error\n%s\n%s^^^\n%v\n",
		//				r.Name, string(line), strings.Repeat(" ", pe.Pos.Offset), err)
		//		}
		//		continue
		//	}
		//	log.Printf("status %+v\n", s)
		//	section.Status = &Status{}
		//	section = nil
		//	continue
		//} else if reCourierStatus.Match(line) {
		//	if debug {
		//		log.Printf("input %q\n", string(line))
		//	}
		//	s, err := status.Expect(r.Name, line)
		//	if err != nil {
		//		if pe := status.ExtractParseError(err); pe == nil {
		//			log.Printf("%s: section: status %v", r.Name, err)
		//		} else {
		//			log.Printf("%s: section: status: parsing error\n%s\n%s^^^\n%v\n",
		//				r.Name, string(line), strings.Repeat(" ", pe.Pos.Offset), err)
		//		}
		//		continue
		//	}
		//	log.Printf("status %+v\n", s)
		//	section.Status = &Status{}
		//	section = nil
		//	continue
		//} else if reElementStatus.Match(line) {
		//	if debug {
		//		log.Printf("input %q\n", string(line))
		//	}
		//	s, err := status.Expect(r.Name, line)
		//	if err != nil {
		//		if pe := scout.ExtractParseError(err); pe == nil {
		//			log.Printf("%s: section: status %v", r.Name, err)
		//		} else {
		//			log.Printf("%s: section: status: parsing error\n%s\n%s^^^\n%v\n",
		//				r.Name, string(line), strings.Repeat(" ", pe.Pos.Offset), err)
		//		}
		//		continue
		//	}
		//	log.Printf("status %+v\n", s)
		//	section.Status = &Status{}
		//	section = nil
		//	continue
		//} else if reFleetStatus.Match(line) {
		//	if debug {
		//		log.Printf("input %q\n", string(line))
		//	}
		//	s, err := status.Expect(r.Name, line)
		//	if err != nil {
		//		if pe := scout.ExtractParseError(err); pe == nil {
		//			log.Printf("%s: section: status %v", r.Name, err)
		//		} else {
		//			log.Printf("%s: section: status: parsing error\n%s\n%s^^^\n%v\n",
		//				r.Name, string(line), strings.Repeat(" ", pe.Pos.Offset), err)
		//		}
		//		continue
		//	}
		//	log.Printf("status %+v\n", s)
		//	section.Status = &Status{}
		//	section = nil
		//	continue
		//} else if reGarrisonStatus.Match(line) {
		//	if debug {
		//		log.Printf("input %q\n", string(line))
		//	}
		//	s, err := status.Expect(r.Name, line)
		//	if err != nil {
		//		if pe := scout.ExtractParseError(err); pe == nil {
		//			log.Printf("%s: section: status %v", r.Name, err)
		//		} else {
		//			log.Printf("%s: section: status: parsing error\n%s\n%s^^^\n%v\n",
		//				r.Name, string(line), strings.Repeat(" ", pe.Pos.Offset), err)
		//		}
		//		continue
		//	}
		//	log.Printf("status %+v\n", s)
		//	section.Status = &Status{}
		//	section = nil
		//	continue
		//} else if reTribeStatus.Match(line) {
		//	if debug {
		//		log.Printf("input %q\n", string(line))
		//	}
		//	s, err := status.Expect(r.Name, line)
		//	if err != nil {
		//		if pe := scout.ExtractParseError(err); pe == nil {
		//			log.Printf("%s: section: status %v", r.Name, err)
		//		} else {
		//			log.Printf("%s: section: status: parsing error\n%s\n%s^^^\n%v\n",
		//				r.Name, string(line), strings.Repeat(" ", pe.Pos.Offset), err)
		//		}
		//		continue
		//	}
		//	log.Printf("status %+v\n", s)
		//	section.Status = &Status{}
		//	section = nil
		//	continue
		//}
	}

	if len(r.Sections) == 0 {
		return Report{}, fmt.Errorf("azul: parse: no sections")
	}

	for _, section := range r.Sections {
		log.Printf("%s: from %-9q to %-9q\n", r.Name, section.PreviousCoords, section.CurrentCoords)
		if r.TurnNo == "" {
			r.TurnNo = section.TurnNo
		}
		//if r.TurnNo != section.TurnNo {
		//	return Report{}, fmt.Errorf("invalid report: multiple turns")
		//}
	}

	if r.TurnNo == "" {
		return Report{}, fmt.Errorf("azul: parse: no turn info")
	}

	return r, nil
}

const (
	CR = 0x0d // carriage return
	LF = 0x0a // line feed
)

var (
	// Tribe 0987, , Current Hex = QQ 0909, (Previous Hex = QQ 1010)
	reClanSection     = regexp.MustCompile(`^Tribe\s(0\d{3}),`)
	reCourierSection  = regexp.MustCompile(`^Courier\s(\d{4}c[1-9]),`)
	reElementSection  = regexp.MustCompile(`^Element\s(\d{4}e[1-9]),`)
	reFleetSection    = regexp.MustCompile(`^Fleet\s(\d{4}f[1-9]),`)
	reGarrisonSection = regexp.MustCompile(`^Garrison\s(\d{4}g[1-9]),`)
	reTribeSection    = regexp.MustCompile(`^Tribe\s(\d{4}),`)

	// Current Turn 899-12 (#0), Winter, FINE Split Turn 900-01 (#1), 28/11/2025
	reCurrentTurn = regexp.MustCompile(`^Current\sTurn\s(\d{3,4}-\d{2})\s\(#\d+\),`)

	// 0987 Scry: QQ 1010:
	reClanScry     = regexp.MustCompile(`^0\d{3}\sScry:`)
	reCourierScry  = regexp.MustCompile(`^\d{4}c[1-9]\sScry:`)
	reElementScry  = regexp.MustCompile(`^\d{4}e[1-9]\sScry:`)
	reFleetScry    = regexp.MustCompile(`^\d{4}f[1-9]\sScry:`)
	reGarrisonScry = regexp.MustCompile(`^\d{4}g[1-9]\sScry:`)
	reTribeScry    = regexp.MustCompile(`^\d{4}\sScry:`)

	// CALM NE Fleet Movement:
	reFleetMovement = regexp.MustCompile(`^(CALM|MILD|STRONG|GALE)\s[NS][EW]?\sFleet\sMovement:`)

	// Tribe Follows 0987e1
	reTribeFollows = regexp.MustCompile(`^Tribe Follows\s`)

	// Tribe Goes to QQ 1010
	reTribeGoesTo = regexp.MustCompile(`^Tribe Goes to\s`)

	// Tribe Movement: Move
	reTribeMovement = regexp.MustCompile(`^Tribe Movement:`)

	// Scout 1:Scout
	reScout = regexp.MustCompile(`^Scout\s\d:Scout`)

	// 0987 Status:
	reClanStatus     = regexp.MustCompile(`^0\d{3}\sStatus:`)
	reCourierStatus  = regexp.MustCompile(`^\d{4}c[1-9]\sStatus:`)
	reElementStatus  = regexp.MustCompile(`^\d{4}e[1-9]\sStatus:`)
	reFleetStatus    = regexp.MustCompile(`^\d{4}f[1-9]\sStatus:`)
	reGarrisonStatus = regexp.MustCompile(`^\d{4}g[1-9]\sStatus:`)
	reTribeStatus    = regexp.MustCompile(`^\d{4}\sStatus:`)
)

type Report struct {
	Path     string
	Name     string
	TurnNo   string
	Sections []*Section
}

type Section struct {
	UnitId         string
	Kind           string
	TurnNo         string // yyyy-mm
	PreviousCoords string
	CurrentCoords  string
	Follows        *Follows
	GoesTo         *GoesTo
	LandMovement   *LandMovement
	WaterMovement  *WaterMovement
	ScoutMovement  *ScoutMovement
	Status         *Status
}

type Follows struct {
	UnitId string
}

type GoesTo struct {
	Coords string
}

type LandMovement struct {
	Line []byte
}

type WaterMovement struct {
	Wind      string
	Direction string
	Line      []byte
}

type ScoutMovement struct {
	Line []byte
	No   int
}

type Status struct {
	Line []byte
}
