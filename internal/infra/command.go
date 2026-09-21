package infra

import (
	"adbcon/internal/domain"
	"bufio"
	"context"
	"log/slog"
	"os/exec"
)

// Command is a wrapper to execute OS command.
type Command struct {
	name string
	args []string
}

// NewCommand creates object with command name and arguments.
// `name` is a name of command. Supported commands are defined in domain by CommandName type.
// `args` is a  command arguments.
func NewCommand(name string, args ...string) *Command {
	return &Command{
		name: name,
		args: args,
	}
}

// Run execute command by sync.
// Return all output bytes when terminated the command.
// The output is merged stdout and stderr.
func (c Command) Run(ctx context.Context) (domain.CommandResult, error) {
	// prepare
	cmd := exec.CommandContext(ctx, c.name, c.args...)
	cmd.Stderr = cmd.Stdout // merge stdout & stderr

	// execute and get output
	output, err := cmd.Output()
	if err != nil {
		slog.Error("command output error", "err", err, "name", c.name, "args", c.args)
		return nil, err
	}

	return output, nil
}

// Run execute command by async.
// Return channel of output bytes at immediatly.
// Notify output byte that is each lines via channel.
// The output is merged stdout and stderr.
func (c Command) Start(ctx context.Context) (domain.CommandCh, error) {
	// prepare
	cmd := exec.CommandContext(ctx, string(c.name), c.args...)
	r, err := cmd.StdoutPipe()
	if err != nil {
		slog.Error("command pipe error", "err", err)
		return nil, err
	}
	cmd.Stderr = cmd.Stdout // merge stdout & stderr

	output := make(domain.CommandCh)

	// execute
	if err := cmd.Start(); err != nil {
		slog.Error("command start error", "err", err, "name", c.name, "args", c.args)
		return nil, err
	}

	// get output
	go func() {
		defer close(output)
		defer r.Close()

		// scan stdout
		scan := bufio.NewScanner(r)
		for scan.Scan() {
			line := append([]byte(nil), scan.Bytes()...)
			output <- line
		}
		if err := scan.Err(); err != nil {
			slog.Error("command scan error", "err", err)
		}

		// discard resources
		if err := cmd.Wait(); err != nil {
			slog.Error("command wait error", "err", err, "name", c.name, "args", c.args)
		}
	}()

	return output, nil
}
