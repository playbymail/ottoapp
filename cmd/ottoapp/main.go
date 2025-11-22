// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package main implements the ottoapp command line tool.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mdhender/phrases/v2"
	"github.com/mitchellh/go-homedir"
	"github.com/playbymail/ottoapp"
	"github.com/playbymail/ottoapp/backend/binder"
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
	cmdAppTestUserProfile.Flags().String("email", "skeener917@gmail.com", "email for authentication")
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

	cmd := cmdGame()
	cmdRoot.AddCommand(cmd)

	var cmdPhrase = &cobra.Command{
		Use:   "phrase",
		Short: "random phrase",
		RunE: func(cmd *cobra.Command, args []string) error {
			length, err := cmd.Flags().GetInt("length")
			if err != nil {
				return err
			}
			if length < 1 {
				length = 1
			} else if length > 16 {
				length = 16
			}
			fmt.Println(phrases.Generate(length))
			return nil
		},
	}
	cmdRoot.AddCommand(cmdPhrase)
	cmdPhrase.Flags().Int("length", 6, "number of words in phrase")

	var cmdReport = &cobra.Command{
		Use:   "report",
		Short: "report management",
	}
	cmdRoot.AddCommand(cmdReport)
	cmdReport.AddCommand(cmdReportExtract)
	cmdReportExtract.Flags().String("output", "report.txt", "file to create")
	cmdReport.AddCommand(cmdReportParse)
	cmdReportParse.Flags().Bool("docxml-only", false, "parse to DocXML only")

	var cmdRun = &cobra.Command{
		Use:   "run",
		Short: "Run utility commands",
	}
	cmdRoot.AddCommand(cmdRun)
	cmdRun.AddCommand(cmdRunGenMake)
	cmdRunGenMake.Flags().String("root", "data/tn3.1", "root directory for data")
	cmdRunGenMake.Flags().String("output", "data/tn3.1/maps.mk", "output makefile path")

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
