package cli

import (
	"io"
)

const (
	ExitOK    = 0
	ExitError = 1
)

type Main struct {
	command Commander
}

type Commander interface {
	SetOutStream(io.Writer)
	SetErrStream(io.Writer)
	//SetFlagSet(*flag.FlagSet)
	//Name() string
	Run([]string) error
}

func NewMain(command Commander, out, err io.Writer) *Main {
	command.SetOutStream(out)
	command.SetErrStream(err)
	return &Main{
		command: command,
	}
}

func (m *Main) Run(args []string) error {
	return m.command.Run(args)
}
