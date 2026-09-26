package handler

import (
	"adbcon/api"
	"adbcon/internal/adapter/controller"
	"adbcon/internal/adapter/presenter"
	"adbcon/internal/adapter/presenter/apiconv"
	"adbcon/internal/adapter/service"
	"adbcon/internal/domain"
	"adbcon/internal/infra"
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
	slog.Debug("EchoHandler::GetMainPage")
	// NOTE: This endpoint is not called, because proceeds by statics middleware.
	return echo.ErrNotFound
}

// GetDevices Execute adb devices command and report device status.
// (GET /api/devices)
func (h *EchoHandler) GetDevices(ctx echo.Context, params api.GetDevicesParams) error {
	slog.Debug("EchoHandler::GetDevices", "params", params)

	// to domain
	interval, ok := controller.IntervalToDomain(params.Interval)

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

	// proceeds polling if interval is specified in query
	if ok {
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

	return nil
}

// ExecuteAdbShell Execute adb shell command and report result.
// (POST /api/adb/shell)
func (h *EchoHandler) ExecuteAdbShell(ctx echo.Context, params api.ExecuteAdbShellParams) error {
	// bind body
	var body api.ExecuteAdbShellJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return ctx.JSON(http.StatusBadRequest, apiconv.MakeBadRequest(err))
	}

	slog.Debug("EchoHandler::ExecuteAdbShell", "params", params, "body", body)

	// convert to domain (with validate)
	args := controller.CommandArgsToDomain(&params.Args)
	targets, err := controller.SerialListToDomain(body)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, apiconv.MakeBadRequest(err))
	}

	// parepare commands
	exec := service.NewMultiCommand()
	for _, serial := range targets {
		cmd := infra.NewAdbCommandWithSerial(domain.CommandShell, serial, args...)
		exec.Add(serial, cmd)
	}

	// prepare reporter
	reporter := presenter.NewSSEReporter(ctx.Response())

	// start command
	ch := exec.Start(ctx.Request().Context())
	for {
		result, ok := <-ch
		if !ok {
			reporter.ReportClose()
			break
		}
		if err := reporter.ReportCommandResult(result); err != nil {
			slog.Error("command result report error", "err", err, "result", result)
			return nil
		}
	}

	return nil
}

// ExecuteAdbReboot Execute adb reroot command and report result.
// (POST /api/adb/reboot)
func (h *EchoHandler) ExecuteAdbReboot(ctx echo.Context, params api.ExecuteAdbRebootParams) error {
	// bind body
	var body api.ExecuteAdbRootJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return ctx.JSON(http.StatusBadRequest, apiconv.MakeBadRequest(err))
	}

	slog.Debug("EchoHandler::ExecuteAdbReboot", "body", body)

	// convert to domain (with validate)
	args := controller.RebootArgsToDomain(params.Args)
	targets, err := controller.SerialListToDomain(body)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, apiconv.MakeBadRequest(err))
	}

	// parepare commands
	exec := service.NewMultiCommand()
	for _, serial := range targets {
		cmd := infra.NewAdbCommandWithSerial(domain.CommandReboot, serial, args...)
		exec.Add(serial, cmd)
	}

	// prepare reporter
	reporter := presenter.NewSSEReporter(ctx.Response())

	// start command
	ch := exec.Start(ctx.Request().Context())
	for {
		result, ok := <-ch
		if !ok {
			reporter.ReportClose()
			break
		}
		if err := reporter.ReportCommandResult(result); err != nil {
			slog.Error("command result report error", "err", err, "result", result)
			return nil
		}
	}

	return nil
}

// ExecuteAdbRoot Execute adb root command and report result.
// (POST /api/adb/root)
func (h *EchoHandler) ExecuteAdbRoot(ctx echo.Context) error {
	// bind body
	var body api.ExecuteAdbRootJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return ctx.JSON(http.StatusBadRequest, apiconv.MakeBadRequest(err))
	}

	slog.Debug("EchoHandler::ExecuteAdbRoot", "body", body)

	// convert to domain (with validate)
	targets, err := controller.SerialListToDomain(body)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, apiconv.MakeBadRequest(err))
	}

	// parepare commands
	exec := service.NewMultiCommand()
	for _, serial := range targets {
		cmd := infra.NewAdbCommandWithSerial(domain.CommandRoot, serial)
		exec.Add(serial, cmd)
	}

	// prepare reporter
	reporter := presenter.NewSSEReporter(ctx.Response())

	// start command
	ch := exec.Start(ctx.Request().Context())
	for {
		result, ok := <-ch
		if !ok {
			reporter.ReportClose()
			break
		}
		if err := reporter.ReportCommandResult(result); err != nil {
			slog.Error("command result report error", "err", err, "result", result)
			return nil
		}
	}

	return nil
}

// ExecuteAdbUnroot Execute adb unroot command and report result.
// (POST /api/adb/unroot)
func (h *EchoHandler) ExecuteAdbUnroot(ctx echo.Context) error {
	// bind body
	var body api.ExecuteAdbUnrootJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return ctx.JSON(http.StatusBadRequest, apiconv.MakeBadRequest(err))
	}

	slog.Debug("EchoHandler::ExecuteAdbUnroot", "body", body)

	// convert to domain (with validate)
	targets, err := controller.SerialListToDomain(body)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, apiconv.MakeBadRequest(err))
	}

	// parepare commands
	exec := service.NewMultiCommand()
	for _, serial := range targets {
		cmd := infra.NewAdbCommandWithSerial(domain.CommandUnroot, serial)
		exec.Add(serial, cmd)
	}

	// prepare reporter
	reporter := presenter.NewSSEReporter(ctx.Response())

	// start command
	ch := exec.Start(ctx.Request().Context())
	for {
		result, ok := <-ch
		if !ok {
			reporter.ReportClose()
			break
		}
		if err := reporter.ReportCommandResult(result); err != nil {
			slog.Error("command result report error", "err", err, "result", result)
			return nil
		}
	}

	return nil
}
