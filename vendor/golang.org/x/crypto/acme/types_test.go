// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package acme

import (
	"errors"
	"net/http"
	"testing"
	"time"
)

func TestRateLimit(t *testing.T) ***REMOVED***
	now := time.Date(2017, 04, 27, 10, 0, 0, 0, time.UTC)
	f := timeNow
	defer func() ***REMOVED*** timeNow = f ***REMOVED***()
	timeNow = func() time.Time ***REMOVED*** return now ***REMOVED***

	h120, hTime := http.Header***REMOVED******REMOVED***, http.Header***REMOVED******REMOVED***
	h120.Set("Retry-After", "120")
	hTime.Set("Retry-After", "Tue Apr 27 11:00:00 2017")

	err1 := &Error***REMOVED***
		ProblemType: "urn:ietf:params:acme:error:nolimit",
		Header:      h120,
	***REMOVED***
	err2 := &Error***REMOVED***
		ProblemType: "urn:ietf:params:acme:error:rateLimited",
		Header:      h120,
	***REMOVED***
	err3 := &Error***REMOVED***
		ProblemType: "urn:ietf:params:acme:error:rateLimited",
		Header:      nil,
	***REMOVED***
	err4 := &Error***REMOVED***
		ProblemType: "urn:ietf:params:acme:error:rateLimited",
		Header:      hTime,
	***REMOVED***

	tt := []struct ***REMOVED***
		err error
		res time.Duration
		ok  bool
	***REMOVED******REMOVED***
		***REMOVED***nil, 0, false***REMOVED***,
		***REMOVED***errors.New("dummy"), 0, false***REMOVED***,
		***REMOVED***err1, 0, false***REMOVED***,
		***REMOVED***err2, 2 * time.Minute, true***REMOVED***,
		***REMOVED***err3, 0, true***REMOVED***,
		***REMOVED***err4, time.Hour, true***REMOVED***,
	***REMOVED***
	for i, test := range tt ***REMOVED***
		res, ok := RateLimit(test.err)
		if ok != test.ok ***REMOVED***
			t.Errorf("%d: RateLimit(%+v): ok = %v; want %v", i, test.err, ok, test.ok)
			continue
		***REMOVED***
		if res != test.res ***REMOVED***
			t.Errorf("%d: RateLimit(%+v) = %v; want %v", i, test.err, res, test.res)
		***REMOVED***
	***REMOVED***
***REMOVED***
