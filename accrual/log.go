package accrual

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/soundrussian/go-practicum-diploma/pkg/logging"
)

// Log returns logger with service field set.
func (acc *Accrual) log(ctx context.Context) *zerolog.Logger {
	_, logger := logging.CtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceNameKey, "accrual").Logger()

	return &logger
}
