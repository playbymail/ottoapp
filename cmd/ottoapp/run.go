// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/playbymail/ottoapp/backend/parsers/reports"
	"github.com/playbymail/ottoapp/backend/services/reports/cst"
	"github.com/playbymail/ottoapp/backend/services/reports/cst2"
	parsers "github.com/playbymail/ottoapp/backend/services/reports/docx"
	"github.com/playbymail/ottoapp/backend/services/reports/lexers"
	"github.com/playbymail/ottoapp/backend/services/reports/lexers/lemon"
	"github.com/playbymail/ottoapp/backend/services/reports/lexers/tokens"
	"github.com/playbymail/ottoapp/backend/services/reports/scrubbers"
	"github.com/spf13/cobra"
)

func cmdRun() *cobra.Command {
	addFlags := func(cmd *cobra.Command) error {
		cmd.PersistentFlags().Bool("show-timing", false, "time command")
		return nil
	}
	var cmd = &cobra.Command{
		Use:   "run",
		Short: "Run utility commands",
	}
	if err := addFlags(cmd); err != nil {
		log.Fatal(err)
	}
	cmd.AddCommand(cmdRunLemon())
	cmd.AddCommand(cmdRunLexer())
	cmd.AddCommand(cmdRunParse())
	cmd.AddCommand(cmdRunWelcomeEmail())
	return cmd
}

func cmdRunLemon() *cobra.Command {
	addFlags := func(cmd *cobra.Command) error {
		return nil
	}
	var cmd = &cobra.Command{
		Use:   "lemon",
		Short: "Run lemon commands",
	}
	cmd.AddCommand(cmdRunLemonLex())
	cmd.AddCommand(cmdRunLemonParse())
	if err := addFlags(cmd); err != nil {
		log.Fatalf("%s: %v\n", cmd.Use, err)
	}
	return cmd
}

func cmdRunLemonLex() *cobra.Command {
	addFlags := func(cmd *cobra.Command) error {
		return nil
	}
	var cmd = &cobra.Command{
		Use:   "lex <report-file>",
		Short: "Run lemon lexer against report file",
		Args:  cobra.ExactArgs(1), // require path to report file
		RunE: func(cmd *cobra.Command, args []string) error {
			showTiming, _ := cmd.Flags().GetBool("show-timing")
			reportFile := args[0]
			if strings.ToLower(filepath.Ext(reportFile)) != ".docx" {
				return fmt.Errorf("report file must be .docx")
			}
			var docx *parsers.Docx
			if input, err := os.ReadFile(reportFile); err != nil {
				log.Fatal(err)
			} else if docx, err = parsers.ParseDocx(bytes.NewReader(input), true, true); err != nil {
				log.Fatal(err)
			}
			startedAt := time.Now()
			lexer := lemon.New(docx.Text)
			totalTokens := 0
			var tok tokens.Token
			for {
				tok, totalTokens = lexer.Next(), totalTokens+1
				log.Printf("lemon: lexer: %5d: %5d: %-18s: %q\n", tok.Line, tok.Col, tok.Kind, tok.Value)
				if tok.Kind == tokens.EOF {
					break
				}
			}
			if showTiming {
				log.Printf("%s: lines %6d: tokens %7d: lexed in %v\n", reportFile, tok.Line, totalTokens, time.Since(startedAt))
				return nil
			}
			return nil
		},
	}
	if err := addFlags(cmd); err != nil {
		log.Fatalf("%s: %v\n", cmd.Use, err)
	}
	return cmd
}

func cmdRunLemonParse() *cobra.Command {
	addFlags := func(cmd *cobra.Command) error {
		return nil
	}
	var cmd = &cobra.Command{
		Use:   "parse <report-file>",
		Short: "Run lemon parser against report file",
		Args:  cobra.ExactArgs(1), // require path to report file
		RunE: func(cmd *cobra.Command, args []string) error {
			showTiming, _ := cmd.Flags().GetBool("show-timing")
			reportFile := args[0]
			if strings.ToLower(filepath.Ext(reportFile)) != ".docx" {
				return fmt.Errorf("report file must be .docx")
			}
			var docx *parsers.Docx
			if input, err := os.ReadFile(reportFile); err != nil {
				log.Fatal(err)
			} else if docx, err = parsers.ParseDocx(bytes.NewReader(input), false, false); err != nil {
				log.Fatal(err)
			}
			startedAt := time.Now()
			lexer := lemon.New(docx.Text)
			totalTokens := 0
			var tok tokens.Token
			for {
				tok, totalTokens = lexer.Next(), totalTokens+1
				//log.Printf("lemon: lexer: %5d: %5d: %-18s: %q\n", tok.Line, tok.Col, tok.Kind, tok.Value)
				if tok.Kind == tokens.EOF {
					break
				}
			}
			if showTiming {
				log.Printf("%s: lines %6d: tokens %7d: lexed in %v\n", reportFile, tok.Line, totalTokens, time.Since(startedAt))
				return nil
			}
			return nil
		},
	}
	if err := addFlags(cmd); err != nil {
		log.Fatalf("%s: %v\n", cmd.Use, err)
	}
	return cmd
}

func cmdRunLexer() *cobra.Command {
	showTiming := false
	addFlags := func(cmd *cobra.Command) error {
		cmd.Flags().BoolVar(&showTiming, "show-timing", showTiming, "time lexer")
		return nil
	}

	var cmd = &cobra.Command{
		Use:   "lexer",
		Short: "Run lexer against file",
		Args:  cobra.ExactArgs(1), // require path to file to lex
		RunE: func(cmd *cobra.Command, args []string) error {
			input := args[0]
			data, err := os.ReadFile(input)
			if err != nil {
				return err
			}
			lines := len(bytes.Split(data, []byte{'\n'}))
			startedAt := time.Now()
			tokens := lexers.Scan(data)
			if showTiming {
				log.Printf("%s: lines %6d: tokens %7d: lexed in %v\n", input, lines, len(tokens), time.Since(startedAt))
				return nil
			}
			for i, token := range tokens {
				line, col := token.Position()
				fmt.Printf("%7d: %6d: %6d: %-20s: %q\n",
					i, line, col, token.Kind.String(), string(token.Bytes()))
			}
			return nil
		},
	}

	if err := addFlags(cmd); err != nil {
		log.Fatal(err)
	}

	return cmd
}

func cmdRunParse() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "parse",
		Short: "Run parser commands",
	}

	cmd.AddCommand(cmdRunParseCst())
	cmd.AddCommand(cmdRunParseReportFile())
	cmd.AddCommand(cmdRunParseTurnReport())

	return cmd
}

func cmdRunParseCst() *cobra.Command {
	showLineNo := false
	showColNo := false
	showTiming := false
	showTokenKind := false
	addFlags := func(cmd *cobra.Command) error {
		cmd.Flags().BoolVar(&showColNo, "col", showColNo, "print column number")
		cmd.Flags().BoolVar(&showLineNo, "line", showLineNo, "print line number")
		cmd.Flags().BoolVar(&showTiming, "show-timing", showTiming, "time lexer")
		cmd.Flags().BoolVar(&showTokenKind, "kind", showTokenKind, "print token kind")
		return nil
	}

	var cmd = &cobra.Command{
		Use:          "cst",
		Short:        "Run cst parser against file",
		SilenceUsage: true,
		Args:         cobra.ExactArgs(1), // require path to file to parse
		RunE: func(cmd *cobra.Command, args []string) error {
			input := args[0]
			data, err := os.ReadFile(input)
			if err != nil {
				return err
			}
			lines := len(bytes.Split(data, []byte{'\n'}))
			startedAt := time.Now()
			tokens := lexers.Scan(data)
			cnodes, err := cst2.Parse(tokens)
			if err != nil {
				return err
			} else if cnodes != nil {
				lineNo := 1
				for _, line := range bytes.Split(cnodes.Source(), []byte{'\n'}) {
					if bytes.HasSuffix(line, []byte{'-', '-', '>', '8'}) {
						fmt.Printf("%s\n", line)
						continue
					}
					fmt.Printf("%5d: %s\n", lineNo, line)
					lineNo++
				}
				fmt.Printf("%s: lines %6d: tokens %7d: parsed in %v\n", input, lines, len(tokens), time.Since(startedAt))
				return nil
			}
			result := cst.Parse(tokens)
			if showTiming {
				log.Printf("%s: lines %6d: tokens %7d: errors %7d: parsed in %v\n", input, lines, len(tokens), len(result.Errors()), time.Since(startedAt))
				return nil
			}
			options := []cst.PrintOption{(cst.WithName(input))}
			if showColNo {
				options = append(options, cst.WithColumnNo(), cst.WithLineNo())
			}
			if showLineNo {
				options = append(options, cst.WithLineNo())
			}
			if showTokenKind {
				options = append(options, cst.WithTokenKind())
			}
			cst.PrettyPrint(os.Stdout, result, options...)

			return nil
		},
	}

	if err := addFlags(cmd); err != nil {
		log.Fatal(err)
	}

	return cmd
}

func cmdRunParseReportFile() *cobra.Command {
	var outputPath string
	trimLeading := true
	trimTrailing := true
	addFlags := func(cmd *cobra.Command) error {
		cmd.Flags().BoolVar(&trimLeading, "trim-leading", trimLeading, "trim leading spaces")
		cmd.Flags().BoolVar(&trimTrailing, "trim-trailing", trimTrailing, "trim trailing spaces")
		cmd.Flags().StringVar(&outputPath, "output", outputPath, "path to save parsed report to")
		return nil
	}

	var cmd = &cobra.Command{
		Use:   "report <turn-report-file-name>",
		Short: "Parse a turn report file",
		Args:  cobra.ExactArgs(1), // require path to turn report file
		RunE: func(cmd *cobra.Command, args []string) error {
			startedAt := time.Now()
			report := args[0]
			if strings.ToLower(filepath.Ext(report)) != ".docx" {
				return fmt.Errorf("turn report file must be .docx")
			}

			var docx *parsers.Docx
			if input, err := os.ReadFile(report); err != nil {
				log.Fatal(err)
			} else if docx, err = parsers.ParseDocx(bytes.NewReader(input), trimLeading, trimTrailing); err != nil {
				log.Fatal(err)
			}

			if outputPath == "" {
				fmt.Println(string(docx.Text))
				return nil
			}
			if err := os.WriteFile(outputPath, docx.Text, 0o644); err != nil {
				log.Fatalf("error: %v\n", err)
			}
			log.Printf("%s: created in %v\n", outputPath, time.Since(startedAt))
			return nil
		},
	}

	if err := addFlags(cmd); err != nil {
		log.Fatal(err)
	}

	return cmd
}

func cmdRunParseTurnReport() *cobra.Command {
	patchNA := false
	rawExtractFile := ""
	scrubbedFile := ""
	showJson := false
	showStats := false
	showTiming := false
	addFlags := func(cmd *cobra.Command) error {
		cmd.Flags().BoolVar(&patchNA, "patch-na", patchNA, "patch N/A")
		cmd.Flags().StringVar(&rawExtractFile, "raw-extract", rawExtractFile, "path to save raw extract to")
		cmd.Flags().StringVar(&scrubbedFile, "scrubbed-extract", scrubbedFile, "path to save scrubbed extract to")
		cmd.Flags().BoolVar(&showJson, "show-json", showJson, "show json after extract")
		cmd.Flags().BoolVar(&showStats, "show-stats", showStats, "show stats after extract")
		cmd.Flags().BoolVar(&showTiming, "show-timing", showStats, "show timing after extract")
		return nil
	}

	var cmd = &cobra.Command{
		Use:          "turn-report <turn-report-file-name>",
		Short:        "Parse a turn report file",
		SilenceUsage: true,
		Args:         cobra.ExactArgs(1), // require path to turn report file
		RunE: func(cmd *cobra.Command, args []string) error {
			startedAt := time.Now()

			var docx *parsers.Docx
			docxFileName := args[0]
			if data, err := os.ReadFile(docxFileName); err != nil {
				return err
			} else if docx, err = parsers.ParseDocx(bytes.NewReader(data), false, false); err != nil {
				return err
			}
			if showTiming {
				log.Printf("%s: %d bytes\n", args[0], len(docx.Text))
			}
			if rawExtractFile != "" {
				err := os.WriteFile(rawExtractFile, docx.Text, 0o644)
				if err != nil {
					return err
				}
				if showTiming {
					log.Printf("%s: wrote raw extract\n", rawExtractFile)
				}
			}

			lines := scrubbers.Scrub(bytes.Split(docx.Text, []byte{'\n'}), patchNA)
			if scrubbedFile != "" {
				output := bytes.Join(lines, []byte{'\n'})
				if len(output) == 0 {
					output = []byte{'\n'}
				} else if output[len(output)-1] != '\n' {
					output = append(output, '\n')
				}
				err := os.WriteFile(scrubbedFile, output, 0o644)
				if err != nil {
					return err
				}
				if showTiming {
					log.Printf("%s: wrote scrubbed extract\n", scrubbedFile)
				}
			}

			stats := reports.Stats{}
			rpt, err := reports.Parse(filepath.Base(docxFileName), bytes.Join(lines, []byte{'\n'}), reports.Statistics(&stats, "no match"))
			if err != nil {
				return err
			}
			if showJson {
				b, err := json.MarshalIndent(rpt, "", "  ")
				if err != nil {
					return err
				}
				log.Printf("rpt: %s\n", string(b))
			}
			if showStats {
				b, err := json.MarshalIndent(stats.ChoiceAltCnt, "", "  ")
				if err != nil {
					return err
				}
				log.Printf("rpt: %s\n", string(b))
			}

			if showTiming {
				log.Printf("parse: turn-report: completed in %v\n", time.Since(startedAt))
			}

			return nil
		},
	}

	if err := addFlags(cmd); err != nil {
		log.Fatal(err)
	}

	return cmd
}

func cmdRunWelcomeEmail() *cobra.Command {
	var gameDataPath string
	var gameID string
	addFlags := func(cmd *cobra.Command) error {
		cmd.Flags().StringVar(&gameDataPath, "game-data", gameDataPath, "path to game data (required)")
		if err := cmd.MarkFlagRequired("game-data"); err != nil {
			return err
		}
		cmd.Flags().StringVar(&gameID, "game", gameID, "Game ID (required)")
		if err := cmd.MarkFlagRequired("game"); err != nil {
			return err
		}
		return nil
	}

	cmd := &cobra.Command{
		Use:          "welcome-email",
		Short:        "Send a welcome email to all players in a game",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			panic("update to pull from database")
			//if gameDataPath == "" {
			//	return fmt.Errorf("--game-data is required")
			//}
			//if gameID == "" {
			//	return fmt.Errorf("--game is required")
			//}
			//
			//// load the game data from the json file (should be in the database?)
			//data, err := games.LoadGameData(gameDataPath)
			//if err != nil {
			//	return err
			//} else if data.Config.Mailgun == nil {
			//	log.Printf("error: missing mailgun config in game data\n")
			//	return nil
			//}
			//
			//ctx := context.Background()
			//
			//emailSvc, err := email.NewMailgun(data.Config.Mailgun.ApiBase, data.Config.Mailgun.Domain, data.Config.Mailgun.ApiKey, data.Config.Mailgun.From, data.Config.Mailgun.ReplyTo)
			//if err != nil {
			//	log.Fatal(err)
			//}
			//
			//// fetch the game details
			//var gd *games.Game_t
			//for _, g := range data.Games {
			//	if g.Id == gameID {
			//		gd = g
			//		break
			//	}
			//}
			//if gd == nil {
			//	log.Printf("error: unknown game %q\n", gameID)
			//	return nil
			//}
			//
			//// fetch the players in the game
			//var players []*games.ImportPlayer
			//for _, p := range data.Players {
			//	if !(p.Config != nil && p.Config.EmailOptIn == true && p.Config.SendWelcomeMail == true) {
			//		continue
			//	}
			//	for _, g := range p.Games {
			//		if g.Id == gd.Id {
			//			players = append(players, p)
			//			break
			//		}
			//	}
			//}
			//
			//log.Printf("%s: sending welcome email to %d players\n", gd.Description, len(players))
			//if len(players) == 0 {
			//	return nil
			//}
			//
			//for _, p := range players {
			//	if p.Email == "" {
			//		continue
			//	}
			//	var pgd *games.ImportPlayerGame
			//	for _, g := range p.Games {
			//		if g.Id == gd.Id {
			//			pgd = g
			//			break
			//		}
			//	}
			//	if pgd == nil { // should never happen
			//		continue
			//	}
			//
			//	// On registration:
			//	if err := emailSvc.SendWelcome(ctx, p.Email, p.Username, gd.Description, pgd.Clan, p.Password); err != nil {
			//		log.Printf("%s: %04d: %s: email %v\n", gd.Description, pgd.Clan, p.Email, err)
			//		continue
			//	}
			//
			//	log.Printf("%s: %04d: %s: email sent\n", gd.Description, pgd.Clan, p.Email)
			//}
			//
			//log.Println("Done.")
			//return nil
		},
	}

	if err := addFlags(cmd); err != nil {
		log.Fatalf("%s: %v\n", cmd.Use, err)
	}

	return cmd
}
