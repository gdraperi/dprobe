// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

// Package debug provides facilities to execute svc.Handler on console.
//
package debug

import (
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sys/windows/svc"
)

// Run executes service name by calling appropriate handler function.
// The process is running on console, unlike real service. Use Ctrl+C to
// send "Stop" command to your service.
func Run(name string, handler svc.Handler) error ***REMOVED***
	cmds := make(chan svc.ChangeRequest)
	changes := make(chan svc.Status)

	sig := make(chan os.Signal)
	signal.Notify(sig)

	go func() ***REMOVED***
		status := svc.Status***REMOVED***State: svc.Stopped***REMOVED***
		for ***REMOVED***
			select ***REMOVED***
			case <-sig:
				cmds <- svc.ChangeRequest***REMOVED***svc.Stop, 0, 0, status***REMOVED***
			case status = <-changes:
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	_, errno := handler.Execute([]string***REMOVED***name***REMOVED***, cmds, changes)
	if errno != 0 ***REMOVED***
		return syscall.Errno(errno)
	***REMOVED***
	return nil
***REMOVED***
