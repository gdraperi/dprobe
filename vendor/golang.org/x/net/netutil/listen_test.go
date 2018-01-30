// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netutil

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"golang.org/x/net/internal/nettest"
)

func TestLimitListener(t *testing.T) ***REMOVED***
	const max = 5
	attempts := (nettest.MaxOpenFiles() - max) / 2
	if attempts > 256 ***REMOVED*** // maximum length of accept queue is 128 by default
		attempts = 256
	***REMOVED***

	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer l.Close()
	l = LimitListener(l, max)

	var open int32
	go http.Serve(l, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		if n := atomic.AddInt32(&open, 1); n > max ***REMOVED***
			t.Errorf("%d open connections, want <= %d", n, max)
		***REMOVED***
		defer atomic.AddInt32(&open, -1)
		time.Sleep(10 * time.Millisecond)
		fmt.Fprint(w, "some body")
	***REMOVED***))

	var wg sync.WaitGroup
	var failed int32
	for i := 0; i < attempts; i++ ***REMOVED***
		wg.Add(1)
		go func() ***REMOVED***
			defer wg.Done()
			c := http.Client***REMOVED***Timeout: 3 * time.Second***REMOVED***
			r, err := c.Get("http://" + l.Addr().String())
			if err != nil ***REMOVED***
				t.Log(err)
				atomic.AddInt32(&failed, 1)
				return
			***REMOVED***
			defer r.Body.Close()
			io.Copy(ioutil.Discard, r.Body)
		***REMOVED***()
	***REMOVED***
	wg.Wait()

	// We expect some Gets to fail as the kernel's accept queue is filled,
	// but most should succeed.
	if int(failed) >= attempts/2 ***REMOVED***
		t.Errorf("%d requests failed within %d attempts", failed, attempts)
	***REMOVED***
***REMOVED***

type errorListener struct ***REMOVED***
	net.Listener
***REMOVED***

func (errorListener) Accept() (net.Conn, error) ***REMOVED***
	return nil, errFake
***REMOVED***

var errFake = errors.New("fake error from errorListener")

// This used to hang.
func TestLimitListenerError(t *testing.T) ***REMOVED***
	donec := make(chan bool, 1)
	go func() ***REMOVED***
		const n = 2
		ll := LimitListener(errorListener***REMOVED******REMOVED***, n)
		for i := 0; i < n+1; i++ ***REMOVED***
			_, err := ll.Accept()
			if err != errFake ***REMOVED***
				t.Fatalf("Accept error = %v; want errFake", err)
			***REMOVED***
		***REMOVED***
		donec <- true
	***REMOVED***()
	select ***REMOVED***
	case <-donec:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout. deadlock?")
	***REMOVED***
***REMOVED***
