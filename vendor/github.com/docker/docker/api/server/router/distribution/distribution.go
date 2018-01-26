package distribution

import "github.com/docker/docker/api/server/router"

// distributionRouter is a router to talk with the registry
type distributionRouter struct ***REMOVED***
	backend Backend
	routes  []router.Route
***REMOVED***

// NewRouter initializes a new distribution router
func NewRouter(backend Backend) router.Router ***REMOVED***
	r := &distributionRouter***REMOVED***
		backend: backend,
	***REMOVED***
	r.initRoutes()
	return r
***REMOVED***

// Routes returns the available routes
func (r *distributionRouter) Routes() []router.Route ***REMOVED***
	return r.routes
***REMOVED***

// initRoutes initializes the routes in the distribution router
func (r *distributionRouter) initRoutes() ***REMOVED***
	r.routes = []router.Route***REMOVED***
		// GET
		router.NewGetRoute("/distribution/***REMOVED***name:.****REMOVED***/json", r.getDistributionInfo),
	***REMOVED***
***REMOVED***
