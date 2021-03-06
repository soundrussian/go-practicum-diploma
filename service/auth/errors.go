package auth

import "errors"

var (
	ErrInvalidLogin              = errors.New("login must be non-empty string")
	ErrInvalidPassword           = errors.New("password must be non-empty string")
	ErrUserAlreadyRegistered     = errors.New("user already registered")
	ErrRegistrationInternalError = errors.New("failed to register user")
	ErrUserNotFound              = errors.New("user not found")
	ErrAuthenticateInternalError = errors.New("internal error while authenticating")
	ErrPasswordIncorrect         = errors.New("incorrect password")
)
