// Package jsonrpc provides JSON RPC utilities for serialization of AWS
// requests and responses.
package jsonrpc

//go:generate go run -tags codegen ../../../models/protocol_tests/generate.go ../../../models/protocol_tests/input/json.json build_test.go
//go:generate go run -tags codegen ../../../models/protocol_tests/generate.go ../../../models/protocol_tests/output/json.json unmarshal_test.go

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/private/protocol/json/jsonutil"
	"github.com/aws/aws-sdk-go/private/protocol/rest"
)

var emptyJSON = []byte("***REMOVED******REMOVED***")

// BuildHandler is a named request handler for building jsonrpc protocol requests
var BuildHandler = request.NamedHandler***REMOVED***Name: "awssdk.jsonrpc.Build", Fn: Build***REMOVED***

// UnmarshalHandler is a named request handler for unmarshaling jsonrpc protocol requests
var UnmarshalHandler = request.NamedHandler***REMOVED***Name: "awssdk.jsonrpc.Unmarshal", Fn: Unmarshal***REMOVED***

// UnmarshalMetaHandler is a named request handler for unmarshaling jsonrpc protocol request metadata
var UnmarshalMetaHandler = request.NamedHandler***REMOVED***Name: "awssdk.jsonrpc.UnmarshalMeta", Fn: UnmarshalMeta***REMOVED***

// UnmarshalErrorHandler is a named request handler for unmarshaling jsonrpc protocol request errors
var UnmarshalErrorHandler = request.NamedHandler***REMOVED***Name: "awssdk.jsonrpc.UnmarshalError", Fn: UnmarshalError***REMOVED***

// Build builds a JSON payload for a JSON RPC request.
func Build(req *request.Request) ***REMOVED***
	var buf []byte
	var err error
	if req.ParamsFilled() ***REMOVED***
		buf, err = jsonutil.BuildJSON(req.Params)
		if err != nil ***REMOVED***
			req.Error = awserr.New("SerializationError", "failed encoding JSON RPC request", err)
			return
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		buf = emptyJSON
	***REMOVED***

	if req.ClientInfo.TargetPrefix != "" || string(buf) != "***REMOVED******REMOVED***" ***REMOVED***
		req.SetBufferBody(buf)
	***REMOVED***

	if req.ClientInfo.TargetPrefix != "" ***REMOVED***
		target := req.ClientInfo.TargetPrefix + "." + req.Operation.Name
		req.HTTPRequest.Header.Add("X-Amz-Target", target)
	***REMOVED***
	if req.ClientInfo.JSONVersion != "" ***REMOVED***
		jsonVersion := req.ClientInfo.JSONVersion
		req.HTTPRequest.Header.Add("Content-Type", "application/x-amz-json-"+jsonVersion)
	***REMOVED***
***REMOVED***

// Unmarshal unmarshals a response for a JSON RPC service.
func Unmarshal(req *request.Request) ***REMOVED***
	defer req.HTTPResponse.Body.Close()
	if req.DataFilled() ***REMOVED***
		err := jsonutil.UnmarshalJSON(req.Data, req.HTTPResponse.Body)
		if err != nil ***REMOVED***
			req.Error = awserr.New("SerializationError", "failed decoding JSON RPC response", err)
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

// UnmarshalMeta unmarshals headers from a response for a JSON RPC service.
func UnmarshalMeta(req *request.Request) ***REMOVED***
	rest.UnmarshalMeta(req)
***REMOVED***

// UnmarshalError unmarshals an error response for a JSON RPC service.
func UnmarshalError(req *request.Request) ***REMOVED***
	defer req.HTTPResponse.Body.Close()
	bodyBytes, err := ioutil.ReadAll(req.HTTPResponse.Body)
	if err != nil ***REMOVED***
		req.Error = awserr.New("SerializationError", "failed reading JSON RPC error response", err)
		return
	***REMOVED***
	if len(bodyBytes) == 0 ***REMOVED***
		req.Error = awserr.NewRequestFailure(
			awserr.New("SerializationError", req.HTTPResponse.Status, nil),
			req.HTTPResponse.StatusCode,
			"",
		)
		return
	***REMOVED***
	var jsonErr jsonErrorResponse
	if err := json.Unmarshal(bodyBytes, &jsonErr); err != nil ***REMOVED***
		req.Error = awserr.New("SerializationError", "failed decoding JSON RPC error response", err)
		return
	***REMOVED***

	codes := strings.SplitN(jsonErr.Code, "#", 2)
	req.Error = awserr.NewRequestFailure(
		awserr.New(codes[len(codes)-1], jsonErr.Message, nil),
		req.HTTPResponse.StatusCode,
		req.RequestID,
	)
***REMOVED***

type jsonErrorResponse struct ***REMOVED***
	Code    string `json:"__type"`
	Message string `json:"message"`
***REMOVED***
