// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package svc_test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

func getState(t *testing.T, s *mgr.Service) svc.State ***REMOVED***
	status, err := s.Query()
	if err != nil ***REMOVED***
		t.Fatalf("Query(%s) failed: %s", s.Name, err)
	***REMOVED***
	return status.State
***REMOVED***

func testState(t *testing.T, s *mgr.Service, want svc.State) ***REMOVED***
	have := getState(t, s)
	if have != want ***REMOVED***
		t.Fatalf("%s state is=%d want=%d", s.Name, have, want)
	***REMOVED***
***REMOVED***

func waitState(t *testing.T, s *mgr.Service, want svc.State) ***REMOVED***
	for i := 0; ; i++ ***REMOVED***
		have := getState(t, s)
		if have == want ***REMOVED***
			return
		***REMOVED***
		if i > 10 ***REMOVED***
			t.Fatalf("%s state is=%d, waiting timeout", s.Name, have)
		***REMOVED***
		time.Sleep(300 * time.Millisecond)
	***REMOVED***
***REMOVED***

func TestExample(t *testing.T) ***REMOVED***
	if testing.Short() ***REMOVED***
		t.Skip("skipping test in short mode - it modifies system services")
	***REMOVED***

	const name = "myservice"

	m, err := mgr.Connect()
	if err != nil ***REMOVED***
		t.Fatalf("SCM connection failed: %s", err)
	***REMOVED***
	defer m.Disconnect()

	dir, err := ioutil.TempDir("", "svc")
	if err != nil ***REMOVED***
		t.Fatalf("failed to create temp directory: %v", err)
	***REMOVED***
	defer os.RemoveAll(dir)

	exepath := filepath.Join(dir, "a.exe")
	o, err := exec.Command("go", "build", "-o", exepath, "golang.org/x/sys/windows/svc/example").CombinedOutput()
	if err != nil ***REMOVED***
		t.Fatalf("failed to build service program: %v\n%v", err, string(o))
	***REMOVED***

	s, err := m.OpenService(name)
	if err == nil ***REMOVED***
		err = s.Delete()
		if err != nil ***REMOVED***
			s.Close()
			t.Fatalf("Delete failed: %s", err)
		***REMOVED***
		s.Close()
	***REMOVED***
	s, err = m.CreateService(name, exepath, mgr.Config***REMOVED***DisplayName: "my service"***REMOVED***, "is", "auto-started")
	if err != nil ***REMOVED***
		t.Fatalf("CreateService(%s) failed: %v", name, err)
	***REMOVED***
	defer s.Close()

	testState(t, s, svc.Stopped)
	err = s.Start("is", "manual-started")
	if err != nil ***REMOVED***
		t.Fatalf("Start(%s) failed: %s", s.Name, err)
	***REMOVED***
	waitState(t, s, svc.Running)
	time.Sleep(1 * time.Second)

	// testing deadlock from issues 4.
	_, err = s.Control(svc.Interrogate)
	if err != nil ***REMOVED***
		t.Fatalf("Control(%s) failed: %s", s.Name, err)
	***REMOVED***
	_, err = s.Control(svc.Interrogate)
	if err != nil ***REMOVED***
		t.Fatalf("Control(%s) failed: %s", s.Name, err)
	***REMOVED***
	time.Sleep(1 * time.Second)

	_, err = s.Control(svc.Stop)
	if err != nil ***REMOVED***
		t.Fatalf("Control(%s) failed: %s", s.Name, err)
	***REMOVED***
	waitState(t, s, svc.Stopped)

	err = s.Delete()
	if err != nil ***REMOVED***
		t.Fatalf("Delete failed: %s", err)
	***REMOVED***
***REMOVED***
