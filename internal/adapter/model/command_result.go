package model

type CommandResult struct {
	serial string
	src    []byte
	err    error
}

func NewCommandResult(serial string, src []byte) *CommandResult {
	return &CommandResult{
		serial: serial,
		src:    src,
	}
}

func NewCommandErrorResult(serial string, err error) *CommandResult {
	return &CommandResult{
		serial: serial,
		err:    err,
	}
}

func (c CommandResult) Serial() string {
	return c.serial
}

func (c CommandResult) Text() string {
	if c.err != nil {
		return c.err.Error()
	}
	return string(c.src)
}

func (c CommandResult) Error() error {
	return c.err
}
