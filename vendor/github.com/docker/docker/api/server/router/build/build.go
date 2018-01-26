package build

import "github.com/docker/docker/api/server/router"

// buildRouter is a router to talk with the build controller
type buildRouter struct ***REMOVED***
	backend Backend
	daemon  experimentalProvider
	routes  []router.Route
***REMOVED***

// NewRouter initializes a new build router
func NewRouter(b Backend, d experimentalProvider) router.Router ***REMOVED***
	r := &buildRouter***REMOVED***backend: b, daemon: d***REMOVED***
	r.initRoutes()
	return r
***REMOVED***

// Routes returns the available routers to the build controller
func (r *buildRouter) Routes() []router.Route ***REMOVED***
	return r.routes
***REMOVED***

func (r *buildRouter) initRoutes() ***REMOVED***
	r.routes = []router.Route***REMOVED***
		router.NewPostRoute("/build", r.postBuild, router.WithCancel),
		router.NewPostRoute("/build/prune", r.postPrune, router.WithCancel),
	***REMOVED***
***REMOVED***
