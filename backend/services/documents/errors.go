// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package documents

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrExists      = Error("document already exists")
	ErrInvalidPath = Error("invalid path")
)
