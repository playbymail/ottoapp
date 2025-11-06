// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package domains

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrHashFailed          = Error("hash failed")
	ErrInvalidPath         = Error("invalid path")
	ErrMissingUserdataPath = Error("missing userdata path")
	ErrNotDirectory        = Error("not a directory")
	ErrNotExists           = Error("does not exist")
	ErrNotImplemented      = Error("not implemented")
	ErrOpenFailed          = Error("open failed")
	ErrReadFailed          = Error("read failed")
	ErrWriteFailed         = Error("write failed")
)
