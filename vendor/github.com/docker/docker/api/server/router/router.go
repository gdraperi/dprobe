package router

import "github.com/docker/docker/api/server/httputils"

// Router defines an interface to specify a group of routes to add to the docker server.
type Router interface ***REMOVED***
	// Routes returns the list of routes to add to the docker server.
	Routes() []Route
***REMOVED***

// Route defines an individual API route in the docker server.
type Route interface ***REMOVED***
	// Handler returns the raw function to create the http handler.
	Handler() httputils.APIFunc
	// Method returns the http method that the route responds to.
	Method() string
	// Path returns the subpath where the route responds to.
	Path() string
***REMOVED***
