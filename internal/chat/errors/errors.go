package errors

import "errors"

var (
	NotAuthorizedError = errors.New("not authorized")
	NotFoundError      = errors.New("not found")
	AlreadyExistsError = errors.New("already exists")
	BannedError        = errors.New("user is banned")
)
