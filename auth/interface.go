package auth

import (
	"context"
	"github.com/soundrussian/go-practicum-diploma/model"
)

type Auth interface {
	Register(ctx context.Context, login string, password string) (*model.User, error)
	Authenticate(ctx context.Context, login string, password string) (*model.User, error)
	AuthToken(ctx context.Context, user *model.User) (*string, error)
}
