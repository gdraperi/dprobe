package transport

import (
	"io"
	"net/http"
	"strings"
)

// VersionMimetype is the Content-Type the engine sends to plugins.
const VersionMimetype = "application/vnd.docker.plugins.v1.2+json"

// RequestFactory defines an interface that
// transports can implement to create new requests.
type RequestFactory interface ***REMOVED***
	NewRequest(path string, data io.Reader) (*http.Request, error)
***REMOVED***

// Transport defines an interface that plugin transports
// must implement.
type Transport interface ***REMOVED***
	http.RoundTripper
	RequestFactory
***REMOVED***

// newHTTPRequest creates a new request with a path and a body.
func newHTTPRequest(path string, data io.Reader) (*http.Request, error) ***REMOVED***
	if !strings.HasPrefix(path, "/") ***REMOVED***
		path = "/" + path
	***REMOVED***
	req, err := http.NewRequest("POST", path, data)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	req.Header.Add("Accept", VersionMimetype)
	return req, nil
***REMOVED***
