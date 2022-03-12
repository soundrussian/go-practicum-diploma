package order

import (
	"context"
	"github.com/soundrussian/go-practicum-diploma/model"
)

type Order interface {
	AcceptOrder(ctx context.Context, userID uint64, orderID string) error
	UserOrders(ctx context.Context, userID uint64) ([]model.Order, error)
}
