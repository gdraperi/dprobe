// +build !go1.8

package request

import "io"

// NoBody is an io.ReadCloser with no bytes. Read always returns EOF
// and Close always returns nil. It can be used in an outgoing client
// request to explicitly signal that a request has zero bytes.
// An alternative, however, is to simply set Request.Body to nil.
//
// Copy of Go 1.8 NoBody type from net/http/http.go
type noBody struct***REMOVED******REMOVED***

func (noBody) Read([]byte) (int, error)         ***REMOVED*** return 0, io.EOF ***REMOVED***
func (noBody) Close() error                     ***REMOVED*** return nil ***REMOVED***
func (noBody) WriteTo(io.Writer) (int64, error) ***REMOVED*** return 0, nil ***REMOVED***

// NoBody is an empty reader that will trigger the Go HTTP client to not include
// and body in the HTTP request.
var NoBody = noBody***REMOVED******REMOVED***

// ResetBody rewinds the request body back to its starting position, and
// set's the HTTP Request body reference. When the body is read prior
// to being sent in the HTTP request it will need to be rewound.
//
// ResetBody will automatically be called by the SDK's build handler, but if
// the request is being used directly ResetBody must be called before the request
// is Sent.  SetStringBody, SetBufferBody, and SetReaderBody will automatically
// call ResetBody.
func (r *Request) ResetBody() ***REMOVED***
	body, err := r.getNextRequestBody()
	if err != nil ***REMOVED***
		r.Error = err
		return
	***REMOVED***

	r.HTTPRequest.Body = body
***REMOVED***
