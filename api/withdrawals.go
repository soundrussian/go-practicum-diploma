package api

import (
	"encoding/json"
	"github.com/soundrussian/go-practicum-diploma/model"
	"github.com/soundrussian/go-practicum-diploma/pkg/curruser"
	"github.com/soundrussian/go-practicum-diploma/pkg/logging"
	"net/http"
)

type withdrawalsJSONResponse []model.Withdrawal

func (api *API) HandleWithdrawals(w http.ResponseWriter, r *http.Request) {
	ctx, logger := logging.CtxLogger(r.Context())
	logger = logger.With().Str(logging.HandlerNameKey, "withdrawals").Logger()
	logger.Info().Msg("handling withdrawals")

	userID, _ := curruser.CurrentUser(ctx)

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

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(&withdrawals); err != nil {
		logger.Err(err).Msgf("failed to encode %+v", withdrawals)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
