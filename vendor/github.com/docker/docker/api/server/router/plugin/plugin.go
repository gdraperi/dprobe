package plugin

import "github.com/docker/docker/api/server/router"

// pluginRouter is a router to talk with the plugin controller
type pluginRouter struct ***REMOVED***
	backend Backend
	routes  []router.Route
***REMOVED***

// NewRouter initializes a new plugin router
func NewRouter(b Backend) router.Router ***REMOVED***
	r := &pluginRouter***REMOVED***
		backend: b,
	***REMOVED***
	r.initRoutes()
	return r
***REMOVED***

// Routes returns the available routers to the plugin controller
func (r *pluginRouter) Routes() []router.Route ***REMOVED***
	return r.routes
***REMOVED***

func (r *pluginRouter) initRoutes() ***REMOVED***
	r.routes = []router.Route***REMOVED***
		router.NewGetRoute("/plugins", r.listPlugins),
		router.NewGetRoute("/plugins/***REMOVED***name:.****REMOVED***/json", r.inspectPlugin),
		router.NewGetRoute("/plugins/privileges", r.getPrivileges),
		router.NewDeleteRoute("/plugins/***REMOVED***name:.****REMOVED***", r.removePlugin),
		router.NewPostRoute("/plugins/***REMOVED***name:.****REMOVED***/enable", r.enablePlugin), // PATCH?
		router.NewPostRoute("/plugins/***REMOVED***name:.****REMOVED***/disable", r.disablePlugin),
		router.NewPostRoute("/plugins/pull", r.pullPlugin, router.WithCancel),
		router.NewPostRoute("/plugins/***REMOVED***name:.****REMOVED***/push", r.pushPlugin, router.WithCancel),
		router.NewPostRoute("/plugins/***REMOVED***name:.****REMOVED***/upgrade", r.upgradePlugin, router.WithCancel),
		router.NewPostRoute("/plugins/***REMOVED***name:.****REMOVED***/set", r.setPlugin),
		router.NewPostRoute("/plugins/create", r.createPlugin),
	***REMOVED***
***REMOVED***
