// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package runners

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
)

// GetUserProfile tests the GET /api/users/me endpoint
// It requires authentication, so it logs in first
func (r *Runner) GetUserProfile(email, password string) error {
	// Create a cookie jar to maintain session
	jar, err := cookiejar.New(nil)
	if err != nil {
		return fmt.Errorf("create cookie jar: %w", err)
	}
	client := &http.Client{Jar: jar}

	// Step 1: Login to get session cookie
	loginUrl := r.apiUrl("/api/login", nil)
	log.Printf("logging in to %s as %s\n", loginUrl, email)

	loginBody := map[string]string{
		"Email":    email,
		"Password": password,
	}
	loginJSON, err := json.Marshal(loginBody)
	if err != nil {
		return fmt.Errorf("marshal login: %w", err)
	}

	loginResp, err := client.Post(loginUrl, "application/json", bytes.NewBuffer(loginJSON))
	if err != nil {
		return fmt.Errorf("login request: %w", err)
	}
	defer loginResp.Body.Close()

	loginRespBody, _ := io.ReadAll(loginResp.Body)
	if loginResp.StatusCode != http.StatusOK {
		log.Printf("login failed: status %d: %s\n", loginResp.StatusCode, string(loginRespBody))
		return fmt.Errorf("login failed: status %d: %s", loginResp.StatusCode, string(loginRespBody))
	}

	log.Printf("login successful: %s\n", string(loginRespBody))

	// Step 2: Call /api/users/me
	usersMeUrl := r.apiUrl("/api/users/me", nil)
	log.Printf("fetching user profile from %s\n", usersMeUrl)

	resp, err := client.Get(usersMeUrl)
	if err != nil {
		return fmt.Errorf("get user profile: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("request failed: status %d: %s\n", resp.StatusCode, string(body))
		return fmt.Errorf("request failed: status %d: %s", resp.StatusCode, string(body))
	}

	// Pretty print the JSON response
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, body, "", "  "); err != nil {
		fmt.Println(string(body))
	} else {
		fmt.Println(prettyJSON.String())
	}

	return nil
}
