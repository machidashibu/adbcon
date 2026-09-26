package service

import (
	"adbcon/internal/adapter/model"
	"adbcon/internal/domain"
	"context"
	"log/slog"
	"sync"
)

type MultiCommand struct {
	cmds map[string]domain.CommandExecuter
}

func NewMultiCommand() *MultiCommand {
	return &MultiCommand{
		cmds: map[string]domain.CommandExecuter{},
	}
}

func (m *MultiCommand) Add(serial string, cmd domain.CommandExecuter) *MultiCommand {
	m.cmds[serial] = cmd
	return m
}

func (m MultiCommand) Start(ctx context.Context) chan domain.CommandResult {
	report := make(chan domain.CommandResult)
	wg := sync.WaitGroup{}

	// execute all command by async.
	for serial, cmd := range m.cmds {
		wg.Add(1)
		go func(serial string, cmd domain.CommandExecuter) {
			defer wg.Done()

			// execute command
			ch, err := cmd.Start(ctx)
			if err != nil {
				slog.Error("command start error", "err", err, "serial", serial)
				report <- model.NewCommandErrorResult(serial, err)
				return
			}

			// receive and report command output
			for {
				select {
				case <-ctx.Done():
					if ctx.Err() != nil {
						report <- model.NewCommandErrorResult(serial, ctx.Err())
					}
					return // cancel command
				case output, ok := <-ch:
					if !ok {
						return // terminate command
					}
					report <- model.NewCommandResult(serial, output)
				}
			}
		}(serial, cmd)
	}

	// wait end of all command.
	go func() {
		wg.Wait()
		close(report)
	}()

	return report
}
