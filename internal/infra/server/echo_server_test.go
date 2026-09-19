package server_test

import (
	"adbcon/internal/infra/server"
	"testing"

	"github.com/go-openapi/testify/v2/require"
)

func TestEchoServer(t *testing.T) {
	srv := server.NewEchoServer()
	require.NotNil(t, srv)
	go func() {
		require.NoError(t, srv.Start())
	}()
	require.NoError(t, srv.Shutdown())
}
