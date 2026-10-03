package usecase

import (
	"adbcon/internal/domain"
	"context"
	"log/slog"
)

type updateDeviceListUsecaseRepository interface {
	UpdateDevice(devs domain.DeviceList) (domain.DeviceList, error)
	CheckAlive(alive domain.Interval) (domain.DeviceList, error)
}

type updateDeviceListReporter interface {
	ReportDeviceList(devs domain.DeviceList) error
}

// ExecuteAdbDevicesUsecase is a usecase.
//   - executes `adb devices` command.
//   - stores result of command to repository.
type UpdateDeviceListUsecase struct {
	repo updateDeviceListUsecaseRepository
	cmd  domain.DeviceListResolver
}

// NewExecuteAdbDevicesUsecase creates object for usecase.
func NewUpdateDeviceInfoUsecase(repo updateDeviceListUsecaseRepository, cmd domain.DeviceListResolver) *UpdateDeviceListUsecase {
	return &UpdateDeviceListUsecase{
		repo: repo,
		cmd:  cmd,
	}
}

// Update stores device information to repository.
// It return list of device that modified status.
func (uc UpdateDeviceListUsecase) Update(ctx context.Context, reporter updateDeviceListReporter, alive domain.Interval) error {
	slog.Debug("UpdateDeviceListUsecase::Update", "alive", alive)

	// get devce list
	devs, err := uc.cmd.GetDeviceList(ctx)
	if err != nil {
		slog.Error("usecase update device list get error", "err", err)
		return err
	}

	// update repository
	updated, err := uc.repo.UpdateDevice(devs)
	if err != nil {
		slog.Error("usecase update device list update error", "err", err, "devs", devs)
		return err
	}

	// check alive
	dead, err := uc.repo.CheckAlive(alive)
	if err != nil {
		slog.Error("usecase update device list check alive error", "err", err, "alive", alive)
		return err
	}

	// TODO: marge updated and dead
	// report
	if err := reporter.ReportDeviceList(updated); err != nil {
		slog.Error("usecase update device list report error", "updated", updated, "dead", dead)
		return nil
	}

	return nil
}
