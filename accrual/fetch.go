package accrual

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	"net/http"
)

func (acc *Accrual) fetch(ctx context.Context, orderID string) (*result, error) {
	var result result

	acc.log(ctx).Info().Msgf("getting accrual for order <%s>", orderID)

	resp, err := retryablehttp.Get(fmt.Sprintf("%s/api/orders/%s", *accrualAddress, orderID))
	if err != nil {
		acc.log(ctx).Err(err).Msgf("failed to fetch accrual for order <%s>", orderID)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		acc.log(ctx).Error().Msgf("fetching order <%s> responded with status %d", orderID, resp.StatusCode)
		return nil, ErrFailedToFetch
	}

	decoder := json.NewDecoder(resp.Body)
	defer resp.Body.Close()

	if err := decoder.Decode(&result); err != nil {
		acc.log(ctx).Err(err).Msg("failed to decode response from accrual service")
		return nil, err
	}

	return &result, nil
}
