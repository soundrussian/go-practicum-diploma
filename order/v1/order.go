package v1

import (
	"context"
	"errors"
	"github.com/rs/zerolog"
	"github.com/soundrussian/go-practicum-diploma/model"
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
	ord := model.Order{OrderID: orderID}
	if err := ord.Validate(); err != nil {
		o.Log(ctx).Err(err).Msgf("failed validating order %s", orderID)
		return order.ErrOrderInvalid
	}

	if _, err := o.storage.AcceptOrder(ctx, userID, orderID); err != nil {
		o.Log(ctx).Err(err).Msgf("error while storing order <%s> from user %d", orderID, userID)
		if errors.Is(err, storage.ErrOrderExistsSameUser) {
			return order.ErrAlreadyAccepted
		}
		if errors.Is(err, storage.ErrOrderExistsAnotherUser) {
			return order.ErrConflict
		}

		return err
	}

	return nil
}

// Log returns logger with service field set.
func (o *Order) Log(ctx context.Context) *zerolog.Logger {
	_, logger := logging.CtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceNameKey, "order").Logger()

	return &logger
}