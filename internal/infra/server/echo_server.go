package server

import (
	"adbcon/api"
	"adbcon/internal/adapter/controller"
	"adbcon/internal/adapter/handler"
	"context"
	_ "embed"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//go:embed assets/cert.pem
var serverCert []byte

//go:embed assets/key.pem
var serverKey []byte

type EchoServer struct {
	srv *echo.Echo
}

func NewEchoServer() *EchoServer {
	return &EchoServer{}
}

func (s *EchoServer) Start() error {
	s.srv = echo.New()
	s.srv.Use(middleware.RequestLogger())
	s.srv.Use(middleware.Recover())

	api.RegisterHandlers(s.srv, handler.NewEchoHandler())

	addr, err := controller.GetServerAddr()
	if err != nil {
		slog.Error("server address error", "err", err)
		return err
	}

	err = s.srv.StartTLS(addr, serverCert, serverKey)
	if isCiticalError(err) {
		slog.Error("server start error", "err", err, "addr", addr)
		return err
	}

	return nil
}

func (s EchoServer) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.srv.Shutdown(ctx); err != nil {
		slog.Error("server shutdown error", "err", err)
		return err
	}
	return nil
}

func isCiticalError(err error) bool {
	if errors.Is(err, http.ErrServerClosed) || errors.Is(err, net.ErrClosed) || errors.Is(err, context.Canceled) {
		return false
	}
	return true
}
