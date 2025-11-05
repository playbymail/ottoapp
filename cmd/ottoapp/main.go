// Package main implements the OttoApp command.
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
		Use:   "ottoapp",
		Short: "OttoMap version",
		Long:  `OttoApp shows the version and build information.`,
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
	cmdRoot.CompletionOptions.DisableDefaultCmd = true
	cmdRoot.Flags().Bool("build-info", false, "show build information")

	if err := cmdRoot.Execute(); err != nil {
		os.Exit(1)
	}
}
