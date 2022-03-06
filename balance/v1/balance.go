package v1

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/soundrussian/go-practicum-diploma/balance"
	"github.com/soundrussian/go-practicum-diploma/model"
	"github.com/soundrussian/go-practicum-diploma/pkg/logging"
	"github.com/soundrussian/go-practicum-diploma/storage"
	"github.com/theplant/luhn"
	"strconv"
)

var _ balance.Balance = (*Balance)(nil)

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
	var res *model.UserBalance
	var err error

	if res, err = b.storage.UserBalance(ctx, userID); err != nil {
		b.Log(ctx).Err(err).Msg("failed to get UserBalance")
		return nil, err
	}

	return res, nil
}

func (b *Balance) Withdraw(ctx context.Context, userID uint64, withdrawal model.Withdrawal) error {
	if withdrawal.Sum <= 0 {
		return balance.ErrInvalidSum
	}

	var orderID int
	var err error

	if orderID, err = strconv.Atoi(withdrawal.Order); err != nil {
		b.Log(ctx).Err(err).Msgf("failed to convert %s to integer", withdrawal.Order)
		return fmt.Errorf("%w: orderID is not a number", balance.ErrInvalidOrder)
	}

	if !luhn.Valid(orderID) {
		b.Log(ctx).Err(err).Msgf("invalid checksum for %d", orderID)
		return fmt.Errorf("%w: orderID checksum is wrong", balance.ErrInvalidOrder)
	}

	if _, err = b.storage.Withdraw(ctx, userID, withdrawal); err != nil {
		b.Log(ctx).Err(err).Msgf("failed to make withdrawal %+v for user %d", withdrawal, userID)
	}

	return err
}

// Log returns logger with service field set.
func (b *Balance) Log(ctx context.Context) *zerolog.Logger {
	_, logger := logging.CtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceNameKey, "balance").Logger()

	return &logger
}
