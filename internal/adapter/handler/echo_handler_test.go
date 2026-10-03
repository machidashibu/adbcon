//go:build use_real_adb

package handler_test

import (
	"adbcon/internal/adapter/handler"
	"adbcon/internal/adapter/model"
	"adbcon/internal/domain"
	"adbcon/internal/infra/config"
	"adbcon/internal/infra/database"
	"bufio"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-openapi/testify/v2/require"
	"github.com/labstack/echo/v4"
)

// --- helpers ---
var errorTimeout = errors.New("timeout")

func makeInfo(args ...any) domain.DeviceInfo {
	info := model.DeviceInfo{}
	for index := 0; index < len(args); index += 2 {
		info[args[index].(string)] = args[index+1]
	}
	return info
}

type sseReceiver struct {
	resp     io.Reader
	received []string
}

func newSseReceiver(resp io.Reader) *sseReceiver {
	return &sseReceiver{
		resp:     resp,
		received: []string{},
	}
}

func (s sseReceiver) wait(timeout time.Duration) error {
	result := make(chan error, 1)
	go func() {
		scan := bufio.NewScanner(s.resp)
		for scan.Scan() {
			s.received = append(s.received, scan.Text())
		}
		result <- scan.Err()
	}()

	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case <-timer.C:
		return errorTimeout

	case err := <-result:
		if err != nil {
			return err
		}
	}

	return nil
}

// testing

func TestGetMainPage(t *testing.T) {
	h := handler.Factory(database.StubDatabase{}, config.Config{})

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/gui", strings.NewReader(""))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// testing
	require.ErrorIs(t, h.GetMainPage(c), echo.ErrNotFound)
}

type testServer struct {
	s *httptest.Server
}

func newTestServe(method string, url string, h func(c echo.Context) error) *testServer {
	e := echo.New()
	switch method {
	case http.MethodGet:
		e.GET(url, h)
	case http.MethodPost:
		e.POST(url, h)
	}
	return &testServer{
		s: httptest.NewServer(e),
	}
}

func (t testServer) Close() {
	t.s.Close()
}

func request(ctx context.Context, method string, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req)
}

func TestGetDevices(t *testing.T) {
	t.SkipNow()

	const url = "/api/devices"
	// h := handler.NewEchoHandler(config.Config{}, usecase.NewUpdateDeviceInfoUsecase(database.StubDatabase{}))

	t.Run("periodical", func(t *testing.T) {
		// // parepare test server
		// server := newTestServe(http.MethodGet, url, func(c echo.Context) error {
		// 	return h.GetDevices(c, &api.GetDevicesParams{
		// 		Interval: api.PollingInterval(1),
		// 	})
		// })
		// defer server.Close()

		// // prepare test client
		// ctx, cancel := context.WithCancel(context.Background())
		// defer cancel()
		// resp, err := request(ctx, http.MethodGet, url)
		// require.NoError(t, err)

		// // receive events
		// rcv := newSseReceiver(resp.Body)
		// rcv.wait(2500 * time.Millisecond)
		// cancel()

		// // validate test
		// require.Equal(t, http.StatusOK, resp.StatusCode)
		// require.Equal(t, "text/event-stream", resp.Header.Get(echo.HeaderContentType))
		// require.Equal(t, "no-cache", resp.Header.Get(echo.HeaderCacheControl))
		// require.Equal(t, "no-cache", resp.Header.Get(echo.HeaderCacheControl))
		// require.Equal(t, "'no'", resp.Header.Get("X-Accel-Buffering"))
		// t.Log(rcv.received)
		// require.NotEmpty(t, rcv.received)
		// for index := range len(tc.want) {
		// 	require.JSONEq(t, tc.want[index], received[index])
		// }
	})
}
