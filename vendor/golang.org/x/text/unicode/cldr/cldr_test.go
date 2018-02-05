// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cldr

import "testing"

func TestParseDraft(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		in    string
		draft Draft
		err   bool
	***REMOVED******REMOVED***
		***REMOVED***"unconfirmed", Unconfirmed, false***REMOVED***,
		***REMOVED***"provisional", Provisional, false***REMOVED***,
		***REMOVED***"contributed", Contributed, false***REMOVED***,
		***REMOVED***"approved", Approved, false***REMOVED***,
		***REMOVED***"", Approved, false***REMOVED***,
		***REMOVED***"foo", Approved, true***REMOVED***,
	***REMOVED***
	for _, tt := range tests ***REMOVED***
		if d, err := ParseDraft(tt.in); d != tt.draft || (err != nil) != tt.err ***REMOVED***
			t.Errorf("%q: was %v, %v; want %v, %v", tt.in, d, err != nil, tt.draft, tt.err)
		***REMOVED***
	***REMOVED***
***REMOVED***
