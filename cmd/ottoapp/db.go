// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/playbymail/ottoapp/backend/stores/sqlite"
	"github.com/spf13/cobra"
)

var cmdDbBackup = &cobra.Command{
	Use:          "backup",
	Short:        "Backup the database",
	Long:         `Backup the database.`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		quiet, _ := cmd.Flags().GetBool("quiet")
		verbose, _ := cmd.Flags().GetBool("verbose")
		debug, _ := cmd.Flags().GetBool("debug")
		if quiet {
			verbose = false
		}

		path, err := cmd.Flags().GetString("db")
		if err != nil {
			return err
		}
		outputPath, err := cmd.Flags().GetString("output")
		if err != nil {
			return err
		}

		name, err := sqlite.Backup(context.Background(), path, outputPath, quiet, verbose, debug)
		if err != nil {
			return fmt.Errorf("backup failed: %w", err)
		}
		log.Printf("db: %s: backup\n", name)
		return nil
	},
}

var cmdDbClone = &cobra.Command{
	Use:          "clone <output-directory>",
	Short:        "Clone the database for testing",
	Long:         `Clone the database to a working copy for testing. Creates ottoapp.db in the output directory.`,
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		quiet, _ := cmd.Flags().GetBool("quiet")
		verbose, _ := cmd.Flags().GetBool("verbose")
		debug, _ := cmd.Flags().GetBool("debug")
		if quiet {
			verbose = false
		}

		path, err := cmd.Flags().GetString("db")
		if err != nil {
			return err
		}
		outputPath := args[0]

		clonePath, err := sqlite.Clone(context.Background(), path, outputPath, quiet, verbose, debug)
		if err != nil {
			return fmt.Errorf("clone failed: %w", err)
		}
		log.Printf("db: %s: cloned\n", clonePath)
		return nil
	},
}

var cmdDbCompact = &cobra.Command{
	Use:   "compact",
	Short: "Compact the database",
	Long:  `Compact the database.`,
	RunE: func(cmd *cobra.Command, args []string) error {
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
		err = sqlite.Compact(context.Background(), dbPath, quiet, verbose, debug)
		if err != nil {
			log.Fatalf("db: compact: %v\n", err)
		}
		log.Printf("db: %s: compacted\n", dbPath)
		return nil
	},
}

var cmdDbCreate = &cobra.Command{
	Use:   "create",
	Short: "Create data base records",
	Long:  `Create new database records.`,
}

var cmdDbInit = &cobra.Command{
	Use:   "init",
	Short: "Initialize the database",
	Long:  `Create the database file if it doesn't exist.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		dbPath, err := cmd.Flags().GetString("db")
		if err != nil {
			return err
		}
		overwrite, err := cmd.Flags().GetBool("overwrite")
		if err != nil {
			return err
		}
		err = sqlite.Init(context.Background(), dbPath, overwrite)
		if err != nil {
			log.Fatalf("db: init: %v\n", err)
		}
		log.Printf("db: %s: initialized\n", dbPath)
		return nil
	},
}

var cmdDbMigrate = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Long:  `Apply schema migrations to the database.`,
}

var cmdDbMigrateUp = &cobra.Command{
	Use:   "up",
	Short: "Run database migration up",
	Long:  `Apply schema migrations to upgrade the database.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		const checkVersion, isInitializing = false, false
		quiet, _ := cmd.Flags().GetBool("quiet")
		verbose, _ := cmd.Flags().GetBool("verbose")
		debug, _ := cmd.Flags().GetBool("debug")
		if quiet {
			verbose = false
		}

		started := time.Now()
		dbPath, err := cmd.Flags().GetString("db")
		if err != nil {
			return err
		}
		err = sqlite.MigrateUp(context.Background(), dbPath, isInitializing, quiet, verbose, debug)
		if err != nil {
			log.Fatalf("db: migrate: up: %v\n", err)
		}
		log.Printf("db: migrate: up: completed in %v\n", time.Since(started))
		return nil
	},
}

var cmdDbMigrateStatus = &cobra.Command{
	Use:   "status",
	Short: "show status of migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		const checkVersion = false
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
		rows, err := db.GetDatabaseMigrationStatus()
		if err != nil {
			log.Fatalf("db: migrate: status %v\n", err)
		}
		fmt.Printf("ID__ Migration_ID_ C Applied________________ File_Name_______________________________\n")
		for _, row := range rows {
			var id string
			if row.Id != 0 {
				id = fmt.Sprintf("%4d", row.Id)
			}
			cf := " "
			if row.IsCurrent {
				cf = "*"
			}
			var appliedAt string
			if !row.AppliedAt.IsZero() {
				appliedAt = row.AppliedAt.Format(time.DateTime) + " UTC"
			}
			fmt.Printf("%s %s %-11s %23s %s\n", id, cf, row.MigrationId, appliedAt, row.FileName)
		}
		return nil
	},
}

var cmdDbVersion = &cobra.Command{
	Use:   "version",
	Short: "show database version",
	RunE: func(cmd *cobra.Command, args []string) error {
		const checkVersion = false
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
		version, err := db.GetDatabaseVersion()
		if err != nil {
			log.Fatalf("db: version %v\n", err)
		}
		fmt.Printf("%s\n", version)
		return nil
	},
}
