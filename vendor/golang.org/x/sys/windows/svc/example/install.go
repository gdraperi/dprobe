// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"
)

func exePath() (string, error) ***REMOVED***
	prog := os.Args[0]
	p, err := filepath.Abs(prog)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	fi, err := os.Stat(p)
	if err == nil ***REMOVED***
		if !fi.Mode().IsDir() ***REMOVED***
			return p, nil
		***REMOVED***
		err = fmt.Errorf("%s is directory", p)
	***REMOVED***
	if filepath.Ext(p) == "" ***REMOVED***
		p += ".exe"
		fi, err := os.Stat(p)
		if err == nil ***REMOVED***
			if !fi.Mode().IsDir() ***REMOVED***
				return p, nil
			***REMOVED***
			err = fmt.Errorf("%s is directory", p)
		***REMOVED***
	***REMOVED***
	return "", err
***REMOVED***

func installService(name, desc string) error ***REMOVED***
	exepath, err := exePath()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	m, err := mgr.Connect()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer m.Disconnect()
	s, err := m.OpenService(name)
	if err == nil ***REMOVED***
		s.Close()
		return fmt.Errorf("service %s already exists", name)
	***REMOVED***
	s, err = m.CreateService(name, exepath, mgr.Config***REMOVED***DisplayName: desc***REMOVED***, "is", "auto-started")
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer s.Close()
	err = eventlog.InstallAsEventCreate(name, eventlog.Error|eventlog.Warning|eventlog.Info)
	if err != nil ***REMOVED***
		s.Delete()
		return fmt.Errorf("SetupEventLogSource() failed: %s", err)
	***REMOVED***
	return nil
***REMOVED***

func removeService(name string) error ***REMOVED***
	m, err := mgr.Connect()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer m.Disconnect()
	s, err := m.OpenService(name)
	if err != nil ***REMOVED***
		return fmt.Errorf("service %s is not installed", name)
	***REMOVED***
	defer s.Close()
	err = s.Delete()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	err = eventlog.Remove(name)
	if err != nil ***REMOVED***
		return fmt.Errorf("RemoveEventLogSource() failed: %s", err)
	***REMOVED***
	return nil
***REMOVED***
