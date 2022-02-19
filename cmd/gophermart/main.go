package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/soundrussian/go-practicum-diploma/api"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serverDone, err := api.RunServer(ctx)

	if err != nil && err != http.ErrServerClosed {
		fmt.Println(err)
		return
	}

	<-serverDone
}
