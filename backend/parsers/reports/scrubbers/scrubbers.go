// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package scrubbers

import (
	"bytes"
	"regexp"
)

var (
	reRunOfSpaces = regexp.MustCompile(` {2,}`)
	reRunOfTabs   = regexp.MustCompile(`\t+`)
)

func Scrub(input [][]byte) [][]byte {
	var lines [][]byte
	for _, line := range input {
		line = bytes.TrimSpace(line)
		line = reRunOfTabs.ReplaceAll(line, []byte{' '})
		line = reRunOfSpaces.ReplaceAll(line, []byte{' '})
		if acceptUnitLocationLine(line) {
			// todo - this introduces a bug if the unit has movement
			lines = append(lines, rePatchNA.ReplaceAll(line, []byte(`Current Hex = $1, (Previous Hex = $1)`)))
		} else if acceptCurrentTurnLine(line) {
			lines = append(lines, line)
		} else if acceptUnitMovementLine(line) {
			lines = append(lines, line)
		} else if acceptScoutLine(line) {
			lines = append(lines, line)
		} else if acceptStatusLine(line) {
			lines = append(lines, line)
		}
	}
	return lines
}

var (
	// Courier 0987c1, , Current Hex = QQ 0408, (Previous Hex = QQ 0206)
	reCourierLocationLine = regexp.MustCompile(`^Courier [\d]{4}c[1-9],`)
	// Element 0987e1, , Current Hex = QQ 1601, (Previous Hex = QQ 1303)
	reElementLocationLine = regexp.MustCompile(`^Element [\d]{4}e[1-9],`)
	// Fleet 0987f1, , Current Hex = QQ 0818, (Previous Hex = QQ 0408)
	reFleetLocationLine = regexp.MustCompile(`^Fleet [\d]{4}f[1-9],`)
	// Garrison 0987g1, , Current Hex = QQ 1408, (Previous Hex = QQ 1408)
	reGarrisonLocationLine = regexp.MustCompile(`^Garrison [\d]{4}g[1-9],`)
	// Tribe 0987, , Current Hex = QQ 1203, (Previous Hex = QQ 1203)
	reTribeLocationLine = regexp.MustCompile(`^Tribe [\d]{4},`)

	// this will break things
	rePatchNA = regexp.MustCompile(`Current Hex = ([A-Z]{2} \d{4}),.*\(Previous Hex = N/A\)`)
)

func acceptUnitLocationLine(line []byte) bool {
	if reCourierLocationLine.Match(line) {
		return true
	} else if reElementLocationLine.Match(line) {
		return true
	} else if reFleetLocationLine.Match(line) {
		return true
	} else if reGarrisonLocationLine.Match(line) {
		return true
	}
	return reTribeLocationLine.Match(line)
}

var (
	// Current Turn 899-12 (#0), Winter, FINE Next Turn 900-01 (#1), 28/11/2025
	// Current Turn 900-01 (#1), Spring, FINE
	reCurrentTurnLine = regexp.MustCompile(`^Current Turn \d`)
)

func acceptCurrentTurnLine(line []byte) bool {
	return reCurrentTurnLine.Match(line)
}

var (
	// Tribe Movement: Move
	// Tribe Movement: Move \
	// Tribe Movement: Move failed due to Insufficient capacity to carry
	// Tribe Movement: Move N-GH, \N-JH, LJm NW,\N-JH, LJm SW, N,,River N NE\
	// Tribe Movement: Move SE-PR, \SE-CH, ,River SE S, The Settlement\\
	reTribeMovement = regexp.MustCompile(`^Tribe Movement: Move`)
)

func acceptUnitMovementLine(line []byte) bool {
	return reTribeMovement.Match(line)
}

var (
	// Scout 1:Scout ,Can't Move on Lake to NW of HEX, Nothing of interest found
	// Scout 1:Scout ,Can't Move on Ocean to N of HEX, Nothing of interest found
	// Scout 1:Scout ,No Ford on River to NE of HEX, Nothing of interest found
	// Scout 1:Scout N-GH, \N-JH, LJm NW,\,Not enough M.P's to move to N into JUNGLE HILLS, Nothing of interest found
	// Scout 1:Scout N-CH, \N-PR, , O N\,Can't Move on Ocean to N of HEX, Nothing of interest found
	// Scout 1:Scout N-GH, \N-PR, \N-PR, , O NW, N\,Can't Move on Ocean to N of HEX, Nothing of interest found
	// Scout 1:Scout N-LCM, Lcm NW,, O NE, SE, N,Find Copper Ore\,Can't Move on Ocean to N of HEX, Nothing of interest found
	// Scout 1:Scout NW-GH, , L SW, NW, N, S,Find Copper Ore\,Can't Move on Lake to NW of HEX,
	// Scout 1:Scout N-PR, , L SW,River N NE\,No Ford on River to N of HEX, Nothing of interest found
	// Scout 1:Scout NE-RH, ,River SE S\NE-RH, , O NE, SE, N,River S\,Can't Move on Ocean to NE of HEX, Nothing of interest found
	// Scout 1:Scout NE-JH, \N-PR, LJm NW,,River NW\ Nothing of interest found
	// Scout 1:Scout SW-CH, ,River SE S SW\,No Ford on River to SW of HEX, Nothing of interest found
	// Scout 1:Scout S-CH, Lcm NE,\S-PR, ,River SE,Ford S\S-PR, ,River SW NW,Ford N\,Not enough M.P's to move to S into PRAIRIE, Nothing of interest found
	// Scout 1:Scout SE-PR, \SE-PR, \SE-PR, \-, 0987, 1987\SE-PR, \SE-PR, \,Not enough M.P's to move to SE into PRAIRIE, Nothing of interest found
	// Scout 2:Scout SE-PR, \SE-CH, ,River SE S, The Settlement\,No Ford on River to SE of HEX, Nothing of interest found
	// Scout 2:Scout N-RH, \NE-BF, , O NE, Mysterious NPC\,Can't Move on Ocean to NE of HEX, Nothing of interest found
	// Scout 5:Scout S-PR, \SE-GH, \SE-JH, LJm NE,,River SE,Find Coal\,No Ford on River to SE of HEX, Nothing of interest found
	// Scout 8:Scout SE-PR, \S-GH, \S-CH, Lcm SE,\-, 1520, 0520\,Not enough M.P's to move to S into CONIFER HILLS, Patrolled and found 1987, 0987
	reScoutLine = regexp.MustCompile(`^Scout [1-8]:Scout `)
)

func acceptScoutLine(line []byte) bool {
	return reScoutLine.Match(line)
}

var (
	// 0987 Status: PRAIRIE, 0987
	// 0987 Status: PRAIRIE, L NW 0987
	// 0987 Status: PRAIRIE, O NE, SE, L SW 0987
	// 0987 Status: CONIFER HILLS, Lcm SW, 0987
	// 0987 Status: CONIFER HILLS, Lcm SE, 1987, 0987
	// 0987 Status: PRAIRIE,,River SW NW 0987
	// 0987 Status: JUNGLE HILLS, LJm NE, SE, N,,Pass NE 0987
	// 0987 Status: CONIFER HILLS, The Settlement,,River SE S 0987
	// 0987 Status: GRASSY HILLS, Copper Ore, L NE, N 0987
	// 1987 Status: PRAIRIE, O SW, L NE 0987, 1987
	// 2987 Status: GRASSY HILLS, 2987, 0987
	reCourierStatusLine  = regexp.MustCompile(`^[\d]{4}c[\d] Status:`)
	reElementStatusLine  = regexp.MustCompile(`^[\d]{4}e[\d] Status:`)
	reFleetStatusLine    = regexp.MustCompile(`^[\d]{4}f[\d] Status:`)
	reGarrisonStatusLine = regexp.MustCompile(`^[\d]{4}g[\d] Status:`)
	reTribeStatusLine    = regexp.MustCompile(`^[\d]{4} Status:`)
)

func acceptStatusLine(line []byte) bool {
	if reCourierStatusLine.Match(line) {
		return true
	} else if reElementStatusLine.Match(line) {
		return true
	} else if reFleetStatusLine.Match(line) {
		return true
	} else if reGarrisonStatusLine.Match(line) {
		return true
	}
	return reTribeStatusLine.Match(line)
}
