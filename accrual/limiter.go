package accrual

import (
	"context"
	"golang.org/x/time/rate"
	"io"
	"net/http"
	"regexp"
	"strconv"
)

func (acc *Accrual) HandleTooManyRequests(ctx context.Context, resp *http.Response) {
	acc.Log(ctx).Info().Msg("handling TooManyRequests")

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		acc.Log(ctx).Err(err).Msg("failed to read body of TooManyRequests")
		return
	}
	defer resp.Body.Close()

	re := regexp.MustCompile(`\d+`)
	limitPerMinute := re.Find(body)

	limitPerMinuteF, err := strconv.ParseFloat(string(limitPerMinute), 64)
	if err != nil {
		acc.Log(ctx).Err(err).Msg("failed to convert limit to float64")
		return
	}

	acc.limiter = rate.NewLimiter(rate.Limit(limitPerMinuteF/60.0), 10)
}
