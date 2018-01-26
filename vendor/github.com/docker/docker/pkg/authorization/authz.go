package authorization

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/docker/docker/pkg/ioutils"
	"github.com/sirupsen/logrus"
)

const maxBodySize = 1048576 // 1MB

// NewCtx creates new authZ context, it is used to store authorization information related to a specific docker
// REST http session
// A context provides two method:
// Authenticate Request:
// Call authZ plugins with current REST request and AuthN response
// Request contains full HTTP packet sent to the docker daemon
// https://docs.docker.com/engine/reference/api/
//
// Authenticate Response:
// Call authZ plugins with full info about current REST request, REST response and AuthN response
// The response from this method may contains content that overrides the daemon response
// This allows authZ plugins to filter privileged content
//
// If multiple authZ plugins are specified, the block/allow decision is based on ANDing all plugin results
// For response manipulation, the response from each plugin is piped between plugins. Plugin execution order
// is determined according to daemon parameters
func NewCtx(authZPlugins []Plugin, user, userAuthNMethod, requestMethod, requestURI string) *Ctx ***REMOVED***
	return &Ctx***REMOVED***
		plugins:         authZPlugins,
		user:            user,
		userAuthNMethod: userAuthNMethod,
		requestMethod:   requestMethod,
		requestURI:      requestURI,
	***REMOVED***
***REMOVED***

// Ctx stores a single request-response interaction context
type Ctx struct ***REMOVED***
	user            string
	userAuthNMethod string
	requestMethod   string
	requestURI      string
	plugins         []Plugin
	// authReq stores the cached request object for the current transaction
	authReq *Request
***REMOVED***

// AuthZRequest authorized the request to the docker daemon using authZ plugins
func (ctx *Ctx) AuthZRequest(w http.ResponseWriter, r *http.Request) error ***REMOVED***
	var body []byte
	if sendBody(ctx.requestURI, r.Header) && r.ContentLength > 0 && r.ContentLength < maxBodySize ***REMOVED***
		var err error
		body, r.Body, err = drainBody(r.Body)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	var h bytes.Buffer
	if err := r.Header.Write(&h); err != nil ***REMOVED***
		return err
	***REMOVED***

	ctx.authReq = &Request***REMOVED***
		User:            ctx.user,
		UserAuthNMethod: ctx.userAuthNMethod,
		RequestMethod:   ctx.requestMethod,
		RequestURI:      ctx.requestURI,
		RequestBody:     body,
		RequestHeaders:  headers(r.Header),
	***REMOVED***

	if r.TLS != nil ***REMOVED***
		for _, c := range r.TLS.PeerCertificates ***REMOVED***
			pc := PeerCertificate(*c)
			ctx.authReq.RequestPeerCertificates = append(ctx.authReq.RequestPeerCertificates, &pc)
		***REMOVED***
	***REMOVED***

	for _, plugin := range ctx.plugins ***REMOVED***
		logrus.Debugf("AuthZ request using plugin %s", plugin.Name())

		authRes, err := plugin.AuthZRequest(ctx.authReq)
		if err != nil ***REMOVED***
			return fmt.Errorf("plugin %s failed with error: %s", plugin.Name(), err)
		***REMOVED***

		if !authRes.Allow ***REMOVED***
			return newAuthorizationError(plugin.Name(), authRes.Msg)
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// AuthZResponse authorized and manipulates the response from docker daemon using authZ plugins
func (ctx *Ctx) AuthZResponse(rm ResponseModifier, r *http.Request) error ***REMOVED***
	ctx.authReq.ResponseStatusCode = rm.StatusCode()
	ctx.authReq.ResponseHeaders = headers(rm.Header())

	if sendBody(ctx.requestURI, rm.Header()) ***REMOVED***
		ctx.authReq.ResponseBody = rm.RawBody()
	***REMOVED***

	for _, plugin := range ctx.plugins ***REMOVED***
		logrus.Debugf("AuthZ response using plugin %s", plugin.Name())

		authRes, err := plugin.AuthZResponse(ctx.authReq)
		if err != nil ***REMOVED***
			return fmt.Errorf("plugin %s failed with error: %s", plugin.Name(), err)
		***REMOVED***

		if !authRes.Allow ***REMOVED***
			return newAuthorizationError(plugin.Name(), authRes.Msg)
		***REMOVED***
	***REMOVED***

	rm.FlushAll()

	return nil
***REMOVED***

// drainBody dump the body (if its length is less than 1MB) without modifying the request state
func drainBody(body io.ReadCloser) ([]byte, io.ReadCloser, error) ***REMOVED***
	bufReader := bufio.NewReaderSize(body, maxBodySize)
	newBody := ioutils.NewReadCloserWrapper(bufReader, func() error ***REMOVED*** return body.Close() ***REMOVED***)

	data, err := bufReader.Peek(maxBodySize)
	// Body size exceeds max body size
	if err == nil ***REMOVED***
		logrus.Warnf("Request body is larger than: '%d' skipping body", maxBodySize)
		return nil, newBody, nil
	***REMOVED***
	// Body size is less than maximum size
	if err == io.EOF ***REMOVED***
		return data, newBody, nil
	***REMOVED***
	// Unknown error
	return nil, newBody, err
***REMOVED***

// sendBody returns true when request/response body should be sent to AuthZPlugin
func sendBody(url string, header http.Header) bool ***REMOVED***
	// Skip body for auth endpoint
	if strings.HasSuffix(url, "/auth") ***REMOVED***
		return false
	***REMOVED***

	// body is sent only for text or json messages
	return header.Get("Content-Type") == "application/json"
***REMOVED***

// headers returns flatten version of the http headers excluding authorization
func headers(header http.Header) map[string]string ***REMOVED***
	v := make(map[string]string)
	for k, values := range header ***REMOVED***
		// Skip authorization headers
		if strings.EqualFold(k, "Authorization") || strings.EqualFold(k, "X-Registry-Config") || strings.EqualFold(k, "X-Registry-Auth") ***REMOVED***
			continue
		***REMOVED***
		for _, val := range values ***REMOVED***
			v[k] = val
		***REMOVED***
	***REMOVED***
	return v
***REMOVED***

// authorizationError represents an authorization deny error
type authorizationError struct ***REMOVED***
	error
***REMOVED***

func (authorizationError) Forbidden() ***REMOVED******REMOVED***

func newAuthorizationError(plugin, msg string) authorizationError ***REMOVED***
	return authorizationError***REMOVED***error: fmt.Errorf("authorization denied by plugin %s: %s", plugin, msg)***REMOVED***
***REMOVED***
