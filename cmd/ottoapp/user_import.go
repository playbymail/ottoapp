// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package main

//import (
//	"context"
//	_ "embed"
//	"encoding/csv"
//	"errors"
//	"fmt"
//	"log"
//	"os"
//	"regexp"
//	"strings"
//	"time"
//
//	"github.com/mdhender/phrases/v2"
//	"github.com/playbymail/ottoapp/backend/auth"
//	"github.com/playbymail/ottoapp/backend/iana"
//	"github.com/playbymail/ottoapp/backend/stores/sqlite"
//	"github.com/playbymail/ottoapp/backend/stores/sqlite/sqlc"
//	"github.com/playbymail/ottoapp/backend/users"
//	"github.com/spf13/cobra"
//	"golang.org/x/crypto/bcrypt"
//)
//
////go:embed testdata/memdb-players.csv
//var memdbPlayersCsvData []byte
//
//var cmdUserImport = &cobra.Command{
//	Use:          "import <csv-file>",
//	Short:        "Import users from CSV file",
//	Long:         `Import users from a CSV file. Creates new users or updates existing ones based on email.`,
//	Args:         cobra.ExactArgs(1),
//	SilenceUsage: true,
//	RunE: func(cmd *cobra.Command, args []string) error {
//		csvPath := args[0]
//
//		// Read and validate CSV file
//		records, err := readAndValidateCSV(csvPath)
//		if err != nil {
//			return err
//		}
//
//		log.Printf("validated %d records from %s\n", len(records), csvPath)
//
//		// Open database
//		dbPath, err := cmd.Flags().GetString("db")
//		if err != nil {
//			return err
//		}
//
//		ctx := context.Background()
//		db, err := sqlite.Open(ctx, dbPath, true, false)
//		if err != nil {
//			return errors.Join(fmt.Errorf("db.open"), err)
//		}
//		defer func() {
//			log.Printf("db: close\n")
//			_ = db.Close()
//		}()
//
//		authSvc := auth.New(db)
//		tzSvc, err := iana.New(db)
//		if err != nil {
//			return errors.Join(fmt.Errorf("iana.new"), err)
//		}
//		usersSvc := users.New(db, authSvc, tzSvc)
//
//		// Track generated passwords
//		passwordUpdates := make(map[int]string) // map of row index to generated password
//
//		// Process each record
//		for i, record := range records {
//			generatedPassword, err := importUser(db, authSvc, usersSvc, record)
//			if err != nil {
//				log.Printf("row %d: email %q: failed: %v\n", i+2, record.Email, err)
//			} else {
//				log.Printf("row %d: email %q: imported/updated\n", i+2, record.Email)
//				if generatedPassword != "" {
//					passwordUpdates[i] = generatedPassword
//				}
//			}
//		}
//
//		// Update CSV with generated passwords
//		if len(passwordUpdates) > 0 {
//			err = updateCSVWithPasswords(csvPath, passwordUpdates)
//			if err != nil {
//				log.Printf("warning: failed to update CSV with generated passwords: %v\n", err)
//			} else {
//				log.Printf("updated %d password(s) in %s\n", len(passwordUpdates), csvPath)
//			}
//		}
//
//		return nil
//	},
//}
//
//type userRecord struct {
//	Clan     string
//	Handle   string
//	Username string
//	Email    string
//	Roles    map[string]bool
//	Timezone string
//	Password string
//}
//
//func readAndValidateCSV(path string) ([]userRecord, error) {
//	f, err := os.Open(path)
//	if err != nil {
//		return nil, fmt.Errorf("open csv: %w", err)
//	}
//	defer f.Close()
//
//	reader := csv.NewReader(f)
//	rows, err := reader.ReadAll()
//	if err != nil {
//		return nil, fmt.Errorf("read csv: %w", err)
//	}
//
//	return parseAndValidateCSVRows(rows)
//}
//
//func readAndValidateCSVData(data []byte) ([]userRecord, error) {
//	reader := csv.NewReader(strings.NewReader(string(data)))
//	rows, err := reader.ReadAll()
//	if err != nil {
//		return nil, fmt.Errorf("read csv: %w", err)
//	}
//
//	return parseAndValidateCSVRows(rows)
//}
//
//func parseAndValidateCSVRows(rows [][]string) ([]userRecord, error) {
//
//	if len(rows) == 0 {
//		return nil, fmt.Errorf("csv is empty")
//	}
//
//	// Validate header
//	header := rows[0]
//	expectedHeaders := []string{"Clan", "Handle", "User Name", "Email", "Roles", "Timezone", "Password"}
//	if len(header) != len(expectedHeaders) {
//		return nil, fmt.Errorf("header has %d columns, expected %d", len(header), len(expectedHeaders))
//	}
//	for i, h := range header {
//		if h != expectedHeaders[i] {
//			return nil, fmt.Errorf("header column %d: got %q, expected %q", i, h, expectedHeaders[i])
//		}
//	}
//
//	// Validate records
//	var records []userRecord
//	clanPattern := regexp.MustCompile(`^0\d{3}$`)
//	emailPattern := regexp.MustCompile(`^[^@]+@[^@]+$`)
//
//	for i, row := range rows[1:] {
//		lineNum := i + 2
//
//		if len(row) != 7 {
//			return nil, fmt.Errorf("row %d: has %d columns, expected 7", lineNum, len(row))
//		}
//
//		var ur = userRecord{
//			Clan:     row[0],
//			Handle:   row[1],
//			Username: row[2],
//			Email:    strings.ToLower(row[3]),
//			Roles:    map[string]bool{"active": true, "user": true}, // all users should be active
//			Timezone: row[5],
//			Password: row[6],
//		}
//		for _, role := range strings.Split(row[4], "+") {
//			role := strings.TrimSpace(role)
//			if role != "" {
//				ur.Roles[role] = true
//			}
//		}
//
//		// Validate Clan
//		if !clanPattern.MatchString(ur.Clan) {
//			return nil, fmt.Errorf("row %d: clan %q is not a 4-digit number starting with 0", lineNum, ur.Clan)
//		}
//
//		// Validate Roles
//		if ur.Username == "Penguin" {
//			if !ur.Roles["admin"] {
//				return nil, fmt.Errorf("row %d: Penguin must have role 'admin'", lineNum)
//			}
//		} else {
//			if ur.Roles["admin"] {
//				return nil, fmt.Errorf("row %d: %q is not allowed role %q", lineNum, ur.Username, "admin")
//			}
//		}
//
//		// Validate Timezone
//		if _, ok := iana.CanonicalName(ur.Timezone); !ok {
//			return nil, fmt.Errorf("row %d: timezone %q is not a valid IANA name", lineNum, ur.Timezone)
//		}
//
//		// Validate Email
//		if !emailPattern.MatchString(ur.Email) {
//			return nil, fmt.Errorf("row %d: email %q does not look like 'x@x'", lineNum, ur.Email)
//		}
//
//		records = append(records, ur)
//	}
//
//	return records, nil
//}
//
//func importUsersFromCSVData(db *sqlite.DB, authSvc *auth.Service, usersSvc *users.Service, csvData []byte) error {
//	records, err := readAndValidateCSVData(csvData)
//	if err != nil {
//		return err
//	}
//
//	log.Printf("[memdb] importing %d test users\n", len(records))
//
//	for i, record := range records {
//		_, err := importUser(db, authSvc, usersSvc, record)
//		if err != nil {
//			log.Printf("[memdb] row %d: email %q: failed: %v\n", i+2, record.Email, err)
//			return fmt.Errorf("import user %q: %w", record.Email, err)
//		} else {
//			log.Printf("[memdb] imported user %q with email %q\n", record.Username, record.Email)
//		}
//	}
//
//	return nil
//}
//
//func importUser(db *sqlite.DB, authSvc *auth.Service, usersSvc *users.Service, record userRecord) (string, error) {
//	// Generate password if empty
//	password := record.Password
//	var generatedPassword string
//	if password == "" {
//		password = phrases.Generate(6)
//		generatedPassword = password
//		// log.Printf("  generated password: %q\n", password)
//	}
//
//	// Get timezone location
//	loc, err := time.LoadLocation(record.Timezone)
//	if err != nil {
//		return "", fmt.Errorf("load timezone %q: %w", record.Timezone, err)
//	}
//
//	// Try to get existing user by email
//	_, err = usersSvc.GetUserByEmail(record.Email)
//
//	if err != nil {
//		// User doesn't exist, create new one via service
//		_, err = usersSvc.CreateUser(record.Username, record.Email, record.Handle, password, loc)
//		if err != nil {
//			return "", fmt.Errorf("create user: %w", err)
//		}
//		log.Printf("  created user %q\n", record.Username)
//	} else {
//		// User exists - for now just log it
//		// TODO: Implement update logic if needed
//		log.Printf("  user %q already exists, skipping\n", record.Username)
//	}
//
//	// Assign roles from CSV
//	// Start transaction for role assignment
//	tx, err := db.Stdlib().BeginTx(db.Context(), nil)
//	if err != nil {
//		return "", err
//	}
//	defer tx.Rollback()
//
//	qtx := db.Queries().WithTx(tx)
//	ctx := db.Context()
//
//	// Get user ID
//	userID, err := qtx.GetUserIDByEmail(ctx, record.Email)
//	if err != nil {
//		return "", fmt.Errorf("get user ID: %w", err)
//	}
//
//	// Get current roles
//	currentRoles, err := qtx.GetUserRoles(ctx, userID)
//	if err != nil {
//		return "", fmt.Errorf("get user roles: %w", err)
//	}
//
//	roleMap := make(map[string]bool)
//	for _, role := range currentRoles {
//		roleMap[role] = true
//	}
//
//	// Ensure role from CSV exists
//	now := time.Now().UTC()
//	for role, required := range record.Roles {
//		if required && !roleMap[role] {
//			err = qtx.AssignUserRole(ctx, sqlc.AssignUserRoleParams{
//				UserID:    userID,
//				RoleID:    role,
//				CreatedAt: now.Unix(),
//				UpdatedAt: now.Unix(),
//			})
//			if err != nil {
//				return "", fmt.Errorf("assign role %q: %w", role, err)
//			}
//			log.Printf("  assigned role %q\n", role)
//			roleMap[role] = true
//		}
//	}
//
//	// Commit transaction
//	err = tx.Commit()
//	if err != nil {
//		return "", fmt.Errorf("commit transaction: %w", err)
//	}
//
//	return generatedPassword, nil
//}
//
//func updateCSVWithPasswords(csvPath string, passwordUpdates map[int]string) error {
//	// Read the CSV file
//	f, err := os.Open(csvPath)
//	if err != nil {
//		return fmt.Errorf("open csv: %w", err)
//	}
//	defer f.Close()
//
//	reader := csv.NewReader(f)
//	rows, err := reader.ReadAll()
//	if err != nil {
//		return fmt.Errorf("read csv: %w", err)
//	}
//
//	// Update the password column (index 6) for the specified rows
//	for rowIdx, password := range passwordUpdates {
//		// rowIdx is 0-based for data rows, but we need to account for the header
//		actualRow := rowIdx + 1 // +1 to skip header
//		if actualRow < len(rows) && len(rows[actualRow]) > 6 {
//			rows[actualRow][6] = password
//		}
//	}
//
//	// Write the updated CSV back to the file
//	tempPath := csvPath + ".tmp"
//	tmpFile, err := os.Create(tempPath)
//	if err != nil {
//		return fmt.Errorf("create temp file: %w", err)
//	}
//	defer tmpFile.Close()
//
//	writer := csv.NewWriter(tmpFile)
//	err = writer.WriteAll(rows)
//	if err != nil {
//		return fmt.Errorf("write csv: %w", err)
//	}
//	writer.Flush()
//
//	if err := writer.Error(); err != nil {
//		return fmt.Errorf("flush csv: %w", err)
//	}
//
//	// Close the temp file before renaming
//	tmpFile.Close()
//
//	// Replace the original file with the updated one
//	err = os.Rename(tempPath, csvPath)
//	if err != nil {
//		return fmt.Errorf("rename temp file: %w", err)
//	}
//
//	return nil
//}
//
//func hashPassword(plainTextPassword string) (string, error) {
//	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), bcrypt.DefaultCost)
//	if err != nil {
//		return "", err
//	}
//	return string(hashedPasswordBytes), nil
//}
