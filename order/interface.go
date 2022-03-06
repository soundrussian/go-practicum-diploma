package order

import "context"

type Order interface {
	AcceptOrder(ctx context.Context, userID uint64, orderID string) error
}
