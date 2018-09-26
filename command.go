package main

import (
	"os"
	"os/exec"
	"strings"
)

type Command struct {
	ExitErrCh chan *exec.ExitError
	ExecCmd   *exec.Cmd
}

func NewCommand(cmdString string, exitErrCh chan *exec.ExitError) *Command {
	cmdSplit := strings.Fields(cmdString)
	cmd := exec.Command(cmdSplit[0], cmdSplit[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return &Command{
		ExitErrCh: exitErrCh,
		ExecCmd:   cmd,
	}
}

func (command *Command) ProxySignal(signalCh chan os.Signal) {
	for sig := range signalCh {
		command.ExecCmd.Process.Signal(sig)
	}
}

func (command *Command) Execute() {
	err := command.ExecCmd.Run()

	if exitError, ok := err.(*exec.ExitError); ok {
		command.ExitErrCh <- exitError
	} else {
		command.ExitErrCh <- nil
	}
}
