// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/playbymail/ottoapp/backend/parsers/azul"
	"github.com/spf13/cobra"
)

func cmdTest() *cobra.Command {
	addFlags := func(cmd *cobra.Command) error {
		return nil
	}
	var cmd = &cobra.Command{
		Use:   "test",
		Short: "test stages",
	}
	cmd.AddCommand(cmdTestAzulParser())
	if err := addFlags(cmd); err != nil {
		log.Fatal(err)
	}
	return cmd
}

func cmdTestAzulParser() *cobra.Command {
	addFlags := func(cmd *cobra.Command) error {
		return nil
	}
	var cmd = &cobra.Command{
		Use:          "azul <turn-report>",
		Short:        "Parse a turn report",
		SilenceUsage: true,
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			started := time.Now()
			quiet, _ := cmd.Flags().GetBool("quiet")
			verbose, _ := cmd.Flags().GetBool("verbose")
			debug, _ := cmd.Flags().GetBool("debug")
			if quiet {
				verbose = false
			}
			path, err := filepath.Abs(args[0])
			if err != nil {
				return err
			}
			input, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			_, err = azul.Parse(path, input, quiet, verbose, debug)
			if err != nil {
				return err
			}

			log.Printf("%s: completed in %v\n", path, time.Since(started))
			return nil
		},
	}
	if err := addFlags(cmd); err != nil {
		log.Fatal(err)
	}
	return cmd
}
