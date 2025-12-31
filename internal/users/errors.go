package users

import "errors"

var (
	ErrNotFound      = errors.New("user not found")
	ErrDuplicateEmail = errors.New("email already in use")
)
