// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.7

package unix_test

import (
	"fmt"
	"testing"

	"golang.org/x/sys/unix"
)

func TestDevices(t *testing.T) ***REMOVED***
	testCases := []struct ***REMOVED***
		path  string
		major uint32
		minor uint32
	***REMOVED******REMOVED***
		// Well-known major/minor numbers on OpenSolaris according to
		// /etc/name_to_major
		***REMOVED***"/dev/zero", 134, 12***REMOVED***,
		***REMOVED***"/dev/null", 134, 2***REMOVED***,
		***REMOVED***"/dev/ptyp0", 172, 0***REMOVED***,
		***REMOVED***"/dev/ttyp0", 175, 0***REMOVED***,
		***REMOVED***"/dev/ttyp1", 175, 1***REMOVED***,
	***REMOVED***
	for _, tc := range testCases ***REMOVED***
		t.Run(fmt.Sprintf("%s %v:%v", tc.path, tc.major, tc.minor), func(t *testing.T) ***REMOVED***
			var stat unix.Stat_t
			err := unix.Stat(tc.path, &stat)
			if err != nil ***REMOVED***
				t.Errorf("failed to stat device: %v", err)
				return
			***REMOVED***

			dev := uint64(stat.Rdev)
			if unix.Major(dev) != tc.major ***REMOVED***
				t.Errorf("for %s Major(%#x) == %d, want %d", tc.path, dev, unix.Major(dev), tc.major)
			***REMOVED***
			if unix.Minor(dev) != tc.minor ***REMOVED***
				t.Errorf("for %s Minor(%#x) == %d, want %d", tc.path, dev, unix.Minor(dev), tc.minor)
			***REMOVED***
			if unix.Mkdev(tc.major, tc.minor) != dev ***REMOVED***
				t.Errorf("for %s Mkdev(%d, %d) == %#x, want %#x", tc.path, tc.major, tc.minor, unix.Mkdev(tc.major, tc.minor), dev)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***
