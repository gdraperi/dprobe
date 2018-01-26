// +build go1.5

package v4

import (
	"net/url"
	"strings"
)

func getURIPath(u *url.URL) string ***REMOVED***
	var uri string

	if len(u.Opaque) > 0 ***REMOVED***
		uri = "/" + strings.Join(strings.Split(u.Opaque, "/")[3:], "/")
	***REMOVED*** else ***REMOVED***
		uri = u.EscapedPath()
	***REMOVED***

	if len(uri) == 0 ***REMOVED***
		uri = "/"
	***REMOVED***

	return uri
***REMOVED***
