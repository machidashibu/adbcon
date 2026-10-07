package handler

import (
	"adbcon/api"
	"adbcon/internal/adapter/controller/todomain"
	"adbcon/internal/adapter/presenter"
	"adbcon/internal/adapter/presenter/toapi"
	"adbcon/internal/adapter/service"
	"adbcon/internal/domain"
	"adbcon/internal/infra"
	"adbcon/internal/usecase"
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type handlerConfig interface {
	ServerSseInterval() time.Duration
	ServerSseLimit() int
}

type EchoHandler struct {
	ucUpdate           *usecase.UpdateDeviceListUsecase
	ucPeriodicalUpdate *usecase.PeriodicalUpdateDeviceListUsecase
	ucExecute          *usecase.ExecuteAdbCommand
	ucVersion          *usecase.GetVersionUsecase
	sseInterval        time.Duration
	sseLimit           int
}

func NewEchoHandler(config handlerConfig,
	ucUpdate *usecase.UpdateDeviceListUsecase,
	ucPeriodicalUpdate *usecase.PeriodicalUpdateDeviceListUsecase,
	ucExecute *usecase.ExecuteAdbCommand,
	ucVersion *usecase.GetVersionUsecase) *EchoHandler {
	return &EchoHandler{
		ucUpdate:           ucUpdate,
		ucPeriodicalUpdate: ucPeriodicalUpdate,
		ucExecute:          ucExecute,
		ucVersion:          ucVersion,
		sseInterval:        config.ServerSseInterval(),
		sseLimit:           config.ServerSseLimit(),
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

	// convert to domain
	interval, priodical := todomain.Interval(params.Interval)

	// prepare reporter
	reporter := presenter.NewSSEReporter(ctx.Response(), h.sseInterval, h.sseLimit)
	defer reporter.ReportClose()

	if priodical {
		// periodical
		if err := h.ucPeriodicalUpdate.Update(ctx.Request().Context(), reporter, interval); err != nil {
			slog.Error("periodical update error", "err", err)
			return nil
		}
	} else {
		// 1 shot
		if err := h.ucUpdate.Update(ctx.Request().Context(), reporter, interval); err != nil {
			slog.Error("periodical update error", "err", err)
			return nil
		}
	}

	return nil
}

// GetVersion Get version of application, ADB and SDK.
// (GET /api/adb/version)
func (h *EchoHandler) GetVersion(ctx echo.Context) error {
	slog.Debug("EchoHandler::GetVersion")

	// prepare reporter
	reporter := presenter.NewSSEReporter(ctx.Response(), h.sseInterval, h.sseLimit)
	defer reporter.ReportClose()

	// get version
	if err := h.ucVersion.Get(ctx.Request().Context(), reporter); err != nil {
		slog.Error("get version error", "err", err)
		return nil
	}

	return nil
}

// ExecuteAdbPush Execute adb push command that push file(s) and report result.
// (POST /api/adb/push)
func (h *EchoHandler) ExecuteAdbPush(ctx echo.Context, params api.ExecuteAdbPushParams) error {
	// bind body
	var body api.ExecuteAdbPullJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return ctx.JSON(http.StatusBadRequest, toapi.MakeBadRequest(err))
	}

	slog.Debug("EchoHandler::ExecuteAdbPush", "params", params, "body", body)

	// convert to domain
	args := todomain.CommandArgs(&params.Args)
	targets, err := todomain.SerialList(body)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, toapi.MakeBadRequest(err))
	}

	// parepare commands
	exec := service.NewMultiCommandProvider()
	for _, serial := range targets {
		exec.Add(serial, infra.NewAdbCommandWithSerial(domain.CommandPush, serial))
	}

	// prepare reporter
	reporter := presenter.NewSSEReporter(ctx.Response(), h.sseInterval, h.sseLimit)
	defer reporter.ReportClose()

	// execute command
	if err := h.ucExecute.Execute(ctx.Request().Context(), reporter, exec, args...); err != nil {
		slog.Error("excute command error", "err", err, "cmd", domain.CommandPush, "args", args)
		return nil
	}

	return nil
}

// ExecuteAdbPull Execute adb pull command that pull file(s) and report result.
// (POST /api/adb/pull)
func (h *EchoHandler) ExecuteAdbPull(ctx echo.Context, params api.ExecuteAdbPullParams) error {
	// bind body
	var body api.ExecuteAdbPullJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return ctx.JSON(http.StatusBadRequest, toapi.MakeBadRequest(err))
	}

	slog.Debug("EchoHandler::ExecuteAdbPull", "params", params, "body", body)

	// convert to domain (with validate)
	args := todomain.CommandArgs(&params.Args)
	targets, err := todomain.SerialList(body)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, toapi.MakeBadRequest(err))
	}

	// parepare commands
	exec := service.NewMultiCommandProvider()
	for _, serial := range targets {
		exec.Add(serial, infra.NewAdbCommandWithSerial(domain.CommandPull, serial))
	}

	// prepare reporter
	reporter := presenter.NewSSEReporter(ctx.Response(), h.sseInterval, h.sseLimit)
	defer reporter.ReportClose()

	// execute command
	if err := h.ucExecute.Execute(ctx.Request().Context(), reporter, exec, args...); err != nil {
		slog.Error("excute command error", "err", err, "cmd", domain.CommandPull, "args", args)
		return nil
	}

	return nil
}

// ExecuteAdbShell Execute adb shell command and report result.
// (POST /api/adb/shell)
func (h *EchoHandler) ExecuteAdbShell(ctx echo.Context, params api.ExecuteAdbShellParams) error {
	// bind body
	var body api.ExecuteAdbShellJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return ctx.JSON(http.StatusBadRequest, toapi.MakeBadRequest(err))
	}

	slog.Debug("EchoHandler::ExecuteAdbShell", "params", params, "body", body)

	// convert to domain (with validate)
	args := todomain.CommandArgs(&params.Args)
	targets, err := todomain.SerialList(body)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, toapi.MakeBadRequest(err))
	}

	// parepare commands
	exec := service.NewMultiCommandProvider()
	for _, serial := range targets {
		exec.Add(serial, infra.NewAdbCommandWithSerial(domain.CommandShell, serial))
	}

	// prepare reporter
	reporter := presenter.NewSSEReporter(ctx.Response(), h.sseInterval, h.sseLimit)
	defer reporter.ReportClose()

	// execute command
	if err := h.ucExecute.Execute(ctx.Request().Context(), reporter, exec, args...); err != nil {
		slog.Error("excute command error", "err", err, "cmd", domain.CommandShell, "args", args)
		return nil
	}

	return nil
}

// ExecuteAdbLogcat Execute adb logcat command and report result.
// (POST /api/adb/logcat)
func (h *EchoHandler) ExecuteAdbLogcat(ctx echo.Context, params api.ExecuteAdbLogcatParams) error {
	// bind body
	var body api.ExecuteAdbLogcatJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return ctx.JSON(http.StatusBadRequest, toapi.MakeBadRequest(err))
	}

	slog.Debug("EchoHandler::ExecuteAdbLogcat", "params", params, "body", body)

	// convert to domain (with validate)
	args := todomain.CommandArgs(params.Args)
	targets, err := todomain.SerialList(body)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, toapi.MakeBadRequest(err))
	}

	// parepare commands
	exec := service.NewMultiCommandProvider()
	for _, serial := range targets {
		exec.Add(serial, infra.NewAdbCommandWithSerial(domain.CommandLogcat, serial))
	}

	// prepare reporter
	reporter := presenter.NewSSEReporter(ctx.Response(), h.sseInterval, h.sseLimit)
	defer reporter.ReportClose()

	// execute command
	if err := h.ucExecute.Execute(ctx.Request().Context(), reporter, exec, args...); err != nil {
		slog.Error("excute command error", "err", err, "cmd", domain.CommandLogcat, "args", args)
		return nil
	}

	return nil
}

// ExecuteAdbReboot Execute adb reroot command and report result.
// (POST /api/adb/reboot)
func (h *EchoHandler) ExecuteAdbReboot(ctx echo.Context, params api.ExecuteAdbRebootParams) error {
	// bind body
	var body api.ExecuteAdbRootJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return ctx.JSON(http.StatusBadRequest, toapi.MakeBadRequest(err))
	}

	slog.Debug("EchoHandler::ExecuteAdbReboot", "body", body)

	// convert to domain (with validate)
	args := todomain.RebootArgs(params.Args)
	targets, err := todomain.SerialList(body)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, toapi.MakeBadRequest(err))
	}

	// parepare commands
	exec := service.NewMultiCommandProvider()
	for _, serial := range targets {
		exec.Add(serial, infra.NewAdbCommandWithSerial(domain.CommandReboot, serial))
	}

	// prepare reporter
	reporter := presenter.NewSSEReporter(ctx.Response(), h.sseInterval, h.sseLimit)
	defer reporter.ReportClose()

	// execute command
	if err := h.ucExecute.Execute(ctx.Request().Context(), reporter, exec, args...); err != nil {
		slog.Error("excute command error", "err", err, "cmd", domain.CommandReboot, "args", args)
		return nil
	}

	return nil
}

// ExecuteAdbRoot Execute adb root command and report result.
// (POST /api/adb/root)
func (h *EchoHandler) ExecuteAdbRoot(ctx echo.Context) error {
	// bind body
	var body api.ExecuteAdbRootJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return ctx.JSON(http.StatusBadRequest, toapi.MakeBadRequest(err))
	}

	slog.Debug("EchoHandler::ExecuteAdbRoot", "body", body)

	// convert to domain (with validate)
	targets, err := todomain.SerialList(body)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, toapi.MakeBadRequest(err))
	}

	// parepare commands
	exec := service.NewMultiCommandProvider()
	for _, serial := range targets {
		exec.Add(serial, infra.NewAdbCommandWithSerial(domain.CommandRoot, serial))
	}

	// prepare reporter
	reporter := presenter.NewSSEReporter(ctx.Response(), h.sseInterval, h.sseLimit)
	defer reporter.ReportClose()

	// execute command
	if err := h.ucExecute.Execute(ctx.Request().Context(), reporter, exec); err != nil {
		slog.Error("excute command error", "err", err, "cmd", domain.CommandRoot)
		return nil
	}

	return nil
}

// ExecuteAdbUnroot Execute adb unroot command and report result.
// (POST /api/adb/unroot)
func (h *EchoHandler) ExecuteAdbUnroot(ctx echo.Context) error {
	// bind body
	var body api.ExecuteAdbUnrootJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return ctx.JSON(http.StatusBadRequest, toapi.MakeBadRequest(err))
	}

	slog.Debug("EchoHandler::ExecuteAdbUnroot", "body", body)

	// convert to domain (with validate)
	targets, err := todomain.SerialList(body)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, toapi.MakeBadRequest(err))
	}

	// parepare commands
	exec := service.NewMultiCommandProvider()
	for _, serial := range targets {
		exec.Add(serial, infra.NewAdbCommandWithSerial(domain.CommandUnroot, serial))
	}

	// prepare reporter
	reporter := presenter.NewSSEReporter(ctx.Response(), h.sseInterval, h.sseLimit)
	defer reporter.ReportClose()

	// execute command
	if err := h.ucExecute.Execute(ctx.Request().Context(), reporter, exec); err != nil {
		slog.Error("excute command error", "err", err, "cmd", domain.CommandUnroot)
		return nil
	}

	return nil
}

// ExecuteAdbStartServer Execute adb start-server command and report result.
// (POST /api/adb/start-server)
func (h *EchoHandler) ExecuteAdbStartServer(ctx echo.Context) error {
	slog.Debug("EchoHandler::ExecuteAdbStartServer")

	// prepare command
	cmd := infra.NewAdbCommand(domain.CommandStartServer)

	// prepare reporter
	reporter := presenter.NewSSEReporter(ctx.Response(), h.sseInterval, h.sseLimit)
	defer reporter.ReportClose()

	// execute command
	if err := h.ucExecute.Execute(ctx.Request().Context(), reporter, cmd); err != nil {
		slog.Error("excute command error", "err", err, "cmd", domain.CommandStartServer)
		return nil
	}

	return nil
}

// ExecuteAdbKillServer Execute adb kill-server command and report result.
// (POST /api/adb/kill-server)
func (h *EchoHandler) ExecuteAdbKillServer(ctx echo.Context) error {
	slog.Debug("EchoHandler::ExecuteAdbKillServer")

	// prepare command
	cmd := infra.NewAdbCommand(domain.CommandKillServer)

	// prepare reporter
	reporter := presenter.NewSSEReporter(ctx.Response(), h.sseInterval, h.sseLimit)
	defer reporter.ReportClose()

	// execute command
	if err := h.ucExecute.Execute(ctx.Request().Context(), reporter, cmd); err != nil {
		slog.Error("excute command error", "err", err, "cmd", domain.CommandKillServer)
		return nil
	}

	return nil
}
