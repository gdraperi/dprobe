// +build appengine plan9

package request

import (
	"strings"
)

func isErrConnectionReset(err error) bool ***REMOVED***
	return strings.Contains(err.Error(), "connection reset")
***REMOVED***
