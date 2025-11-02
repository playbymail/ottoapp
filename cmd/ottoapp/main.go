// Package main implements the backend server and CLI for OttoApp.
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
	"github.com/playbymail/ottoapp/backend/servers/rest"
	"github.com/playbymail/ottoapp/backend/stores/sqlite"
	"github.com/spf13/cobra"
)

func main() {
	log.SetFlags(log.Lshortfile)

	var cmdRoot = &cobra.Command{
		Use:   "ottoapp",
		Short: "OttoApp web server",
		Long:  `OttoApp is a web server for OttoMap.`,
	}
	cmdRoot.CompletionOptions.DisableDefaultCmd = true

	cmdRoot.AddCommand(cmdDb)
	cmdDb.PersistentFlags().String("db", ".", "path to the database file")

	cmdDb.AddCommand(cmdDbBackup)
	cmdDb.AddCommand(cmdDbCompact)
	cmdDb.AddCommand(cmdDbCreate)
	cmdDbCreate.AddCommand(cmdDbCreateUser)
	cmdDbCreateUser.Flags().String("email", "", "email address for user")
	cmdDbCreateUser.Flags().String("password", "", "password for user (generates random if not provided)")
	cmdDbCreateUser.Flags().String("tz", "UTC", "IANA timezone for user")
	cmdDb.AddCommand(cmdDbInit)
	cmdDbInit.Flags().Bool("overwrite", false, "overwrite existing database")
	cmdDb.AddCommand(cmdDbMigrate)
	cmdDbMigrate.AddCommand(cmdDbMigrateUp)
	//cmdDb.AddCommand(cmdDbSeed)
	cmdDb.AddCommand(cmdDbUpdate)
	cmdDbUpdate.AddCommand(cmdDbUpdateUser)
	cmdDbUpdateUser.Flags().Bool("active", true, "active flag user")
	cmdDbUpdateUser.Flags().String("email", "", "email address for user")
	cmdDbUpdateUser.Flags().String("password", "", "password for user (generates random if \":\")")
	cmdDbUpdateUser.Flags().String("tz", "", "IANA timezone for user")

	cmdRoot.AddCommand(cmdServe)
	cmdServe.PersistentFlags().Bool("csrf-guard", false, "enable csrf guards")
	cmdServe.PersistentFlags().String("db", ".", "path to the database file")
	cmdServe.PersistentFlags().Bool("debug", false, "enable debugging options")
	cmdServe.PersistentFlags().Bool("dev", false, "enable development mode")
	cmdServe.PersistentFlags().Bool("enable-catbird", false, "enable catbird testing")
	cmdServe.PersistentFlags().Bool("log-routes", false, "enable route logging")
	cmdServe.PersistentFlags().Duration("shutdown-delay", 30*time.Second, "delay for services to close during shutdown")
	cmdServe.PersistentFlags().Duration("shutdown-timer", 0, "timer to shut server down")

	cmdRoot.AddCommand(cmdVersion)
	cmdVersion.Flags().Bool("build-info", false, "show build information")

	if err := cmdRoot.Execute(); err != nil {
		os.Exit(1)
	}
}

var cmdDb = &cobra.Command{
	Use:   "db",
	Short: "Database management commands",
	Long:  `Manage the OttoApp database including migrations and seeding.`,
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

var cmdDbCreateUser = &cobra.Command{
	Use:   "user <handle>",
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
			log.Fatalf("db: create: user %q: %v\n", handle, err)
		}

		log.Printf("db: create: user %q: email %q: tz %q: password %q)", handle, email, loc.String(), password)

		return nil
	},
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

var cmdDbUpdate = &cobra.Command{
	Use:   "update",
	Short: "Update database records",
	Long:  `Update existing database records.`,
}

var cmdDbUpdateUser = &cobra.Command{
	Use:   "user <handle>",
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
			log.Fatalf("db: update: user %q: %v\n", handle, err)
		}
		if newPassword != nil {
			err = db.UpdateUserPassword(handle, *newPassword)
			if err != nil {
				log.Fatalf("db: update: user %q: password %v\n", handle, err)
			}
		}

		log.Printf("db: update: user %q: completed", handle)
		return nil
	},
}

var cmdServe = &cobra.Command{
	Use:   "serve",
	Short: "start the server",
	Long:  `Start the REST server.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := cmd.Flags().GetString("db")
		if err != nil {
			return err
		}

		var options []rest.Option
		if value, err := cmd.Flags().GetBool("enable-catbird"); err != nil {
			return err
		} else {
			options = append(options, rest.WithCatbird(value))
		}
		if value, err := cmd.Flags().GetBool("csrf-guard"); err != nil {
			return err
		} else {
			options = append(options, rest.WithCsrfGuard(value))
		}
		if value, err := cmd.Flags().GetBool("log-routes"); err != nil {
			return err
		} else {
			options = append(options, rest.WithRouteLogging(value))
		}
		if timer, err := cmd.Flags().GetDuration("shutdown-delay"); err != nil {
			return err
		} else if timer != 0 {
			options = append(options, rest.WithGrace(timer))
		}
		if timer, err := cmd.Flags().GetDuration("shutdown-timer"); err != nil {
			return err
		} else if timer != 0 {
			options = append(options, rest.WithTimer(timer))
		}

		log.Printf("[serve] db %q\n", path)

		ctx := context.Background()
		db, err := sqlite.Open(ctx, path, true)
		if err != nil {
			log.Fatalf("[serve] db: open: %v\n", err)
		}
		defer func() {
			log.Printf("[serve] db: close\n")
			_ = db.Close()
		}()

		s, err := rest.New(db, options...)
		if err != nil {
			_ = db.Close()
			log.Fatalf("[serve] rest: %v\n", err)
		}
		err = s.Run()
		if err != nil {
			_ = db.Close()
			log.Fatalf("[serve] rest: %v\n", err)
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
