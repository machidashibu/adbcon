package todomain

import (
	"adbcon/api"
	"adbcon/internal/domain"
	"fmt"
	"time"
)

// CommandArgs maps api value to domain entity.
func CommandArgs(args *api.CommandArgs) []string {
	if args == nil {
		return []string{}
	}
	return []string(*args)
}

// RebootArgs maps api value to domain entity.
func RebootArgs(args *api.RebootArgs) []string {
	if args == nil {
		return []string{}
	}
	return []string{string(*args)}
}

// DeviceStatusToDomain converts status of command result to domain entity.
// Normally, casts status by the string as is a domain entity.
// But, only `device` replaces to `online`.
// `device` is indicated online device by ADB.
func DeviceStatus(status string) domain.DeviceStatus {
	switch status {
	case string(api.Offline):
		return domain.Offline
	case string(api.Online), "device": // adb devices rule
		return domain.Online
	case string(api.Unauthorized):
		return domain.Unauthorized
	default: // api.Unknown
		return domain.UnknownDevice
	}
}

// Interval maps api value to domain entity.
// It returns true if interval is valid value.
// It returns false if interval is invalid value. (nil or <= 0)
func Interval(interval *api.PollingInterval) (domain.Interval, bool) {
	if interval == nil || *interval <= 0 {
		return domain.Interval(0), false
	}
	return domain.Interval(*interval * api.PollingInterval(time.Second)), true
}

// SerialList maps api value to domain entity.
// It errors if list is empty.
func SerialList(list api.SerialList) (domain.SerialList, error) {
	if len(list) == 0 {
		return nil, fmt.Errorf("not allowed empty serial list")
	}
	serials := make(domain.SerialList, len(list))
	for index, serial := range list {
		serials[index] = string(serial)
	}
	return serials, nil
}
