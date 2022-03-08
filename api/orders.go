package api

import (
	"github.com/soundrussian/go-practicum-diploma/model"
	"github.com/soundrussian/go-practicum-diploma/pkg/curruser"
	"github.com/soundrussian/go-practicum-diploma/pkg/logging"
	"net/http"
)

func (api *API) HandleOrders(w http.ResponseWriter, r *http.Request) {
	ctx, logger := logging.CtxLogger(r.Context())
	logger = logger.With().Str(logging.HandlerNameKey, "orders").Logger()
	logger.Info().Msg("handling orders")

	userID, _ := curruser.CurrentUser(ctx)

	var orders []model.Order
	var err error

	if orders, err = api.orderService.UserOrders(ctx, userID); err != nil {
		logger.Err(err).Msgf("failed getting orders for user %d", userID)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
}
