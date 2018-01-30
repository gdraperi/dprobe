// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !plan9,!go1.7

package ctxhttp

import (
	"net"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"golang.org/x/net/context"
)

// golang.org/issue/14065
func TestClosesResponseBodyOnCancel(t *testing.T) ***REMOVED***
	defer func() ***REMOVED*** testHookContextDoneBeforeHeaders = nop ***REMOVED***()
	defer func() ***REMOVED*** testHookDoReturned = nop ***REMOVED***()
	defer func() ***REMOVED*** testHookDidBodyClose = nop ***REMOVED***()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED******REMOVED***))
	defer ts.Close()

	ctx, cancel := context.WithCancel(context.Background())

	// closed when Do enters select case <-ctx.Done()
	enteredDonePath := make(chan struct***REMOVED******REMOVED***)

	testHookContextDoneBeforeHeaders = func() ***REMOVED***
		close(enteredDonePath)
	***REMOVED***

	testHookDoReturned = func() ***REMOVED***
		// We now have the result (the Flush'd headers) at least,
		// so we can cancel the request.
		cancel()

		// But block the client.Do goroutine from sending
		// until Do enters into the <-ctx.Done() path, since
		// otherwise if both channels are readable, select
		// picks a random one.
		<-enteredDonePath
	***REMOVED***

	sawBodyClose := make(chan struct***REMOVED******REMOVED***)
	testHookDidBodyClose = func() ***REMOVED*** close(sawBodyClose) ***REMOVED***

	tr := &http.Transport***REMOVED******REMOVED***
	defer tr.CloseIdleConnections()
	c := &http.Client***REMOVED***Transport: tr***REMOVED***
	req, _ := http.NewRequest("GET", ts.URL, nil)
	_, doErr := Do(ctx, c, req)

	select ***REMOVED***
	case <-sawBodyClose:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for body to close")
	***REMOVED***

	if doErr != ctx.Err() ***REMOVED***
		t.Errorf("Do error = %v; want %v", doErr, ctx.Err())
	***REMOVED***
***REMOVED***

type noteCloseConn struct ***REMOVED***
	net.Conn
	onceClose sync.Once
	closefn   func()
***REMOVED***

func (c *noteCloseConn) Close() error ***REMOVED***
	c.onceClose.Do(c.closefn)
	return c.Conn.Close()
***REMOVED***
