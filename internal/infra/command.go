package infra

import (
	"adbcon/internal/adapter/model"
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

// NewAdbCommand creates object for ADB command with sub command name and arguments.
// `name` is a name of sub command for ADN. Supported commands are defined in domain by CommandName type.
// `args` is a  command arguments.
func NewAdbCommand(name domain.CommandName, args ...string) *Command {
	obj := &Command{
		name: string(domain.CommandAdb),
		args: []string{string(name)},
	}
	obj.args = append(obj.args, args...)
	return obj
}

// NewAdbCommandWithSerial creates object for ADB command with sub command name, serial and arguments.
// `name` is a name of sub command for ADN. Supported commands are defined in domain by CommandName type.
// `serial` is a serial of devie. It is a `-s` option of ADB command.
// `args` is a  command arguments.
func NewAdbCommandWithSerial(name domain.CommandName, serial string, args ...string) *Command {
	obj := &Command{
		name: string(domain.CommandAdb),
		args: []string{"-s", serial, string(name)},
	}
	obj.args = append(obj.args, args...)
	return obj
}

// Run execute command by sync.
// Return all output bytes when terminated the command.
// The output is merged stdout and stderr.
func (c Command) Run(ctx context.Context, args ...string) (domain.CommandOutput, error) {
	// prepare
	params := append(c.args, args...)
	cmd := exec.CommandContext(ctx, c.name, params...)
	cmd.Stderr = cmd.Stdout // merge stdout & stderr

	// execute and get output
	slog.Info("run command", "name", c.name, "args", params)
	output, err := cmd.Output()
	if err != nil {
		slog.Error("command output error", "err", err, "name", c.name, "args", c.args)
		return nil, err
	}

	return model.CommandOutput(output), nil
}

// Run execute command by async.
// Return channel of output bytes at immediatly.
// Notify output byte that is each lines via channel.
// The output is merged stdout and stderr.
func (c Command) Start(ctx context.Context, args ...string) (domain.CommandCh, error) {
	// prepare
	params := append(c.args, args...)
	cmd := exec.CommandContext(ctx, string(c.name), params...)
	r, err := cmd.StdoutPipe()
	if err != nil {
		slog.Error("command pipe error", "err", err)
		return nil, err
	}
	cmd.Stderr = cmd.Stdout // merge stdout & stderr

	output := make(domain.CommandCh)

	// execute
	slog.Info("start command", "name", c.name, "args", params)
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
			output <- model.NewCommandResult("", line)
		}
		if err := scan.Err(); err != nil {
			output <- model.NewCommandErrorResult("", err)
		}

		// dispose resources
		if err := cmd.Wait(); err != nil {
			output <- model.NewCommandErrorResult("", err)
		}
	}()

	return output, nil
}
