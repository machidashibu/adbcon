package handler

import "github.com/labstack/echo/v4"

type EchoHandler struct{}

func NewEchoHandler() *EchoHandler {
	return &EchoHandler{}
}

// GetMainPage Get GUI
// (GET /gui)
func (h *EchoHandler) GetMainPage(ctx echo.Context) error {
	// NOTE: This endpoint is not called, because proceeds by statics middleware.
	return echo.ErrNotFound
}
