// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package jsondb

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mdhender/phrases/v2"
)

// Users is a map of Handle to User data
type Users map[string]*User

type User struct {
	Handle   string
	UserName string
	Email    string
	Tz       *time.Location
	Password struct {
		Password string
		Update   bool
	}
	Roles struct {
		Active     bool
		EmailOptIn bool
		Roles      []string
	}
}

func LoadUsers(path string) (Users, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	jsonUsers := map[string]*struct {
		Handle         string   `json:"handle"`
		UserName       string   `json:"user-name"`
		Email          string   `json:"email"`
		Tz             string   `json:"tz"` // IANA time zone name
		Password       string   `json:"password"`
		CreatePassword bool     `json:"create-password"`
		UpdatePassword bool     `json:"update-password"`
		ChangePassword bool     `json:"change-password"`
		Roles          []string `json:"roles"`
	}{}
	err = json.Unmarshal(data, &jsonUsers)
	if err != nil {
		return nil, err
	}

	users := map[string]*User{}
	for handle, jsonUser := range jsonUsers {
		user := &User{
			Handle:   strings.ToLower(handle),
			UserName: jsonUser.UserName,
			Email:    strings.ToLower(jsonUser.Email),
		}
		// load the timezone, returning any errors
		user.Tz, err = time.LoadLocation(jsonUser.Tz)
		if err != nil {
			return nil, fmt.Errorf("iana: %s: %w", jsonUser.Tz, err)
		}
		user.Password.Password = jsonUser.Password
		if user.Password.Password == "" {
			user.Password.Password = phrases.Generate(6, ".")
			user.Password.Update = true
		}
		if jsonUser.CreatePassword || jsonUser.ChangePassword {
			user.Password.Password = phrases.Generate(6, ".")
			user.Password.Update = true
		}
		if jsonUser.UpdatePassword {
			user.Password.Update = true
		}
		for _, role := range jsonUser.Roles {
			switch role {
			case "active":
				user.Roles.Active = true
			case "inactive":
				user.Roles.Active = false
			case "email-opt-in":
				user.Roles.EmailOptIn = true
			case "email-opt-out":
				user.Roles.EmailOptIn = false
			default:
				user.Roles.Roles = append(user.Roles.Roles)
			}
		}
		users[user.Handle] = user
	}
	return users, nil
}

func LoadUser(path string, handle string) (*User, error) {
	users, err := LoadUsers(path)
	if err != nil {
		return nil, err
	}
	user, ok := users[handle]
	if !ok {
		return nil, fmt.Errorf("%s: not found", handle)
	}
	return user, nil
}
