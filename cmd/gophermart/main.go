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

	ctx, logger := logging.CtxLogger(nil)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serverDone, err := api.RunServer(ctx)

	if err != nil && err != http.ErrServerClosed {
		logger.Err(err).Msg("error starting server")
		return
	}

	<-serverDone
}
