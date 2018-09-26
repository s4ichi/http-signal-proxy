package main

import (
	"fmt"
	"net/http"
	"os"
	"syscall"
)

const portMaxNum = 65535

var (
	supportSignals = []struct {
		pathStr string
		signal  syscall.Signal
	}{
		{
			pathStr: "sigint",
			signal:  syscall.SIGINT,
		},
		{
			pathStr: "sigterm",
			signal:  syscall.SIGTERM,
		},
		{
			pathStr: "sighup",
			signal:  syscall.SIGHUP,
		},
		{
			pathStr: "sigquit",
			signal:  syscall.SIGQUIT,
		},
		{
			pathStr: "sigusr2",
			signal:  syscall.SIGUSR2,
		},
	}
)

type Httpd struct {
	Port     uint
	Prefix   string
	Callback func(sig syscall.Signal) error
}

func NewHttpd(port uint, prefix string) (*Httpd, error) {
	if port > portMaxNum {
		return nil, fmt.Errorf("port number shold be smaller than %d", portMaxNum)
	}

	return &Httpd{
		Port:   port,
		Prefix: prefix,
	}, nil
}

func (httpd *Httpd) Run() {
	for i, _ := range supportSignals {
		s := supportSignals[i]
		http.HandleFunc(httpd.Prefix+"/"+s.pathStr, func(w http.ResponseWriter, r *http.Request) {
			err := httpd.Callback(s.signal)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Failed to poroxy %s to destination command", s.pathStr)
			} else {
				fmt.Fprintf(w, "Successed to proxy %s to destination command", s.pathStr)
			}
		})
	}

	http.ListenAndServe(fmt.Sprintf(":%d", httpd.Port), nil)
}
