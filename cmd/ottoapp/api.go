// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/playbymail/ottoapp"
	"github.com/playbymail/ottoapp/backend/iana"
	"github.com/playbymail/ottoapp/backend/servers/rest"
	"github.com/playbymail/ottoapp/backend/services/authn"
	"github.com/playbymail/ottoapp/backend/services/authz"
	"github.com/playbymail/ottoapp/backend/services/documents"
	"github.com/playbymail/ottoapp/backend/services/games"
	"github.com/playbymail/ottoapp/backend/services/users"
	"github.com/playbymail/ottoapp/backend/sessions"
	"github.com/playbymail/ottoapp/backend/stores/sqlite"
	"github.com/playbymail/ottoapp/backend/versions"
	"github.com/spf13/cobra"
)

var cmdApiServe = &cobra.Command{
	Use:          "serve",
	Short:        "start the API server",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		const checkVersion = true
		path, err := cmd.Flags().GetString("db")
		if err != nil {
			return err
		}
		quiet, _ := cmd.Flags().GetBool("quiet")
		verbose, _ := cmd.Flags().GetBool("verbose")
		debug, _ := cmd.Flags().GetBool("debug")
		if quiet {
			verbose = false
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
		//if value, err := cmd.Flags().GetString("userdata"); err != nil {
		//	return err
		//} else {
		//	options = append(options, rest.WithUserData(value))
		//}

		log.Printf("[serve] db %q\n", path)
		ctx := context.Background()
		var db *sqlite.DB
		if path == ":memory:" {
			// server has the ability to use a temporary database for testing.
			db, err = sqlite.OpenTempDB(ctx)
		} else {
			db, err = sqlite.Open(ctx, path, checkVersion, quiet, verbose, debug)
		}
		if err != nil {
			return errors.Join(fmt.Errorf("db.open"), err)
		}
		defer func() {
			log.Printf("[serve] db: close\n")
			_ = db.Close()
		}()

		authzSvc := authz.New(db)
		authnSvc := authn.New(db, authzSvc)
		tzSvc, err := iana.New(db, quiet, verbose, debug)
		if err != nil {
			return errors.Join(fmt.Errorf("iana.new"), err)
		}
		usersSvc := users.New(db, authnSvc, authzSvc, tzSvc) // uses sqlite + domains
		documentsSvc, err := documents.New(db, authzSvc, usersSvc, quiet, verbose, debug)
		if err != nil {
			return errors.Join(fmt.Errorf("sessions.new"), err)
		}
		sessionsSvc, err := sessions.New(db, authnSvc, authzSvc, usersSvc, 24*time.Hour, 15*time.Minute)
		if err != nil {
			return errors.Join(fmt.Errorf("sessions.new"), err)
		}
		gamesSvc, err := games.New(db, authnSvc, authzSvc, usersSvc, quiet, verbose, debug)
		if err != nil {
			return err
		}
		versionSvc := versions.New(ottoapp.Version())

		// Import test users for in-memory database
		if path == ":memory:" {
			panic("obsolete: replace with sync.Service")
			//var data games.ImportFile
			//err = json.Unmarshal(memdbPlayersJsonData, &data)
			//if err != nil {
			//	log.Printf("[memdb] warning: failed to import test users: %v\n", err)
			//} else if err = gamesSvc.Import(&data); err != nil {
			//	log.Printf("[memdb] warning: failed to import test users: %v\n", err)
			//}
		}

		s, err := rest.New(authnSvc, authzSvc, documentsSvc, gamesSvc, sessionsSvc, tzSvc, usersSvc, versionSvc, options...)
		if err != nil {
			return errors.Join(fmt.Errorf("rest.new"), err)
		}
		err = s.Run()
		if err != nil {
			return errors.Join(fmt.Errorf("rest.run"), err)
		}

		return nil
	},
}
