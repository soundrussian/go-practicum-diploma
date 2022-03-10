package accrual

import (
	"context"
	"github.com/soundrussian/go-practicum-diploma/model"
	"strings"
)

func (acc *Accrual) NextBatch(ctx context.Context) ([]string, error) {
	var orders []string
	var err error

	if orders, err = acc.storage.OrdersWithStatus(ctx, model.OrderNew, acc.batch); err != nil {
		acc.Log(ctx).Err(err).Msg("failed to fetch next portion of orders to process")
		return nil, err
	}

	return orders, nil
}

func (acc *Accrual) Process(ctx context.Context, orderID string) error {
	var res *Result
	var err error
	var finalResult bool

	// Mark order as processing so that it won't go into next batch
	// while being processed
	if err = acc.storage.UpdateOrderStatus(ctx, orderID, model.OrderProcessing); err != nil {
		acc.Log(ctx).Err(err).Msgf("failed to mark order <%s> as processing", orderID)
		return err
	}

	// Check if we got final result in defer. If not, make order NEW again
	defer func() {
		if !finalResult {
			if err = acc.storage.UpdateOrderStatus(ctx, orderID, model.OrderNew); err != nil {
				acc.Log(ctx).Err(err).Msgf("failed to mark order <%s> as new", orderID)
			}
		}
	}()

	if res, err = acc.Fetch(ctx, orderID); err != nil {
		acc.Log(ctx).Err(err).Msgf("failed to process order <%s>", orderID)
		return err
	}

	acc.Log(ctx).Info().Msgf("got response from accrual service: %+v", res)

	if strings.ToLower(res.Status) == "invalid" {
		acc.Log(ctx).Info().Msgf("order <%s> has been marked as invalid", orderID)
		if err = acc.storage.UpdateOrderStatus(ctx, orderID, model.OrderInvalid); err != nil {
			acc.Log(ctx).Err(err).Msgf("failed to mark order <%s> as invalid", orderID)
			return err
		}
		finalResult = true
		return nil
	}

	if strings.ToLower(res.Status) == "processed" {
		acc.Log(ctx).Info().Msgf("order <%s> has been processed with accrual %f", orderID, res.Accrual)
		if err = acc.storage.AddAccrual(ctx, orderID, model.OrderProcessed, res.Accrual); err != nil {
			acc.Log(ctx).Err(err).Msgf("failed to mark order <%s> as processed", orderID)
			return err
		}
		finalResult = true
		return nil
	}

	return nil
}
