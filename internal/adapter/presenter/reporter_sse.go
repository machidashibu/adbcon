package presenter

import (
	"adbcon/internal/adapter/presenter/apiconv"
	"adbcon/internal/domain"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"

	"github.com/labstack/echo/v4"
)

// SSEReporter reports domain data to user via SSE stream
type SSEReporter struct {
	mu   *sync.Mutex
	resp *echo.Response
}

// NewSSEReporter creates object.
func NewSSEReporter(resp *echo.Response) *SSEReporter {
	resp.Header().Set(echo.HeaderContentType, "text/event-stream")
	resp.Header().Set(echo.HeaderCacheControl, "no-cache")
	resp.Header().Set(echo.HeaderConnection, "keep-alive")
	resp.Header().Set("X-Accel-Buffering", "'no'")

	return &SSEReporter{
		mu:   &sync.Mutex{},
		resp: resp,
	}
}

// ReportDeviceList reports dvice list.
func (r SSEReporter) ReportDeviceList(devs domain.DeviceList) error {
	// marshal to JSON
	apiData := apiconv.DeviceListToApi(devs)
	data, err := json.Marshal(apiData)
	if err != nil {
		slog.Error("reporter json marshal error", "err", err, "data", apiData)
		return err
	}

	// report
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, err := fmt.Fprintf(r.resp.Writer, "data: %s\n\n", string(data)); err != nil {
		slog.Error("reporter print error", "err", err)
		return err
	}
	r.resp.Flush()

	return nil
}

func (r SSEReporter) ReportClose() error {
	// report
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, err := fmt.Fprintf(r.resp.Writer, "event: close\ndata: SSE connection closed.\n\n"); err != nil {
		slog.Error("reporter print error", "err", err)
		return err
	}
	r.resp.Flush()

	return nil
}
