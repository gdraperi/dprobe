package authorization

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/sirupsen/logrus"
)

// ResponseModifier allows authorization plugins to read and modify the content of the http.response
type ResponseModifier interface ***REMOVED***
	http.ResponseWriter
	http.Flusher
	http.CloseNotifier

	// RawBody returns the current http content
	RawBody() []byte

	// RawHeaders returns the current content of the http headers
	RawHeaders() ([]byte, error)

	// StatusCode returns the current status code
	StatusCode() int

	// OverrideBody replaces the body of the HTTP reply
	OverrideBody(b []byte)

	// OverrideHeader replaces the headers of the HTTP reply
	OverrideHeader(b []byte) error

	// OverrideStatusCode replaces the status code of the HTTP reply
	OverrideStatusCode(statusCode int)

	// FlushAll flushes all data to the HTTP response
	FlushAll() error

	// Hijacked indicates the response has been hijacked by the Docker daemon
	Hijacked() bool
***REMOVED***

// NewResponseModifier creates a wrapper to an http.ResponseWriter to allow inspecting and modifying the content
func NewResponseModifier(rw http.ResponseWriter) ResponseModifier ***REMOVED***
	return &responseModifier***REMOVED***rw: rw, header: make(http.Header)***REMOVED***
***REMOVED***

// responseModifier is used as an adapter to http.ResponseWriter in order to manipulate and explore
// the http request/response from docker daemon
type responseModifier struct ***REMOVED***
	// The original response writer
	rw http.ResponseWriter
	// body holds the response body
	body []byte
	// header holds the response header
	header http.Header
	// statusCode holds the response status code
	statusCode int
	// hijacked indicates the request has been hijacked
	hijacked bool
***REMOVED***

func (rm *responseModifier) Hijacked() bool ***REMOVED***
	return rm.hijacked
***REMOVED***

// WriterHeader stores the http status code
func (rm *responseModifier) WriteHeader(s int) ***REMOVED***

	// Use original request if hijacked
	if rm.hijacked ***REMOVED***
		rm.rw.WriteHeader(s)
		return
	***REMOVED***

	rm.statusCode = s
***REMOVED***

// Header returns the internal http header
func (rm *responseModifier) Header() http.Header ***REMOVED***

	// Use original header if hijacked
	if rm.hijacked ***REMOVED***
		return rm.rw.Header()
	***REMOVED***

	return rm.header
***REMOVED***

// StatusCode returns the http status code
func (rm *responseModifier) StatusCode() int ***REMOVED***
	return rm.statusCode
***REMOVED***

// OverrideBody replaces the body of the HTTP response
func (rm *responseModifier) OverrideBody(b []byte) ***REMOVED***
	rm.body = b
***REMOVED***

// OverrideStatusCode replaces the status code of the HTTP response
func (rm *responseModifier) OverrideStatusCode(statusCode int) ***REMOVED***
	rm.statusCode = statusCode
***REMOVED***

// OverrideHeader replaces the headers of the HTTP response
func (rm *responseModifier) OverrideHeader(b []byte) error ***REMOVED***
	header := http.Header***REMOVED******REMOVED***
	if err := json.Unmarshal(b, &header); err != nil ***REMOVED***
		return err
	***REMOVED***
	rm.header = header
	return nil
***REMOVED***

// Write stores the byte array inside content
func (rm *responseModifier) Write(b []byte) (int, error) ***REMOVED***

	if rm.hijacked ***REMOVED***
		return rm.rw.Write(b)
	***REMOVED***

	rm.body = append(rm.body, b...)
	return len(b), nil
***REMOVED***

// Body returns the response body
func (rm *responseModifier) RawBody() []byte ***REMOVED***
	return rm.body
***REMOVED***

func (rm *responseModifier) RawHeaders() ([]byte, error) ***REMOVED***
	var b bytes.Buffer
	if err := rm.header.Write(&b); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return b.Bytes(), nil
***REMOVED***

// Hijack returns the internal connection of the wrapped http.ResponseWriter
func (rm *responseModifier) Hijack() (net.Conn, *bufio.ReadWriter, error) ***REMOVED***

	rm.hijacked = true
	rm.FlushAll()

	hijacker, ok := rm.rw.(http.Hijacker)
	if !ok ***REMOVED***
		return nil, nil, fmt.Errorf("Internal response writer doesn't support the Hijacker interface")
	***REMOVED***
	return hijacker.Hijack()
***REMOVED***

// CloseNotify uses the internal close notify API of the wrapped http.ResponseWriter
func (rm *responseModifier) CloseNotify() <-chan bool ***REMOVED***
	closeNotifier, ok := rm.rw.(http.CloseNotifier)
	if !ok ***REMOVED***
		logrus.Error("Internal response writer doesn't support the CloseNotifier interface")
		return nil
	***REMOVED***
	return closeNotifier.CloseNotify()
***REMOVED***

// Flush uses the internal flush API of the wrapped http.ResponseWriter
func (rm *responseModifier) Flush() ***REMOVED***
	flusher, ok := rm.rw.(http.Flusher)
	if !ok ***REMOVED***
		logrus.Error("Internal response writer doesn't support the Flusher interface")
		return
	***REMOVED***

	rm.FlushAll()
	flusher.Flush()
***REMOVED***

// FlushAll flushes all data to the HTTP response
func (rm *responseModifier) FlushAll() error ***REMOVED***
	// Copy the header
	for k, vv := range rm.header ***REMOVED***
		for _, v := range vv ***REMOVED***
			rm.rw.Header().Add(k, v)
		***REMOVED***
	***REMOVED***

	// Copy the status code
	// Also WriteHeader needs to be done after all the headers
	// have been copied (above).
	if rm.statusCode > 0 ***REMOVED***
		rm.rw.WriteHeader(rm.statusCode)
	***REMOVED***

	var err error
	if len(rm.body) > 0 ***REMOVED***
		// Write body
		_, err = rm.rw.Write(rm.body)
	***REMOVED***

	// Clean previous data
	rm.body = nil
	rm.statusCode = 0
	rm.header = http.Header***REMOVED******REMOVED***
	return err
***REMOVED***
