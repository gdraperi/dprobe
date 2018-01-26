package sockets

import (
	"net"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/proxy"
)

// GetProxyEnv allows access to the uppercase and the lowercase forms of
// proxy-related variables.  See the Go specification for details on these
// variables. https://golang.org/pkg/net/http/
func GetProxyEnv(key string) string ***REMOVED***
	proxyValue := os.Getenv(strings.ToUpper(key))
	if proxyValue == "" ***REMOVED***
		return os.Getenv(strings.ToLower(key))
	***REMOVED***
	return proxyValue
***REMOVED***

// DialerFromEnvironment takes in a "direct" *net.Dialer and returns a
// proxy.Dialer which will route the connections through the proxy using the
// given dialer.
func DialerFromEnvironment(direct *net.Dialer) (proxy.Dialer, error) ***REMOVED***
	allProxy := GetProxyEnv("all_proxy")
	if len(allProxy) == 0 ***REMOVED***
		return direct, nil
	***REMOVED***

	proxyURL, err := url.Parse(allProxy)
	if err != nil ***REMOVED***
		return direct, err
	***REMOVED***

	proxyFromURL, err := proxy.FromURL(proxyURL, direct)
	if err != nil ***REMOVED***
		return direct, err
	***REMOVED***

	noProxy := GetProxyEnv("no_proxy")
	if len(noProxy) == 0 ***REMOVED***
		return proxyFromURL, nil
	***REMOVED***

	perHost := proxy.NewPerHost(proxyFromURL, direct)
	perHost.AddFromString(noProxy)

	return perHost, nil
***REMOVED***
