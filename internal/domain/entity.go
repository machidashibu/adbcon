package domain

import "time"

// Supported command name
type CommandName string

const (
	CommandAdb     CommandName = "adb"
	CommandDevices CommandName = "devices"
)

// Polling interval
type Interval time.Duration

// device status
type DeviceStatus string

const (
	Unknown = "unknown"

	Online        DeviceStatus = "online"
	Offline       DeviceStatus = "offline"
	UnknownDevice DeviceStatus = "unknown"
	Unauthorized  DeviceStatus = "unauthorized"
)

// device einformation
type DeviceInfo interface {
	Serial() string
	Status() DeviceStatus
	Product() string
	Model() string
	Device() string
	TransportId() int
}

// list of device information
type DeviceList []DeviceInfo
