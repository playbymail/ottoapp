// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package rest

import (
	"encoding/json"
	"net/http"
)

// handleGetProfile returns the current user's profile information
// using the users service to fetch fresh data from the database.
func (s *Server) handleGetProfile(w http.ResponseWriter, r *http.Request) {
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
	profileData := struct {
		Handle   string `json:"handle"`
		Email    string `json:"email"`
		Timezone string `json:"timezone"`
	}{
		Handle:   user.Username,
		Email:    user.Email,
		Timezone: user.Locale.Timezone.Location.String(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(profileData)
}
