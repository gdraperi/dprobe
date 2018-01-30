// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http2

import (
	"fmt"
	"strings"
	"testing"
)

func TestGoroutineLock(t *testing.T) ***REMOVED***
	oldDebug := DebugGoroutines
	DebugGoroutines = true
	defer func() ***REMOVED*** DebugGoroutines = oldDebug ***REMOVED***()

	g := newGoroutineLock()
	g.check()

	sawPanic := make(chan interface***REMOVED******REMOVED***)
	go func() ***REMOVED***
		defer func() ***REMOVED*** sawPanic <- recover() ***REMOVED***()
		g.check() // should panic
	***REMOVED***()
	e := <-sawPanic
	if e == nil ***REMOVED***
		t.Fatal("did not see panic from check in other goroutine")
	***REMOVED***
	if !strings.Contains(fmt.Sprint(e), "wrong goroutine") ***REMOVED***
		t.Errorf("expected on see panic about running on the wrong goroutine; got %v", e)
	***REMOVED***
***REMOVED***
