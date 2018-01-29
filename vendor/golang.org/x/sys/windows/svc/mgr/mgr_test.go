// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package mgr_test

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"testing"
	"time"

	"golang.org/x/sys/windows/svc/mgr"
)

func TestOpenLanManServer(t *testing.T) ***REMOVED***
	m, err := mgr.Connect()
	if err != nil ***REMOVED***
		if errno, ok := err.(syscall.Errno); ok && errno == syscall.ERROR_ACCESS_DENIED ***REMOVED***
			t.Skip("Skipping test: we don't have rights to manage services.")
		***REMOVED***
		t.Fatalf("SCM connection failed: %s", err)
	***REMOVED***
	defer m.Disconnect()

	s, err := m.OpenService("LanmanServer")
	if err != nil ***REMOVED***
		t.Fatalf("OpenService(lanmanserver) failed: %s", err)
	***REMOVED***
	defer s.Close()

	_, err = s.Config()
	if err != nil ***REMOVED***
		t.Fatalf("Config failed: %s", err)
	***REMOVED***
***REMOVED***

func install(t *testing.T, m *mgr.Mgr, name, exepath string, c mgr.Config) ***REMOVED***
	// Sometimes it takes a while for the service to get
	// removed after previous test run.
	for i := 0; ; i++ ***REMOVED***
		s, err := m.OpenService(name)
		if err != nil ***REMOVED***
			break
		***REMOVED***
		s.Close()

		if i > 10 ***REMOVED***
			t.Fatalf("service %s already exists", name)
		***REMOVED***
		time.Sleep(300 * time.Millisecond)
	***REMOVED***

	s, err := m.CreateService(name, exepath, c)
	if err != nil ***REMOVED***
		t.Fatalf("CreateService(%s) failed: %v", name, err)
	***REMOVED***
	defer s.Close()
***REMOVED***

func depString(d []string) string ***REMOVED***
	if len(d) == 0 ***REMOVED***
		return ""
	***REMOVED***
	for i := range d ***REMOVED***
		d[i] = strings.ToLower(d[i])
	***REMOVED***
	ss := sort.StringSlice(d)
	ss.Sort()
	return strings.Join([]string(ss), " ")
***REMOVED***

func testConfig(t *testing.T, s *mgr.Service, should mgr.Config) mgr.Config ***REMOVED***
	is, err := s.Config()
	if err != nil ***REMOVED***
		t.Fatalf("Config failed: %s", err)
	***REMOVED***
	if should.DisplayName != is.DisplayName ***REMOVED***
		t.Fatalf("config mismatch: DisplayName is %q, but should have %q", is.DisplayName, should.DisplayName)
	***REMOVED***
	if should.StartType != is.StartType ***REMOVED***
		t.Fatalf("config mismatch: StartType is %v, but should have %v", is.StartType, should.StartType)
	***REMOVED***
	if should.Description != is.Description ***REMOVED***
		t.Fatalf("config mismatch: Description is %q, but should have %q", is.Description, should.Description)
	***REMOVED***
	if depString(should.Dependencies) != depString(is.Dependencies) ***REMOVED***
		t.Fatalf("config mismatch: Dependencies is %v, but should have %v", is.Dependencies, should.Dependencies)
	***REMOVED***
	return is
***REMOVED***

func remove(t *testing.T, s *mgr.Service) ***REMOVED***
	err := s.Delete()
	if err != nil ***REMOVED***
		t.Fatalf("Delete failed: %s", err)
	***REMOVED***
***REMOVED***

func TestMyService(t *testing.T) ***REMOVED***
	if testing.Short() ***REMOVED***
		t.Skip("skipping test in short mode - it modifies system services")
	***REMOVED***

	const name = "myservice"

	m, err := mgr.Connect()
	if err != nil ***REMOVED***
		if errno, ok := err.(syscall.Errno); ok && errno == syscall.ERROR_ACCESS_DENIED ***REMOVED***
			t.Skip("Skipping test: we don't have rights to manage services.")
		***REMOVED***
		t.Fatalf("SCM connection failed: %s", err)
	***REMOVED***
	defer m.Disconnect()

	c := mgr.Config***REMOVED***
		StartType:    mgr.StartDisabled,
		DisplayName:  "my service",
		Description:  "my service is just a test",
		Dependencies: []string***REMOVED***"LanmanServer", "W32Time"***REMOVED***,
	***REMOVED***

	exename := os.Args[0]
	exepath, err := filepath.Abs(exename)
	if err != nil ***REMOVED***
		t.Fatalf("filepath.Abs(%s) failed: %s", exename, err)
	***REMOVED***

	install(t, m, name, exepath, c)

	s, err := m.OpenService(name)
	if err != nil ***REMOVED***
		t.Fatalf("service %s is not installed", name)
	***REMOVED***
	defer s.Close()

	c.BinaryPathName = exepath
	c = testConfig(t, s, c)

	c.StartType = mgr.StartManual
	err = s.UpdateConfig(c)
	if err != nil ***REMOVED***
		t.Fatalf("UpdateConfig failed: %v", err)
	***REMOVED***

	testConfig(t, s, c)

	svcnames, err := m.ListServices()
	if err != nil ***REMOVED***
		t.Fatalf("ListServices failed: %v", err)
	***REMOVED***
	var myserviceIsInstalled bool
	for _, sn := range svcnames ***REMOVED***
		if sn == name ***REMOVED***
			myserviceIsInstalled = true
			break
		***REMOVED***
	***REMOVED***
	if !myserviceIsInstalled ***REMOVED***
		t.Errorf("ListServices failed to find %q service", name)
	***REMOVED***

	remove(t, s)
***REMOVED***
