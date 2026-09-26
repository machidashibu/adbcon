package controller

import (
	"adbcon/api"
	"adbcon/internal/domain"
	"time"
)

// IntervalToDomain casts api value to domain entity.
// It returns true if interval is valid value.
// It returns false if interval is invalid value. (nil or <= 0)
func IntervalToDomain(interval *api.PollingInterval) (domain.Interval, bool) {
	if interval == nil || *interval <= 0 {
		return domain.Interval(api.PollingInterval(0)), false
	}
	return domain.Interval(*interval * api.PollingInterval(time.Second)), true
}

// DeviceStatusToDomain converts status of command result to domain entity.
// Normally, casts status by the string as is a domain entity.
// But, only `device` replaces to `online`.
// `device` is indicated online device by ADB.
func DeviceStatusToDomain(status string) domain.DeviceStatus {
	if status == "device" {
		return domain.Online // adb devices rule
	}
	return domain.DeviceStatus(status)
}
