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
	var body []byte
	var limitPerMinuteF float64
	var err error

	acc.Log(ctx).Info().Msg("handling TooManyRequests")

	if body, err = io.ReadAll(resp.Body); err != nil {
		acc.Log(ctx).Err(err).Msg("failed to read body of TooManyRequests")
		return
	}
	defer resp.Body.Close()

	re := regexp.MustCompile(`\d+`)
	limitPerMinute := re.Find(body)

	if limitPerMinuteF, err = strconv.ParseFloat(string(limitPerMinute), 64); err != nil {
		acc.Log(ctx).Err(err).Msg("failed to convert limit to float64")
		return
	}

	acc.limiter = rate.NewLimiter(rate.Limit(limitPerMinuteF/60.0), 10)
}
