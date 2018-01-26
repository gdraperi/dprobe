// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package proxy provides support for a variety of protocols to proxy network
// data.
package proxy // import "golang.org/x/net/proxy"

import (
	"errors"
	"net"
	"net/url"
	"os"
)

// A Dialer is a means to establish a connection.
type Dialer interface ***REMOVED***
	// Dial connects to the given address via the proxy.
	Dial(network, addr string) (c net.Conn, err error)
***REMOVED***

// Auth contains authentication parameters that specific Dialers may require.
type Auth struct ***REMOVED***
	User, Password string
***REMOVED***

// FromEnvironment returns the dialer specified by the proxy related variables in
// the environment.
func FromEnvironment() Dialer ***REMOVED***
	allProxy := os.Getenv("all_proxy")
	if len(allProxy) == 0 ***REMOVED***
		return Direct
	***REMOVED***

	proxyURL, err := url.Parse(allProxy)
	if err != nil ***REMOVED***
		return Direct
	***REMOVED***
	proxy, err := FromURL(proxyURL, Direct)
	if err != nil ***REMOVED***
		return Direct
	***REMOVED***

	noProxy := os.Getenv("no_proxy")
	if len(noProxy) == 0 ***REMOVED***
		return proxy
	***REMOVED***

	perHost := NewPerHost(proxy, Direct)
	perHost.AddFromString(noProxy)
	return perHost
***REMOVED***

// proxySchemes is a map from URL schemes to a function that creates a Dialer
// from a URL with such a scheme.
var proxySchemes map[string]func(*url.URL, Dialer) (Dialer, error)

// RegisterDialerType takes a URL scheme and a function to generate Dialers from
// a URL with that scheme and a forwarding Dialer. Registered schemes are used
// by FromURL.
func RegisterDialerType(scheme string, f func(*url.URL, Dialer) (Dialer, error)) ***REMOVED***
	if proxySchemes == nil ***REMOVED***
		proxySchemes = make(map[string]func(*url.URL, Dialer) (Dialer, error))
	***REMOVED***
	proxySchemes[scheme] = f
***REMOVED***

// FromURL returns a Dialer given a URL specification and an underlying
// Dialer for it to make network requests.
func FromURL(u *url.URL, forward Dialer) (Dialer, error) ***REMOVED***
	var auth *Auth
	if u.User != nil ***REMOVED***
		auth = new(Auth)
		auth.User = u.User.Username()
		if p, ok := u.User.Password(); ok ***REMOVED***
			auth.Password = p
		***REMOVED***
	***REMOVED***

	switch u.Scheme ***REMOVED***
	case "socks5":
		return SOCKS5("tcp", u.Host, auth, forward)
	***REMOVED***

	// If the scheme doesn't match any of the built-in schemes, see if it
	// was registered by another package.
	if proxySchemes != nil ***REMOVED***
		if f, ok := proxySchemes[u.Scheme]; ok ***REMOVED***
			return f(u, forward)
		***REMOVED***
	***REMOVED***

	return nil, errors.New("proxy: unknown scheme: " + u.Scheme)
***REMOVED***
