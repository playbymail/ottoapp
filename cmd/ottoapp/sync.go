// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/playbymail/ottoapp"
	"github.com/playbymail/ottoapp/backend/make"
	"github.com/playbymail/ottoapp/backend/services/config"
	"github.com/playbymail/ottoapp/backend/services/sync"
	"github.com/playbymail/ottoapp/backend/stores/sqlite"
	"github.com/spf13/cobra"
)

func cmdSync() *cobra.Command {
	addFlags := func(cmd *cobra.Command) error {
		cmd.PersistentFlags().String("game", "0301", "game")
		cmd.PersistentFlags().Bool("show-timing", true, "time command")
		return nil
	}
	var cmd = &cobra.Command{
		Use:   "sync",
		Short: "Make things right",
	}
	if err := addFlags(cmd); err != nil {
		log.Fatal(err)
	}
	cmd.AddCommand(cmdSyncConfigFile())
	cmd.AddCommand(cmdSyncExport())
	cmd.AddCommand(cmdSyncImport())
	return cmd
}

func cmdSyncConfigFile() *cobra.Command {
	const checkVersion = true
	configFileName := filepath.Join("config", "configuration.json")
	addFlags := func(cmd *cobra.Command) error {
		cmd.Flags().StringVar(&configFileName, "input", configFileName, "name of config file to sync")
		return nil
	}
	cmd := &cobra.Command{
		Use:          "config-file <path-to-tribenet-data>",
		Short:        "update database from configuration file",
		Long:         "Sync the database with the configuration file.",
		SilenceUsage: true,
		Args:         cobra.ExactArgs(1), // require path to TN3.1 root
		RunE: func(cmd *cobra.Command, args []string) error {
			started := time.Now()
			path, err := filepath.Abs(args[0])
			if err != nil {
				return err
			}
			quiet, _ := cmd.Flags().GetBool("quiet")
			verbose, _ := cmd.Flags().GetBool("verbose")
			if quiet {
				verbose = false
			}

			dbPath, err := cmd.Flags().GetString("db")
			if err != nil {
				return err
			}
			debug, err := cmd.Flags().GetBool("debug")
			if err != nil {
				return err
			}
			ctx := context.Background()
			db, err := sqlite.Open(ctx, dbPath, checkVersion, quiet, verbose, debug)
			if err != nil {
				log.Fatalf("db: open: %v\n", err)
			}
			defer func() {
				_ = db.Close()
			}()
			log.Printf("%s: connected\n", dbPath)

			err = make.SyncConfigFile(db, path, configFileName, quiet, verbose, debug)
			if err != nil {
				return err
			}
			log.Printf("%s: synchronized to database\n", configFileName)

			if showTiming, _ := cmd.Flags().GetBool("show-timing"); showTiming {
				log.Printf("sync: config-file: completed in %v\n", time.Since(started))
			}
			return nil
		},
	}
	if err := addFlags(cmd); err != nil {
		log.Fatalf("%s: %v\n", cmd.Use, err)
	}
	return cmd
}

func cmdSyncExport() *cobra.Command {
	addFlags := func(cmd *cobra.Command) error {
		return nil
	}
	var cmd = &cobra.Command{
		Use:   "export",
		Short: "export files",
	}
	cmd.AddCommand(cmdSyncExportReportExtractFiles())
	if err := addFlags(cmd); err != nil {
		log.Fatal(err)
	}
	return cmd
}

func cmdSyncExportReportExtractFiles() *cobra.Command {
	addFlags := func(cmd *cobra.Command) error {
		return nil
	}
	var cmd = &cobra.Command{
		Use:          "report-extract-files",
		Short:        "export report extract files",
		Long:         "export report export files",
		SilenceUsage: true,
		Args:         cobra.ExactArgs(1), // require path to TN3.1 root
		RunE: func(cmd *cobra.Command, args []string) error {
			const checkVersion = true
			quiet, _ := cmd.Flags().GetBool("quiet")
			verbose, _ := cmd.Flags().GetBool("verbose")
			debug, _ := cmd.Flags().GetBool("debug")
			if quiet {
				verbose = false
			}

			started := time.Now()
			path, err := filepath.Abs(args[0])
			if err != nil {
				return err
			}

			dbPath, err := cmd.Flags().GetString("db")
			if err != nil {
				return err
			}
			ctx := context.Background()
			db, err := sqlite.Open(ctx, dbPath, checkVersion, quiet, verbose, debug)
			if err != nil {
				log.Fatalf("db: open: %v\n", err)
			}
			defer func() {
				_ = db.Close()
			}()
			log.Printf("%s: connected\n", dbPath)

			err = make.ExportExtractFiles(db, path, quiet, verbose, debug)
			if err != nil {
				return err
			}

			if showTiming, _ := cmd.Flags().GetBool("show-timing"); showTiming {
				log.Printf("export: report-extract-files: completed in %v\n", time.Since(started))
			}
			return nil
		},
	}
	if err := addFlags(cmd); err != nil {
		log.Fatal(err)
	}
	return cmd
}

func cmdSyncImport() *cobra.Command {
	addFlags := func(cmd *cobra.Command) error {
		return nil
	}
	var cmd = &cobra.Command{
		Use:   "import",
		Short: "import documents",
	}
	cmd.AddCommand(cmdSyncImportGames())
	cmd.AddCommand(cmdSyncImportMapFiles())
	cmd.AddCommand(cmdSyncImportReportFiles())
	cmd.AddCommand(cmdSyncImportReportExtractFiles())
	cmd.AddCommand(cmdSyncImportTurnReportFiles())
	cmd.AddCommand(cmdSyncImportUsers())
	if err := addFlags(cmd); err != nil {
		log.Fatal(err)
	}
	return cmd
}

func cmdSyncImportGames() *cobra.Command {
	var handles []string
	addFlags := func(cmd *cobra.Command) error {
		cmd.Flags().StringSliceVar(&handles, "handle", handles, "games to import")
		return nil
	}
	cmd := &cobra.Command{
		Use:          "games <path-to-games-file>",
		Short:        "update database with new (or changed) games",
		SilenceUsage: true,
		Args:         cobra.ExactArgs(1), // require path to file
		RunE: func(cmd *cobra.Command, args []string) error {
			const checkVersion = true
			quiet, _ := cmd.Flags().GetBool("quiet")
			verbose, _ := cmd.Flags().GetBool("verbose")
			debug, _ := cmd.Flags().GetBool("debug")
			if quiet {
				verbose = false
			}

			dbPath, err := cmd.Flags().GetString("db")
			if err != nil {
				return err
			}
			ctx := context.Background()
			db, err := sqlite.Open(ctx, dbPath, checkVersion, quiet, verbose, debug)
			if err != nil {
				log.Fatalf("db: open: %v\n", err)
			}
			defer func() {
				_ = db.Close()
			}()

			configSvc, err := config.New(db)
			if err != nil {
				return err
			}
			syncSvc, err := sync.New(db, nil, configSvc, nil)
			if err != nil {
				return err
			}

			path, err := filepath.Abs(args[0])
			if err != nil {
				return err
			}

			log.Printf("%s: importing games\n", path)
			err = syncSvc.ImportGames(path, handles...)
			return err
		},
	}
	if err := addFlags(cmd); err != nil {
		log.Fatalf("%s: %v\n", cmd.Use, err)
	}
	return cmd
}

func cmdSyncImportMapFiles() *cobra.Command {
	addFlags := func(cmd *cobra.Command) error {
		return nil
	}
	cmd := &cobra.Command{
		Use:          "map-files",
		Short:        "import map files",
		SilenceUsage: true,
		Args:         cobra.ExactArgs(1), // require path to TN3.1 root
		RunE: func(cmd *cobra.Command, args []string) error {
			const checkVersion = true
			quiet, _ := cmd.Flags().GetBool("quiet")
			verbose, _ := cmd.Flags().GetBool("verbose")
			debug, _ := cmd.Flags().GetBool("debug")
			if quiet {
				verbose = false
			}

			started := time.Now()
			path, err := filepath.Abs(args[0])
			if err != nil {
				return err
			}

			dbPath, err := cmd.Flags().GetString("db")
			if err != nil {
				return err
			}
			ctx := context.Background()
			db, err := sqlite.Open(ctx, dbPath, checkVersion, quiet, verbose, debug)
			if err != nil {
				log.Fatalf("db: open: %v\n", err)
			}
			defer func() {
				_ = db.Close()
			}()
			log.Printf("%s: connected\n", dbPath)

			err = make.ImportMapFiles(db, path, quiet, verbose, debug)
			if err != nil {
				return err
			}

			if showTiming, _ := cmd.Flags().GetBool("show-timing"); showTiming {
				log.Printf("import: map: completed in %v\n", time.Since(started))
			}
			return nil
		},
	}
	if err := addFlags(cmd); err != nil {
		log.Fatalf("%s: %v\n", cmd.Use, err)
	}
	return cmd
}

func cmdSyncImportReportFiles() *cobra.Command {
	addFlags := func(cmd *cobra.Command) error {
		return nil
	}
	cmd := &cobra.Command{
		Use:          "report-files <path-to-tribenet-data>",
		Short:        "update database with new (or changed) turn report files",
		Long:         "Sync the database with the turn report files on the file system.",
		SilenceUsage: true,
		Args:         cobra.ExactArgs(1), // require path to TN3.1 root
		RunE: func(cmd *cobra.Command, args []string) error {
			const checkVersion = true
			quiet, _ := cmd.Flags().GetBool("quiet")
			verbose, _ := cmd.Flags().GetBool("verbose")
			debug, _ := cmd.Flags().GetBool("debug")
			if quiet {
				verbose = false
			}

			started := time.Now()
			path, err := filepath.Abs(args[0])
			if err != nil {
				return err
			}

			dbPath, err := cmd.Flags().GetString("db")
			if err != nil {
				return err
			}
			ctx := context.Background()
			db, err := sqlite.Open(ctx, dbPath, checkVersion, quiet, verbose, debug)
			if err != nil {
				log.Fatalf("db: open: %v\n", err)
			}
			defer func() {
				_ = db.Close()
			}()

			err = make.SyncReportFiles(db, path, quiet, verbose, debug)
			if err != nil {
				return err
			}

			if showTiming, _ := cmd.Flags().GetBool("show-timing"); showTiming {
				log.Printf("sync: report-files: completed in %v\n", time.Since(started))
			}
			return nil
		},
	}
	if err := addFlags(cmd); err != nil {
		log.Fatalf("%s: %v\n", cmd.Use, err)
	}
	return cmd
}

func cmdSyncImportReportExtractFiles() *cobra.Command {
	addFlags := func(cmd *cobra.Command) error {
		return nil
	}
	cmd := &cobra.Command{
		Use:          "report-extract-files <path-to-tribenet-data>",
		Short:        "import report extract files",
		SilenceUsage: true,
		Args:         cobra.ExactArgs(1), // require path to TN3.1 root
		RunE: func(cmd *cobra.Command, args []string) error {
			const checkVersion = true
			quiet, _ := cmd.Flags().GetBool("quiet")
			verbose, _ := cmd.Flags().GetBool("verbose")
			debug, _ := cmd.Flags().GetBool("debug")
			if quiet {
				verbose = false
			}

			started := time.Now()
			path, err := filepath.Abs(args[0])
			if err != nil {
				return err
			}

			dbPath, err := cmd.Flags().GetString("db")
			if err != nil {
				return err
			}
			ctx := context.Background()
			db, err := sqlite.Open(ctx, dbPath, checkVersion, quiet, verbose, debug)
			if err != nil {
				log.Fatalf("db: open: %v\n", err)
			}
			defer func() {
				_ = db.Close()
			}()

			err = make.ImportReportExtractFiles(db, path, quiet, verbose, debug)
			if err != nil {
				return err
			}

			if showTiming, _ := cmd.Flags().GetBool("show-timing"); showTiming {
				log.Printf("sync: report-extract-files: completed in %v\n", time.Since(started))
			}
			return nil
		},
	}
	if err := addFlags(cmd); err != nil {
		log.Fatalf("%s: %v\n", cmd.Use, err)
	}
	return cmd
}

func cmdSyncImportTurnReportFiles() *cobra.Command {
	addFlags := func(cmd *cobra.Command) error {
		return nil
	}
	var cmd = &cobra.Command{
		Use:          "turn-report-files",
		Short:        "import turn report files",
		SilenceUsage: true,
		Args:         cobra.ExactArgs(1), // require path to TN3.1 root
		RunE: func(cmd *cobra.Command, args []string) error {
			const checkVersion = true
			quiet, _ := cmd.Flags().GetBool("quiet")
			verbose, _ := cmd.Flags().GetBool("verbose")
			debug, _ := cmd.Flags().GetBool("debug")
			if quiet {
				verbose = false
			}

			started := time.Now()
			path, err := filepath.Abs(args[0])
			if err != nil {
				return err
			}

			dbPath, err := cmd.Flags().GetString("db")
			if err != nil {
				return err
			}
			ctx := context.Background()
			db, err := sqlite.Open(ctx, dbPath, checkVersion, quiet, verbose, debug)
			if err != nil {
				log.Fatalf("db: open: %v\n", err)
			}
			defer func() {
				_ = db.Close()
			}()
			log.Printf("%s: connected\n", dbPath)

			err = make.ImportTurnReportFiles(db, path, quiet, verbose, debug)
			if err != nil {
				return err
			}

			if showTiming, _ := cmd.Flags().GetBool("show-timing"); showTiming {
				log.Printf("import: turn-report-files: completed in %v\n", time.Since(started))
			}
			return nil
		},
	}
	if err := addFlags(cmd); err != nil {
		log.Fatal(err)
	}
	return cmd
}

func cmdSyncImportUsers() *cobra.Command {
	var handles []string
	addFlags := func(cmd *cobra.Command) error {
		cmd.Flags().StringSliceVar(&handles, "handle", handles, "user handle to import")
		return nil
	}
	cmd := &cobra.Command{
		Use:          "users <path-to-users-file>",
		Short:        "update database with new (or changed) users",
		SilenceUsage: true,
		Args:         cobra.ExactArgs(1), // require path to file
		RunE: func(cmd *cobra.Command, args []string) error {
			const checkVersion = true
			quiet, _ := cmd.Flags().GetBool("quiet")
			verbose, _ := cmd.Flags().GetBool("verbose")
			debug, _ := cmd.Flags().GetBool("debug")
			if quiet {
				verbose = false
			}

			dbPath, err := cmd.Flags().GetString("db")
			if err != nil {
				return err
			}
			ctx := context.Background()
			db, err := sqlite.Open(ctx, dbPath, checkVersion, quiet, verbose, debug)
			if err != nil {
				log.Fatalf("db: open: %v\n", err)
			}
			defer func() {
				_ = db.Close()
			}()

			configSvc, err := config.New(db)
			if err != nil {
				return err
			}
			syncSvc, err := sync.New(db, nil, configSvc, nil)
			if err != nil {
				return err
			}

			path, err := filepath.Abs(args[0])
			if err != nil {
				return err
			}

			if !quiet {
				log.Printf("%s: importing users\n", path)
			}
			err = syncSvc.ImportUsers(path, handles...)
			return err
		},
	}
	if err := addFlags(cmd); err != nil {
		log.Fatalf("%s: %v\n", cmd.Use, err)
	}
	return cmd
}

func cmdVersion() *cobra.Command {
	showBuildInfo := false
	addFlags := func(cmd *cobra.Command) error {
		cmd.Flags().BoolVar(&showBuildInfo, "build-info", showBuildInfo, "show build information")
		return nil
	}
	var cmd = &cobra.Command{
		Use:   "version",
		Short: "display the application's version number",
		RunE: func(cmd *cobra.Command, args []string) error {
			if showBuildInfo {
				fmt.Println(ottoapp.Version().String())
				return nil
			}
			fmt.Println(ottoapp.Version().Core())
			return nil
		},
	}
	if err := addFlags(cmd); err != nil {
		log.Fatal(err)
	}
	return cmd
}
