package main

import (
	"context"
	"fmt"
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
	logw, err := os.Create("adbcom.log")
	if err != nil {
		fmt.Println("ERROR: Failed to create log file: ", err)
		panic(err)
	}
	defer logw.Close()
	slog.SetDefault(slog.New(slog.NewJSONHandler(logw, &slog.HandlerOptions{Level: slog.LevelDebug})))

	// run
	os.Exit(run())
}

func run() int {
	// create application context
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	slog.Info("start application")
	// TODO

	// wait CTRL+C
	<-ctx.Done()
	slog.Info("terminate application")

	return 0
}
