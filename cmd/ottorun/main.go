// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package main implements a command line tool that uses the
// API to run commands on the server.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/playbymail/ottoapp"
	"github.com/spf13/cobra"
)

func main() {
	log.SetFlags(log.Lshortfile)

	var cmdRoot = &cobra.Command{
		Use:   "ottorun",
		Short: "run commands on the server",
		Long:  `OttoRun runs commands on the server.`,
	}
	cmdRoot.CompletionOptions.DisableDefaultCmd = true
	var cmdVersion = &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Long:  `Display the current version of OttoMap.`,
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
	cmdRoot.AddCommand(cmdVersion)
	cmdVersion.Flags().Bool("build-info", false, "show build information")

	if err := cmdRoot.Execute(); err != nil {
		os.Exit(1)
	}
}
