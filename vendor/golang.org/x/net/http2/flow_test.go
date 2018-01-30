// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http2

import "testing"

func TestFlow(t *testing.T) ***REMOVED***
	var st flow
	var conn flow
	st.add(3)
	conn.add(2)

	if got, want := st.available(), int32(3); got != want ***REMOVED***
		t.Errorf("available = %d; want %d", got, want)
	***REMOVED***
	st.setConnFlow(&conn)
	if got, want := st.available(), int32(2); got != want ***REMOVED***
		t.Errorf("after parent setup, available = %d; want %d", got, want)
	***REMOVED***

	st.take(2)
	if got, want := conn.available(), int32(0); got != want ***REMOVED***
		t.Errorf("after taking 2, conn = %d; want %d", got, want)
	***REMOVED***
	if got, want := st.available(), int32(0); got != want ***REMOVED***
		t.Errorf("after taking 2, stream = %d; want %d", got, want)
	***REMOVED***
***REMOVED***

func TestFlowAdd(t *testing.T) ***REMOVED***
	var f flow
	if !f.add(1) ***REMOVED***
		t.Fatal("failed to add 1")
	***REMOVED***
	if !f.add(-1) ***REMOVED***
		t.Fatal("failed to add -1")
	***REMOVED***
	if got, want := f.available(), int32(0); got != want ***REMOVED***
		t.Fatalf("size = %d; want %d", got, want)
	***REMOVED***
	if !f.add(1<<31 - 1) ***REMOVED***
		t.Fatal("failed to add 2^31-1")
	***REMOVED***
	if got, want := f.available(), int32(1<<31-1); got != want ***REMOVED***
		t.Fatalf("size = %d; want %d", got, want)
	***REMOVED***
	if f.add(1) ***REMOVED***
		t.Fatal("adding 1 to max shouldn't be allowed")
	***REMOVED***

***REMOVED***
