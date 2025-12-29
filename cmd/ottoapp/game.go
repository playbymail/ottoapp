// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package main

import (
	_ "embed"
	"log"

	"github.com/spf13/cobra"
)

//go:embed testdata/memdb-players.json
var memdbPlayersJsonData []byte

func cmdGame() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "game",
		Short:        "game management",
		SilenceUsage: true,
	}

	cmd.AddCommand(cmdGameImport())

	cmd.AddCommand(cmdGameUpload)
	cmdGameUpload.Flags().Bool("can-delete", true, "delete flag")
	cmdGameUpload.Flags().Bool("can-read", true, "read flag")
	cmdGameUpload.Flags().Bool("can-share", true, "share flag")
	cmdGameUpload.Flags().Bool("can-write", true, "write flag")
	cmdGameUpload.Flags().Int("clan", 0, "clan to assign ownership to")
	cmdGameUpload.Flags().String("game", "0301", "game to upload to")
	if err := cmdGameUpload.MarkFlagRequired("game"); err != nil {
		log.Fatalf("game: upload: game: MarkFlagRequired %v\n", err)
	}
	cmdGameUpload.Flags().String("handle", "", "handle to assign ownership to")
	cmdGameUpload.Flags().String("name", "", "overwrite the file name after uploading")
	cmdGameUpload.MarkFlagsMutuallyExclusive("clan", "handle")

	return cmd
}

func cmdGameImport() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "import <path>",
		Short:        "import game data from JSON data file",
		SilenceUsage: true,
		Args:         cobra.ExactArgs(1), // require path
		RunE: func(cmd *cobra.Command, args []string) error {
			panic("obsolete: replace with sync.Service")
			//const checkVersion = true
			//quiet, _ := cmd.Flags().GetBool("quiet")
			//verbose, _ := cmd.Flags().GetBool("verbose")
			//debug, _ := cmd.Flags().GetBool("debug")
			//if quiet {
			//	verbose = false
			//}
			//
			//started := time.Now()
			//path := args[0]
			//if sb, err := os.Stat(path); err != nil {
			//	if errors.Is(err, os.ErrNotExist) {
			//		return errors.Join(domains.ErrInvalidPath, domains.ErrNotExists)
			//	}
			//	return errors.Join(domains.ErrInvalidPath, err)
			//} else if sb.IsDir() || !sb.Mode().IsRegular() {
			//	return errors.Join(domains.ErrInvalidPath, domains.ErrNotFile)
			//}
			//input, err := os.ReadFile(path)
			//if err != nil {
			//	return err
			//}
			//var data games.ImportFile
			//err = json.Unmarshal(input, &data)
			//if err != nil {
			//	return err
			//}
			//
			//dbPath, err := cmd.Flags().GetString("db")
			//if err != nil {
			//	return err
			//}
			//ctx := context.Background()
			//db, err := sqlite.Open(ctx, dbPath, checkVersion, quiet, verbose, debug)
			//if err != nil {
			//	log.Fatalf("db: open: %v\n", err)
			//}
			//defer func() {
			//	_ = db.Close()
			//}()
			//
			//gameSvc, err := games.New(db, nil, nil, nil)
			//if err != nil {
			//	return fmt.Errorf("games: new service %w", err)
			//}
			//err = gameSvc.Import(&data)
			//if err != nil {
			//	return fmt.Errorf("games: import %w", err)
			//}
			//
			//fmt.Printf("%s: imported in %v\n", path, time.Since(started))
			//return nil
		},
	}

	return cmd
}

var cmdGameUpload = &cobra.Command{
	Use:   "upload <document>",
	Short: "Upload a new game document (report, extract, or map)",
	Long:  `Upload game documents to the server.`,
	Args:  cobra.ExactArgs(1), // require path to document to upload
	RunE: func(cmd *cobra.Command, args []string) error {
		panic("obsolete: replace with sync.Service")
		//const checkVersion = true
		//quiet, _ := cmd.Flags().GetBool("quiet")
		//verbose, _ := cmd.Flags().GetBool("verbose")
		//debug, _ := cmd.Flags().GetBool("debug")
		//if quiet {
		//	verbose = false
		//}
		//
		//startedAt := time.Now()
		//
		//dbPath, err := cmd.Flags().GetString("db")
		//if err != nil {
		//	return err
		//}
		//path := args[0]
		////log.Printf("game: upload: path %q\n", path)
		//var mimeType domains.MimeType
		//ext := strings.ToLower(filepath.Ext(path))
		//switch ext {
		//case ".docx":
		//	mimeType = domains.DOCXMimeType
		//case ".txt":
		//	mimeType = domains.ReportMimeType
		//case ".wxx":
		//	mimeType = domains.WXXMimeType
		//default:
		//	return fmt.Errorf("unknown file type %q", ext)
		//}
		//var name string
		//if cmd.Flags().Changed("name") {
		//	if value, err := cmd.Flags().GetString("name"); err != nil {
		//		return err
		//	} else {
		//		name = value
		//	}
		//} else {
		//	name = path
		//}
		//name = filepath.Base(filepath.Clean(name))
		////log.Printf("game: upload: name %q\n", name)
		//gameId, err := cmd.Flags().GetString("game")
		//if err != nil {
		//	return err
		//}
		////log.Printf("game: upload: game %q\n", gameId)
		//
		//// enforce oneOf for clan and handle
		//var clanNo int
		//var handle string
		//clanSet := cmd.Flags().Changed("clan")
		//handleSet := cmd.Flags().Changed("handle")
		//if (clanSet && handleSet) || !(clanSet || handleSet) {
		//	return errors.New("must specify either --clan or --handle, but not both")
		//} else if clanSet {
		//	clanNo, err = cmd.Flags().GetInt("clan")
		//	if err != nil {
		//		return err
		//	} else if !(0 < clanNo && clanNo <= 999) {
		//		return fmt.Errorf("clan %d: must be between 1 and 999", err)
		//	}
		//	//log.Printf("game: upload: owner: clan %d\n", clanNo)
		//} else {
		//	handle, err = cmd.Flags().GetString("handle")
		//	if err != nil {
		//		return err
		//	}
		//	//log.Printf("game: upload: owner: handle %q\n", handle)
		//}
		//
		//ctx := context.Background()
		//db, err := sqlite.Open(ctx, dbPath, checkVersion, quiet, verbose, debug)
		//if err != nil {
		//	return errors.Join(fmt.Errorf("game: db open failed"), err)
		//}
		//defer func() {
		//	_ = db.Close()
		//}()
		//
		//authzSvc := authz.New(db)
		//authnSvc := authn.New(db, authzSvc)
		//ianaSvc, err := iana.New(db)
		//if err != nil {
		//	return err
		//}
		//usersSvc := users.New(db, authnSvc, authzSvc, ianaSvc)
		//docSvc, err := documents.New(db, authzSvc, usersSvc)
		//if err != nil {
		//	return errors.Join(fmt.Errorf("sessions.new"), err)
		//}
		//gamesSvc, err := games.New(db, authnSvc, authzSvc, usersSvc)
		//if err != nil {
		//	return err
		//}
		//
		//actor := &domains.Actor{ID: authz.SysopId, Sysop: true}
		//var clan *domains.Clan
		//if clanSet {
		//	clan, err = gamesSvc.ReadClan(domains.GameID(gameId), clanNo)
		//	if err != nil {
		//		return errors.Join(fmt.Errorf("gameId %q: clanNo %d: invalid", gameId, clanNo), err)
		//	}
		//} else {
		//	owner, err := authzSvc.GetActorByHandle(handle)
		//	if err != nil {
		//		return errors.Join(fmt.Errorf("handle %q: invalid", handle), err)
		//	}
		//	clan, err = gamesSvc.ReadClanByGameIdAndUserId(domains.GameID(gameId), owner.ID)
		//	if err != nil {
		//		return errors.Join(fmt.Errorf("handle %q: invalid", handle), err)
		//	}
		//}
		////log.Printf("game: upload: clan %d: (%q, %d, %d) d:%v r:%v s:%v w:%v\n", clan.ClanID, clan.GameID, clan.UserID, clan.ClanNo, canDelete, canRead, canShare, canWrite)
		//
		//var docId domains.ID
		//switch mimeType {
		//case domains.DOCXMimeType:
		//	docId, err = docSvc.LoadDocxFromFS(actor, clan, path, name, quiet, verbose, debug)
		//	if err != nil {
		//		return errors.Join(fmt.Errorf("%q", path), err)
		//	}
		//case domains.ReportMimeType:
		//	docId, err = docSvc.LoadReportFromFS(actor, clan, path, name, quiet, verbose, debug)
		//	if err != nil {
		//		return errors.Join(fmt.Errorf("%q", path), err)
		//	}
		//case domains.WXXMimeType:
		//	docId, err = docSvc.LoadMapFromFS(actor, clan, path, name, quiet, verbose, debug)
		//	if err != nil {
		//		return errors.Join(fmt.Errorf("%q", path), err)
		//	}
		//default:
		//	panic("!implemented")
		//}
		//
		//log.Printf("game: upload: docId %d: completed in %v\n", docId, time.Since(startedAt))
		//return nil
	},
}
