package infra_test

import (
	"adbcon/internal/infra"
	"context"
	"testing"
	"time"

	"github.com/go-openapi/testify/v2/require"
)

func TestCommandRun(t *testing.T) {
	// testcase
	type testcase struct {
		name     string
		command  string
		args     []string
		failcase bool
	}
	testcases := []testcase{
		{
			name:    "success",
			command: "go",
			args:    []string{"version"},
		},
		{
			name:     "unknown command",
			command:  "abb",
			args:     []string{},
			failcase: true,
		},
	}

	ctx := context.TODO()

	// testing
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := infra.NewCommand(tc.command, tc.args...)
			require.NotNil(t, cmd)

			result, err := cmd.Run(ctx)
			if tc.failcase {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, result)
				t.Log(result.String())
			}
		})
	}
}

func TestCommandStart(t *testing.T) {
	// testcase
	type testcase struct {
		name     string
		command  string
		args     []string
		failcase bool
	}
	testcases := []testcase{
		{
			name:    "success",
			command: "go",
			args: []string{"run", "-e", `
package main
import ("fmt"; "time")
func main() {
	for index := range 4 {
		time.Sleep(500 * time.Millisecond)
		fmt.Println("ASYNC: ", index, " th")
	}
}
`},
		},
		{
			name:     "unknown command",
			command:  "abb",
			args:     []string{},
			failcase: true,
		},
		// TODO: cancel by context
	}

	limit := 10 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), limit)
	defer cancel()

	// testing
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := infra.NewCommand(tc.command, tc.args...)
			require.NotNil(t, cmd)

			// testing
			ch, err := cmd.Start(ctx)
			if tc.failcase {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				runing := true
				output := []byte{}
				for runing {
					select {
					case <-ctx.Done():
						t.Fatalf("test is timeout. (limit: %f sec)", limit.Seconds())
					case out, ok := <-ch:
						if !ok {
							runing = false
							break
						}
						output = append(output, out...)
					}
				}
				require.NotEmpty(t, output)
				t.Log(string(output))
			}
		})
	}
}
