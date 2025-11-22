// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package restapi

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrInvalidPageNumber = Error("invalid page number")
	ErrInvalidPageSize   = Error("invalid page size")
)
