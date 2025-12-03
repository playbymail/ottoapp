// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package games

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/mdhender/phrases/v2"
	"github.com/playbymail/ottoapp/backend/domains"
	"github.com/playbymail/ottoapp/backend/iana"
	"github.com/playbymail/ottoapp/backend/services/authz"
	"github.com/playbymail/ottoapp/backend/stores/sqlite/sqlc"
)

type ImportFile struct {
	Config  *ImportConfig `json:"config,omitempty"`
	Games   []*ImportGame
	Players map[string]*ImportPlayer
}

type ImportConfig struct {
	Mailgun *ImportConfigMailgun `json:"mailgun,omitempty"`
}

type ImportConfigMailgun struct {
	Domain  string `json:"domain,omitempty"`
	From    string `json:"from,omitempty"`
	ReplyTo string `json:"reply-to,omitempty"`
	ApiKey  string `json:"api-key,omitempty"`
	ApiBase string `json:"api-base,omitempty"`
}

type ImportGame struct {
	Id          string `json:"id"`
	Description string `json:"description"`
	SetupTurn   ImportGameSetupTurn
}

type ImportGameSetupTurn struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	No    int `json:"no"`
}

type ImportPlayer struct {
	Handle            string              `json:"handle"`
	Active            bool                `json:"active,omitempty"`
	Username          string              `json:"username,omitempty"`
	Email             string              `json:"email,omitempty"`
	Timezone          string              `json:"tz,omitempty"`
	Config            *ImportPlayerConfig `json:"config,omitempty"`
	Roles             []string            `json:"roles,omitempty"`
	Games             []*ImportPlayerGame `json:"games,omitempty"`
	Password          string              `json:"password,omitempty"`
	actor             *domains.Actor
	loc               *time.Location
	generatedPassword bool
}

func (player ImportPlayer) ChangedPassword() bool {
	return player.generatedPassword
}

type ImportPlayerConfig struct {
	EmailOptIn      bool `json:"email-opt-in,omitempty"`
	SendWelcomeMail bool `json:"send-welcome-mail,omitempty"`
}

type ImportPlayerGame struct {
	Id          string `json:"id"`
	Clan        int    `json:"clan"`
	SetupTurnNo int    `json:"setupTurnNo,omitempty"`
}

func LoadGameData(path string) (*ImportFile, error) {
	// load the game data from the json file (should be in the database?)
	if sb, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, errors.Join(domains.ErrInvalidPath, domains.ErrNotExists)
		}
		return nil, errors.Join(domains.ErrInvalidPath, err)
	} else if sb.IsDir() || !sb.Mode().IsRegular() {
		return nil, errors.Join(domains.ErrInvalidPath, domains.ErrNotFile)
	}
	input, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var data ImportFile
	err = json.Unmarshal(input, &data)
	if err != nil {
		return nil, err
	}
	// force default configs for backwards compatibility
	if data.Config == nil {
		data.Config = &ImportConfig{}
	}
	for _, player := range data.Players {
		if player.Config == nil {
			player.Config = &ImportPlayerConfig{}
		}
	}
	return &data, nil
}

// Import expects to run from the command line, so log errors.
func (s *Service) Import(data *ImportFile) error {
	actor := &domains.Actor{ID: authz.SysopId, Sysop: true}

	// json import won't update the player handle, so we must make that update
	for handle, player := range data.Players {
		player.Handle = strings.ToLower(handle)
	}

	// by convention, players should be sorted by handle, with penguin and catbird first
	var players []*ImportPlayer
	for _, player := range data.Players {
		players = append(players, player)
	}
	sort.Slice(players, func(i, j int) bool {
		a, b := players[i], players[j]
		if a.Handle == "penguin" {
			return true
		} else if b.Handle == "penguin" {
			return false
		} else if a.Handle == "catbird" {
			return true
		} else if b.Handle == "catbird" {
			return false
		}
		return a.Handle < b.Handle
	})

	var err error

	// perform basic validation
	for _, game := range data.Games {
		if errId := domains.ValidateGameID(game.Id); errId != nil {
			log.Printf("error: game %q: id %q: %v\n", game.Id, game.Id, errId)
			err = errId
		}
		if errDescr := domains.ValidateGameDescription(game.Description); errDescr != nil {
			log.Printf("error: game %q: description %q: %v\n", game.Id, game.Description, errDescr)
			err = errDescr
		}
		if errSetup := domains.ValidateGameTurn(game.SetupTurn.Year, game.SetupTurn.Month); err != nil {
			log.Printf("error: game %q: setup %d-%d: %v\n", game.Id, game.SetupTurn.Year, game.SetupTurn.Month, errSetup)
			err = errSetup
		}
	}

	// perform basic validation
	for _, player := range players {
		handle := player.Handle
		if errHandle := domains.ValidateHandle(handle); errHandle != nil {
			log.Printf("error: player %q: handle %q: %v\n", player.Handle, player.Handle, errHandle)
			err = errHandle
		}
		if errUsername := domains.ValidateUsername(player.Username); errUsername != nil {
			log.Printf("error: player %q: userName %q: %v\n", player.Handle, player.Username, errUsername)
			err = errUsername
		}
		if errEmail := domains.ValidateEmail(player.Email); errEmail != nil {
			log.Printf("error: player %q: email %q: %v\n", player.Handle, player.Email, errEmail)
			err = errEmail
		}
		if player.Password == "" {
			player.generatedPassword, player.Password = true, phrases.Generate(6)
		}
		if errPassword := domains.ValidatePassword(player.Password); errPassword != nil {
			log.Printf("error: player %q: password %q: %v\n", player.Handle, player.Password, errPassword)
			err = errPassword
		}
		if tz, ok := iana.Normalize(player.Timezone); !ok {
			err = domains.ErrInvalidTimezone
			log.Printf("error: player %q: tz %q: %v\n", player.Handle, player.Timezone, err)
		} else if loc, errTimezone := time.LoadLocation(player.Timezone); err != nil {
			log.Printf("error: player %q: tz %q: %v\n", player.Handle, player.Timezone, errTimezone)
			err = errTimezone
		} else {
			player.Timezone, player.loc = tz, loc
		}
		if len(player.Roles) == 0 {
			err = domains.ErrInvalidRole
			log.Printf("error: player %q: roles %v: %v\n", player.Handle, player.Roles, err)
		} else {
			for _, role := range player.Roles {
				if errRole := domains.ValidateRole(role); errRole != nil {
					log.Printf("error: player %q: role %q: %v\n", player.Handle, role, err)
					err = errRole
				}
			}
		}

		// perform basic validation
		for _, game := range player.Games {
			if errId := domains.ValidateGameID(game.Id); errId != nil {
				log.Printf("error: player %q: game %q: %v\n", player.Handle, game.Id, errId)
				err = errId
			}
			if errClan := domains.ValidateClan(game.Clan); err != nil {
				log.Printf("error: player %q: game %q: clan %d: %v\n", player.Handle, game.Id, game.Clan, errClan)
				err = errClan
			}
		}
	}
	// return if the basic validation found errors
	if err != nil {
		log.Printf("error: not importing due to validation errors")
		return err
	}

	q := s.db.Queries()
	ctx := s.db.Context()
	now := time.Now().UTC()

	// create or update games
	for _, game := range data.Games {
		if _, err := q.UpsertGame(ctx, sqlc.UpsertGameParams{
			GameID:         game.Id,
			Description:    game.Description,
			SetupTurnNo:    int64(game.SetupTurn.No),
			SetupTurnYear:  int64(game.SetupTurn.Year),
			SetupTurnMonth: int64(game.SetupTurn.Month),
			IsActive:       true,
			CreatedAt:      now.Unix(),
			UpdatedAt:      now.Unix(),
		}); err != nil {
			log.Printf("error: game %q: upsert %v\n", game.Id, err)
			return err
		}
	}

	// create or update players
	for _, player := range players {
		user, err := s.usersSvc.UpsertUser(player.Handle, player.Email, player.Username, player.loc)
		if err != nil {
			log.Printf("error: player %q: upsert %v\n", player.Handle, err)
			return err
		}
		player.actor = &domains.Actor{ID: user.ID, User: true}

		// create or update the player's password
		//log.Printf(" info: player %q: password %q\n", player.Handle, player.Password)
		_, err = s.authnSvc.UpdateCredentials(actor, player.actor, "", player.Password)
		if err != nil {
			log.Printf("error: player %q: password %q: upsert %v\n", player.Handle, player.Password, err)
			return err
		}
		if player.generatedPassword {
			fmt.Printf("%s: password %q\n", player.Handle, player.Password)
		}

		// create or update the player's roles
		for _, role := range player.Roles {
			// the assign is implemented as an "upsert."
			err = s.authzSvc.AssignRole(user.ID, role)
			if err != nil {
				log.Printf("error: player %q: role %q: upsert %v\n", player.Handle, role, err)
				return err
			}
		}

		// create or update the player's games
		for _, game := range player.Games {
			_, err := q.UpsertGameUserClan(ctx, sqlc.UpsertGameUserClanParams{
				GameID:      game.Id,
				UserID:      int64(user.ID),
				Clan:        int64(game.Clan),
				SetupTurnNo: int64(game.SetupTurnNo),
				CreatedAt:   now.Unix(),
				UpdatedAt:   now.Unix(),
			})
			if err != nil {
				log.Printf("error: player %q: game %q: upsert %v\n", player.Handle, game.Id, err)
				return err
			}
		}
	}

	return err
}
