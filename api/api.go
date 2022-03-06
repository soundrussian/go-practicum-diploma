package api

import (
	"errors"
	"github.com/soundrussian/go-practicum-diploma/auth"
	"github.com/soundrussian/go-practicum-diploma/balance"
)

type API struct {
	authService    auth.Auth
	balanceService balance.Balance
}

func New(auth auth.Auth, balance balance.Balance) (*API, error) {
	if auth == nil {
		return nil, errors.New("nil auth service passed to API constructor")
	}
	if balance == nil {
		return nil, errors.New("nil balance service passed to API constructor")
	}

	api := &API{
		authService:    auth,
		balanceService: balance,
	}

	return api, nil
}
