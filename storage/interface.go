package storage

import (
	"context"
	"github.com/soundrussian/go-practicum-diploma/model"
)

type Storage interface {
	CreateUser(ctx context.Context, login string, password string) (*model.User, error)
	FetchUser(ctx context.Context, login string) (*model.User, error)
	UserBalance(ctx context.Context, userID uint64) (*model.UserBalance, error)
	Withdraw(ctx context.Context, userID uint64, withdrawal model.Withdrawal) (*model.Withdrawal, error)
	UserWithdrawals(ctx context.Context, userID uint64) ([]model.Withdrawal, error)
	AcceptOrder(ctx context.Context, userID uint64, orderID string) (*model.Order, error)
	UserOrders(ctx context.Context, userID uint64) ([]model.Order, error)
	OrdersWithStatus(ctx context.Context, status model.OrderStatus, limit int) ([]string, error)
	UpdateOrder(ctx context.Context, orderID string, status model.OrderStatus, accrual float32) (*model.Order, error)
	Close()
}
