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
		logger.Fatal().Msgf("failed to initialize storage: %s", err.Error())
		return
	}

	authService, err := auth.New(store)
	if err != nil {
		logger.Fatal().Msgf("error initializing auth service: %s", err.Error())
		return
	}

	balanceService, err := balance.New(store)
	if err != nil {
		logger.Fatal().Msgf("error initializing balance service: %s", err.Error())
		return
	}

	orderService, err := order.New(store)
	if err != nil {
		logger.Fatal().Msgf("error initializing order service: %s", err.Error())
		return
	}

	a, err := api.New(authService, balanceService, orderService)
	if err != nil {
		logger.Fatal().Msgf("error intializing API: %s", err.Error())
		return
	}

	processor, err := accrual.New(store)
	if err != nil {
		logger.Fatal().Msgf("failed to start accrual processor: %s", err.Error())
		return
	}

	processor.Run(ctx)

	serverDone, err := a.RunServer(ctx)

	if err != nil && err != http.ErrServerClosed {
		logger.Fatal().Msgf("error starting server: %s", err)
		return
	}

	<-serverDone
}
