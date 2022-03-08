package accrual

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/soundrussian/go-practicum-diploma/model"
	"github.com/soundrussian/go-practicum-diploma/storage"
	"net/http"
	"strings"
	"sync"
)

type Accrual struct {
	storage storage.Storage
	batch   int
}

type Result struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float32 `json:"accrual"`
}

func New(store storage.Storage) (*Accrual, error) {
	if store == nil {
		return nil, errors.New("nil storage passed to Processor constructor")
	}

	if accrualAddress == nil {
		return nil, errors.New("accrualAddress has not been configured")
	}

	return &Accrual{storage: store, batch: 10}, nil
}

func (acc *Accrual) Tick(ctx context.Context) error {
	var orders []string
	var err error

	if orders, err = acc.NextBatch(ctx); err != nil {
		acc.Log(ctx).Err(err).Msg("failed to get next batch of records to process")
		return err
	}

	if len(orders) == 0 {
		acc.Log(ctx).Info().Msg("no orders to process")
		return nil
	}

	var wg sync.WaitGroup
	for _, order := range orders {
		wg.Add(1)
		go func(order string) {
			defer func() { wg.Done() }()
			if err := acc.Process(ctx, order); err != nil {
				acc.Log(ctx).Err(err).Msgf("error processing order <%s>", order)
			}
		}(order)
	}
	wg.Wait()

	return nil
}

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

	// Mark order as processing so that it won't go into next batch
	// while being processed
	if err = acc.storage.UpdateOrderStatus(ctx, orderID, model.OrderProcessing); err != nil {
		acc.Log(ctx).Err(err).Msgf("failed to mark order <%s> as processing", orderID)
		return err
	}

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
		return nil
	}

	if strings.ToLower(res.Status) == "processed" {
		acc.Log(ctx).Info().Msgf("order <%s> has been processed with accrual %f", orderID, res.Accrual)
		if err = acc.storage.UpdateOrderStatus(ctx, orderID, model.OrderProcessed); err != nil {
			acc.Log(ctx).Err(err).Msgf("failed to mark order <%s> as processed", orderID)
			return err
		}
		return nil
	}

	// If we are this far, it means that the order is not processed by accrual yet,
	// So mark it as new
	if err = acc.storage.UpdateOrderStatus(ctx, orderID, model.OrderNew); err != nil {
		acc.Log(ctx).Err(err).Msgf("failed to mark order <%s> as new", orderID)
		return err
	}

	return nil
}

func (acc *Accrual) Fetch(ctx context.Context, orderID string) (*Result, error) {
	var result Result

	acc.Log(ctx).Info().Msgf("getting accrual for order <%s>", orderID)
	resp, err := http.Get(fmt.Sprintf("%s/api/orders/%s", *accrualAddress, orderID))
	if err != nil {
		acc.Log(ctx).Err(err).Msgf("failed to fetch accrual for order <%s>", orderID)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		acc.Log(ctx).Error().Msgf("fetching order <%s> responded with status %d", orderID, resp.StatusCode)
		return nil, ErrFailedToFetch
	}

	decoder := json.NewDecoder(resp.Body)
	defer resp.Body.Close()

	if err = decoder.Decode(&result); err != nil {
		acc.Log(ctx).Err(err).Msg("failed to decode response from accrual service")
		return nil, err
	}

	return &result, nil
}
