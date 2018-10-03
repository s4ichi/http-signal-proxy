package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	flags "github.com/jessevdk/go-flags"
)

var (
	cmdString string
	logLevel  string
	port      int

	opts struct {
		CmdString string `short:"c" long:"command" description:"Destination command to proxying signals" required:"true"`
		Port      uint   `short:"p" long:"port" description:"Port number to listen http" required:"true"`
		Prefix    string `short:"f" long:"prefix" description:"Prefix path of url to listen http" default:"/http-signal"`
	}
)

func main() {
	os.Exit(runMain())
}

func runMain() int {
	_, err := flags.Parse(&opts)
	if err != nil {
		return 1
	}

	pid := syscall.Getpid()
	signalCh := make(chan os.Signal)
	readyCh := make(chan error)
	exitCh := make(chan error)

	defer close(signalCh)
	defer close(readyCh)
	defer close(exitCh)

	// Create new command with option
	command, err := NewCommand(opts.CmdString, readyCh, exitCh)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Failed to create command: %s\n", err.Error())
		return 1
	}

	// Execute command and wait to start them.
	go command.Execute()
	if err = <-readyCh; err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Failed to execute command: %s\n", err.Error())
		return 1
	}

	// Hadling signal in this process and command
	go command.HandleSignal(signalCh)
	handleSignal(signalCh)

	// Setup http server to proxy signal to command
	httpd, err := NewHttpd(opts.Port, opts.Prefix, proxySignalFunc(pid))
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Failed to initialize Httpd: %s\n", err.Error())
		return 1
	}

	// Start http server
	go httpd.Run()

	// Waiting to exit command
	if err = <-exitCh; err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Failed to exit command: %s\n", err.Error())
		return 1
	}

	fmt.Fprintf(os.Stdout, "Successed to exit command. Bye.\n")
	return 0
}

// handling all signal and pass to signalCh
func handleSignal(signalCh chan os.Signal) {
	signal.Notify(signalCh)
}

// proxySignalFunc return fanction that receive signal name and send it to pid.
func proxySignalFunc(pid int) func(syscall.Signal) error {
	return func(sig syscall.Signal) error {
		return syscall.Kill(pid, sig)
	}
}
