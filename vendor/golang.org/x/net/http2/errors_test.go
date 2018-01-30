// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http2

import "testing"

func TestErrCodeString(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		err  ErrCode
		want string
	***REMOVED******REMOVED***
		***REMOVED***ErrCodeProtocol, "PROTOCOL_ERROR"***REMOVED***,
		***REMOVED***0xd, "HTTP_1_1_REQUIRED"***REMOVED***,
		***REMOVED***0xf, "unknown error code 0xf"***REMOVED***,
	***REMOVED***
	for i, tt := range tests ***REMOVED***
		got := tt.err.String()
		if got != tt.want ***REMOVED***
			t.Errorf("%d. Error = %q; want %q", i, got, tt.want)
		***REMOVED***
	***REMOVED***
***REMOVED***
