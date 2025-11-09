// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package rest

import (
	"encoding/json"
	"net/http"
	"time"
)

// handleGetProfile returns the current user's profile information
// using the users service to fetch fresh data from the database.
func (s *Server) handleGetProfile(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Username string   `json:"username,omitempty"`
		Email    string   `json:"email,omitempty"`
		Timezone string   `json:"timezone,omitempty"`
		Errors   []string `json:"errors,omitempty"`
	}

	// Get the current session to identify the user
	sess, err := s.services.sessionsSvc.GetCurrentSession(r)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	// Fetch user data from the users service
	user, err := s.services.usersSvc.GetUserByID(sess.User.ID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Create a response payload with the fields the frontend needs
	profileData := response{
		Username: user.Username,
		Email:    user.Email,
		Timezone: user.Locale.Timezone.Location.String(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(profileData)
}

// handlePostProfile updates the current user's profile information.
// Users can update their email and timezone, but not their username.
func (s *Server) handlePostProfile(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email    *string `json:"email"`
		Timezone *string `json:"timezone"`
	}
	type response struct {
		Username string   `json:"username,omitempty"`
		Email    string   `json:"email,omitempty"`
		Timezone string   `json:"timezone,omitempty"`
		Errors   []string `json:"errors,omitempty"`
	}

	// Get the current session to identify the user
	sess, err := s.services.sessionsSvc.GetCurrentSession(r)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	// Parse the request body
	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(response{
			Errors: []string{"Invalid request body"},
		})
		return
	}

	var rsp response

	// Validate email string if provided
	if req.Email != nil && *req.Email != "" {
		if !s.services.usersSvc.ValidateEmail(*req.Email) {
			rsp.Errors = append(rsp.Errors, "Invalid email")
		}
	}

	// Convert timezone string to *time.Location if provided
	var timezone *time.Location
	if req.Timezone != nil && *req.Timezone != "" {
		timezone, err = s.services.ianaSvc.Location(*req.Timezone)
		if err != nil {
			rsp.Errors = append(rsp.Errors, "Invalid timezone")
		}
	}

	// Return if there were validation errors
	if len(rsp.Errors) != 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(rsp)
		return
	}

	// Update the user (userName is nil since users can't change their username)
	err = s.services.usersSvc.UpdateUser(sess.User.ID, nil, req.Email, timezone)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		rsp.Errors = []string{err.Error()}
		_ = json.NewEncoder(w).Encode(rsp)
		return
	}

	// Fetch and return the updated profile
	user, err := s.services.usersSvc.GetUserByID(sess.User.ID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	rsp.Username = user.Username
	rsp.Email = user.Email
	rsp.Timezone = user.Locale.Timezone.Location.String()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(rsp)
}
