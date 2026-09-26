package controller

import (
	"adbcon/api"
	"adbcon/internal/domain"
	"fmt"
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

func CommandArgsToComain(args *api.CommandArgs) []string {
	if args == nil {
		return []string{}
	}
	return []string(*args)
}

func SerialListToDomain(list api.SerialList) (domain.SerialList, error) {
	if len(list) == 0 {
		return nil, fmt.Errorf("not allowed empty serial list")
	}
	serials := make(domain.SerialList, len(list))
	for index, serial := range list {
		serials[index] = string(serial)
	}
	return serials, nil
}
