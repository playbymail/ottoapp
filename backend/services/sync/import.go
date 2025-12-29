// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package sync

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/mdhender/phrases/v2"
	"github.com/playbymail/ottoapp/backend/domains"
	"github.com/playbymail/ottoapp/backend/services/authz"
	"github.com/playbymail/ottoapp/backend/stores/jsondb"
)

func (s *Service) ImportGames(path string, handles ...string) error {
	gamesData, err := jsondb.LoadGames(path)
	if err != nil {
		log.Printf("import: %s: failed %v\n", path, err)
		return fmt.Errorf("import: games %w", err)
	}

	if len(handles) == 0 {
		for _, game := range gamesData {
			handles = append(handles, game.Code)
		}
	}
	if len(handles) == 0 {
		return nil
	}
	sort.Strings(handles)

	missingGames := 0
	for _, handle := range handles {
		game, ok := gamesData[handle]
		if !ok {
			log.Printf("import: %s: not found\n", handle)
			missingGames++
			continue
		}
		// call the service to create or update the game
		err = s.importGame(game)
		if err != nil {
			log.Printf("import: %s: %s: failed %v\n", path, handle, err)
			return fmt.Errorf("import: games %w", err)
		}
	}
	if missingGames == 1 {
		return fmt.Errorf("game not found")
	} else if missingGames > 1 {
		return fmt.Errorf("%d games not found", missingGames)
	}

	return nil
}

func (s *Service) importGame(game *jsondb.Game) error {
	var err error
	if errHandle := domains.ValidateGameCode(game.Code); errHandle != nil {
		log.Printf("error: game %q: code %q: %v\n", game.Code, game.Code, errHandle)
		err = errHandle
	}
	if err != nil {
		log.Printf("import: game %s: validation failed\n", game.Code)
		return fmt.Errorf("import: game %q: %w", game.Code, err)
	}

	// call the service to create or update the game

	// call the service to create or update the game turns

	// call the service to create or update the game clans

	return domains.ErrNotImplemented
}

func (s *Service) ImportOttoAppConfig(path string, quiet, verbose, debug bool) error {
	if debug {
		log.Printf("sync: import: %q\n", path)
	}

	oac, err := jsondb.LoadOttoAppConfig(path)
	if err != nil {
		if debug {
			log.Printf("sync: import: load %s: failed %v\n", path, err)
		}
		return fmt.Errorf("sync: import: config %w", err)
	}

	err = s.configSvc.UpdateKeyValuePairs(
		"mailgun.domain", oac.Mailgun.Domain,
		"mailgun.from", oac.Mailgun.From,
		"mailgun.api.base", oac.Mailgun.ApiBase,
		"mailgun.api.key", oac.Mailgun.ApiKey,
	)
	if err != nil {
		if debug {
			log.Printf("sync: import: %s: failed %v\n", path, err)
		}
		return fmt.Errorf("sync: import: config %w", err)
	}

	return nil
}

func (s *Service) ImportUsers(path string, handles ...string) error {
	usersData, err := jsondb.LoadUsers(path)
	if err != nil {
		log.Printf("import: %s: failed %v\n", path, err)
		return fmt.Errorf("import: users %w", err)
	}

	if len(handles) == 0 {
		for _, user := range usersData {
			handles = append(handles, user.Handle)
		}
	}
	if len(handles) == 0 {
		return nil
	}
	for n := range handles {
		handles[n] = strings.ToLower(handles[n])
	}
	sort.Strings(handles)

	missingUsers := 0
	for _, handle := range handles {
		user, ok := usersData[handle]
		if !ok {
			log.Printf("import: %s: not found\n", handle)
			missingUsers++
			continue
		}
		// call the service to create or update the user, password, and roles
		err = s.importUser(user)
		if err != nil {
			log.Printf("import: %s: %s: failed %v\n", path, handle, err)
			return fmt.Errorf("import: user %w", err)
		}
	}
	if missingUsers == 1 {
		return fmt.Errorf("user not found")
	} else if missingUsers > 1 {
		return fmt.Errorf("%d users not found", missingUsers)
	}

	return nil
}

func (s *Service) importUser(u *jsondb.User) error {
	var err error
	if errHandle := domains.ValidateHandle(u.Handle); errHandle != nil {
		log.Printf("error: u %q: handle %q: %v\n", u.Handle, u.Handle, errHandle)
		err = errHandle
	}
	if errUsername := domains.ValidateUsername(u.UserName); errUsername != nil {
		log.Printf("error: u %q: userName %q: %v\n", u.Handle, u.UserName, errUsername)
		err = errUsername
	}
	if errEmail := domains.ValidateEmail(u.Email); errEmail != nil {
		log.Printf("error: u %q: email %q: %v\n", u.Handle, u.Email, errEmail)
		err = errEmail
	}
	if u.Password.Password == "" {
		u.Password.Update, u.Password.Password = true, phrases.Generate(6)
	}
	if u.Password.Update {
		if errPassword := domains.ValidatePassword(u.Password.Password); errPassword != nil {
			log.Printf("error: u %q: password %q: %v\n", u.Handle, u.Password, errPassword)
			err = errPassword
		}
	}
	for _, role := range u.Roles.Roles {
		if errRole := domains.ValidateRole(role); errRole != nil {
			log.Printf("error: u %q: role %q: %v\n", u.Handle, role, err)
			err = errRole
		}
	}
	if err != nil {
		log.Printf("import: u %s: validation failed\n", u.Handle)
		return fmt.Errorf("import: u %q: %w", u.Handle, err)
	}

	user, err := s.usersSvc.GetUserByHandle(u.Handle)
	if err != nil {
		if !(errors.Is(err, sql.ErrNoRows) || errors.Is(err, domains.ErrNotFound) || errors.Is(err, domains.ErrNotExists)) {
			log.Printf("import: u %s: getUser failed %v\n", u.Handle, err)
			return fmt.Errorf("import: u %q: %w", u.Handle, err)
		}
		user, err = s.usersSvc.CreateUser(user.Handle, user.Email, user.Username, user.Locale.Timezone.Location)
		if err != nil {
			log.Printf("import: u %s: createUser failed %v\n", user.Handle, err)
			return fmt.Errorf("import: u %q: %w", u.Handle, err)
		}
	} else {
		user.Username = u.UserName
		user.Email = u.Email
		user.Locale.DateFormat = "2006-01-02"
		user.Locale.Timezone.Location = u.Tz
		user.Updated = time.Now().UTC()
	}
	userActor := &domains.Actor{ID: user.ID}

	// todo: implement roles update correctly
	for _, role := range u.Roles.Roles {
		user.Roles[domains.Role(role)] = true
	}
	err = s.usersSvc.UpdateUser(user)
	if err != nil {
		log.Printf("import: u %s: updateUser %v\n", u.Handle, err)
		return err
	}

	// create or update the user's password
	if u.Password.Update {
		actor := &domains.Actor{ID: authz.SysopId, Sysop: true}
		_, err = s.authnSvc.UpdateCredentials(actor, userActor, "", u.Password.Password)
		if err != nil {
			log.Printf("import: u %s: password %q: upsert %v\n", user.Handle, u.Password.Password, err)
			return err
		}
		fmt.Printf("%s: password %q\n", user.Handle, u.Password.Password)
	}

	return nil
}
