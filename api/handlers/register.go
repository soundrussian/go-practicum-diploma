package handlers

import (
	"github.com/soundrussian/go-practicum-diploma/pkg/logging"
	"net/http"
)

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	_, logger := logging.CtxLogger(r.Context())
	logger.Info().Msg("handling register")

	w.WriteHeader(http.StatusInternalServerError)
}
