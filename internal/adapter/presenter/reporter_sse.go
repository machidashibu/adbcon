package presenter

import (
	"adbcon/internal/adapter/presenter/toapi"
	"adbcon/internal/domain"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

// SSEReporter reports domain data to user via SSE stream
type SSEReporter struct {
	mu   *sync.Mutex
	resp *echo.Response

	buffer   []string
	latest   time.Time
	interval time.Duration
	stay     int
}

// NewSSEReporter creates object.
func NewSSEReporter(resp *echo.Response, interval time.Duration, stay int) *SSEReporter {
	resp.Header().Set(echo.HeaderContentType, "text/event-stream")
	resp.Header().Set(echo.HeaderCacheControl, "no-cache")
	resp.Header().Set(echo.HeaderConnection, "keep-alive")
	resp.Header().Set("X-Accel-Buffering", "'no'")

	return &SSEReporter{
		mu:       &sync.Mutex{},
		resp:     resp,
		buffer:   make([]string, 0, stay),
		interval: interval,
		stay:     stay,
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
	r.buffer = append(r.buffer, fmt.Sprintf("data: %s", string(data)))

	if len(r.buffer) >= r.stay || time.Since(r.latest) > r.interval {
		for _, line := range r.buffer {
			fmt.Fprintln(r.resp.Writer, line)
		}
		fmt.Fprintln(r.resp.Writer)
		r.resp.Flush()
		r.latest = time.Now()
		r.buffer = make([]string, 0, r.stay)
	}

	return nil
}

// ReportDeviceList reports dvice list.
func (r SSEReporter) ReportDeviceList(devs domain.DeviceList) error {
	return r.report(toapi.DeviceList(devs))
}

// ReportCommandResult reports command result.
func (r SSEReporter) ReportCommandResult(result domain.CommandResult) error {
	// marshal to JSON
	return r.report(toapi.CommandResult(result))
}

func (r SSEReporter) ReportVersion(version domain.Version) error {
	// TODO:
	return r.report(toapi.Version(version))
}

// ReportClose reports to close SSE connection.
func (r SSEReporter) ReportClose() error {
	// report
	r.mu.Lock()
	defer r.mu.Unlock()
	// send data if rest data in buffer
	if len(r.buffer) > 0 {
		for _, line := range r.buffer {
			fmt.Fprintln(r.resp.Writer, line)
		}
		fmt.Fprintln(r.resp.Writer)
		r.resp.Flush()
		r.latest = time.Now()
		r.buffer = make([]string, 0, r.stay)
	}

	// send close event with data
	fmt.Fprintf(r.resp.Writer, "event: close\ndata: SSE connection closed.\n\n")
	r.resp.Flush()

	return nil
}
