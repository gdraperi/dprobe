// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package eventlog_test

import (
	"testing"

	"golang.org/x/sys/windows/svc/eventlog"
)

func TestLog(t *testing.T) ***REMOVED***
	if testing.Short() ***REMOVED***
		t.Skip("skipping test in short mode - it modifies system logs")
	***REMOVED***

	const name = "mylog"
	const supports = eventlog.Error | eventlog.Warning | eventlog.Info
	err := eventlog.InstallAsEventCreate(name, supports)
	if err != nil ***REMOVED***
		t.Fatalf("Install failed: %s", err)
	***REMOVED***
	defer func() ***REMOVED***
		err = eventlog.Remove(name)
		if err != nil ***REMOVED***
			t.Fatalf("Remove failed: %s", err)
		***REMOVED***
	***REMOVED***()

	l, err := eventlog.Open(name)
	if err != nil ***REMOVED***
		t.Fatalf("Open failed: %s", err)
	***REMOVED***
	defer l.Close()

	err = l.Info(1, "info")
	if err != nil ***REMOVED***
		t.Fatalf("Info failed: %s", err)
	***REMOVED***
	err = l.Warning(2, "warning")
	if err != nil ***REMOVED***
		t.Fatalf("Warning failed: %s", err)
	***REMOVED***
	err = l.Error(3, "error")
	if err != nil ***REMOVED***
		t.Fatalf("Error failed: %s", err)
	***REMOVED***
***REMOVED***
