package storage

import (
	"context"
	"github.com/soundrussian/go-practicum-diploma/model"
)

type Store interface {
	CreateUser(ctx context.Context, login string, password string) (*model.User, error)
	Close()
}
