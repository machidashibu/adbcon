package database

import "adbcon/internal/domain"

type StubDatabase struct{}

func (db StubDatabase) UpdateDevice(devs domain.DeviceList) (domain.DeviceList, error) {
	return devs, nil
}
