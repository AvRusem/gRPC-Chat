package errors

import "errors"

var (
	ErrNotAuthorized = errors.New("not authorized")
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
	ErrBanned        = errors.New("user is banned")
)
