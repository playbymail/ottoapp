// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/mdhender/phrases/v2"
	"github.com/playbymail/ottoapp/backend/domains"
	"github.com/playbymail/ottoapp/backend/iana"
	"github.com/playbymail/ottoapp/backend/services/authn"
	"github.com/playbymail/ottoapp/backend/services/authz"
	"github.com/playbymail/ottoapp/backend/services/users"
	"github.com/playbymail/ottoapp/backend/stores/sqlite"
	"github.com/spf13/cobra"
)

var cmdUserCreate = &cobra.Command{
	Use:   "create <handle>",
	Short: "Create a new user",
	Long:  `Create a new user.`,
	Args:  cobra.ExactArgs(1), // require handle
	RunE: func(cmd *cobra.Command, args []string) error {
		const checkVersion = true
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
		handle := strings.ToLower(args[0])
		if err := domains.ValidateHandle(handle); err != nil {
			return errors.Join(domains.ErrInvalidUsername, err)
		}
		userNameSet, userName := cmd.Flags().Changed("username"), handle
		if userNameSet {
			userName, err = cmd.Flags().GetString("username")
			if err != nil {
				return err
			}
		}
		if err := domains.ValidateUsername(userName); err != nil {
			return errors.Join(domains.ErrInvalidUsername, err)
		}
		emailSet, email := cmd.Flags().Changed("email"), handle+"@ottoapp"
		if emailSet {
			email, err = cmd.Flags().GetString("email")
			if err != nil {
				return err
			} else if err := domains.ValidateEmail(email); err != nil {
				return errors.Join(domains.ErrInvalidEmail, err)
			}
			email = strings.ToLower(email)
		}
		passwordSet, password := cmd.Flags().Changed("password"), phrases.Generate(6)
		if passwordSet {
			password, err = cmd.Flags().GetString("password")
			if err != nil {
				return err
			} else if err := domains.ValidatePassword(password); err != nil {
				return errors.Join(domains.ErrInvalidPassword, err)
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
		db, err := sqlite.Open(ctx, path, checkVersion, quiet, verbose, debug)
		if err != nil {
			log.Fatalf("db: open: %v\n", err)
		}
		defer func() {
			_ = db.Close()
		}()

		authzSvc := authz.New(db)
		authnSvc := authn.New(db, authzSvc)
		tzSvc, err := iana.New(db, quiet, verbose, debug)
		usersSvc := users.New(db, authnSvc, authzSvc, tzSvc)

		// For the user create command, use userName as the handle
		user, err := usersSvc.CreateUser(handle, email, userName, loc)
		if err != nil {
			return errors.Join(fmt.Errorf("user %q", handle), err)
		}
		actor, err := authzSvc.GetActorById(user.ID)
		if err != nil {
			return errors.Join(fmt.Errorf("user %q", handle), err)
		}
		_, err = authnSvc.UpdateCredentials(&domains.Actor{ID: authz.SysopId, Roles: domains.Roles{Sysop: true}}, actor, "", password)
		if err != nil {
			return errors.Join(fmt.Errorf("user %q: secret %q", handle, password), err)
		}
		for _, role := range []string{"active", "user"} {
			err = authzSvc.AssignRole(user.ID, role)
			if err != nil {
				return errors.Join(fmt.Errorf("user %q: role %q", handle, role), err)
			}
		}

		log.Printf("user %q: email %q: tz %q: password %q: created\n", userName, email, loc.String(), password)

		return nil
	},
}

var cmdUserUpdate = &cobra.Command{
	Use:          "update <username>",
	Short:        "UpdatePassword user record",
	Long:         `UpdatePassword fields for a specific user. At least one update flag must be provided.`,
	Args:         cobra.ExactArgs(1), // require username
	SilenceUsage: true,
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

		authzSvc := authz.New(db)
		authnSvc := authn.New(db, authzSvc)
		tzSvc, err := iana.New(db, quiet, verbose, debug)
		if err != nil {
			return err
		}
		usersSvc := users.New(db, authnSvc, authzSvc, tzSvc)

		handle := strings.ToLower(args[0])
		user, err := usersSvc.GetUserByHandle(handle)
		if err != nil {
			return err
		}
		if cmd.Flags().Changed("username") {
			value, err := cmd.Flags().GetString("username")
			if err != nil {
				return err
			} else if err := domains.ValidateUsername(value); err != nil {
				return errors.Join(domains.ErrInvalidUsername, err)
			}
			user.Username = value
		}
		if err := domains.ValidateUsername(user.Username); err != nil {
			return errors.Join(domains.ErrInvalidUsername, err)
		}
		if cmd.Flags().Changed("email") {
			value, err := cmd.Flags().GetString("email")
			if err != nil {
				return err
			} else if err := domains.ValidateEmail(value); err != nil {
				return errors.Join(domains.ErrInvalidEmail, domains.ErrBadInput)
			}
			user.Email = value
		}
		var password string
		if cmd.Flags().Changed("password") {
			value, err := cmd.Flags().GetString("password")
			if err != nil {
				return err
			} else if value == "+" {
				value = phrases.Generate(6)
				log.Printf("generated random password: %q", value)
			}
			if err := domains.ValidatePassword(value); err != nil {
				return errors.Join(domains.ErrInvalidPassword, domains.ErrBadInput)
			}
			password = value
		}
		var loc *time.Location
		if cmd.Flags().Changed("tz") {
			if tz, err := cmd.Flags().GetString("tz"); err != nil {
				return err
			} else if ctz, ok := iana.CanonicalName(tz); !ok {
				return fmt.Errorf("invalid time zone")
			} else if loc, err = time.LoadLocation(ctz); err != nil {
				return err
			} else {
				user.Locale.Timezone.Location = loc
			}
		}

		actor, err := authzSvc.GetActorById(user.ID)
		if err != nil {
			return errors.Join(fmt.Errorf("user %q", handle), err)
		}

		updatedUser := &domains.User_t{
			ID:         user.ID,
			Username:   user.Username,
			Email:      user.Email,
			EmailOptIn: user.EmailOptIn,
			Handle:     user.Handle,
			Locale: domains.UserLocale_t{
				DateFormat: user.Locale.DateFormat,
				Timezone: domains.UserTimezone_t{
					Location: user.Locale.Timezone.Location,
				},
			},
			Roles:   user.Roles,
			Created: user.Created,
			Updated: time.Now().UTC(),
		}
		err = usersSvc.UpdateUser(updatedUser)
		if err != nil {
			return fmt.Errorf("user: %q: update %v\n", user.Handle, err)
		} else {
			log.Printf("user %q: email %q: updated", user.Handle, user.Email)
			log.Printf("user %q: tz %q: updated", user.Handle, user.Locale.Timezone.Location.String())
		}

		if password != "" {
			_, err = authnSvc.UpdateCredentials(&domains.Actor{ID: authz.SysopId, Roles: domains.Roles{Sysop: true}}, actor, "", password)
			if err != nil {
				return fmt.Errorf("user %q: password %q: update %v\n", updatedUser.Handle, password, err)
			}
			log.Printf("user %q: password %q: updated", updatedUser.Handle, password)
		}

		return nil
	},
}

var cmdUserRole = &cobra.Command{
	Use:          "role <handle>",
	Short:        "Manage user roles",
	Long:         `Add or remove roles for a specific user. At least one --add or --remove flag must be provided.`,
	Args:         cobra.ExactArgs(1), // require handle
	SilenceUsage: true,
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
		handle := strings.ToLower(args[0])
		if err := domains.ValidateHandle(handle); err != nil {
			return errors.Join(domains.ErrInvalidHandle, domains.ErrBadInput)
		}

		rolesToAdd, err := cmd.Flags().GetStringSlice("add")
		if err != nil {
			return err
		}
		rolesToRemove, err := cmd.Flags().GetStringSlice("remove")
		if err != nil {
			return err
		}

		if len(rolesToAdd) == 0 && len(rolesToRemove) == 0 {
			return fmt.Errorf("must specify at least one role to add or remove")
		}

		ctx := context.Background()
		db, err := sqlite.Open(ctx, dbPath, checkVersion, quiet, verbose, debug)
		if err != nil {
			log.Fatalf("db: open: %v\n", err)
		}
		defer func() {
			_ = db.Close()
		}()

		authzSvc := authz.New(db)
		authnSvc := authn.New(db, authzSvc)
		tzSvc, err := iana.New(db, quiet, verbose, debug)
		if err != nil {
			return err
		}
		usersSvc := users.New(db, authnSvc, authzSvc, tzSvc)

		user, err := usersSvc.GetUserByHandle(handle)
		if err != nil {
			return fmt.Errorf("user: %q: not found: %v", handle, err)
		}

		// Add roles
		for _, roleID := range rolesToAdd {
			err = authzSvc.AssignRole(user.ID, roleID)
			if err != nil {
				log.Printf("user %q: role %q: failed to add: %v", user.Handle, roleID, err)
			} else {
				log.Printf("user %q: role %q: added", user.Handle, roleID)
			}
		}

		// Remove roles
		for _, roleID := range rolesToRemove {
			err = authzSvc.RemoveRole(user.ID, roleID)
			if err != nil {
				log.Printf("user %q: role %q: failed to remove: %v", user.Handle, roleID, err)
			} else {
				log.Printf("user %q: role %q: removed", user.Handle, roleID)
			}
		}

		return nil
	},
}
