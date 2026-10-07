package domain

import (
	"context"
)

// CommandExecuter is an interface of service that execute OS command by sync
type SyncCommandExecuter interface {
	// Run executes OS command by sync.
	Run(ctx context.Context, args ...string) (CommandOutput, error)
}

// CommandExecuter is an interface of service that execute OS command by async.
type AsyncCommandExecuter interface {
	// Start executes OS command by async.
	Start(ctx context.Context, args ...string) (CommandCh, error)
}

// CommandExecuter is an interface of service that execute OS command by sync or async.
type CommandExecuter interface {
	SyncCommandExecuter
	AsyncCommandExecuter
}

type DeviceListResolver interface {
	GetDeviceList(ctx context.Context) (DeviceList, error)
}

type VersionResolver interface {
	GetVersion(ctx context.Context) (Version, error)
}
