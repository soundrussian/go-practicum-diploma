package storage

import "errors"

var (
	ErrLoginAlreadyExists = errors.New("login already exists")
	ErrNotFound           = errors.New("not found")
	ErrNotEnoughBalance   = errors.New("not enough balance")

	ErrOrderExistsSameUser    = errors.New("order has already been uploaded by the same user")
	ErrOrderExistsAnotherUser = errors.New("order has already been uploaded by another user")
)
