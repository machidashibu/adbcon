package toapi_test

import (
	"adbcon/api"
	"adbcon/internal/adapter/model"
	"adbcon/internal/adapter/presenter/toapi"
	"adbcon/internal/domain"
	"errors"
	"net/http"
	"testing"

	"github.com/go-openapi/testify/v2/require"
)

var errorTestDummy = errors.New("error test dummy")

func CommandResult(t *testing.T) {
	t.Run("CommandResult", func(t *testing.T) {
		result := toapi.CommandResult(model.NewCommandResult("serial1", []byte("test result")))
		require.Equal(t, "serial1", result.Serial)
		require.Equal(t, "test result", result.Result)
	})
	t.Run("CommandErrorResult", func(t *testing.T) {
		result := toapi.CommandResult(model.NewCommandErrorResult("serial1", errorTestDummy))
		require.Equal(t, "serial1", result.Serial)
		require.Equal(t, "error test dummy", result.Result)
	})
}

func TestDeviceInfo(t *testing.T) {
	info := toapi.DeviceInfo(model.DeviceInfo{
		model.LabelSerial:  "serial1",
		model.LabelStatus:  domain.Online,
		model.LabelProduct: "product1",
		model.LabelModel:   "model1",
		model.LabelDecide:  "device1",
		model.LabelTid:     1,
	})
	require.Equal(t, api.DeviceInfo{
		Serial:  "serial1",
		Status:  api.Online,
		Product: new("product1"),
		Model:   new("model1"),
		Device:  new("device1"),
		Tid:     new(1)}, info)
}

func TestShortDeviceInfo(t *testing.T) {
	info := toapi.ShortDeviceInfo(model.DeviceInfo{
		model.LabelSerial:  "serial1",
		model.LabelStatus:  domain.Online,
		model.LabelProduct: "product1",
		model.LabelModel:   "model1",
		model.LabelDecide:  "device1",
		model.LabelTid:     1,
	})
	require.Equal(t, api.DeviceInfo{Serial: "serial1", Status: api.Online, Model: new("model1")}, info)
}

func TestDeviceList(t *testing.T) {
	list := toapi.DeviceList(domain.DeviceList{
		model.DeviceInfo{
			model.LabelSerial:  "serial1",
			model.LabelStatus:  domain.Online,
			model.LabelProduct: "product1",
			model.LabelModel:   "model1",
			model.LabelDecide:  "device1",
			model.LabelTid:     1,
		},
		model.DeviceInfo{
			model.LabelSerial:  "serial2",
			model.LabelStatus:  domain.Offline,
			model.LabelProduct: "product2",
			model.LabelModel:   "model2",
			model.LabelDecide:  "device2",
			model.LabelTid:     2,
		},
		model.DeviceInfo{
			model.LabelSerial:  "serial3",
			model.LabelStatus:  domain.UnknownDevice,
			model.LabelProduct: "product3",
			model.LabelModel:   "model3",
			model.LabelDecide:  "device3",
			model.LabelTid:     2,
		},
	})
	require.Equal(t, api.DeviceList{
		api.DeviceInfo{Serial: "serial1", Status: api.Online, Model: new("model1")},
		api.DeviceInfo{Serial: "serial2", Status: api.Offline, Model: new("model2")},
		api.DeviceInfo{Serial: "serial3", Status: api.Unknown, Model: new("model3")},
	}, list)
}

func TestMakeProblemDetails(t *testing.T) {
	api := toapi.MakeProblemDetails(http.StatusInternalServerError, errorTestDummy)
	require.Equal(t, 500, *api.Status)
	require.Equal(t, "Internal Server Error", *api.Title)
	require.Equal(t, "error test dummy", *api.Detail)
}

func TestMakeBadRequest(t *testing.T) {
	api := toapi.MakeBadRequest(errorTestDummy)
	require.Equal(t, 400, *api.Status)
	require.Equal(t, "Bad Request", *api.Title)
	require.Equal(t, "error test dummy", *api.Detail)
}

func TestMakeInternalServerError(t *testing.T) {
	api := toapi.MakeInternalServerError(errorTestDummy)
	require.Equal(t, 500, *api.Status)
	require.Equal(t, "Internal Server Error", *api.Title)
	require.Equal(t, "error test dummy", *api.Detail)
}
