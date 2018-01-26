package router

import (
	"net/http"

	"golang.org/x/net/context"

	"github.com/docker/docker/api/server/httputils"
)

// ExperimentalRoute defines an experimental API route that can be enabled or disabled.
type ExperimentalRoute interface ***REMOVED***
	Route

	Enable()
	Disable()
***REMOVED***

// experimentalRoute defines an experimental API route that can be enabled or disabled.
// It implements ExperimentalRoute
type experimentalRoute struct ***REMOVED***
	local   Route
	handler httputils.APIFunc
***REMOVED***

// Enable enables this experimental route
func (r *experimentalRoute) Enable() ***REMOVED***
	r.handler = r.local.Handler()
***REMOVED***

// Disable disables the experimental route
func (r *experimentalRoute) Disable() ***REMOVED***
	r.handler = experimentalHandler
***REMOVED***

type notImplementedError struct***REMOVED******REMOVED***

func (notImplementedError) Error() string ***REMOVED***
	return "This experimental feature is disabled by default. Start the Docker daemon in experimental mode in order to enable it."
***REMOVED***

func (notImplementedError) NotImplemented() ***REMOVED******REMOVED***

func experimentalHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	return notImplementedError***REMOVED******REMOVED***
***REMOVED***

// Handler returns returns the APIFunc to let the server wrap it in middlewares.
func (r *experimentalRoute) Handler() httputils.APIFunc ***REMOVED***
	return r.handler
***REMOVED***

// Method returns the http method that the route responds to.
func (r *experimentalRoute) Method() string ***REMOVED***
	return r.local.Method()
***REMOVED***

// Path returns the subpath where the route responds to.
func (r *experimentalRoute) Path() string ***REMOVED***
	return r.local.Path()
***REMOVED***

// Experimental will mark a route as experimental.
func Experimental(r Route) Route ***REMOVED***
	return &experimentalRoute***REMOVED***
		local:   r,
		handler: experimentalHandler,
	***REMOVED***
***REMOVED***
