package handler_test

import (
	"adbcon/api"
	"adbcon/internal/adapter/handler"
	"adbcon/internal/adapter/model"
	"adbcon/internal/domain"
	"adbcon/internal/infra/database"
	"adbcon/internal/usecase"
	"bufio"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-openapi/testify/v2/require"
	"github.com/labstack/echo/v4"
)

func TestEchoHandler(t *testing.T) {
	// testcase
	type testcase struct {
		name string
		sub  func(t *testing.T)
	}
	testcases := []testcase{
		{
			name: "GetMainPage",
			sub:  testGetMainPage,
		},
		{
			name: "GetDevices",
			sub:  testGetDevices,
		},
	}

	// testing
	for _, tc := range testcases {
		t.Run(tc.name, tc.sub)
	}
}

// --- GetMainPage ---
func testGetMainPage(t *testing.T) {
	h := handler.Factory(database.StubDatabase{})

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/gui", strings.NewReader(""))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// testing
	require.ErrorIs(t, h.GetMainPage(c), echo.ErrNotFound)
}

// --- GetDevices ---
type stubGetDevicesCommand struct{}

func (s stubGetDevicesCommand) Run(ctx context.Context) (domain.CommandResult, error) {
	return []byte{}, nil
}

type stubGetDevicesPaerser struct {
	list  []domain.DeviceList
	count int
}

func (s *stubGetDevicesPaerser) Parse(result domain.CommandResult) (domain.DeviceList, error) {
	if s.count >= len(s.list) {
		return s.list[s.count-1], nil
	}
	parsed := s.list[s.count]
	s.count++
	return parsed, nil
}

func makeInfo(args ...any) domain.DeviceInfo {
	info := model.DeviceInfo{}
	for index := 0; index < len(args); index += 2 {
		info[args[index].(string)] = args[index+1]
	}
	return info
}

func testGetDevices(t *testing.T) {
	// TODO: implements failed case
	expecteded := []string{
		`[{"serial":"S1","status":"unknown","model":"model1"}]`,
		`[{"serial":"S1","status":"offline","model":"model1"},{"serial":"S2","status":"offline","model":"model2"}]`,
		`[{"serial":"S1","status":"online","model":"model1"},{"serial":"S2","status":"unauthorized","model":"model2"}]`,
	}
	timeout := time.Duration(len(expecteded)+1) * time.Second

	ucAdbDevices := usecase.NewExecuteAdbDevicesUsecase(
		stubGetDevicesCommand{},
		&stubGetDevicesPaerser{
			list: []domain.DeviceList{
				{
					makeInfo("serial", "S1", "status", domain.Unknown, "model", "model1"),
				},
				{
					makeInfo("serial", "S1", "status", domain.Offline, "model", "model1"),
					makeInfo("serial", "S2", "status", domain.Offline, "model", "model2"),
				},
				{
					makeInfo("serial", "S1", "status", domain.Online, "model", "model1"),
					makeInfo("serial", "S2", "status", domain.Unauthorized, "model", "model2"),
				},
			},
		},
		database.StubDatabase{})
	h := handler.NewEchoHandler(ucAdbDevices)

	e := echo.New()
	e.GET("/api/devices", func(c echo.Context) error {
		interval := api.PollingInterval(1)
		return h.GetDevices(c, api.GetDevicesParams{
			Interval: &interval,
		})
	})

	server := httptest.NewServer(e)
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, server.URL+"/api/devices", nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// receive events
	type scanResult struct {
		line string
		err  error
	}
	results := make(chan scanResult, 1)
	go func() {
		scan := bufio.NewScanner(resp.Body)
		for scan.Scan() {
			results <- scanResult{line: scan.Text()}
		}
		results <- scanResult{err: scan.Err()}
	}()

	received := []string{}
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	for len(received) < len(expecteded) {
		select {
		case <-timer.C:
			t.Fatalf("SSE event timeout: received %d/%d", len(received), len(expecteded))

		case result := <-results:
			if result.err != nil {
				require.NoError(t, result.err)
			}

			if result.line == "" {
				continue
			}

			t.Log("received:", result.line)
			require.True(t, strings.HasPrefix(result.line, "data:"))
			received = append(received, strings.TrimPrefix(result.line, "data:"))
		}

	}
	cancel()

	// validate test
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, "text/event-stream", resp.Header.Get(echo.HeaderContentType))
	require.Equal(t, "no-cache", resp.Header.Get(echo.HeaderCacheControl))
	require.Equal(t, "no-cache", resp.Header.Get(echo.HeaderCacheControl))
	require.Equal(t, "'no'", resp.Header.Get("X-Accel-Buffering"))
	for index := range len(expecteded) {
		require.JSONEq(t, expecteded[index], received[index])
	}
}
