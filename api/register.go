package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/soundrussian/go-practicum-diploma/auth"
	"github.com/soundrussian/go-practicum-diploma/model"
	"github.com/soundrussian/go-practicum-diploma/pkg/logging"
	"net/http"
)

type registerJSONRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type registerJSONResponse struct {
	Token string `json:"token"`
}

func (api *API) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var jsonRequest registerJSONRequest
	var user *model.User
	var err error

	_, logger := logging.CtxLogger(r.Context())
	logger = logger.With().Str(logging.HandlerNameKey, "register").Logger()
	logger.Info().Msg("handling register")

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	if err = decoder.Decode(&jsonRequest); err != nil {
		logger.Err(err).Msgf("failed to parse request body as JSON")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if user, err = api.authService.Register(r.Context(), jsonRequest.Login, jsonRequest.Password); err != nil {
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

	token := api.authService.AuthToken(r.Context(), user)

	cookie := &http.Cookie{
		Name:  "token",
		Value: token,
	}
	http.SetCookie(w, cookie)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Authentication", fmt.Sprintf("Bearer %s", token))

	response := registerJSONResponse{Token: token}
	encoder := json.NewEncoder(w)
	if err = encoder.Encode(&response); err != nil {
		logger.Err(err).Msgf("failed to encode json response from %+v", response)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
