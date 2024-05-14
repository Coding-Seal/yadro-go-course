package ports

import "errors"

var (
	ErrNotFound = errors.New("comic not found")
	ErrInternal = errors.New("internal error")
)
