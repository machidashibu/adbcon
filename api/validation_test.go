package api_test

import (
	"adbcon/api"
	"testing"

	"github.com/go-openapi/testify/v2/require"
)

func alloc[T any](value T) *T {
	return &value
}

func TestValidate(t *testing.T) {
	type validator interface {
		Validate() error
	}
	// testcase
	type testcase struct {
		name  string
		param validator
		err   bool
	}
	testcases := []testcase{
		{
			name: "PostAsbDevicesParams success",
			param: api.GetDevicesParams{
				Interval: alloc(api.PollingInterval(5)),
			},
		},
		{
			name: "PostAsbDevicesParams fail (negative value)",
			param: api.GetDevicesParams{
				Interval: alloc(api.PollingInterval(-1)),
			},
			err: true,
		},
	}

	// testing
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.err {
				require.Error(t, tc.param.Validate())
			} else {
				require.NoError(t, tc.param.Validate())
			}
		})
	}
}
