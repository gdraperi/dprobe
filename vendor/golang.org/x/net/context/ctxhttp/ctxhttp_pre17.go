// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !go1.7

package ctxhttp // import "golang.org/x/net/context/ctxhttp"

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/context"
)

func nop() ***REMOVED******REMOVED***

var (
	testHookContextDoneBeforeHeaders = nop
	testHookDoReturned               = nop
	testHookDidBodyClose             = nop
)

// Do sends an HTTP request with the provided http.Client and returns an HTTP response.
// If the client is nil, http.DefaultClient is used.
// If the context is canceled or times out, ctx.Err() will be returned.
func Do(ctx context.Context, client *http.Client, req *http.Request) (*http.Response, error) ***REMOVED***
	if client == nil ***REMOVED***
		client = http.DefaultClient
	***REMOVED***

	// TODO(djd): Respect any existing value of req.Cancel.
	cancel := make(chan struct***REMOVED******REMOVED***)
	req.Cancel = cancel

	type responseAndError struct ***REMOVED***
		resp *http.Response
		err  error
	***REMOVED***
	result := make(chan responseAndError, 1)

	// Make local copies of test hooks closed over by goroutines below.
	// Prevents data races in tests.
	testHookDoReturned := testHookDoReturned
	testHookDidBodyClose := testHookDidBodyClose

	go func() ***REMOVED***
		resp, err := client.Do(req)
		testHookDoReturned()
		result <- responseAndError***REMOVED***resp, err***REMOVED***
	***REMOVED***()

	var resp *http.Response

	select ***REMOVED***
	case <-ctx.Done():
		testHookContextDoneBeforeHeaders()
		close(cancel)
		// Clean up after the goroutine calling client.Do:
		go func() ***REMOVED***
			if r := <-result; r.resp != nil ***REMOVED***
				testHookDidBodyClose()
				r.resp.Body.Close()
			***REMOVED***
		***REMOVED***()
		return nil, ctx.Err()
	case r := <-result:
		var err error
		resp, err = r.resp, r.err
		if err != nil ***REMOVED***
			return resp, err
		***REMOVED***
	***REMOVED***

	c := make(chan struct***REMOVED******REMOVED***)
	go func() ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
			close(cancel)
		case <-c:
			// The response's Body is closed.
		***REMOVED***
	***REMOVED***()
	resp.Body = &notifyingReader***REMOVED***resp.Body, c***REMOVED***

	return resp, nil
***REMOVED***

// Get issues a GET request via the Do function.
func Get(ctx context.Context, client *http.Client, url string) (*http.Response, error) ***REMOVED***
	req, err := http.NewRequest("GET", url, nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return Do(ctx, client, req)
***REMOVED***

// Head issues a HEAD request via the Do function.
func Head(ctx context.Context, client *http.Client, url string) (*http.Response, error) ***REMOVED***
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return Do(ctx, client, req)
***REMOVED***

// Post issues a POST request via the Do function.
func Post(ctx context.Context, client *http.Client, url string, bodyType string, body io.Reader) (*http.Response, error) ***REMOVED***
	req, err := http.NewRequest("POST", url, body)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	req.Header.Set("Content-Type", bodyType)
	return Do(ctx, client, req)
***REMOVED***

// PostForm issues a POST request via the Do function.
func PostForm(ctx context.Context, client *http.Client, url string, data url.Values) (*http.Response, error) ***REMOVED***
	return Post(ctx, client, url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
***REMOVED***

// notifyingReader is an io.ReadCloser that closes the notify channel after
// Close is called or a Read fails on the underlying ReadCloser.
type notifyingReader struct ***REMOVED***
	io.ReadCloser
	notify chan<- struct***REMOVED******REMOVED***
***REMOVED***

func (r *notifyingReader) Read(p []byte) (int, error) ***REMOVED***
	n, err := r.ReadCloser.Read(p)
	if err != nil && r.notify != nil ***REMOVED***
		close(r.notify)
		r.notify = nil
	***REMOVED***
	return n, err
***REMOVED***

func (r *notifyingReader) Close() error ***REMOVED***
	err := r.ReadCloser.Close()
	if r.notify != nil ***REMOVED***
		close(r.notify)
		r.notify = nil
	***REMOVED***
	return err
***REMOVED***
