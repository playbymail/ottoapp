// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package main

import (
	"fmt"

	"github.com/playbymail/ottoapp/backend/runners"
	"github.com/spf13/cobra"
)

var cmdAppVersion = &cobra.Command{
	Use:           "version",
	Short:         "display the application's version number",
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		showBuildInfo, err := cmd.Flags().GetBool("show-build-info")
		if err != nil {
			return err
		}
		r := runners.New("http", "127.0.0.1", "8181")
		err = r.GetVersion(showBuildInfo)
		if err != nil {
			return err
		}
		return nil
	},
}

var cmdAppTestUserProfile = &cobra.Command{
	Use:          "test-user-profile",
	Short:        "Test the GET /api/users/me endpoint",
	Long:         `Test the GET /api/users/me endpoint by logging in and fetching the user profile.`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		host, err := cmd.Flags().GetString("host")
		if err != nil {
			return err
		}
		port, err := cmd.Flags().GetString("port")
		if err != nil {
			return err
		}
		email, err := cmd.Flags().GetString("email")
		if err != nil {
			return err
		}
		password, err := cmd.Flags().GetString("password")
		if err != nil {
			return err
		}
		if password == "" {
			return fmt.Errorf("password is required")
		}

		r := runners.New("http", host, port)
		return r.GetUserProfile(email, password)
	},
}
