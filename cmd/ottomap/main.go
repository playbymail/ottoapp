// Package main implements the OttoApp command.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/playbymail/ottoapp"
	"github.com/playbymail/ottoapp/backend/domains"
	"github.com/playbymail/ottoapp/backend/stores/sqlite"
	"github.com/spf13/cobra"
)

func main() {
	log.SetFlags(log.Lshortfile)

	var cmdRoot = &cobra.Command{
		Use:   "ottoapp",
		Short: "OttoMap application",
		Long:  `OttoMap reads turn reports and renders maps.`,
	}
	cmdRoot.CompletionOptions.DisableDefaultCmd = true
	cmdRoot.PersistentFlags().String("db", ".", "path to the database file")

	var cmdDocument = &cobra.Command{
		Use:   "document",
		Short: "Manage documents reports",
		Long:  `Commands to manage documents on the server.`,
	}
	cmdRoot.AddCommand(cmdDocument)
	cmdDocument.AddCommand(cmdDocumentParse)

	var cmdReport = &cobra.Command{
		Use:   "report",
		Short: "Manage turn reports",
		Long:  `Commands to upload, delete, and parse turn reports.`,
	}
	cmdRoot.AddCommand(cmdReport)
	cmdReport.AddCommand(cmdReportUpload)

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

var cmdDocumentParse = &cobra.Command{
	Use:   "parse <document>",
	Short: "Parse a turn report document",
	Long:  `Parse a turn report that has been uploaded to the server.`,
	Args:  cobra.ExactArgs(1), // require document id
	
}

var cmdReportUpload = &cobra.Command{
	Use:   "upload <document>",
	Short: "Upload a new turn report",
	Long:  `Upload turn reports to the server.`,
	Args:  cobra.ExactArgs(1), // require path to turn report
	RunE: func(cmd *cobra.Command, args []string) error {
		dbPath, err := cmd.Flags().GetString("db")
		if err != nil {
			return err
		}
		path := args[0]
		log.Printf("report: upload %q\n", path)

		startedAt := time.Now()

		d, err := NewDocument(path)
		if err != nil {
			log.Fatalf("%q: %v\n", path, err)
		}
		log.Printf("%q: read %d bytes\n", d.SourcePath, d.Length)
		log.Printf("%q: hash %q\n", d.SourcePath, d.ID)

		ctx := context.Background()
		db, err := sqlite.Open(ctx, dbPath, true)
		if err != nil {
			log.Fatalf("db: open: %v\n", err)
		}
		defer func() {
			_ = db.Close()
		}()

		// fetch the location of the documents folder
		documentRoot, err := db.GetKeyValue("documents.root")
		if err != nil {
			log.Printf("config %q: %v\n", "documents.root", err)
		} else if documentRoot == "" {
			log.Printf("error: document.root: not configured\n")
			return nil
		}
		log.Printf("report: documentRoot %q\n", documentRoot)

		// create a path to store the document
		if d.Path == "" {
			d.Path = filepath.Join(documentRoot, d.ID)
		}

		// write the document to the folder
		if err := os.WriteFile(d.Path, d.Data, 0o644); err != nil {
			log.Printf("report: upload: write failed %v\n", err)
			return nil
		} else {
			log.Printf("report: upload: saved %q\n", d.Path)
		}

		// update the documents table
		log.Printf("todo: create document record\n")

		log.Printf("report: %q: upload: completed in %v\n", d.Path, time.Since(startedAt))
		return domains.ErrNotImplemented
	},
}
