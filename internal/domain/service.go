package domain

import (
	"context"
)

type CommandResult interface {
	Bytes() []byte
	String() string
	Lines() []string
	IsEmpty() bool
}

type CommandCh chan []byte

type CommandExecuter interface {
	Run(ctx context.Context) (CommandResult, error)
	Start(ctx context.Context) (CommandCh, error)
}
