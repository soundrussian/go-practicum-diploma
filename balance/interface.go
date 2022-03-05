package balance

import "context"

type UserBalance struct {
	Current   uint64
	Withdrawn uint64
}

type Balance interface {
	UserBalance(ctx context.Context, userID uint64) (*UserBalance, error)
}
