// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package main implements the ottoapp command line tool.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/mdhender/phrases/v2"
	"github.com/mitchellh/go-homedir"
	"github.com/playbymail/ottoapp"
	"github.com/playbymail/ottoapp/backend/binder"
	"github.com/playbymail/ottoapp/backend/make"
	"github.com/playbymail/ottoapp/backend/stores/sqlite"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	log.SetFlags(log.Lshortfile)
	home, err := homedir.Dir()
	if err != nil {
		log.Fatalf("can't find home")
	}

	var cmdRoot = &cobra.Command{
		Use:   "ottoapp",
		Short: "OttoMap command runner",
		Long:  `OttoApp runs commands for OttoMap.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			ignoreConfigFile, _ := cmd.Flags().GetBool("ignore-config-file")
			if ignoreConfigFile {
				return nil
			}
			debugConfig := false
			if cmd.Flags().Changed("debug-config") {
				debugConfig = true
			}
			err := binder.Bind(cmd, binder.Options{
				EnvPrefix: "OTTOAPP",
				//ConfigFile:  ".",
				ConfigName:  "ottoapp",
				ConfigPaths: []string{".", home},
				ConfigType:  "json",
				DebugConfig: debugConfig,
			})
			if cmd.Flags().Changed("dump-config") {
				fmt.Println("---- Effective Config ----")
				data, _ := json.MarshalIndent(viper.GetViper().AllSettings(), "", "  ")
				fmt.Println(string(data))
				fmt.Println("--------------------------")
			}
			if cmd.Flags().Changed("dump-resolved-config") {
				resolved := binder.DumpResolved(cmd, viper.GetViper())
				fmt.Println("---- Resolved  Config ----")
				b, _ := json.MarshalIndent(resolved, "", "  ")
				fmt.Println(string(b))
				fmt.Println("--------------------------")
			}
			return err
		},
	}
	cmdRoot.CompletionOptions.DisableDefaultCmd = true
	cmdRoot.PersistentFlags().String("db", ".", "path to the database file")
	cmdRoot.PersistentFlags().Bool("debug", false, "enable debugging options")
	cmdRoot.PersistentFlags().Bool("dev", false, "enable development mode")
	cmdRoot.PersistentFlags().Bool("debug-config", false, "show config binding sources")
	cmdRoot.PersistentFlags().Bool("dump-config", false, "dump config after binding")
	cmdRoot.PersistentFlags().Bool("dump-resolved-config", false, "dump resolved config after binding")
	cmdRoot.PersistentFlags().BoolP("ignore-config-file", "N", false, "ignore ottoapp.json file")
	cmdRoot.PersistentFlags().Bool("quiet", false, "log less information")
	cmdRoot.PersistentFlags().Bool("verbose", false, "log more information")

	var cmdApi = &cobra.Command{
		Use:   "api",
		Short: "API server commands",
	}
	cmdRoot.AddCommand(cmdApi)
	cmdApi.AddCommand(cmdApiServe)
	cmdApiServe.Flags().Bool("csrf-guard", false, "enable csrf guards")
	cmdApiServe.Flags().String("host", "localhost", "change the bind network")
	cmdApiServe.Flags().Bool("log-routes", false, "enable route logging")
	cmdApiServe.Flags().String("port", "8181", "change the bind port")
	cmdApiServe.Flags().Duration("sessions-reap-interval", 15*time.Minute, "interval to remove expired sessions")
	cmdApiServe.Flags().Duration("sessions-ttl", 24*time.Hour, "session duration")
	cmdApiServe.Flags().Duration("shutdown-delay", 30*time.Second, "delay for services to close during shutdown")
	cmdApiServe.Flags().String("shutdown-key", "", "api key authorizing shutdown")
	cmdApiServe.Flags().Duration("shutdown-timer", 0, "timer to shut server down")
	cmdApiServe.Flags().String("userdata", "userdata", "path to user data")

	var cmdApp = &cobra.Command{
		Use:   "app",
		Short: "application management commands",
	}
	cmdRoot.AddCommand(cmdApp)
	cmdApp.AddCommand(cmdAppVersion)
	cmdAppVersion.Flags().Bool("show-build-info", false, "show build information")
	cmdApp.AddCommand(cmdAppTestUserProfile)
	cmdAppTestUserProfile.Flags().String("host", "localhost", "API server host")
	cmdAppTestUserProfile.Flags().String("port", "8181", "API server port")
	cmdAppTestUserProfile.Flags().String("email", "ottoapp@ottoapp.example.com", "email for authentication")
	cmdAppTestUserProfile.Flags().String("password", "", "password for authentication (required)")

	var cmdDb = &cobra.Command{
		Use:   "db",
		Short: "Database management commands",
		Long:  `Manage the OttoMap database including migrations and seeding.`,
	}
	cmdRoot.AddCommand(cmdDb)
	cmdDb.AddCommand(cmdDbBackup)
	cmdDbBackup.Flags().String("output", "", "path to directory for backup file (must exist)")
	cmdDb.AddCommand(cmdDbClone)
	cmdDb.AddCommand(cmdDbCompact)
	cmdDb.AddCommand(cmdDbCreate)
	cmdDb.AddCommand(cmdDbInit)
	cmdDbInit.Flags().Bool("overwrite", false, "overwrite existing database")
	cmdDb.AddCommand(cmdDbMigrate)
	cmdDbMigrate.AddCommand(cmdDbMigrateStatus)
	cmdDbMigrate.AddCommand(cmdDbMigrateUp)
	//cmdDb.AddCommand(cmdDbSeed)
	cmdDb.AddCommand(cmdDbVersion)

	cmdRoot.AddCommand(cmdGame())

	cmdRoot.AddCommand(cmdGenerate())

	cmdRoot.AddCommand(cmdPhrase())

	var cmdReport = &cobra.Command{
		Use:   "report",
		Short: "report management",
	}
	cmdRoot.AddCommand(cmdReport)
	cmdReport.AddCommand(cmdReportExtract)
	cmdReportExtract.Flags().String("output", "report.txt", "file to create")
	cmdReport.AddCommand(cmdReportParse)
	cmdReportParse.Flags().Bool("docxml-only", false, "parse to DocXML only")

	cmdRoot.AddCommand(cmdRun())

	cmdRoot.AddCommand(cmdSync())

	var cmdUser = &cobra.Command{
		Use:   "user",
		Short: "user management",
		Long:  `Commands to create, replace, update, and delete user records.`,
	}
	cmdRoot.AddCommand(cmdUser)
	cmdUser.AddCommand(cmdUserCreate)
	cmdUserCreate.Flags().String("email", "", "email address for user")
	cmdUserCreate.Flags().String("password", "", "password for user (generates random if not provided)")
	cmdUserCreate.Flags().String("tz", "UTC", "IANA timezone for user")
	cmdUserCreate.Flags().String("username", "", "user name")
	cmdUser.AddCommand(cmdUserUpdate)
	cmdUserUpdate.Flags().Bool("active", true, "active flag user")
	cmdUserUpdate.Flags().String("email", "", "email address for user")
	cmdUserUpdate.Flags().String("username", "", "user name")
	cmdUserUpdate.Flags().String("password", "", "password for user (generates random if \":\")")
	cmdUserUpdate.Flags().String("tz", "", "IANA timezone for user")
	cmdUser.AddCommand(cmdUserRole)
	cmdUserRole.Flags().StringSlice("add", []string{}, "roles to add (comma-separated: user,admin,player)")
	cmdUserRole.Flags().StringSlice("remove", []string{}, "roles to remove (comma-separated: user,admin,player)")

	cmdRoot.AddCommand(cmdVersion())

	if err := cmdRoot.Execute(); err != nil {
		os.Exit(1)
	}
}

func cmdGenerate() *cobra.Command {
	addFlags := func(cmd *cobra.Command) error {
		return nil
	}
	var cmd = &cobra.Command{
		Use:   "generate",
		Short: "generate artifacts",
	}
	cmd.AddCommand(cmdGenerateMakefile())
	if err := addFlags(cmd); err != nil {
		log.Fatal(err)
	}
	return cmd
}

func cmdGenerateMakefile() *cobra.Command {
	gameId := "0301"
	makefileName := "maps.mk"
	addFlags := func(cmd *cobra.Command) error {
		cmd.Flags().StringVar(&gameId, "game", gameId, "game to create makefile for")
		cmd.Flags().StringVar(&makefileName, "output", makefileName, "name of makefile to create")
		return nil
	}
	var cmd = &cobra.Command{
		Use:          "makefile",
		Short:        "Generate a Makefile for map generation",
		Long:         `Scans the data directory for turn reports and generates a Makefile to build maps.`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			quiet, _ := cmd.Flags().GetBool("quiet")
			verbose, _ := cmd.Flags().GetBool("verbose")
			debug, _ := cmd.Flags().GetBool("debug")
			if quiet {
				verbose = false
			}
			return make.Makefile(makefileName, gameId, quiet, verbose, debug)
		},
	}
	if err := addFlags(cmd); err != nil {
		log.Fatal(err)
	}
	return cmd
}

func cmdPhrase() *cobra.Command {
	length := 6
	addFlags := func(cmd *cobra.Command) error {
		cmd.Flags().IntVar(&length, "length", length, "number of words in phrase")
		return nil
	}
	var cmd = &cobra.Command{
		Use:   "phrase",
		Short: "random phrase",
		RunE: func(cmd *cobra.Command, args []string) error {
			if length < 1 {
				length = 1
			} else if length > 16 {
				length = 16
			}
			fmt.Println(phrases.Generate(length))
			return nil
		},
	}
	if err := addFlags(cmd); err != nil {
		log.Fatal(err)
	}
	return cmd
}

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
	cmd.AddCommand(cmdSyncImportMapFiles())
	cmd.AddCommand(cmdSyncImportReportFiles())
	cmd.AddCommand(cmdSyncImportReportExtractFiles())
	cmd.AddCommand(cmdSyncImportTurnReportFiles())
	if err := addFlags(cmd); err != nil {
		log.Fatal(err)
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
