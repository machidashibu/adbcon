package presenter_test

import (
	"adbcon/internal/adapter/model"
	"adbcon/internal/adapter/presenter"
	"adbcon/internal/domain"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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

func TestReportDeviceList(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	r := presenter.NewSSEReporter(c.Response())

	devs := domain.DeviceList{
		makeInfo("serial", "s1", "status", domain.Online, "product", "product1", "model", "model1", "device", "device1", "tid", 1),
		makeInfo("serial", "s2", "status", domain.Offline),
	}
	require.NoError(t, r.ReportDeviceList(devs))

	// expected:
	// data: [{...},{...}]
	//
	resp := domain.CommandResult(rec.Body.Bytes())
	lines := resp.Lines()
	require.Len(t, lines, 2)
	require.True(t, strings.HasPrefix(lines[0], "data:"))
	require.Empty(t, lines[1])
	data, _ := strings.CutPrefix(lines[0], "data:")
	require.JSONEq(t, `[{"serial":"s1","status":"online","model":"model1"},{"serial":"s2","status":"offline","model":"unknown"}]`, data)
}

func TestReportClose(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	r := presenter.NewSSEReporter(c.Response())

	expected := []string{
		"event: close",
		"data: SSE connection closed.",
		"",
	}
	require.NoError(t, r.ReportClose())
	resp := domain.CommandResult(rec.Body.Bytes())
	lines := resp.Lines()
	require.Len(t, lines, 3)
	require.Equal(t, expected, lines)
}
