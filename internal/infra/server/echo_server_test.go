package server_test

import (
	"adbcon/internal/infra/server"
	"testing"

	"github.com/go-openapi/testify/v2/require"
)

type stubConfig struct {
	port string
}

func (c stubConfig) ServerBind() string {
	return ""
}

func (c stubConfig) ServerPort() string {
	return c.port
}

func TestEchoServer(t *testing.T) {
	// testing: server start
	srv := server.NewEchoServer()
	require.NotNil(t, srv)
	go func() {
		require.NoError(t, srv.Start(stubConfig{port: "8080"}))
	}()

	// testing: confrict port
	srv2nd := server.NewEchoServer()
	require.NotNil(t, srv2nd)
	require.Error(t, srv2nd.Start(stubConfig{port: "8080"}))

	// testing: server shutdown
	require.NoError(t, srv.Shutdown())
}
