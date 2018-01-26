// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

// Package eventlog implements access to Windows event log.
//
package eventlog

import (
	"errors"
	"syscall"

	"golang.org/x/sys/windows"
)

// Log provides access to the system log.
type Log struct ***REMOVED***
	Handle windows.Handle
***REMOVED***

// Open retrieves a handle to the specified event log.
func Open(source string) (*Log, error) ***REMOVED***
	return OpenRemote("", source)
***REMOVED***

// OpenRemote does the same as Open, but on different computer host.
func OpenRemote(host, source string) (*Log, error) ***REMOVED***
	if source == "" ***REMOVED***
		return nil, errors.New("Specify event log source")
	***REMOVED***
	var s *uint16
	if host != "" ***REMOVED***
		s = syscall.StringToUTF16Ptr(host)
	***REMOVED***
	h, err := windows.RegisterEventSource(s, syscall.StringToUTF16Ptr(source))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &Log***REMOVED***Handle: h***REMOVED***, nil
***REMOVED***

// Close closes event log l.
func (l *Log) Close() error ***REMOVED***
	return windows.DeregisterEventSource(l.Handle)
***REMOVED***

func (l *Log) report(etype uint16, eid uint32, msg string) error ***REMOVED***
	ss := []*uint16***REMOVED***syscall.StringToUTF16Ptr(msg)***REMOVED***
	return windows.ReportEvent(l.Handle, etype, 0, eid, 0, 1, 0, &ss[0], nil)
***REMOVED***

// Info writes an information event msg with event id eid to the end of event log l.
// When EventCreate.exe is used, eid must be between 1 and 1000.
func (l *Log) Info(eid uint32, msg string) error ***REMOVED***
	return l.report(windows.EVENTLOG_INFORMATION_TYPE, eid, msg)
***REMOVED***

// Warning writes an warning event msg with event id eid to the end of event log l.
// When EventCreate.exe is used, eid must be between 1 and 1000.
func (l *Log) Warning(eid uint32, msg string) error ***REMOVED***
	return l.report(windows.EVENTLOG_WARNING_TYPE, eid, msg)
***REMOVED***

// Error writes an error event msg with event id eid to the end of event log l.
// When EventCreate.exe is used, eid must be between 1 and 1000.
func (l *Log) Error(eid uint32, msg string) error ***REMOVED***
	return l.report(windows.EVENTLOG_ERROR_TYPE, eid, msg)
***REMOVED***
