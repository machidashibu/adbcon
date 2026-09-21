package apiconv

import (
	"adbcon/api"
	"adbcon/internal/domain"
)

func DeviceInfoToApi(info domain.DeviceInfo) api.DeviceInfo {
	product := info.Product()
	model := info.Model()
	device := info.Device()
	tid := info.TransportId()
	return api.DeviceInfo{
		Serial:  info.Serial(),
		Status:  api.DeviceStatus(info.Status()),
		Product: &product,
		Model:   &model,
		Device:  &device,
		Tid:     &tid,
	}
}

func DeviceInfoShortToApi(info domain.DeviceInfo) api.DeviceInfo {
	model := info.Model()
	return api.DeviceInfo{
		Serial: info.Serial(),
		Status: api.DeviceStatus(info.Status()),
		Model:  &model,
	}
}
