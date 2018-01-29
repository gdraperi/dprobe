// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build solaris

package unix_test

import (
	"os/exec"
	"testing"
	"time"

	"golang.org/x/sys/unix"
)

func TestSelect(t *testing.T) ***REMOVED***
	err := unix.Select(0, nil, nil, nil, &unix.Timeval***REMOVED***Sec: 0, Usec: 0***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("Select: %v", err)
	***REMOVED***

	dur := 150 * time.Millisecond
	tv := unix.NsecToTimeval(int64(dur))
	start := time.Now()
	err = unix.Select(0, nil, nil, nil, &tv)
	took := time.Since(start)
	if err != nil ***REMOVED***
		t.Fatalf("Select: %v", err)
	***REMOVED***

	if took < dur ***REMOVED***
		t.Errorf("Select: timeout should have been at least %v, got %v", dur, took)
	***REMOVED***
***REMOVED***

func TestStatvfs(t *testing.T) ***REMOVED***
	if err := unix.Statvfs("", nil); err == nil ***REMOVED***
		t.Fatal(`Statvfs("") expected failure`)
	***REMOVED***

	statvfs := unix.Statvfs_t***REMOVED******REMOVED***
	if err := unix.Statvfs("/", &statvfs); err != nil ***REMOVED***
		t.Errorf(`Statvfs("/") failed: %v`, err)
	***REMOVED***

	if t.Failed() ***REMOVED***
		mount, err := exec.Command("mount").CombinedOutput()
		if err != nil ***REMOVED***
			t.Logf("mount: %v\n%s", err, mount)
		***REMOVED*** else ***REMOVED***
			t.Logf("mount: %s", mount)
		***REMOVED***
	***REMOVED***
***REMOVED***
