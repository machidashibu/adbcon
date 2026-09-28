package presenter_test

import (
	"adbcon/internal/adapter/model"
	"adbcon/internal/adapter/presenter"
	"adbcon/internal/domain"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/go-openapi/testify/v2/require"
	"github.com/labstack/echo/v4"
)

func makeInfo(args ...any) domain.DeviceInfo {
	info := model.DeviceInfo{}
	for index := 0; index < len(args); index += 2 {
		info[args[index].(string)] = args[index+1]
	}
	return info
}

func TestSSEReport(t *testing.T) {
	// testcase
	type testcase struct {
		name string
		sub  func(t *testing.T)
	}
	testcases := []testcase{
		{
			name: "ReportDeviceList",
			sub:  testReportDeviceList,
		},
		{
			name: "ReportClose",
			sub:  testReportClose,
		},
		{
			name: "ReportDevicesListMultiThread",
			sub:  testReportDevicesListMultiThread,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, tc.sub)
	}
}

func testReportDeviceList(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	r := presenter.NewSSEReporter(c.Response(), 500*time.Millisecond, 100)

	devs := domain.DeviceList{
		makeInfo("serial", "s1", "status", domain.Online, "product", "product1", "model", "model1", "device", "device1", "tid", 1),
		makeInfo("serial", "s2", "status", domain.Offline),
	}
	require.NoError(t, r.ReportDeviceList(devs))

	// expected:
	// data: [{...},{...}]
	//
	resp := model.CommandOutput(rec.Body.Bytes())
	lines := resp.Lines()
	require.Len(t, lines, 2)
	require.True(t, strings.HasPrefix(lines[0], "data:"))
	require.Empty(t, lines[1])
	data, _ := strings.CutPrefix(lines[0], "data:")
	require.JSONEq(t, `[{"serial":"s1","status":"online","model":"model1"},{"serial":"s2","status":"offline","model":"unknown"}]`, data)
	require.Equal(t, "text/event-stream", rec.Header().Get(echo.HeaderContentType))
	require.Equal(t, "no-cache", rec.Header().Get(echo.HeaderCacheControl))
	require.Equal(t, "keep-alive", rec.Header().Get(echo.HeaderConnection))
	require.Equal(t, "'no'", rec.Header().Get("X-Accel-Buffering"))
}

func testReportClose(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	r := presenter.NewSSEReporter(c.Response(), 500*time.Millisecond, 100)

	expected := []string{
		"event: close",
		"data: SSE connection closed.",
		"",
	}
	require.NoError(t, r.ReportClose())
	resp := model.CommandOutput(rec.Body.Bytes())
	lines := resp.Lines()
	require.Len(t, lines, 3)
	require.Equal(t, expected, lines)
	require.Equal(t, "text/event-stream", rec.Header().Get(echo.HeaderContentType))
	require.Equal(t, "no-cache", rec.Header().Get(echo.HeaderCacheControl))
	require.Equal(t, "keep-alive", rec.Header().Get(echo.HeaderConnection))
	require.Equal(t, "'no'", rec.Header().Get("X-Accel-Buffering"))
}

func testReportDevicesListMultiThread(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	r := presenter.NewSSEReporter(c.Response(), 500*time.Millisecond, 100)

	devs := domain.DeviceList{
		makeInfo("serial", "s1", "status", domain.Online),
	}

	var wg sync.WaitGroup
	for range 10 {
		wg.Go(func() {
			for range 5 {
				r.ReportDeviceList(devs)
				time.Sleep(1 * time.Second)
			}
		})
	}

	wg.Wait()

	resp := model.CommandOutput(rec.Body.Bytes())
	lines := resp.Lines()
	require.NotEmpty(t, lines)
	require.Empty(t, lines[len(lines)-1]) // check last line is emty

	inBlock := false
	blockCount := 0
	for index, line := range resp.Lines() {
		if line != "" {
			require.Truef(t, strings.HasPrefix(line, "data:"), "line %dth must start with 'data:': %s", index+1, line)
			if !inBlock {
				inBlock = true
				blockCount++
			}
		} else {
			inBlock = false
		}
	}

	require.NotEqual(t, 0, blockCount)
}
