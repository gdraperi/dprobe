package image

import (
	"github.com/docker/docker/api/server/httputils"
	"github.com/docker/docker/api/server/router"
)

// imageRouter is a router to talk with the image controller
type imageRouter struct ***REMOVED***
	backend Backend
	decoder httputils.ContainerDecoder
	routes  []router.Route
***REMOVED***

// NewRouter initializes a new image router
func NewRouter(backend Backend, decoder httputils.ContainerDecoder) router.Router ***REMOVED***
	r := &imageRouter***REMOVED***
		backend: backend,
		decoder: decoder,
	***REMOVED***
	r.initRoutes()
	return r
***REMOVED***

// Routes returns the available routes to the image controller
func (r *imageRouter) Routes() []router.Route ***REMOVED***
	return r.routes
***REMOVED***

// initRoutes initializes the routes in the image router
func (r *imageRouter) initRoutes() ***REMOVED***
	r.routes = []router.Route***REMOVED***
		// GET
		router.NewGetRoute("/images/json", r.getImagesJSON),
		router.NewGetRoute("/images/search", r.getImagesSearch),
		router.NewGetRoute("/images/get", r.getImagesGet),
		router.NewGetRoute("/images/***REMOVED***name:.****REMOVED***/get", r.getImagesGet),
		router.NewGetRoute("/images/***REMOVED***name:.****REMOVED***/history", r.getImagesHistory),
		router.NewGetRoute("/images/***REMOVED***name:.****REMOVED***/json", r.getImagesByName),
		// POST
		router.NewPostRoute("/commit", r.postCommit),
		router.NewPostRoute("/images/load", r.postImagesLoad),
		router.NewPostRoute("/images/create", r.postImagesCreate, router.WithCancel),
		router.NewPostRoute("/images/***REMOVED***name:.****REMOVED***/push", r.postImagesPush, router.WithCancel),
		router.NewPostRoute("/images/***REMOVED***name:.****REMOVED***/tag", r.postImagesTag),
		router.NewPostRoute("/images/prune", r.postImagesPrune, router.WithCancel),
		// DELETE
		router.NewDeleteRoute("/images/***REMOVED***name:.****REMOVED***", r.deleteImages),
	***REMOVED***
***REMOVED***
