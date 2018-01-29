// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd openbsd

package unix_test

import (
	"os/exec"
	"runtime"
	"testing"
	"time"

	"golang.org/x/sys/unix"
)

const MNT_WAIT = 1
const MNT_NOWAIT = 2

func TestGetfsstat(t *testing.T) ***REMOVED***
	const flags = MNT_NOWAIT // see golang.org/issue/16937
	n, err := unix.Getfsstat(nil, flags)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	data := make([]unix.Statfs_t, n)
	n2, err := unix.Getfsstat(data, flags)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if n != n2 ***REMOVED***
		t.Errorf("Getfsstat(nil) = %d, but subsequent Getfsstat(slice) = %d", n, n2)
	***REMOVED***
	for i, stat := range data ***REMOVED***
		if stat == (unix.Statfs_t***REMOVED******REMOVED***) ***REMOVED***
			t.Errorf("index %v is an empty Statfs_t struct", i)
		***REMOVED***
	***REMOVED***
	if t.Failed() ***REMOVED***
		for i, stat := range data[:n2] ***REMOVED***
			t.Logf("data[%v] = %+v", i, stat)
		***REMOVED***
		mount, err := exec.Command("mount").CombinedOutput()
		if err != nil ***REMOVED***
			t.Logf("mount: %v\n%s", err, mount)
		***REMOVED*** else ***REMOVED***
			t.Logf("mount: %s", mount)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestSelect(t *testing.T) ***REMOVED***
	err := unix.Select(0, nil, nil, nil, &unix.Timeval***REMOVED***Sec: 0, Usec: 0***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("Select: %v", err)
	***REMOVED***

	dur := 250 * time.Millisecond
	tv := unix.NsecToTimeval(int64(dur))
	start := time.Now()
	err = unix.Select(0, nil, nil, nil, &tv)
	took := time.Since(start)
	if err != nil ***REMOVED***
		t.Fatalf("Select: %v", err)
	***REMOVED***

	// On some BSDs the actual timeout might also be slightly less than the requested.
	// Add an acceptable margin to avoid flaky tests.
	if took < dur*2/3 ***REMOVED***
		t.Errorf("Select: timeout should have been at least %v, got %v", dur, took)
	***REMOVED***
***REMOVED***

func TestSysctlRaw(t *testing.T) ***REMOVED***
	if runtime.GOOS == "openbsd" ***REMOVED***
		t.Skip("kern.proc.pid does not exist on OpenBSD")
	***REMOVED***

	_, err := unix.SysctlRaw("kern.proc.pid", unix.Getpid())
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func TestSysctlUint32(t *testing.T) ***REMOVED***
	maxproc, err := unix.SysctlUint32("kern.maxproc")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	t.Logf("kern.maxproc: %v", maxproc)
***REMOVED***
