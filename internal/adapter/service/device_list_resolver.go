package service

import (
	"adbcon/internal/domain"
	"context"
	"log/slog"
)

type deviceListparser interface {
	Parse(result domain.CommandOutput) (domain.DeviceList, error)
}

type DeviceListResulver struct {
	cmd    domain.SyncCommandExecuter
	parser deviceListparser
}

func NewDeviceListResulver(cmd domain.SyncCommandExecuter, parser deviceListparser) *DeviceListResulver {
	return &DeviceListResulver{
		cmd:    cmd,
		parser: parser,
	}
}

func (d DeviceListResulver) GetDeviceList(ctx context.Context) (domain.DeviceList, error) {
	out, err := d.cmd.Run(ctx)
	if err != nil {
		slog.Error("command error", "err", err)
		return nil, err
	}

	list, err := d.parser.Parse(out)
	if err != nil {
		slog.Error("parse error", "err", err, "out", out)
		return nil, err
	}

	return list, nil
}
