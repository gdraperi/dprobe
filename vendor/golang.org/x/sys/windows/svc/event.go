// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package svc

import (
	"errors"

	"golang.org/x/sys/windows"
)

// event represents auto-reset, initially non-signaled Windows event.
// It is used to communicate between go and asm parts of this package.
type event struct ***REMOVED***
	h windows.Handle
***REMOVED***

func newEvent() (*event, error) ***REMOVED***
	h, err := windows.CreateEvent(nil, 0, 0, nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &event***REMOVED***h: h***REMOVED***, nil
***REMOVED***

func (e *event) Close() error ***REMOVED***
	return windows.CloseHandle(e.h)
***REMOVED***

func (e *event) Set() error ***REMOVED***
	return windows.SetEvent(e.h)
***REMOVED***

func (e *event) Wait() error ***REMOVED***
	s, err := windows.WaitForSingleObject(e.h, windows.INFINITE)
	switch s ***REMOVED***
	case windows.WAIT_OBJECT_0:
		break
	case windows.WAIT_FAILED:
		return err
	default:
		return errors.New("unexpected result from WaitForSingleObject")
	***REMOVED***
	return nil
***REMOVED***
