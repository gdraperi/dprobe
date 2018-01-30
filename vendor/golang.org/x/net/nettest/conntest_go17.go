// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.7

package nettest

import "testing"

func testConn(t *testing.T, mp MakePipe) ***REMOVED***
	// Use subtests on Go 1.7 and above since it is better organized.
	t.Run("BasicIO", func(t *testing.T) ***REMOVED*** timeoutWrapper(t, mp, testBasicIO) ***REMOVED***)
	t.Run("PingPong", func(t *testing.T) ***REMOVED*** timeoutWrapper(t, mp, testPingPong) ***REMOVED***)
	t.Run("RacyRead", func(t *testing.T) ***REMOVED*** timeoutWrapper(t, mp, testRacyRead) ***REMOVED***)
	t.Run("RacyWrite", func(t *testing.T) ***REMOVED*** timeoutWrapper(t, mp, testRacyWrite) ***REMOVED***)
	t.Run("ReadTimeout", func(t *testing.T) ***REMOVED*** timeoutWrapper(t, mp, testReadTimeout) ***REMOVED***)
	t.Run("WriteTimeout", func(t *testing.T) ***REMOVED*** timeoutWrapper(t, mp, testWriteTimeout) ***REMOVED***)
	t.Run("PastTimeout", func(t *testing.T) ***REMOVED*** timeoutWrapper(t, mp, testPastTimeout) ***REMOVED***)
	t.Run("PresentTimeout", func(t *testing.T) ***REMOVED*** timeoutWrapper(t, mp, testPresentTimeout) ***REMOVED***)
	t.Run("FutureTimeout", func(t *testing.T) ***REMOVED*** timeoutWrapper(t, mp, testFutureTimeout) ***REMOVED***)
	t.Run("CloseTimeout", func(t *testing.T) ***REMOVED*** timeoutWrapper(t, mp, testCloseTimeout) ***REMOVED***)
	t.Run("ConcurrentMethods", func(t *testing.T) ***REMOVED*** timeoutWrapper(t, mp, testConcurrentMethods) ***REMOVED***)
***REMOVED***
