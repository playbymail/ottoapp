// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package jsondb_test

import (
	"fmt"
	"testing"

	"github.com/playbymail/ottoapp/backend/stores/jsondb"
)

func TestLoadGames(t *testing.T) {
	games, err := jsondb.LoadGames("testdata/games.json")
	if err != nil {
		t.Fatal(err)
	}
	if want := 2; want != len(games) {
		t.Errorf("len: want %d, got %d\n", want, len(games))
	}
	code := "0300"
	if game, ok := games[code]; !ok {
		t.Errorf("%s: want *Game, got nil\n", code)
	} else {
		testGame0300(t, code, game)
	}
	code = "0301"
	if game, ok := games["0301"]; !ok {
		t.Errorf("%s: want *Game, got nil\n", code)
	} else {
		testGame0301(t, code, game)
	}
}

func testGame0300(t *testing.T, id string, game *jsondb.Game) {
	if game == nil {
		t.Errorf("%s: want *Game, got nil\n", id)
		return
	}
	if want := id; want != game.Code {
		t.Errorf("%s: Code: want %q, got %q\n", id, want, game.Code)
	}
	if want := "TN3"; want != game.Description {
		t.Errorf("%s: Description: want %q, got %q\n", id, want, game.Description)
	}
	wantClans := []int{1}
	if want := len(wantClans); want != len(game.Clans) {
		t.Errorf("%s: clans: len: want %d, got %d\n", id, want, len(game.Clans))
	}
	for _, no := range wantClans {
		handle, want := fmt.Sprintf("clan%04d", no), no
		if got, ok := game.Clans[handle]; !ok {
			t.Errorf("%s: clans: %q: want ok, got !ok\n", id, handle)
		} else if want != got {
			t.Errorf("%s: clans: %q: want %d, got %d\n", id, handle, want, got)
		}
	}
	wantTurns := []string{"0899-12"}
	if want := len(wantTurns); want != len(game.Turns) {
		t.Errorf("%s: turns: len: want %d, got %d\n", id, want, len(game.Turns))
	}
	for _, turnNo := range wantTurns {
		if turn, ok := game.Turns[turnNo]; !ok {
			t.Errorf("%s: turns: %q: want ok, got !ok\n", id, turn)
		} else if nil == turn {
			t.Errorf("%s: turns: %q: want *Turn, got nil\n", id, turn)
		} else {
			gotYearMonth := fmt.Sprintf("%04d-%02d", turn.Year, turn.Month)
			if want := turnNo; want != gotYearMonth {
				t.Errorf("%s: turns: %q: yyyy-mm: want %q, got %q\n", id, turnNo, want, gotYearMonth)
			}
		}
	}
}

func testGame0301(t *testing.T, id string, game *jsondb.Game) {
	if game == nil {
		t.Errorf("%s: want *Game, got nil\n", id)
		return
	}
	if want := id; want != game.Code {
		t.Errorf("%s: Code: want %q, got %q\n", id, want, game.Code)
	}
	if want := "TN3.1"; want != game.Description {
		t.Errorf("%s: Description: want %q, got %q\n", id, want, game.Description)
	}
	wantClans := []int{1, 999}
	if want := len(wantClans); want != len(game.Clans) {
		t.Errorf("%s: clans: len: want %d, got %d\n", id, want, len(game.Clans))
	}
	for _, no := range wantClans {
		handle, want := fmt.Sprintf("clan%04d", no), no
		if got, ok := game.Clans[handle]; !ok {
			t.Errorf("%s: clans: %q: want ok, got !ok\n", id, handle)
		} else if want != got {
			t.Errorf("%s: clans: %q: want %d, got %d\n", id, handle, want, got)
		}
	}
	wantTurns := []string{"0899-12", "0900-01"}
	if want := len(wantTurns); want != len(game.Turns) {
		t.Errorf("%s: turns: len: want %d, got %d\n", id, want, len(game.Turns))
	}
	for _, turnNo := range wantTurns {
		if turn, ok := game.Turns[turnNo]; !ok {
			t.Errorf("%s: turns: %q: want ok, got !ok\n", id, turnNo)
		} else if nil == turn {
			t.Errorf("%s: turns: %q: want *Turn, got nil\n", id, turnNo)
		} else {
			gotYearMonth := fmt.Sprintf("%04d-%02d", turn.Year, turn.Month)
			if want := turnNo; want != gotYearMonth {
				t.Errorf("%s: turns: %q: yyyy-mm: want %q, got %q\n", id, turnNo, want, gotYearMonth)
			}
			if turnNo == "0900-01" {
				wantOrdersDue := "2025-07-04 21:00:00 +0000 UTC"
				gotOrdersDue := turn.OrdersDue.String()
				if wantOrdersDue != gotOrdersDue {
					t.Errorf("%s: turns: %q: orders-due: want %q, got %q\n", id, turnNo, wantOrdersDue, gotOrdersDue)
				}
			}
		}
	}

}
