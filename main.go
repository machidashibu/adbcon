package main

import (
	"adbcon/internal/adapter/handler"
	"adbcon/internal/infra/config"
	"adbcon/internal/infra/database"
	"adbcon/internal/infra/server"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
)

func abort(err error) int {
	slog.Error(err.Error())
	return 1
}

func main() {
	// Setup debug logging
	logf, err := os.Create("adbcom.log")
	if err != nil {
		fmt.Println("ERROR: Failed to create log file: ", err)
		panic(err)
	}
	defer logf.Close()
	logw := io.MultiWriter(os.Stdout, logf)
	slog.SetDefault(slog.New(slog.NewJSONHandler(logw, &slog.HandlerOptions{Level: slog.LevelDebug})))

	// run
	os.Exit(run())
}

func run() int {
	// read configuration
	cfg := new(config.Config)
	if err := cfg.Read("config.yaml"); err != nil {
		slog.Warn("config file open error", "err", err)
		// do not abort here, use default configrations.
	}

	slog.Info("start application")

	// create application context
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// parepa database
	dbStatus := new(database.StubDatabase)

	// start API server
	apiServer := server.NewEchoServer()
	apiError := make(chan error, 1)
	go func() {
		if err := apiServer.Start(ctx, cfg, handler.Factory(dbStatus)); err != nil {
			apiError <- err
		}
	}()

	// wait CTRL+C
	select {
	case <-ctx.Done(): // CTRL+C
		slog.Info("post processing...")
		// shutdown server
		if err := apiServer.Shutdown(ctx); err != nil {
			return abort(err)
		}
	case err := <-apiError:
		return abort(err)
	}
	slog.Info("terminate application")

	return 0
}
