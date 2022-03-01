package mock

import (
	"context"
	"github.com/soundrussian/go-practicum-diploma/auth"
	"github.com/soundrussian/go-practicum-diploma/model"
)

var _ auth.Auth = (*SuccessfulRegistration)(nil)

const AuthToken = "123456qwer"

type SuccessfulRegistration struct {
}

func (s SuccessfulRegistration) Register(ctx context.Context, login string, password string) (*model.User, error) {
	user := model.User{
		ID:    100,
		Login: login,
	}

	return &user, nil
}

func (s SuccessfulRegistration) Authenticate(ctx context.Context, login string, password string) (*model.User, error) {
	panic("Authenticate(login string, password string) is not implemented in SuccessfulRegistration mock")
}

func (s SuccessfulRegistration) AuthToken(ctx context.Context, user *model.User) string {
	return AuthToken
}
