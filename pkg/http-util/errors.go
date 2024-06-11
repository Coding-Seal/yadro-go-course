package http_util

import "errors"

var (
	ErrNotFound          = errors.New("not found")
	ErrInternal          = errors.New("internal server error")
	ErrBadRequest        = errors.New("bad request")
	ErrNoLoginOrPassword = errors.New("no login or password")
	ErrForbidden         = errors.New("forbidden")
)
