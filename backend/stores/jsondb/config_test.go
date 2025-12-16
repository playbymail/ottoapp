// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package jsondb_test

import (
	"testing"

	"github.com/playbymail/ottoapp/backend/stores/jsondb"
)

func TestLoadOttoAppConfig(t *testing.T) {
	oac, err := jsondb.LoadOttoAppConfig("testdata/config.json")
	if err != nil {
		t.Fatal(err)
	}
	if want := "ottoapp.example.com"; want != oac.Mailgun.Domain {
		t.Errorf("oac: mailgun: domain: want %q, got %q\n", want, oac.Mailgun.Domain)
	}
	if want := "ottoapp@ottoapp.example.com"; want != oac.Mailgun.From {
		t.Errorf("oac: mailgun: from: want %q, got %q\n", want, oac.Mailgun.From)
	}
	if want := "https://api.mailgun.net/v3"; want != oac.Mailgun.ApiBase {
		t.Errorf("oac: mailgun: from: want %q, got %q\n", want, oac.Mailgun.ApiBase)
	}
	if want := "not-an-api-key"; want != oac.Mailgun.ApiKey {
		t.Errorf("oac: mailgun: from: want %q, got %q\n", want, oac.Mailgun.ApiKey)
	}
}
