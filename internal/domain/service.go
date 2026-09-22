package domain

import (
	"context"
)

// CommandResult is an output of command.
// It has converting method that from bytes (received raw) to string, string array (each lines).
type CommandResult interface {
	// Bytes provides received raw data by byte array.
	Bytes() []byte
	// String provides strings that is casted raw data.
	String() string
	// Lines provides lines that split by return code.
	Lines() []string
	// IsEmpty returns result is empty or not.
	IsEmpty() bool
}

// CommandCh is a channel of command output.
type CommandCh chan []byte

// CommandExecuter is an interface of service that execute OS command by sync or async.
type CommandExecuter interface {
	// Run executes OS command by sync.
	Run(ctx context.Context) (CommandResult, error)
	// Start executes OS command by async.
	Start(ctx context.Context) (CommandCh, error)
}
