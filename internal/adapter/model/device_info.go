package model

import "adbcon/internal/domain"

const (
	LabelSerial      = "serial"
	LabelStatus      = "status"
	LabelProduct     = "product"
	LabelModel       = "model"
	LabelDecide      = "device"
	LabelTransportId = "transport_id"
	LabelTid         = "tid"
)

type DeviceInfo map[string]any

func pick[T any](info DeviceInfo, defaultValue T, labels ...string) T {
	for _, label := range labels {
		val, ok := info[label]
		if !ok {
			continue // not existing
		}
		v, ok := val.(T)
		if !ok {
			continue // unmatched type
		}
		return v
	}
	return defaultValue
}

func (d DeviceInfo) Serial() string {
	return pick(d, "", LabelSerial)
}

func (d DeviceInfo) Status() domain.DeviceStatus {
	return pick(d, domain.UnknownDevice, LabelStatus)
}

func (d DeviceInfo) Product() string {
	return pick(d, domain.Unknown, LabelProduct)
}

func (d DeviceInfo) Model() string {
	return pick(d, domain.Unknown, LabelModel)
}

func (d DeviceInfo) Device() string {
	return pick(d, domain.Unknown, LabelDecide)
}

func (d DeviceInfo) TransportId() int {
	return pick(d, 0, LabelDecide, LabelTransportId, LabelTid)
}
