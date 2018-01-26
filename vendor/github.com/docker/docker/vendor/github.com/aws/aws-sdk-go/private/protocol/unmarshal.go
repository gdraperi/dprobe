package protocol

import (
	"io"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws/request"
)

// UnmarshalDiscardBodyHandler is a named request handler to empty and close a response's body
var UnmarshalDiscardBodyHandler = request.NamedHandler***REMOVED***Name: "awssdk.shared.UnmarshalDiscardBody", Fn: UnmarshalDiscardBody***REMOVED***

// UnmarshalDiscardBody is a request handler to empty a response's body and closing it.
func UnmarshalDiscardBody(r *request.Request) ***REMOVED***
	if r.HTTPResponse == nil || r.HTTPResponse.Body == nil ***REMOVED***
		return
	***REMOVED***

	io.Copy(ioutil.Discard, r.HTTPResponse.Body)
	r.HTTPResponse.Body.Close()
***REMOVED***
