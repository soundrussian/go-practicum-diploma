package api

import (
	"encoding/json"
	"errors"
	"github.com/soundrussian/go-practicum-diploma/auth"
	"github.com/soundrussian/go-practicum-diploma/pkg/logging"
	"net/http"
)

type loginJSONRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (api *API) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var jsonRequest loginJSONRequest

	_, logger := logging.CtxLogger(r.Context())
	logger = logger.With().Str(logging.HandlerNameKey, "login").Logger()
	logger.Info().Msg("handling login")

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	if err := decoder.Decode(&jsonRequest); err != nil {
		logger.Err(err).Msgf("failed to parse request body as JSON")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := api.authService.Authenticate(r.Context(), jsonRequest.Login, jsonRequest.Password)
	if err != nil {
		logger.Err(err).Msgf("failed to authenticate user %s", jsonRequest.Login)
		if errors.Is(err, auth.ErrUserNotFound) || errors.Is(err, auth.ErrPasswordIncorrect) {
			http.Error(w, "incorrect login or password", http.StatusUnauthorized)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if user == nil {
		logger.Err(errors.New("api.authService.Authenticate returned nil user")).Msgf("failed to login %+v", jsonRequest)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	logger.Info().Msgf("authenticated user id=%d with login %s", user.ID, user.Login)

	api.authenticate(user, w, r)
}
