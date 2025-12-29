// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package domains

import (
	"fmt"
	"time"
)

// Document represents a document that has been saved to the server.
//
// The Documents service must untaint or reject the Path.
type Document struct {
	ID             ID           // unique identifier for document
	GameID         GameID       // game associated with the document
	ClanId         ID           // clan associated with the document
	TurnNo         TurnID       // turn associated with the document
	ClanNo         int          // clan associated with the document
	UnitId         string       // unit id associated with the document
	Path           string       // Client's original name for the source
	Type           DocumentType // for categorizing on the dashboard
	ContentType    string       // for downloading
	ContentsLength int          // File size in bytes
	ContentsHash   string       // SHA-256 hash of document contents (hex, 64 chars) – used as dedupe key
	Contents       []byte       // Current file contents
	ModifiedAt     time.Time    // Document modification time, stored in UTC
	CreatedAt      time.Time    // record creation time, stored in UTC
	UpdatedAt      time.Time    // record update time, stored in UTC
}

type DocumentType string

const (
	TurnReportExtract DocumentType = "turn-report-extract"
	TurnReportFile    DocumentType = "turn-report-file"
	WorldographerMap  DocumentType = "worldographer-map"
)

func (dt DocumentType) IsValid() bool {
	switch dt {
	case TurnReportExtract:
		return true
	case TurnReportFile:
		return true
	case WorldographerMap:
		return true
	}
	return false
}

var (
	documentExtToDocumentType = map[string]DocumentType{
		"docx": TurnReportFile,
		"txt":  TurnReportExtract,
		"wxx":  WorldographerMap,
	}
	documentTypeToDocumentExt = map[DocumentType]string{
		TurnReportFile:    "docx",
		TurnReportExtract: "txt",
		WorldographerMap:  "wxx",
	}
)

func DocumentExtToDocumentType(ext string) DocumentType {
	dt, ok := documentExtToDocumentType[ext]
	if !ok {
		panic(fmt.Sprintf("assert(ext != %q)", ext))
	}
	return dt
}

func DocumentTypeToDocumentExt(dt DocumentType) string {
	ext, ok := documentTypeToDocumentExt[dt]
	if !ok {
		panic(fmt.Sprintf("assert(type != %d)", dt))
	}
	return ext
}

type MimeType string

const (
	// We can use a parameter to specify the format version when we
	// serve the Worldographer file:
	//  - Classic Map: Content-Type: application/vnd.worldographer.map; version=classic
	//  - 2025 Map: Content-Type: application/vnd.worldographer.map; version=2025

	DOCXMimeType   MimeType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	ReportMimeType MimeType = "text/vnd.tribenet-turnreport.data"
	WXXMimeType    MimeType = "application/vnd.worldographer.map"
)

func (mt MimeType) IsValid() bool {
	switch mt {
	case DOCXMimeType:
		return true
	case ReportMimeType:
		return true
	case WXXMimeType:
		return true
	}
	return false
}

type ProcessingStatus string

const (
	PSNone    ProcessingStatus = ""
	PSError   ProcessingStatus = "error"
	PSFailed  ProcessingStatus = "failed"
	PSRunning ProcessingStatus = "running"
	PSSuccess ProcessingStatus = "success"
	PSWaiting ProcessingStatus = "waiting"
)
