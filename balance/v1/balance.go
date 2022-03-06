package v1

import (
	"context"
	"errors"
	"github.com/rs/zerolog"
	"github.com/soundrussian/go-practicum-diploma/balance"
	"github.com/soundrussian/go-practicum-diploma/model"
	"github.com/soundrussian/go-practicum-diploma/pkg/logging"
	"github.com/soundrussian/go-practicum-diploma/storage"
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

// Log returns logger with service field set.
func (b *Balance) Log(ctx context.Context) *zerolog.Logger {
	_, logger := logging.CtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceNameKey, "balance").Logger()

	return &logger
}
