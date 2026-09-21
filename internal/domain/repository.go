package domain

type DeviceStatusRepository interface {
	UpdateDevice(devs DeviceList) (DeviceList, error)
}
