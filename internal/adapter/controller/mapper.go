package controller

import (
	"adbcon/api"
	"adbcon/internal/domain"
	"time"
)

func IntervalToDomain(interval api.PollingInterval) domain.Interval {
	return domain.Interval(interval * api.PollingInterval(time.Second))
}

func DeviceStatusToDomain(status string) domain.DeviceStatus {
	if status == "device" {
		return domain.Online // adb devices rule
	}
	return domain.DeviceStatus(status)
}
