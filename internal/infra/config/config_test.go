package config_test

import (
	"adbcon/internal/infra/config"
	"testing"
	"time"

	"github.com/go-openapi/testify/v2/require"
)

func TestConfig(t *testing.T) {
	// test case
	type testcase struct {
		name      string
		path      string
		readError bool
		result    map[string]any
	}
	testcases := []testcase{
		{
			name: "read error",
			path: "config.yaml",
			// readError: true,	// Change spec. : It did not error if file is not existing
			result: map[string]any{
				"filepath":     "adbcon.log",
				"bind":         "localhost",
				"port":         "8080",
				"sse_interval": 100 * time.Millisecond,
				"sse_limit":    30,
			},
		},
		{
			name: "read full setting",
			path: "testdata/config.yaml",
			result: map[string]any{
				"filepath":     "logfile.log",
				"bind":         "",
				"port":         "12345",
				"sse_interval": 1000 * time.Millisecond,
				"sse_limit":    300,
			},
		},
		{
			name: "read minimum setting (check default value)",
			path: "testdata/config_min.yaml",
			result: map[string]any{
				"filepath":     "adbcon.log",
				"bind":         "localhost",
				"port":         "8080",
				"sse_interval": 100 * time.Millisecond,
				"sse_limit":    30,
			},
		},
		{
			name:      "read invalid setting",
			path:      "testdata/invalid.txt",
			readError: true,
		},
	}

	// testing
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := new(config.Config)
			if tc.readError {
				// testing: file not exist
				require.Error(t, cfg.Read(tc.path))
			} else {
				// testing: read and value
				require.NoError(t, cfg.Read(tc.path))
				require.Equal(t, tc.result["filepath"], cfg.LogFilePath())
				require.Equal(t, tc.result["bind"], cfg.ServerBind())
				require.Equal(t, tc.result["port"], cfg.ServerPort())
				require.Equal(t, tc.result["sse_interval"], cfg.ServerSseInterval())
				require.Equal(t, tc.result["sse_limit"], cfg.ServerSseLimit())
			}
		})
	}
}
