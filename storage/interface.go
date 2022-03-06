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
	Close()
}
