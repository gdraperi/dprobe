// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build plan9

package plan9_test

import (
	"testing"

	"golang.org/x/sys/plan9"
)

func testSetGetenv(t *testing.T, key, value string) ***REMOVED***
	err := plan9.Setenv(key, value)
	if err != nil ***REMOVED***
		t.Fatalf("Setenv failed to set %q: %v", value, err)
	***REMOVED***
	newvalue, found := plan9.Getenv(key)
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
