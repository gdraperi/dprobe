package swarm

import "github.com/docker/docker/api/server/router"

// swarmRouter is a router to talk with the build controller
type swarmRouter struct ***REMOVED***
	backend Backend
	routes  []router.Route
***REMOVED***

// NewRouter initializes a new build router
func NewRouter(b Backend) router.Router ***REMOVED***
	r := &swarmRouter***REMOVED***
		backend: b,
	***REMOVED***
	r.initRoutes()
	return r
***REMOVED***

// Routes returns the available routers to the swarm controller
func (sr *swarmRouter) Routes() []router.Route ***REMOVED***
	return sr.routes
***REMOVED***

func (sr *swarmRouter) initRoutes() ***REMOVED***
	sr.routes = []router.Route***REMOVED***
		router.NewPostRoute("/swarm/init", sr.initCluster),
		router.NewPostRoute("/swarm/join", sr.joinCluster),
		router.NewPostRoute("/swarm/leave", sr.leaveCluster),
		router.NewGetRoute("/swarm", sr.inspectCluster),
		router.NewGetRoute("/swarm/unlockkey", sr.getUnlockKey),
		router.NewPostRoute("/swarm/update", sr.updateCluster),
		router.NewPostRoute("/swarm/unlock", sr.unlockCluster),

		router.NewGetRoute("/services", sr.getServices),
		router.NewGetRoute("/services/***REMOVED***id***REMOVED***", sr.getService),
		router.NewPostRoute("/services/create", sr.createService),
		router.NewPostRoute("/services/***REMOVED***id***REMOVED***/update", sr.updateService),
		router.NewDeleteRoute("/services/***REMOVED***id***REMOVED***", sr.removeService),
		router.NewGetRoute("/services/***REMOVED***id***REMOVED***/logs", sr.getServiceLogs, router.WithCancel),

		router.NewGetRoute("/nodes", sr.getNodes),
		router.NewGetRoute("/nodes/***REMOVED***id***REMOVED***", sr.getNode),
		router.NewDeleteRoute("/nodes/***REMOVED***id***REMOVED***", sr.removeNode),
		router.NewPostRoute("/nodes/***REMOVED***id***REMOVED***/update", sr.updateNode),

		router.NewGetRoute("/tasks", sr.getTasks),
		router.NewGetRoute("/tasks/***REMOVED***id***REMOVED***", sr.getTask),
		router.NewGetRoute("/tasks/***REMOVED***id***REMOVED***/logs", sr.getTaskLogs, router.WithCancel),

		router.NewGetRoute("/secrets", sr.getSecrets),
		router.NewPostRoute("/secrets/create", sr.createSecret),
		router.NewDeleteRoute("/secrets/***REMOVED***id***REMOVED***", sr.removeSecret),
		router.NewGetRoute("/secrets/***REMOVED***id***REMOVED***", sr.getSecret),
		router.NewPostRoute("/secrets/***REMOVED***id***REMOVED***/update", sr.updateSecret),

		router.NewGetRoute("/configs", sr.getConfigs),
		router.NewPostRoute("/configs/create", sr.createConfig),
		router.NewDeleteRoute("/configs/***REMOVED***id***REMOVED***", sr.removeConfig),
		router.NewGetRoute("/configs/***REMOVED***id***REMOVED***", sr.getConfig),
		router.NewPostRoute("/configs/***REMOVED***id***REMOVED***/update", sr.updateConfig),
	***REMOVED***
***REMOVED***
