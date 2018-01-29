// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

// Example service program that beeps.
//
// The program demonstrates how to create Windows service and
// install / remove it on a computer. It also shows how to
// stop / start / pause / continue any service, and how to
// write to event log. It also shows how to use debug
// facilities available in debug package.
//
package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/sys/windows/svc"
)

func usage(errmsg string) ***REMOVED***
	fmt.Fprintf(os.Stderr,
		"%s\n\n"+
			"usage: %s <command>\n"+
			"       where <command> is one of\n"+
			"       install, remove, debug, start, stop, pause or continue.\n",
		errmsg, os.Args[0])
	os.Exit(2)
***REMOVED***

func main() ***REMOVED***
	const svcName = "myservice"

	isIntSess, err := svc.IsAnInteractiveSession()
	if err != nil ***REMOVED***
		log.Fatalf("failed to determine if we are running in an interactive session: %v", err)
	***REMOVED***
	if !isIntSess ***REMOVED***
		runService(svcName, false)
		return
	***REMOVED***

	if len(os.Args) < 2 ***REMOVED***
		usage("no command specified")
	***REMOVED***

	cmd := strings.ToLower(os.Args[1])
	switch cmd ***REMOVED***
	case "debug":
		runService(svcName, true)
		return
	case "install":
		err = installService(svcName, "my service")
	case "remove":
		err = removeService(svcName)
	case "start":
		err = startService(svcName)
	case "stop":
		err = controlService(svcName, svc.Stop, svc.Stopped)
	case "pause":
		err = controlService(svcName, svc.Pause, svc.Paused)
	case "continue":
		err = controlService(svcName, svc.Continue, svc.Running)
	default:
		usage(fmt.Sprintf("invalid command %s", cmd))
	***REMOVED***
	if err != nil ***REMOVED***
		log.Fatalf("failed to %s %s: %v", cmd, svcName, err)
	***REMOVED***
	return
***REMOVED***
