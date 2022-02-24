package main

import (
	"context"
	"flag"
	"github.com/soundrussian/go-practicum-diploma/api"
	"github.com/soundrussian/go-practicum-diploma/pkg/logging"
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
	serverDone, err := api.RunServer(ctx)

	if err != nil && err != http.ErrServerClosed {
		logger.Err(err).Msg("error starting server")
		return
	}

	<-serverDone
}
