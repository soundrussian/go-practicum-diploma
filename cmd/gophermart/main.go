package main

import (
	"context"
	"flag"
	"github.com/soundrussian/go-practicum-diploma/accrual"
	"github.com/soundrussian/go-practicum-diploma/api"
	"github.com/soundrussian/go-practicum-diploma/pkg/logging"
	auth "github.com/soundrussian/go-practicum-diploma/service/auth/v1"
	balance "github.com/soundrussian/go-practicum-diploma/service/balance/v1"
	order "github.com/soundrussian/go-practicum-diploma/service/order/v1"
	db "github.com/soundrussian/go-practicum-diploma/storage/psql"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	ctx, logger := logging.CtxLogger(ctx)

	store, err := db.New()
	defer func() {
		if store != nil {
			store.Close()
		}
	}()

	if err != nil {
		logger.Err(err).Msg("failed to initialize storage")
		return
	}

	authService, err := auth.New(store)
	if err != nil {
		logger.Err(err).Msg("error initializing auth service")
		return
	}

	balanceService, err := balance.New(store)
	if err != nil {
		logger.Err(err).Msg("error initializing balance service")
		return
	}

	orderService, err := order.New(store)
	if err != nil {
		logger.Err(err).Msg("error initializing order service")
		return
	}

	a, err := api.New(authService, balanceService, orderService)
	if err != nil {
		logger.Err(err).Msg("error intializing API")
		return
	}

	processor, err := accrual.New(store)
	if err != nil {
		logger.Err(err).Msg("failed to start accrual processor")
	}

	processor.Run(ctx)

	serverDone, err := a.RunServer(ctx)

	if err != nil && err != http.ErrServerClosed {
		logger.Err(err).Msg("error starting server")
		return
	}

	<-serverDone
}
