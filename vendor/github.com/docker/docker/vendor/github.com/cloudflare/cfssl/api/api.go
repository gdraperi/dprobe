// Package api implements an HTTP-based API and server for CFSSL.
package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/cloudflare/cfssl/errors"
	"github.com/cloudflare/cfssl/log"
)

// Handler is an interface providing a generic mechanism for handling HTTP requests.
type Handler interface ***REMOVED***
	Handle(w http.ResponseWriter, r *http.Request) error
***REMOVED***

// HTTPHandler is a wrapper that encapsulates Handler interface as http.Handler.
// HTTPHandler also enforces that the Handler only responds to requests with registered HTTP methods.
type HTTPHandler struct ***REMOVED***
	Handler          // CFSSL handler
	Methods []string // The associated HTTP methods
***REMOVED***

// HandlerFunc is similar to the http.HandlerFunc type; it serves as
// an adapter allowing the use of ordinary functions as Handlers. If
// f is a function with the appropriate signature, HandlerFunc(f) is a
// Handler object that calls f.
type HandlerFunc func(http.ResponseWriter, *http.Request) error

// Handle calls f(w, r)
func (f HandlerFunc) Handle(w http.ResponseWriter, r *http.Request) error ***REMOVED***
	w.Header().Set("Content-Type", "application/json")
	return f(w, r)
***REMOVED***

// handleError is the centralised error handling and reporting.
func handleError(w http.ResponseWriter, err error) (code int) ***REMOVED***
	if err == nil ***REMOVED***
		return http.StatusOK
	***REMOVED***
	msg := err.Error()
	httpCode := http.StatusInternalServerError

	// If it is recognized as HttpError emitted from cfssl,
	// we rewrite the status code accordingly. If it is a
	// cfssl error, set the http status to StatusBadRequest
	switch err := err.(type) ***REMOVED***
	case *errors.HTTPError:
		httpCode = err.StatusCode
		code = err.StatusCode
	case *errors.Error:
		httpCode = http.StatusBadRequest
		code = err.ErrorCode
		msg = err.Message
	***REMOVED***

	response := NewErrorResponse(msg, code)
	jsonMessage, err := json.Marshal(response)
	if err != nil ***REMOVED***
		log.Errorf("Failed to marshal JSON: %v", err)
	***REMOVED*** else ***REMOVED***
		msg = string(jsonMessage)
	***REMOVED***
	http.Error(w, msg, httpCode)
	return code
***REMOVED***

// ServeHTTP encapsulates the call to underlying Handler to handle the request
// and return the response with proper HTTP status code
func (h HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) ***REMOVED***
	var err error
	var match bool
	// Throw 405 when requested with an unsupported verb.
	for _, m := range h.Methods ***REMOVED***
		if m == r.Method ***REMOVED***
			match = true
		***REMOVED***
	***REMOVED***
	if match ***REMOVED***
		err = h.Handle(w, r)
	***REMOVED*** else ***REMOVED***
		err = errors.NewMethodNotAllowed(r.Method)
	***REMOVED***
	status := handleError(w, err)
	log.Infof("%s - \"%s %s\" %d", r.RemoteAddr, r.Method, r.URL, status)
***REMOVED***

// readRequestBlob takes a JSON-blob-encoded response body in the form
// map[string]string and returns it, the list of keywords presented,
// and any error that occurred.
func readRequestBlob(r *http.Request) (map[string]string, error) ***REMOVED***
	var blob map[string]string

	body, err := ioutil.ReadAll(r.Body)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	r.Body.Close()

	err = json.Unmarshal(body, &blob)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return blob, nil
***REMOVED***

// ProcessRequestOneOf reads a JSON blob for the request and makes
// sure it contains one of a set of keywords. For example, a request
// might have the ('foo' && 'bar') keys, OR it might have the 'baz'
// key.  In either case, we want to accept the request; however, if
// none of these sets shows up, the request is a bad request, and it
// should be returned.
func ProcessRequestOneOf(r *http.Request, keywordSets [][]string) (map[string]string, []string, error) ***REMOVED***
	blob, err := readRequestBlob(r)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	var matched []string
	for _, set := range keywordSets ***REMOVED***
		if matchKeywords(blob, set) ***REMOVED***
			if matched != nil ***REMOVED***
				return nil, nil, errors.NewBadRequestString("mismatched parameters")
			***REMOVED***
			matched = set
		***REMOVED***
	***REMOVED***
	if matched == nil ***REMOVED***
		return nil, nil, errors.NewBadRequestString("no valid parameter sets found")
	***REMOVED***
	return blob, matched, nil
***REMOVED***

// ProcessRequestFirstMatchOf reads a JSON blob for the request and returns
// the first match of a set of keywords. For example, a request
// might have one of the following combinations: (foo=1, bar=2), (foo=1), and (bar=2)
// By giving a specific ordering of those combinations, we could decide how to accept
// the request.
func ProcessRequestFirstMatchOf(r *http.Request, keywordSets [][]string) (map[string]string, []string, error) ***REMOVED***
	blob, err := readRequestBlob(r)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	for _, set := range keywordSets ***REMOVED***
		if matchKeywords(blob, set) ***REMOVED***
			return blob, set, nil
		***REMOVED***
	***REMOVED***
	return nil, nil, errors.NewBadRequestString("no valid parameter sets found")
***REMOVED***

func matchKeywords(blob map[string]string, keywords []string) bool ***REMOVED***
	for _, keyword := range keywords ***REMOVED***
		if _, ok := blob[keyword]; !ok ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

// ResponseMessage implements the standard for response errors and
// messages. A message has a code and a string message.
type ResponseMessage struct ***REMOVED***
	Code    int    `json:"code"`
	Message string `json:"message"`
***REMOVED***

// Response implements the CloudFlare standard for API
// responses.
type Response struct ***REMOVED***
	Success  bool              `json:"success"`
	Result   interface***REMOVED******REMOVED***       `json:"result"`
	Errors   []ResponseMessage `json:"errors"`
	Messages []ResponseMessage `json:"messages"`
***REMOVED***

// NewSuccessResponse is a shortcut for creating new successul API
// responses.
func NewSuccessResponse(result interface***REMOVED******REMOVED***) Response ***REMOVED***
	return Response***REMOVED***
		Success:  true,
		Result:   result,
		Errors:   []ResponseMessage***REMOVED******REMOVED***,
		Messages: []ResponseMessage***REMOVED******REMOVED***,
	***REMOVED***
***REMOVED***

// NewSuccessResponseWithMessage is a shortcut for creating new successul API
// responses that includes a message.
func NewSuccessResponseWithMessage(result interface***REMOVED******REMOVED***, message string, code int) Response ***REMOVED***
	return Response***REMOVED***
		Success:  true,
		Result:   result,
		Errors:   []ResponseMessage***REMOVED******REMOVED***,
		Messages: []ResponseMessage***REMOVED******REMOVED***code, message***REMOVED******REMOVED***,
	***REMOVED***
***REMOVED***

// NewErrorResponse is a shortcut for creating an error response for a
// single error.
func NewErrorResponse(message string, code int) Response ***REMOVED***
	return Response***REMOVED***
		Success:  false,
		Result:   nil,
		Errors:   []ResponseMessage***REMOVED******REMOVED***code, message***REMOVED******REMOVED***,
		Messages: []ResponseMessage***REMOVED******REMOVED***,
	***REMOVED***
***REMOVED***

// SendResponse builds a response from the result, sets the JSON
// header, and writes to the http.ResponseWriter.
func SendResponse(w http.ResponseWriter, result interface***REMOVED******REMOVED***) error ***REMOVED***
	response := NewSuccessResponse(result)
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	err := enc.Encode(response)
	return err
***REMOVED***

// SendResponseWithMessage builds a response from the result and the
// provided message, sets the JSON header, and writes to the
// http.ResponseWriter.
func SendResponseWithMessage(w http.ResponseWriter, result interface***REMOVED******REMOVED***, message string, code int) error ***REMOVED***
	response := NewSuccessResponseWithMessage(result, message, code)
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	err := enc.Encode(response)
	return err
***REMOVED***
