package middleware

import (
	"net/http"

	"golang.org/x/net/context"
)

// ExperimentalMiddleware is a the middleware in charge of adding the
// 'Docker-Experimental' header to every outgoing request
type ExperimentalMiddleware struct ***REMOVED***
	experimental string
***REMOVED***

// NewExperimentalMiddleware creates a new ExperimentalMiddleware
func NewExperimentalMiddleware(experimentalEnabled bool) ExperimentalMiddleware ***REMOVED***
	if experimentalEnabled ***REMOVED***
		return ExperimentalMiddleware***REMOVED***"true"***REMOVED***
	***REMOVED***
	return ExperimentalMiddleware***REMOVED***"false"***REMOVED***
***REMOVED***

// WrapHandler returns a new handler function wrapping the previous one in the request chain.
func (e ExperimentalMiddleware) WrapHandler(handler func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error) func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
		w.Header().Set("Docker-Experimental", e.experimental)
		return handler(ctx, w, r, vars)
	***REMOVED***
***REMOVED***
