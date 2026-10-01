package logger_test

import (
	"adbcon/internal/infra/logger"
	"bytes"
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/go-openapi/testify/v2/require"
)

type stubFaileHandler struct {
	slog.JSONHandler
}

func (h *stubFaileHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return true
}

func (h *stubFaileHandler) Handle(ctx context.Context, r slog.Record) error {
	return errors.ErrUnsupported
}

func (h *stubFaileHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *stubFaileHandler) WithGroup(name string) slog.Handler {
	return h
}

func TestMultiHandler(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		outInfo := bytes.Buffer{}
		outError := bytes.Buffer{}
		mh := logger.NewMultiHandler(
			slog.NewJSONHandler(&outInfo, &slog.HandlerOptions{Level: slog.LevelInfo}),
			slog.NewJSONHandler(&outError, &slog.HandlerOptions{Level: slog.LevelError}),
		)
		slog.SetDefault(slog.New(mh).WithGroup("groupNormal").With(slog.String("keyNormal", "valueNormal")))

		slog.Debug("test debug message", "key1", "value1")
		slog.Info("test info message", "key2", "value2")
		slog.Warn("test warn message", "key3", "value3")
		slog.Error("test error message", "key4", "value4")

		// TODO: check output
	})
	t.Run("fail Handle()", func(t *testing.T) {
		outInfo := bytes.Buffer{}
		mh := logger.NewMultiHandler(
			slog.NewJSONHandler(&outInfo, &slog.HandlerOptions{Level: slog.LevelInfo}),
			&stubFaileHandler{},
		)
		slog.SetDefault(slog.New(mh))

		require.Error(t, mh.Handle(context.TODO(), slog.Record{Level: slog.LevelError, Message: "test error message"}))
	})
}
