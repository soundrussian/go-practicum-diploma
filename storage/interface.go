package storage

import "github.com/soundrussian/go-practicum-diploma/model"

type Store interface {
	CreateUser(login string, password string) (*model.User, error)
	Close()
}
