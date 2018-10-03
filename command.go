package main

import (
	"errors"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// Command struct expressses background process that constructed by exec.Cmd
// If start Command, Command struct sends status to ReadyCh (error OR nil)
// If exit Command, Command struct sends status to ExitCh (error OR nil)
type Command struct {
	ExecCmd *exec.Cmd
	ReadyCh chan error
	ExitCh  chan error
}

// NewCommand crates new command instance with option
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

// HandleSignal function handles signal that sent through signalCh.
func (command *Command) HandleSignal(signalCh chan os.Signal) {
	for sig := range signalCh {
		// If command.ExecCmd.Process is running, command.ExecCmd.ProcessState = nil.
		// After call command.ExecCmd.Process.Wait(), command.ExecCmd.ProcessState != nil
		if command.ExecCmd.ProcessState == nil {
			command.ExecCmd.Process.Signal(sig)
		}
	}
}

// Execute function Start and Wait command
// If exit by non signal cause, send error through ExitCh
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
