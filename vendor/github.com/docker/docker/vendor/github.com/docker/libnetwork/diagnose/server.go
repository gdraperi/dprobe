package diagnose

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"

	stackdump "github.com/docker/docker/pkg/signal"
	"github.com/docker/libnetwork/common"
	"github.com/sirupsen/logrus"
)

// HTTPHandlerFunc TODO
type HTTPHandlerFunc func(interface***REMOVED******REMOVED***, http.ResponseWriter, *http.Request)

type httpHandlerCustom struct ***REMOVED***
	ctx interface***REMOVED******REMOVED***
	F   func(interface***REMOVED******REMOVED***, http.ResponseWriter, *http.Request)
***REMOVED***

// ServeHTTP TODO
func (h httpHandlerCustom) ServeHTTP(w http.ResponseWriter, r *http.Request) ***REMOVED***
	h.F(h.ctx, w, r)
***REMOVED***

var diagPaths2Func = map[string]HTTPHandlerFunc***REMOVED***
	"/":          notImplemented,
	"/help":      help,
	"/ready":     ready,
	"/stackdump": stackTrace,
***REMOVED***

// Server when the debug is enabled exposes a
// This data structure is protected by the Agent mutex so does not require and additional mutex here
type Server struct ***REMOVED***
	enable            int32
	srv               *http.Server
	port              int
	mux               *http.ServeMux
	registeredHanders map[string]bool
	sync.Mutex
***REMOVED***

// New creates a new diagnose server
func New() *Server ***REMOVED***
	return &Server***REMOVED***
		registeredHanders: make(map[string]bool),
	***REMOVED***
***REMOVED***

// Init initialize the mux for the http handling and register the base hooks
func (s *Server) Init() ***REMOVED***
	s.mux = http.NewServeMux()

	// Register local handlers
	s.RegisterHandler(s, diagPaths2Func)
***REMOVED***

// RegisterHandler allows to register new handlers to the mux and to a specific path
func (s *Server) RegisterHandler(ctx interface***REMOVED******REMOVED***, hdlrs map[string]HTTPHandlerFunc) ***REMOVED***
	s.Lock()
	defer s.Unlock()
	for path, fun := range hdlrs ***REMOVED***
		if _, ok := s.registeredHanders[path]; ok ***REMOVED***
			continue
		***REMOVED***
		s.mux.Handle(path, httpHandlerCustom***REMOVED***ctx, fun***REMOVED***)
		s.registeredHanders[path] = true
	***REMOVED***
***REMOVED***

// ServeHTTP this is the method called bu the ListenAndServe, and is needed to allow us to
// use our custom mux
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) ***REMOVED***
	s.mux.ServeHTTP(w, r)
***REMOVED***

// EnableDebug opens a TCP socket to debug the passed network DB
func (s *Server) EnableDebug(ip string, port int) ***REMOVED***
	s.Lock()
	defer s.Unlock()

	s.port = port

	if s.enable == 1 ***REMOVED***
		logrus.Info("The server is already up and running")
		return
	***REMOVED***

	logrus.Infof("Starting the diagnose server listening on %d for commands", port)
	srv := &http.Server***REMOVED***Addr: fmt.Sprintf("%s:%d", ip, port), Handler: s***REMOVED***
	s.srv = srv
	s.enable = 1
	go func(n *Server) ***REMOVED***
		// Ingore ErrServerClosed that is returned on the Shutdown call
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed ***REMOVED***
			logrus.Errorf("ListenAndServe error: %s", err)
			atomic.SwapInt32(&n.enable, 0)
		***REMOVED***
	***REMOVED***(s)
***REMOVED***

// DisableDebug stop the dubug and closes the tcp socket
func (s *Server) DisableDebug() ***REMOVED***
	s.Lock()
	defer s.Unlock()

	s.srv.Shutdown(context.Background())
	s.srv = nil
	s.enable = 0
	logrus.Info("Disabling the diagnose server")
***REMOVED***

// IsDebugEnable returns true when the debug is enabled
func (s *Server) IsDebugEnable() bool ***REMOVED***
	s.Lock()
	defer s.Unlock()
	return s.enable == 1
***REMOVED***

func notImplemented(ctx interface***REMOVED******REMOVED***, w http.ResponseWriter, r *http.Request) ***REMOVED***
	r.ParseForm()
	_, json := ParseHTTPFormOptions(r)
	rsp := WrongCommand("not implemented", fmt.Sprintf("URL path: %s no method implemented check /help\n", r.URL.Path))

	// audit logs
	log := logrus.WithFields(logrus.Fields***REMOVED***"component": "diagnose", "remoteIP": r.RemoteAddr, "method": common.CallerName(0), "url": r.URL.String()***REMOVED***)
	log.Info("command not implemented done")

	HTTPReply(w, rsp, json)
***REMOVED***

func help(ctx interface***REMOVED******REMOVED***, w http.ResponseWriter, r *http.Request) ***REMOVED***
	r.ParseForm()
	_, json := ParseHTTPFormOptions(r)

	// audit logs
	log := logrus.WithFields(logrus.Fields***REMOVED***"component": "diagnose", "remoteIP": r.RemoteAddr, "method": common.CallerName(0), "url": r.URL.String()***REMOVED***)
	log.Info("help done")

	n, ok := ctx.(*Server)
	var result string
	if ok ***REMOVED***
		for path := range n.registeredHanders ***REMOVED***
			result += fmt.Sprintf("%s\n", path)
		***REMOVED***
		HTTPReply(w, CommandSucceed(&StringCmd***REMOVED***Info: result***REMOVED***), json)
	***REMOVED***
***REMOVED***

func ready(ctx interface***REMOVED******REMOVED***, w http.ResponseWriter, r *http.Request) ***REMOVED***
	r.ParseForm()
	_, json := ParseHTTPFormOptions(r)

	// audit logs
	log := logrus.WithFields(logrus.Fields***REMOVED***"component": "diagnose", "remoteIP": r.RemoteAddr, "method": common.CallerName(0), "url": r.URL.String()***REMOVED***)
	log.Info("ready done")
	HTTPReply(w, CommandSucceed(&StringCmd***REMOVED***Info: "OK"***REMOVED***), json)
***REMOVED***

func stackTrace(ctx interface***REMOVED******REMOVED***, w http.ResponseWriter, r *http.Request) ***REMOVED***
	r.ParseForm()
	_, json := ParseHTTPFormOptions(r)

	// audit logs
	log := logrus.WithFields(logrus.Fields***REMOVED***"component": "diagnose", "remoteIP": r.RemoteAddr, "method": common.CallerName(0), "url": r.URL.String()***REMOVED***)
	log.Info("stack trace")

	path, err := stackdump.DumpStacks("/tmp/")
	if err != nil ***REMOVED***
		log.WithError(err).Error("failed to write goroutines dump")
		HTTPReply(w, FailCommand(err), json)
	***REMOVED*** else ***REMOVED***
		log.Info("stack trace done")
		HTTPReply(w, CommandSucceed(&StringCmd***REMOVED***Info: fmt.Sprintf("goroutine stacks written to %s", path)***REMOVED***), json)
	***REMOVED***
***REMOVED***

// DebugHTTPForm helper to print the form url parameters
func DebugHTTPForm(r *http.Request) ***REMOVED***
	for k, v := range r.Form ***REMOVED***
		logrus.Debugf("Form[%q] = %q\n", k, v)
	***REMOVED***
***REMOVED***

// JSONOutput contains details on JSON output printing
type JSONOutput struct ***REMOVED***
	enable      bool
	prettyPrint bool
***REMOVED***

// ParseHTTPFormOptions easily parse the JSON printing options
func ParseHTTPFormOptions(r *http.Request) (bool, *JSONOutput) ***REMOVED***
	_, unsafe := r.Form["unsafe"]
	v, json := r.Form["json"]
	var pretty bool
	if len(v) > 0 ***REMOVED***
		pretty = v[0] == "pretty"
	***REMOVED***
	return unsafe, &JSONOutput***REMOVED***enable: json, prettyPrint: pretty***REMOVED***
***REMOVED***

// HTTPReply helper function that takes care of sending the message out
func HTTPReply(w http.ResponseWriter, r *HTTPResult, j *JSONOutput) (int, error) ***REMOVED***
	var response []byte
	if j.enable ***REMOVED***
		w.Header().Set("Content-Type", "application/json")
		var err error
		if j.prettyPrint ***REMOVED***
			response, err = json.MarshalIndent(r, "", "  ")
			if err != nil ***REMOVED***
				response, _ = json.MarshalIndent(FailCommand(err), "", "  ")
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			response, err = json.Marshal(r)
			if err != nil ***REMOVED***
				response, _ = json.Marshal(FailCommand(err))
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		response = []byte(r.String())
	***REMOVED***
	return fmt.Fprint(w, string(response))
***REMOVED***
