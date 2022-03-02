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

type loginJSONRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (api *API) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var jsonRequest loginJSONRequest
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

	if user, err = api.authService.Authenticate(r.Context(), jsonRequest.Login, jsonRequest.Password); err != nil {
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

	var token *string
	if token, err = api.authService.AuthToken(r.Context(), user); err != nil {
		logger.Err(err).Msg("failed to get auth token for user")
		http.Error(w, "user has been registered, but failed to log in", http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:  "jwt",
		Value: *token,
	}
	http.SetCookie(w, cookie)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Authentication", fmt.Sprintf("Bearer %s", *token))

	response := authenticateJSONResponse{Token: *token}
	encoder := json.NewEncoder(w)
	if err = encoder.Encode(&response); err != nil {
		logger.Err(err).Msgf("failed to encode json response from %+v", response)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
