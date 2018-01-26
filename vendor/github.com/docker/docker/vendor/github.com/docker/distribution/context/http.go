package context

import (
	"errors"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/docker/distribution/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// Common errors used with this package.
var (
	ErrNoRequestContext        = errors.New("no http request in context")
	ErrNoResponseWriterContext = errors.New("no http response in context")
)

func parseIP(ipStr string) net.IP ***REMOVED***
	ip := net.ParseIP(ipStr)
	if ip == nil ***REMOVED***
		log.Warnf("invalid remote IP address: %q", ipStr)
	***REMOVED***
	return ip
***REMOVED***

// RemoteAddr extracts the remote address of the request, taking into
// account proxy headers.
func RemoteAddr(r *http.Request) string ***REMOVED***
	if prior := r.Header.Get("X-Forwarded-For"); prior != "" ***REMOVED***
		proxies := strings.Split(prior, ",")
		if len(proxies) > 0 ***REMOVED***
			remoteAddr := strings.Trim(proxies[0], " ")
			if parseIP(remoteAddr) != nil ***REMOVED***
				return remoteAddr
			***REMOVED***
		***REMOVED***
	***REMOVED***
	// X-Real-Ip is less supported, but worth checking in the
	// absence of X-Forwarded-For
	if realIP := r.Header.Get("X-Real-Ip"); realIP != "" ***REMOVED***
		if parseIP(realIP) != nil ***REMOVED***
			return realIP
		***REMOVED***
	***REMOVED***

	return r.RemoteAddr
***REMOVED***

// RemoteIP extracts the remote IP of the request, taking into
// account proxy headers.
func RemoteIP(r *http.Request) string ***REMOVED***
	addr := RemoteAddr(r)

	// Try parsing it as "IP:port"
	if ip, _, err := net.SplitHostPort(addr); err == nil ***REMOVED***
		return ip
	***REMOVED***

	return addr
***REMOVED***

// WithRequest places the request on the context. The context of the request
// is assigned a unique id, available at "http.request.id". The request itself
// is available at "http.request". Other common attributes are available under
// the prefix "http.request.". If a request is already present on the context,
// this method will panic.
func WithRequest(ctx Context, r *http.Request) Context ***REMOVED***
	if ctx.Value("http.request") != nil ***REMOVED***
		// NOTE(stevvooe): This needs to be considered a programming error. It
		// is unlikely that we'd want to have more than one request in
		// context.
		panic("only one request per context")
	***REMOVED***

	return &httpRequestContext***REMOVED***
		Context:   ctx,
		startedAt: time.Now(),
		id:        uuid.Generate().String(),
		r:         r,
	***REMOVED***
***REMOVED***

// GetRequest returns the http request in the given context. Returns
// ErrNoRequestContext if the context does not have an http request associated
// with it.
func GetRequest(ctx Context) (*http.Request, error) ***REMOVED***
	if r, ok := ctx.Value("http.request").(*http.Request); r != nil && ok ***REMOVED***
		return r, nil
	***REMOVED***
	return nil, ErrNoRequestContext
***REMOVED***

// GetRequestID attempts to resolve the current request id, if possible. An
// error is return if it is not available on the context.
func GetRequestID(ctx Context) string ***REMOVED***
	return GetStringValue(ctx, "http.request.id")
***REMOVED***

// WithResponseWriter returns a new context and response writer that makes
// interesting response statistics available within the context.
func WithResponseWriter(ctx Context, w http.ResponseWriter) (Context, http.ResponseWriter) ***REMOVED***
	if closeNotifier, ok := w.(http.CloseNotifier); ok ***REMOVED***
		irwCN := &instrumentedResponseWriterCN***REMOVED***
			instrumentedResponseWriter: instrumentedResponseWriter***REMOVED***
				ResponseWriter: w,
				Context:        ctx,
			***REMOVED***,
			CloseNotifier: closeNotifier,
		***REMOVED***

		return irwCN, irwCN
	***REMOVED***

	irw := instrumentedResponseWriter***REMOVED***
		ResponseWriter: w,
		Context:        ctx,
	***REMOVED***
	return &irw, &irw
***REMOVED***

// GetResponseWriter returns the http.ResponseWriter from the provided
// context. If not present, ErrNoResponseWriterContext is returned. The
// returned instance provides instrumentation in the context.
func GetResponseWriter(ctx Context) (http.ResponseWriter, error) ***REMOVED***
	v := ctx.Value("http.response")

	rw, ok := v.(http.ResponseWriter)
	if !ok || rw == nil ***REMOVED***
		return nil, ErrNoResponseWriterContext
	***REMOVED***

	return rw, nil
***REMOVED***

// getVarsFromRequest let's us change request vars implementation for testing
// and maybe future changes.
var getVarsFromRequest = mux.Vars

// WithVars extracts gorilla/mux vars and makes them available on the returned
// context. Variables are available at keys with the prefix "vars.". For
// example, if looking for the variable "name", it can be accessed as
// "vars.name". Implementations that are accessing values need not know that
// the underlying context is implemented with gorilla/mux vars.
func WithVars(ctx Context, r *http.Request) Context ***REMOVED***
	return &muxVarsContext***REMOVED***
		Context: ctx,
		vars:    getVarsFromRequest(r),
	***REMOVED***
***REMOVED***

// GetRequestLogger returns a logger that contains fields from the request in
// the current context. If the request is not available in the context, no
// fields will display. Request loggers can safely be pushed onto the context.
func GetRequestLogger(ctx Context) Logger ***REMOVED***
	return GetLogger(ctx,
		"http.request.id",
		"http.request.method",
		"http.request.host",
		"http.request.uri",
		"http.request.referer",
		"http.request.useragent",
		"http.request.remoteaddr",
		"http.request.contenttype")
***REMOVED***

// GetResponseLogger reads the current response stats and builds a logger.
// Because the values are read at call time, pushing a logger returned from
// this function on the context will lead to missing or invalid data. Only
// call this at the end of a request, after the response has been written.
func GetResponseLogger(ctx Context) Logger ***REMOVED***
	l := getLogrusLogger(ctx,
		"http.response.written",
		"http.response.status",
		"http.response.contenttype")

	duration := Since(ctx, "http.request.startedat")

	if duration > 0 ***REMOVED***
		l = l.WithField("http.response.duration", duration.String())
	***REMOVED***

	return l
***REMOVED***

// httpRequestContext makes information about a request available to context.
type httpRequestContext struct ***REMOVED***
	Context

	startedAt time.Time
	id        string
	r         *http.Request
***REMOVED***

// Value returns a keyed element of the request for use in the context. To get
// the request itself, query "request". For other components, access them as
// "request.<component>". For example, r.RequestURI
func (ctx *httpRequestContext) Value(key interface***REMOVED******REMOVED***) interface***REMOVED******REMOVED*** ***REMOVED***
	if keyStr, ok := key.(string); ok ***REMOVED***
		if keyStr == "http.request" ***REMOVED***
			return ctx.r
		***REMOVED***

		if !strings.HasPrefix(keyStr, "http.request.") ***REMOVED***
			goto fallback
		***REMOVED***

		parts := strings.Split(keyStr, ".")

		if len(parts) != 3 ***REMOVED***
			goto fallback
		***REMOVED***

		switch parts[2] ***REMOVED***
		case "uri":
			return ctx.r.RequestURI
		case "remoteaddr":
			return RemoteAddr(ctx.r)
		case "method":
			return ctx.r.Method
		case "host":
			return ctx.r.Host
		case "referer":
			referer := ctx.r.Referer()
			if referer != "" ***REMOVED***
				return referer
			***REMOVED***
		case "useragent":
			return ctx.r.UserAgent()
		case "id":
			return ctx.id
		case "startedat":
			return ctx.startedAt
		case "contenttype":
			ct := ctx.r.Header.Get("Content-Type")
			if ct != "" ***REMOVED***
				return ct
			***REMOVED***
		***REMOVED***
	***REMOVED***

fallback:
	return ctx.Context.Value(key)
***REMOVED***

type muxVarsContext struct ***REMOVED***
	Context
	vars map[string]string
***REMOVED***

func (ctx *muxVarsContext) Value(key interface***REMOVED******REMOVED***) interface***REMOVED******REMOVED*** ***REMOVED***
	if keyStr, ok := key.(string); ok ***REMOVED***
		if keyStr == "vars" ***REMOVED***
			return ctx.vars
		***REMOVED***

		if strings.HasPrefix(keyStr, "vars.") ***REMOVED***
			keyStr = strings.TrimPrefix(keyStr, "vars.")
		***REMOVED***

		if v, ok := ctx.vars[keyStr]; ok ***REMOVED***
			return v
		***REMOVED***
	***REMOVED***

	return ctx.Context.Value(key)
***REMOVED***

// instrumentedResponseWriterCN provides response writer information in a
// context. It implements http.CloseNotifier so that users can detect
// early disconnects.
type instrumentedResponseWriterCN struct ***REMOVED***
	instrumentedResponseWriter
	http.CloseNotifier
***REMOVED***

// instrumentedResponseWriter provides response writer information in a
// context. This variant is only used in the case where CloseNotifier is not
// implemented by the parent ResponseWriter.
type instrumentedResponseWriter struct ***REMOVED***
	http.ResponseWriter
	Context

	mu      sync.Mutex
	status  int
	written int64
***REMOVED***

func (irw *instrumentedResponseWriter) Write(p []byte) (n int, err error) ***REMOVED***
	n, err = irw.ResponseWriter.Write(p)

	irw.mu.Lock()
	irw.written += int64(n)

	// Guess the likely status if not set.
	if irw.status == 0 ***REMOVED***
		irw.status = http.StatusOK
	***REMOVED***

	irw.mu.Unlock()

	return
***REMOVED***

func (irw *instrumentedResponseWriter) WriteHeader(status int) ***REMOVED***
	irw.ResponseWriter.WriteHeader(status)

	irw.mu.Lock()
	irw.status = status
	irw.mu.Unlock()
***REMOVED***

func (irw *instrumentedResponseWriter) Flush() ***REMOVED***
	if flusher, ok := irw.ResponseWriter.(http.Flusher); ok ***REMOVED***
		flusher.Flush()
	***REMOVED***
***REMOVED***

func (irw *instrumentedResponseWriter) Value(key interface***REMOVED******REMOVED***) interface***REMOVED******REMOVED*** ***REMOVED***
	if keyStr, ok := key.(string); ok ***REMOVED***
		if keyStr == "http.response" ***REMOVED***
			return irw
		***REMOVED***

		if !strings.HasPrefix(keyStr, "http.response.") ***REMOVED***
			goto fallback
		***REMOVED***

		parts := strings.Split(keyStr, ".")

		if len(parts) != 3 ***REMOVED***
			goto fallback
		***REMOVED***

		irw.mu.Lock()
		defer irw.mu.Unlock()

		switch parts[2] ***REMOVED***
		case "written":
			return irw.written
		case "status":
			return irw.status
		case "contenttype":
			contentType := irw.Header().Get("Content-Type")
			if contentType != "" ***REMOVED***
				return contentType
			***REMOVED***
		***REMOVED***
	***REMOVED***

fallback:
	return irw.Context.Value(key)
***REMOVED***

func (irw *instrumentedResponseWriterCN) Value(key interface***REMOVED******REMOVED***) interface***REMOVED******REMOVED*** ***REMOVED***
	if keyStr, ok := key.(string); ok ***REMOVED***
		if keyStr == "http.response" ***REMOVED***
			return irw
		***REMOVED***
	***REMOVED***

	return irw.instrumentedResponseWriter.Value(key)
***REMOVED***
