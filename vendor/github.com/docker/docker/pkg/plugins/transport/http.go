package transport

import (
	"io"
	"net/http"
)

// httpTransport holds an http.RoundTripper
// and information about the scheme and address the transport
// sends request to.
type httpTransport struct ***REMOVED***
	http.RoundTripper
	scheme string
	addr   string
***REMOVED***

// NewHTTPTransport creates a new httpTransport.
func NewHTTPTransport(r http.RoundTripper, scheme, addr string) Transport ***REMOVED***
	return httpTransport***REMOVED***
		RoundTripper: r,
		scheme:       scheme,
		addr:         addr,
	***REMOVED***
***REMOVED***

// NewRequest creates a new http.Request and sets the URL
// scheme and address with the transport's fields.
func (t httpTransport) NewRequest(path string, data io.Reader) (*http.Request, error) ***REMOVED***
	req, err := newHTTPRequest(path, data)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	req.URL.Scheme = t.scheme
	req.URL.Host = t.addr
	return req, nil
***REMOVED***
