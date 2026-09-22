package model_test

import (
	"adbcon/internal/adapter/model"
	"testing"

	"github.com/go-openapi/testify/v2/require"
)

func TestCommandResult(t *testing.T) {
	// testcase
	type testcase struct {
		name      string
		in        []byte
		outEmpty  bool
		outBytes  []byte
		outString string
		outLines  []string
	}
	testcases := []testcase{
		{
			name:      "normal",
			in:        []byte("test1\r\ntest2\r\n\r\n"),
			outBytes:  []byte{0x74, 0x65, 0x73, 0x74, 0x31, 0x0d, 0x0a, 0x74, 0x65, 0x73, 0x74, 0x32, 0x0d, 0x0a, 0x0d, 0x0a},
			outString: "test1\r\ntest2\r\n\r\n",
			outLines:  []string{"test1", "test2", ""},
		},
		{
			name:     "empty",
			in:       []byte(""),
			outEmpty: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			result := model.CommandResult(tc.in)
			require.Equal(t, tc.outEmpty, result.IsEmpty())
			if !tc.outEmpty {
				require.Equal(t, tc.outBytes, result.Bytes())
				require.Equal(t, tc.outString, result.String())
				require.Equal(t, tc.outLines, result.Lines())
			}
		})
	}

}
