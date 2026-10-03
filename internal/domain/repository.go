package domain

type DeviceStatusRepository interface {
	// UpdateDevice updates device information to repository.
	// It returns list of device that modified status.
	UpdateDevice(devs DeviceList) (DeviceList, error)
	// CheckAlive checks alive device in repository.
	// It returns list of device that is not alive. (i.e. modifed offline)
	CheckAlive(alive Interval) (DeviceList, error)
}
