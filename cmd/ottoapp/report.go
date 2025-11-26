// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/playbymail/ottoapp/backend/parsers"
	"github.com/playbymail/ottoapp/backend/parsers/office"
	"github.com/spf13/cobra"
)

var cmdReportExtract = &cobra.Command{
	Use:   "extract <documentID>",
	Short: "Extract text from a turn report document",
	Args:  cobra.ExactArgs(1), // require path
	RunE: func(cmd *cobra.Command, args []string) error {
		//startedAt := time.Now()
		path := args[0]
		rptPath, err := cmd.Flags().GetString("output")
		if err != nil {
			return err
		}

		var docx *parsers.Docx
		if input, err := os.ReadFile(path); err != nil {
			log.Fatal(err)
		} else if docx, err = parsers.ParseDocx(bytes.NewReader(input), false, false); err != nil {
			log.Fatal(err)
		}

		output := &bytes.Buffer{}
		for _, line := range bytes.Split(docx.Text, []byte{'\n'}) {
			output.Write(bytes.TrimSpace(line))
			output.WriteByte('\n')
		}
		if err := os.WriteFile(rptPath, output.Bytes(), 0o644); err != nil {
			log.Fatalf("error: %v\n", err)
		}

		return nil
	},
}

var cmdReportParse = &cobra.Command{
	Use:   "parse <documentID>",
	Short: "Parse a turn report document",
	Long:  `Parse a turn report that has been uploaded to the server.`,
	Args:  cobra.ExactArgs(1), // require document id
	RunE: func(cmd *cobra.Command, args []string) error {
		startedAt := time.Now()
		toDocXmlOnly, err := cmd.Flags().GetBool("docxml-only")
		if err != nil {
			return err
		}
		path := args[0]

		if input, err := os.ReadFile(path); err != nil {
			log.Fatal(err)
		} else if toDocXmlOnly {
			p, err := office.DocXMLPath(path)
			if err != nil {
				return errors.Join(fmt.Errorf("parser: parse file"), err)
			}
			fmt.Printf("%s\n", string(p))
		} else {
			docx, err := parsers.ParseDocx(bytes.NewReader(input), true, true)
			if err != nil {
				return errors.Join(fmt.Errorf("parser: parse file"), err)
			}
			fmt.Printf("%s\n", string(docx.Text))
		}

		fmt.Printf("report: parse %q: completed in %v\n", path, time.Since(startedAt))
		return nil
	},
}
