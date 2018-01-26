package corehandlers

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"runtime"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
)

// Interface for matching types which also have a Len method.
type lener interface ***REMOVED***
	Len() int
***REMOVED***

// BuildContentLengthHandler builds the content length of a request based on the body,
// or will use the HTTPRequest.Header's "Content-Length" if defined. If unable
// to determine request body length and no "Content-Length" was specified it will panic.
//
// The Content-Length will only be added to the request if the length of the body
// is greater than 0. If the body is empty or the current `Content-Length`
// header is <= 0, the header will also be stripped.
var BuildContentLengthHandler = request.NamedHandler***REMOVED***Name: "core.BuildContentLengthHandler", Fn: func(r *request.Request) ***REMOVED***
	var length int64

	if slength := r.HTTPRequest.Header.Get("Content-Length"); slength != "" ***REMOVED***
		length, _ = strconv.ParseInt(slength, 10, 64)
	***REMOVED*** else ***REMOVED***
		switch body := r.Body.(type) ***REMOVED***
		case nil:
			length = 0
		case lener:
			length = int64(body.Len())
		case io.Seeker:
			r.BodyStart, _ = body.Seek(0, 1)
			end, _ := body.Seek(0, 2)
			body.Seek(r.BodyStart, 0) // make sure to seek back to original location
			length = end - r.BodyStart
		default:
			panic("Cannot get length of body, must provide `ContentLength`")
		***REMOVED***
	***REMOVED***

	if length > 0 ***REMOVED***
		r.HTTPRequest.ContentLength = length
		r.HTTPRequest.Header.Set("Content-Length", fmt.Sprintf("%d", length))
	***REMOVED*** else ***REMOVED***
		r.HTTPRequest.ContentLength = 0
		r.HTTPRequest.Header.Del("Content-Length")
	***REMOVED***
***REMOVED******REMOVED***

// SDKVersionUserAgentHandler is a request handler for adding the SDK Version to the user agent.
var SDKVersionUserAgentHandler = request.NamedHandler***REMOVED***
	Name: "core.SDKVersionUserAgentHandler",
	Fn: request.MakeAddToUserAgentHandler(aws.SDKName, aws.SDKVersion,
		runtime.Version(), runtime.GOOS, runtime.GOARCH),
***REMOVED***

var reStatusCode = regexp.MustCompile(`^(\d***REMOVED***3***REMOVED***)`)

// ValidateReqSigHandler is a request handler to ensure that the request's
// signature doesn't expire before it is sent. This can happen when a request
// is built and signed significantly before it is sent. Or significant delays
// occur when retrying requests that would cause the signature to expire.
var ValidateReqSigHandler = request.NamedHandler***REMOVED***
	Name: "core.ValidateReqSigHandler",
	Fn: func(r *request.Request) ***REMOVED***
		// Unsigned requests are not signed
		if r.Config.Credentials == credentials.AnonymousCredentials ***REMOVED***
			return
		***REMOVED***

		signedTime := r.Time
		if !r.LastSignedAt.IsZero() ***REMOVED***
			signedTime = r.LastSignedAt
		***REMOVED***

		// 10 minutes to allow for some clock skew/delays in transmission.
		// Would be improved with aws/aws-sdk-go#423
		if signedTime.Add(10 * time.Minute).After(time.Now()) ***REMOVED***
			return
		***REMOVED***

		fmt.Println("request expired, resigning")
		r.Sign()
	***REMOVED***,
***REMOVED***

// SendHandler is a request handler to send service request using HTTP client.
var SendHandler = request.NamedHandler***REMOVED***
	Name: "core.SendHandler",
	Fn: func(r *request.Request) ***REMOVED***
		sender := sendFollowRedirects
		if r.DisableFollowRedirects ***REMOVED***
			sender = sendWithoutFollowRedirects
		***REMOVED***

		if request.NoBody == r.HTTPRequest.Body ***REMOVED***
			// Strip off the request body if the NoBody reader was used as a
			// place holder for a request body. This prevents the SDK from
			// making requests with a request body when it would be invalid
			// to do so.
			//
			// Use a shallow copy of the http.Request to ensure the race condition
			// of transport on Body will not trigger
			reqOrig, reqCopy := r.HTTPRequest, *r.HTTPRequest
			reqCopy.Body = nil
			r.HTTPRequest = &reqCopy
			defer func() ***REMOVED***
				r.HTTPRequest = reqOrig
			***REMOVED***()
		***REMOVED***

		var err error
		r.HTTPResponse, err = sender(r)
		if err != nil ***REMOVED***
			handleSendError(r, err)
		***REMOVED***
	***REMOVED***,
***REMOVED***

func sendFollowRedirects(r *request.Request) (*http.Response, error) ***REMOVED***
	return r.Config.HTTPClient.Do(r.HTTPRequest)
***REMOVED***

func sendWithoutFollowRedirects(r *request.Request) (*http.Response, error) ***REMOVED***
	transport := r.Config.HTTPClient.Transport
	if transport == nil ***REMOVED***
		transport = http.DefaultTransport
	***REMOVED***

	return transport.RoundTrip(r.HTTPRequest)
***REMOVED***

func handleSendError(r *request.Request, err error) ***REMOVED***
	// Prevent leaking if an HTTPResponse was returned. Clean up
	// the body.
	if r.HTTPResponse != nil ***REMOVED***
		r.HTTPResponse.Body.Close()
	***REMOVED***
	// Capture the case where url.Error is returned for error processing
	// response. e.g. 301 without location header comes back as string
	// error and r.HTTPResponse is nil. Other URL redirect errors will
	// comeback in a similar method.
	if e, ok := err.(*url.Error); ok && e.Err != nil ***REMOVED***
		if s := reStatusCode.FindStringSubmatch(e.Err.Error()); s != nil ***REMOVED***
			code, _ := strconv.ParseInt(s[1], 10, 64)
			r.HTTPResponse = &http.Response***REMOVED***
				StatusCode: int(code),
				Status:     http.StatusText(int(code)),
				Body:       ioutil.NopCloser(bytes.NewReader([]byte***REMOVED******REMOVED***)),
			***REMOVED***
			return
		***REMOVED***
	***REMOVED***
	if r.HTTPResponse == nil ***REMOVED***
		// Add a dummy request response object to ensure the HTTPResponse
		// value is consistent.
		r.HTTPResponse = &http.Response***REMOVED***
			StatusCode: int(0),
			Status:     http.StatusText(int(0)),
			Body:       ioutil.NopCloser(bytes.NewReader([]byte***REMOVED******REMOVED***)),
		***REMOVED***
	***REMOVED***
	// Catch all other request errors.
	r.Error = awserr.New("RequestError", "send request failed", err)
	r.Retryable = aws.Bool(true) // network errors are retryable

	// Override the error with a context canceled error, if that was canceled.
	ctx := r.Context()
	select ***REMOVED***
	case <-ctx.Done():
		r.Error = awserr.New(request.CanceledErrorCode,
			"request context canceled", ctx.Err())
		r.Retryable = aws.Bool(false)
	default:
	***REMOVED***
***REMOVED***

// ValidateResponseHandler is a request handler to validate service response.
var ValidateResponseHandler = request.NamedHandler***REMOVED***Name: "core.ValidateResponseHandler", Fn: func(r *request.Request) ***REMOVED***
	if r.HTTPResponse.StatusCode == 0 || r.HTTPResponse.StatusCode >= 300 ***REMOVED***
		// this may be replaced by an UnmarshalError handler
		r.Error = awserr.New("UnknownError", "unknown error", nil)
	***REMOVED***
***REMOVED******REMOVED***

// AfterRetryHandler performs final checks to determine if the request should
// be retried and how long to delay.
var AfterRetryHandler = request.NamedHandler***REMOVED***Name: "core.AfterRetryHandler", Fn: func(r *request.Request) ***REMOVED***
	// If one of the other handlers already set the retry state
	// we don't want to override it based on the service's state
	if r.Retryable == nil || aws.BoolValue(r.Config.EnforceShouldRetryCheck) ***REMOVED***
		r.Retryable = aws.Bool(r.ShouldRetry(r))
	***REMOVED***

	if r.WillRetry() ***REMOVED***
		r.RetryDelay = r.RetryRules(r)

		if sleepFn := r.Config.SleepDelay; sleepFn != nil ***REMOVED***
			// Support SleepDelay for backwards compatibility and testing
			sleepFn(r.RetryDelay)
		***REMOVED*** else if err := aws.SleepWithContext(r.Context(), r.RetryDelay); err != nil ***REMOVED***
			r.Error = awserr.New(request.CanceledErrorCode,
				"request context canceled", err)
			r.Retryable = aws.Bool(false)
			return
		***REMOVED***

		// when the expired token exception occurs the credentials
		// need to be expired locally so that the next request to
		// get credentials will trigger a credentials refresh.
		if r.IsErrorExpired() ***REMOVED***
			r.Config.Credentials.Expire()
		***REMOVED***

		r.RetryCount++
		r.Error = nil
	***REMOVED***
***REMOVED******REMOVED***

// ValidateEndpointHandler is a request handler to validate a request had the
// appropriate Region and Endpoint set. Will set r.Error if the endpoint or
// region is not valid.
var ValidateEndpointHandler = request.NamedHandler***REMOVED***Name: "core.ValidateEndpointHandler", Fn: func(r *request.Request) ***REMOVED***
	if r.ClientInfo.SigningRegion == "" && aws.StringValue(r.Config.Region) == "" ***REMOVED***
		r.Error = aws.ErrMissingRegion
	***REMOVED*** else if r.ClientInfo.Endpoint == "" ***REMOVED***
		r.Error = aws.ErrMissingEndpoint
	***REMOVED***
***REMOVED******REMOVED***
