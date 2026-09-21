package handler

import (
	"adbcon/api"
	"adbcon/internal/adapter/controller"
	"adbcon/internal/adapter/presenter"
	"adbcon/internal/adapter/presenter/apiconv"
	"adbcon/internal/usecase"
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type EchoHandler struct {
	ucAdbDevices *usecase.ExecuteAdbDevicesUsecase
}

func NewEchoHandler(ucAdbDevices *usecase.ExecuteAdbDevicesUsecase) *EchoHandler {
	return &EchoHandler{
		ucAdbDevices: ucAdbDevices,
	}
}

// GetMainPage Get GUI
// (GET /gui)
func (h *EchoHandler) GetMainPage(ctx echo.Context) error {
	// NOTE: This endpoint is not called, because proceeds by statics middleware.
	return echo.ErrNotFound
}

// GetDevices Execute adb devices command.
// (GET /api/devices)
func (h *EchoHandler) GetDevices(ctx echo.Context, params api.GetDevicesParams) error {
	// validate
	if err := params.Validate(); err != nil {
		return ctx.JSON(http.StatusBadRequest, apiconv.MakeBadRequest(err))
	}

	// to domain
	interval := controller.IntervalToDomain(*params.Interval)

	// prepare reporter
	reporter := presenter.NewSSEReporter(ctx.Response())

	// call usecase: execute command (immediate update current status)
	updated, err := h.ucAdbDevices.Execute(ctx.Request().Context())
	if err != nil {
		slog.Error("adb devices usecase error", "err", err)
		return nil // disconnect from server
		// return ctx.JSON(http.StatusInternalServerError, apiconv.MakeInternalServerError(err))
	}
	// report to client
	if err := reporter.ReportDeviceList(updated); err != nil {
		slog.Error("device list report error", "err", err, "updated", updated)
		return nil // disconnect from server
		// return ctx.JSON(http.StatusInternalServerError, apiconv.MakeInternalServerError(err))
	}

	// start monitoring
	polling := time.NewTicker(time.Duration(interval))
	defer polling.Stop()
	for {
		select {
		case <-ctx.Request().Context().Done():
			return nil // disconnect
		case <-polling.C:
			// call usecase: execute command (periodical update status)
			updated, err := h.ucAdbDevices.Execute(ctx.Request().Context())
			if err != nil {
				slog.Error("adb devices usecase error", "err", err)
				return nil // disconnect from server
				// return ctx.JSON(http.StatusInternalServerError, apiconv.MakeInternalServerError(err))
			}
			// report to client
			if err := reporter.ReportDeviceList(updated); err != nil {
				slog.Error("device list report error", "err", err, "updated", updated)
				return nil // disconnect from server
				// return ctx.JSON(http.StatusInternalServerError, apiconv.MakeInternalServerError(err))
			}
		}
	}
}
