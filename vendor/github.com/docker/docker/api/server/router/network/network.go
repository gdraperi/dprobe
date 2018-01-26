package network

import (
	"github.com/docker/docker/api/server/router"
	"github.com/docker/docker/daemon/cluster"
)

// networkRouter is a router to talk with the network controller
type networkRouter struct ***REMOVED***
	backend Backend
	cluster *cluster.Cluster
	routes  []router.Route
***REMOVED***

// NewRouter initializes a new network router
func NewRouter(b Backend, c *cluster.Cluster) router.Router ***REMOVED***
	r := &networkRouter***REMOVED***
		backend: b,
		cluster: c,
	***REMOVED***
	r.initRoutes()
	return r
***REMOVED***

// Routes returns the available routes to the network controller
func (r *networkRouter) Routes() []router.Route ***REMOVED***
	return r.routes
***REMOVED***

func (r *networkRouter) initRoutes() ***REMOVED***
	r.routes = []router.Route***REMOVED***
		// GET
		router.NewGetRoute("/networks", r.getNetworksList),
		router.NewGetRoute("/networks/", r.getNetworksList),
		router.NewGetRoute("/networks/***REMOVED***id:.+***REMOVED***", r.getNetwork),
		// POST
		router.NewPostRoute("/networks/create", r.postNetworkCreate),
		router.NewPostRoute("/networks/***REMOVED***id:.****REMOVED***/connect", r.postNetworkConnect),
		router.NewPostRoute("/networks/***REMOVED***id:.****REMOVED***/disconnect", r.postNetworkDisconnect),
		router.NewPostRoute("/networks/prune", r.postNetworksPrune, router.WithCancel),
		// DELETE
		router.NewDeleteRoute("/networks/***REMOVED***id:.****REMOVED***", r.deleteNetwork),
	***REMOVED***
***REMOVED***
