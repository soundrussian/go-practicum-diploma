package main

import (
	"context"
	"flag"
	"github.com/soundrussian/go-practicum-diploma/api"
	auth "github.com/soundrussian/go-practicum-diploma/auth/v1"
	"github.com/soundrussian/go-practicum-diploma/pkg/logging"
	storage "github.com/soundrussian/go-practicum-diploma/storage"
	db "github.com/soundrussian/go-practicum-diploma/storage/v1"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var a *api.API
	var authService *auth.Auth
	var store storage.Store
	var err error

	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	ctx, logger := logging.CtxLogger(ctx)

	store, err = db.New()
	defer func() {
		if store != nil {
			store.Close()
		}
	}()

	if err != nil {
		logger.Err(err).Msg("failed to initialize storage")
		return
	}

	if authService, err = auth.New(store); err != nil {
		logger.Err(err).Msg("error initializing auth service")
		return
	}

	if a, err = api.New(authService); err != nil {
		logger.Err(err).Msg("error intializing API")
		return
	}

	serverDone, err := a.RunServer(ctx)

	if err != nil && err != http.ErrServerClosed {
		logger.Err(err).Msg("error starting server")
		return
	}

	<-serverDone
}
