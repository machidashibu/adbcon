package usecase

import (
	"adbcon/internal/domain"
	"context"
	"log/slog"
	"time"
)

type PeriodicalUpdateDeviceListUsecase struct {
	ucUpdate *UpdateDeviceListUsecase
}

func NewPeriodicalUpdateDeviceListUsecase(uc *UpdateDeviceListUsecase) *PeriodicalUpdateDeviceListUsecase {
	return &PeriodicalUpdateDeviceListUsecase{
		ucUpdate: uc,
	}
}

func (uc PeriodicalUpdateDeviceListUsecase) Update(ctx context.Context, reporter updateDeviceListReporter, interval domain.Interval) error {
	slog.Debug("PeriodicalUpdateDeviceListUsecase::Update", "interval", interval)

	// 1st update (immediate)
	if err := uc.ucUpdate.Update(ctx, reporter, interval); err != nil {
		return err
	}

	// start polling
	polling := time.NewTicker(time.Duration(interval))
	defer polling.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-polling.C:
			// periodical update
			if err := uc.ucUpdate.Update(ctx, reporter, interval); err != nil {
				return err
			}
		}
	}
}
