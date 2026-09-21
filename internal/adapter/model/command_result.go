package model

import (
	"bufio"
	"bytes"
)

type CommandResult []byte

func (cr CommandResult) Bytes() []byte {
	return cr
}

func (cr CommandResult) String() string {
	return string(cr)
}

func (cr CommandResult) Lines() []string {
	lines := []string{}
	scan := bufio.NewScanner(bytes.NewReader(cr))
	for scan.Scan() {
		lines = append(lines, scan.Text())
	}
	return lines
}

func (cr CommandResult) IsEmpty() bool {
	return len(cr) == 0
}
