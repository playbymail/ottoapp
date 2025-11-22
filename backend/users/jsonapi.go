// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package users

import "time"

// UserView is the JSON:API view for a user
type UserView struct {
	ID          string          `jsonapi:"primary,user"` // singular when sending a payload
	Username    string          `jsonapi:"attr,username"`
	Email       string          `jsonapi:"attr,email"`
	Handle      string          `jsonapi:"attr,handle"`
	Timezone    string          `jsonapi:"attr,timezone"`
	Roles       []string        `jsonapi:"attr,roles,omitempty"`
	Permissions map[string]bool `jsonapi:"attr,permissions,omitempty"`
	CreatedAt   time.Time       `jsonapi:"attr,created-at,iso8601"`
	UpdatedAt   time.Time       `jsonapi:"attr,updated-at,iso8601"`
}

// UserCreateRequest is the JSON:API request payload for creating a user
type UserCreateRequest struct {
	ID       string   `jsonapi:"primary,users"` // plural when receiving from Ember Data
	Handle   string   `jsonapi:"attr,handle,omitempty"`
	Email    string   `jsonapi:"attr,email,omitempty"`
	Username string   `jsonapi:"attr,username,omitempty"`
	Password string   `jsonapi:"attr,password,omitempty"`
	Timezone string   `jsonapi:"attr,timezone,omitempty"`
	Roles    []string `json:"attr,roles,omitempty"`
}

// UserPatchRequest is the JSON:API request payload for a user patch
type UserPatchRequest struct {
	ID       string `jsonapi:"primary,users"` // plural when receiving from Ember Data
	Username string `jsonapi:"attr,username,omitempty"`
	Email    string `jsonapi:"attr,email,omitempty"`
	Timezone string `jsonapi:"attr,timezone,omitempty"`
}
