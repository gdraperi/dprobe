// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package main

import (
	"fmt"
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

func startService(name string) error ***REMOVED***
	m, err := mgr.Connect()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer m.Disconnect()
	s, err := m.OpenService(name)
	if err != nil ***REMOVED***
		return fmt.Errorf("could not access service: %v", err)
	***REMOVED***
	defer s.Close()
	err = s.Start("is", "manual-started")
	if err != nil ***REMOVED***
		return fmt.Errorf("could not start service: %v", err)
	***REMOVED***
	return nil
***REMOVED***

func controlService(name string, c svc.Cmd, to svc.State) error ***REMOVED***
	m, err := mgr.Connect()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer m.Disconnect()
	s, err := m.OpenService(name)
	if err != nil ***REMOVED***
		return fmt.Errorf("could not access service: %v", err)
	***REMOVED***
	defer s.Close()
	status, err := s.Control(c)
	if err != nil ***REMOVED***
		return fmt.Errorf("could not send control=%d: %v", c, err)
	***REMOVED***
	timeout := time.Now().Add(10 * time.Second)
	for status.State != to ***REMOVED***
		if timeout.Before(time.Now()) ***REMOVED***
			return fmt.Errorf("timeout waiting for service to go to state=%d", to)
		***REMOVED***
		time.Sleep(300 * time.Millisecond)
		status, err = s.Query()
		if err != nil ***REMOVED***
			return fmt.Errorf("could not retrieve service status: %v", err)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***
