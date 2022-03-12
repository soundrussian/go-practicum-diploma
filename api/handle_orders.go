package api

import (
	"encoding/json"
	"github.com/soundrussian/go-practicum-diploma/pkg/curruser"
	"github.com/soundrussian/go-practicum-diploma/pkg/logging"
	"net/http"
)

func (api *API) HandleOrders(w http.ResponseWriter, r *http.Request) {
	ctx, logger := logging.CtxLogger(r.Context())
	logger = logger.With().Str(logging.HandlerNameKey, "orders").Logger()
	logger.Info().Msg("handling orders")

	userID, err := curruser.CurrentUser(ctx)
	if err != nil {
		logger.Err(err).Msg("failed to get current user from context")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	orders, err := api.orderService.UserOrders(ctx, userID)
	if err != nil {
		logger.Err(err).Msgf("failed getting orders for user %d", userID)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(&orders); err != nil {
		logger.Err(err).Msgf("failed to encode json response from %+v", orders)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
