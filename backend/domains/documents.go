// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package domains

import "time"

// Document represents a document that has been saved to the server.
type Document struct {
	ID            ID     // unique identifier for document
	Path          string // Client's original name for the source (tainted!)
	MimeType      string
	ContentLength int64     // File size in bytes
	Contents      []byte    // Current file contents
	ContentsHash  string    // SHA-256 hash of document contents (hex, 64 chars) – used as dedupe key
	CreatedAt     time.Time // Document creation time, stored in UTC
	UpdatedAt     time.Time // Document modification time, stored in UTC
}
