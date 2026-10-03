package usecase

import (
	"adbcon/internal/domain"
	"context"
	"log/slog"
)

type executeAdbCommandRepoter interface {
	ReportCommandResult(result domain.CommandResult) error
}

type ExecuteAdbCommand struct{}

func NewExecuteAdbCommand() *ExecuteAdbCommand {
	return &ExecuteAdbCommand{}
}

func (uc ExecuteAdbCommand) Execute(ctx context.Context, reporter executeAdbCommandRepoter, cmd domain.AsyncCommandExecuter, args ...string) error {
	slog.Debug("ExecuteAdbCommand::Execute", "cmd", cmd, "args", args)

	// start command
	out, err := cmd.Start(ctx, args...)
	if err != nil {
		slog.Error("usecase execute adb command start error", "args", args)
		return err
	}

	// receive command output by each line, and wait exit.
	for {
		select {
		case <-ctx.Done():
			return nil
		case result, ok := <-out:
			if !ok {
				return nil
			}
			if err := reporter.ReportCommandResult(result); err != nil {
				slog.Error("usecase execute adb command report error", "err", err, "result", result)
				return err
			}
		}

	}
}
