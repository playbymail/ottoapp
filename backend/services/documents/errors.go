// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package documents

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrInvalidPath = Error("invalid path")
)
