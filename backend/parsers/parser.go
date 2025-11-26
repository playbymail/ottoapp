// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package parser implements parsers for TribeNet reports.
package parsers

import (
	"bytes"

	"github.com/playbymail/ottoapp/backend/parsers/office"
)

type Docx struct {
	Text []byte
}

func ParseDocx(r *bytes.Reader) (*Docx, error) {
	doc, err := office.Parse(r)
	if err != nil {
		return nil, err
	}
	return &Docx{Text: doc.Text}, nil
}

// ParseClanHeading returns the first heading in the report or an error.
//
// Note: hard coded values to use while developing the GM upload page.
// We will implement a real parser soon.
func ParseClanHeading(doc *Docx) (*ElementHeader_t, error) {
	if doc == nil {
		return nil, ErrBadInput
	}
	if !bytes.HasPrefix(doc.Text, []byte(`Tribe 0`)) {
		return nil, ErrNotATurnReport
	}
	return &ElementHeader_t{
		Id: "0511",
		Turn: &Turn_t{
			Year:  899,
			Month: 12,
			No:    0,
		},
	}, nil
}

type ReportParser struct{}

type ElementHeader_t struct {
	Id   string
	Turn *Turn_t
}

type Turn_t struct {
	Year  int
	Month int
	No    int
}
