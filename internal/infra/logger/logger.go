package logger

import (
	"log/slog"
	"os"
)

type loggerConfig interface {
	LogFilePath() string
}

// Setup set up logger for slog by log configurations.
// It merge handler to write file more than Debug level and handler to output stdout more than Info level.
func Setup(config loggerConfig) error {
	// open log file
	f, err := os.Create(config.LogFilePath())
	if err != nil {
		slog.Error("logger setup open error", "err", err, "path", config.LogFilePath())
		return err
	}

	// prepare log handler that outputs all log to file
	handlerFile := slog.NewJSONHandler(f, &slog.HandlerOptions{Level: slog.LevelDebug})
	// prepare log handler that outputs information log to stdout
	handlerStdout := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	// create logger with multi handler
	slog.SetDefault(slog.New(NewMultiHandler(handlerFile, handlerStdout)))

	return nil
}
