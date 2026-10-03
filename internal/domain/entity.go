package domain

import "time"

// CommandName is a supported command name.
type CommandName string

const (
	CommandAdb         CommandName = "adb"
	CommandDevices     CommandName = "devices"
	CommandPush        CommandName = "push"
	CommandPull        CommandName = "pull"
	CommandShell       CommandName = "shell"
	CommandLogcat      CommandName = "logcat"
	CommandReboot      CommandName = "reboot"
	CommandRoot        CommandName = "root"
	CommandUnroot      CommandName = "unroot"
	CommandStartServer CommandName = "start-server"
	CommandKillServer  CommandName = "kill-server"
)

// Interval is a polling interval to execute command.
type Interval time.Duration

// DeviceStatus is a status of device.
type DeviceStatus string

const (
	Online  DeviceStatus = "online"
	Offline DeviceStatus = "offline"
	// UnknownDeviceis a unknown status for device.
	UnknownDevice DeviceStatus = "unknown"
	Unauthorized  DeviceStatus = "unauthorized"
)

// Unknown is a common unknown status.
const Unknown = "unknown"

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

// CommandOutput is an output of command.
// It has converting method that from bytes (received raw) to string, string array (each lines).
type CommandOutput interface {
	// Bytes provides received raw data by byte array.
	Bytes() []byte
	// String provides strings that is casted raw data.
	String() string
	// Lines provides lines that split by return code.
	Lines() []string
	// IsEmpty returns result is empty or not.
	IsEmpty() bool
}

// CommandCh is a channel of command output.
type CommandCh chan CommandResult
