package api

import (
	"encoding/json"
	"errors"
	"github.com/shopspring/decimal"
	"github.com/soundrussian/go-practicum-diploma/model"
	"github.com/soundrussian/go-practicum-diploma/pkg/curruser"
	"github.com/soundrussian/go-practicum-diploma/pkg/logging"
	"github.com/soundrussian/go-practicum-diploma/service/balance"
	"net/http"
)

type withdrawJSONRequest struct {
	Order string          `json:"order"`
	Sum   decimal.Decimal `json:"sum"`
}

func (api *API) HandleWithdraw(w http.ResponseWriter, r *http.Request) {
	var jsonRequest withdrawJSONRequest

	ctx, logger := logging.CtxLogger(r.Context())
	logger = logger.With().Str(logging.HandlerNameKey, "withdraw").Logger()
	logger.Info().Msg("handling withdraw")

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	if err := decoder.Decode(&jsonRequest); err != nil {
		logger.Err(err).Msgf("failed to parse request body as JSON")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := curruser.CurrentUser(ctx)
	if err != nil {
		logger.Err(err).Msg("failed to get current user from context")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err := api.balanceService.Withdraw(ctx, userID, requestToWithdrawal(jsonRequest)); err != nil {
		logger.Err(err).Msgf("failed to withdraw %s point from user %d for order %s", jsonRequest.Sum, userID, jsonRequest.Order)
		if errors.Is(err, balance.ErrNotEnoughBalance) {
			http.Error(w, err.Error(), http.StatusPaymentRequired)
			return
		}
		if errors.Is(err, balance.ErrInvalidSum) {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		if errors.Is(err, balance.ErrInvalidOrder) {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func requestToWithdrawal(req withdrawJSONRequest) model.Withdrawal {
	return model.Withdrawal{
		Order: req.Order,
		Sum:   req.Sum,
	}
}
