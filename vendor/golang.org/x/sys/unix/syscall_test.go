// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package unix_test

import (
	"fmt"
	"testing"

	"golang.org/x/sys/unix"
)

func testSetGetenv(t *testing.T, key, value string) ***REMOVED***
	err := unix.Setenv(key, value)
	if err != nil ***REMOVED***
		t.Fatalf("Setenv failed to set %q: %v", value, err)
	***REMOVED***
	newvalue, found := unix.Getenv(key)
	if !found ***REMOVED***
		t.Fatalf("Getenv failed to find %v variable (want value %q)", key, value)
	***REMOVED***
	if newvalue != value ***REMOVED***
		t.Fatalf("Getenv(%v) = %q; want %q", key, newvalue, value)
	***REMOVED***
***REMOVED***

func TestEnv(t *testing.T) ***REMOVED***
	testSetGetenv(t, "TESTENV", "AVALUE")
	// make sure TESTENV gets set to "", not deleted
	testSetGetenv(t, "TESTENV", "")
***REMOVED***

func TestItoa(t *testing.T) ***REMOVED***
	// Make most negative integer: 0x8000...
	i := 1
	for i<<1 != 0 ***REMOVED***
		i <<= 1
	***REMOVED***
	if i >= 0 ***REMOVED***
		t.Fatal("bad math")
	***REMOVED***
	s := unix.Itoa(i)
	f := fmt.Sprint(i)
	if s != f ***REMOVED***
		t.Fatalf("itoa(%d) = %s, want %s", i, s, f)
	***REMOVED***
***REMOVED***

func TestUname(t *testing.T) ***REMOVED***
	var utsname unix.Utsname
	err := unix.Uname(&utsname)
	if err != nil ***REMOVED***
		t.Fatalf("Uname: %v", err)
	***REMOVED***

	t.Logf("OS: %s/%s %s", utsname.Sysname[:], utsname.Machine[:], utsname.Release[:])
***REMOVED***
