package api

import (
	"context"
	"fmt"
	"github.com/soundrussian/go-practicum-diploma/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestRunServer(t *testing.T) {
	t.Run("it runs server on specified address", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		api, err := New(new(mocks.Auth), new(mocks.Balance), new(mocks.Order))
		require.NoError(t, err)

		runServerOnFreePort(t, ctx, api)

		err = pingServer()
		require.NoError(t, err)
	})
}

func TestRunServer_StopThroughContext(t *testing.T) {
	t.Run("it stops server by cancelling context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		api, err := New(new(mocks.Auth), new(mocks.Balance), new(mocks.Order))
		require.NoError(t, err)

		runServerOnFreePort(t, ctx, api)

		cancel()

		err = pingServer()
		require.Error(t, err)
	})
}

func getFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func runServerOnFreePort(t *testing.T, ctx context.Context, api *API) {
	freePort, err := getFreePort()
	require.NoError(t, err)

	config.RunAddress = fmt.Sprintf("localhost:%d", freePort)

	go func() {
		_, err = api.RunServer(ctx)
		assert.ErrorIs(t, err, http.ErrServerClosed)
	}()

	time.Sleep(500 * time.Millisecond) // Wait for server to start
}

func pingServer() error {
	timeout := 1 * time.Second
	conn, err := net.DialTimeout("tcp", config.RunAddress, timeout)
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()

	return err
}
