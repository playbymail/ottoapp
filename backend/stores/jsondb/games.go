// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package jsondb

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

// Games is a map of Code to Game data
type Games map[string]*Game

type Game struct {
	// Code is the unique identifier for a game
	Code string
	// Description is a short description of the game
	Description string
	// Clans is a map of Handle to ClanNo
	Clans map[string]int
	// Turns is a map of TurnId (yyyy-mm) to Turn data
	Turns map[string]*Turn
}

type Turn struct {
	Id        string
	Year      int
	Month     int
	No        int
	OrdersDue time.Time
}

func LoadGames(path string) (Games, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	jsonGames := map[string]*struct {
		Code        string         `json:"code"`
		Description string         `json:"description"`
		Clans       map[string]int `json:"clans"`
		Turns       map[string]*struct {
			Id        string `json:"id"`
			Year      int    `json:"year"`
			Month     int    `json:"month"`
			No        int    `json:"no"`
			OrdersDue string `json:"orders-due"`
		} `json:"turns"`
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
			Clans:       map[string]int{},
			Turns:       map[string]*Turn{},
		}
		for handle, clan := range jsonGame.Clans {
			game.Clans[handle] = clan
		}
		for _, jsonGameTurn := range jsonGame.Turns {
			turn := &Turn{
				Id:    fmt.Sprintf("%04d-%02d", jsonGameTurn.Year, jsonGameTurn.Month),
				Year:  jsonGameTurn.Year,
				Month: jsonGameTurn.Month,
				No:    jsonGameTurn.No,
			}
			if jsonGameTurn.OrdersDue != "" {
				turn.OrdersDue, err = ianaTimezoneHelper(jsonGameTurn.OrdersDue)
				if err != nil {
					return nil, fmt.Errorf("%s: %w", game.Code, err)
				}
			}
			game.Turns[turn.Id] = turn
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
