// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package jsondb

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// Users is a map of Handle to User data
type Users map[string]*User

type User struct {
	Handle     string
	UserName   string
	Email      string
	EmailOptIn bool
	Tz         *time.Location
	Password   struct {
		Password       string
		CreatePassword bool
		UpdatePassword bool
		ChangePassword bool
	}
	Roles map[string]bool
}

func LoadUsers(path string) (Users, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	jsonUsers := map[string]*struct {
		Handle   string   `json:"handle"`
		UserName string   `json:"user-name"`
		Email    string   `json:"email"`
		Tz       string   `json:"tz"` // IANA time zone name
		Password string   `json:"password"`
		Roles    []string `json:"roles"`
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
			Roles:    map[string]bool{},
		}
		// load the timezone, returning any errors
		user.Tz, err = time.LoadLocation(jsonUser.Tz)
		if err != nil {
			return nil, fmt.Errorf("iana: %s: %w", jsonUser.Tz, err)
		}
		user.Password.Password = jsonUser.Password
		for _, role := range jsonUser.Roles {
			role = strings.ToLower(role)
			switch role {
			case "active":
				user.Roles[role] = true
			case "admin":
				user.Roles[role] = true
			case "change-password":
				user.Password.ChangePassword = true
			case "create-password":
				user.Password.CreatePassword = true
			case "email-opt-in":
				user.EmailOptIn = true
			case "email-opt-out":
				user.EmailOptIn = false
			case "gm":
				user.Roles[role] = true
			case "guest":
				user.Roles[role] = true
			case "inactive":
				user.Roles[role] = true
			case "player":
				user.Roles[role] = true
			case "service":
				user.Roles[role] = true
			case "update-password":
				user.Password.UpdatePassword = true
			case "user":
				user.Roles[role] = true
			default:
				log.Printf("jsondb: import: user: role %q: ignoring\n", role)
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
