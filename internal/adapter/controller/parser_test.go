package controller_test

import (
	"adbcon/internal/adapter/controller"
	"adbcon/internal/domain"
	"testing"

	"github.com/go-openapi/testify/v2/require"
)

func TestAdbDevicesParser(t *testing.T) {
	// slog.SetLogLoggerLevel(slog.LevelDebug) // for debug

	var testdata domain.CommandResult = []byte(`List of devices attached
ABC123DEF        device product:product1 model:model1 device:device1 transport_id:1
HIJ456KLM        unauthorized transport_id:6
OPQR789STU       offline transport_id:4

`)
	parser := new(controller.AdbDeviceParser)
	devs, err := parser.Parse(testdata)
	require.NoError(t, err)
	require.Len(t, devs, 3)
	require.Equal(t, "ABC123DEF", devs[0].Serial())
	require.Equal(t, domain.Online, devs[0].Status())
	require.Equal(t, "product1", devs[0].Product())
	require.Equal(t, "model1", devs[0].Model())
	require.Equal(t, "device1", devs[0].Device())
	require.Equal(t, 1, devs[0].TransportId())
	require.Equal(t, "HIJ456KLM", devs[1].Serial())
	require.Equal(t, domain.Unauthorized, devs[1].Status())
	require.Equal(t, 6, devs[1].TransportId())
	require.Equal(t, "OPQR789STU", devs[2].Serial())
	require.Equal(t, domain.Offline, devs[2].Status())
	require.Equal(t, 4, devs[2].TransportId())
}
