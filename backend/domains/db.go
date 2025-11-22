// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package domains

const (
	ErrCreateSchema          = Error("create schema")
	ErrCreateMeta            = Error("create metadata")
	ErrDatabaseError         = Error("database error")
	ErrDatabaseExists        = Error("database exists")
	ErrForeignKeysDisabled   = Error("foreign keys disabled")
	ErrPragmaFailed          = Error("pragma failed")
	ErrPragmaReturnedNil     = Error("pragma returned nil")
	ErrSchemaVersionMismatch = Error("schema version mismatch")
)
