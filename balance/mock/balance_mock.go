package mock

import (
	"context"
	"github.com/soundrussian/go-practicum-diploma/balance"
)

var _ balance.Balance = (*BalanceMock)(nil)

type Balance struct {
	Balance   uint64
	Withdrawn uint64
}

type BalanceMock struct {
	ByUser map[uint64]Balance
}

func (b BalanceMock) Balance(ctx context.Context, userID uint64) (uint64, error) {
	var bal Balance
	var ok bool

	if bal, ok = b.ByUser[userID]; !ok {
		return 0, balance.ErrUserNotFound
	}

	return bal.Balance, nil
}

func (b BalanceMock) Withdrawn(ctx context.Context, userID uint64) (uint64, error) {
	var bal Balance
	var ok bool

	if bal, ok = b.ByUser[userID]; !ok {
		return 0, balance.ErrUserNotFound
	}

	return bal.Withdrawn, nil
}

