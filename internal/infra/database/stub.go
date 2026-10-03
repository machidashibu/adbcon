package database

import (
	"adbcon/internal/domain"
)

type StubDatabase struct{}

func (db StubDatabase) UpdateDevice(devs domain.DeviceList) (domain.DeviceList, error) {
	return devs, nil
}

func (db StubDatabase) CheckAlive(alive domain.Interval) (domain.DeviceList, error) {
	return domain.DeviceList{}, nil
}
