// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package route

import (
	"testing"
	"unsafe"
)

func TestFetchAndParseRIBOnFreeBSD(t *testing.T) ***REMOVED***
	for _, typ := range []RIBType***REMOVED***sysNET_RT_IFMALIST***REMOVED*** ***REMOVED***
		var lastErr error
		var ms []Message
		for _, af := range []int***REMOVED***sysAF_UNSPEC, sysAF_INET, sysAF_INET6***REMOVED*** ***REMOVED***
			rs, err := fetchAndParseRIB(af, typ)
			if err != nil ***REMOVED***
				lastErr = err
				continue
			***REMOVED***
			ms = append(ms, rs...)
		***REMOVED***
		if len(ms) == 0 && lastErr != nil ***REMOVED***
			t.Error(typ, lastErr)
			continue
		***REMOVED***
		ss, err := msgs(ms).validate()
		if err != nil ***REMOVED***
			t.Error(typ, err)
			continue
		***REMOVED***
		for _, s := range ss ***REMOVED***
			t.Log(s)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestFetchAndParseRIBOnFreeBSD10AndAbove(t *testing.T) ***REMOVED***
	if _, err := FetchRIB(sysAF_UNSPEC, sysNET_RT_IFLISTL, 0); err != nil ***REMOVED***
		t.Skip("NET_RT_IFLISTL not supported")
	***REMOVED***
	var p uintptr
	if kernelAlign != int(unsafe.Sizeof(p)) ***REMOVED***
		t.Skip("NET_RT_IFLIST vs. NET_RT_IFLISTL doesn't work for 386 emulation on amd64")
	***REMOVED***

	var tests = [2]struct ***REMOVED***
		typ  RIBType
		b    []byte
		msgs []Message
		ss   []string
	***REMOVED******REMOVED***
		***REMOVED***typ: sysNET_RT_IFLIST***REMOVED***,
		***REMOVED***typ: sysNET_RT_IFLISTL***REMOVED***,
	***REMOVED***
	for i := range tests ***REMOVED***
		var lastErr error
		for _, af := range []int***REMOVED***sysAF_UNSPEC, sysAF_INET, sysAF_INET6***REMOVED*** ***REMOVED***
			rs, err := fetchAndParseRIB(af, tests[i].typ)
			if err != nil ***REMOVED***
				lastErr = err
				continue
			***REMOVED***
			tests[i].msgs = append(tests[i].msgs, rs...)
		***REMOVED***
		if len(tests[i].msgs) == 0 && lastErr != nil ***REMOVED***
			t.Error(tests[i].typ, lastErr)
			continue
		***REMOVED***
		tests[i].ss, lastErr = msgs(tests[i].msgs).validate()
		if lastErr != nil ***REMOVED***
			t.Error(tests[i].typ, lastErr)
			continue
		***REMOVED***
		for _, s := range tests[i].ss ***REMOVED***
			t.Log(s)
		***REMOVED***
	***REMOVED***
	for i := len(tests) - 1; i > 0; i-- ***REMOVED***
		if len(tests[i].ss) != len(tests[i-1].ss) ***REMOVED***
			t.Errorf("got %v; want %v", tests[i].ss, tests[i-1].ss)
			continue
		***REMOVED***
		for j, s1 := range tests[i].ss ***REMOVED***
			s0 := tests[i-1].ss[j]
			if s1 != s0 ***REMOVED***
				t.Errorf("got %s; want %s", s1, s0)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
