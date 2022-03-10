// Package accrual fetches records to process from storage,
// makes requests to external service and updates user balance
// with data received from that external service
package accrual

import (
	"context"
	"errors"
	"github.com/soundrussian/go-practicum-diploma/storage"
	"sync"
	"time"
)

// Accrual contains settings for async accrual processing
type Accrual struct {
	// storage is storage service to read and write from DB
	storage storage.Storage
	// batch - how many records to process in one tick
	batch int
	// interval - how often to check for new records to process
	interval time.Duration
}

// Result contains order status received from external service
type Result struct {
	// Order is ID of processed order
	Order string `json:"order"`
	// Status contains status of accrual, see external service docs
	Status string `json:"status"`
	// Accrual contains reward points granted for the given order (in rubles)
	Accrual float64 `json:"accrual"`
}

// New initializes new Accrual processor. Returns err in passed storage is nil
// or accrualAddress of the external service is not set.
func New(store storage.Storage) (*Accrual, error) {
	if store == nil {
		return nil, errors.New("nil storage passed to Processor constructor")
	}

	if accrualAddress == nil {
		return nil, errors.New("accrualAddress has not been configured")
	}

	// batch and interval are not configurable in the current version
	return &Accrual{storage: store, batch: 10, interval: time.Second}, nil
}

// Run spins up timer ticking every acc.interval.
// On each tick processor fetches new batch of records from storage,
// makes async requests to external service and updates accrual data
// accordingly. It is stopped if passed context is done.
func (acc *Accrual) Run(ctx context.Context) {
	timer := time.NewTicker(acc.interval)

	go func() {
		for {
			select {
			case <-timer.C:
				// Each tick runs with its own context
				tickCtx := context.Background()
				if err := acc.Tick(tickCtx); err != nil {
					acc.Log(tickCtx).Err(err).Msg("error during processor tick")
				}
			case <-ctx.Done():
				timer.Stop()
				acc.Log(ctx).Info().Msg("shutting down processor")
				return
			}
		}
	}()
}

// Tick fetches next batch or records to process,
// fires up goroutines to check each of them
// and waits for completion of each.
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
			if err := acc.Process(context.Background(), order); err != nil {
				acc.Log(ctx).Err(err).Msgf("error processing order <%s>", order)
			}
		}(order)
	}
	wg.Wait()

	return nil
}
