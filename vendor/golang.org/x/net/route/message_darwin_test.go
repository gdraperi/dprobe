// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package route

import "testing"

func TestFetchAndParseRIBOnDarwin(t *testing.T) ***REMOVED***
	for _, typ := range []RIBType***REMOVED***sysNET_RT_FLAGS, sysNET_RT_DUMP2, sysNET_RT_IFLIST2***REMOVED*** ***REMOVED***
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
