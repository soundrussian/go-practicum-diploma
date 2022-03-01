package mock

import (
	"context"
	"github.com/soundrussian/go-practicum-diploma/auth"
	"github.com/soundrussian/go-practicum-diploma/model"
)

var _ auth.Auth = (*FailedValidation)(nil)

type FailedValidation struct {
}

func (s FailedValidation) Register(ctx context.Context, login string, password string) (*model.User, error) {
	return nil, auth.ErrInvalidLogin
}

func (s FailedValidation) Authenticate(ctx context.Context, login string, password string) (*model.User, error) {
	panic("Authenticate(login string, password string) is not implemented in FailedValidation mock")
}

func (s FailedValidation) AuthToken(ctx context.Context, user *model.User) string {
	panic("AuthToken(user *auth.User) is not implemented in FailedValidation mock")
}
