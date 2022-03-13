package psql

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/soundrussian/go-practicum-diploma/pkg/logging"
)

// Log returns logger with service field set.
func (s *Storage) Log(ctx context.Context) *zerolog.Logger {
	_, logger := logging.CtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceNameKey, "storage").Logger()

	return &logger
}
