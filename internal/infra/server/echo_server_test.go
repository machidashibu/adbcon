package server_test

import (
	"adbcon/internal/infra/server"
	"net"
	"testing"
	"time"

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

// utility: wait server is established
func waitListening(t *testing.T, address string) {
	t.Helper()

	deadline := time.Now().Add(3 * time.Second)

	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", address, 100*time.Millisecond)
		if err == nil {
			_ = conn.Close()
			return
		}

		time.Sleep(10 * time.Millisecond)
	}

	t.Fatalf("server did not start listening on %s", address)
}

func TestEchoServer(t *testing.T) {
	// testing: server start
	srv := server.NewEchoServer()
	require.NotNil(t, srv)
	go func() {
		require.NoError(t, srv.Start(stubConfig{port: "8080"}))
	}()

	// wait 1st server is established
	waitListening(t, "127.0.0.1:8080")

	// testing: confrict port
	srv2nd := server.NewEchoServer()
	require.NotNil(t, srv2nd)
	require.Error(t, srv2nd.Start(stubConfig{port: "8080"}))

	// testing: server shutdown
	require.NoError(t, srv.Shutdown())
}
