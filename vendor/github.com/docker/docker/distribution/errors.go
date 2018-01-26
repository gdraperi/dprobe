package distribution

import (
	"fmt"
	"net/url"
	"strings"
	"syscall"

	"github.com/docker/distribution"
	"github.com/docker/distribution/reference"
	"github.com/docker/distribution/registry/api/errcode"
	"github.com/docker/distribution/registry/api/v2"
	"github.com/docker/distribution/registry/client"
	"github.com/docker/distribution/registry/client/auth"
	"github.com/docker/docker/distribution/xfer"
	"github.com/docker/docker/errdefs"
	"github.com/sirupsen/logrus"
)

// ErrNoSupport is an error type used for errors indicating that an operation
// is not supported. It encapsulates a more specific error.
type ErrNoSupport struct***REMOVED*** Err error ***REMOVED***

func (e ErrNoSupport) Error() string ***REMOVED***
	if e.Err == nil ***REMOVED***
		return "not supported"
	***REMOVED***
	return e.Err.Error()
***REMOVED***

// fallbackError wraps an error that can possibly allow fallback to a different
// endpoint.
type fallbackError struct ***REMOVED***
	// err is the error being wrapped.
	err error
	// confirmedV2 is set to true if it was confirmed that the registry
	// supports the v2 protocol. This is used to limit fallbacks to the v1
	// protocol.
	confirmedV2 bool
	// transportOK is set to true if we managed to speak HTTP with the
	// registry. This confirms that we're using appropriate TLS settings
	// (or lack of TLS).
	transportOK bool
***REMOVED***

// Error renders the FallbackError as a string.
func (f fallbackError) Error() string ***REMOVED***
	return f.Cause().Error()
***REMOVED***

func (f fallbackError) Cause() error ***REMOVED***
	return f.err
***REMOVED***

// shouldV2Fallback returns true if this error is a reason to fall back to v1.
func shouldV2Fallback(err errcode.Error) bool ***REMOVED***
	switch err.Code ***REMOVED***
	case errcode.ErrorCodeUnauthorized, v2.ErrorCodeManifestUnknown, v2.ErrorCodeNameUnknown:
		return true
	***REMOVED***
	return false
***REMOVED***

type notFoundError struct ***REMOVED***
	cause errcode.Error
	ref   reference.Named
***REMOVED***

func (e notFoundError) Error() string ***REMOVED***
	switch e.cause.Code ***REMOVED***
	case errcode.ErrorCodeDenied:
		// ErrorCodeDenied is used when access to the repository was denied
		return fmt.Sprintf("pull access denied for %s, repository does not exist or may require 'docker login'", reference.FamiliarName(e.ref))
	case v2.ErrorCodeManifestUnknown:
		return fmt.Sprintf("manifest for %s not found", reference.FamiliarString(e.ref))
	case v2.ErrorCodeNameUnknown:
		return fmt.Sprintf("repository %s not found", reference.FamiliarName(e.ref))
	***REMOVED***
	// Shouldn't get here, but this is better than returning an empty string
	return e.cause.Message
***REMOVED***

func (e notFoundError) NotFound() ***REMOVED******REMOVED***

func (e notFoundError) Cause() error ***REMOVED***
	return e.cause
***REMOVED***

// TranslatePullError is used to convert an error from a registry pull
// operation to an error representing the entire pull operation. Any error
// information which is not used by the returned error gets output to
// log at info level.
func TranslatePullError(err error, ref reference.Named) error ***REMOVED***
	switch v := err.(type) ***REMOVED***
	case errcode.Errors:
		if len(v) != 0 ***REMOVED***
			for _, extra := range v[1:] ***REMOVED***
				logrus.Infof("Ignoring extra error returned from registry: %v", extra)
			***REMOVED***
			return TranslatePullError(v[0], ref)
		***REMOVED***
	case errcode.Error:
		switch v.Code ***REMOVED***
		case errcode.ErrorCodeDenied, v2.ErrorCodeManifestUnknown, v2.ErrorCodeNameUnknown:
			return notFoundError***REMOVED***v, ref***REMOVED***
		***REMOVED***
	case xfer.DoNotRetry:
		return TranslatePullError(v.Err, ref)
	***REMOVED***

	return errdefs.Unknown(err)
***REMOVED***

// continueOnError returns true if we should fallback to the next endpoint
// as a result of this error.
func continueOnError(err error, mirrorEndpoint bool) bool ***REMOVED***
	switch v := err.(type) ***REMOVED***
	case errcode.Errors:
		if len(v) == 0 ***REMOVED***
			return true
		***REMOVED***
		return continueOnError(v[0], mirrorEndpoint)
	case ErrNoSupport:
		return continueOnError(v.Err, mirrorEndpoint)
	case errcode.Error:
		return mirrorEndpoint || shouldV2Fallback(v)
	case *client.UnexpectedHTTPResponseError:
		return true
	case ImageConfigPullError:
		// ImageConfigPullError only happens with v2 images, v1 fallback is
		// unnecessary.
		// Failures from a mirror endpoint should result in fallback to the
		// canonical repo.
		return mirrorEndpoint
	case error:
		return !strings.Contains(err.Error(), strings.ToLower(syscall.ESRCH.Error()))
	***REMOVED***
	// let's be nice and fallback if the error is a completely
	// unexpected one.
	// If new errors have to be handled in some way, please
	// add them to the switch above.
	return true
***REMOVED***

// retryOnError wraps the error in xfer.DoNotRetry if we should not retry the
// operation after this error.
func retryOnError(err error) error ***REMOVED***
	switch v := err.(type) ***REMOVED***
	case errcode.Errors:
		if len(v) != 0 ***REMOVED***
			return retryOnError(v[0])
		***REMOVED***
	case errcode.Error:
		switch v.Code ***REMOVED***
		case errcode.ErrorCodeUnauthorized, errcode.ErrorCodeUnsupported, errcode.ErrorCodeDenied, errcode.ErrorCodeTooManyRequests, v2.ErrorCodeNameUnknown:
			return xfer.DoNotRetry***REMOVED***Err: err***REMOVED***
		***REMOVED***
	case *url.Error:
		switch v.Err ***REMOVED***
		case auth.ErrNoBasicAuthCredentials, auth.ErrNoToken:
			return xfer.DoNotRetry***REMOVED***Err: v.Err***REMOVED***
		***REMOVED***
		return retryOnError(v.Err)
	case *client.UnexpectedHTTPResponseError:
		return xfer.DoNotRetry***REMOVED***Err: err***REMOVED***
	case error:
		if err == distribution.ErrBlobUnknown ***REMOVED***
			return xfer.DoNotRetry***REMOVED***Err: err***REMOVED***
		***REMOVED***
		if strings.Contains(err.Error(), strings.ToLower(syscall.ENOSPC.Error())) ***REMOVED***
			return xfer.DoNotRetry***REMOVED***Err: err***REMOVED***
		***REMOVED***
	***REMOVED***
	// let's be nice and fallback if the error is a completely
	// unexpected one.
	// If new errors have to be handled in some way, please
	// add them to the switch above.
	return err
***REMOVED***

type invalidManifestClassError struct ***REMOVED***
	mediaType string
	class     string
***REMOVED***

func (e invalidManifestClassError) Error() string ***REMOVED***
	return fmt.Sprintf("Encountered remote %q(%s) when fetching", e.mediaType, e.class)
***REMOVED***

func (e invalidManifestClassError) InvalidParameter() ***REMOVED******REMOVED***

type invalidManifestFormatError struct***REMOVED******REMOVED***

func (invalidManifestFormatError) Error() string ***REMOVED***
	return "unsupported manifest format"
***REMOVED***

func (invalidManifestFormatError) InvalidParameter() ***REMOVED******REMOVED***

type reservedNameError string

func (e reservedNameError) Error() string ***REMOVED***
	return "'" + string(e) + "' is a reserved name"
***REMOVED***

func (e reservedNameError) Forbidden() ***REMOVED******REMOVED***
