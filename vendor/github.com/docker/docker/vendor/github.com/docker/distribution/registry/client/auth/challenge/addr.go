package challenge

import (
	"net/url"
	"strings"
)

// FROM: https://golang.org/src/net/http/http.go
// Given a string of the form "host", "host:port", or "[ipv6::address]:port",
// return true if the string includes a port.
func hasPort(s string) bool ***REMOVED*** return strings.LastIndex(s, ":") > strings.LastIndex(s, "]") ***REMOVED***

// FROM: http://golang.org/src/net/http/transport.go
var portMap = map[string]string***REMOVED***
	"http":  "80",
	"https": "443",
***REMOVED***

// canonicalAddr returns url.Host but always with a ":port" suffix
// FROM: http://golang.org/src/net/http/transport.go
func canonicalAddr(url *url.URL) string ***REMOVED***
	addr := url.Host
	if !hasPort(addr) ***REMOVED***
		return addr + ":" + portMap[url.Scheme]
	***REMOVED***
	return addr
***REMOVED***
