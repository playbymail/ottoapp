// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package make

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/playbymail/ottoapp/backend/services/documents"
	"github.com/playbymail/ottoapp/backend/stores/sqlite"
)

// ExportExtractFiles creates missing folders and files.
func ExportExtractFiles(db *sqlite.DB, path string, quiet, verbose, debug bool) error {
	if verbose {
		log.Printf("export: path %q\n", path)
	}

	var ottomapTemplate struct {
		Clan        string `json:"Clan,omitempty"`
		AllowConfig bool   `json:"AllowConfig,omitempty"`
		DebugFlags  struct {
			LogFile bool `json:"LogFile,omitempty"`
		} `json:"DebugFlags"`
		Worldographer struct {
			Map struct {
				Zoom   int `json:"Zoom,omitempty"`
				Layers struct {
					LargeCoords bool `json:"LargeCoords,omitempty"`
					MPCost      bool `json:"MPCost,omitempty"`
				} `json:"Layers"`
				Units struct {
					Default  string `json:"Default,omitempty"`
					Clan     string `json:"Clan,omitempty"`
					Courier  string `json:"Courier,omitempty"`
					Fleet    string `json:"Fleet,omitempty"`
					Garrison string `json:"Garrison,omitempty"`
					Multiple string `json:"Multiple,omitempty"`
					Tribe    string `json:"Tribe,omitempty"`
				} `json:"Units"`
				Terrain struct {
					Lake  string `json:"Lake,omitempty"`
					Ocean string `json:"Ocean,omitempty"`
				} `json:"Terrain"`
			} `json:"Map"`
		} `json:"Worldographer"`
	}
	ottomapTemplate.AllowConfig = true
	ottomapTemplate.DebugFlags.LogFile = true
	ottomapTemplate.Worldographer.Map.Zoom = 1
	ottomapTemplate.Worldographer.Map.Layers.LargeCoords = true
	ottomapTemplate.Worldographer.Map.Layers.MPCost = true
	ottomapTemplate.Worldographer.Map.Units.Default = "Military Ancient Soldier"
	ottomapTemplate.Worldographer.Map.Units.Clan = "Military Ancient Soldier"
	ottomapTemplate.Worldographer.Map.Units.Courier = "Military Knight"
	ottomapTemplate.Worldographer.Map.Units.Fleet = "Military Sailship"
	ottomapTemplate.Worldographer.Map.Units.Garrison = "Military Camp"
	ottomapTemplate.Worldographer.Map.Units.Multiple = "Military Ancient Soldier"
	ottomapTemplate.Worldographer.Map.Units.Tribe = "Military Ancient Soldier"
	ottomapTemplate.Worldographer.Map.Terrain.Lake = "Water Sea"
	ottomapTemplate.Worldographer.Map.Terrain.Ocean = "Water Ocean"

	// initialize services needed to export the report extract files
	documentsSvc, err := documents.New(db, nil, nil)
	if err != nil {
		if debug {
			log.Printf("make: export: extract-files: documents: new %v\n", err)
		}
		return fmt.Errorf("export: %w", err)
	}

	reportExtractFiles, err := documentsSvc.ReadReportExtractMeta()
	if err != nil {
		return err
	}

	for _, reportExtractFile := range reportExtractFiles {
		gameId := string(reportExtractFile.GameID)
		// turnNo := string(reportExtractFile.TurnNo)
		clanNo := fmt.Sprintf("%04d", reportExtractFile.ClanNo)
		// log.Printf("export: %s %s %s\n", gameId, turnNo, clanNo)

		gameFolder := filepath.Join(path, "files", gameId, "ottomap")
		clanFolder := filepath.Join(gameFolder, clanNo)
		clanDataFolder := filepath.Join(clanFolder, "data")

		// create any missing clan data sub-folders
		for _, folder := range []string{
			"errors", "input", "logs", "output",
		} {
			dataFolder := filepath.Join(clanDataFolder, folder)
			if isdir(dataFolder) {
				continue
			}
			err = os.MkdirAll(dataFolder, 0o755)
			if err != nil {
				return err
			}
			log.Printf("export: created %q\n", dataFolder)
		}

		// create any missing ottomap configuration files
		ottomapJsonFile := filepath.Join(clanDataFolder, "input", "ottomap.json")
		if !isfile(ottomapJsonFile) {
			ottomapData := ottomapTemplate
			ottomapData.Clan = clanNo
			data, err := json.MarshalIndent(ottomapData, "", "  ")
			if err != nil {
				return err
			}
			err = os.WriteFile(ottomapJsonFile, data, 0o644)
			if err != nil {
				return err
			}
			log.Printf("export: created %q\n", ottomapJsonFile)
		}

		// export all extract files (files/{gameId}/ottomap/{clanNo}/data/input/{extractFileName}

		// fetch the contents of the report extract file
		reportExtractFile.Contents, err = documentsSvc.ReadReportExtractContents(reportExtractFile.ID)
		if err != nil {
			log.Printf("export: %s: %v\n", reportExtractFile.Path)
		}
		_, contentsHash, err := documents.Hash(reportExtractFile.Contents)
		if err != nil {
			return err
		}

		extractFileName := filepath.Join(clanDataFolder, "input", reportExtractFile.Path)

		// export only if the file is missing or the contents don't match the database
		mustExport := false
		if !isfile(extractFileName) {
			mustExport = true
		} else if data, err := os.ReadFile(extractFileName); err != nil {
			return err
		} else if _, fsContentsHash, err := documents.Hash(data); err != nil {
			return err
		} else if contentsHash != fsContentsHash {
			mustExport = true
		}
		if !mustExport {
			continue
		}
		err = os.WriteFile(extractFileName, reportExtractFile.Contents, 0o644)
		if err != nil {
			return err
		}

		// restore timestamps
		err = os.Chtimes(extractFileName, reportExtractFile.ModifiedAt, reportExtractFile.ModifiedAt)
		if err != nil {
			return err
		}

		// release memory to the GC
		reportExtractFile.Contents = nil

		log.Printf("export: created %q\n", extractFileName)
	}
	return nil
}
