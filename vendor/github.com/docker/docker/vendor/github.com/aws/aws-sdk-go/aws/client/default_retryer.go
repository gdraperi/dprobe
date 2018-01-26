package client

import (
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws/request"
)

// DefaultRetryer implements basic retry logic using exponential backoff for
// most services. If you want to implement custom retry logic, implement the
// request.Retryer interface or create a structure type that composes this
// struct and override the specific methods. For example, to override only
// the MaxRetries method:
//
//		type retryer struct ***REMOVED***
//      client.DefaultRetryer
//***REMOVED***
//
//    // This implementation always has 100 max retries
//    func (d retryer) MaxRetries() int ***REMOVED*** return 100 ***REMOVED***
type DefaultRetryer struct ***REMOVED***
	NumMaxRetries int
***REMOVED***

// MaxRetries returns the number of maximum returns the service will use to make
// an individual API request.
func (d DefaultRetryer) MaxRetries() int ***REMOVED***
	return d.NumMaxRetries
***REMOVED***

var seededRand = rand.New(&lockedSource***REMOVED***src: rand.NewSource(time.Now().UnixNano())***REMOVED***)

// RetryRules returns the delay duration before retrying this request again
func (d DefaultRetryer) RetryRules(r *request.Request) time.Duration ***REMOVED***
	// Set the upper limit of delay in retrying at ~five minutes
	minTime := 30
	throttle := d.shouldThrottle(r)
	if throttle ***REMOVED***
		if delay, ok := getRetryDelay(r); ok ***REMOVED***
			return delay
		***REMOVED***

		minTime = 500
	***REMOVED***

	retryCount := r.RetryCount
	if throttle && retryCount > 8 ***REMOVED***
		retryCount = 8
	***REMOVED*** else if retryCount > 13 ***REMOVED***
		retryCount = 13
	***REMOVED***

	delay := (1 << uint(retryCount)) * (seededRand.Intn(minTime) + minTime)
	return time.Duration(delay) * time.Millisecond
***REMOVED***

// ShouldRetry returns true if the request should be retried.
func (d DefaultRetryer) ShouldRetry(r *request.Request) bool ***REMOVED***
	// If one of the other handlers already set the retry state
	// we don't want to override it based on the service's state
	if r.Retryable != nil ***REMOVED***
		return *r.Retryable
	***REMOVED***

	if r.HTTPResponse.StatusCode >= 500 ***REMOVED***
		return true
	***REMOVED***
	return r.IsErrorRetryable() || d.shouldThrottle(r)
***REMOVED***

// ShouldThrottle returns true if the request should be throttled.
func (d DefaultRetryer) shouldThrottle(r *request.Request) bool ***REMOVED***
	switch r.HTTPResponse.StatusCode ***REMOVED***
	case 429:
	case 502:
	case 503:
	case 504:
	default:
		return r.IsErrorThrottle()
	***REMOVED***

	return true
***REMOVED***

// This will look in the Retry-After header, RFC 7231, for how long
// it will wait before attempting another request
func getRetryDelay(r *request.Request) (time.Duration, bool) ***REMOVED***
	if !canUseRetryAfterHeader(r) ***REMOVED***
		return 0, false
	***REMOVED***

	delayStr := r.HTTPResponse.Header.Get("Retry-After")
	if len(delayStr) == 0 ***REMOVED***
		return 0, false
	***REMOVED***

	delay, err := strconv.Atoi(delayStr)
	if err != nil ***REMOVED***
		return 0, false
	***REMOVED***

	return time.Duration(delay) * time.Second, true
***REMOVED***

// Will look at the status code to see if the retry header pertains to
// the status code.
func canUseRetryAfterHeader(r *request.Request) bool ***REMOVED***
	switch r.HTTPResponse.StatusCode ***REMOVED***
	case 429:
	case 503:
	default:
		return false
	***REMOVED***

	return true
***REMOVED***

// lockedSource is a thread-safe implementation of rand.Source
type lockedSource struct ***REMOVED***
	lk  sync.Mutex
	src rand.Source
***REMOVED***

func (r *lockedSource) Int63() (n int64) ***REMOVED***
	r.lk.Lock()
	n = r.src.Int63()
	r.lk.Unlock()
	return
***REMOVED***

func (r *lockedSource) Seed(seed int64) ***REMOVED***
	r.lk.Lock()
	r.src.Seed(seed)
	r.lk.Unlock()
***REMOVED***
