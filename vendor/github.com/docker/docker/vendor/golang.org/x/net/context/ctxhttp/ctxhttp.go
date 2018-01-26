// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.7

// Package ctxhttp provides helper functions for performing context-aware HTTP requests.
package ctxhttp // import "golang.org/x/net/context/ctxhttp"

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/context"
)

// Do sends an HTTP request with the provided http.Client and returns
// an HTTP response.
//
// If the client is nil, http.DefaultClient is used.
//
// The provided ctx must be non-nil. If it is canceled or times out,
// ctx.Err() will be returned.
func Do(ctx context.Context, client *http.Client, req *http.Request) (*http.Response, error) ***REMOVED***
	if client == nil ***REMOVED***
		client = http.DefaultClient
	***REMOVED***
	resp, err := client.Do(req.WithContext(ctx))
	// If we got an error, and the context has been canceled,
	// the context's error is probably more useful.
	if err != nil ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
			err = ctx.Err()
		default:
		***REMOVED***
	***REMOVED***
	return resp, err
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
