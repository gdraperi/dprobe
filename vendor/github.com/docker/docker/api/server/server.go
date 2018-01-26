package server

import (
	"crypto/tls"
	"net"
	"net/http"
	"strings"

	"github.com/docker/docker/api/server/httputils"
	"github.com/docker/docker/api/server/middleware"
	"github.com/docker/docker/api/server/router"
	"github.com/docker/docker/api/server/router/debug"
	"github.com/docker/docker/dockerversion"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// versionMatcher defines a variable matcher to be parsed by the router
// when a request is about to be served.
const versionMatcher = "/v***REMOVED***version:[0-9.]+***REMOVED***"

// Config provides the configuration for the API server
type Config struct ***REMOVED***
	Logging     bool
	CorsHeaders string
	Version     string
	SocketGroup string
	TLSConfig   *tls.Config
***REMOVED***

// Server contains instance details for the server
type Server struct ***REMOVED***
	cfg           *Config
	servers       []*HTTPServer
	routers       []router.Router
	routerSwapper *routerSwapper
	middlewares   []middleware.Middleware
***REMOVED***

// New returns a new instance of the server based on the specified configuration.
// It allocates resources which will be needed for ServeAPI(ports, unix-sockets).
func New(cfg *Config) *Server ***REMOVED***
	return &Server***REMOVED***
		cfg: cfg,
	***REMOVED***
***REMOVED***

// UseMiddleware appends a new middleware to the request chain.
// This needs to be called before the API routes are configured.
func (s *Server) UseMiddleware(m middleware.Middleware) ***REMOVED***
	s.middlewares = append(s.middlewares, m)
***REMOVED***

// Accept sets a listener the server accepts connections into.
func (s *Server) Accept(addr string, listeners ...net.Listener) ***REMOVED***
	for _, listener := range listeners ***REMOVED***
		httpServer := &HTTPServer***REMOVED***
			srv: &http.Server***REMOVED***
				Addr: addr,
			***REMOVED***,
			l: listener,
		***REMOVED***
		s.servers = append(s.servers, httpServer)
	***REMOVED***
***REMOVED***

// Close closes servers and thus stop receiving requests
func (s *Server) Close() ***REMOVED***
	for _, srv := range s.servers ***REMOVED***
		if err := srv.Close(); err != nil ***REMOVED***
			logrus.Error(err)
		***REMOVED***
	***REMOVED***
***REMOVED***

// serveAPI loops through all initialized servers and spawns goroutine
// with Serve method for each. It sets createMux() as Handler also.
func (s *Server) serveAPI() error ***REMOVED***
	var chErrors = make(chan error, len(s.servers))
	for _, srv := range s.servers ***REMOVED***
		srv.srv.Handler = s.routerSwapper
		go func(srv *HTTPServer) ***REMOVED***
			var err error
			logrus.Infof("API listen on %s", srv.l.Addr())
			if err = srv.Serve(); err != nil && strings.Contains(err.Error(), "use of closed network connection") ***REMOVED***
				err = nil
			***REMOVED***
			chErrors <- err
		***REMOVED***(srv)
	***REMOVED***

	for range s.servers ***REMOVED***
		err := <-chErrors
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// HTTPServer contains an instance of http server and the listener.
// srv *http.Server, contains configuration to create an http server and a mux router with all api end points.
// l   net.Listener, is a TCP or Socket listener that dispatches incoming request to the router.
type HTTPServer struct ***REMOVED***
	srv *http.Server
	l   net.Listener
***REMOVED***

// Serve starts listening for inbound requests.
func (s *HTTPServer) Serve() error ***REMOVED***
	return s.srv.Serve(s.l)
***REMOVED***

// Close closes the HTTPServer from listening for the inbound requests.
func (s *HTTPServer) Close() error ***REMOVED***
	return s.l.Close()
***REMOVED***

func (s *Server) makeHTTPHandler(handler httputils.APIFunc) http.HandlerFunc ***REMOVED***
	return func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		// Define the context that we'll pass around to share info
		// like the docker-request-id.
		//
		// The 'context' will be used for global data that should
		// apply to all requests. Data that is specific to the
		// immediate function being called should still be passed
		// as 'args' on the function call.
		ctx := context.WithValue(context.Background(), dockerversion.UAStringKey, r.Header.Get("User-Agent"))
		handlerFunc := s.handlerWithGlobalMiddlewares(handler)

		vars := mux.Vars(r)
		if vars == nil ***REMOVED***
			vars = make(map[string]string)
		***REMOVED***

		if err := handlerFunc(ctx, w, r, vars); err != nil ***REMOVED***
			statusCode := httputils.GetHTTPErrorStatusCode(err)
			if statusCode >= 500 ***REMOVED***
				logrus.Errorf("Handler for %s %s returned error: %v", r.Method, r.URL.Path, err)
			***REMOVED***
			httputils.MakeErrorHandler(err)(w, r)
		***REMOVED***
	***REMOVED***
***REMOVED***

// InitRouter initializes the list of routers for the server.
// This method also enables the Go profiler if enableProfiler is true.
func (s *Server) InitRouter(routers ...router.Router) ***REMOVED***
	s.routers = append(s.routers, routers...)

	m := s.createMux()
	s.routerSwapper = &routerSwapper***REMOVED***
		router: m,
	***REMOVED***
***REMOVED***

type pageNotFoundError struct***REMOVED******REMOVED***

func (pageNotFoundError) Error() string ***REMOVED***
	return "page not found"
***REMOVED***

func (pageNotFoundError) NotFound() ***REMOVED******REMOVED***

// createMux initializes the main router the server uses.
func (s *Server) createMux() *mux.Router ***REMOVED***
	m := mux.NewRouter()

	logrus.Debug("Registering routers")
	for _, apiRouter := range s.routers ***REMOVED***
		for _, r := range apiRouter.Routes() ***REMOVED***
			f := s.makeHTTPHandler(r.Handler())

			logrus.Debugf("Registering %s, %s", r.Method(), r.Path())
			m.Path(versionMatcher + r.Path()).Methods(r.Method()).Handler(f)
			m.Path(r.Path()).Methods(r.Method()).Handler(f)
		***REMOVED***
	***REMOVED***

	debugRouter := debug.NewRouter()
	s.routers = append(s.routers, debugRouter)
	for _, r := range debugRouter.Routes() ***REMOVED***
		f := s.makeHTTPHandler(r.Handler())
		m.Path("/debug" + r.Path()).Handler(f)
	***REMOVED***

	notFoundHandler := httputils.MakeErrorHandler(pageNotFoundError***REMOVED******REMOVED***)
	m.HandleFunc(versionMatcher+"/***REMOVED***path:.****REMOVED***", notFoundHandler)
	m.NotFoundHandler = notFoundHandler

	return m
***REMOVED***

// Wait blocks the server goroutine until it exits.
// It sends an error message if there is any error during
// the API execution.
func (s *Server) Wait(waitChan chan error) ***REMOVED***
	if err := s.serveAPI(); err != nil ***REMOVED***
		logrus.Errorf("ServeAPI error: %v", err)
		waitChan <- err
		return
	***REMOVED***
	waitChan <- nil
***REMOVED***
