// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package office implements a parser to read a .docx file and
// return a []byte with whitespace preserved, mostly. We inject
// a line-feed at the end of every paragraph to help parsing.
package office

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"
)

// Parse a reader that's loaded a .docx file. Returns the body text with
// whitespace preserved from <w:t xml:space="preserve"> plus tabs and
// line breaks. Injects a line-feed at the end of every paragraph.
func Parse(r *bytes.Reader) ([]byte, error) {
	zr, err := zip.NewReader(r, r.Size())
	if err != nil {
		return nil, fmt.Errorf("open zip: %w", err)
	}
	for _, file := range zr.File {
		if file.Name == "word/document.xml" {
			rc, err := file.Open()
			if err != nil {
				return nil, fmt.Errorf("open document.xml: %w", err)
			}
			defer rc.Close()
			return parseWordXML(rc)
		}
	}
	return nil, errors.New("word/document.xml not found")
}

// ParsePath is a helper function. It opens and parses a .docx file. Returns
// the body text with whitespace preserved from <w:t xml:space="preserve">
// plus tabs and line breaks.  Injects a line-feed at the end of every paragraph.
func ParsePath(path string) ([]byte, error) {
	zr, err := zip.OpenReader(path)
	if err != nil {
		return nil, fmt.Errorf("open docx: %w", err)
	}
	defer zr.Close()
	for _, file := range zr.File {
		if file.Name == "word/document.xml" {
			rc, err := file.Open()
			if err != nil {
				return nil, fmt.Errorf("open document.xml: %w", err)
			}
			defer rc.Close()
			return parseWordXML(rc)
		}
	}
	return nil, errors.New("word/document.xml not found")
}

// ParseBufferPreserveWhitespace from a reader that's loaded a .docx file and returns the body text
// with whitespace preserved from <w:t xml:space="preserve"> plus tabs and line breaks.
func ParseBufferPreserveWhitespace(r *bytes.Reader) ([]byte, error) {
	zr, err := zip.NewReader(r, r.Size())
	if err != nil {
		return nil, fmt.Errorf("open zip: %w", err)
	}
	for _, file := range zr.File {
		if file.Name == "word/document.xml" {
			rc, err := file.Open()
			if err != nil {
				return nil, fmt.Errorf("open document.xml: %w", err)
			}
			defer rc.Close()
			return parseWordXML(rc)
		}
	}
	return nil, errors.New("word/document.xml not found")
}

// parseWordXML actually walks the XML and builds the text.
func parseWordXML(r io.Reader) ([]byte, error) {
	dec := xml.NewDecoder(r)

	var (
		buf         = &bytes.Buffer{}
		inP         = false // we're inside a paragraph
		inT         = false // we're inside a <w:t>
		preserve    = false // xml:space="preserve" on current <w:t>
		seenTextInP = false // to know if paragraph has any text
	)
	_, _ = inP, seenTextInP

	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("xml decode: %w", err)
		}

		switch se := tok.(type) {
		case xml.StartElement:
			local := se.Name.Local

			switch local {
			case "p":
				inP = true
				seenTextInP = false

			case "t":
				inT = true
				preserve = false
				// check attributes for xml:space="preserve"
				for _, a := range se.Attr {
					if (a.Name.Local == "space" || strings.HasSuffix(a.Name.Local, "space")) &&
						a.Value == "preserve" {
						preserve = true
						break
					}
				}

			case "tab":
				buf.WriteByte('\t')
				seenTextInP = true

			case "br":
				buf.WriteByte('\n')
				seenTextInP = true
			}

		case xml.EndElement:
			local := se.Name.Local
			switch local {
			case "t":
				inT = false
				preserve = false
			case "p":
				// inject a line-feed at the end of every paragraph
				buf.WriteByte('\n')
				inP = false
				seenTextInP = false
			}

		case xml.CharData:
			if inT {
				text := string(se)
				if preserve {
					buf.WriteString(text)
				} else {
					// most Word text here is already “normal”, but just in case,
					// don’t trim aggressively — just write it.
					buf.WriteString(text)
				}
				seenTextInP = true
			}
		}
	}

	return buf.Bytes(), nil
}
