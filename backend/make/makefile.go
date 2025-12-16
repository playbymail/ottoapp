// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package make

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var (
	// Regex to match turn report files (e.g., 0301.0899-12.0500.docx)
	reTurnReportFile = regexp.MustCompile(`^(\d{4})\.(\d{4}-\d{2})\.(0\d{3})\.docx$`)
)

// Makefile scans the data directory for turn reports and generates a Makefile to build maps.
//
//	makefileName := "maps.mk"
//	oldBehavior := false
func Makefile(makefileName string, gameId string, quiet, verbose, debug bool) error {
	err := os.Chdir(filepath.Join("files", gameId))
	if err != nil {
		return errors.Join(fmt.Errorf("invalid path"), err)
	}
	turnReportsPath, ottomapPath := "turn-reports", "ottomap"
	if sb, err := os.Stat(turnReportsPath); err != nil {
		return errors.Join(fmt.Errorf("invalid path"), err)
	} else if !sb.IsDir() {
		return fmt.Errorf("%s: not a directory", turnReportsPath)
	}
	log.Printf("%s\n", turnReportsPath)
	log.Printf("%s\n", ottomapPath)

	if filepath.Base(makefileName) != makefileName {
		return fmt.Errorf("%s: not a filename", makefileName)
	} else if filepath.Ext(makefileName) != ".mk" {
		makefileName += ".mk"
	}
	log.Printf("%s: make file\n", makefileName)

	type FileData struct {
		Id     string // file name
		Game   string
		TurnNo string
		ClanNo string
	}
	type ClanData struct {
		Id    string // clanNo
		Files map[string]*FileData
	}
	clans := make(map[string]*ClanData)
	files := make(map[string]*FileData)

	// Walk the directory structure
	entries, err := os.ReadDir(turnReportsPath)
	if err != nil {
		return fmt.Errorf("%s: %w", turnReportsPath, err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		matches := reTurnReportFile.FindStringSubmatch(entry.Name())
		if matches == nil || matches[1] != gameId {
			continue
		}
		fileData := &FileData{
			Id:     entry.Name(),
			Game:   matches[1],
			TurnNo: matches[2],
			ClanNo: matches[3],
		}
		files[fileData.Id] = fileData
		clan, ok := clans[fileData.ClanNo]
		if !ok {
			clan = &ClanData{
				Id:    fileData.ClanNo,
				Files: make(map[string]*FileData),
			}
			clans[fileData.ClanNo] = clan
		}
		clan.Files[fileData.Id] = fileData
	}

	// sort clans
	var clanNos []string
	for id := range clans {
		clanNos = append(clanNos, id)
	}
	sort.Strings(clanNos)
	if debug {
		log.Printf("makefile: clans %+v\n", clanNos)
	}

	// ../../../0301/turn-reports/
	gameDataPath := filepath.Join("..", "..", "turn-reports")
	inputPath := filepath.Join("data", "input")
	errorsPath := filepath.Join("data", "errors")
	logsPath := filepath.Join("data", "logs")
	outputPath := filepath.Join("data", "output")

	for _, clanNo := range clanNos {
		clan := clans[clanNo]

		// sort all clan files by turn
		var files []*FileData
		for _, file := range clan.Files {
			files = append(files, file)
		}
		sort.Slice(files, func(i, j int) bool {
			a, b := files[i], files[j]
			return a.Id < b.Id
		})

		// Open output file
		fmt.Printf("creating file: %q\n", filepath.Join("ottomap", clan.Id, makefileName))
		f, err := os.Create(filepath.Join("ottomap", clan.Id, makefileName))
		if err != nil {
			return fmt.Errorf("%s: %w", filepath.Join("ottomap", clan.Id, makefileName), err)
		}
		defer f.Close()

		// Write Makefile header
		_, _ = fmt.Fprintf(f, "######################################################################\n")
		_, _ = fmt.Fprintf(f, "# Generated for clan %s by ottoapp generate makefile\n", clan.Id)
		_, _ = fmt.Fprintf(f, "\n")

		// Define tool paths - using relative paths from the makefile execution context
		_, _ = fmt.Fprintf(f, "######################################################################\n")
		_, _ = fmt.Fprintf(f, "OTTOAPP  := ../../../../bin/ottoapp\n")
		_, _ = fmt.Fprintf(f, "OTTOMAP  := ../../../../bin/ottomap\n")
		_, _ = fmt.Fprintf(f, "INPUTS   := %s\n", inputPath)
		_, _ = fmt.Fprintf(f, "ERRORS   := %s\n", errorsPath)
		_, _ = fmt.Fprintf(f, "LOGS     := %s\n", logsPath)
		_, _ = fmt.Fprintf(f, "OUTPUTS  := %s\n", outputPath)
		_, _ = fmt.Fprintf(f, "\n")

		_, _ = fmt.Fprintf(f, "######################################################################\n")
		_, _ = fmt.Fprintf(f, ".PHONY: all maps\n\n")
		_, _ = fmt.Fprintf(f, "all: maps\n")
		_, _ = fmt.Fprintf(f, "maps:")
		for _, file := range files {
			_, _ = fmt.Fprintf(f, " %s/%s.%s.wxx", outputPath, file.TurnNo, clan.Id)
		}
		_, _ = fmt.Fprintf(f, "\n")

		var allReportExtracts []string
		for _, file := range files {
			// artifacts and dependencies for this turn
			turnReportFile := filepath.Join(gameDataPath, file.Id)
			rawExtractFile := fmt.Sprintf("%s/%s.%s.report.txt", outputPath, file.TurnNo, clan.Id)
			scrubbedExtractFile := fmt.Sprintf("%s/%s.%s.report.txt", inputPath, file.TurnNo, clan.Id)
			extractErrorFile := fmt.Sprintf("%s/%s.%s.extract.log", errorsPath, file.TurnNo, clan.Id)
			extractLogFile := fmt.Sprintf("%s/%s.%s.extract.log", logsPath, file.TurnNo, clan.Id)
			reportExtractFile := fmt.Sprintf("%s/%s.%s.report.txt", inputPath, file.TurnNo, clan.Id)
			mapFile := fmt.Sprintf("%s/%s.%s.wxx", outputPath, file.TurnNo, clan.Id)
			renderErrorFile := fmt.Sprintf("%s/%s.%s.render.log", errorsPath, file.TurnNo, clan.Id)
			renderLogFile := fmt.Sprintf("%s/%s.%s.render.log", logsPath, file.TurnNo, clan.Id)

			_, _ = fmt.Fprintf(f, "\n")
			_, _ = fmt.Fprintf(f, "######################################################################\n")
			_, _ = fmt.Fprintf(f, "# turn %s\n", file.TurnNo)

			// extract turn report file
			_, _ = fmt.Fprintf(f, "%s: %s\n", reportExtractFile, turnReportFile)
			_, _ = fmt.Fprintf(f, "\t@echo \"extracting $@...\"\n")
			_, _ = fmt.Fprintf(f, "\t@rm -f %s %s %s\n", rawExtractFile, extractErrorFile, extractLogFile)
			_, _ = fmt.Fprintf(f, "\t@$(OTTOAPP) run parse turn-report %s --raw-extract %s --scrubbed-extract %s 2> %s || (mv %s %s && exit 1)\n",
				turnReportFile,
				rawExtractFile,
				scrubbedExtractFile,
				extractLogFile,
				extractLogFile,
				extractErrorFile,
			)
			_, _ = fmt.Fprintf(f, "\n")
			allReportExtracts = append([]string{reportExtractFile}, allReportExtracts...)

			// create map file
			_, _ = fmt.Fprintf(f, "%s: %s\n", mapFile, strings.Join(allReportExtracts, " "))
			_, _ = fmt.Fprintf(f, "\t@echo \"rendering $@...\"\n")
			_, _ = fmt.Fprintf(f, "\t@rm -f %s %s\n", renderErrorFile, renderLogFile)
			_, _ = fmt.Fprintf(f, "\t@$(OTTOMAP) render --clan-id %s --max-turn %s --show-grid-coords --save-with-turn-id 2> %s || (mv %s %s && exit 1)\n",
				file.ClanNo,
				file.TurnNo,
				renderLogFile,
				renderLogFile,
				renderErrorFile,
			)
		}
		_, _ = fmt.Fprintf(f, "\n")
	}

	//if len(validationErrors) > 0 {
	//	for _, errMsg := range validationErrors {
	//		fmt.Fprintln(os.Stderr, errMsg)
	//	}
	//	// Clean up partial output file
	//	f.Close()
	//	os.Remove(makefileName)
	//	return fmt.Errorf("validation failed with %d errors", len(validationErrors))
	//}
	//
	//fmt.Fprintf(f, "maps: %s\n\n", strings.Join(allMaps, " "))
	//fmt.Fprintf(f, "%s", specificRules.String())
	//
	//fmt.Printf("Generated Makefile at %s with %d map targets\n", makefileName, len(allMaps))

	return nil
}
