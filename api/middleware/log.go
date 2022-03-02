package middleware

import (
	"bytes"
	"github.com/rs/zerolog/log"
	"github.com/soundrussian/go-practicum-diploma/pkg/logging"
	"io/ioutil"
	"net/http"
	"strings"
)

type logEntry struct {
	body   string
	method string
	addr   string
	path   string
}

func fromRequest(r *http.Request) logEntry {

	le := logEntry{
		body:   r.Header.Get("user-agent"),
		method: r.Method,
		addr:   r.RemoteAddr,
		path:   r.URL.Path,
	}

	return le
}

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, logger := logging.CtxLogger(r.Context())

		entry := fromRequest(r)

		buf, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error().Msgf("[%s] %s Error reading request body: %v", strings.ToUpper(entry.method), entry.path, err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		entry.body = string(buf)

		// Make sure body can be read again in next handlers
		// See https://gist.github.com/briangershon/fa9feb08e6a65d52bdc35c738d8cf104?permalink_comment_id=3719593#gistcomment-3719593
		reader := ioutil.NopCloser(bytes.NewBuffer(buf))
		r.Body = reader

		logger = logger.With().
			Str("method", entry.method).
			Str("addr", entry.addr).
			Str("path", entry.path).
			Str("body", entry.body).
			Logger()

		logger.Info().Msgf("[%s] %s", strings.ToUpper(entry.method), entry.path)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
