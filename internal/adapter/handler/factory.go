package handler

import (
	"adbcon/internal/adapter/controller"
	"adbcon/internal/adapter/service"
	"adbcon/internal/domain"
	"adbcon/internal/infra"
	"adbcon/internal/usecase"
	"time"
)

type factoryConfig interface {
	ServerSseInterval() time.Duration
	ServerSseLimit() int
}

// CreateEchoHandler creates handler for echo server.
func Factory(repo domain.DeviceStatusRepository, config factoryConfig, version string) *EchoHandler {
	// prepare usecases
	cmdDevInfo := service.NewDeviceListResulver(infra.NewAdbCommand(domain.CommandDevices, "-l"), controller.AdbDeviceParser{})
	ucUpdate := usecase.NewUpdateDeviceInfoUsecase(repo, cmdDevInfo)
	ucPeriodicalUpdate := usecase.NewPeriodicalUpdateDeviceListUsecase(ucUpdate)
	ucExecute := usecase.NewExecuteAdbCommand()
	cmdVersion := service.NewVersionProvider(infra.NewAdbCommand(domain.CommandVersion), controller.NewAdbVersionParser(version))
	ucVersion := usecase.NewGetVersionUsecase(cmdVersion)

	return NewEchoHandler(config, ucUpdate, ucPeriodicalUpdate, ucExecute, ucVersion)
}
