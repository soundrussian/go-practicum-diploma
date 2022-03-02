package mock

import (
	"context"
	"github.com/soundrussian/go-practicum-diploma/model"
	"github.com/soundrussian/go-practicum-diploma/storage"
)

var _ storage.Store = (*MemoryStorage)(nil)

type MemoryStorage struct {
}

func (m MemoryStorage) CreateUser(ctx context.Context, login string, password string) (*model.User, error) {
	user := model.User{ID: 1, Login: login}
	return &user, nil
}

func (m MemoryStorage) FetchUser(ctx context.Context, login string) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (m MemoryStorage) Close() {

}
