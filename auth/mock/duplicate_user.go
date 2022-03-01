package mock

import "github.com/soundrussian/go-practicum-diploma/auth"

var _ auth.Auth = (*DuplicateUser)(nil)

type DuplicateUser struct {
}

func (s DuplicateUser) Register(login string, password string) (*auth.User, error) {
	return nil, auth.ErrUserAlreadyRegistered
}

func (s DuplicateUser) Authenticate(login string, password string) (*auth.User, error) {
	panic("Authenticate(login string, password string) is not implemented in DuplicateUser mock")
}
