package request

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
)

// Retryer is an interface to control retry logic for a given service.
// The default implementation used by most services is the client.DefaultRetryer
// structure, which contains basic retry logic using exponential backoff.
type Retryer interface ***REMOVED***
	RetryRules(*Request) time.Duration
	ShouldRetry(*Request) bool
	MaxRetries() int
***REMOVED***

// WithRetryer sets a config Retryer value to the given Config returning it
// for chaining.
func WithRetryer(cfg *aws.Config, retryer Retryer) *aws.Config ***REMOVED***
	cfg.Retryer = retryer
	return cfg
***REMOVED***

// retryableCodes is a collection of service response codes which are retry-able
// without any further action.
var retryableCodes = map[string]struct***REMOVED******REMOVED******REMOVED***
	"RequestError":            ***REMOVED******REMOVED***,
	"RequestTimeout":          ***REMOVED******REMOVED***,
	ErrCodeResponseTimeout:    ***REMOVED******REMOVED***,
	"RequestTimeoutException": ***REMOVED******REMOVED***, // Glacier's flavor of RequestTimeout
***REMOVED***

var throttleCodes = map[string]struct***REMOVED******REMOVED******REMOVED***
	"ProvisionedThroughputExceededException": ***REMOVED******REMOVED***,
	"Throttling":                             ***REMOVED******REMOVED***,
	"ThrottlingException":                    ***REMOVED******REMOVED***,
	"RequestLimitExceeded":                   ***REMOVED******REMOVED***,
	"RequestThrottled":                       ***REMOVED******REMOVED***,
	"TooManyRequestsException":               ***REMOVED******REMOVED***, // Lambda functions
	"PriorRequestNotComplete":                ***REMOVED******REMOVED***, // Route53
***REMOVED***

// credsExpiredCodes is a collection of error codes which signify the credentials
// need to be refreshed. Expired tokens require refreshing of credentials, and
// resigning before the request can be retried.
var credsExpiredCodes = map[string]struct***REMOVED******REMOVED******REMOVED***
	"ExpiredToken":          ***REMOVED******REMOVED***,
	"ExpiredTokenException": ***REMOVED******REMOVED***,
	"RequestExpired":        ***REMOVED******REMOVED***, // EC2 Only
***REMOVED***

func isCodeThrottle(code string) bool ***REMOVED***
	_, ok := throttleCodes[code]
	return ok
***REMOVED***

func isCodeRetryable(code string) bool ***REMOVED***
	if _, ok := retryableCodes[code]; ok ***REMOVED***
		return true
	***REMOVED***

	return isCodeExpiredCreds(code)
***REMOVED***

func isCodeExpiredCreds(code string) bool ***REMOVED***
	_, ok := credsExpiredCodes[code]
	return ok
***REMOVED***

var validParentCodes = map[string]struct***REMOVED******REMOVED******REMOVED***
	ErrCodeSerialization: ***REMOVED******REMOVED***,
	ErrCodeRead:          ***REMOVED******REMOVED***,
***REMOVED***

type temporaryError interface ***REMOVED***
	Temporary() bool
***REMOVED***

func isNestedErrorRetryable(parentErr awserr.Error) bool ***REMOVED***
	if parentErr == nil ***REMOVED***
		return false
	***REMOVED***

	if _, ok := validParentCodes[parentErr.Code()]; !ok ***REMOVED***
		return false
	***REMOVED***

	err := parentErr.OrigErr()
	if err == nil ***REMOVED***
		return false
	***REMOVED***

	if aerr, ok := err.(awserr.Error); ok ***REMOVED***
		return isCodeRetryable(aerr.Code())
	***REMOVED***

	if t, ok := err.(temporaryError); ok ***REMOVED***
		return t.Temporary()
	***REMOVED***

	return isErrConnectionReset(err)
***REMOVED***

// IsErrorRetryable returns whether the error is retryable, based on its Code.
// Returns false if error is nil.
func IsErrorRetryable(err error) bool ***REMOVED***
	if err != nil ***REMOVED***
		if aerr, ok := err.(awserr.Error); ok ***REMOVED***
			return isCodeRetryable(aerr.Code()) || isNestedErrorRetryable(aerr)
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// IsErrorThrottle returns whether the error is to be throttled based on its code.
// Returns false if error is nil.
func IsErrorThrottle(err error) bool ***REMOVED***
	if err != nil ***REMOVED***
		if aerr, ok := err.(awserr.Error); ok ***REMOVED***
			return isCodeThrottle(aerr.Code())
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// IsErrorExpiredCreds returns whether the error code is a credential expiry error.
// Returns false if error is nil.
func IsErrorExpiredCreds(err error) bool ***REMOVED***
	if err != nil ***REMOVED***
		if aerr, ok := err.(awserr.Error); ok ***REMOVED***
			return isCodeExpiredCreds(aerr.Code())
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// IsErrorRetryable returns whether the error is retryable, based on its Code.
// Returns false if the request has no Error set.
//
// Alias for the utility function IsErrorRetryable
func (r *Request) IsErrorRetryable() bool ***REMOVED***
	return IsErrorRetryable(r.Error)
***REMOVED***

// IsErrorThrottle returns whether the error is to be throttled based on its code.
// Returns false if the request has no Error set
//
// Alias for the utility function IsErrorThrottle
func (r *Request) IsErrorThrottle() bool ***REMOVED***
	return IsErrorThrottle(r.Error)
***REMOVED***

// IsErrorExpired returns whether the error code is a credential expiry error.
// Returns false if the request has no Error set.
//
// Alias for the utility function IsErrorExpiredCreds
func (r *Request) IsErrorExpired() bool ***REMOVED***
	return IsErrorExpiredCreds(r.Error)
***REMOVED***
