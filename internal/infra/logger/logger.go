package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"
)

var logw io.WriteCloser

type loggerConfig interface {
	LogFilePath() string
}

// Setup set up logger for slog by log configurations.
// It merge handler to write file more than Debug level and handler to output stdout more than Info level.
func Setup(config loggerConfig) error {
	// open log file
	if err := os.MkdirAll(path.Dir(config.LogFilePath()), 888); err != nil {
		Fatal(fmt.Errorf("%s: %s", err.Error(), path.Dir(config.LogFilePath())))
		return err
	}
	f, err := os.Create(config.LogFilePath())
	if err != nil {
		Fatal(fmt.Errorf("%s: %s", err.Error(), config.LogFilePath()))
		return err
	}
	logw = f

	// prepare log handler that outputs all log to file
	handlerFile := slog.NewJSONHandler(logw, &slog.HandlerOptions{Level: slog.LevelDebug})
	// prepare log handler that outputs information log to stdout
	handlerStdout := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	// create logger with multi handler
	slog.SetDefault(slog.New(NewMultiHandler(handlerFile, handlerStdout)))

	return nil
}

// Close closes log file.
func Close() {
	if logw != nil {
		logw.Close()
		logw = nil
	}
}

// Fatal outputs error message to  stderr and return exit code 1.
func Fatal(err error) int {
	fmt.Fprintln(os.Stderr, err.Error())
	return 1
}
