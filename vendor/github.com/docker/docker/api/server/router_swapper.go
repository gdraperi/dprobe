package server

import (
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

// routerSwapper is an http.Handler that allows you to swap
// mux routers.
type routerSwapper struct ***REMOVED***
	mu     sync.Mutex
	router *mux.Router
***REMOVED***

// Swap changes the old router with the new one.
func (rs *routerSwapper) Swap(newRouter *mux.Router) ***REMOVED***
	rs.mu.Lock()
	rs.router = newRouter
	rs.mu.Unlock()
***REMOVED***

// ServeHTTP makes the routerSwapper to implement the http.Handler interface.
func (rs *routerSwapper) ServeHTTP(w http.ResponseWriter, r *http.Request) ***REMOVED***
	rs.mu.Lock()
	router := rs.router
	rs.mu.Unlock()
	router.ServeHTTP(w, r)
***REMOVED***
