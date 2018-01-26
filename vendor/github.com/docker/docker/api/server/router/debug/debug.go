package debug

import (
	"expvar"
	"net/http"
	"net/http/pprof"

	"github.com/docker/docker/api/server/httputils"
	"github.com/docker/docker/api/server/router"
	"golang.org/x/net/context"
)

// NewRouter creates a new debug router
// The debug router holds endpoints for debug the daemon, such as those for pprof.
func NewRouter() router.Router ***REMOVED***
	r := &debugRouter***REMOVED******REMOVED***
	r.initRoutes()
	return r
***REMOVED***

type debugRouter struct ***REMOVED***
	routes []router.Route
***REMOVED***

func (r *debugRouter) initRoutes() ***REMOVED***
	r.routes = []router.Route***REMOVED***
		router.NewGetRoute("/vars", frameworkAdaptHandler(expvar.Handler())),
		router.NewGetRoute("/pprof/", frameworkAdaptHandlerFunc(pprof.Index)),
		router.NewGetRoute("/pprof/cmdline", frameworkAdaptHandlerFunc(pprof.Cmdline)),
		router.NewGetRoute("/pprof/profile", frameworkAdaptHandlerFunc(pprof.Profile)),
		router.NewGetRoute("/pprof/symbol", frameworkAdaptHandlerFunc(pprof.Symbol)),
		router.NewGetRoute("/pprof/trace", frameworkAdaptHandlerFunc(pprof.Trace)),
		router.NewGetRoute("/pprof/***REMOVED***name***REMOVED***", handlePprof),
	***REMOVED***
***REMOVED***

func (r *debugRouter) Routes() []router.Route ***REMOVED***
	return r.routes
***REMOVED***

func frameworkAdaptHandler(handler http.Handler) httputils.APIFunc ***REMOVED***
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
		handler.ServeHTTP(w, r)
		return nil
	***REMOVED***
***REMOVED***

func frameworkAdaptHandlerFunc(handler http.HandlerFunc) httputils.APIFunc ***REMOVED***
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
		handler(w, r)
		return nil
	***REMOVED***
***REMOVED***
