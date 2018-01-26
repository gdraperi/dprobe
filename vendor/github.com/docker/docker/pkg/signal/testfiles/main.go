package main

import (
	"os"
	"syscall"
	"time"

	"github.com/docker/docker/pkg/signal"
	"github.com/sirupsen/logrus"
)

func main() ***REMOVED***
	sigmap := map[string]os.Signal***REMOVED***
		"TERM": syscall.SIGTERM,
		"QUIT": syscall.SIGQUIT,
		"INT":  os.Interrupt,
	***REMOVED***
	signal.Trap(func() ***REMOVED***
		time.Sleep(time.Second)
		os.Exit(99)
	***REMOVED***, logrus.StandardLogger())
	go func() ***REMOVED***
		p, err := os.FindProcess(os.Getpid())
		if err != nil ***REMOVED***
			panic(err)
		***REMOVED***
		s := os.Getenv("SIGNAL_TYPE")
		multiple := os.Getenv("IF_MULTIPLE")
		switch s ***REMOVED***
		case "TERM", "INT":
			if multiple == "1" ***REMOVED***
				for ***REMOVED***
					p.Signal(sigmap[s])
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				p.Signal(sigmap[s])
			***REMOVED***
		case "QUIT":
			p.Signal(sigmap[s])
		***REMOVED***
	***REMOVED***()
	time.Sleep(2 * time.Second)
***REMOVED***
