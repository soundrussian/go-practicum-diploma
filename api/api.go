package api

import (
	"errors"
	"github.com/soundrussian/go-practicum-diploma/service/auth"
	"github.com/soundrussian/go-practicum-diploma/service/balance"
	"github.com/soundrussian/go-practicum-diploma/service/order"
)

type API struct {
	authService    auth.Auth
	balanceService balance.Balance
	orderService   order.Order
}

func New(auth auth.Auth, balance balance.Balance, order order.Order) (*API, error) {
	if auth == nil {
		return nil, errors.New("nil auth service passed to API constructor")
	}
	if balance == nil {
		return nil, errors.New("nil balance service passed to API constructor")
	}
	if order == nil {
		return nil, errors.New("nil order service passed to API constructor")
	}

	api := &API{
		authService:    auth,
		balanceService: balance,
		orderService:   order,
	}

	return api, nil
}
