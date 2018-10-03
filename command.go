package main

import (
	"errors"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

type Command struct {
	ExecCmd *exec.Cmd

	ReadyCh chan error
	ExitCh  chan error
}

func NewCommand(cmdString string, readyCh chan error, exitCh chan error) (*Command, error) {
	if cmdString == "" {
		return nil, errors.New("command must be non empty value")
	}

	cmdSplit := strings.Fields(cmdString)
	cmd := exec.Command(cmdSplit[0], cmdSplit[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return &Command{
		ExecCmd: cmd,
		ReadyCh: readyCh,
		ExitCh:  exitCh,
	}, nil
}

func (command *Command) HandleSignal(signalCh chan os.Signal) {
	for sig := range signalCh {
		// If command.ExecCmd.Process is running, command.ExecCmd.ProcessState = nil.
		// After call command.ExecCmd.Process.Wait(), command.ExecCmd.ProcessState != nil
		if command.ExecCmd.ProcessState == nil {
			command.ExecCmd.Process.Signal(sig)
		}
	}
}

func (command *Command) Execute() {
	err := command.ExecCmd.Start()
	command.ReadyCh <- err

	err = command.ExecCmd.Wait()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus := exitError.Sys().(syscall.WaitStatus)
			if !waitStatus.Signaled() {
				command.ExitCh <- err
				return
			}
		}
	}

	command.ExitCh <- nil
}
