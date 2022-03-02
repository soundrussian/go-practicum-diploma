package api

import (
	"encoding/json"
	"fmt"
	"github.com/soundrussian/go-practicum-diploma/model"
	"github.com/soundrussian/go-practicum-diploma/pkg/logging"
	"net/http"
)

type authenticateJSONResponse struct {
	Token string `json:"token"`
}

func (api *API) authenticate(user *model.User, w http.ResponseWriter, r *http.Request) {
	var token *string
	var err error

	_, logger := logging.CtxLogger(r.Context())
	logger = logger.With().Str(logging.HandlerNameKey, "authenticate").Logger()

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
