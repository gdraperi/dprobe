package router

import (
	"net/http"

	"github.com/docker/docker/api/server/httputils"
	"golang.org/x/net/context"
)

// RouteWrapper wraps a route with extra functionality.
// It is passed in when creating a new route.
type RouteWrapper func(r Route) Route

// localRoute defines an individual API route to connect
// with the docker daemon. It implements Route.
type localRoute struct ***REMOVED***
	method  string
	path    string
	handler httputils.APIFunc
***REMOVED***

// Handler returns the APIFunc to let the server wrap it in middlewares.
func (l localRoute) Handler() httputils.APIFunc ***REMOVED***
	return l.handler
***REMOVED***

// Method returns the http method that the route responds to.
func (l localRoute) Method() string ***REMOVED***
	return l.method
***REMOVED***

// Path returns the subpath where the route responds to.
func (l localRoute) Path() string ***REMOVED***
	return l.path
***REMOVED***

// NewRoute initializes a new local route for the router.
func NewRoute(method, path string, handler httputils.APIFunc, opts ...RouteWrapper) Route ***REMOVED***
	var r Route = localRoute***REMOVED***method, path, handler***REMOVED***
	for _, o := range opts ***REMOVED***
		r = o(r)
	***REMOVED***
	return r
***REMOVED***

// NewGetRoute initializes a new route with the http method GET.
func NewGetRoute(path string, handler httputils.APIFunc, opts ...RouteWrapper) Route ***REMOVED***
	return NewRoute("GET", path, handler, opts...)
***REMOVED***

// NewPostRoute initializes a new route with the http method POST.
func NewPostRoute(path string, handler httputils.APIFunc, opts ...RouteWrapper) Route ***REMOVED***
	return NewRoute("POST", path, handler, opts...)
***REMOVED***

// NewPutRoute initializes a new route with the http method PUT.
func NewPutRoute(path string, handler httputils.APIFunc, opts ...RouteWrapper) Route ***REMOVED***
	return NewRoute("PUT", path, handler, opts...)
***REMOVED***

// NewDeleteRoute initializes a new route with the http method DELETE.
func NewDeleteRoute(path string, handler httputils.APIFunc, opts ...RouteWrapper) Route ***REMOVED***
	return NewRoute("DELETE", path, handler, opts...)
***REMOVED***

// NewOptionsRoute initializes a new route with the http method OPTIONS.
func NewOptionsRoute(path string, handler httputils.APIFunc, opts ...RouteWrapper) Route ***REMOVED***
	return NewRoute("OPTIONS", path, handler, opts...)
***REMOVED***

// NewHeadRoute initializes a new route with the http method HEAD.
func NewHeadRoute(path string, handler httputils.APIFunc, opts ...RouteWrapper) Route ***REMOVED***
	return NewRoute("HEAD", path, handler, opts...)
***REMOVED***

func cancellableHandler(h httputils.APIFunc) httputils.APIFunc ***REMOVED***
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
		if notifier, ok := w.(http.CloseNotifier); ok ***REMOVED***
			notify := notifier.CloseNotify()
			notifyCtx, cancel := context.WithCancel(ctx)
			finished := make(chan struct***REMOVED******REMOVED***)
			defer close(finished)
			ctx = notifyCtx
			go func() ***REMOVED***
				select ***REMOVED***
				case <-notify:
					cancel()
				case <-finished:
				***REMOVED***
			***REMOVED***()
		***REMOVED***
		return h(ctx, w, r, vars)
	***REMOVED***
***REMOVED***

// WithCancel makes new route which embeds http.CloseNotifier feature to
// context.Context of handler.
func WithCancel(r Route) Route ***REMOVED***
	return localRoute***REMOVED***
		method:  r.Method(),
		path:    r.Path(),
		handler: cancellableHandler(r.Handler()),
	***REMOVED***
***REMOVED***
