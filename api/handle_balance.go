package api

import (
	"encoding/json"
	"github.com/shopspring/decimal"
	"github.com/soundrussian/go-practicum-diploma/model"
	"github.com/soundrussian/go-practicum-diploma/pkg/curruser"
	"github.com/soundrussian/go-practicum-diploma/pkg/logging"
	"net/http"
)

type balanceJSONResponse struct {
	Current   decimal.Decimal `json:"current"`
	Withdrawn decimal.Decimal `json:"withdrawn"`
}

func (api *API) HandleBalance(w http.ResponseWriter, r *http.Request) {
	ctx, logger := logging.CtxLogger(r.Context())
	logger = logger.With().Str(logging.HandlerNameKey, "balance").Logger()
	logger.Info().Msg("handling balance")

	userID, err := curruser.CurrentUser(r.Context())
	if err != nil {
		logger.Err(err).Msg("failed to get current user from context")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userBalance, err := api.balanceService.UserBalance(ctx, userID)
	if err != nil {
		logger.Err(err).Msg("failed to get balance for user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	resp := respFromModel(userBalance)

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(&resp); err != nil {
		logger.Err(err).Msgf("failed to encode json response from %+v", resp)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func respFromModel(balance *model.UserBalance) balanceJSONResponse {
	return balanceJSONResponse{
		Current:   balance.Current.RoundBank(2),
		Withdrawn: balance.Withdrawn.RoundBank(2),
	}
}
