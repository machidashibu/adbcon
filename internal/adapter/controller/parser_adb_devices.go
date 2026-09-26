package controller

import (
	"adbcon/internal/adapter/model"
	"adbcon/internal/domain"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
)

// AdbDeviceParseris a parser that parse output of `adb devices` commnd.
type AdbDeviceParser struct{}

// Parse parses command result to domain data.
func (a AdbDeviceParser) Parse(result domain.CommandOutput) (domain.DeviceList, error) {
	var errorInvalidformat = fmt.Errorf("invalid format: adb devices")

	// validate
	if result.IsEmpty() {
		return nil, errorInvalidformat
	}

	// prepare
	devs := domain.DeviceList{}
	lines := result.Lines()

	// parse 1st line (fixed string)
	if lines[0] != "List of devices attached" {
		slog.Error("invalid format at 1st line", "text", lines[0])
		return nil, errorInvalidformat
	}

	// parse device information lines
	// ex. ABC123DEF             device product:F-52G model:F_52G device:F-52G transport_id:1
	// TODO: support when not established adb server
	// 	* daemon not running; starting now at tcp:5037
	//  * daemon started successfully
	for _, text := range lines[1:] {
		// split by spece
		fields := strings.Fields(text)
		if len(fields) < 2 {
			slog.Warn("invalid format", "text", text)
			continue // skip when unmatched line
		}

		info := model.DeviceInfo{}
		// get serial
		info[model.LabelSerial] = fields[0]
		// get status
		info[model.LabelStatus] = DeviceStatusToDomain(fields[1])

		if len(fields) < 3 {
			continue // short format
		}

		// get option params
		for _, field := range fields[2:] {
			// spit by collon
			vals := strings.Split(field, ":")
			if len(vals) < 2 {
				slog.Warn("unknown field", "field", field)
				continue // value only field
			}

			// repaire value field
			val := strings.Join(vals[1:], ":")
			// store number if convertable
			num, err := strconv.Atoi(val)
			if err != nil {
				info[vals[0]] = val
			} else {
				info[vals[0]] = num
			}
		}

		slog.Debug("add device infromation", "info", info)
		devs = append(devs, info)
	}

	return devs, nil
}
