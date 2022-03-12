package api

import (
	"errors"
	"github.com/soundrussian/go-practicum-diploma/pkg/curruser"
	"github.com/soundrussian/go-practicum-diploma/pkg/logging"
	"github.com/soundrussian/go-practicum-diploma/service/order"
	"io"
	"net/http"
)

func (api *API) HandleOrder(w http.ResponseWriter, r *http.Request) {
	ctx, logger := logging.CtxLogger(r.Context())
	logger = logger.With().Str(logging.HandlerNameKey, "order").Logger()
	logger.Info().Msg("handling order")

	orderID, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Err(err).Msg("failed to read request body")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer r.Body.Close()

	userID, _ := curruser.CurrentUser(ctx)

	if err := api.orderService.AcceptOrder(ctx, userID, string(orderID)); err != nil {
		logger.Err(err).Msgf("error accepting order <%s> for user %d", string(orderID), userID)
		status := http.StatusInternalServerError
		if errors.Is(err, order.ErrOrderInvalid) {
			status = http.StatusUnprocessableEntity
		}
		if errors.Is(err, order.ErrConflict) {
			status = http.StatusConflict
		}
		if errors.Is(err, order.ErrAlreadyAccepted) {
			return // Implicit 200 OK
		}
		http.Error(w, http.StatusText(status), status)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
