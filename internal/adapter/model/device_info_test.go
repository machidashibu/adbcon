package model_test

import (
	"adbcon/internal/adapter/model"
	"adbcon/internal/domain"
	"testing"

	"github.com/go-openapi/testify/v2/require"
)

func TestDeviceInfo(t *testing.T) {
	// testcase
	type testcase struct {
		name                string
		info                model.DeviceInfo
		expectedSerial      string
		expectedStatus      domain.DeviceStatus
		expectedProduct     string
		expectedModel       string
		expectedDevice      string
		expectedTransportId int
	}
	testcases := []testcase{
		{
			name: "empty",
			info: model.DeviceInfo{},
		},
		{
			name: "all set",
			info: model.DeviceInfo{
				model.LabelSerial:      "serial",
				model.LabelStatus:      domain.Online,
				model.LabelProduct:     "product",
				model.LabelModel:       "model",
				model.LabelDecide:      "device",
				model.LabelTransportId: 1,
			},
			expectedSerial:      "serial",
			expectedStatus:      domain.Online,
			expectedProduct:     "product",
			expectedModel:       "model",
			expectedDevice:      "device",
			expectedTransportId: 1,
		},
		{
			name: "alt. transport_id (tid)",
			info: model.DeviceInfo{
				model.LabelTid: 1,
			},
			expectedTransportId: 1,
		},
	}

	// testing
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expectedSerial, tc.info.Serial()) // default is ""
			if tc.expectedStatus == "" {
				require.Equal(t, domain.UnknownDevice, tc.info.Status()) // default
			} else {
				require.Equal(t, tc.expectedStatus, tc.info.Status())
			}
			if tc.expectedProduct == "" {
				require.Equal(t, domain.Unknown, tc.info.Product()) // default
			} else {
				require.Equal(t, tc.expectedProduct, tc.info.Product())
			}
			if tc.expectedModel == "" {
				require.Equal(t, domain.Unknown, tc.info.Model()) // default
			} else {
				require.Equal(t, tc.expectedModel, tc.info.Model())
			}
			if tc.expectedDevice == "" {
				require.Equal(t, domain.Unknown, tc.info.Device()) // default
			} else {
				require.Equal(t, tc.expectedDevice, tc.info.Device())
			}
			require.Equal(t, tc.expectedTransportId, tc.info.TransportId()) // default is 0
		})
	}
}
