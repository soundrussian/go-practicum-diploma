package middleware

import (
	"github.com/soundrussian/go-practicum-diploma/pkg/logging"
	"net/http"
	"strings"
)

type logEntry struct {
	userAgent string
	method    string
	addr      string
	path      string
}

func fromRequest(r *http.Request) logEntry {
	le := logEntry{
		userAgent: r.Header.Get("user-agent"),
		method:    r.Method,
		addr:      r.RemoteAddr,
		path:      r.URL.Path,
	}

	return le
}

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, logger := logging.CtxLogger(r.Context())

		logEntry := fromRequest(r)
		logger = logger.With().
			Str("method", logEntry.method).
			Str("addr", logEntry.addr).
			Str("path", logEntry.path).
			Str("ua", logEntry.userAgent).
			Logger()

		logger.Info().Msgf("[%s] %s", strings.ToUpper(logEntry.method), logEntry.path)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
