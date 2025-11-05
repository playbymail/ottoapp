// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package main implements the ottosrv web server.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/playbymail/ottoapp"
	"github.com/playbymail/ottoapp/backend/servers/rest"
	ssi "github.com/playbymail/ottoapp/backend/services/sessions"
	"github.com/playbymail/ottoapp/backend/stores/sqlite"
	"github.com/spf13/cobra"
)

func main() {
	log.SetFlags(log.Lshortfile)

	var cmdRoot = &cobra.Command{
		Use:   "ottosrv",
		Short: "OttoMap web server",
		Long:  `OttoSrv is a web server for OttoMap.`,
	}
	cmdRoot.CompletionOptions.DisableDefaultCmd = true

	cmdRoot.AddCommand(cmdServe)
	cmdServe.PersistentFlags().Bool("csrf-guard", false, "enable csrf guards")
	cmdServe.PersistentFlags().String("db", ".", "path to the database file")
	cmdServe.PersistentFlags().Bool("debug", false, "enable debugging options")
	cmdServe.PersistentFlags().Bool("dev", false, "enable development mode")
	cmdServe.PersistentFlags().Bool("log-routes", false, "enable route logging")
	cmdServe.PersistentFlags().Duration("sessions-reap-interval", 15*time.Minute, "interval to remove expired sessions")
	cmdServe.PersistentFlags().Duration("sessions-ttl", 24*time.Hour, "session duration")
	cmdServe.PersistentFlags().Duration("shutdown-delay", 30*time.Second, "delay for services to close during shutdown")
	cmdServe.PersistentFlags().Duration("shutdown-timer", 0, "timer to shut server down")

	cmdRoot.AddCommand(cmdVersion)
	cmdVersion.Flags().Bool("build-info", false, "show build information")

	if err := cmdRoot.Execute(); err != nil {
		os.Exit(1)
	}
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

		sessionManager, err := ssi.NewSessionManager(db, db, 24*time.Hour, 15*time.Minute)
		if err != nil {
			_ = db.Close()
			log.Fatalf("[serve] sessionManager: %v\n", err)
		}

		s, err := rest.New(sessionManager, options...)
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
