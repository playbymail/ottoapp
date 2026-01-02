// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package jsondb

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/playbymail/ottoapp/backend/domains"
)

// Games is a map of Code to Game data
type Games map[string]*Game

type Game struct {
	// Code is the unique identifier for a game
	Code string
	// Description is a short description of the game
	Description           string
	SetupTurn, ActiveTurn struct {
		Year, Month int
	}
	OrdersDue time.Time
	// Clans is a map of Handle to ClanSetup
	Clans map[string]ClanSetup
}

type ClanSetup struct {
	ClanNo    int    `json:"clan-no,omitempty"`
	SetupTurn string `json:"setup-turn,omitempty"`
}

func LoadGames(path string) (Games, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	jsonGames := map[string]*struct {
		Code        string               `json:"code"`
		Description string               `json:"description"`
		SetupTurn   string               `json:"setup-turn"`
		ActiveTurn  string               `json:"active-turn"`
		OrdersDue   string               `json:"orders-due,omitempty"`
		Clans       map[string]ClanSetup `json:"clans"`
	}{}
	err = json.Unmarshal(data, &jsonGames)
	if err != nil {
		return nil, err
	}

	games := map[string]*Game{}
	for code, jsonGame := range jsonGames {
		game := &Game{
			Code:        strings.ToUpper(code),
			Description: jsonGame.Description,
			Clans:       map[string]ClanSetup{},
		}
		if jsonGame.SetupTurn == "" {
			game.SetupTurn.Year, game.SetupTurn.Month = 899, 12
		} else if game.SetupTurn.Year, game.SetupTurn.Month, err = yearMonthHelper(jsonGame.SetupTurn); err != nil {
			return nil, fmt.Errorf("%s: setup %q: invalid", game.Code, jsonGame.SetupTurn)
		}
		if jsonGame.ActiveTurn == "" {
			game.ActiveTurn = game.SetupTurn
		} else if game.ActiveTurn.Year, game.ActiveTurn.Month, err = yearMonthHelper(jsonGame.ActiveTurn); err != nil {
			return nil, fmt.Errorf("%s: active %q: invalid", game.Code, jsonGame.ActiveTurn)
		}
		if jsonGame.OrdersDue != "" {
			game.OrdersDue, err = ianaTimezoneHelper(jsonGame.OrdersDue)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", game.Code, err)
			}
		}
		for handle, clan := range jsonGame.Clans {
			if clan.SetupTurn == "" {
				clan.SetupTurn = fmt.Sprintf("%04d-%02d", game.SetupTurn.Year, game.SetupTurn.Month)
			}
			game.Clans[handle] = clan
		}
		games[game.Code] = game
	}
	return games, nil
}

func LoadGame(path string, code string) (*Game, error) {
	games, err := LoadGames(path)
	if err != nil {
		return nil, err
	}
	game, ok := games[code]
	if !ok {
		return nil, fmt.Errorf("%s: not found", code)
	}
	return game, nil
}

func ianaTimezoneHelper(s string) (time.Time, error) {
	// split timestamp and zone
	parts := strings.LastIndex(s, " ")
	if parts == -1 {
		return time.Time{}, fmt.Errorf("%s: invalid datetime format", s)
	}
	ts, tz := s[:parts], s[parts+1:]

	// load the timezone, returning any errors
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return time.Time{}, fmt.Errorf("%s: iana: %w", s, err)
	}

	// parse the time, returning any errors
	t, err := time.ParseInLocation("2006/01/02 15:04:05", ts, loc)
	if err != nil {
		return time.Time{}, fmt.Errorf("%s: parse: %w", s, err)
	}

	return t.UTC(), nil
}

func yearMonthHelper(yearMonth string) (year, month int, err error) {
	if len(yearMonth) != 7 || yearMonth[4] != '-' {
		return 0, 0, domains.ErrBadInput
	}
	year, err = strconv.Atoi(yearMonth[:4])
	if err != nil {
		return 0, 0, fmt.Errorf("bad year: %w", err)
	}
	month, err = strconv.Atoi(yearMonth[5:])
	if err != nil {
		return 0, 0, fmt.Errorf("bad month: %w", err)
	}
	if year < 899 || (year == 899 && month != 12) || (month < 1 || month > 12) {
		return 0, 0, domains.ErrBadInput
	}
	return year, month, nil
}
