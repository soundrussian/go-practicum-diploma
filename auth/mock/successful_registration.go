package mock

import (
	"context"
	"github.com/soundrussian/go-practicum-diploma/auth"
	v1 "github.com/soundrussian/go-practicum-diploma/auth/v1"
	"github.com/soundrussian/go-practicum-diploma/model"
)

var _ auth.Auth = (*Successful)(nil)

const Password = "topsecret"
const UserID = 100

type Successful struct {
}

func (s Successful) Register(ctx context.Context, login string, password string) (*model.User, error) {
	user := model.User{
		ID:    UserID,
		Login: login,
	}

	return &user, nil
}

func (s Successful) Authenticate(ctx context.Context, login string, password string) (*model.User, error) {
	if password != Password {
		return nil, auth.ErrInvalidPassword
	}

	user := model.User{
		ID:    UserID,
		Login: login,
	}

	return &user, nil
}

func (s Successful) AuthToken(ctx context.Context, user *model.User) (*string, error) {
	token := Token()
	return &token, nil
}

func Token() string {
	a := &v1.Auth{}
	token, _ := a.AuthToken(context.Background(), &model.User{ID: UserID})
	return *token
}
