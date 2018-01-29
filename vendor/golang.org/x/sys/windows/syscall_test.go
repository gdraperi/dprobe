// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package windows_test

import (
	"syscall"
	"testing"

	"golang.org/x/sys/windows"
)

func testSetGetenv(t *testing.T, key, value string) ***REMOVED***
	err := windows.Setenv(key, value)
	if err != nil ***REMOVED***
		t.Fatalf("Setenv failed to set %q: %v", value, err)
	***REMOVED***
	newvalue, found := windows.Getenv(key)
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

func TestGetProcAddressByOrdinal(t *testing.T) ***REMOVED***
	// Attempt calling shlwapi.dll:IsOS, resolving it by ordinal, as
	// suggested in
	// https://msdn.microsoft.com/en-us/library/windows/desktop/bb773795.aspx
	h, err := windows.LoadLibrary("shlwapi.dll")
	if err != nil ***REMOVED***
		t.Fatalf("Failed to load shlwapi.dll: %s", err)
	***REMOVED***
	procIsOS, err := windows.GetProcAddressByOrdinal(h, 437)
	if err != nil ***REMOVED***
		t.Fatalf("Could not find shlwapi.dll:IsOS by ordinal: %s", err)
	***REMOVED***
	const OS_NT = 1
	r, _, _ := syscall.Syscall(procIsOS, 1, OS_NT, 0, 0)
	if r == 0 ***REMOVED***
		t.Error("shlwapi.dll:IsOS(OS_NT) returned 0, expected non-zero value")
	***REMOVED***
***REMOVED***
