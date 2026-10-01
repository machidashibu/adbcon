package logger_test

import (
	"adbcon/internal/infra/config"
	"adbcon/internal/infra/logger"
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/go-openapi/testify/v2/require"
)

type stubConfig struct{}

func TestSetup(t *testing.T) {
	// remove previous test result
	os.RemoveAll("testdata/test.log")

	t.Run("normal", func(t *testing.T) {
		// testing: close
		require.NotPanics(t, func() { logger.Close() })

		// testing: open
		require.NoError(t, logger.Setup(
			config.Config{Log: config.LogConfig{FilePath: "testdata/test.log"}},
		))
		defer logger.Close()

		// testing: write
		slog.Debug("test debug message")
		slog.Info("test info message")
		slog.Warn("test warning message")
		slog.Error("test error message")

		// check
		require.FileExists(t, "testdata/test.log")
		require.FileNotEmpty(t, "testdata/test.log")
	})

	t.Run("file open error", func(t *testing.T) {
		require.Error(t, logger.Setup(
			config.Config{Log: config.LogConfig{FilePath: "testdata"}},
		))
	})

	t.Run("mk dir all fail", func(t *testing.T) {
		require.Error(t, logger.Setup(
			config.Config{Log: config.LogConfig{FilePath: `testdata/test.log/tst.log`}},
		))
	})
}

func TestFatal(t *testing.T) {
	require.Equal(t, 1, logger.Fatal(errors.New("test error")))
}
