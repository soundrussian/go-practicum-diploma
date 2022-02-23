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

// CtxLogger stores zerolog logger and correlation ID in the passed context,
// or fetches logger already stored in the context.
//
// If nil context is given, it uses background context.
func CtxLogger(ctx context.Context) (context.Context, zerolog.Logger) {
	if ctx == nil {
		ctx = context.Background()
	}

	if ctxValue := ctx.Value(contextKeyLogger); ctxValue != nil {
		if ctxLogger, ok := ctxValue.(zerolog.Logger); ok {
			return ctx, ctxLogger
		}
	}

	var correlationID string
	var err error
	// If there is no correlation ID in the context, set it
	if correlationID, err = CorrelationID(ctx); err != nil {
		correlationUUID, _ := uuid.NewUUID()
		correlationID = correlationUUID.String()
		ctx = SetCorrelationID(ctx, correlationID)
	}

	// Create new logger and store it in the context
	logger := NewLogger().With().Str(CorrelationIDKey, correlationID).Logger()
	return SetCtxLogger(ctx, logger), logger
}

// SetCtxLogger stores logger in the ctx
func SetCtxLogger(ctx context.Context, logger zerolog.Logger) context.Context {
	return context.WithValue(ctx, contextKeyLogger, logger)
}

// CorrelationID returns correlation ID stored in the context or empty string
// and error if passed context does not contain correlation ID
func CorrelationID(ctx context.Context) (string, error) {
	id, ok := ctx.Value(contextKeyCorrelationID).(string)
	if !ok {
		return "", fmt.Errorf("correlation id not found in context")
	}

	return id, nil
}

// SetCorrelationID sets passed id as correlation ID in the given context
func SetCorrelationID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, contextKeyCorrelationID, id)
}
