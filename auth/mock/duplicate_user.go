package mock

import (
	"context"
	"github.com/soundrussian/go-practicum-diploma/auth"
	"github.com/soundrussian/go-practicum-diploma/model"
)

var _ auth.Auth = (*DuplicateUser)(nil)

type DuplicateUser struct {
}

func (s DuplicateUser) Register(ctx context.Context, login string, password string) (*model.User, error) {
	return nil, auth.ErrUserAlreadyRegistered
}

func (s DuplicateUser) Authenticate(ctx context.Context, login string, password string) (*model.User, error) {
	panic("Authenticate(login string, password string) is not implemented in DuplicateUser mock")
}

func (s DuplicateUser) AuthToken(ctx context.Context, user *model.User) (*string, error) {
	panic("Token(user *auth.User) is not implemented in DuplicateUser mock")
}
