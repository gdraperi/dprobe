package session

import "github.com/docker/docker/api/server/router"

// sessionRouter is a router to talk with the session controller
type sessionRouter struct ***REMOVED***
	backend Backend
	routes  []router.Route
***REMOVED***

// NewRouter initializes a new session router
func NewRouter(b Backend) router.Router ***REMOVED***
	r := &sessionRouter***REMOVED***
		backend: b,
	***REMOVED***
	r.initRoutes()
	return r
***REMOVED***

// Routes returns the available routers to the session controller
func (r *sessionRouter) Routes() []router.Route ***REMOVED***
	return r.routes
***REMOVED***

func (r *sessionRouter) initRoutes() ***REMOVED***
	r.routes = []router.Route***REMOVED***
		router.Experimental(router.NewPostRoute("/session", r.startSession)),
	***REMOVED***
***REMOVED***
