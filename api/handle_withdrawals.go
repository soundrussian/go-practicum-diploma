package api

import (
	"encoding/json"
	"github.com/shopspring/decimal"
	"github.com/soundrussian/go-practicum-diploma/model"
	"github.com/soundrussian/go-practicum-diploma/pkg/curruser"
	"github.com/soundrussian/go-practicum-diploma/pkg/logging"
	"net/http"
	"time"
)

type withdrawalsResponse struct {
	Order       string          `json:"order"`
	Sum         decimal.Decimal `json:"sum"`
	ProcessedAt string          `json:"processed_at"`
}

func withdrawalsResponseFromModel(model model.Withdrawal) withdrawalsResponse {
	return withdrawalsResponse{
		Order:       model.Order,
		Sum:         model.Sum.RoundBank(2),
		ProcessedAt: model.ProcessedAt.Format(time.RFC3339),
	}
}

func (api *API) HandleWithdrawals(w http.ResponseWriter, r *http.Request) {
	ctx, logger := logging.CtxLogger(r.Context())
	logger = logger.With().Str(logging.HandlerNameKey, "withdrawals").Logger()
	logger.Info().Msg("handling withdrawals")

	userID, err := curruser.CurrentUser(ctx)
	if err != nil {
		logger.Err(err).Msg("failed to get current user from context")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	withdrawals, err := api.balanceService.Withdrawals(ctx, userID)
	if err != nil {
		logger.Err(err).Msg("failed to fetch withdrawals")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(withdrawals) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	response := make([]withdrawalsResponse, 0, len(withdrawals))
	for _, withdrawal := range withdrawals {
		response = append(response, withdrawalsResponseFromModel(withdrawal))
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(&response); err != nil {
		logger.Err(err).Msgf("failed to encode %+v", withdrawals)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
