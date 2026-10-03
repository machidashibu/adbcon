package service

import (
	"adbcon/internal/adapter/model"
	"adbcon/internal/domain"
	"context"
	"log/slog"
	"sync"
)

type MultiCommandProvider struct {
	cmds map[string]domain.AsyncCommandExecuter
}

func NewMultiCommandProvider() *MultiCommandProvider {
	return &MultiCommandProvider{
		cmds: map[string]domain.AsyncCommandExecuter{},
	}
}

func (m *MultiCommandProvider) Add(serial string, cmd domain.AsyncCommandExecuter) *MultiCommandProvider {
	m.cmds[serial] = cmd
	return m
}

func (m MultiCommandProvider) Start(ctx context.Context, args ...string) (domain.CommandCh, error) {
	report := make(domain.CommandCh)
	wg := sync.WaitGroup{}

	// execute all command by async.
	for serial, cmd := range m.cmds {
		wg.Add(1)
		go func(serial string, name domain.AsyncCommandExecuter) {
			defer wg.Done()

			// execute command
			ch, err := cmd.Start(ctx, args...)
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
				case result, ok := <-ch:
					if !ok {
						return // terminate command
					}
					report <- model.NewCommandResult(serial, []byte(result.Text()))
				}
			}
		}(serial, cmd)
	}

	// wait end of all command.
	go func() {
		wg.Wait()
		close(report)
	}()

	return report, nil
}
