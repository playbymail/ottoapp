// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package office

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
)

// DocXML extracts the word/document.xml from a .docx file.
func DocXML(r *bytes.Reader) ([]byte, error) {
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
			w := &bytes.Buffer{}
			_, err = io.Copy(w, rc)
			if err != nil {
				panic(err)
			}
			return w.Bytes(), nil
		}
	}
	return nil, errors.New("word/document.xml not found")
}

// DocXMLPath loads a .docx file and extracts the word/document.xml from it.
func DocXMLPath(path string) ([]byte, error) {
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
			w := &bytes.Buffer{}
			_, err = io.Copy(w, rc)
			if err != nil {
				panic(err)
			}
			return w.Bytes(), nil
		}
	}
	return nil, errors.New("word/document.xml not found")
}
