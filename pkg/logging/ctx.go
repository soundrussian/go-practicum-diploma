package logging

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/soundrussian/go-practicum-diploma/pkg"
)

const (
	contextKeyLogger        = pkg.ContextKey("Logger")
	contextKeyCorrelationID = pkg.ContextKey("CorrelationID")
)

func CtxLogger(ctx context.Context) (context.Context, zerolog.Logger) {
	if ctx == nil {
		ctx = context.Background()
	}

	if ctxValue := ctx.Value(contextKeyLogger); ctxValue != nil {
		if ctxLogger, ok := ctxValue.(zerolog.Logger); ok {
			return ctx, ctxLogger
		}
	}

	correlationID, _ := uuid.NewUUID()
	logger := NewLogger().With().Str(CorrelationIDKey, correlationID.String()).Logger()
	ctx = context.WithValue(ctx, contextKeyCorrelationID, correlationID.String())

	return SetCtxLogger(ctx, logger), logger
}

func SetCtxLogger(ctx context.Context, logger zerolog.Logger) context.Context {
	return context.WithValue(ctx, contextKeyLogger, logger)
}

func CorrelationID(ctx context.Context) (string, error) {
	id, ok := ctx.Value(contextKeyCorrelationID).(string)
	if !ok {
		return "", fmt.Errorf("correlation id not found in context")
	}

	return id, nil
}

func SetCorrelationID(ctx context.Context, id string) (context.Context, zerolog.Logger) {
	ctx, logger := CtxLogger(ctx)
	logger = logger.With().Str(CorrelationIDKey, id).Logger()

	ctx = context.WithValue(ctx, CorrelationIDKey, id)

	return SetCtxLogger(ctx, logger), logger
}
