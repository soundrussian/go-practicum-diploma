package api

import (
	"encoding/json"
	"github.com/soundrussian/go-practicum-diploma/pkg/logging"
	"net/http"
)

type registerJSONRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (api *API) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var jsonRequest registerJSONRequest

	_, logger := logging.CtxLogger(r.Context())
	logger.Info().Msg("handling register")

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	if err := decoder.Decode(&jsonRequest); err != nil {
		logger.Err(err).Msgf("failed to parse request body as JSON")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err := api.authService.Register(jsonRequest.Login, jsonRequest.Password); err != nil {
		logger.Err(err).Msgf("failed to register user %s", jsonRequest.Login)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
