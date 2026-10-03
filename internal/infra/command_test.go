package infra_test

import (
	"adbcon/internal/domain"
	"adbcon/internal/infra"
	"context"
	"testing"
	"time"

	"github.com/go-openapi/testify/v2/require"
)

// This code simulates an async command.
var testcode = `
package main
import ("fmt"; "time")
func main() {
	for index := range 4 {
		time.Sleep(500 * time.Millisecond)
		fmt.Println("ASYNC: ", index, " th")
	}
}
`

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
			args:    []string{"run", "-e", testcode},
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
				results := []domain.CommandResult{}
				for runing {
					select {
					case <-ctx.Done():
						t.Fatalf("test is timeout. (limit: %f sec)", limit.Seconds())
					case result, ok := <-ch:
						if !ok {
							runing = false
							break
						}
						results = append(results, result)
					}
				}
				for index, result := range results {
					require.NotEmpty(t, result)
					t.Logf("[%d] %s", index, result.Text())
				}

			}
		})
	}
}
