package main

import (
	"fmt"
	"os"
	"os/exec"
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

func init() {

}

func main() {
	os.Exit(runMain())
}

func runMain() int {
	_, err := flags.Parse(&opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Failed to parse option values: %s\n", err.Error())
		return 1
	}

	pid := syscall.Getpid()
	signalCh := make(chan os.Signal, 1)
	exitErrCh := make(chan *exec.ExitError)

	command := NewCommand(opts.CmdString, exitErrCh)
	go command.Execute()

	handleSignal(signalCh)
	go command.ProxySignal(signalCh)

	httpd, err := NewHttpd(opts.Port, opts.Prefix)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Failed to initialize Httpd: %s\n", err.Error())
		return 1
	}

	httpd.Callback = func(sig syscall.Signal) error {
		return syscall.Kill(pid, sig)
	}
	go httpd.Run()

	if exitErr := <-exitErrCh; exitErr != nil {
		waitStatus := exitErr.Sys().(syscall.WaitStatus)
		fmt.Fprintf(os.Stderr, "[ERROR] Failed to exit command: %s : exitcode=%d\n", exitErr.Error(), waitStatus.ExitStatus())
		return 1
	}

	return 0
}

func handleSignal(signalCh chan os.Signal) {
	signal.Notify(signalCh)
}
