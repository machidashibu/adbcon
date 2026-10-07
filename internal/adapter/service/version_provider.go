package service

import (
	"adbcon/internal/domain"
	"context"
	"log/slog"
)

type versionparser interface {
	Parse(result domain.CommandOutput) (domain.Version, error)
}

type VersionProvider struct {
	cmd    domain.SyncCommandExecuter
	parser versionparser
}

func NewVersionProvider(cmd domain.SyncCommandExecuter, parser versionparser) *VersionProvider {
	return &VersionProvider{
		cmd:    cmd,
		parser: parser,
	}
}

func (v VersionProvider) GetVersion(ctx context.Context) (domain.Version, error) {
	out, err := v.cmd.Run(ctx)
	if err != nil {
		slog.Error("command error", "err", err)
		return nil, err
	}

	version, err := v.parser.Parse(out)
	if err != nil {
		slog.Error("parse error", "err", err, "out", out)
		return nil, err
	}

	return version, nil
}
