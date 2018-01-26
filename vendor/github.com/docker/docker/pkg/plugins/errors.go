package plugins

import (
	"fmt"
	"net/http"
)

type statusError struct ***REMOVED***
	status int
	method string
	err    string
***REMOVED***

// Error returns a formatted string for this error type
func (e *statusError) Error() string ***REMOVED***
	return fmt.Sprintf("%s: %v", e.method, e.err)
***REMOVED***

// IsNotFound indicates if the passed in error is from an http.StatusNotFound from the plugin
func IsNotFound(err error) bool ***REMOVED***
	return isStatusError(err, http.StatusNotFound)
***REMOVED***

func isStatusError(err error, status int) bool ***REMOVED***
	if err == nil ***REMOVED***
		return false
	***REMOVED***
	e, ok := err.(*statusError)
	if !ok ***REMOVED***
		return false
	***REMOVED***
	return e.status == status
***REMOVED***
