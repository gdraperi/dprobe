package query

import (
	"encoding/xml"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
)

type xmlErrorResponse struct ***REMOVED***
	XMLName   xml.Name `xml:"ErrorResponse"`
	Code      string   `xml:"Error>Code"`
	Message   string   `xml:"Error>Message"`
	RequestID string   `xml:"RequestId"`
***REMOVED***

type xmlServiceUnavailableResponse struct ***REMOVED***
	XMLName xml.Name `xml:"ServiceUnavailableException"`
***REMOVED***

// UnmarshalErrorHandler is a name request handler to unmarshal request errors
var UnmarshalErrorHandler = request.NamedHandler***REMOVED***Name: "awssdk.query.UnmarshalError", Fn: UnmarshalError***REMOVED***

// UnmarshalError unmarshals an error response for an AWS Query service.
func UnmarshalError(r *request.Request) ***REMOVED***
	defer r.HTTPResponse.Body.Close()

	bodyBytes, err := ioutil.ReadAll(r.HTTPResponse.Body)
	if err != nil ***REMOVED***
		r.Error = awserr.New("SerializationError", "failed to read from query HTTP response body", err)
		return
	***REMOVED***

	// First check for specific error
	resp := xmlErrorResponse***REMOVED******REMOVED***
	decodeErr := xml.Unmarshal(bodyBytes, &resp)
	if decodeErr == nil ***REMOVED***
		reqID := resp.RequestID
		if reqID == "" ***REMOVED***
			reqID = r.RequestID
		***REMOVED***
		r.Error = awserr.NewRequestFailure(
			awserr.New(resp.Code, resp.Message, nil),
			r.HTTPResponse.StatusCode,
			reqID,
		)
		return
	***REMOVED***

	// Check for unhandled error
	servUnavailResp := xmlServiceUnavailableResponse***REMOVED******REMOVED***
	unavailErr := xml.Unmarshal(bodyBytes, &servUnavailResp)
	if unavailErr == nil ***REMOVED***
		r.Error = awserr.NewRequestFailure(
			awserr.New("ServiceUnavailableException", "service is unavailable", nil),
			r.HTTPResponse.StatusCode,
			r.RequestID,
		)
		return
	***REMOVED***

	// Failed to retrieve any error message from the response body
	r.Error = awserr.New("SerializationError",
		"failed to decode query XML error response", decodeErr)
***REMOVED***
