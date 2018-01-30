// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build solaris

// Package lif provides basic functions for the manipulation of
// logical network interfaces and interface addresses on Solaris.
//
// The package supports Solaris 11 or above.
package lif

import "syscall"

type endpoint struct ***REMOVED***
	af int
	s  uintptr
***REMOVED***

func (ep *endpoint) close() error ***REMOVED***
	return syscall.Close(int(ep.s))
***REMOVED***

func newEndpoints(af int) ([]endpoint, error) ***REMOVED***
	var lastErr error
	var eps []endpoint
	afs := []int***REMOVED***sysAF_INET, sysAF_INET6***REMOVED***
	if af != sysAF_UNSPEC ***REMOVED***
		afs = []int***REMOVED***af***REMOVED***
	***REMOVED***
	for _, af := range afs ***REMOVED***
		s, err := syscall.Socket(af, sysSOCK_DGRAM, 0)
		if err != nil ***REMOVED***
			lastErr = err
			continue
		***REMOVED***
		eps = append(eps, endpoint***REMOVED***af: af, s: uintptr(s)***REMOVED***)
	***REMOVED***
	if len(eps) == 0 ***REMOVED***
		return nil, lastErr
	***REMOVED***
	return eps, nil
***REMOVED***
