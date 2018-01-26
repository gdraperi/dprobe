// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package oauth2

import (
	"errors"
	"io"
	"net/http"
	"sync"
)

// Transport is an http.RoundTripper that makes OAuth 2.0 HTTP requests,
// wrapping a base RoundTripper and adding an Authorization header
// with a token from the supplied Sources.
//
// Transport is a low-level mechanism. Most code will use the
// higher-level Config.Client method instead.
type Transport struct ***REMOVED***
	// Source supplies the token to add to outgoing requests'
	// Authorization headers.
	Source TokenSource

	// Base is the base RoundTripper used to make HTTP requests.
	// If nil, http.DefaultTransport is used.
	Base http.RoundTripper

	mu     sync.Mutex                      // guards modReq
	modReq map[*http.Request]*http.Request // original -> modified
***REMOVED***

// RoundTrip authorizes and authenticates the request with an
// access token. If no token exists or token is expired,
// tries to refresh/fetch a new token.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) ***REMOVED***
	if t.Source == nil ***REMOVED***
		return nil, errors.New("oauth2: Transport's Source is nil")
	***REMOVED***
	token, err := t.Source.Token()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	req2 := cloneRequest(req) // per RoundTripper contract
	token.SetAuthHeader(req2)
	t.setModReq(req, req2)
	res, err := t.base().RoundTrip(req2)
	if err != nil ***REMOVED***
		t.setModReq(req, nil)
		return nil, err
	***REMOVED***
	res.Body = &onEOFReader***REMOVED***
		rc: res.Body,
		fn: func() ***REMOVED*** t.setModReq(req, nil) ***REMOVED***,
	***REMOVED***
	return res, nil
***REMOVED***

// CancelRequest cancels an in-flight request by closing its connection.
func (t *Transport) CancelRequest(req *http.Request) ***REMOVED***
	type canceler interface ***REMOVED***
		CancelRequest(*http.Request)
	***REMOVED***
	if cr, ok := t.base().(canceler); ok ***REMOVED***
		t.mu.Lock()
		modReq := t.modReq[req]
		delete(t.modReq, req)
		t.mu.Unlock()
		cr.CancelRequest(modReq)
	***REMOVED***
***REMOVED***

func (t *Transport) base() http.RoundTripper ***REMOVED***
	if t.Base != nil ***REMOVED***
		return t.Base
	***REMOVED***
	return http.DefaultTransport
***REMOVED***

func (t *Transport) setModReq(orig, mod *http.Request) ***REMOVED***
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.modReq == nil ***REMOVED***
		t.modReq = make(map[*http.Request]*http.Request)
	***REMOVED***
	if mod == nil ***REMOVED***
		delete(t.modReq, orig)
	***REMOVED*** else ***REMOVED***
		t.modReq[orig] = mod
	***REMOVED***
***REMOVED***

// cloneRequest returns a clone of the provided *http.Request.
// The clone is a shallow copy of the struct and its Header map.
func cloneRequest(r *http.Request) *http.Request ***REMOVED***
	// shallow copy of the struct
	r2 := new(http.Request)
	*r2 = *r
	// deep copy of the Header
	r2.Header = make(http.Header, len(r.Header))
	for k, s := range r.Header ***REMOVED***
		r2.Header[k] = append([]string(nil), s...)
	***REMOVED***
	return r2
***REMOVED***

type onEOFReader struct ***REMOVED***
	rc io.ReadCloser
	fn func()
***REMOVED***

func (r *onEOFReader) Read(p []byte) (n int, err error) ***REMOVED***
	n, err = r.rc.Read(p)
	if err == io.EOF ***REMOVED***
		r.runFunc()
	***REMOVED***
	return
***REMOVED***

func (r *onEOFReader) Close() error ***REMOVED***
	err := r.rc.Close()
	r.runFunc()
	return err
***REMOVED***

func (r *onEOFReader) runFunc() ***REMOVED***
	if fn := r.fn; fn != nil ***REMOVED***
		fn()
		r.fn = nil
	***REMOVED***
***REMOVED***
