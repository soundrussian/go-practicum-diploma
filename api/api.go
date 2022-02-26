package api

import (
	"errors"
	"github.com/soundrussian/go-practicum-diploma/auth"
)

type API struct {
	authService auth.Auth
}

func New(auth auth.Auth) (*API, error) {
	if auth == nil {
		return nil, errors.New("nil auth service passed to API constructor")
	}

	api := &API{
		authService: auth,
	}

	return api, nil
}
