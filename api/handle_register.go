package api

import (
	"encoding/json"
	"errors"
	"github.com/soundrussian/go-practicum-diploma/pkg/logging"
	"github.com/soundrussian/go-practicum-diploma/service/auth"
	"net/http"
)

type registerJSONRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (api *API) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var jsonRequest registerJSONRequest

	_, logger := logging.CtxLogger(r.Context())
	logger = logger.With().Str(logging.HandlerNameKey, "register").Logger()
	logger.Info().Msg("handling register")

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	if err := decoder.Decode(&jsonRequest); err != nil {
		logger.Err(err).Msgf("failed to parse request body as JSON")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := api.authService.Register(r.Context(), jsonRequest.Login, jsonRequest.Password)
	if err != nil {
		logger.Err(err).Msgf("failed to register user %s", jsonRequest.Login)
		status := http.StatusBadRequest
		if errors.Is(err, auth.ErrUserAlreadyRegistered) {
			status = http.StatusConflict
		}
		http.Error(w, err.Error(), status)
		return
	}

	if user == nil {
		logger.Err(errors.New("api.authService.Register returned nil user")).Msgf("failed to register %+v", jsonRequest)
		http.Error(w, "registration failed", http.StatusInternalServerError)
		return
	}

	logger.Info().Msgf("saved user id=%d with login %s", user.ID, user.Login)

	api.authenticate(user, w, r)
}
