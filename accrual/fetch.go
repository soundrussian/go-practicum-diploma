package accrual

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (acc *Accrual) Fetch(ctx context.Context, orderID string) (*Result, error) {
	var result Result

	if err := acc.limiter.Wait(ctx); err != nil {
		acc.Log(ctx).Err(err).Msg("error waiting for limiter")
	}

	acc.Log(ctx).Info().Msgf("getting accrual for order <%s>", orderID)
	resp, err := http.Get(fmt.Sprintf("%s/api/orders/%s", *accrualAddress, orderID))
	if err != nil {
		acc.Log(ctx).Err(err).Msgf("failed to fetch accrual for order <%s>", orderID)
		return nil, err
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		acc.HandleTooManyRequests(ctx, resp)
		return nil, ErrFailedToFetch
	}

	if resp.StatusCode != http.StatusOK {
		acc.Log(ctx).Error().Msgf("fetching order <%s> responded with status %d", orderID, resp.StatusCode)
		return nil, ErrFailedToFetch
	}

	decoder := json.NewDecoder(resp.Body)
	defer resp.Body.Close()

	if err := decoder.Decode(&result); err != nil {
		acc.Log(ctx).Err(err).Msg("failed to decode response from accrual service")
		return nil, err
	}

	return &result, nil
}
