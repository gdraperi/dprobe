package checkpoint

import (
	"github.com/docker/docker/api/server/httputils"
	"github.com/docker/docker/api/server/router"
)

// checkpointRouter is a router to talk with the checkpoint controller
type checkpointRouter struct ***REMOVED***
	backend Backend
	decoder httputils.ContainerDecoder
	routes  []router.Route
***REMOVED***

// NewRouter initializes a new checkpoint router
func NewRouter(b Backend, decoder httputils.ContainerDecoder) router.Router ***REMOVED***
	r := &checkpointRouter***REMOVED***
		backend: b,
		decoder: decoder,
	***REMOVED***
	r.initRoutes()
	return r
***REMOVED***

// Routes returns the available routers to the checkpoint controller
func (r *checkpointRouter) Routes() []router.Route ***REMOVED***
	return r.routes
***REMOVED***

func (r *checkpointRouter) initRoutes() ***REMOVED***
	r.routes = []router.Route***REMOVED***
		router.NewGetRoute("/containers/***REMOVED***name:.****REMOVED***/checkpoints", r.getContainerCheckpoints, router.Experimental),
		router.NewPostRoute("/containers/***REMOVED***name:.****REMOVED***/checkpoints", r.postContainerCheckpoint, router.Experimental),
		router.NewDeleteRoute("/containers/***REMOVED***name***REMOVED***/checkpoints/***REMOVED***checkpoint***REMOVED***", r.deleteContainerCheckpoint, router.Experimental),
	***REMOVED***
***REMOVED***
