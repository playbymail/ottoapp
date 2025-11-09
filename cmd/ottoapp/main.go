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
	"strings"
	"time"

	"github.com/mdhender/phrases/v2"
	"github.com/mitchellh/go-homedir"
	"github.com/playbymail/ottoapp"
	"github.com/playbymail/ottoapp/backend/auth"
	"github.com/playbymail/ottoapp/backend/binder"
	"github.com/playbymail/ottoapp/backend/documents"
	"github.com/playbymail/ottoapp/backend/domains"
	"github.com/playbymail/ottoapp/backend/iana"
	"github.com/playbymail/ottoapp/backend/runners"
	"github.com/playbymail/ottoapp/backend/servers/rest"
	"github.com/playbymail/ottoapp/backend/sessions"
	"github.com/playbymail/ottoapp/backend/stores/sqlite"
	"github.com/playbymail/ottoapp/backend/users"
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

	var cmdApp = &cobra.Command{
		Use:   "app",
		Short: "application management commands",
	}
	cmdRoot.AddCommand(cmdApp)
	cmdApp.AddCommand(cmdAppVersion)
	cmdAppVersion.Flags().Bool("show-build-info", false, "show build information")

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

	var cmdReport = &cobra.Command{
		Use:   "report",
		Short: "report management",
	}
	cmdRoot.AddCommand(cmdReport)
	cmdReport.AddCommand(cmdReportParse)
	cmdReport.AddCommand(cmdReportUpload)
	cmdReportUpload.Flags().String("name", "", "overwrite the file name after uploading")
	cmdReportUpload.Flags().String("owner", "sysop", "user to assign ownership to")

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
	cmdUser.AddCommand(cmdUserUpdate)
	cmdUserUpdate.Flags().Bool("active", true, "active flag user")
	cmdUserUpdate.Flags().String("email", "", "email address for user")
	cmdUserUpdate.Flags().String("password", "", "password for user (generates random if \":\")")
	cmdUserUpdate.Flags().String("tz", "", "IANA timezone for user")

	var cmdVersion = &cobra.Command{
		Use:   "version",
		Short: "display the application's version number",
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
	cmdVersion.Flags().Bool("build-info", false, "show build information")
	cmdRoot.AddCommand(cmdVersion)

	if err := cmdRoot.Execute(); err != nil {
		os.Exit(1)
	}
}

var cmdApiServe = &cobra.Command{
	Use:   "serve",
	Short: "start the API server",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := cmd.Flags().GetString("db")
		if err != nil {
			return err
		}

		var options []rest.Option
		if value, err := cmd.Flags().GetBool("csrf-guard"); err != nil {
			return err
		} else {
			options = append(options, rest.WithCsrfGuard(value))
		}
		if value, err := cmd.Flags().GetString("host"); err != nil {
			return err
		} else {
			options = append(options, rest.WithHost(value))
		}
		if value, err := cmd.Flags().GetBool("log-routes"); err != nil {
			return err
		} else {
			options = append(options, rest.WithRouteLogging(value))
		}
		if value, err := cmd.Flags().GetString("port"); err != nil {
			return err
		} else {
			options = append(options, rest.WithPort(value))
		}
		if value, err := cmd.Flags().GetDuration("shutdown-delay"); err != nil {
			return err
		} else if value != 0 {
			options = append(options, rest.WithGrace(value))
		}
		if cmd.Flags().Changed("shutdown-key") {
			if value, err := cmd.Flags().GetString("shutdown-key"); err != nil {
				return err
			} else {
				options = append(options, rest.WithShutdownKey(value))
			}
		}
		if value, err := cmd.Flags().GetDuration("shutdown-timer"); err != nil {
			return err
		} else if value != 0 {
			options = append(options, rest.WithTimer(value))
		}

		log.Printf("[serve] db %q\n", path)
		ctx := context.Background()
		var db *sqlite.DB
		if path == ":memory:" {
			// server has the ability to use a temporary database for testing.
			db, err = sqlite.OpenTempDB(ctx)
		} else {
			db, err = sqlite.Open(ctx, path, true, false)
		}
		if err != nil {
			log.Fatalf("[serve] db: open: %v\n", err)
		}
		defer func() {
			log.Printf("[serve] db: close\n")
			_ = db.Close()
		}()

		authSvc := auth.New(db) // uses sqlite + domains
		tzSvc, err := iana.New(db)
		usersSvc := users.New(db, authSvc, tzSvc) // uses sqlite + domains

		sessionsSvc, err := sessions.New(db, authSvc, usersSvc, 24*time.Hour, 15*time.Minute)
		if err != nil {
			_ = db.Close()
			log.Fatalf("[serve] sessionManager: %v\n", err)
		}

		options = append(options, rest.WithIanaService(tzSvc))
		options = append(options, rest.WithUsersService(usersSvc))
		s, err := rest.New(sessionsSvc, options...)
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

var cmdAppVersion = &cobra.Command{
	Use:   "version",
	Short: "display the application's version number",
	RunE: func(cmd *cobra.Command, args []string) error {
		showBuildInfo, err := cmd.Flags().GetBool("show-build-info")
		if err != nil {
			return err
		}
		r := runners.New("http", "127.0.0.1", "8181")
		err = r.GetVersion(showBuildInfo)
		if err != nil {
			log.Fatalf("%v\n", err)
		}
		return nil
	},
}

var cmdDbBackup = &cobra.Command{
	Use:          "backup",
	Short:        "Backup the database",
	Long:         `Backup the database.`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := cmd.Flags().GetString("db")
		if err != nil {
			return err
		}
		outputPath, err := cmd.Flags().GetString("output")
		if err != nil {
			return err
		}

		name, err := sqlite.Backup(context.Background(), path, outputPath, false)
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
		path, err := cmd.Flags().GetString("db")
		if err != nil {
			return err
		}
		outputPath := args[0]

		clonePath, err := sqlite.Clone(context.Background(), path, outputPath, false)
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
		dbPath, err := cmd.Flags().GetString("db")
		if err != nil {
			return err
		}
		err = sqlite.Compact(context.Background(), dbPath, false)
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
		started := time.Now()
		dbPath, err := cmd.Flags().GetString("db")
		if err != nil {
			return err
		}
		debug, err := cmd.Flags().GetBool("debug")
		if err != nil {
			return err
		}
		err = sqlite.MigrateUp(context.Background(), dbPath, false, debug)
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
		dbPath, err := cmd.Flags().GetString("db")
		if err != nil {
			return err
		}
		debug, err := cmd.Flags().GetBool("debug")
		if err != nil {
			return err
		}
		ctx := context.Background()
		db, err := sqlite.Open(ctx, dbPath, false, debug)
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
		dbPath, err := cmd.Flags().GetString("db")
		if err != nil {
			return err
		}
		debug, err := cmd.Flags().GetBool("debug")
		if err != nil {
			return err
		}
		ctx := context.Background()
		db, err := sqlite.Open(ctx, dbPath, false, debug)
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

var cmdReportParse = &cobra.Command{
	Use:   "parse <documentID>",
	Short: "Parse a turn report document",
	Long:  `Parse a turn report that has been uploaded to the server.`,
	Args:  cobra.ExactArgs(1), // require document id
}

var cmdReportUpload = &cobra.Command{
	Use:   "upload <document>",
	Short: "Upload a new turn report",
	Long:  `Upload turn reports to the server.`,
	Args:  cobra.ExactArgs(1), // require path to turn report
	RunE: func(cmd *cobra.Command, args []string) error {
		startedAt := time.Now()

		dbPath, err := cmd.Flags().GetString("db")
		if err != nil {
			return err
		}
		debug, err := cmd.Flags().GetBool("debug")
		if err != nil {
			return err
		}
		path := args[0]
		log.Printf("report: upload %q\n", path)
		if ext := filepath.Ext(path); strings.ToLower(ext) != ".docx" {
			log.Fatalf("report: ext %q: not a DOCX file\n", path)
		}
		var name string
		if cmd.Flags().Changed("name") {
			if value, err := cmd.Flags().GetString("name"); err != nil {
				return err
			} else {
				name = value
			}
		}
		owner, err := cmd.Flags().GetString("owner")
		if err != nil {
			return err
		}

		ctx := context.Background()
		db, err := sqlite.Open(ctx, dbPath, true, debug)
		if err != nil {
			log.Fatalf("db: open: %v\n", err)
		}
		defer func() {
			_ = db.Close()
		}()

		if name == "" {
			name = path
		}
		if owner == "" {
			owner = "sysop"
		}

		authSvc := auth.New(db)
		tzSvc, err := iana.New(db)
		usersSvc := users.New(db, authSvc, tzSvc)
		docSvc := documents.New(db, usersSvc)

		docId, err := docSvc.LoadDocxFromFS(path, name, owner)
		if err != nil {
			log.Fatalf("%q: %v\n", path, err)
		}

		log.Printf("report: docId %d: upload: completed in %v\n", docId, time.Since(startedAt))
		return domains.ErrNotImplemented
	},
}

var cmdUserCreate = &cobra.Command{
	Use:   "create <username>",
	Short: "Create a new user",
	Long:  `Create a new user with specified name.`,
	Args:  cobra.ExactArgs(1), // require username
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := cmd.Flags().GetString("db")
		if err != nil {
			return err
		}
		userName := strings.ToLower(args[0])
		if !users.ValidateUsername(userName) {
			return domains.ErrInvalidUsername
		}
		emailSet, email := cmd.Flags().Changed("email"), userName+"@ottoapp"
		if emailSet {
			email, err = cmd.Flags().GetString("email")
			if err != nil {
				return err
			} else if !users.ValidateEmail(email) {
				return domains.ErrInvalidEmail
			}
			email = strings.ToLower(email)
		}
		passwordSet, password := cmd.Flags().Changed("password"), phrases.Generate(6)
		if passwordSet {
			password, err = cmd.Flags().GetString("password")
			if err != nil {
				return err
			} else if !auth.ValidatePassword(password) {
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
		db, err := sqlite.Open(ctx, path, true, false)
		if err != nil {
			log.Fatalf("db: open: %v\n", err)
		}
		defer func() {
			_ = db.Close()
		}()

		authSvc := auth.New(db)
		tzSvc, err := iana.New(db)
		usersSvc := users.New(db, authSvc, tzSvc)

		_, err = usersSvc.CreateUser(userName, email, password, loc)
		if err != nil {
			log.Fatalf("user %q: create: %v\n", userName, err)
		}

		log.Printf("user %q: email %q: tz %q: password %q: created\n", userName, email, loc.String(), password)

		return nil
	},
}

var cmdUserUpdate = &cobra.Command{
	Use:          "update <username>",
	Short:        "Update user record",
	Long:         `Update fields for a specific user. At least one update flag must be provided.`,
	Args:         cobra.ExactArgs(1), // require username
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		dbPath, err := cmd.Flags().GetString("db")
		if err != nil {
			return err
		}
		userName := strings.ToLower(args[0])
		if !users.ValidateUsername(userName) {
			return domains.ErrInvalidUsername
		}
		var newUserName *string
		userNameSet := cmd.Flags().Changed("username")
		if userNameSet {
			value, err := cmd.Flags().GetString("username")
			if err != nil {
				return err
			} else if !users.ValidateUsername(value) {
				return fmt.Errorf("invalid new username")
			}
			newUserName = &value
		}
		var newEmail *string
		emailSet := cmd.Flags().Changed("email")
		if emailSet {
			value, err := cmd.Flags().GetString("email")
			if err != nil {
				return err
			} else if !users.ValidateEmail(value) {
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
			if !auth.ValidatePassword(value) {
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
		db, err := sqlite.Open(ctx, dbPath, true, false)
		if err != nil {
			log.Fatalf("db: open: %v\n", err)
		}
		defer func() {
			_ = db.Close()
		}()

		authSvc := auth.New(db)
		tzSvc, err := iana.New(db)
		if err != nil {
			return err
		}
		usersSvc := users.New(db, authSvc, tzSvc)

		user, err := usersSvc.GetUserByUsername(userName)
		if err != nil {
			return fmt.Errorf("user: %q: update %v\n", userName, err)
		}

		err = usersSvc.UpdateUser(user.ID, newUserName, newEmail, newTimeZone)
		if err != nil {
			return fmt.Errorf("user: %q: update %v\n", userName, err)
		} else {
			if emailSet {
				log.Printf("user %q: email %q: updated", userName, *newEmail)
			}
			if tzSet {
				log.Printf("user %q: tz %q: updated", userName, newTimeZone.String())
			}
		}

		if newPassword != nil {
			err = authSvc.UpdateUserSecret(user.ID, *newPassword)
			if err != nil {
				return fmt.Errorf("user %q: password %q: update %v\n", userName, *newPassword, err)
			}
			log.Printf("user %q: password %q: updated", userName, *newPassword)
		}

		return nil
	},
}
