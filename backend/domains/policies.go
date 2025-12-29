// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package domains

import (
	"strconv"
	"strings"
)

// Actor is:
// - A representation of the authenticated principal
// - With roles and capabilities already attached
// So that:
// - authorization logic becomes simple and pure-Go
// - ID returns to being just a database primary key
//
// It gives you a clean, scalable authz core without forcing you to build a full RBAC engine.

type Actor struct {
	ID ID // the numeric user_id from the DB
	// one or more roles: admin, gm, sysop, service, player
	Admin   bool
	GM      bool
	Service bool
	Sysop   bool
	User    bool
}

func (a *Actor) IsValid() bool   { return a != nil && a.ID != InvalidID }
func (a *Actor) IsAdmin() bool   { return a != nil && a.Admin }
func (a *Actor) IsGM() bool      { return a != nil && a.GM }
func (a *Actor) IsService() bool { return a != nil && a.Service }
func (a *Actor) IsSysop() bool   { return a != nil && a.Sysop }
func (a *Actor) IsUser() bool    { return a != nil && a.User }

type Role string

const (
	RoleAdmin   Role = "admin"
	RoleGM      Role = "gm"
	RoleService Role = "service"
	RoleSysop   Role = "sysop"
	RoleUser    Role = "user"
)

type Roles map[Role]bool

func (r Roles) IsAdmin() bool {
	return r != nil && r[RoleAdmin]
}
func (r Roles) IsSysop() bool {
	return r != nil && r[RoleSysop]
}
func (r Roles) IsUser() bool {
	return r != nil && r[RoleUser]
}

// ActorAuthorizations represents the roles and permissions an actor has
// for a target user. Used to build API responses for frontend authorization.
type ActorAuthorizations struct {
	Roles       []string        `json:"roles"`
	Permissions map[string]bool `json:"permissions"`
}

// Validators are policy helpers that return true if item passes validation.

func ValidateClan(clan int) error {
	if !(0 <= clan && clan <= 999) {
		return ErrBadInput
	}
	return nil
}

func ValidateGameDescription(descr string) error {
	if descr != strings.TrimSpace(descr) {
		return ErrBadInput
	} else if descr == "" {
		return ErrBadInput
	}
	return nil
}

func ValidateGameID(id string) error {
	if len(id) != 4 {
		return ErrBadInput
	} else if n, err := strconv.Atoi(id); err != nil {
		return err
	} else if !(0 <= n && n <= 999) {
		return ErrBadInput
	}
	return nil
}

func ValidateGameTurn(year, month int) error {
	if !(899 <= year && year <= 9999) {
		return ErrBadInput
	} else if year == 899 && month != 12 {
		return ErrBadInput
	} else if !(1 <= month && month <= 12) {
		return ErrBadInput
	}
	return nil
}

func ValidateEmail(email string) error {
	if email != strings.TrimSpace(email) {
		return ErrBadInput
	} else if !(4 < len(email) && len(email) < 48) {
		return ErrBadInput
	} else if !strings.Contains(email, "@") {
		return ErrBadInput
	} else if email != strings.ToLower(email) {
		return ErrBadInput
	}
	return nil
}

func ValidateGameCode(code string) error {
	if code != strings.TrimSpace(code) {
		return ErrBadInput
	} else if !(1 < len(code) && len(code) <= 4) {
		return ErrBadInput
	} else if code != strings.ToUpper(code) {
		return ErrBadInput
	}
	return nil
}

func ValidateHandle(handle string) error {
	// Must not be empty
	if len(handle) == 0 || len(handle) > 14 {
		return ErrBadInput
	}
	// Must be lowercase (no uppercase allowed)
	if handle != strings.ToLower(handle) {
		return ErrBadInput
	}
	// Must start with a letter
	if handle[0] < 'a' || handle[0] > 'z' {
		return ErrBadInput
	}
	// All characters must be lowercase letters, digits, or underscores
	for _, ch := range handle {
		if !((ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') || ch == '_') {
			return ErrBadInput
		}
	}
	return nil
}

func ValidatePassword(plainTextPassword string) error {
	// no leading/trailing whitespace
	if strings.TrimSpace(plainTextPassword) != plainTextPassword {
		return ErrBadInput
	}
	// length bounds 8...128 bytes
	l := len(plainTextPassword)
	if !(8 <= l && l <= 128) {
		return ErrBadInput
	}
	return nil
}

func ValidateRole(role string) error {
	if !containsWord(role, "active", "sysop", "admin", "gm", "player", "guest", "user", "tn3", "tn3.1") {
		return ErrBadInput
	}
	return nil
}

func ValidateUsername(userName string) error {
	if userName != strings.TrimSpace(userName) {
		return ErrBadInput
	} else if !(3 <= len(userName) && len(userName) < 35) {
		return ErrBadInput
	} else if strings.Contains(userName, "@") {
		return ErrBadInput
	}
	return nil
}

func containsWord(word string, list ...string) bool {
	for _, elem := range list {
		if elem == word {
			return true
		}
	}
	return false
}
