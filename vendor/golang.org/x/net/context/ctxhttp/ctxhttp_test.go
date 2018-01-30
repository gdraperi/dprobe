// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !plan9

package ctxhttp

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"golang.org/x/net/context"
)

const (
	requestDuration = 100 * time.Millisecond
	requestBody     = "ok"
)

func okHandler(w http.ResponseWriter, r *http.Request) ***REMOVED***
	time.Sleep(requestDuration)
	io.WriteString(w, requestBody)
***REMOVED***

func TestNoTimeout(t *testing.T) ***REMOVED***
	ts := httptest.NewServer(http.HandlerFunc(okHandler))
	defer ts.Close()

	ctx := context.Background()
	res, err := Get(ctx, nil, ts.URL)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer res.Body.Close()
	slurp, err := ioutil.ReadAll(res.Body)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if string(slurp) != requestBody ***REMOVED***
		t.Errorf("body = %q; want %q", slurp, requestBody)
	***REMOVED***
***REMOVED***

func TestCancelBeforeHeaders(t *testing.T) ***REMOVED***
	ctx, cancel := context.WithCancel(context.Background())

	blockServer := make(chan struct***REMOVED******REMOVED***)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		cancel()
		<-blockServer
		io.WriteString(w, requestBody)
	***REMOVED***))
	defer ts.Close()
	defer close(blockServer)

	res, err := Get(ctx, nil, ts.URL)
	if err == nil ***REMOVED***
		res.Body.Close()
		t.Fatal("Get returned unexpected nil error")
	***REMOVED***
	if err != context.Canceled ***REMOVED***
		t.Errorf("err = %v; want %v", err, context.Canceled)
	***REMOVED***
***REMOVED***

func TestCancelAfterHangingRequest(t *testing.T) ***REMOVED***
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		w.WriteHeader(http.StatusOK)
		w.(http.Flusher).Flush()
		<-w.(http.CloseNotifier).CloseNotify()
	***REMOVED***))
	defer ts.Close()

	ctx, cancel := context.WithCancel(context.Background())
	resp, err := Get(ctx, nil, ts.URL)
	if err != nil ***REMOVED***
		t.Fatalf("unexpected error in Get: %v", err)
	***REMOVED***

	// Cancel befer reading the body.
	// Reading Request.Body should fail, since the request was
	// canceled before anything was written.
	cancel()

	done := make(chan struct***REMOVED******REMOVED***)

	go func() ***REMOVED***
		b, err := ioutil.ReadAll(resp.Body)
		if len(b) != 0 || err == nil ***REMOVED***
			t.Errorf(`Read got (%q, %v); want ("", error)`, b, err)
		***REMOVED***
		close(done)
	***REMOVED***()

	select ***REMOVED***
	case <-time.After(1 * time.Second):
		t.Errorf("Test timed out")
	case <-done:
	***REMOVED***
***REMOVED***
