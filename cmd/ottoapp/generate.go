// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package main

import (
	"log"

	"github.com/playbymail/ottoapp/backend/make"
	"github.com/spf13/cobra"
)

func cmdGenerate() *cobra.Command {
	addFlags := func(cmd *cobra.Command) error {
		return nil
	}
	var cmd = &cobra.Command{
		Use:   "generate",
		Short: "generate artifacts",
	}
	cmd.AddCommand(cmdGenerateMakefile())
	if err := addFlags(cmd); err != nil {
		log.Fatal(err)
	}
	return cmd
}

func cmdGenerateMakefile() *cobra.Command {
	gameId := "0301"
	makefileName := "maps.mk"
	addFlags := func(cmd *cobra.Command) error {
		cmd.Flags().StringVar(&gameId, "game", gameId, "game to create makefile for")
		cmd.Flags().StringVar(&makefileName, "output", makefileName, "name of makefile to create")
		return nil
	}
	var cmd = &cobra.Command{
		Use:          "makefile",
		Short:        "Generate a Makefile for map generation",
		Long:         `Scans the data directory for turn reports and generates a Makefile to build maps.`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			quiet, _ := cmd.Flags().GetBool("quiet")
			verbose, _ := cmd.Flags().GetBool("verbose")
			debug, _ := cmd.Flags().GetBool("debug")
			if quiet {
				verbose = false
			}
			return make.Makefile(makefileName, gameId, quiet, verbose, debug)
		},
	}
	if err := addFlags(cmd); err != nil {
		log.Fatal(err)
	}
	return cmd
}
