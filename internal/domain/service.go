package domain

import (
	"bufio"
	"bytes"
	"context"
)

type CommandResult []byte

func (cr CommandResult) Lines() []string {
	lines := []string{}
	scan := bufio.NewScanner(bytes.NewReader(cr))
	for scan.Scan() {
		lines = append(lines, scan.Text())
	}
	return lines
}

type CommandCh chan []byte

type CommandExecuter interface {
	Run(ctx context.Context) (CommandResult, error)
	Start(ctx context.Context) (CommandCh, error)
}
