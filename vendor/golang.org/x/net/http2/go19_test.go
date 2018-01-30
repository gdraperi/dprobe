// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.9

package http2

import (
	"context"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestServerGracefulShutdown(t *testing.T) ***REMOVED***
	var st *serverTester
	handlerDone := make(chan struct***REMOVED******REMOVED***)
	st = newServerTester(t, func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		defer close(handlerDone)
		go st.ts.Config.Shutdown(context.Background())

		ga := st.wantGoAway()
		if ga.ErrCode != ErrCodeNo ***REMOVED***
			t.Errorf("GOAWAY error = %v; want ErrCodeNo", ga.ErrCode)
		***REMOVED***
		if ga.LastStreamID != 1 ***REMOVED***
			t.Errorf("GOAWAY LastStreamID = %v; want 1", ga.LastStreamID)
		***REMOVED***

		w.Header().Set("x-foo", "bar")
	***REMOVED***)
	defer st.Close()

	st.greet()
	st.bodylessReq1()

	select ***REMOVED***
	case <-handlerDone:
	case <-time.After(5 * time.Second):
		t.Fatalf("server did not shutdown?")
	***REMOVED***
	hf := st.wantHeaders()
	goth := st.decodeHeader(hf.HeaderBlockFragment())
	wanth := [][2]string***REMOVED***
		***REMOVED***":status", "200"***REMOVED***,
		***REMOVED***"x-foo", "bar"***REMOVED***,
		***REMOVED***"content-length", "0"***REMOVED***,
	***REMOVED***
	if !reflect.DeepEqual(goth, wanth) ***REMOVED***
		t.Errorf("Got headers %v; want %v", goth, wanth)
	***REMOVED***

	n, err := st.cc.Read([]byte***REMOVED***0***REMOVED***)
	if n != 0 || err == nil ***REMOVED***
		t.Errorf("Read = %v, %v; want 0, non-nil", n, err)
	***REMOVED***
***REMOVED***
