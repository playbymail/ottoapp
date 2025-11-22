// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package domains

import "time"

// Document represents a document that has been saved to the server.
//
// The Documents service must untaint or reject the Path.
type Document struct {
	ID            ID           // unique identifier for document
	ClanId        ID           // clan associated with the document
	Path          string       // Client's original name for the source
	Type          DocumentType // for categorizing on the dashboard
	MimeType      MimeType     // for downloading
	ContentLength int64        // File size in bytes
	Contents      []byte       // Current file contents
	ContentsHash  string       // SHA-256 hash of document contents (hex, 64 chars) – used as dedupe key
	CanDelete     bool
	CanRead       bool
	CanShare      bool
	CanWrite      bool
	CreatedAt     time.Time // Document creation time, stored in UTC
	UpdatedAt     time.Time // Document modification time, stored in UTC
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
