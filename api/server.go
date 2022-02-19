package api

import (
	"context"
	"fmt"
	"net/http"
)

func RunServer(ctx context.Context) (<-chan struct{}, error) {
	server := http.Server{Addr: config.RunAddress, Handler: http.DefaultServeMux}
	c := make(chan struct{})

	go func() {
		<-ctx.Done()

		shutdownTimeout, cancel := context.WithTimeout(context.Background(), config.ServerShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(shutdownTimeout); err != nil {
			fmt.Println(err)
		}

		c <- struct{}{}
	}()

	fmt.Println("Starting gophermart on address", config.RunAddress)
	return c, server.ListenAndServe()
}
