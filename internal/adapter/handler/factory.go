package handler

import (
	"adbcon/internal/adapter/controller"
	"adbcon/internal/domain"
	"adbcon/internal/infra"
	"adbcon/internal/usecase"
)

// CreateEchoHandler creates handler for echo server.
func Factory(repo domain.DeviceStatusRepository) *EchoHandler {
	// prepare usecases
	ucAdbDevices := usecase.NewExecuteAdbDevicesUsecase(
		infra.NewCommand(string(domain.CommandAdb), string(domain.CommandDevices), "-l"),
		new(controller.AdbDeviceParser),
		repo,
	)

	return NewEchoHandler(ucAdbDevices)
}
