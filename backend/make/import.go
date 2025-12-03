// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package make

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"time"

	"github.com/playbymail/ottoapp/backend/domains"
	"github.com/playbymail/ottoapp/backend/services/authz"
	"github.com/playbymail/ottoapp/backend/services/documents"
	"github.com/playbymail/ottoapp/backend/services/games"
	"github.com/playbymail/ottoapp/backend/stores/sqlite"
)

var (
	// map file name is {turnNo}.{unitId}.wxx
	reMapFile = regexp.MustCompile(`^(\d{4}-\d{2})\.(0\d{3})\.wxx$`)
	// report extract file name is {turnNo}.{unitId}.report.txt
	reReportExtractFile = regexp.MustCompile(`^(\d{4}-\d{2})\.(0\d{3})\.report\.txt$`)
)

type MapRenderFile struct {
	Name    string // {turn}.{clan}.wxx
	GameID  domains.GameID
	TurnNo  domains.TurnNo
	Clan    *domains.Clan
	UnitId  domains.UnitId
	Path    string    // full path to map file
	ModTime time.Time // from os.Stat, assumed UTC?
}

// sort map render files by game, turn, clan, and unitId
func sortMapRenderFiles(files []*MapRenderFile) []*MapRenderFile {
	sort.Slice(files, func(i, j int) bool {
		a, b := files[i], files[j]
		if a.GameID < b.GameID {
			return true
		} else if a.GameID > b.GameID {
			return false
		}
		if a.TurnNo < b.TurnNo {
			return true
		} else if a.TurnNo > b.TurnNo {
			return false
		}
		if a.Clan.ClanNo < b.Clan.ClanNo {
			return true
		} else if a.Clan.ClanNo > b.Clan.ClanNo {
			return false
		}
		return a.UnitId < b.UnitId
	})
	return files
}

type ReportExtractFile struct {
	Name    string // {turn}.{clan}.report.txt
	GameID  domains.GameID
	TurnNo  domains.TurnNo
	Clan    *domains.Clan
	UnitId  domains.UnitId
	Path    string    // full path to report extract file
	ModTime time.Time // from os.Stat, assumed UTC?
}

// sort report extract files by game, turn, clan, and unitId
func sortReportExtractFiles(files []*ReportExtractFile) []*ReportExtractFile {
	sort.Slice(files, func(i, j int) bool {
		a, b := files[i], files[j]
		if a.GameID < b.GameID {
			return true
		} else if a.GameID > b.GameID {
			return false
		}
		if a.TurnNo < b.TurnNo {
			return true
		} else if a.TurnNo > b.TurnNo {
			return false
		}
		if a.Clan.ClanNo < b.Clan.ClanNo {
			return true
		} else if a.Clan.ClanNo > b.Clan.ClanNo {
			return false
		}
		return a.UnitId < b.UnitId
	})
	return files
}

// ImportMapFiles uploads all map files to the database.
// It replaces any existing files. If a file is in the database but not
// in the file system, it will not be touched.
func ImportMapFiles(db *sqlite.DB, path string, quiet, verbose, debug bool) error {
	if verbose {
		log.Printf("import: path %q\n", path)
	}

	// find the files directory
	filesPath := filepath.Join(path, "files")
	if verbose {
		log.Printf("import: filesPath %q\n", filesPath)
	}
	if !isdir(filesPath) {
		if debug {
			log.Printf("import: filesPath %q: %v\n", filesPath, "not a folder")
		}
		return fmt.Errorf("%s: not a folder", filesPath)
	}

	// initialize services needed to import the map render files
	authzSvc := authz.New(db)
	documentsSvc, err := documents.New(db, authzSvc, nil)
	if err != nil {
		if debug {
			log.Printf("import: documents: new %v\n", err)
		}
		return fmt.Errorf("import: %w", err)
	}
	gamesSvc, err := games.New(db, nil, authzSvc, nil)
	if err != nil {
		if debug {
			log.Printf("import: games: new %v\n", err)
		}
		return fmt.Errorf("import: %w", err)
	}

	// fetch the list of games from the database
	gamesList, err := gamesSvc.GetGamesList()
	if err != nil {
		if debug {
			log.Printf("import: games: GetGamesList %v\n", err)
		}
		return fmt.Errorf("import: GetGamesList %w", err)
	}
	sort.Slice(gamesList, func(i, j int) bool {
		return gamesList[i].ID < gamesList[j].ID
	})

	// locate game and map output directories; quit if any are missing
	type gameData_t struct {
		GameId         domains.GameID
		Clans          map[string]*domains.Clan
		MapRenderFiles []*MapRenderFile
	}
	gameData := map[domains.GameID]*gameData_t{}
	for _, game := range gamesList {
		gameDir := filepath.Join(filesPath, string(game.ID))
		if verbose {
			log.Printf("import: gameDir %q\n", gameDir)
		}
		if !isdir(gameDir) {
			log.Printf("import: game %s: missing game folder\n", game.ID)
			return fmt.Errorf("missing folders")
		}
		ottomapDir := filepath.Join(gameDir, "ottomap")
		if verbose {
			log.Printf("import: ottomapDir %q\n", ottomapDir)
		}
		if !isdir(ottomapDir) {
			log.Printf("import: game %s: missing game ottomap folder\n", game.ID)
			return fmt.Errorf("missing folders")
		}

		thisGameData := &gameData_t{
			GameId: game.ID,
			Clans:  map[string]*domains.Clan{},
		}
		gameData[game.ID] = thisGameData

		// fetch clans in this game
		clanList, err := gamesSvc.ReadClansByGame(game.ID, quiet, verbose, debug)
		if err != nil {
			if debug {
				log.Printf("import: ReadClansByGames %v\n", err)
			}
			return fmt.Errorf("import: ReadClansByGames %w", err)
		}

		// create a list of files in the map file render directories
		for _, clan := range clanList {
			clanNo := fmt.Sprintf("%04d", clan.ClanNo)
			thisGameData.Clans[clanNo] = clan
			mapRenderDir := filepath.Join(ottomapDir, clanNo, "data", "output")
			if !isdir(mapRenderDir) {
				log.Printf("import: game %q: clan %q: missing map render output folder\n", game.ID, clanNo)
				continue
			}
			entries, err := os.ReadDir(mapRenderDir)
			if err != nil {
				log.Printf("import: game %q: clan %q: walk %v\n", game.ID, clanNo, err)
				return err
			}
			for _, entry := range entries {
				if entry.IsDir() {
					continue
				}
				if debug {
					log.Printf("import: game %q: clan %q: %q\n", game.ID, clanNo, entry.Name())
				}
				matches := reMapFile.FindStringSubmatch(entry.Name())
				if matches == nil { // not a map file
					continue
				} else if matches[2] != clanNo { // map file from different clan
					log.Printf("import: game %q: clan %q: %q\n", game.ID, clanNo, entry.Name())
					continue
				}
				fi, err := entry.Info()
				if err != nil || !fi.Mode().IsRegular() {
					continue
				}
				mapRenderFile := &MapRenderFile{
					Name:    filepath.Join(mapRenderDir, entry.Name()),
					GameID:  domains.GameID(matches[1]),
					TurnNo:  domains.TurnNo(matches[2]),
					Clan:    clan,
					UnitId:  domains.UnitId(clanNo),
					Path:    fmt.Sprintf("%s.%s", game.ID, entry.Name()),
					ModTime: fi.ModTime().UTC(),
				}
				thisGameData.MapRenderFiles = append(thisGameData.MapRenderFiles, mapRenderFile)
			}
		}
		thisGameData.MapRenderFiles = sortMapRenderFiles(thisGameData.MapRenderFiles)

		if verbose {
			log.Printf("import: game %q: clans %8d\n", game.ID, len(thisGameData.Clans))
			log.Printf("import: game %q: files %8d\n", game.ID, len(thisGameData.MapRenderFiles))
		}
	}

	// upload all map render files
	actor := &domains.Actor{ID: authz.SysopId, Sysop: true}
	for _, thisGameData := range gameData {
		for _, file := range thisGameData.MapRenderFiles {
			// load the map render file from the file system
			contents, err := os.ReadFile(file.Name)
			if err != nil {
				log.Printf("import: game %q: %q %v\n", thisGameData.GameId, file.Name, err)
				return err
			}
			doc := &domains.Document{
				GameID:     file.GameID,
				ClanId:     file.Clan.ClanID,
				TurnNo:     file.TurnNo,
				ClanNo:     file.Clan.ClanNo,
				UnitId:     string(file.UnitId),
				Path:       file.Path,
				Type:       domains.WorldographerMap,
				Contents:   contents,
				ModifiedAt: file.ModTime,
				CreatedAt:  file.ModTime,
				UpdatedAt:  file.ModTime,
			}
			documentId, err := documentsSvc.ReplaceDocument(actor, file.Clan, doc, quiet, verbose, debug)
			if err != nil {
				log.Printf("import: game %q: %q %v\n", thisGameData.GameId, file.Name, err)
				continue
			}
			if verbose {
				log.Printf("import: game %q: %q %d\n", thisGameData.GameId, doc.Path, documentId)
			}
		}
	}

	return nil
}

// ImportReportExtractFiles uploads all report extract files to the database.
// It replaces any existing files. If a file is in the database but not
// in the file system, it will not be touched.
func ImportReportExtractFiles(db *sqlite.DB, path string, quiet, verbose, debug bool) error {
	if verbose {
		log.Printf("import: path %q\n", path)
	}

	// find the files directory
	filesPath := filepath.Join(path, "files")
	if verbose {
		log.Printf("import: filesPath %q\n", filesPath)
	}
	if !isdir(filesPath) {
		if debug {
			log.Printf("import: filesPath %q: %v\n", filesPath, "not a folder")
		}
		return fmt.Errorf("%s: not a folder", filesPath)
	}

	// initialize services needed to import the report extract files
	authzSvc := authz.New(db)
	documentsSvc, err := documents.New(db, authzSvc, nil)
	if err != nil {
		if debug {
			log.Printf("import: documents: new %v\n", err)
		}
		return fmt.Errorf("import: %w", err)
	}
	gamesSvc, err := games.New(db, nil, authzSvc, nil)
	if err != nil {
		if debug {
			log.Printf("import: games: new %v\n", err)
		}
		return fmt.Errorf("import: %w", err)
	}

	// fetch the list of games from the database
	gamesList, err := gamesSvc.GetGamesList()
	if err != nil {
		if debug {
			log.Printf("import: games: GetGamesList %v\n", err)
		}
		return fmt.Errorf("import: GetGamesList %w", err)
	}
	sort.Slice(gamesList, func(i, j int) bool {
		return gamesList[i].ID < gamesList[j].ID
	})

	// locate game and report extract directories; quit if any are missing
	type gameData_t struct {
		GameId             domains.GameID
		Clans              map[string]*domains.Clan
		ReportExtractFiles []*ReportExtractFile
	}
	gameData := map[domains.GameID]*gameData_t{}
	for _, game := range gamesList {
		gameDir := filepath.Join(filesPath, string(game.ID))
		if verbose {
			log.Printf("import: gameDir %q\n", gameDir)
		}
		if !isdir(gameDir) {
			log.Printf("import: game %s: missing game folder\n", game.ID)
			return fmt.Errorf("missing folders")
		}
		ottomapDir := filepath.Join(gameDir, "ottomap")
		if verbose {
			log.Printf("import: ottomapDir %q\n", ottomapDir)
		}
		if !isdir(ottomapDir) {
			log.Printf("import: game %s: missing game ottomap folder\n", game.ID)
			return fmt.Errorf("missing folders")
		}

		thisGameData := &gameData_t{
			GameId: game.ID,
			Clans:  map[string]*domains.Clan{},
		}
		gameData[game.ID] = thisGameData

		// fetch clans in this game
		clanList, err := gamesSvc.ReadClansByGame(game.ID, quiet, verbose, debug)
		if err != nil {
			if debug {
				log.Printf("import: ReadClansByGames %v\n", err)
			}
			return fmt.Errorf("import: ReadClansByGames %w", err)
		}

		// create a list of files in the map file render directories
		for _, clan := range clanList {
			clanNo := fmt.Sprintf("%04d", clan.ClanNo)
			thisGameData.Clans[clanNo] = clan
			reportExtractsDir := filepath.Join(ottomapDir, clanNo, "data", "input")
			if !isdir(reportExtractsDir) {
				log.Printf("import: game %q: clan %q: missing report extracts folder\n", game.ID, clanNo)
				continue
			}
			entries, err := os.ReadDir(reportExtractsDir)
			if err != nil {
				log.Printf("import: game %q: clan %q: walk %v\n", game.ID, clanNo, err)
				return err
			}
			for _, entry := range entries {
				if entry.IsDir() {
					continue
				}
				if debug {
					log.Printf("import: game %q: clan %q: %q\n", game.ID, clanNo, entry.Name())
				}
				matches := reReportExtractFile.FindStringSubmatch(entry.Name())
				if matches == nil { // not a map file
					continue
				} else if matches[2] != clanNo { // map file from different clan
					log.Printf("import: game %q: clan %q: %q\n", game.ID, clanNo, entry.Name())
					continue
				}
				fi, err := entry.Info()
				if err != nil || !fi.Mode().IsRegular() {
					continue
				}
				reportExtractFile := &ReportExtractFile{
					Name:    filepath.Join(reportExtractsDir, entry.Name()),
					GameID:  domains.GameID(matches[1]),
					TurnNo:  domains.TurnNo(matches[2]),
					Clan:    clan,
					UnitId:  domains.UnitId(clanNo),
					Path:    fmt.Sprintf("%s.%s", game.ID, entry.Name()),
					ModTime: fi.ModTime().UTC(),
				}
				thisGameData.ReportExtractFiles = append(thisGameData.ReportExtractFiles, reportExtractFile)
			}
		}
		thisGameData.ReportExtractFiles = sortReportExtractFiles(thisGameData.ReportExtractFiles)

		if verbose {
			log.Printf("import: game %q: clans %8d\n", game.ID, len(thisGameData.Clans))
			log.Printf("import: game %q: files %8d\n", game.ID, len(thisGameData.ReportExtractFiles))
		}
	}

	// upload all report extract files
	actor := &domains.Actor{ID: authz.SysopId, Sysop: true}
	for _, thisGameData := range gameData {
		for _, file := range thisGameData.ReportExtractFiles {
			// load the report extract file from the file system
			contents, err := os.ReadFile(file.Name)
			if err != nil {
				log.Printf("import: game %q: %q %v\n", thisGameData.GameId, file.Name, err)
				return err
			}
			doc := &domains.Document{
				GameID:     file.GameID,
				ClanId:     file.Clan.ClanID,
				TurnNo:     file.TurnNo,
				ClanNo:     file.Clan.ClanNo,
				UnitId:     string(file.UnitId),
				Path:       file.Path,
				Type:       domains.TurnReportExtract,
				Contents:   contents,
				ModifiedAt: file.ModTime,
				CreatedAt:  file.ModTime,
				UpdatedAt:  file.ModTime,
			}
			documentId, err := documentsSvc.ReplaceDocument(actor, file.Clan, doc, quiet, verbose, debug)
			if err != nil {
				log.Printf("import: game %q: %q %v\n", thisGameData.GameId, file.Name, err)
				continue
			}
			if verbose {
				log.Printf("import: game %q: %q %d\n", thisGameData.GameId, doc.Path, documentId)
			}
		}
	}

	return nil
}

// ImportTurnReportFiles uploads all turn report files to the database.
// It replaces any existing files. If a file is in the database but not
// in the file system, it will not be touched.
func ImportTurnReportFiles(db *sqlite.DB, path string, quiet, verbose, debug bool) error {
	if verbose {
		log.Printf("import: path %q\n", path)
	}

	// find the files directory
	filesPath := filepath.Join(path, "files")
	if verbose {
		log.Printf("import: filesPath %q\n", filesPath)
	}
	if !isdir(filesPath) {
		if debug {
			log.Printf("import: filesPath %q: %v\n", filesPath, "not a folder")
		}
		return fmt.Errorf("%s: not a folder", filesPath)
	}

	// initialize services needed to import the turn report files
	authzSvc := authz.New(db)
	documentsSvc, err := documents.New(db, authzSvc, nil)
	if err != nil {
		if debug {
			log.Printf("import: documents: new %v\n", err)
		}
		return fmt.Errorf("import: %w", err)
	}
	gameSvc, err := games.New(db, nil, authzSvc, nil)
	if err != nil {
		if debug {
			log.Printf("import: games: new %v\n", err)
		}
		return fmt.Errorf("import: %w", err)
	}

	// fetch the list of games from the database
	gamesList, err := gameSvc.GetGamesList()
	if err != nil {
		if debug {
			log.Printf("import: games: GetGamesList %v\n", err)
		}
		return fmt.Errorf("import: GetGamesList %w", err)
	}
	sort.Slice(gamesList, func(i, j int) bool {
		return gamesList[i].ID < gamesList[j].ID
	})

	// locate game and turn report directories; quit if any are missing
	gamePaths, turnReportFilesPath := map[domains.GameID]string{}, map[domains.GameID]string{}
	missingGameDirs, missingTurnReportFilesPath := false, false
	for _, game := range gamesList {
		gameDir := filepath.Join(filesPath, string(game.ID))
		if !isdir(gameDir) {
			log.Printf("import: game %q: not a folder\n", gameDir)
			missingGameDirs = true
			continue
		}
		gamePaths[game.ID] = gameDir
		gameTurnReportFilesPath := filepath.Join(gameDir, "turn-reports")
		if !isdir(gameTurnReportFilesPath) {
			log.Printf("import: game: %q: not a folder\n", gameTurnReportFilesPath)
			missingTurnReportFilesPath = true
			continue
		}
		turnReportFilesPath[game.ID] = gameTurnReportFilesPath
	}
	if missingGameDirs || missingTurnReportFilesPath {
		if debug && missingGameDirs {
			log.Printf("import: games: missing game folders\n")
		}
		if debug && missingTurnReportFilesPath {
			log.Printf("import: games: missing game turn-reports folders\n")
		}
		return fmt.Errorf("missing folders")
	}

	// create a list of files in the turn report files directory
	var turnReportFiles []*ReportFile
	for _, game := range gamesList {
		fileNames, err := walkTurnReportsPath(game.ID, turnReportFilesPath[game.ID], quiet, verbose, debug)
		if err != nil {
			if debug {
				log.Printf("import: walkTurnReportPath %v\n", err)
			}
			return err
		}
		if verbose {
			log.Printf("import: game %s: found %8d turn report files\n", game.ID, len(fileNames))
		}
		turnReportFiles = append(turnReportFiles, fileNames...)
	}
	turnReportFiles = sortReportFiles(turnReportFiles)
	if verbose {
		log.Printf("import: found %8d turn report files\n", len(turnReportFiles))
	}

	// cache actors and clans for uploading documents
	actorCache, clanCache := map[domains.ID]*domains.Actor{}, map[string]*domains.Clan{}
	actorCache[authz.SysopId] = &domains.Actor{ID: authz.SysopId, Sysop: true}
	for _, file := range turnReportFiles {
		gameClan := fmt.Sprintf("%s:%04d", file.GameID, file.Clan.ClanNo)
		clan, ok := clanCache[gameClan]
		if !ok {
			clan, err = gameSvc.GetClan(file.GameID, file.Clan.ClanNo)
			if err != nil {
				log.Printf("import: cache: game %q: clan %4d: %q %v\n", file.GameID, file.Clan.ClanNo, gameClan, err)
				continue
			}
			clanCache[gameClan] = clan
		}
		file.Clan = clan
	}
	if verbose {
		log.Printf("import: cached %8d actors\n", len(actorCache))
		log.Printf("import: cached %8d clans\n", len(clanCache))
	}

	// upload all turn report files
	for _, file := range turnReportFiles {
		actor := actorCache[authz.SysopId]
		// load the turn report file from the file system
		contents, err := os.ReadFile(file.Path)
		if err != nil {
			log.Printf("import: %s: %v\n", file.Path, err)
			return err
		}
		doc := &domains.Document{
			GameID:     file.GameID,
			ClanId:     file.Clan.ClanID,
			TurnNo:     file.TurnNo,
			ClanNo:     file.Clan.ClanNo,
			UnitId:     string(file.UnitId),
			Path:       fmt.Sprintf("%s.%s.%04d.docx", file.GameID, file.TurnNo, file.Clan.ClanNo),
			Type:       domains.TurnReportFile,
			Contents:   contents,
			ModifiedAt: file.ModTime,
			CreatedAt:  file.ModTime,
			UpdatedAt:  file.ModTime,
		}
		documentId, err := documentsSvc.ReplaceDocument(actor, file.Clan, doc, quiet, false, debug)
		if err != nil {
			log.Printf("import: ReplaceDocument(%q): %v\n", file.Path, err)
			continue
		}
		if verbose {
			log.Printf("import: %q: %q %d\n", file.Path, doc.Path, documentId)
			continue
		}
	}

	return nil
}
