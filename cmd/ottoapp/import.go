// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package main

import (
	"context"
	"log"
	"path/filepath"

	"github.com/playbymail/ottoapp/backend/services/config"
	"github.com/playbymail/ottoapp/backend/services/sync"
	"github.com/playbymail/ottoapp/backend/stores/sqlite"
	"github.com/spf13/cobra"
)

func cmdImport() *cobra.Command {
	addFlags := func(cmd *cobra.Command) error {
		return nil
	}
	var cmd = &cobra.Command{
		Use:   "import",
		Short: "import data from the file system",
	}
	cmd.AddCommand(cmdSyncImportGames())
	cmd.AddCommand(cmdSyncImportMapFiles())
	cmd.AddCommand(cmdImportOttoAppConfig())
	cmd.AddCommand(cmdSyncImportReportFiles())
	cmd.AddCommand(cmdSyncImportReportExtractFiles())
	cmd.AddCommand(cmdSyncImportTurnReportFiles())
	cmd.AddCommand(cmdSyncImportUsers())
	if err := addFlags(cmd); err != nil {
		log.Fatal(err)
	}
	return cmd
}

func cmdImportOttoAppConfig() *cobra.Command {
	addFlags := func(cmd *cobra.Command) error {
		return nil
	}
	cmd := &cobra.Command{
		Use:          "ottoapp-config-file <path-to-config-file>",
		Short:        "update database with new (or changed) configuration",
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
			if verbose {
				log.Printf("%s: importing configuration\n", path)
			}

			return syncSvc.ImportOttoAppConfig(path, quiet, verbose, debug)
		},
	}
	if err := addFlags(cmd); err != nil {
		log.Fatalf("%s: %v\n", cmd.Use, err)
	}
	return cmd
}
