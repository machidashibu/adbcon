package domain

import "time"

// CommandName is a supported command name.
type CommandName string

const (
	CommandAdb     CommandName = "adb"
	CommandDevices CommandName = "devices"
	CommandShell   CommandName = "shell"
	CommandRoot    CommandName = "root"
	CommandUnroot  CommandName = "unroot"
)

// Interval is a polling interval to execute command.
type Interval time.Duration

// DeviceStatus is a status of device.
type DeviceStatus string

const (
	// Unknown is a common unknown status.
	Unknown = "unknown"

	Online  DeviceStatus = "online"
	Offline DeviceStatus = "offline"
	// UnknownDeviceis a unknown status for device.
	UnknownDevice DeviceStatus = "unknown"
	Unauthorized  DeviceStatus = "unauthorized"
)

// DeviceInfo is a device information.
type DeviceInfo interface {
	Serial() string
	Status() DeviceStatus
	Product() string
	Model() string
	Device() string
	TransportId() int
}

// DeviceList is a list of device information.
type DeviceList []DeviceInfo

// SerialList is a list of device serial.
type SerialList []string

// CommandResult is a result of command with device serial.
// It indicates an error by `Error` and `Text` is an error message if error is occurred.
type CommandResult interface {
	Serial() string
	Text() string
	Error() error
}
