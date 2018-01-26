package request

import (
	"io"
	"net/http"
	"net/url"
)

func copyHTTPRequest(r *http.Request, body io.ReadCloser) *http.Request ***REMOVED***
	req := new(http.Request)
	*req = *r
	req.URL = &url.URL***REMOVED******REMOVED***
	*req.URL = *r.URL
	req.Body = body

	req.Header = http.Header***REMOVED******REMOVED***
	for k, v := range r.Header ***REMOVED***
		for _, vv := range v ***REMOVED***
			req.Header.Add(k, vv)
		***REMOVED***
	***REMOVED***

	return req
***REMOVED***
