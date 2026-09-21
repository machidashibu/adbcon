package apiconv

import (
	"adbcon/api"
	"adbcon/internal/domain"
)

func DeviceListToApi(list domain.DeviceList) api.DeviceList {
	apiList := make(api.DeviceList, len(list))
	for index, info := range list {
		apiList[index] = DeviceInfoShortToApi(info)
	}
	return apiList
}
