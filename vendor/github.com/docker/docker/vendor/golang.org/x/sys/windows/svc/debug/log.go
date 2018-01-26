// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package debug

import (
	"os"
	"strconv"
)

// Log interface allows different log implementations to be used.
type Log interface ***REMOVED***
	Close() error
	Info(eid uint32, msg string) error
	Warning(eid uint32, msg string) error
	Error(eid uint32, msg string) error
***REMOVED***

// ConsoleLog provides access to the console.
type ConsoleLog struct ***REMOVED***
	Name string
***REMOVED***

// New creates new ConsoleLog.
func New(source string) *ConsoleLog ***REMOVED***
	return &ConsoleLog***REMOVED***Name: source***REMOVED***
***REMOVED***

// Close closes console log l.
func (l *ConsoleLog) Close() error ***REMOVED***
	return nil
***REMOVED***

func (l *ConsoleLog) report(kind string, eid uint32, msg string) error ***REMOVED***
	s := l.Name + "." + kind + "(" + strconv.Itoa(int(eid)) + "): " + msg + "\n"
	_, err := os.Stdout.Write([]byte(s))
	return err
***REMOVED***

// Info writes an information event msg with event id eid to the console l.
func (l *ConsoleLog) Info(eid uint32, msg string) error ***REMOVED***
	return l.report("info", eid, msg)
***REMOVED***

// Warning writes an warning event msg with event id eid to the console l.
func (l *ConsoleLog) Warning(eid uint32, msg string) error ***REMOVED***
	return l.report("warn", eid, msg)
***REMOVED***

// Error writes an error event msg with event id eid to the console l.
func (l *ConsoleLog) Error(eid uint32, msg string) error ***REMOVED***
	return l.report("error", eid, msg)
***REMOVED***
