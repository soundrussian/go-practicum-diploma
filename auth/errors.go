package auth

import "errors"

var (
	ErrInvalidLogin    = errors.New("login must be non-empty string")
	ErrInvalidPassword = errors.New("password must be non-empty string")
)
