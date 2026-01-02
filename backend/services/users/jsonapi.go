// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package users

import (
	"fmt"
	"time"

	"github.com/playbymail/ottoapp/backend/domains"
)

// UserView is the JSON:API view for a user
type UserView struct {
	ID          string          `jsonapi:"primary,user"` // singular when sending a payload
	Handle      string          `jsonapi:"attr,handle"`
	Username    string          `jsonapi:"attr,username"`
	Email       string          `jsonapi:"attr,email"`
	EmailOptIn  bool            `jsonapi:"attr,email-opt-in"`
	Timezone    string          `jsonapi:"attr,timezone"`
	Roles       []string        `jsonapi:"attr,roles,omitempty"`
	Permissions map[string]bool `jsonapi:"attr,permissions,omitempty"`
	CreatedAt   time.Time       `jsonapi:"attr,created-at,iso8601"`
	UpdatedAt   time.Time       `jsonapi:"attr,updated-at,iso8601"`
}

/*
| Scenario                                                                                 | Meaning                                | Status                       |
| ---------------------------------------------------------------------------------------- | -------------------------------------- | ---------------------------- |
| **Unauthenticated** (no valid session/JWT)                                               | “You must sign in first.”              | **401 Unauthorized**         |
| **Authenticated but not allowed** (role too low, wrong resource, bad old password, etc.) | “You’re signed in, but can’t do this.” | **403 Forbidden**            |
| **Wrong resource identifier** (the IDs don’t match, or you’re editing someone else)      | “You can’t edit this user.”            | **403 Forbidden**            |
| **User doesn’t exist / soft-deleted**                                                    | “Resource missing.”                    | **404 Not Found**            |
| **Validation or business rule fails** (password too short, missing field)                | Client error but not authorization     | **422 Unprocessable Entity** |
*/

// UserView constructs a UserView with permissions based on actor's privileges
func (s *Service) UserView(user *domains.User_t, actor, target *domains.Actor) *UserView {
	aa := s.authzSvc.BuildActorAuth(actor, target)
	return &UserView{
		ID:          fmt.Sprintf("%d", user.ID),
		Handle:      user.Handle,
		Username:    user.Username,
		Email:       user.Email,
		EmailOptIn:  user.EmailOptIn,
		Timezone:    user.Locale.Timezone.Location.String(),
		Roles:       aa.Roles,
		Permissions: aa.Permissions,
		CreatedAt:   user.Created.UTC(),
		UpdatedAt:   user.Updated.UTC(),
	}
}
