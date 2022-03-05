package mock

import (
	"context"
	"github.com/soundrussian/go-practicum-diploma/balance"
)

var _ balance.Balance = (*BalanceMock)(nil)

type BalanceMock struct {
	ByUser map[uint64]balance.UserBalance
}

func (b BalanceMock) UserBalance(ctx context.Context, userID uint64) (*balance.UserBalance, error) {
	var bal balance.UserBalance
	var ok bool

	if bal, ok = b.ByUser[userID]; !ok {
		return nil, balance.ErrUserNotFound
	}

	return &bal, nil
}
