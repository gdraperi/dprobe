package middleware

import (
	"net/http"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// CORSMiddleware injects CORS headers to each request
// when it's configured.
type CORSMiddleware struct ***REMOVED***
	defaultHeaders string
***REMOVED***

// NewCORSMiddleware creates a new CORSMiddleware with default headers.
func NewCORSMiddleware(d string) CORSMiddleware ***REMOVED***
	return CORSMiddleware***REMOVED***defaultHeaders: d***REMOVED***
***REMOVED***

// WrapHandler returns a new handler function wrapping the previous one in the request chain.
func (c CORSMiddleware) WrapHandler(handler func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error) func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
		// If "api-cors-header" is not given, but "api-enable-cors" is true, we set cors to "*"
		// otherwise, all head values will be passed to HTTP handler
		corsHeaders := c.defaultHeaders
		if corsHeaders == "" ***REMOVED***
			corsHeaders = "*"
		***REMOVED***

		logrus.Debugf("CORS header is enabled and set to: %s", corsHeaders)
		w.Header().Add("Access-Control-Allow-Origin", corsHeaders)
		w.Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, X-Registry-Auth")
		w.Header().Add("Access-Control-Allow-Methods", "HEAD, GET, POST, DELETE, PUT, OPTIONS")
		return handler(ctx, w, r, vars)
	***REMOVED***
***REMOVED***
