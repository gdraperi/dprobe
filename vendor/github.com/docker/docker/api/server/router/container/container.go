package container

import (
	"github.com/docker/docker/api/server/httputils"
	"github.com/docker/docker/api/server/router"
)

// containerRouter is a router to talk with the container controller
type containerRouter struct ***REMOVED***
	backend Backend
	decoder httputils.ContainerDecoder
	routes  []router.Route
***REMOVED***

// NewRouter initializes a new container router
func NewRouter(b Backend, decoder httputils.ContainerDecoder) router.Router ***REMOVED***
	r := &containerRouter***REMOVED***
		backend: b,
		decoder: decoder,
	***REMOVED***
	r.initRoutes()
	return r
***REMOVED***

// Routes returns the available routes to the container controller
func (r *containerRouter) Routes() []router.Route ***REMOVED***
	return r.routes
***REMOVED***

// initRoutes initializes the routes in container router
func (r *containerRouter) initRoutes() ***REMOVED***
	r.routes = []router.Route***REMOVED***
		// HEAD
		router.NewHeadRoute("/containers/***REMOVED***name:.****REMOVED***/archive", r.headContainersArchive),
		// GET
		router.NewGetRoute("/containers/json", r.getContainersJSON),
		router.NewGetRoute("/containers/***REMOVED***name:.****REMOVED***/export", r.getContainersExport),
		router.NewGetRoute("/containers/***REMOVED***name:.****REMOVED***/changes", r.getContainersChanges),
		router.NewGetRoute("/containers/***REMOVED***name:.****REMOVED***/json", r.getContainersByName),
		router.NewGetRoute("/containers/***REMOVED***name:.****REMOVED***/top", r.getContainersTop),
		router.NewGetRoute("/containers/***REMOVED***name:.****REMOVED***/logs", r.getContainersLogs, router.WithCancel),
		router.NewGetRoute("/containers/***REMOVED***name:.****REMOVED***/stats", r.getContainersStats, router.WithCancel),
		router.NewGetRoute("/containers/***REMOVED***name:.****REMOVED***/attach/ws", r.wsContainersAttach),
		router.NewGetRoute("/exec/***REMOVED***id:.****REMOVED***/json", r.getExecByID),
		router.NewGetRoute("/containers/***REMOVED***name:.****REMOVED***/archive", r.getContainersArchive),
		// POST
		router.NewPostRoute("/containers/create", r.postContainersCreate),
		router.NewPostRoute("/containers/***REMOVED***name:.****REMOVED***/kill", r.postContainersKill),
		router.NewPostRoute("/containers/***REMOVED***name:.****REMOVED***/pause", r.postContainersPause),
		router.NewPostRoute("/containers/***REMOVED***name:.****REMOVED***/unpause", r.postContainersUnpause),
		router.NewPostRoute("/containers/***REMOVED***name:.****REMOVED***/restart", r.postContainersRestart),
		router.NewPostRoute("/containers/***REMOVED***name:.****REMOVED***/start", r.postContainersStart),
		router.NewPostRoute("/containers/***REMOVED***name:.****REMOVED***/stop", r.postContainersStop),
		router.NewPostRoute("/containers/***REMOVED***name:.****REMOVED***/wait", r.postContainersWait, router.WithCancel),
		router.NewPostRoute("/containers/***REMOVED***name:.****REMOVED***/resize", r.postContainersResize),
		router.NewPostRoute("/containers/***REMOVED***name:.****REMOVED***/attach", r.postContainersAttach),
		router.NewPostRoute("/containers/***REMOVED***name:.****REMOVED***/copy", r.postContainersCopy), // Deprecated since 1.8, Errors out since 1.12
		router.NewPostRoute("/containers/***REMOVED***name:.****REMOVED***/exec", r.postContainerExecCreate),
		router.NewPostRoute("/exec/***REMOVED***name:.****REMOVED***/start", r.postContainerExecStart),
		router.NewPostRoute("/exec/***REMOVED***name:.****REMOVED***/resize", r.postContainerExecResize),
		router.NewPostRoute("/containers/***REMOVED***name:.****REMOVED***/rename", r.postContainerRename),
		router.NewPostRoute("/containers/***REMOVED***name:.****REMOVED***/update", r.postContainerUpdate),
		router.NewPostRoute("/containers/prune", r.postContainersPrune, router.WithCancel),
		// PUT
		router.NewPutRoute("/containers/***REMOVED***name:.****REMOVED***/archive", r.putContainersArchive),
		// DELETE
		router.NewDeleteRoute("/containers/***REMOVED***name:.****REMOVED***", r.deleteContainers),
	***REMOVED***
***REMOVED***
