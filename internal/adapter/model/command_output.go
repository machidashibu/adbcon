package model

import (
	"bufio"
	"bytes"
)

type CommandOutput []byte

func (c CommandOutput) Bytes() []byte {
	return c
}

func (c CommandOutput) String() string {
	return string(c)
}

func (c CommandOutput) Lines() []string {
	lines := []string{}
	scan := bufio.NewScanner(bytes.NewReader(c))
	for scan.Scan() {
		lines = append(lines, scan.Text())
	}
	return lines
}

func (c CommandOutput) IsEmpty() bool {
	return len(c) == 0
}
