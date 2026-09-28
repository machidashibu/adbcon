package toapi

import (
	"adbcon/api"
	"adbcon/internal/domain"
	"net/http"
)

// CommandResult maps domain entity to api value.
func CommandResult(result domain.CommandResult) api.CommandResult {
	return api.CommandResult{
		Serial: result.Serial(),
		Result: result.Text(),
	}
}

// DeviceInfo maps domain entity to api value.
func DeviceInfo(info domain.DeviceInfo) api.DeviceInfo {
	product := info.Product()
	model := info.Model()
	device := info.Device()
	tid := info.TransportId()
	return api.DeviceInfo{
		Serial:  info.Serial(),
		Status:  api.DeviceStatus(info.Status()),
		Product: &product,
		Model:   &model,
		Device:  &device,
		Tid:     &tid,
	}
}

// ShortDeviceInfo makes a part of api value from domain entity.
// It target are `serial`, `status` and `model`.
func ShortDeviceInfo(info domain.DeviceInfo) api.DeviceInfo {
	model := info.Model()
	return api.DeviceInfo{
		Serial: info.Serial(),
		Status: api.DeviceStatus(info.Status()),
		Model:  &model,
	}
}

// DeviceList maps domain entity to api value.
// DeviceInfo applies short varsion. (See: ShortDeviceInfo)
func DeviceList(list domain.DeviceList) api.DeviceList {
	apiList := make(api.DeviceList, len(list))
	for index, info := range list {
		apiList[index] = ShortDeviceInfo(info)
	}
	return apiList
}

// MakeProblemDetails makes ProblemDetails pbject.
func MakeProblemDetails(code int, err error) api.ProblemDetails {
	status := code
	title := http.StatusText(status)
	detail := err.Error()

	return api.ProblemDetails{
		Status: &status,
		Title:  &title,
		Detail: &detail,
	}
}

// MakeBadRequest makes ProblemDetails object with status code = 400.
func MakeBadRequest(err error) api.ProblemDetails {
	return MakeProblemDetails(http.StatusBadRequest, err)
}

// MakeInternalServerError makes ProblemDetails object with status code = 500.
func MakeInternalServerError(err error) api.ProblemDetails {
	return MakeProblemDetails(http.StatusInternalServerError, err)
}
