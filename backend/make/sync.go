// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package make

import (
	"log"
	"os"
	"regexp"
	"sort"
	"time"

	"github.com/playbymail/ottoapp/backend/domains"
	"github.com/playbymail/ottoapp/backend/stores/sqlite"
)

// Sync the database and file system (config, reports, maps)

// sync extract files to the database
// * generate missing extract files
// * warning about any "deleted" extract files
// sync map files to the database
// * generate missing map files
// * warn about any "deleted" map files

// config/tn3.1.json               ->  upload
// archive/{turn}/{clan}.docx      ->  {db}/{turn}.{clan}.report.docx
// {db}/{turn}.{clan}.report.docx  ->  {db}/{turn}.{clan}.extract.txt
// {db}/{turn}.{clan}.extract.txt  ->  {db}/{turn}.{clan}.wxx
// {db}/{turn}.{clan}.wxx          ->  {clan}/data/outut/{turn}.{clan}.wxx

// Dependencies
// * {db}/{turn}.{clan}.report.docx       <-  archive/{turn}/{clan}.docx
// * {db}/{turn}.{clan}.extract.txt       <-  {db}/{turn}.{clan}.report.docx
// * {db}/{turn}.{clan}.wxx               <-  {db}/{turn*}.{clan}.extract.txt

// SyncConfigFile to the database.
func SyncConfigFile(db *sqlite.DB, path, name string, quiet, verbose, debug bool) error {
	log.Printf("sync: config-file: path %q\n", path)
	log.Printf("sync: config-file: name %q\n", name)
	panic("obsolete: call sync.Service directly")
	//if verbose {
	//	log.Printf("sync: config-file: path %q\n", path)
	//	log.Printf("sync: config-file: name %q\n", name)
	//}
	//
	//// load the configuration file
	//var data games.ImportFile
	//if input, err := os.ReadFile(filepath.Join(path, name)); err != nil {
	//	if debug {
	//		log.Printf("sync: SyncConfigFile: readFile(%q) %v\n", name, err)
	//	}
	//	return err
	//} else if err = json.Unmarshal(input, &data); err != nil {
	//	if debug {
	//		log.Printf("sync: SyncConfigFile: unmarshal(%q) %v\n", name, err)
	//	}
	//	return err
	//}
	//
	//// initialize services needed to import the configuration file
	//gameSvc, err := games.New(db, nil, nil, nil)
	//if err != nil {
	//	if debug {
	//		log.Printf("sync: SyncConfigFile: games: new() %v\n", err)
	//	}
	//	return fmt.Errorf("sync: %w", err)
	//}
	//
	//// load changes from the configuration file
	//err = gameSvc.Import(&data)
	//if err != nil {
	//	if debug {
	//		log.Printf("sync: SyncConfigFile: games: import() %v\n", err)
	//	}
	//	return fmt.Errorf("sync: import: %w", err)
	//}
	//
	//// display newly created passwords
	//for _, player := range data.Players {
	//	if player.ChangedPassword() {
	//		fmt.Printf("%s: new password %q\n", player.Handle, player.Password)
	//	}
	//}
	//
	//return nil
}

// SyncReportFiles to the database.
func SyncReportFiles(db *sqlite.DB, path string, quiet, verbose, debug bool) error {
	log.Printf("sync: report-files: path %q\n", path)
	panic("obsolete: call sync.Service directly")
	//if verbose {
	//	log.Printf("sync: report-files: path %q\n", path)
	//}
	//
	//// find the files directory
	//filesPath := filepath.Join(path, "files")
	//if verbose {
	//	log.Printf("sync: report-files: filesPath %q\n", filesPath)
	//}
	//if !isdir(filesPath) {
	//	if debug {
	//		log.Printf("sync: report-files: %s: %v\n", filesPath, "not a folder")
	//	}
	//	return fmt.Errorf("%s: not a folder", filesPath)
	//}
	//
	//// initialize services needed to sync the report files
	//authzSvc := authz.New(db)
	//documentsSvc, err := documents.New(db, authzSvc, nil)
	//if err != nil {
	//	if debug {
	//		log.Printf("sync: report-files: documents: new %v\n", err)
	//	}
	//	return fmt.Errorf("sync: %w", err)
	//}
	//gameSvc, err := games.New(db, nil, authzSvc, nil)
	//if err != nil {
	//	if debug {
	//		log.Printf("sync: report-files: games: new %v\n", err)
	//	}
	//	return fmt.Errorf("sync: %w", err)
	//}
	//
	//// fetch the list of games from the database
	//gamesList, err := gameSvc.ReadGames()
	//if err != nil {
	//	if debug {
	//		log.Printf("sync: report-files: games: ReadGames %v\n", err)
	//	}
	//	return fmt.Errorf("sync: getGamesList %w", err)
	//}
	//sort.Slice(gamesList, func(i, j int) bool {
	//	return gamesList[i].ID < gamesList[j].ID
	//})
	//
	//// locate game and turn report directories; quit if any are missing
	//gamePaths, turnReportPaths := map[domains.GameID]string{}, map[domains.GameID]string{}
	//missingGameDirs, missingTurnReportsDirs := false, false
	//for _, game := range gamesList {
	//	gameDir := filepath.Join(filesPath, string(game.ID))
	//	if !isdir(gameDir) {
	//		log.Printf("sync: game: %q: not a folder\n", gameDir)
	//		missingGameDirs = true
	//		continue
	//	}
	//	gamePaths[game.ID] = gameDir
	//	turnReportsDir := filepath.Join(gameDir, "turn-reports")
	//	if !isdir(turnReportsDir) {
	//		log.Printf("sync: game: %q: not a folder\n", turnReportsDir)
	//		missingTurnReportsDirs = true
	//		continue
	//	}
	//	turnReportPaths[game.ID] = turnReportsDir
	//}
	//if missingGameDirs || missingTurnReportsDirs {
	//	if debug && missingGameDirs {
	//		log.Printf("sync: report-files: games: missing game folders\n")
	//	}
	//	if debug && missingTurnReportsDirs {
	//		log.Printf("sync: report-files: games: missing game turn-reports folders\n")
	//	}
	//	return fmt.Errorf("missing folders")
	//}
	//
	//// create a list of files in the turn files directory ({filesPath}/{game}/turn-reports/{game}.{turn}.{clan}.docx)
	//var turnReportFiles []*ReportFile
	//for _, game := range gamesList {
	//	reportFiles, err := walkTurnReportsPath(game.ID, turnReportPaths[game.ID], quiet, verbose, debug)
	//	if err != nil {
	//		if debug {
	//			log.Printf("sync: report-files: walkTurnReportPath %v\n", err)
	//		}
	//		return err
	//	}
	//	if verbose {
	//		log.Printf("sync: report-files: game %s: found %8d report files\n", game.ID, len(reportFiles))
	//	}
	//	turnReportFiles = append(turnReportFiles, reportFiles...)
	//}
	//turnReportFiles = sortReportFiles(turnReportFiles)
	//if verbose {
	//	log.Printf("sync: report-files: found %8d report files\n", len(turnReportFiles))
	//}
	//
	//// cache actors and clans for uploading documents
	//actorCache, clanCache := map[domains.ID]*domains.Actor{}, map[string]*domains.Clan{}
	//actorCache[authz.SysopId] = &domains.Actor{ID: authz.SysopId, Sysop: true}
	//for _, file := range turnReportFiles {
	//	gameClan := fmt.Sprintf("%s:%04d", file.GameID, file.Clan.ClanNo)
	//	clan, ok := clanCache[gameClan]
	//	if !ok {
	//		clan, err = gameSvc.ReadClan(file.GameID, file.Clan.ClanNo)
	//		if err != nil {
	//			log.Printf("sync: cache: game %q: clan %4d: %q %v\n", file.GameID, file.Clan.ClanNo, gameClan, err)
	//			continue
	//		}
	//		clanCache[gameClan] = clan
	//	}
	//	file.Clan = clan
	//	//if debug {
	//	//	log.Printf("sync: cache: game %q: turn %q: clan %4d\n", file.GameID, file.TurnID, file.Clan.ClanNo)
	//	//	log.Printf("sync: cache: game %q: clan %4d: %q %8d\n", file.GameID, file.Clan.ClanNo, gameClan, clan.ClanID)
	//	//}
	//}
	//if verbose {
	//	log.Printf("sync: cache: cached %8d actors\n", len(actorCache))
	//	log.Printf("sync: cache: cached %8d clans\n", len(clanCache))
	//}
	//
	//// upload any new files ({game}.{turn}.{clan}.docx ->  {game}.{turn}.{clan}.docx)
	//for _, file := range turnReportFiles {
	//	actor := actorCache[authz.SysopId]
	//	// load the turn report file from the file system
	//	contents, err := os.ReadFile(file.Path)
	//	if err != nil {
	//		log.Printf("sync: report-files: %s: %v\n", file.Path, err)
	//		return err
	//	}
	//	doc := &domains.Document{
	//		GameID:     file.GameID,
	//		ClanId:     file.Clan.ClanID,
	//		TurnNo:     file.TurnNo,
	//		ClanNo:     file.Clan.ClanNo,
	//		UnitId:     string(file.UnitId),
	//		Path:       filepath.Base(file.Path),
	//		Type:       domains.TurnReportFile,
	//		Contents:   contents,
	//		ModifiedAt: file.ModTime,
	//		CreatedAt:  file.ModTime,
	//		UpdatedAt:  file.ModTime,
	//	}
	//	_, err = documentsSvc.SyncDocument(actor, file.Clan, doc, true, false, debug)
	//	if err != nil {
	//		log.Printf("sync: report-files: %q: %v\n", file.Path, err)
	//		continue
	//	}
	//	if verbose {
	//		log.Printf("sync: report-files: %q: %d\n", doc.Path, doc.ID)
	//		continue
	//	}
	//}
	//
	//// warn about any "deleted" files
	//
	//return nil
}

// SyncReportExtractFiles to the database.
func SyncReportExtractFiles(db *sqlite.DB, path string, quiet, verbose, debug bool) error {
	log.Printf("sync: report-extract-files: path %q\n", path)
	panic("obsolete: call sync.Service directly")
	//if verbose {
	//	log.Printf("sync: report-extract-files: path %q\n", path)
	//}
	//
	//// find the files directory
	//filesPath := filepath.Join(path, "files")
	//if verbose {
	//	log.Printf("sync: filesPath %q\n", filesPath)
	//}
	//if !isdir(filesPath) {
	//	if debug {
	//		log.Printf("sync: filesPath %s: %v\n", filesPath, "not a folder")
	//	}
	//	return fmt.Errorf("%s: not a folder", filesPath)
	//}
	//
	//// initialize services needed to sync the report extract files
	//authzSvc := authz.New(db)
	//documentsSvc, err := documents.New(db, authzSvc, nil)
	//if err != nil {
	//	if debug {
	//		log.Printf("sync: documents: new %v\n", err)
	//	}
	//	return fmt.Errorf("sync: %w", err)
	//}
	//gamesSvc, err := games.New(db, nil, authzSvc, nil)
	//if err != nil {
	//	if debug {
	//		log.Printf("sync: games: new %v\n", err)
	//	}
	//	return fmt.Errorf("sync: %w", err)
	//}
	//
	//// fetch the list of games from the database
	//gamesList, err := gamesSvc.ReadGames()
	//if err != nil {
	//	if debug {
	//		log.Printf("sync: games: ReadGames %v\n", err)
	//	}
	//	return fmt.Errorf("sync: ReadGames %w", err)
	//}
	//sort.Slice(gamesList, func(i, j int) bool {
	//	return gamesList[i].ID < gamesList[j].ID
	//})
	//
	//// locate game and report extract directories; quit if any are missing
	//type gameData_t struct {
	//	GameId             domains.GameID
	//	Clans              map[string]*domains.Clan
	//	ReportExtractFiles []*ReportFile
	//}
	//gameData := map[domains.GameID]*gameData_t{}
	//for _, game := range gamesList {
	//	gameDir := filepath.Join(filesPath, string(game.ID))
	//	if verbose {
	//		log.Printf("sync: gameDir %q\n", gameDir)
	//	}
	//	if !isdir(gameDir) {
	//		log.Printf("sync: game %s: missing game folder\n", game.ID)
	//		return fmt.Errorf("missing folders")
	//	}
	//	ottomapDir := filepath.Join(gameDir, "ottomap")
	//	if verbose {
	//		log.Printf("sync: ottomapDir %q\n", ottomapDir)
	//	}
	//	if !isdir(ottomapDir) {
	//		log.Printf("sync: game %s: missing game ottomap folder\n", game.ID)
	//		return fmt.Errorf("missing folders")
	//	}
	//
	//	thisGameData := &gameData_t{
	//		GameId: game.ID,
	//		Clans:  map[string]*domains.Clan{},
	//	}
	//	gameData[game.ID] = thisGameData
	//
	//	// fetch clans in this game
	//	clanList, err := gamesSvc.ReadClansByGame(game.ID, quiet, verbose, debug)
	//	if err != nil {
	//		if debug {
	//			log.Printf("sync: ReadClansByGames %v\n", err)
	//		}
	//		return fmt.Errorf("sync: ReadClansByGames %w", err)
	//	}
	//
	//	// create a list of files in the report file extract directories
	//	for _, clan := range clanList {
	//		clanNo := fmt.Sprintf("%04d", clan.ClanNo)
	//		thisGameData.Clans[clanNo] = clan
	//		reportExtractsDir := filepath.Join(ottomapDir, clanNo, "data", "input")
	//		if !isdir(reportExtractsDir) {
	//			log.Printf("sync: game %q: clan %q: missing report extracts folder\n", game.ID, clanNo)
	//			continue
	//		}
	//		entries, err := os.ReadDir(reportExtractsDir)
	//		if err != nil {
	//			log.Printf("sync: game %q: clan %q: walk %v\n", game.ID, clanNo, err)
	//			return err
	//		}
	//		for _, entry := range entries {
	//			if entry.IsDir() {
	//				continue
	//			}
	//			if debug {
	//				log.Printf("sync: game %q: clan %q: %q\n", game.ID, clanNo, entry.Name())
	//			}
	//			matches := reReportExtractFile.FindStringSubmatch(entry.Name())
	//			if matches == nil { // not a report file
	//				continue
	//			} else if matches[2] != clanNo { // report file from different clan
	//				log.Printf("sync: game %q: clan %q: %q\n", game.ID, clanNo, entry.Name())
	//				continue
	//			}
	//			fi, err := entry.Info()
	//			if err != nil || !fi.Mode().IsRegular() {
	//				continue
	//			}
	//			reportExtractFile := &ReportFile{
	//				Name:    filepath.Join(reportExtractsDir, entry.Name()),
	//				GameID:  domains.GameID(matches[1]),
	//				TurnNo:  domains.TurnID(matches[2]),
	//				Clan:    clan,
	//				UnitId:  domains.UnitId(clanNo),
	//				Path:    filepath.Join(path, entry.Name()),
	//				ModTime: fi.ModTime().UTC(),
	//			}
	//			thisGameData.ReportExtractFiles = append(thisGameData.ReportExtractFiles, reportExtractFile)
	//		}
	//	}
	//	thisGameData.ReportExtractFiles = sortReportFiles(thisGameData.ReportExtractFiles)
	//
	//	if verbose {
	//		log.Printf("sync: game %q: clans %8d\n", game.ID, len(thisGameData.Clans))
	//		log.Printf("sync: game %q: files %8d\n", game.ID, len(thisGameData.ReportExtractFiles))
	//	}
	//}
	//
	//// upload any new files
	//actor := &domains.Actor{ID: authz.SysopId, Sysop: true}
	//for _, thisGameData := range gameData {
	//	for _, file := range thisGameData.ReportExtractFiles {
	//		// load the turn report file from the file system
	//		contents, err := os.ReadFile(file.Name)
	//		if err != nil {
	//			log.Printf("sync: game %q: %q %v\n", thisGameData.GameId, file.Name, err)
	//			return err
	//		}
	//		doc := &domains.Document{
	//			GameID:     file.GameID,
	//			ClanId:     file.Clan.ClanID,
	//			TurnNo:     file.TurnNo,
	//			ClanNo:     file.Clan.ClanNo,
	//			UnitId:     string(file.UnitId),
	//			Path:       file.Path,
	//			Type:       domains.TurnReportExtract,
	//			Contents:   contents,
	//			ModifiedAt: file.ModTime,
	//			CreatedAt:  file.ModTime,
	//			UpdatedAt:  file.ModTime,
	//		}
	//		documentId, err := documentsSvc.SyncDocument(actor, file.Clan, doc, quiet, verbose, debug)
	//		if err != nil {
	//			log.Printf("sync: game %q: %q %v\n", thisGameData.GameId, file.Name, err)
	//			continue
	//		}
	//		if verbose {
	//			log.Printf("sync: game %q: %q %d\n", thisGameData.GameId, doc.Path, documentId)
	//		}
	//	}
	//}
	//
	//return nil
}

func isdir(path string) bool {
	sb, err := os.Stat(path)
	if err != nil {
		return false
	}
	return sb.IsDir()
}

func isfile(path string) bool {
	sb, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !sb.IsDir() && sb.Mode().IsRegular()
}

type ReportFile struct {
	Name    string // {game}.{turn}.{clan}.docx
	GameID  domains.GameID
	TurnNo  domains.TurnID
	Clan    *domains.Clan
	UnitId  domains.UnitId
	Path    string    // full path to report file
	ModTime time.Time // from os.Stat, assumed UTC?
}

// sort report files by game, turn, clan, and unitId
func sortReportFiles(files []*ReportFile) []*ReportFile {
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

var (
	// report file name is {game}.{turnNo}.{unitId}.docx
	reReportFile = regexp.MustCompile(`^(\d{4})\.(\d{4}-\d{2})\.(\d{4}([cefg][1-9])?)\.docx$`)
)

func walkTurnReportsPath(gameId domains.GameID, path string, quiet, verbose, debug bool) ([]*ReportFile, error) {
	log.Printf("sync: walk: %s\n", path)
	panic("obsolete")
	//if verbose {
	//	log.Printf("sync: walk: %s\n", path)
	//}
	//
	//entries, err := os.ReadDir(path)
	//if err != nil {
	//	if debug {
	//		log.Printf("sync: walk: %s: %v\n", path, err)
	//	}
	//	return nil, err
	//}
	//
	//var reportFiles []*ReportFile
	//for _, entry := range entries {
	//	if debug {
	//		log.Printf("sync: walk: %s\n", entry.Name())
	//	}
	//	if entry.IsDir() {
	//		if debug {
	//			log.Printf("sync: walk: %s: is dir\n", entry.Name())
	//		}
	//		continue
	//	}
	//	matches := reReportFile.FindStringSubmatch(entry.Name())
	//	if matches == nil { // not a report file
	//		continue
	//	} else if matches[1] != string(gameId) { // report file from different game
	//		if debug {
	//			log.Printf("sync: walk: %s: game %s: match[1] %q\n", entry.Name(), gameId, matches[1])
	//		}
	//		continue
	//	}
	//	fi, err := entry.Info()
	//	if err != nil {
	//		return nil, err
	//	} else if !fi.Mode().IsRegular() {
	//		continue
	//	}
	//
	//	// extract clan from unit id
	//	clanNo, _ := strconv.Atoi(matches[3][:4])
	//	if clanNo > 999 {
	//		if debug {
	//			log.Printf("sync: walk: %s\n", matches[3])
	//		}
	//		clanNo = clanNo % 1000
	//	}
	//
	//	reportFiles = append(reportFiles, &ReportFile{
	//		Name:    entry.Name(),
	//		GameID:  domains.GameID(matches[1]),
	//		TurnNo:  domains.TurnID(matches[2]),
	//		Clan:    &domains.Clan{GameID: matches[1], ClanNo: clanNo},
	//		UnitId:  domains.UnitId(matches[3]),
	//		Path:    filepath.Join(path, entry.Name()),
	//		ModTime: fi.ModTime().UTC(),
	//	})
	//}
	//return reportFiles, nil
}
