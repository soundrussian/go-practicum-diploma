package v1

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/soundrussian/go-practicum-diploma/model"
	"github.com/soundrussian/go-practicum-diploma/pkg/logging"
	balance2 "github.com/soundrussian/go-practicum-diploma/service/balance"
	"github.com/soundrussian/go-practicum-diploma/storage"
	"github.com/theplant/luhn"
	"strconv"
)

var _ balance2.Balance = (*Balance)(nil)

type Balance struct {
	storage storage.Storage
}

func New(storage storage.Storage) (*Balance, error) {
	if storage == nil {
		return nil, errors.New("nil storage passed to Balance service constructor")
	}

	auth := &Balance{storage: storage}

	return auth, nil
}

func (b *Balance) UserBalance(ctx context.Context, userID uint64) (*model.UserBalance, error) {
	res, err := b.storage.UserBalance(ctx, userID)
	if err != nil {
		b.Log(ctx).Err(err).Msg("failed to get UserBalance")
		return nil, err
	}

	return res, nil
}

func (b *Balance) Withdraw(ctx context.Context, userID uint64, withdrawal model.Withdrawal) error {
	if withdrawal.Sum <= 0 {
		return balance2.ErrInvalidSum
	}

	orderID, err := strconv.Atoi(withdrawal.Order)
	if err != nil {
		b.Log(ctx).Err(err).Msgf("failed to convert %s to integer", withdrawal.Order)
		return fmt.Errorf("%w: orderID is not a number", balance2.ErrInvalidOrder)
	}

	if !luhn.Valid(orderID) {
		b.Log(ctx).Err(err).Msgf("invalid checksum for %d", orderID)
		return fmt.Errorf("%w: orderID checksum is wrong", balance2.ErrInvalidOrder)
	}

	if _, err := b.storage.Withdraw(ctx, userID, withdrawal); err != nil {
		b.Log(ctx).Err(err).Msgf("failed to make withdrawal %+v for user %d", withdrawal, userID)
		if errors.Is(err, storage.ErrNotEnoughBalance) {
			return balance2.ErrNotEnoughBalance
		}
		return balance2.ErrInternalError
	}

	return nil
}

func (b *Balance) Withdrawals(ctx context.Context, userID uint64) ([]model.Withdrawal, error) {
	withdrawals, err := b.storage.UserWithdrawals(ctx, userID)
	if err != nil {
		return []model.Withdrawal{}, err
	}

	return withdrawals, nil
}

// Log returns logger with service field set.
func (b *Balance) Log(ctx context.Context) *zerolog.Logger {
	_, logger := logging.CtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceNameKey, "balance").Logger()

	return &logger
}
