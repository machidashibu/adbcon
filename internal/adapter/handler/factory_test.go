package handler_test

import (
	"adbcon/internal/adapter/handler"
	"adbcon/internal/infra/database"
	"testing"

	"github.com/go-openapi/testify/v2/require"
)

func TestCreateEchoServer(t *testing.T) {
	require.NotNil(t, handler.Factory(new(database.StubDatabase)))
}
