// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package parser implements parsers for TribeNet reports.
package parser

import "github.com/playbymail/ottoapp/backend/domains"

type DocxParser struct {
	input []byte
}

func NewDocxParser(input []byte) (*DocxParser, error) {
	d := &DocxParser{input: input}
	return d, nil
}

func (p *DocxParser) Parse(input []byte) error {
	return domains.ErrNotImplemented
}

type ReportParser struct{}
