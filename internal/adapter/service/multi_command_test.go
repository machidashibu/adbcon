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

func TestMultiCommand(t *testing.T) {
	multi := service.NewMultiCommand().
		Add("S1", infra.NewCommand("ping", "-n", "4", "8.8.8.8")).
		Add("S2", infra.NewCommand("ping", "-n", "4", "8.8.8.8")).
		Add("S3", infra.NewCommand("ping", "-n", "4", "8.8.8.8"))

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
