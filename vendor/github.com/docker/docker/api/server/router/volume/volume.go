package volume

import "github.com/docker/docker/api/server/router"

// volumeRouter is a router to talk with the volumes controller
type volumeRouter struct ***REMOVED***
	backend Backend
	routes  []router.Route
***REMOVED***

// NewRouter initializes a new volume router
func NewRouter(b Backend) router.Router ***REMOVED***
	r := &volumeRouter***REMOVED***
		backend: b,
	***REMOVED***
	r.initRoutes()
	return r
***REMOVED***

// Routes returns the available routes to the volumes controller
func (r *volumeRouter) Routes() []router.Route ***REMOVED***
	return r.routes
***REMOVED***

func (r *volumeRouter) initRoutes() ***REMOVED***
	r.routes = []router.Route***REMOVED***
		// GET
		router.NewGetRoute("/volumes", r.getVolumesList),
		router.NewGetRoute("/volumes/***REMOVED***name:.****REMOVED***", r.getVolumeByName),
		// POST
		router.NewPostRoute("/volumes/create", r.postVolumesCreate),
		router.NewPostRoute("/volumes/prune", r.postVolumesPrune, router.WithCancel),
		// DELETE
		router.NewDeleteRoute("/volumes/***REMOVED***name:.****REMOVED***", r.deleteVolumes),
	***REMOVED***
***REMOVED***
