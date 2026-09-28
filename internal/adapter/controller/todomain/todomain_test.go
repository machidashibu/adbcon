package todomain_test

import (
	"adbcon/api"
	"adbcon/internal/adapter/controller/todomain"
	"adbcon/internal/domain"
	"testing"
	"time"

	"github.com/go-openapi/testify/v2/require"
)

func TestCommandArgs(t *testing.T) {
	type testcase struct {
		name string
		in   *api.CommandArgs
		out  []string
	}
	testcases := []testcase{
		{
			name: "normal",
			in:   &api.CommandArgs{"arg1", "arg2", "arg3"},
			out:  []string{"arg1", "arg2", "arg3"},
		},
		{
			name: "empty",
			in:   &api.CommandArgs{},
			out:  []string{},
		},
		{
			name: "nil",
			in:   nil,
			out:  []string{},
		},
	}

	// testing
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			list := todomain.CommandArgs(tc.in)
			require.Equal(t, tc.out, list)
		})
	}
}

func TestRebootArgs(t *testing.T) {
	type testcase struct {
		name string
		in   *api.RebootArgs
		out  []string
	}
	testcases := []testcase{
		{
			name: "bootloader",
			in:   new(api.Bootloader),
			out:  []string{"bootloader"},
		},
		{
			name: "recovery",
			in:   new(api.Recovery),
			out:  []string{"recovery"},
		},
		{
			name: "sideload",
			in:   new(api.Sideload),
			out:  []string{"sideload"},
		},
		{
			name: "sideload-auto-reboot",
			in:   new(api.SideloadAutoReboot),
			out:  []string{"sideload-auto-reboot"},
		},
		{
			name: "empty",
			in:   new(api.RebootArgs("")),
			out:  []string{""},
		},
		{
			name: "nil",
			in:   nil,
			out:  []string{},
		},
	}

	// testing
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			list := todomain.RebootArgs(tc.in)
			require.Equal(t, tc.out, list)
		})
	}
}

func TestDeviceStatus(t *testing.T) {
	type testcase struct {
		name string
		in   string
		out  domain.DeviceStatus
	}
	testcases := []testcase{
		{
			name: "offline",
			in:   string(api.Offline),
			out:  domain.Offline,
		},
		{
			name: "online",
			in:   string(api.Online),
			out:  domain.Online,
		},
		{
			name: "device",
			in:   "device",
			out:  domain.Online,
		},
		{
			name: "unauthorized",
			in:   string(api.Unauthorized),
			out:  domain.Unauthorized,
		},
		{
			name: "unknown",
			in:   string(api.Unknown),
			out:  domain.UnknownDevice,
		},
		{
			name: "empty",
			in:   "",
			out:  domain.UnknownDevice,
		},
	}

	// testing
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			status := todomain.DeviceStatus(tc.in)
			require.Equal(t, tc.out, status)
		})
	}
}

func TestInterval(t *testing.T) {
	type testcase struct {
		name        string
		in          *api.PollingInterval
		outInterval domain.Interval
		outRepeat   bool
	}
	testcases := []testcase{
		{
			name:        "positive",
			in:          new(api.PollingInterval(10)),
			outInterval: domain.Interval(10 * time.Second),
			outRepeat:   true,
		},
		{
			name:        "negative",
			in:          new(api.PollingInterval(-10)),
			outInterval: domain.Interval(0),
			outRepeat:   false,
		},
		{
			name:        "zero",
			in:          new(api.PollingInterval(0)),
			outInterval: domain.Interval(0),
			outRepeat:   false,
		},
	}

	// testing
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			interval, repeat := todomain.Interval(tc.in)
			require.Equal(t, tc.outInterval, interval)
			require.Equal(t, tc.outRepeat, repeat)
		})
	}
}

func TestSerialList(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		list, err := todomain.SerialList(api.SerialList{
			api.DeviceSerial("serial1"),
			api.DeviceSerial("serial2"),
			api.DeviceSerial("serial3"),
		})
		require.NoError(t, err)
		require.Equal(t, domain.SerialList{"serial1", "serial2", "serial3"}, list)
	})
	t.Run("empty", func(t *testing.T) {
		_, err := todomain.SerialList(api.SerialList{})
		require.Error(t, err)
	})
}
