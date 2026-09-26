package server

import (
	"adbcon/api"
	"adbcon/internal/adapter/handler"
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	oapimiddleware "github.com/oapi-codegen/echo-middleware"
)

type echoServerConfig interface {
	ServerBind() string
	ServerPort() string
}

type EchoServer struct {
	srv *echo.Echo
}

func NewEchoServer() *EchoServer {
	return &EchoServer{
		srv: echo.New(),
	}
}

func (s *EchoServer) Start(ctx context.Context, config echoServerConfig, h api.ServerInterface) error {
	slog.Info("start echo server")

	s.srv.Use(middleware.RequestLogger())

	// apply error handler
	s.srv.HTTPErrorHandler = handler.EchoErrorHandler
	s.srv.Use(middleware.Recover())

	// apply middleware for validator with Open API definition
	spec, err := api.GetSpec()
	if err != nil {
		slog.Error("get spec error", "err", err)
		return err
	}
	s.srv.Use(oapimiddleware.OapiRequestValidator(spec))

	// apply middleware for assets that provides static files
	s.srv.Pre(AssetsWithConfig(AssetsConfig{
		Assets: assetsTable,
	}))

	api.RegisterHandlers(s.srv, h)

	// apply application context
	s.srv.Server.BaseContext = func(net.Listener) context.Context {
		return ctx
	}

	addr := net.JoinHostPort(config.ServerBind(), config.ServerPort())
	if err := s.srv.StartTLS(addr, serverCert, serverKey); isCiticalError(err) {
		slog.Error("server start error", "err", err, "addr", addr)
		return err
	}

	return nil
}

func (s EchoServer) Shutdown(ctx context.Context) error {
	slog.Info("shutdown echo server")

	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()
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
