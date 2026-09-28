package model_test

import (
	"adbcon/internal/adapter/model"
	"errors"
	"testing"

	"github.com/go-openapi/testify/v2/require"
)

func TestCommandResult(t *testing.T) {
	t.Run("CommandResult", func(t *testing.T) {
		r := model.NewCommandResult("serial", []byte("test result"))
		require.NotNil(t, r)
		require.Equal(t, "serial", r.Serial())
		require.Equal(t, "test result", r.Text())
		require.NoError(t, r.Error())
	})
	t.Run("CommandErrorResult", func(t *testing.T) {
		r := model.NewCommandErrorResult("serial", errors.New("test error"))
		require.NotNil(t, r)
		require.Equal(t, "serial", r.Serial())
		require.Error(t, r.Error())
		require.Equal(t, "test error", r.Text())
	})
}
