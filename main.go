package main

import (
	"adbcon/internal/adapter/handler"
	"adbcon/internal/infra/config"
	"adbcon/internal/infra/database"
	"adbcon/internal/infra/logger"
	"adbcon/internal/infra/server"
	"context"
	"log/slog"
	"os"
	"os/signal"
)

func main() {
	// read configuration
	cfg := new(config.Config)
	if err := cfg.Read("config.yaml"); err != nil {
		os.Exit(logger.Fatal(err))
	}

	// setup logger
	if err := logger.Setup(cfg); err != nil {
		os.Exit(logger.Fatal(err))
	}
	defer logger.Close()

	// run
	os.Exit(run(cfg))
}

func run(cfg *config.Config) int {
	slog.Info("start application")

	// create application context
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// prepare database
	dbStatus := new(database.StubDatabase)

	// start API server
	apiServer := server.NewEchoServer()
	apiError := make(chan error, 1)
	go func() {
		if err := apiServer.Start(ctx, cfg, handler.Factory(dbStatus, cfg)); err != nil {
			apiError <- err
		}
	}()

	// wait CTRL+C
	select {
	case <-ctx.Done(): // CTRL+C
		slog.Info("post processing...")
		// shutdown server
		if err := apiServer.Shutdown(ctx); err != nil {
			return logger.Fatal(err)
		}
	case err := <-apiError:
		return logger.Fatal(err)
	}
	slog.Info("terminate application")

	return 0
}
