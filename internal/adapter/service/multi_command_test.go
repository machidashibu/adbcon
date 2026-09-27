package service_test

import (
	"adbcon/internal/adapter/service"
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

func TestMultiCommand(t *testing.T) {
	multi := service.NewMultiCommand().
		Add("S1", infra.NewCommand("go", "run", "-e", testcode)).
		Add("S2", infra.NewCommand("go", "run", "-e", testcode)).
		Add("S3", infra.NewCommand("go", "run", "-e", testcode))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// testing
	results := []domain.CommandResult{}
	ch := multi.Start(ctx)
	testing := true
	for testing {
		select {
		case <-ctx.Done():
			t.Fatal("test is timeover")
		case result, ok := <-ch:
			if !ok {
				testing = false
				continue // terminate texting
			}
			require.Contains(t, []string{"S1", "S2", "S3"}, result.Serial())
			require.NoError(t, result.Error())
			t.Log(result.Text())
			results = append(results, result)
		}
	}
	// check not received result
	require.NotEmpty(t, results)
}
