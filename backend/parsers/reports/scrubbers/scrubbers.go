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
		} else if acceptTribeMovementLine(line) {
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
	reCurrentTurnLine = regexp.MustCompile(`^Current Turn .*Next Turn`)
)

func acceptCurrentTurnLine(line []byte) bool {
	return reCurrentTurnLine.Match(line)
}

var (
	reScoutLine = regexp.MustCompile(`^Scout [1-8]:Scout `)
)

func acceptScoutLine(line []byte) bool {
	return reScoutLine.Match(line)
}

var (
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

func acceptTribeMovementLine(line []byte) bool {
	return bytes.HasPrefix(line, []byte(`Tribe Movement:`))
}

var (
	reCourierLocationLine  = regexp.MustCompile(`^Courier [\d]{4}c[1-9],`)
	reElementLocationLine  = regexp.MustCompile(`^Element [\d]{4}e[1-9],`)
	reFleetLocationLine    = regexp.MustCompile(`^Fleet [\d]{4}f[1-9],`)
	reGarrisonLocationLine = regexp.MustCompile(`^Garrison [\d]{4}g[1-9],`)
	reTribeLocationLine    = regexp.MustCompile(`^Tribe [\d]{4},`)

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
