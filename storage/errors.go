package storage

import "errors"

var (
	ErrLoginAlreadyExists = errors.New("login already exists")
	ErrNotFound           = errors.New("not found")
	ErrNotEnoughBalance   = errors.New("not enough balance")
)
