package accrual

import (
	"context"
	"github.com/soundrussian/go-practicum-diploma/accrual/status"
	"github.com/soundrussian/go-practicum-diploma/model"
)

func (acc *Accrual) nextBatch(ctx context.Context) ([]string, error) {
	orders, err := acc.storage.OrdersWithStatus(ctx, model.OrderNew, acc.batch)
	if err != nil {
		acc.log(ctx).Err(err).Msg("failed to fetch next portion of orders to process")
		return nil, err
	}

	return orders, nil
}

func (acc *Accrual) process(ctx context.Context, orderID string) error {
	var finalResult bool

	// Mark order as processing so that it won't go into next batch
	// while being processed
	if err := acc.storage.UpdateOrderStatus(ctx, orderID, model.OrderProcessing); err != nil {
		acc.log(ctx).Err(err).Msgf("failed to mark order <%s> as processing", orderID)
		return err
	}

	// Check if we got final result in defer. If not, make order NEW again
	defer func() {
		if !finalResult {
			if err := acc.storage.UpdateOrderStatus(ctx, orderID, model.OrderNew); err != nil {
				acc.log(ctx).Err(err).Msgf("failed to mark order <%s> as new", orderID)
			}
		}
	}()

	res, err := acc.fetch(ctx, orderID)
	if err != nil {
		acc.log(ctx).Err(err).Msgf("failed to process order <%s>", orderID)
		return err
	}

	acc.log(ctx).Info().Msgf("got response from accrual service: %+v", res)

	switch status.New(res.Status) {
	case status.Invalid:
		acc.log(ctx).Info().Msgf("order <%s> has been marked as invalid", orderID)
		if err := acc.storage.UpdateOrderStatus(ctx, orderID, model.OrderInvalid); err != nil {
			acc.log(ctx).Err(err).Msgf("failed to mark order <%s> as invalid", orderID)
			return err
		}
		finalResult = true
		return nil
	case status.Processed:
		acc.log(ctx).Info().Msgf("order <%s> has been processed with accrual %s", orderID, res.Accrual)
		if err := acc.storage.AddAccrual(ctx, orderID, model.OrderProcessed, res.Accrual); err != nil {
			acc.log(ctx).Err(err).Msgf("failed to mark order <%s> as processed", orderID)
			return err
		}
		finalResult = true
		return nil
	}

	return nil
}
