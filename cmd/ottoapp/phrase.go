// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package main

import (
	"fmt"
	"log"

	"github.com/mdhender/phrases/v2"
	"github.com/spf13/cobra"
)

func cmdPhrase() *cobra.Command {
	length := 6
	addFlags := func(cmd *cobra.Command) error {
		cmd.Flags().IntVar(&length, "length", length, "number of words in phrase")
		return nil
	}
	var cmd = &cobra.Command{
		Use:   "phrase",
		Short: "random phrase",
		RunE: func(cmd *cobra.Command, args []string) error {
			if length < 1 {
				length = 1
			} else if length > 16 {
				length = 16
			}
			fmt.Println(phrases.Generate(length))
			return nil
		},
	}
	if err := addFlags(cmd); err != nil {
		log.Fatal(err)
	}
	return cmd
}
