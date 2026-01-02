// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package status_test

import (
	"testing"

	"github.com/playbymail/ottoapp/backend/parsers/azul/status"
)

func TestNext(t *testing.T) {
	input := "0987 Status: GRASSY HILLS, River S E 0987"
	rest := []byte(input)
	for _, tc := range []struct {
		id   int
		want status.Token
	}{
		{1, status.Token{Lexeme: "0987", Trivia: " "}},
		{2, status.Token{Lexeme: "Status"}},
		{3, status.Token{Lexeme: ":", Trivia: " "}},
		{4, status.Token{Lexeme: "GRASSY", Trivia: " "}},
		{5, status.Token{Lexeme: "HILLS"}},
		{6, status.Token{Lexeme: ",", Trivia: " "}},
		{7, status.Token{Lexeme: "River", Trivia: " "}},
		{8, status.Token{Lexeme: "S", Trivia: " "}},
		{9, status.Token{Lexeme: "E", Trivia: " "}},
		{10, status.Token{Lexeme: "0987"}},
	} {
		var got status.Token
		got, rest = status.Next(rest)
		if tc.want.Lexeme != got.Lexeme {
			t.Errorf("%d: lexeme: want %q: got %q\n", tc.id, tc.want.Lexeme, got.Lexeme)
		}
		if tc.want.Trivia != got.Trivia {
			t.Errorf("%d: trivia: want %q: got %q\n", tc.id, tc.want.Trivia, got.Trivia)
		}
	}
}

func TestSplit(t *testing.T) {
	input := "ab: cd, ef h,i 0"
	var left, right []byte
	for n, tc := range []struct {
		id          int
		left, right string
	}{
		{id: 1, left: "ab", right: ": cd, ef h,i 0"},
		{id: 2, left: ":", right: " cd, ef h,i 0"},
		{id: 3, left: " ", right: "cd, ef h,i 0"},
		{id: 4, left: "cd", right: ", ef h,i 0"},
		{id: 5, left: ",", right: " ef h,i 0"},
		{id: 6, left: " ", right: "ef h,i 0"},
		{id: 7, left: "ef", right: " h,i 0"},
		{id: 8, left: " ", right: "h,i 0"},
		{id: 9, left: "h", right: ",i 0"},
		{id: 10, left: ",", right: "i 0"},
		{id: 11, left: "i", right: " 0"},
		{id: 12, left: " ", right: "0"},
		{id: 13, left: "0", right: ""},
		{id: 14, left: "", right: ""},
	} {
		if n == 0 {
			left, right = status.Split([]byte(input))
		} else {
			left, right = status.Split(right)
		}
		if tc.left != string(left) {
			t.Errorf("%d:  left: want %q: got %q\n", tc.id, tc.left, left)
		}
		if tc.right != string(right) {
			t.Errorf("%d: right: want %q: got %q\n", tc.id, tc.right, right)
		}
	}
}
