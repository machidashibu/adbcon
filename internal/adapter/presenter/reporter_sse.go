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

func (r SSEReporter) report(report any) error {
	// marshal to JSON
	data, err := json.Marshal(report)
	if err != nil {
		slog.Error("reporter json marshal error", "err", err, "report", report)
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

// ReportDeviceList reports dvice list.
func (r SSEReporter) ReportDeviceList(devs domain.DeviceList) error {
	return r.report(apiconv.DeviceListToApi(devs))
}

// ReportCommandResult reports command result.
func (r SSEReporter) ReportCommandResult(result domain.CommandResult) error {
	// marshal to JSON
	return r.report(apiconv.CommandResult(result))
}

// ReportClose reports to close SSE connection.
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
