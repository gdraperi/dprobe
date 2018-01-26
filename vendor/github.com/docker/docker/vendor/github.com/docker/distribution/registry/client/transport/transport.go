package transport

import (
	"io"
	"net/http"
	"sync"
)

// RequestModifier represents an object which will do an inplace
// modification of an HTTP request.
type RequestModifier interface ***REMOVED***
	ModifyRequest(*http.Request) error
***REMOVED***

type headerModifier http.Header

// NewHeaderRequestModifier returns a new RequestModifier which will
// add the given headers to a request.
func NewHeaderRequestModifier(header http.Header) RequestModifier ***REMOVED***
	return headerModifier(header)
***REMOVED***

func (h headerModifier) ModifyRequest(req *http.Request) error ***REMOVED***
	for k, s := range http.Header(h) ***REMOVED***
		req.Header[k] = append(req.Header[k], s...)
	***REMOVED***

	return nil
***REMOVED***

// NewTransport creates a new transport which will apply modifiers to
// the request on a RoundTrip call.
func NewTransport(base http.RoundTripper, modifiers ...RequestModifier) http.RoundTripper ***REMOVED***
	return &transport***REMOVED***
		Modifiers: modifiers,
		Base:      base,
	***REMOVED***
***REMOVED***

// transport is an http.RoundTripper that makes HTTP requests after
// copying and modifying the request
type transport struct ***REMOVED***
	Modifiers []RequestModifier
	Base      http.RoundTripper

	mu     sync.Mutex                      // guards modReq
	modReq map[*http.Request]*http.Request // original -> modified
***REMOVED***

// RoundTrip authorizes and authenticates the request with an
// access token. If no token exists or token is expired,
// tries to refresh/fetch a new token.
func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) ***REMOVED***
	req2 := cloneRequest(req)
	for _, modifier := range t.Modifiers ***REMOVED***
		if err := modifier.ModifyRequest(req2); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

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
func (t *transport) CancelRequest(req *http.Request) ***REMOVED***
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

func (t *transport) base() http.RoundTripper ***REMOVED***
	if t.Base != nil ***REMOVED***
		return t.Base
	***REMOVED***
	return http.DefaultTransport
***REMOVED***

func (t *transport) setModReq(orig, mod *http.Request) ***REMOVED***
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
