package query

//go:generate go run -tags codegen ../../../models/protocol_tests/generate.go ../../../models/protocol_tests/output/query.json unmarshal_test.go

import (
	"encoding/xml"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/private/protocol/xml/xmlutil"
)

// UnmarshalHandler is a named request handler for unmarshaling query protocol requests
var UnmarshalHandler = request.NamedHandler***REMOVED***Name: "awssdk.query.Unmarshal", Fn: Unmarshal***REMOVED***

// UnmarshalMetaHandler is a named request handler for unmarshaling query protocol request metadata
var UnmarshalMetaHandler = request.NamedHandler***REMOVED***Name: "awssdk.query.UnmarshalMeta", Fn: UnmarshalMeta***REMOVED***

// Unmarshal unmarshals a response for an AWS Query service.
func Unmarshal(r *request.Request) ***REMOVED***
	defer r.HTTPResponse.Body.Close()
	if r.DataFilled() ***REMOVED***
		decoder := xml.NewDecoder(r.HTTPResponse.Body)
		err := xmlutil.UnmarshalXML(r.Data, decoder, r.Operation.Name+"Result")
		if err != nil ***REMOVED***
			r.Error = awserr.New("SerializationError", "failed decoding Query response", err)
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

// UnmarshalMeta unmarshals header response values for an AWS Query service.
func UnmarshalMeta(r *request.Request) ***REMOVED***
	r.RequestID = r.HTTPResponse.Header.Get("X-Amzn-Requestid")
***REMOVED***
