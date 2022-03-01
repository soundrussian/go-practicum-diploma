package main

import (
	"context"
	"flag"
	"github.com/soundrussian/go-practicum-diploma/api"
	"github.com/soundrussian/go-practicum-diploma/auth/mock"
	"github.com/soundrussian/go-practicum-diploma/pkg/logging"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var a *api.API
	var err error

	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	ctx, logger := logging.CtxLogger(ctx)

	if a, err = api.New(mock.SuccessfulRegistration{}); err != nil {
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
