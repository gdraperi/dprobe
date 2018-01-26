package request

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/awsutil"
)

// WaiterResourceNotReadyErrorCode is the error code returned by a waiter when
// the waiter's max attempts have been exhausted.
const WaiterResourceNotReadyErrorCode = "ResourceNotReady"

// A WaiterOption is a function that will update the Waiter value's fields to
// configure the waiter.
type WaiterOption func(*Waiter)

// WithWaiterMaxAttempts returns the maximum number of times the waiter should
// attempt to check the resource for the target state.
func WithWaiterMaxAttempts(max int) WaiterOption ***REMOVED***
	return func(w *Waiter) ***REMOVED***
		w.MaxAttempts = max
	***REMOVED***
***REMOVED***

// WaiterDelay will return a delay the waiter should pause between attempts to
// check the resource state. The passed in attempt is the number of times the
// Waiter has checked the resource state.
//
// Attempt is the number of attempts the Waiter has made checking the resource
// state.
type WaiterDelay func(attempt int) time.Duration

// ConstantWaiterDelay returns a WaiterDelay that will always return a constant
// delay the waiter should use between attempts. It ignores the number of
// attempts made.
func ConstantWaiterDelay(delay time.Duration) WaiterDelay ***REMOVED***
	return func(attempt int) time.Duration ***REMOVED***
		return delay
	***REMOVED***
***REMOVED***

// WithWaiterDelay will set the Waiter to use the WaiterDelay passed in.
func WithWaiterDelay(delayer WaiterDelay) WaiterOption ***REMOVED***
	return func(w *Waiter) ***REMOVED***
		w.Delay = delayer
	***REMOVED***
***REMOVED***

// WithWaiterLogger returns a waiter option to set the logger a waiter
// should use to log warnings and errors to.
func WithWaiterLogger(logger aws.Logger) WaiterOption ***REMOVED***
	return func(w *Waiter) ***REMOVED***
		w.Logger = logger
	***REMOVED***
***REMOVED***

// WithWaiterRequestOptions returns a waiter option setting the request
// options for each request the waiter makes. Appends to waiter's request
// options already set.
func WithWaiterRequestOptions(opts ...Option) WaiterOption ***REMOVED***
	return func(w *Waiter) ***REMOVED***
		w.RequestOptions = append(w.RequestOptions, opts...)
	***REMOVED***
***REMOVED***

// A Waiter provides the functionality to perform a blocking call which will
// wait for a resource state to be satisfied by a service.
//
// This type should not be used directly. The API operations provided in the
// service packages prefixed with "WaitUntil" should be used instead.
type Waiter struct ***REMOVED***
	Name      string
	Acceptors []WaiterAcceptor
	Logger    aws.Logger

	MaxAttempts int
	Delay       WaiterDelay

	RequestOptions   []Option
	NewRequest       func([]Option) (*Request, error)
	SleepWithContext func(aws.Context, time.Duration) error
***REMOVED***

// ApplyOptions updates the waiter with the list of waiter options provided.
func (w *Waiter) ApplyOptions(opts ...WaiterOption) ***REMOVED***
	for _, fn := range opts ***REMOVED***
		fn(w)
	***REMOVED***
***REMOVED***

// WaiterState are states the waiter uses based on WaiterAcceptor definitions
// to identify if the resource state the waiter is waiting on has occurred.
type WaiterState int

// String returns the string representation of the waiter state.
func (s WaiterState) String() string ***REMOVED***
	switch s ***REMOVED***
	case SuccessWaiterState:
		return "success"
	case FailureWaiterState:
		return "failure"
	case RetryWaiterState:
		return "retry"
	default:
		return "unknown waiter state"
	***REMOVED***
***REMOVED***

// States the waiter acceptors will use to identify target resource states.
const (
	SuccessWaiterState WaiterState = iota // waiter successful
	FailureWaiterState                    // waiter failed
	RetryWaiterState                      // waiter needs to be retried
)

// WaiterMatchMode is the mode that the waiter will use to match the WaiterAcceptor
// definition's Expected attribute.
type WaiterMatchMode int

// Modes the waiter will use when inspecting API response to identify target
// resource states.
const (
	PathAllWaiterMatch  WaiterMatchMode = iota // match on all paths
	PathWaiterMatch                            // match on specific path
	PathAnyWaiterMatch                         // match on any path
	PathListWaiterMatch                        // match on list of paths
	StatusWaiterMatch                          // match on status code
	ErrorWaiterMatch                           // match on error
)

// String returns the string representation of the waiter match mode.
func (m WaiterMatchMode) String() string ***REMOVED***
	switch m ***REMOVED***
	case PathAllWaiterMatch:
		return "pathAll"
	case PathWaiterMatch:
		return "path"
	case PathAnyWaiterMatch:
		return "pathAny"
	case PathListWaiterMatch:
		return "pathList"
	case StatusWaiterMatch:
		return "status"
	case ErrorWaiterMatch:
		return "error"
	default:
		return "unknown waiter match mode"
	***REMOVED***
***REMOVED***

// WaitWithContext will make requests for the API operation using NewRequest to
// build API requests. The request's response will be compared against the
// Waiter's Acceptors to determine the successful state of the resource the
// waiter is inspecting.
//
// The passed in context must not be nil. If it is nil a panic will occur. The
// Context will be used to cancel the waiter's pending requests and retry delays.
// Use aws.BackgroundContext if no context is available.
//
// The waiter will continue until the target state defined by the Acceptors,
// or the max attempts expires.
//
// Will return the WaiterResourceNotReadyErrorCode error code if the waiter's
// retryer ShouldRetry returns false. This normally will happen when the max
// wait attempts expires.
func (w Waiter) WaitWithContext(ctx aws.Context) error ***REMOVED***

	for attempt := 1; ; attempt++ ***REMOVED***
		req, err := w.NewRequest(w.RequestOptions)
		if err != nil ***REMOVED***
			waiterLogf(w.Logger, "unable to create request %v", err)
			return err
		***REMOVED***
		req.Handlers.Build.PushBack(MakeAddToUserAgentFreeFormHandler("Waiter"))
		err = req.Send()

		// See if any of the acceptors match the request's response, or error
		for _, a := range w.Acceptors ***REMOVED***
			if matched, matchErr := a.match(w.Name, w.Logger, req, err); matched ***REMOVED***
				return matchErr
			***REMOVED***
		***REMOVED***

		// The Waiter should only check the resource state MaxAttempts times
		// This is here instead of in the for loop above to prevent delaying
		// unnecessary when the waiter will not retry.
		if attempt == w.MaxAttempts ***REMOVED***
			break
		***REMOVED***

		// Delay to wait before inspecting the resource again
		delay := w.Delay(attempt)
		if sleepFn := req.Config.SleepDelay; sleepFn != nil ***REMOVED***
			// Support SleepDelay for backwards compatibility and testing
			sleepFn(delay)
		***REMOVED*** else ***REMOVED***
			sleepCtxFn := w.SleepWithContext
			if sleepCtxFn == nil ***REMOVED***
				sleepCtxFn = aws.SleepWithContext
			***REMOVED***

			if err := sleepCtxFn(ctx, delay); err != nil ***REMOVED***
				return awserr.New(CanceledErrorCode, "waiter context canceled", err)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return awserr.New(WaiterResourceNotReadyErrorCode, "exceeded wait attempts", nil)
***REMOVED***

// A WaiterAcceptor provides the information needed to wait for an API operation
// to complete.
type WaiterAcceptor struct ***REMOVED***
	State    WaiterState
	Matcher  WaiterMatchMode
	Argument string
	Expected interface***REMOVED******REMOVED***
***REMOVED***

// match returns if the acceptor found a match with the passed in request
// or error. True is returned if the acceptor made a match, error is returned
// if there was an error attempting to perform the match.
func (a *WaiterAcceptor) match(name string, l aws.Logger, req *Request, err error) (bool, error) ***REMOVED***
	result := false
	var vals []interface***REMOVED******REMOVED***

	switch a.Matcher ***REMOVED***
	case PathAllWaiterMatch, PathWaiterMatch:
		// Require all matches to be equal for result to match
		vals, _ = awsutil.ValuesAtPath(req.Data, a.Argument)
		if len(vals) == 0 ***REMOVED***
			break
		***REMOVED***
		result = true
		for _, val := range vals ***REMOVED***
			if !awsutil.DeepEqual(val, a.Expected) ***REMOVED***
				result = false
				break
			***REMOVED***
		***REMOVED***
	case PathAnyWaiterMatch:
		// Only a single match needs to equal for the result to match
		vals, _ = awsutil.ValuesAtPath(req.Data, a.Argument)
		for _, val := range vals ***REMOVED***
			if awsutil.DeepEqual(val, a.Expected) ***REMOVED***
				result = true
				break
			***REMOVED***
		***REMOVED***
	case PathListWaiterMatch:
		// ignored matcher
	case StatusWaiterMatch:
		s := a.Expected.(int)
		result = s == req.HTTPResponse.StatusCode
	case ErrorWaiterMatch:
		if aerr, ok := err.(awserr.Error); ok ***REMOVED***
			result = aerr.Code() == a.Expected.(string)
		***REMOVED***
	default:
		waiterLogf(l, "WARNING: Waiter %s encountered unexpected matcher: %s",
			name, a.Matcher)
	***REMOVED***

	if !result ***REMOVED***
		// If there was no matching result found there is nothing more to do
		// for this response, retry the request.
		return false, nil
	***REMOVED***

	switch a.State ***REMOVED***
	case SuccessWaiterState:
		// waiter completed
		return true, nil
	case FailureWaiterState:
		// Waiter failure state triggered
		return true, awserr.New(WaiterResourceNotReadyErrorCode,
			"failed waiting for successful resource state", err)
	case RetryWaiterState:
		// clear the error and retry the operation
		return false, nil
	default:
		waiterLogf(l, "WARNING: Waiter %s encountered unexpected state: %s",
			name, a.State)
		return false, nil
	***REMOVED***
***REMOVED***

func waiterLogf(logger aws.Logger, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if logger != nil ***REMOVED***
		logger.Log(fmt.Sprintf(msg, args...))
	***REMOVED***
***REMOVED***
