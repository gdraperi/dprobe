// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package httpproxy provides support for HTTP proxy determination
// based on environment variables, as provided by net/http's
// ProxyFromEnvironment function.
//
// The API is not subject to the Go 1 compatibility promise and may change at
// any time.
package httpproxy

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
	"unicode/utf8"

	"golang.org/x/net/idna"
)

// Config holds configuration for HTTP proxy settings. See
// FromEnvironment for details.
type Config struct ***REMOVED***
	// HTTPProxy represents the value of the HTTP_PROXY or
	// http_proxy environment variable. It will be used as the proxy
	// URL for HTTP requests and HTTPS requests unless overridden by
	// HTTPSProxy or NoProxy.
	HTTPProxy string

	// HTTPSProxy represents the HTTPS_PROXY or https_proxy
	// environment variable. It will be used as the proxy URL for
	// HTTPS requests unless overridden by NoProxy.
	HTTPSProxy string

	// NoProxy represents the NO_PROXY or no_proxy environment
	// variable. It specifies URLs that should be excluded from
	// proxying as a comma-separated list of domain names or a
	// single asterisk (*) to indicate that no proxying should be
	// done. A domain name matches that name and all subdomains. A
	// domain name with a leading "." matches subdomains only. For
	// example "foo.com" matches "foo.com" and "bar.foo.com";
	// ".y.com" matches "x.y.com" but not "y.com".
	NoProxy string

	// CGI holds whether the current process is running
	// as a CGI handler (FromEnvironment infers this from the
	// presence of a REQUEST_METHOD environment variable).
	// When this is set, ProxyForURL will return an error
	// when HTTPProxy applies, because a client could be
	// setting HTTP_PROXY maliciously. See https://golang.org/s/cgihttpproxy.
	CGI bool
***REMOVED***

// FromEnvironment returns a Config instance populated from the
// environment variables HTTP_PROXY, HTTPS_PROXY and NO_PROXY (or the
// lowercase versions thereof). HTTPS_PROXY takes precedence over
// HTTP_PROXY for https requests.
//
// The environment values may be either a complete URL or a
// "host[:port]", in which case the "http" scheme is assumed. An error
// is returned if the value is a different form.
func FromEnvironment() *Config ***REMOVED***
	return &Config***REMOVED***
		HTTPProxy:  getEnvAny("HTTP_PROXY", "http_proxy"),
		HTTPSProxy: getEnvAny("HTTPS_PROXY", "https_proxy"),
		NoProxy:    getEnvAny("NO_PROXY", "no_proxy"),
		CGI:        os.Getenv("REQUEST_METHOD") != "",
	***REMOVED***
***REMOVED***

func getEnvAny(names ...string) string ***REMOVED***
	for _, n := range names ***REMOVED***
		if val := os.Getenv(n); val != "" ***REMOVED***
			return val
		***REMOVED***
	***REMOVED***
	return ""
***REMOVED***

// ProxyFunc returns a function that determines the proxy URL to use for
// a given request URL. Changing the contents of cfg will not affect
// proxy functions created earlier.
//
// A nil URL and nil error are returned if no proxy is defined in the
// environment, or a proxy should not be used for the given request, as
// defined by NO_PROXY.
//
// As a special case, if req.URL.Host is "localhost" (with or without a
// port number), then a nil URL and nil error will be returned.
func (cfg *Config) ProxyFunc() func(reqURL *url.URL) (*url.URL, error) ***REMOVED***
	// Prevent Config changes from affecting the function calculation.
	// TODO Preprocess proxy settings for more efficient evaluation.
	cfg1 := *cfg
	return cfg1.proxyForURL
***REMOVED***

func (cfg *Config) proxyForURL(reqURL *url.URL) (*url.URL, error) ***REMOVED***
	var proxy string
	if reqURL.Scheme == "https" ***REMOVED***
		proxy = cfg.HTTPSProxy
	***REMOVED***
	if proxy == "" ***REMOVED***
		proxy = cfg.HTTPProxy
		if proxy != "" && cfg.CGI ***REMOVED***
			return nil, errors.New("refusing to use HTTP_PROXY value in CGI environment; see golang.org/s/cgihttpproxy")
		***REMOVED***
	***REMOVED***
	if proxy == "" ***REMOVED***
		return nil, nil
	***REMOVED***
	if !cfg.useProxy(canonicalAddr(reqURL)) ***REMOVED***
		return nil, nil
	***REMOVED***
	proxyURL, err := url.Parse(proxy)
	if err != nil ||
		(proxyURL.Scheme != "http" &&
			proxyURL.Scheme != "https" &&
			proxyURL.Scheme != "socks5") ***REMOVED***
		// proxy was bogus. Try prepending "http://" to it and
		// see if that parses correctly. If not, we fall
		// through and complain about the original one.
		if proxyURL, err := url.Parse("http://" + proxy); err == nil ***REMOVED***
			return proxyURL, nil
		***REMOVED***
	***REMOVED***
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("invalid proxy address %q: %v", proxy, err)
	***REMOVED***
	return proxyURL, nil
***REMOVED***

// useProxy reports whether requests to addr should use a proxy,
// according to the NO_PROXY or no_proxy environment variable.
// addr is always a canonicalAddr with a host and port.
func (cfg *Config) useProxy(addr string) bool ***REMOVED***
	if len(addr) == 0 ***REMOVED***
		return true
	***REMOVED***
	host, _, err := net.SplitHostPort(addr)
	if err != nil ***REMOVED***
		return false
	***REMOVED***
	if host == "localhost" ***REMOVED***
		return false
	***REMOVED***
	if ip := net.ParseIP(host); ip != nil ***REMOVED***
		if ip.IsLoopback() ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	noProxy := cfg.NoProxy
	if noProxy == "*" ***REMOVED***
		return false
	***REMOVED***

	addr = strings.ToLower(strings.TrimSpace(addr))
	if hasPort(addr) ***REMOVED***
		addr = addr[:strings.LastIndex(addr, ":")]
	***REMOVED***

	for _, p := range strings.Split(noProxy, ",") ***REMOVED***
		p = strings.ToLower(strings.TrimSpace(p))
		if len(p) == 0 ***REMOVED***
			continue
		***REMOVED***
		if hasPort(p) ***REMOVED***
			p = p[:strings.LastIndex(p, ":")]
		***REMOVED***
		if addr == p ***REMOVED***
			return false
		***REMOVED***
		if len(p) == 0 ***REMOVED***
			// There is no host part, likely the entry is malformed; ignore.
			continue
		***REMOVED***
		if p[0] == '.' && (strings.HasSuffix(addr, p) || addr == p[1:]) ***REMOVED***
			// no_proxy ".foo.com" matches "bar.foo.com" or "foo.com"
			return false
		***REMOVED***
		if p[0] != '.' && strings.HasSuffix(addr, p) && addr[len(addr)-len(p)-1] == '.' ***REMOVED***
			// no_proxy "foo.com" matches "bar.foo.com"
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

var portMap = map[string]string***REMOVED***
	"http":   "80",
	"https":  "443",
	"socks5": "1080",
***REMOVED***

// canonicalAddr returns url.Host but always with a ":port" suffix
func canonicalAddr(url *url.URL) string ***REMOVED***
	addr := url.Hostname()
	if v, err := idnaASCII(addr); err == nil ***REMOVED***
		addr = v
	***REMOVED***
	port := url.Port()
	if port == "" ***REMOVED***
		port = portMap[url.Scheme]
	***REMOVED***
	return net.JoinHostPort(addr, port)
***REMOVED***

// Given a string of the form "host", "host:port", or "[ipv6::address]:port",
// return true if the string includes a port.
func hasPort(s string) bool ***REMOVED*** return strings.LastIndex(s, ":") > strings.LastIndex(s, "]") ***REMOVED***

func idnaASCII(v string) (string, error) ***REMOVED***
	// TODO: Consider removing this check after verifying performance is okay.
	// Right now punycode verification, length checks, context checks, and the
	// permissible character tests are all omitted. It also prevents the ToASCII
	// call from salvaging an invalid IDN, when possible. As a result it may be
	// possible to have two IDNs that appear identical to the user where the
	// ASCII-only version causes an error downstream whereas the non-ASCII
	// version does not.
	// Note that for correct ASCII IDNs ToASCII will only do considerably more
	// work, but it will not cause an allocation.
	if isASCII(v) ***REMOVED***
		return v, nil
	***REMOVED***
	return idna.Lookup.ToASCII(v)
***REMOVED***

func isASCII(s string) bool ***REMOVED***
	for i := 0; i < len(s); i++ ***REMOVED***
		if s[i] >= utf8.RuneSelf ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***
