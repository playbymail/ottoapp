// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package dag

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/playbymail/ottoapp/backend/services/documents"
)

// Node is a file in the dependency graph
type Node struct {
	// ID is the hash of the file contents
	ID string

	// Path is the relative path to the file
	Path string

	// Game is the game from the file name
	Game string

	// TurnNo is yyyy-mm from the file name
	TurnNo string

	// Clan is from the file name
	Clan int

	// Type is turn-report, report-extract, map
	Type string

	// ModifiedAt is the ModTime of the file
	ModifiedAt time.Time

	// Outgoing edges: this node depends on these nodes.
	DependsOn []string
}

type Graph struct {
	Nodes map[string]*Node
}

var (
	// turn report file name is {game}.{turnNo}.{unitId}.docx
	reTurnReportFile = regexp.MustCompile(`^(\d{4})\.(\d{4}-\d{2})\.(\d{4}([cefg][1-9])?)\.docx$`)

	// report extract file name is {game}.{turnNo}.{unitId}.scrubbed.txt
	reReportExtractFile = regexp.MustCompile(`^(\d{4})\.(\d{4}-\d{2})\.(\d{4}([cefg][1-9])?)\.scrubbed\.txt$`)
)

// ReadReportExtractFiles walks the report extract files path and creates nodes
// for every report extract we find.
func (g *Graph) ReadReportExtractFiles(gameId, path string, quiet, verbose, debug bool) error {
	if verbose {
		log.Printf("dag: %s: walk %q\n", gameId, path)
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		if debug {
			log.Printf("dag: %s: walk %v\n", gameId, err)
		}
		return err
	}

	for _, entry := range entries {
		if debug {
			log.Printf("dag: %s: walk %q\n", gameId, entry.Name())
		}
		if entry.IsDir() {
			if debug {
				log.Printf("dag: %s: walk %q: is dir\n", gameId, entry.Name())
			}
			continue
		}
		matches := reReportExtractFile.FindStringSubmatch(entry.Name())
		if matches == nil { // not a report extract file
			continue
		} else if matches[1] != gameId { // file from different game
			if debug {
				log.Printf("dag: %s: walk %q: game %q\n", gameId, entry.Name(), matches[1])
			}
			continue
		}
		fi, err := entry.Info()
		if err != nil {
			return err
		} else if !fi.Mode().IsRegular() {
			continue
		}

		node := &Node{
			Path:       filepath.Join(path, entry.Name()),
			Game:       gameId,
			TurnNo:     matches[2],
			Type:       "report-extract",
			ModifiedAt: fi.ModTime(),
		}

		// extract clan from unit id in file name
		node.Clan, _ = strconv.Atoi(matches[3][:4])
		if node.Clan > 999 {
			if debug {
				log.Printf("dag: %s: walk: %s: %q\n", gameId, entry.Name(), matches[3])
			}
			node.Clan = node.Clan % 1000
		}

		if data, err := os.ReadFile(node.Path); err != nil {
			return err
		} else if _, node.ID, err = documents.Hash(data); err != nil {
			return err
		}

		g.Nodes[node.ID] = node
	}

	return nil
}

// ReadTurnReportFiles walks the turn report files path and creates nodes
// for every turn report we find.
func (g *Graph) ReadTurnReportFiles(gameId, path string, quiet, verbose, debug bool) error {
	if verbose {
		log.Printf("dag: %s: walk %q\n", gameId, path)
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		if debug {
			log.Printf("dag: %s: walk %v\n", gameId, err)
		}
		return err
	}

	for _, entry := range entries {
		if debug {
			log.Printf("dag: %s: walk %q\n", gameId, entry.Name())
		}
		if entry.IsDir() {
			if debug {
				log.Printf("dag: %s: walk %q: is dir\n", gameId, entry.Name())
			}
			continue
		}
		matches := reTurnReportFile.FindStringSubmatch(entry.Name())
		if matches == nil { // not a turn report  file
			continue
		} else if matches[1] != gameId { // file from different game
			if debug {
				log.Printf("dag: %s: walk %q: game %q\n", gameId, entry.Name(), matches[1])
			}
			continue
		}
		fi, err := entry.Info()
		if err != nil {
			return err
		} else if !fi.Mode().IsRegular() {
			continue
		}

		node := &Node{
			Path:       filepath.Join(path, entry.Name()),
			Game:       gameId,
			TurnNo:     matches[2],
			Type:       "turn-report",
			ModifiedAt: fi.ModTime(),
		}

		// extract clan from unit id in file name
		node.Clan, _ = strconv.Atoi(matches[3][:4])
		if node.Clan > 999 {
			if debug {
				log.Printf("dag: %s: walk: %s: %q\n", gameId, entry.Name(), matches[3])
			}
			node.Clan = node.Clan % 1000
		}

		if data, err := os.ReadFile(node.Path); err != nil {
			return err
		} else if _, node.ID, err = documents.Hash(data); err != nil {
			return err
		}

		g.Nodes[node.ID] = node
	}

	return nil
}
