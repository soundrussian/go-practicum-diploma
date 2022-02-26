package mock

import "github.com/soundrussian/go-practicum-diploma/auth"

var _ auth.Auth = (*SuccessfulRegistration)(nil)

type SuccessfulRegistration struct {
}

func (s SuccessfulRegistration) Register(login string, password string) (*auth.User, error) {
	user := auth.User{
		ID:    100,
		Login: login,
	}

	return &user, nil
}

func (s SuccessfulRegistration) Authenticate(login string, password string) (*auth.User, error) {
	panic("Authenticate(login string, password string) is not implemented in SuccessfulRegistration mock")
}

