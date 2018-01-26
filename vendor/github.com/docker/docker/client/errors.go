package client

import (
	"fmt"

	"net/http"

	"github.com/docker/docker/api/types/versions"
	"github.com/pkg/errors"
)

// errConnectionFailed implements an error returned when connection failed.
type errConnectionFailed struct ***REMOVED***
	host string
***REMOVED***

// Error returns a string representation of an errConnectionFailed
func (err errConnectionFailed) Error() string ***REMOVED***
	if err.host == "" ***REMOVED***
		return "Cannot connect to the Docker daemon. Is the docker daemon running on this host?"
	***REMOVED***
	return fmt.Sprintf("Cannot connect to the Docker daemon at %s. Is the docker daemon running?", err.host)
***REMOVED***

// IsErrConnectionFailed returns true if the error is caused by connection failed.
func IsErrConnectionFailed(err error) bool ***REMOVED***
	_, ok := errors.Cause(err).(errConnectionFailed)
	return ok
***REMOVED***

// ErrorConnectionFailed returns an error with host in the error message when connection to docker daemon failed.
func ErrorConnectionFailed(host string) error ***REMOVED***
	return errConnectionFailed***REMOVED***host: host***REMOVED***
***REMOVED***

type notFound interface ***REMOVED***
	error
	NotFound() bool // Is the error a NotFound error
***REMOVED***

// IsErrNotFound returns true if the error is a NotFound error, which is returned
// by the API when some object is not found.
func IsErrNotFound(err error) bool ***REMOVED***
	te, ok := err.(notFound)
	return ok && te.NotFound()
***REMOVED***

type objectNotFoundError struct ***REMOVED***
	object string
	id     string
***REMOVED***

func (e objectNotFoundError) NotFound() bool ***REMOVED***
	return true
***REMOVED***

func (e objectNotFoundError) Error() string ***REMOVED***
	return fmt.Sprintf("Error: No such %s: %s", e.object, e.id)
***REMOVED***

func wrapResponseError(err error, resp serverResponse, object, id string) error ***REMOVED***
	switch ***REMOVED***
	case err == nil:
		return nil
	case resp.statusCode == http.StatusNotFound:
		return objectNotFoundError***REMOVED***object: object, id: id***REMOVED***
	case resp.statusCode == http.StatusNotImplemented:
		return notImplementedError***REMOVED***message: err.Error()***REMOVED***
	default:
		return err
	***REMOVED***
***REMOVED***

// unauthorizedError represents an authorization error in a remote registry.
type unauthorizedError struct ***REMOVED***
	cause error
***REMOVED***

// Error returns a string representation of an unauthorizedError
func (u unauthorizedError) Error() string ***REMOVED***
	return u.cause.Error()
***REMOVED***

// IsErrUnauthorized returns true if the error is caused
// when a remote registry authentication fails
func IsErrUnauthorized(err error) bool ***REMOVED***
	_, ok := err.(unauthorizedError)
	return ok
***REMOVED***

type pluginPermissionDenied struct ***REMOVED***
	name string
***REMOVED***

func (e pluginPermissionDenied) Error() string ***REMOVED***
	return "Permission denied while installing plugin " + e.name
***REMOVED***

// IsErrPluginPermissionDenied returns true if the error is caused
// when a user denies a plugin's permissions
func IsErrPluginPermissionDenied(err error) bool ***REMOVED***
	_, ok := err.(pluginPermissionDenied)
	return ok
***REMOVED***

type notImplementedError struct ***REMOVED***
	message string
***REMOVED***

func (e notImplementedError) Error() string ***REMOVED***
	return e.message
***REMOVED***

func (e notImplementedError) NotImplemented() bool ***REMOVED***
	return true
***REMOVED***

// IsErrNotImplemented returns true if the error is a NotImplemented error.
// This is returned by the API when a requested feature has not been
// implemented.
func IsErrNotImplemented(err error) bool ***REMOVED***
	te, ok := err.(notImplementedError)
	return ok && te.NotImplemented()
***REMOVED***

// NewVersionError returns an error if the APIVersion required
// if less than the current supported version
func (cli *Client) NewVersionError(APIrequired, feature string) error ***REMOVED***
	if cli.version != "" && versions.LessThan(cli.version, APIrequired) ***REMOVED***
		return fmt.Errorf("%q requires API version %s, but the Docker daemon API version is %s", feature, APIrequired, cli.version)
	***REMOVED***
	return nil
***REMOVED***
