// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package sync

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"time"

	"github.com/mdhender/phrases/v2"
	"github.com/playbymail/ottoapp/backend/domains"
	"github.com/playbymail/ottoapp/backend/services/authz"
	"github.com/playbymail/ottoapp/backend/stores/jsondb"
	"github.com/playbymail/ottoapp/backend/stores/sqlite/sqlc"
)

// ImportGames imports the games from the configuration file.
// If codes are given, only the games that are in the list of
// codes will be imported.
func (s *Service) ImportGames(path string, codes ...string) error {
	// load the game data then force a validation on the codes
	// for all the games in it, even if we're not going to load them.
	gamesData, err := jsondb.LoadGames(path)
	if err != nil {
		log.Printf("sync: import: games: %s: failed %v\n", path, err)
		return fmt.Errorf("sync: import: games: %w", err)
	}
	for _, code := range codes {
		game, ok := gamesData[code]
		if !ok {
			continue
		}
		err = domains.ValidateGameCode(game.Code)
		if err != nil {
			log.Printf("sync: import: games: %q: %v\n", game.Code, err)
		}
	}
	if err != nil {
		return fmt.Errorf("sync: import: game %w", err)
	}

	// if codes is empty, load all the games. otherwise, load only
	// the games in the list of codes.
	var codesToLoad []string
	if len(codes) > 0 {
		for _, code := range codes {
			codesToLoad = append(codesToLoad, code)
		}
	} else {
		for _, game := range gamesData {
			codesToLoad = append(codesToLoad, game.Code)
		}
	}
	if len(codesToLoad) == 0 {
		return nil
	}
	sort.Strings(codesToLoad)

	// return an error if we're missing any games
	missingGames := 0
	for _, code := range codesToLoad {
		_, ok := gamesData[code]
		if !ok {
			missingGames++
			log.Printf("sync: import: games: %q: missing\n", code)
		}
	}
	if missingGames == 1 {
		return fmt.Errorf("sync: import: game not found")
	} else if missingGames > 1 {
		return fmt.Errorf("sync: import: %d games not found", missingGames)
	}

	gameCache := map[string]int64{}

	ctx := s.db.Context()
	q := s.db.Queries()
	now := time.Now().UTC()
	createdAt, updatedAt := now.Unix(), now.Unix()

	for _, code := range codesToLoad {
		data := gamesData[code]
		err = domains.ValidateGameCode(data.Code)
		if err != nil {
			log.Printf("sync: import: games: %q: %v\n", data.Code, err)
			return fmt.Errorf("sync: import: game %q: %w", data.Code, err)
		}
		gameId, err := q.CreateGame(ctx, sqlc.CreateGameParams{
			Code:        data.Code,
			Description: data.Description,
			SetupTurn:   fmt.Sprintf("%04d-%02d", data.SetupTurn.Year, data.SetupTurn.Month),
			ActiveTurn:  fmt.Sprintf("%04d-%02d", data.ActiveTurn.Year, data.ActiveTurn.Month),
			OrdersDue:   data.OrdersDue.UTC().Unix(),
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		})
		if err != nil {
			log.Printf("sync: import: games: %q: %v\n", data.Code, err)
			return fmt.Errorf("sync: import: game %q: %w", data.Code, err)
		}
		gameCache[data.Code] = gameId

		// import the game turns
		year, month, turnNo := data.SetupTurn.Year, data.SetupTurn.Month, 0
		for year*100+month <= data.ActiveTurn.Year*100+data.ActiveTurn.Month {
			err = q.CreateGameTurn(ctx, sqlc.CreateGameTurnParams{
				GameID:    gameId,
				Turn:      fmt.Sprintf("%04d-%02d", year, month),
				TurnYear:  int64(year),
				TurnMonth: int64(month),
				TurnNo:    int64(turnNo),
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			})
			if err != nil {
				log.Printf("sync: import: games: %q: %v\n", data.Code, err)
				return fmt.Errorf("sync: import: game %q: %w", data.Code, err)
			}
			if month = month + 1; month > 12 {
				year, month = year+1, 1
			}
			turnNo++
		}

		// import the game clans
		for handle, clan := range data.Clans {
			user, err := q.ReadUserByHandle(ctx, handle)
			if err != nil {
				log.Printf("sync: import: games: %q: handle %q: %v\n", data.Code, handle, err)
				return err
			}
			_, err = q.CreateGameUserClan(ctx, sqlc.CreateGameUserClanParams{
				GameID:    gameId,
				UserID:    user.UserID,
				Clan:      int64(clan.ClanNo),
				SetupTurn: clan.SetupTurn,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			})
			if err != nil {
				log.Printf("sync: import: games: %q: handle %q: %v\n", data.Code, handle, err)
				return err
			}
		}
	}
	// log.Printf("sync: import: games: imported %d\n", len(gameCache))

	return nil
}

// ImportMapFiles uploads all worldographer map files to the database.
// It replaces any existing files. If a file is in the database but not
// in the file system, it will not be touched.
func (s *Service) ImportMapFiles(path string, quiet, verbose, debug bool) error {
	if verbose {
		log.Printf("import: path %q\n", path)
	}

	// find the files directory
	filesPath := filepath.Join(path, "files")
	if verbose {
		log.Printf("sync: import: filesPath %q\n", filesPath)
	}
	if !isdir(filesPath) {
		if debug {
			log.Printf("sync: import: filesPath %q: %s\n", filesPath, "not a folder")
		}
		return fmt.Errorf("%s: not a folder", filesPath)
	}

	// fetch the list of games from the database
	gamesList, err := s.gameSvc.ReadGames()
	if err != nil {
		if debug {
			log.Printf("sync: import: games: ReadGames %v\n", err)
		}
		return fmt.Errorf("sync: import: ReadGames %w", err)
	}
	sort.Slice(gamesList, func(i, j int) bool {
		return gamesList[i].Code < gamesList[j].Code
	})

	type mapFile struct {
		Path    string
		Name    string
		code    string
		Turn    string
		Clan    domains.Clan
		ModTime time.Time
	}
	var mapFiles []*mapFile

	// walk game directories
	missingDirs := false
	for _, game := range gamesList {
		gameDir := filepath.Join(filesPath, game.Code)
		if !isdir(gameDir) {
			log.Printf("sync: import: game %q: not a folder\n", gameDir)
			missingDirs = true
			continue
		}
		ottomapDir := filepath.Join(gameDir, "ottomap")
		if !isdir(ottomapDir) {
			log.Printf("sync: import: game %q: missing ottomap folder\n", gameDir)
			missingDirs = true
			continue
		}

		// walk ottomap directories
		ottomapEntries, err := os.ReadDir(ottomapDir)
		if err != nil {
			if debug {
				log.Printf("sync: import: walk: %s: %v\n", ottomapDir, err)
			}
			return err
		}
		for _, ottomapEntry := range ottomapEntries {
			if !ottomapEntry.IsDir() {
				continue
			}
			clan := ottomapEntry.Name()
			if !reClanDataFolder.MatchString(clan) {
				continue
			}
			clanNo, _ := strconv.Atoi(clan)
			// walk the clan output directory, gathering report extract files
			clanOutputPath := filepath.Join(ottomapDir, clan, "data", "output")
			entries, err := os.ReadDir(clanOutputPath)
			if err != nil {
				if debug {
					log.Printf("sync: import: walk: %s: %v\n", clanOutputPath, err)
				}
				return err
			}
			for _, entry := range entries {
				if entry.IsDir() {
					continue
				}
				// clan map file name is {turnNo}.{clanNo}.wxx
				clanMapFileName := entry.Name()
				matches := reClanMapFile.FindStringSubmatch(clanMapFileName)
				if matches == nil {
					log.Printf("%q %+v\n", clanMapFileName, matches)
					continue
				}
				if clan != matches[2] {
					// ignore non-clan files
					continue
				}
				modTime := time.Now().UTC()
				if sb, err := entry.Info(); err == nil {
					modTime = sb.ModTime().UTC()
				}
				mapFiles = append(mapFiles, &mapFile{
					Path: filepath.Join(clanOutputPath, clanMapFileName),
					Name: clanMapFileName,
					code: game.Code,
					Turn: matches[1],
					Clan: domains.Clan{
						GameID: game.ID,
						ClanNo: clanNo,
					},
					ModTime: modTime,
				})
			}
		}
		if missingDirs {
			if debug {
				log.Printf("sync: import: games: missing folders\n")
			}
			return fmt.Errorf("missing folders")
		}
		if verbose {
			log.Printf("sync: import: found %8d report extract files\n", len(mapFiles))
		}

		// sort the list so that we have a deterministic load order for files
		sort.Slice(mapFiles, func(i, j int) bool {
			a, b := mapFiles[i], mapFiles[j]
			return a.Path < b.Path
		})

		// cache actors and clans for uploading documents
		actorCache := map[domains.ID]*domains.Actor{
			authz.SysopId: &domains.Actor{ID: authz.SysopId, Roles: domains.Roles{Sysop: true}},
		}
		clanCache := map[string]*domains.Clan{}
		for _, file := range mapFiles {
			gameClan := fmt.Sprintf("%s.%04d", file.code, file.Clan.ClanNo)
			clan, ok := clanCache[gameClan]
			if !ok {
				clan, err = s.gameSvc.GameIdClanNoToClan(file.Clan.GameID, file.Clan.ClanNo)
				if err != nil {
					log.Printf("sync: import: cache: game %q: clan %4d: %q %v\n", file.code, file.Clan.ClanNo, gameClan, err)
					return err
				}
				clanCache[gameClan] = clan
			}
			file.Clan = *clan
		}
		if verbose {
			log.Printf("import: cached %8d actors\n", len(actorCache))
			log.Printf("import: cached %8d clans\n", len(clanCache))
		}

		for _, file := range mapFiles {
			actor := actorCache[authz.SysopId]
			// load the clan map file from the file system
			contents, err := os.ReadFile(file.Path)
			if err != nil {
				log.Printf("sync: import: %s: %v\n", file.Path, err)
				return err
			}
			doc := &domains.Document{
				GameID:     file.Clan.GameID,
				ClanId:     file.Clan.ClanID,
				Turn:       file.Turn,
				ClanNo:     file.Clan.ClanNo,
				UnitId:     fmt.Sprintf("%4d", file.Clan.ClanNo),
				Path:       fmt.Sprintf("%s.%s.%04d.wxx", file.code, file.Turn, file.Clan.ClanNo),
				Type:       domains.WorldographerMap,
				Contents:   contents,
				ModifiedAt: file.ModTime,
				CreatedAt:  file.ModTime,
				UpdatedAt:  file.ModTime,
			}
			documentId, err := s.documentsSvc.ReplaceDocument(actor, &file.Clan, doc, quiet, false, debug)
			if err != nil {
				log.Printf("sync: import: ReplaceDocument(%q): %v\n", file.Path, err)
				return err
			}
			if verbose {
				log.Printf("sync: import: %q %d\n", doc.Path, documentId)
				continue
			}
		}
	}

	return nil
}

func (s *Service) ImportOttoAppConfig(path string, quiet, verbose, debug bool) error {
	if debug {
		log.Printf("sync: import: %q\n", path)
	}

	oac, err := jsondb.LoadOttoAppConfig(path)
	if err != nil {
		if debug {
			log.Printf("sync: import: load %s: failed %v\n", path, err)
		}
		return fmt.Errorf("sync: import: config %w", err)
	}

	err = s.configSvc.UpdateKeyValuePairs(
		"mailgun.domain", oac.Mailgun.Domain,
		"mailgun.from", oac.Mailgun.From,
		"mailgun.api.base", oac.Mailgun.ApiBase,
		"mailgun.api.key", oac.Mailgun.ApiKey,
	)
	if err != nil {
		if debug {
			log.Printf("sync: import: %s: failed %v\n", path, err)
		}
		return fmt.Errorf("sync: import: config %w", err)
	}

	return nil
}

// ImportReportExtractFiles uploads all report extract files to the database.
// It replaces any existing files. If a file is in the database but not
// in the file system, it will not be touched.
func (s *Service) ImportReportExtractFiles(path string, quiet, verbose, debug bool) error {
	if verbose {
		log.Printf("import: path %q\n", path)
	}

	// find the files directory
	filesPath := filepath.Join(path, "files")
	if verbose {
		log.Printf("sync: import: filesPath %q\n", filesPath)
	}
	if !isdir(filesPath) {
		if debug {
			log.Printf("sync: import: filesPath %q: %s\n", filesPath, "not a folder")
		}
		return fmt.Errorf("%s: not a folder", filesPath)
	}

	// fetch the list of games from the database
	gamesList, err := s.gameSvc.ReadGames()
	if err != nil {
		if debug {
			log.Printf("sync: import: games: ReadGames %v\n", err)
		}
		return fmt.Errorf("sync: import: ReadGames %w", err)
	}
	sort.Slice(gamesList, func(i, j int) bool {
		return gamesList[i].Code < gamesList[j].Code
	})

	type reportExtractFile struct {
		Path    string
		Name    string
		code    string
		Turn    string
		Clan    domains.Clan
		ModTime time.Time
	}
	var reportExtractFiles []*reportExtractFile

	// walk game directories
	missingDirs := false
	for _, game := range gamesList {
		gameDir := filepath.Join(filesPath, game.Code)
		if !isdir(gameDir) {
			log.Printf("sync: import: game %q: not a folder\n", gameDir)
			missingDirs = true
			continue
		}
		ottomapDir := filepath.Join(gameDir, "ottomap")
		if !isdir(ottomapDir) {
			log.Printf("sync: import: game %q: missing ottomap folder\n", gameDir)
			missingDirs = true
			continue
		}

		// walk ottomap directories
		ottomapEntries, err := os.ReadDir(ottomapDir)
		if err != nil {
			if debug {
				log.Printf("sync: import: walk: %s: %v\n", ottomapDir, err)
			}
			return err
		}
		for _, ottomapEntry := range ottomapEntries {
			if !ottomapEntry.IsDir() {
				continue
			}
			clan := ottomapEntry.Name()
			if !reClanDataFolder.MatchString(clan) {
				continue
			}
			clanNo, _ := strconv.Atoi(clan)
			// walk the clan input directory, gathering report extract files
			clanInputPath := filepath.Join(ottomapDir, clan, "data", "input")
			entries, err := os.ReadDir(clanInputPath)
			if err != nil {
				if debug {
					log.Printf("sync: import: walk: %s: %v\n", clanInputPath, err)
				}
				return err
			}
			for _, entry := range entries {
				if entry.IsDir() {
					continue
				}
				// report extract file Name is {turnNo}.{clanNo}.report.txt
				reportExtractName := entry.Name()
				matches := reClanReportExtractFile.FindStringSubmatch(reportExtractName)
				if matches == nil {
					log.Printf("%q %+v\n", reportExtractName, matches)
					continue
				}
				if clan != matches[2] {
					// ignore non-clan files
					continue
				}
				modTime := time.Now().UTC()
				if sb, err := entry.Info(); err == nil {
					modTime = sb.ModTime().UTC()
				}
				reportExtractFiles = append(reportExtractFiles, &reportExtractFile{
					Path: filepath.Join(clanInputPath, reportExtractName),
					Name: reportExtractName,
					code: game.Code,
					Turn: matches[1],
					Clan: domains.Clan{
						GameID: game.ID,
						ClanNo: clanNo,
					},
					ModTime: modTime,
				})
			}
		}
		if missingDirs {
			if debug {
				log.Printf("sync: import: games: missing folders\n")
			}
			return fmt.Errorf("missing folders")
		}
		if verbose {
			log.Printf("sync: import: found %8d report extract files\n", len(reportExtractFiles))
		}

		// sort the list so that we have a deterministic load order for reports
		sort.Slice(reportExtractFiles, func(i, j int) bool {
			a, b := reportExtractFiles[i], reportExtractFiles[j]
			return a.Path < b.Path
		})

		// cache actors and clans for uploading documents
		actorCache := map[domains.ID]*domains.Actor{
			authz.SysopId: &domains.Actor{ID: authz.SysopId, Roles: domains.Roles{Sysop: true}},
		}
		clanCache := map[string]*domains.Clan{}
		for _, file := range reportExtractFiles {
			gameClan := fmt.Sprintf("%s.%04d", file.code, file.Clan.ClanNo)
			clan, ok := clanCache[gameClan]
			if !ok {
				clan, err = s.gameSvc.GameIdClanNoToClan(file.Clan.GameID, file.Clan.ClanNo)
				if err != nil {
					log.Printf("sync: import: cache: game %q: clan %4d: %q %v\n", file.code, file.Clan.ClanNo, gameClan, err)
					return err
				}
				clanCache[gameClan] = clan
			}
			file.Clan = *clan
		}
		if verbose {
			log.Printf("import: cached %8d actors\n", len(actorCache))
			log.Printf("import: cached %8d clans\n", len(clanCache))
		}

		for _, file := range reportExtractFiles {
			actor := actorCache[authz.SysopId]
			// load the Turn report file from the file system
			contents, err := os.ReadFile(file.Path)
			if err != nil {
				log.Printf("sync: import: %s: %v\n", file.Path, err)
				return err
			}
			doc := &domains.Document{
				GameID:     file.Clan.GameID,
				ClanId:     file.Clan.ClanID,
				Turn:       file.Turn,
				ClanNo:     file.Clan.ClanNo,
				UnitId:     fmt.Sprintf("%4d", file.Clan.ClanNo),
				Path:       fmt.Sprintf("%s.%s.%04d.report.txt", file.code, file.Turn, file.Clan.ClanNo),
				Type:       domains.TurnReportExtract,
				Contents:   contents,
				ModifiedAt: file.ModTime,
				CreatedAt:  file.ModTime,
				UpdatedAt:  file.ModTime,
			}
			documentId, err := s.documentsSvc.ReplaceDocument(actor, &file.Clan, doc, quiet, false, debug)
			if err != nil {
				log.Printf("sync: import: ReplaceDocument(%q): %v\n", file.Path, err)
				return err
			}
			if verbose {
				log.Printf("sync: import: %q %d\n", doc.Path, documentId)
				continue
			}
		}
	}

	return nil
}

var (
	// clan folder is 0987
	reClanDataFolder = regexp.MustCompile(`^0\d{3}$`)

	// clan map file Name is {turnNo}.{clan}.wxx
	reClanMapFile = regexp.MustCompile(`^(\d{4}-\d{2})\.(0\d{3})\.wxx$`)

	// report extract file Name is {turnNo}.{clan}.report.txt
	reClanReportExtractFile = regexp.MustCompile(`^(\d{4}-\d{2})\.(0\d{3})\.report\.txt$`)
)

// ImportTurnReportFiles uploads all Turn report files to the database.
// It replaces any existing files. If a file is in the database but not
// in the file system, it will not be touched.
func (s *Service) ImportTurnReportFiles(path string, quiet, verbose, debug bool) error {
	if verbose {
		log.Printf("sync: import: Path %q\n", path)
	}

	// find the files directory
	filesPath := filepath.Join(path, "files")
	if verbose {
		log.Printf("sync: import: filesPath %q\n", filesPath)
	}
	if !isdir(filesPath) {
		if debug {
			log.Printf("sync: import: filesPath %q: %s\n", filesPath, "not a folder")
		}
		return fmt.Errorf("%s: not a folder", filesPath)
	}

	// fetch the list of games from the database
	gamesList, err := s.gameSvc.ReadGames()
	if err != nil {
		if debug {
			log.Printf("sync: import: games: ReadGames %v\n", err)
		}
		return fmt.Errorf("sync: import: ReadGames %w", err)
	}
	sort.Slice(gamesList, func(i, j int) bool {
		return gamesList[i].Code < gamesList[j].Code
	})

	type turnReportFile struct {
		Path    string
		Name    string
		code    string
		Turn    string
		Clan    domains.Clan
		ModTime time.Time
	}
	var turnReportFiles []*turnReportFile

	// walk game and Turn report directories, gathering Turn report files
	missingDirs := false
	for _, game := range gamesList {
		gameDir := filepath.Join(filesPath, game.Code)
		if !isdir(gameDir) {
			log.Printf("sync: import: game %q: not a folder\n", gameDir)
			missingDirs = true
			continue
		}
		turnReportFilePath := filepath.Join(gameDir, "turn-reports")
		if !isdir(turnReportFilePath) {
			log.Printf("sync: import: game: %q: not a folder\n", turnReportFilePath)
			missingDirs = true
			continue
		}
		entries, err := os.ReadDir(turnReportFilePath)
		if err != nil {
			if debug {
				log.Printf("sync: import: walk: %s: %v\n", turnReportFilePath, err)
			}
			return err
		}
		// create a list of Turn report files in the Turn report files directory
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			// report file Name is {code}.{turnNo}.{clanNo}.docx
			reportPath, reportName := turnReportFilePath, entry.Name()
			matches := reClanReportFile.FindStringSubmatch(reportName)
			if matches == nil {
				continue
			}
			code, turn, clan := matches[1], matches[2], matches[3]
			if code != game.Code {
				// ignore non-game files
				continue
			}
			clanNo, _ := strconv.Atoi(clan)
			modTime := time.Now().UTC()
			if sb, err := entry.Info(); err == nil {
				modTime = sb.ModTime().UTC()
			}
			turnReportFiles = append(turnReportFiles, &turnReportFile{
				Path: filepath.Join(reportPath, reportName),
				Name: reportName,
				code: code,
				Turn: turn,
				Clan: domains.Clan{
					GameID: game.ID,
					ClanNo: clanNo,
				},
				ModTime: modTime,
			})
		}
	}
	if missingDirs {
		if debug {
			log.Printf("sync: import: games: missing folders\n")
		}
		return fmt.Errorf("missing folders")
	}
	if verbose {
		log.Printf("sync: import: found %8d turn report files\n", len(turnReportFiles))
	}

	// sort the list so that we have a deterministic load order for reports
	sort.Slice(turnReportFiles, func(i, j int) bool {
		a, b := turnReportFiles[i], turnReportFiles[j]
		return a.Path < b.Path
	})

	// cache actors and clans for uploading documents
	actorCache := map[domains.ID]*domains.Actor{
		authz.SysopId: &domains.Actor{ID: authz.SysopId, Roles: domains.Roles{Sysop: true}},
	}
	clanCache := map[string]*domains.Clan{}
	for _, file := range turnReportFiles {
		gameClan := fmt.Sprintf("%s.%04d", file.code, file.Clan.ClanNo)
		clan, ok := clanCache[gameClan]
		if !ok {
			clan, err = s.gameSvc.GameIdClanNoToClan(file.Clan.GameID, file.Clan.ClanNo)
			if err != nil {
				log.Printf("sync: import: cache: game %q: clan %4d: %q %v\n", file.code, file.Clan.ClanNo, gameClan, err)
				return err
			}
			clanCache[gameClan] = clan
		}
		file.Clan = *clan
	}
	if verbose {
		log.Printf("import: cached %8d actors\n", len(actorCache))
		log.Printf("import: cached %8d clans\n", len(clanCache))
	}

	// upload all Turn report files
	for _, file := range turnReportFiles {
		actor := actorCache[authz.SysopId]
		// load the Turn report file from the file system
		contents, err := os.ReadFile(file.Path)
		if err != nil {
			log.Printf("sync: import: %s: %v\n", file.Path, err)
			return err
		}
		doc := &domains.Document{
			GameID:     file.Clan.GameID,
			ClanId:     file.Clan.ClanID,
			Turn:       file.Turn,
			ClanNo:     file.Clan.ClanNo,
			UnitId:     fmt.Sprintf("%4d", file.Clan.ClanNo),
			Path:       fmt.Sprintf("%s.%s.%04d.docx", file.code, file.Turn, file.Clan.ClanNo),
			Type:       domains.TurnReportFile,
			Contents:   contents,
			ModifiedAt: file.ModTime,
			CreatedAt:  file.ModTime,
			UpdatedAt:  file.ModTime,
		}
		documentId, err := s.documentsSvc.ReplaceDocument(actor, &file.Clan, doc, quiet, false, debug)
		if err != nil {
			log.Printf("sync: import: ReplaceDocument(%q): %v\n", file.Path, err)
			return err
		}
		if verbose {
			log.Printf("sync: import: %q %d\n", doc.Path, documentId)
			continue
		}
	}

	return nil
}

var (
	// report file Name is {game}.{turnNo}.{Clan}.docx
	reClanReportFile = regexp.MustCompile(`^(\d{4})\.(\d{4}-\d{2})\.(0\d{3})\.docx$`)

	// report file Name is {game}.{turnNo}.{unitId}.docx
	reUnitReportFile = regexp.MustCompile(`^(\d{4})\.(\d{4}-\d{2})\.(\d{4}([cefg][1-9])?)\.docx$`)
)

func (s *Service) ImportUsers(path string, handles ...string) error {
	// load the users data then force a validation for all the records,
	// even if we're not going to load them.
	usersData, err := jsondb.LoadUsers(path)
	if err != nil {
		log.Printf("import: %s: failed %v\n", path, err)
		return fmt.Errorf("import: users: %w", err)
	}
	for _, user := range usersData {
		if errHandle := domains.ValidateHandle(user.Handle); errHandle != nil {
			log.Printf("error: u %q: handle %q: %v\n", user.Handle, user.Handle, errHandle)
			err = errHandle
		}
		if errUsername := domains.ValidateUsername(user.UserName); errUsername != nil {
			log.Printf("error: u %q: userName %q: %v\n", user.Handle, user.UserName, errUsername)
			err = errUsername
		}
		if errEmail := domains.ValidateEmail(user.Email); errEmail != nil {
			log.Printf("error: u %q: email %q: %v\n", user.Handle, user.Email, errEmail)
			err = errEmail
		}
	}
	if err != nil {
		return domains.ErrBadInput
	}

	// if handles is empty, load all the users. otherwise, load only
	// the users in the list of handles.
	var handlesToLoad []string
	if len(handles) > 0 {
		for _, handle := range handles {
			handlesToLoad = append(handlesToLoad, handle)
		}
	} else {
		for _, user := range usersData {
			handlesToLoad = append(handlesToLoad, user.Handle)
		}
	}
	if len(handlesToLoad) == 0 {
		return nil
	}
	sort.Strings(handlesToLoad)

	// return an error if we're missing any users
	missingUsers := 0
	for _, handle := range handlesToLoad {
		_, ok := usersData[handle]
		if !ok {
			missingUsers++
			log.Printf("sync: import: users: %q: missing\n", handle)
		}
	}
	if missingUsers == 1 {
		return fmt.Errorf("sync: import: user not found")
	} else if missingUsers > 1 {
		return fmt.Errorf("sync: import: %d users not found", missingUsers)
	}

	ctx := s.db.Context()
	q := s.db.Queries()
	now := time.Now().UTC()
	createdAt, updatedAt := now.Unix(), now.Unix()

	for _, handle := range handlesToLoad {
		u := usersData[handle]

		id, err := q.CreateUser(ctx, sqlc.CreateUserParams{
			Handle:     u.Handle,
			Username:   u.UserName,
			Email:      u.Email,
			EmailOptIn: u.EmailOptIn,
			Timezone:   u.Tz.String(),
			IsActive:   u.Roles["active"],
			IsAdmin:    u.Roles["admin"],
			IsGm:       u.Roles["gm"],
			IsGuest:    u.Roles["guest"],
			IsPlayer:   u.Roles["player"],
			IsService:  u.Roles["service"],
			IsUser:     u.Roles["user"],
			CreatedAt:  createdAt,
			UpdatedAt:  updatedAt,
		})

		user, err := q.ReadUserByUserId(ctx, id)
		if err != nil {
			log.Printf("import: u %s: getUser failed %v\n", u.Handle, err)
			return fmt.Errorf("import: u %q: %w", u.Handle, err)
		}

		// create or update the user's password
		if u.Password.CreatePassword {
			u.Password.Password = phrases.Generate(6)
		}
		if u.Password.ChangePassword {
			u.Password.Password = phrases.Generate(6)
		}
		if errPassword := domains.ValidatePassword(u.Password.Password); errPassword != nil {
			log.Printf("error: u %q: password %q: %v\n", u.Handle, u.Password.Password, errPassword)
			return fmt.Errorf("import: u %q: %w", u.Handle, errPassword)
		}
		log.Printf("import: u %s: password.update %v\n", u.UserName, u.Password.UpdatePassword)
		_, err = s.authnSvc.UpdateCredentials(
			&domains.Actor{ID: authz.SysopId, Roles: domains.Roles{Sysop: true}},
			&domains.Actor{ID: domains.ID(user.UserID)},
			"",
			u.Password.Password)
		if err != nil {
			log.Printf("import: u %s: password %q: upsert %v\n", user.Handle, u.Password.Password, err)
			return err
		}
		fmt.Printf("%s: password %q\n", user.Handle, u.Password.Password)
	}

	return nil
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
	Name    string // {game}.{Turn}.{Clan}.docx
	Game    string // 0301
	Turn    string // YYYY-MM
	Clan    *domains.Clan
	UnitId  domains.UnitId
	Path    string    // full Path to report file
	ModTime time.Time // from os.Stat, assumed UTC?
}

// sort report files by game, Turn, Clan, and unitId
func sortReportFiles(files []*ReportFile) []*ReportFile {
	sort.Slice(files, func(i, j int) bool {
		a, b := files[i], files[j]
		if a.Game < b.Game {
			return true
		} else if a.Game > b.Game {
			return false
		}
		if a.Turn < b.Turn {
			return true
		} else if a.Turn > b.Turn {
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
