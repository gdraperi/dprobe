package request

import (
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
)

var timeoutErr = awserr.New(
	ErrCodeResponseTimeout,
	"read on body has reached the timeout limit",
	nil,
)

type readResult struct ***REMOVED***
	n   int
	err error
***REMOVED***

// timeoutReadCloser will handle body reads that take too long.
// We will return a ErrReadTimeout error if a timeout occurs.
type timeoutReadCloser struct ***REMOVED***
	reader   io.ReadCloser
	duration time.Duration
***REMOVED***

// Read will spin off a goroutine to call the reader's Read method. We will
// select on the timer's channel or the read's channel. Whoever completes first
// will be returned.
func (r *timeoutReadCloser) Read(b []byte) (int, error) ***REMOVED***
	timer := time.NewTimer(r.duration)
	c := make(chan readResult, 1)

	go func() ***REMOVED***
		n, err := r.reader.Read(b)
		timer.Stop()
		c <- readResult***REMOVED***n: n, err: err***REMOVED***
	***REMOVED***()

	select ***REMOVED***
	case data := <-c:
		return data.n, data.err
	case <-timer.C:
		return 0, timeoutErr
	***REMOVED***
***REMOVED***

func (r *timeoutReadCloser) Close() error ***REMOVED***
	return r.reader.Close()
***REMOVED***

const (
	// HandlerResponseTimeout is what we use to signify the name of the
	// response timeout handler.
	HandlerResponseTimeout = "ResponseTimeoutHandler"
)

// adaptToResponseTimeoutError is a handler that will replace any top level error
// to a ErrCodeResponseTimeout, if its child is that.
func adaptToResponseTimeoutError(req *Request) ***REMOVED***
	if err, ok := req.Error.(awserr.Error); ok ***REMOVED***
		aerr, ok := err.OrigErr().(awserr.Error)
		if ok && aerr.Code() == ErrCodeResponseTimeout ***REMOVED***
			req.Error = aerr
		***REMOVED***
	***REMOVED***
***REMOVED***

// WithResponseReadTimeout is a request option that will wrap the body in a timeout read closer.
// This will allow for per read timeouts. If a timeout occurred, we will return the
// ErrCodeResponseTimeout.
//
//     svc.PutObjectWithContext(ctx, params, request.WithTimeoutReadCloser(30 * time.Second)
func WithResponseReadTimeout(duration time.Duration) Option ***REMOVED***
	return func(r *Request) ***REMOVED***

		var timeoutHandler = NamedHandler***REMOVED***
			HandlerResponseTimeout,
			func(req *Request) ***REMOVED***
				req.HTTPResponse.Body = &timeoutReadCloser***REMOVED***
					reader:   req.HTTPResponse.Body,
					duration: duration,
				***REMOVED***
			***REMOVED******REMOVED***

		// remove the handler so we are not stomping over any new durations.
		r.Handlers.Send.RemoveByName(HandlerResponseTimeout)
		r.Handlers.Send.PushBackNamed(timeoutHandler)

		r.Handlers.Unmarshal.PushBack(adaptToResponseTimeoutError)
		r.Handlers.UnmarshalError.PushBack(adaptToResponseTimeoutError)
	***REMOVED***
***REMOVED***
