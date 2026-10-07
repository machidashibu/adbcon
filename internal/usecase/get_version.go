package usecase

import (
	"adbcon/internal/domain"
	"context"
	"log/slog"
)

type versionReporter interface {
	ReportVersion(version domain.Version) error
}

type GetVersionUsecase struct {
	cmd domain.VersionResolver
}

func NewGetVersionUsecase(cmd domain.VersionResolver) *GetVersionUsecase {
	return &GetVersionUsecase{
		cmd: cmd,
	}
}

func (uc GetVersionUsecase) Get(ctx context.Context, reporter versionReporter) error {
	slog.Error("GetVersionUsecase::Get")

	// get version
	version, err := uc.cmd.GetVersion(ctx)
	if err != nil {
		slog.Error("usecase get version error", "err", err)
		return err
	}

	// report
	if err := reporter.ReportVersion(version); err != nil {
		slog.Error("usecase report version error", "err", err)
		return err
	}

	return nil
}
