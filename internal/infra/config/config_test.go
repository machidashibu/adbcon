package config_test

import (
	"adbcon/internal/infra/config"
	"testing"

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
			name:      "read error",
			path:      "config.yaml",
			readError: true,
		},
		{
			name: "read full setting",
			path: "testdata/config.yaml",
			result: map[string]any{
				"bind": "",
				"port": "12345",
			},
		},
		{
			name: "read minimum setting (check default value)",
			path: "testdata/config_min.yaml",
			result: map[string]any{
				"bind": "localhost",
				"port": "8080",
			},
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
				require.Equal(t, tc.result["bind"], cfg.ServerBind())
				require.Equal(t, tc.result["port"], cfg.ServerPort())
			}
		})
	}
}
