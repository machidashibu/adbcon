package controller_test

import (
	"adbcon/internal/adapter/controller"
	"testing"

	"github.com/go-openapi/testify/v2/require"
)

func TestGetServerAddr(t *testing.T) {
	addr, err := controller.GetServerAddr()
	require.NoError(t, err)
	require.Equal(t, "localhost:8080", addr)
}
