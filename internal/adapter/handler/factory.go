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
func Factory(repo domain.DeviceStatusRepository, config factoryConfig) *EchoHandler {
	// prepare usecases
	prov := service.NewDeviceListResulver(infra.NewAdbCommand(domain.CommandDevices, "-l"), controller.AdbDeviceParser{})
	ucUpdate := usecase.NewUpdateDeviceInfoUsecase(repo, prov)
	ucPeriodicalUpdate := usecase.NewPeriodicalUpdateDeviceListUsecase(ucUpdate)
	ucExecute := usecase.NewExecuteAdbCommand()

	return NewEchoHandler(config, ucUpdate, ucPeriodicalUpdate, ucExecute)
}
