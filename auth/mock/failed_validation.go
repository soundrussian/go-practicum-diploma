package mock

import "github.com/soundrussian/go-practicum-diploma/auth"

var _ auth.Auth = (*FailedValidation)(nil)

type FailedValidation struct {
}

func (s FailedValidation) Register(login string, password string) (*auth.User, error) {
	return nil, auth.ErrInvalidLogin
}

func (s FailedValidation) Authenticate(login string, password string) (*auth.User, error) {
	panic("Authenticate(login string, password string) is not implemented in SuccessfulRegistration mock")
}
