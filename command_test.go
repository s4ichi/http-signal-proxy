package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"syscall"
	"time"

	"testing"
)

func TestNewCommand(t *testing.T) {
	dummyCh := make(chan error)
	defer close(dummyCh)

	cases := []struct {
		cmdString string
		success   bool
	}{
		{
			cmdString: "",
			success:   false,
		},
		{
			cmdString: "tail -f production.log",
			success:   true,
		},
		{
			cmdString: "make",
			success:   true,
		},
	}

	for _, c := range cases {
		_, err := NewCommand(c.cmdString, dummyCh, dummyCh)
		if (err == nil) != c.success {
			t.Fatalf("unexpected result with error: %s", err.Error())
		}
	}
}

func TestHandleSignal(t *testing.T) {
	timeout := 2 * time.Second
	signalCh := make(chan os.Signal)
	defer close(signalCh)

	readyCh := make(chan error)
	exitCh := make(chan error)
	defer close(readyCh)
	defer close(exitCh)

	command, _ := NewCommand(fmt.Sprintf("sleep %d", int(timeout.Seconds())+1), readyCh, exitCh)

	go command.Execute()
	if err := <-readyCh; err != nil {
		t.Fatalf("failed to execute command: %s", err.Error())
	}

	go command.HandleSignal(signalCh)
	signalCh <- syscall.SIGTERM

	select {
	case <-time.After(timeout): // Fail
	case <-exitCh: // Success
		return
	}

	t.Fatalf("cannot proxy signal to command within %s", timeout.String())
}

func TestExecute(t *testing.T) {
	cases := []struct {
		cmdStr string
		ready  bool
		exit   bool
	}{
		{
			cmdStr: "sleep 1",
			ready:  true,
			exit:   true,
		},
		{
			cmdStr: "dummy_dummy_command",
			ready:  false,
			exit:   false, // Both are available true or false
		},
		{
			cmdStr: "dummy_dummy_command",
			ready:  false,
			exit:   true, // Both are available true or false
		},
		{
			cmdStr: "ln", // It will be exit with statu 1.
			ready:  true,
			exit:   false,
		},
	}

	for _, c := range cases {
		var err error
		timeout := 2 * time.Second
		signalCh := make(chan os.Signal)

		readyCh := make(chan error)
		exitCh := make(chan error)
		command, _ := NewCommand(c.cmdStr, readyCh, exitCh)

		// Discard output of command for testing
		command.ExecCmd.Stdout = ioutil.Discard
		command.ExecCmd.Stderr = ioutil.Discard

		go command.Execute()
		if err = <-readyCh; (err == nil) != c.ready {
			if c.ready {
				t.Fatalf("unexpected ready state\nactual: %s\nexpected: %s\ncommand: %s\nerror: %s", strconv.FormatBool(!c.ready), strconv.FormatBool(c.ready), c.cmdStr, err.Error())
			} else {
				t.Fatalf("unexpected ready state\nactual: %s\nexpected: %s\ncommand: %s", strconv.FormatBool(!c.ready), strconv.FormatBool(c.ready), c.cmdStr)
			}
		}

		// If failed to start command, we cannot process continuation.
		if err != nil {
			continue
		}

		go command.HandleSignal(signalCh)
		time.Sleep(100 * time.Millisecond) // Warm up command to get exitcode=1
		signalCh <- syscall.SIGINT

		select {
		case <-time.After(timeout):
			t.Fatalf("unexpected timeout to exit command: %s", c.cmdStr)
		case err = <-exitCh:
			if (err == nil) != c.exit {
				if c.exit {
					t.Fatalf("unexpected exit state\nactual: %s\nexpected: %s\ncommand: %s\nerror: %s", strconv.FormatBool(!c.exit), strconv.FormatBool(c.exit), c.cmdStr, err.Error())
				} else {
					t.Fatalf("unexpected exit state\nactual: %s\nexpected: %s\ncommand: %s", strconv.FormatBool(!c.exit), strconv.FormatBool(c.exit), c.cmdStr)
				}
			}
		}

		close(signalCh)
		close(readyCh)
		close(exitCh)
	}
}
