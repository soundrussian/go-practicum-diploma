package v1

import (
	"context"
	"errors"
	"github.com/rs/zerolog"
	"github.com/soundrussian/go-practicum-diploma/order"
	"github.com/soundrussian/go-practicum-diploma/pkg/logging"
	"github.com/soundrussian/go-practicum-diploma/storage"
)

var _ order.Order = (*Order)(nil)

type Order struct {
	storage storage.Storage
}

func New(storage storage.Storage) (*Order, error) {
	if storage == nil {
		return nil, errors.New("nil storage passed to Order service constructor")
	}

	auth := &Order{storage: storage}

	return auth, nil
}

func (o *Order) AcceptOrder(ctx context.Context, userID uint64, orderID string) error {
	return nil
}

// Log returns logger with service field set.
func (o *Order) Log(ctx context.Context) *zerolog.Logger {
	_, logger := logging.CtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceNameKey, "order").Logger()

	return &logger
}
