package usecase

import (
	"adbcon/internal/domain"
	"context"
	"log/slog"
)

type execAdbDevicesExecuter interface {
	Run(ctx context.Context) (domain.CommandOutput, error)
}

type execAdbDevicesRepository interface {
	UpdateDevice(devs domain.DeviceList) (domain.DeviceList, error)
}

type execAdbDevicesParser interface {
	Parse(result domain.CommandOutput) (domain.DeviceList, error)
}

// ExecuteAdbDevicesUsecase is a usecase.
//   - executes `adb devices` command.
//   - stores result of command to repository.
type ExecuteAdbDevicesUsecase struct {
	cmd    execAdbDevicesExecuter
	parser execAdbDevicesParser
	repo   execAdbDevicesRepository
}

// NewExecuteAdbDevicesUsecase creates object for usecase.
func NewExecuteAdbDevicesUsecase(cmd execAdbDevicesExecuter, parser execAdbDevicesParser, repo execAdbDevicesRepository) *ExecuteAdbDevicesUsecase {
	return &ExecuteAdbDevicesUsecase{
		cmd:    cmd,
		parser: parser,
		repo:   repo,
	}
}

// Execute is a usecase that executes `adb devices` command, and stores device information to repository.
// It return list of device that modified status.
func (uc ExecuteAdbDevicesUsecase) Execute(ctx context.Context) (domain.DeviceList, error) {
	// execute command
	result, err := uc.cmd.Run(ctx)
	if err != nil {
		slog.Error("usecase execute adb devices run error", "err", err)
		return nil, err
	}
	// parse result
	devs, err := uc.parser.Parse(result)
	if err != nil {
		slog.Error("usecase execute adb devices parse error", "err", err, "result", result)
		return nil, err
	}

	// update repository
	updated, err := uc.repo.UpdateDevice(devs)
	if err != nil {
		slog.Error("usecase execute adb devices update error", "err", err, "devs", devs)
		return nil, err
	}

	// TODO: check onnline

	return updated, nil
}
