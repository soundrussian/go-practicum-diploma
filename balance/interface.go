package balance

import "context"

type Balance interface {
	Balance(ctx context.Context, userID uint64) (uint64, error)
	Withdrawn(ctx context.Context, userID uint64) (uint64, error)
}
