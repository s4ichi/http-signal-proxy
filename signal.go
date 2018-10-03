package main

import (
	"syscall"
)

var (
	Signals = []syscall.Signal{
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
		syscall.SIGQUIT,
		syscall.SIGWINCH,
		syscall.SIGUSR1,
		syscall.SIGUSR2,
		syscall.SIGTTIN,
		syscall.SIGTTOU,
	}
	SignalStrs = []string{
		"sigint",
		"sigterm",
		"sighup",
		"sigquit",
		"sigwinch",
		"sigusr1",
		"sigusr2",
		"sigttin",
		"sigttou",
	}
)
