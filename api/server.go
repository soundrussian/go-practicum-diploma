package api

import (
	"context"
	"github.com/soundrussian/go-practicum-diploma/pkg/logging"
	"net/http"
)

func RunServer(ctx context.Context) (<-chan struct{}, error) {
	ctx, logger := logging.CtxLogger(ctx)
	server := http.Server{Addr: config.RunAddress, Handler: routes()}
	c := make(chan struct{})

	go func() {
		<-ctx.Done()

		shutdownTimeout, cancel := context.WithTimeout(context.Background(), config.ServerShutdownTimeout)
		defer cancel()

		logger.Info().Msg("waiting for all connections to shut down...")
		
		if err := server.Shutdown(shutdownTimeout); err != nil {
			logger.Err(err).Msg("error while shutting down server")
		}

		c <- struct{}{}
	}()

	logger.Info().Msgf("starting server on address %s", config.RunAddress)

	return c, server.ListenAndServe()
}
