// Package query provides serialization of AWS query requests, and responses.
package query

//go:generate go run -tags codegen ../../../models/protocol_tests/generate.go ../../../models/protocol_tests/input/query.json build_test.go

import (
	"net/url"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/private/protocol/query/queryutil"
)

// BuildHandler is a named request handler for building query protocol requests
var BuildHandler = request.NamedHandler***REMOVED***Name: "awssdk.query.Build", Fn: Build***REMOVED***

// Build builds a request for an AWS Query service.
func Build(r *request.Request) ***REMOVED***
	body := url.Values***REMOVED***
		"Action":  ***REMOVED***r.Operation.Name***REMOVED***,
		"Version": ***REMOVED***r.ClientInfo.APIVersion***REMOVED***,
	***REMOVED***
	if err := queryutil.Parse(body, r.Params, false); err != nil ***REMOVED***
		r.Error = awserr.New("SerializationError", "failed encoding Query request", err)
		return
	***REMOVED***

	if r.ExpireTime == 0 ***REMOVED***
		r.HTTPRequest.Method = "POST"
		r.HTTPRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
		r.SetBufferBody([]byte(body.Encode()))
	***REMOVED*** else ***REMOVED*** // This is a pre-signed request
		r.HTTPRequest.Method = "GET"
		r.HTTPRequest.URL.RawQuery = body.Encode()
	***REMOVED***
***REMOVED***
