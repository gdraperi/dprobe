package cli

import (
	"fmt"
	"strings"
)

// Errors is a list of errors.
// Useful in a loop if you don't want to return the error right away and you want to display after the loop,
// all the errors that happened during the loop.
type Errors []error

func (errList Errors) Error() string ***REMOVED***
	if len(errList) < 1 ***REMOVED***
		return ""
	***REMOVED***

	out := make([]string, len(errList))
	for i := range errList ***REMOVED***
		out[i] = errList[i].Error()
	***REMOVED***
	return strings.Join(out, ", ")
***REMOVED***

// StatusError reports an unsuccessful exit by a command.
type StatusError struct ***REMOVED***
	Status     string
	StatusCode int
***REMOVED***

func (e StatusError) Error() string ***REMOVED***
	return fmt.Sprintf("Status: %s, Code: %d", e.Status, e.StatusCode)
***REMOVED***
