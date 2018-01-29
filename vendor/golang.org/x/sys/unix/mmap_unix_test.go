// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package unix_test

import (
	"testing"

	"golang.org/x/sys/unix"
)

func TestMmap(t *testing.T) ***REMOVED***
	b, err := unix.Mmap(-1, 0, unix.Getpagesize(), unix.PROT_NONE, unix.MAP_ANON|unix.MAP_PRIVATE)
	if err != nil ***REMOVED***
		t.Fatalf("Mmap: %v", err)
	***REMOVED***
	if err := unix.Mprotect(b, unix.PROT_READ|unix.PROT_WRITE); err != nil ***REMOVED***
		t.Fatalf("Mprotect: %v", err)
	***REMOVED***

	b[0] = 42

	if err := unix.Msync(b, unix.MS_SYNC); err != nil ***REMOVED***
		t.Fatalf("Msync: %v", err)
	***REMOVED***
	if err := unix.Madvise(b, unix.MADV_DONTNEED); err != nil ***REMOVED***
		t.Fatalf("Madvise: %v", err)
	***REMOVED***
	if err := unix.Munmap(b); err != nil ***REMOVED***
		t.Fatalf("Munmap: %v", err)
	***REMOVED***
***REMOVED***
