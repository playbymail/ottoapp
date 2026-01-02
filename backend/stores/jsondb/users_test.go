// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package jsondb_test

import (
	"testing"

	"github.com/playbymail/ottoapp/backend/stores/jsondb"
)

func TestLoadUsers(t *testing.T) {
	users, err := jsondb.LoadUsers("testdata/users.json")
	if err != nil {
		t.Fatal(err)
	}
	if want := 4; want != len(users) {
		t.Errorf("len: want %d, got %d\n", want, len(users))
	}

	if user, ok := users["catbird"]; !ok {
		t.Errorf("catbird: want *User, got nil\n")
	} else {
		testUserCatbird(t, user)
	}

	if user, ok := users["frojo"]; !ok {
		t.Errorf("frojo: want *User, got nil\n")
	} else {
		testUserFrojo(t, user)
	}

	if user, ok := users["penguin"]; !ok {
		t.Errorf("penguin: want *User, got nil\n")
	} else {
		testUserPenguin(t, user)
	}

	if user, ok := users["sambo"]; !ok {
		t.Errorf("sambo: want *User, got nil\n")
	} else {
		testUserSambo(t, user)
	}
}

func TestLoadUser(t *testing.T) {
	if user, err := jsondb.LoadUser("testdata/users.json", "catbird"); err != nil {
		t.Fatal(err)
	} else {
		testUserCatbird(t, user)
	}
	if user, err := jsondb.LoadUser("testdata/users.json", "frojo"); err != nil {
		t.Fatal(err)
	} else {
		testUserFrojo(t, user)
	}
	if user, err := jsondb.LoadUser("testdata/users.json", "penguin"); err != nil {
		t.Fatal(err)
	} else {
		testUserPenguin(t, user)
	}
	if user, err := jsondb.LoadUser("testdata/users.json", "sambo"); err != nil {
		t.Fatal(err)
	} else {
		testUserSambo(t, user)
	}
}

func testUserCatbird(t *testing.T, user *jsondb.User) {
	if user == nil {
		t.Errorf("catbird: want *User, got nil\n")
		return
	}
	if want := "catbird"; want != user.Handle {
		t.Errorf("catbird: Handle: want %q, got %q\n", want, user.Handle)
	}
	if want := "catbird@ottoapp"; want != user.Email {
		t.Errorf("catbird: Email: want %q, got %q\n", want, user.Email)
	}
	if nil == user.Tz {
		t.Errorf("catbird: Tz: want *time.Location, got nil\n")
	} else if want := "America/Denver"; want != user.Tz.String() {
		t.Errorf("catbird: Tz: want %q, got %q\n", want, user.Tz.String())
	}
	if doNotWant := ""; doNotWant == user.Password.Password {
		t.Errorf("catbird: Password: do not want %q, got %q\n", doNotWant, user.Password.Password)
	}
	if want := true; want != user.Password.UpdatePassword {
		t.Errorf("catbird: Password.UpdatePassword: want %v, got %v\n", want, user.Password.UpdatePassword)
	}
	if want := true; want != user.Roles.Active {
		t.Errorf("catbird: Roles.Active: want %v, got %v\n", want, user.Roles.Active)
	}
}

func testUserFrojo(t *testing.T, user *jsondb.User) {
	if user == nil {
		t.Errorf("frojo: want *User, got nil\n")
		return
	}
	if want := "frojo"; want != user.Handle {
		t.Errorf("frojo: Handle: want %q, got %q\n", want, user.Handle)
	}
	if want := "frojo@ottoapp"; want != user.Email {
		t.Errorf("frojo: Email: want %q, got %q\n", want, user.Email)
	}
	if nil == user.Tz {
		t.Errorf("frojo: Tz: want *time.Location, got nil\n")
	} else if want := "America/Chicago"; want != user.Tz.String() {
		t.Errorf("frojo: Tz: want %q, got %q\n", want, user.Tz.String())
	}
	if doNotWant := "change-me"; doNotWant == user.Password.Password {
		t.Errorf("frojo: Password: do not want %q, got %q\n", doNotWant, user.Password.Password)
	}
	if want := true; want != user.Password.UpdatePassword {
		t.Errorf("frojo: Password.UpdatePassword: want %v, got %v\n", want, user.Password.UpdatePassword)
	}
	if want := false; want != user.Roles.Active {
		t.Errorf("frojo: Roles.Active: want %v, got %v\n", want, user.Roles.Active)
	}
}

func testUserPenguin(t *testing.T, user *jsondb.User) {
	if user == nil {
		t.Errorf("penguin: want *User, got nil\n")
		return
	}
	if want := "penguin"; want != user.Handle {
		t.Errorf("penguin: Handle: want %q, got %q\n", want, user.Handle)
	}
	if want := "penguin@ottoapp"; want != user.Email {
		t.Errorf("penguin: Email: want %q, got %q\n", want, user.Email)
	}
	if nil == user.Tz {
		t.Errorf("penguin: Tz: want *time.Location, got nil\n")
	} else if want := "Antarctica/Palmer"; want != user.Tz.String() {
		t.Errorf("penguin: Tz: want %q, got %q\n", want, user.Tz.String())
	}
	if want := "happy hoppy chevy levy"; want != user.Password.Password {
		t.Errorf("penguin: Password: want %q, got %q\n", want, user.Password.Password)
	}
	if want := false; want != user.Password.UpdatePassword {
		t.Errorf("penguin: Password.UpdatePassword: want %v, got %v\n", want, user.Password.UpdatePassword)
	}
	if want := true; want != user.Roles.Active {
		t.Errorf("penguin: Roles.Active: want %v, got %v\n", want, user.Roles.Active)
	}
}

func testUserSambo(t *testing.T, user *jsondb.User) {
	if user == nil {
		t.Errorf("sambo: want *User, got nil\n")
		return
	}
	if want := "sambo"; want != user.Handle {
		t.Errorf("sambo: Handle: want %q, got %q\n", want, user.Handle)
	}
	if want := "sambo@ottoapp"; want != user.Email {
		t.Errorf("sambo: Email: want %q, got %q\n", want, user.Email)
	}
	if nil == user.Tz {
		t.Errorf("sambo: Tz: want *time.Location, got nil\n")
	} else if want := "America/Chicago"; want != user.Tz.String() {
		t.Errorf("sambo: Tz: want %q, got %q\n", want, user.Tz.String())
	}
	if want := "update-me"; want != user.Password.Password {
		t.Errorf("sambo: Password: want %q, got %q\n", want, user.Password.Password)
	}
	if want := true; want != user.Password.UpdatePassword {
		t.Errorf("sambo: Password.UpdatePassword: want %v, got %v\n", want, user.Password.UpdatePassword)
	}
	if want := true; want != user.Roles.Active {
		t.Errorf("sambo: Roles.Active: want %v, got %v\n", want, user.Roles.Active)
	}
}
