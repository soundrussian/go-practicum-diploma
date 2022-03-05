package api

import (
	"github.com/soundrussian/go-practicum-diploma/pkg/logging"
	"net/http"
)

func (api *API) HandleBalance(w http.ResponseWriter, r *http.Request) {
	_, logger := logging.CtxLogger(r.Context())
	logger = logger.With().Str(logging.HandlerNameKey, "balance").Logger()
	logger.Info().Msg("handling balance")
}
