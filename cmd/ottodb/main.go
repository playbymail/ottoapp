// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package main implements ottodb, the database management tool
// for OttoMap.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/mdhender/phrases/v2"
	"github.com/playbymail/ottoapp"
	"github.com/playbymail/ottoapp/backend/domains"
	"github.com/playbymail/ottoapp/backend/iana"
	"github.com/playbymail/ottoapp/backend/stores/sqlite"
	"github.com/spf13/cobra"
)

func main() {
	log.SetFlags(log.Lshortfile)

	var cmdRoot = &cobra.Command{
		Use:   "ottodb",
		Short: "OttoMap database management tool",
		Long:  `OttoDB is a database management tool for OttoMap.`,
	}
	cmdRoot.CompletionOptions.DisableDefaultCmd = true
	cmdRoot.PersistentFlags().String("db", ".", "path to the database file")

	cmdRoot.AddCommand(cmdDb)

	cmdDb.AddCommand(cmdDbBackup)
	cmdDb.AddCommand(cmdDbCompact)
	cmdDb.AddCommand(cmdDbCreate)
	cmdDb.AddCommand(cmdDbInit)
	cmdDbInit.Flags().Bool("overwrite", false, "overwrite existing database")
	cmdDb.AddCommand(cmdDbMigrate)
	cmdDbMigrate.AddCommand(cmdDbMigrateUp)
	//cmdDb.AddCommand(cmdDbSeed)

	cmdRoot.AddCommand(cmdUser)
	cmdUser.AddCommand(cmdUserCreate)
	cmdUserCreate.Flags().String("email", "", "email address for user")
	cmdUserCreate.Flags().String("password", "", "password for user (generates random if not provided)")
	cmdUserCreate.Flags().String("tz", "UTC", "IANA timezone for user")
	cmdUser.AddCommand(cmdUserUpdate)
	cmdUserUpdate.Flags().Bool("active", true, "active flag user")
	cmdUserUpdate.Flags().String("email", "", "email address for user")
	cmdUserUpdate.Flags().String("password", "", "password for user (generates random if \":\")")
	cmdUserUpdate.Flags().String("tz", "", "IANA timezone for user")

	cmdRoot.AddCommand(cmdVersion)
	cmdVersion.Flags().Bool("build-info", false, "show build information")

	if err := cmdRoot.Execute(); err != nil {
		os.Exit(1)
	}
}

var cmdDb = &cobra.Command{
	Use:   "db",
	Short: "Database management commands",
	Long:  `Manage the OttoMap database including migrations and seeding.`,
}

var cmdDbBackup = &cobra.Command{
	Use:   "backup",
	Short: "Backup the database",
	Long:  `Backup the database.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := cmd.Flags().GetString("db")
		if err != nil {
			return err
		}
		name, err := sqlite.Backup(context.Background(), path)
		if err != nil {
			log.Fatalf("db: backup: %v\n", err)
		}
		log.Printf("db: %s: backup\n", name)
		return nil
	},
}

var cmdDbCompact = &cobra.Command{
	Use:   "compact",
	Short: "Compact the database",
	Long:  `Compact the database.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := cmd.Flags().GetString("db")
		if err != nil {
			return err
		}
		err = sqlite.Compact(context.Background(), path)
		if err != nil {
			log.Fatalf("db: compact: %v\n", err)
		}
		log.Printf("db: %s: compacted\n", path)
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
		path, err := cmd.Flags().GetString("db")
		if err != nil {
			return err
		}
		overwrite, err := cmd.Flags().GetBool("overwrite")
		if err != nil {
			return err
		}
		err = sqlite.Init(context.Background(), path, overwrite)
		if err != nil {
			log.Fatalf("db: init: %v\n", err)
		}
		log.Printf("db: %s: initialized\n", path)
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
		started := time.Now()
		path, err := cmd.Flags().GetString("db")
		if err != nil {
			return err
		}
		err = sqlite.MigrateUp(context.Background(), path, false)
		if err != nil {
			log.Fatalf("db: migrate: up: %v\n", err)
		}
		log.Printf("db: migrate: up: completed in %v\n", time.Since(started))
		return nil
	},
}

var cmdUser = &cobra.Command{
	Use:   "user",
	Short: "Manage data base records",
	Long:  `Commands to create, replace, update, and delete user records.`,
}

var cmdUserCreate = &cobra.Command{
	Use:   "create <handle>",
	Short: "Create a new user",
	Long:  `Create a new user with specified handle.`,
	Args:  cobra.ExactArgs(1), // require handle
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := cmd.Flags().GetString("db")
		if err != nil {
			return err
		}
		handle := strings.ToLower(args[0])
		if !sqlite.ValidateHandle(handle) {
			return domains.ErrInvalidHandle
		}
		emailSet, email := cmd.Flags().Changed("email"), handle+"@ottoapp"
		if emailSet {
			email, err = cmd.Flags().GetString("email")
			if err != nil {
				return err
			} else if !sqlite.ValidateEmail(email) {
				return domains.ErrInvalidEmail
			}
			email = strings.ToLower(email)
		}
		passwordSet, password := cmd.Flags().Changed("password"), phrases.Generate(6)
		if passwordSet {
			password, err = cmd.Flags().GetString("password")
			if err != nil {
				return err
			} else if !sqlite.ValidatePassword(password) {
				return domains.ErrInvalidPassword
			}
		}
		var loc *time.Location
		if tz, err := cmd.Flags().GetString("tz"); err != nil {
			return err
		} else if ctz, ok := iana.CanonicalName(tz); !ok {
			return fmt.Errorf("invalid time zone")
		} else if loc, err = time.LoadLocation(ctz); err != nil {
			return err
		}
		if loc == nil { // default to UTC
			if loc, err = time.LoadLocation("UTC"); err != nil {
				return err
			}
		}

		ctx := context.Background()
		db, err := sqlite.Open(ctx, path, true)
		if err != nil {
			log.Fatalf("db: open: %v\n", err)
		}
		defer func() {
			_ = db.Close()
		}()

		_, err = db.CreateUser(handle, email, password, loc)
		if err != nil {
			log.Fatalf("user %q: create: %v\n", handle, err)
		}

		log.Printf("user %q: email %q: tz %q: password %q: created\n", handle, email, loc.String(), password)

		return nil
	},
}

var cmdUserUpdate = &cobra.Command{
	Use:   "update <handle>",
	Short: "Update user record",
	Long:  `Update fields for a specific user. At least one update flag must be provided.`,
	Args:  cobra.ExactArgs(1), // require handle
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := cmd.Flags().GetString("db")
		if err != nil {
			return err
		}
		handle := strings.ToLower(args[0])
		if !sqlite.ValidateHandle(handle) {
			return domains.ErrInvalidHandle
		}
		var newEmail *string
		emailSet := cmd.Flags().Changed("email")
		if emailSet {
			value, err := cmd.Flags().GetString("email")
			if err != nil {
				return err
			} else if !sqlite.ValidateEmail(value) {
				return fmt.Errorf("invalid new email")
			}
			newEmail = &value
		}
		var newPassword *string
		passwordSet := cmd.Flags().Changed("password")
		if passwordSet {
			value, err := cmd.Flags().GetString("password")
			if err != nil {
				return err
			} else if value == "+" {
				value = phrases.Generate(6)
				log.Printf("generated random password: %q", value)
			}
			if !sqlite.ValidatePassword(value) {
				return fmt.Errorf("invalid new password")
			}
			newPassword = &value
		}
		var newTimeZone *time.Location
		tzSet := cmd.Flags().Changed("tz")
		if tzSet {
			if tz, err := cmd.Flags().GetString("tz"); err != nil {
				return err
			} else if ctz, ok := iana.CanonicalName(tz); !ok {
				return fmt.Errorf("invalid time zone")
			} else if newTimeZone, err = time.LoadLocation(ctz); err != nil {
				return err
			}
		}
		if !(emailSet || passwordSet || tzSet) {
			return fmt.Errorf("must update at least one field")
		}

		ctx := context.Background()
		db, err := sqlite.Open(ctx, path, true)
		if err != nil {
			log.Fatalf("db: open: %v\n", err)
		}
		defer func() {
			_ = db.Close()
		}()

		err = db.UpdateUser(handle, newEmail, newTimeZone)
		if err != nil {
			log.Fatalf("user: %q: update %v\n", handle, err)
		} else {
			if emailSet {
				log.Printf("user %q: email %q: updated", handle, *newEmail)
			}
			if tzSet {
				log.Printf("user %q: tz %q: updated", handle, newTimeZone.String())
			}
		}

		if newPassword != nil {
			err = db.UpdateUserPassword(handle, *newPassword)
			if err != nil {
				log.Fatalf("user %q: password %q: update %v\n", handle, *newPassword, err)
			}
			log.Printf("user %q: password %q: updated", handle, *newPassword)
		}

		return nil
	},
}

var cmdVersion = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `Display the current version of OttoApp.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if showBuildInfo, err := cmd.Flags().GetBool("build-info"); err != nil {
			return err
		} else if showBuildInfo {
			fmt.Println(ottoapp.Version().String())
		} else {
			fmt.Println(ottoapp.Version().Core())
		}
		return nil
	},
}
