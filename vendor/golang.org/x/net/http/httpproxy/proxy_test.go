// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httpproxy_test

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"testing"

	"golang.org/x/net/http/httpproxy"
)

// setHelper calls t.Helper() for Go 1.9+ (see go19_test.go) and does nothing otherwise.
var setHelper = func(t *testing.T) ***REMOVED******REMOVED***

type proxyForURLTest struct ***REMOVED***
	cfg     httpproxy.Config
	req     string // URL to fetch; blank means "http://example.com"
	want    string
	wanterr error
***REMOVED***

func (t proxyForURLTest) String() string ***REMOVED***
	var buf bytes.Buffer
	space := func() ***REMOVED***
		if buf.Len() > 0 ***REMOVED***
			buf.WriteByte(' ')
		***REMOVED***
	***REMOVED***
	if t.cfg.HTTPProxy != "" ***REMOVED***
		fmt.Fprintf(&buf, "http_proxy=%q", t.cfg.HTTPProxy)
	***REMOVED***
	if t.cfg.HTTPSProxy != "" ***REMOVED***
		space()
		fmt.Fprintf(&buf, "https_proxy=%q", t.cfg.HTTPSProxy)
	***REMOVED***
	if t.cfg.NoProxy != "" ***REMOVED***
		space()
		fmt.Fprintf(&buf, "no_proxy=%q", t.cfg.NoProxy)
	***REMOVED***
	req := "http://example.com"
	if t.req != "" ***REMOVED***
		req = t.req
	***REMOVED***
	space()
	fmt.Fprintf(&buf, "req=%q", req)
	return strings.TrimSpace(buf.String())
***REMOVED***

var proxyForURLTests = []proxyForURLTest***REMOVED******REMOVED***
	cfg: httpproxy.Config***REMOVED***
		HTTPProxy: "127.0.0.1:8080",
	***REMOVED***,
	want: "http://127.0.0.1:8080",
***REMOVED***, ***REMOVED***
	cfg: httpproxy.Config***REMOVED***
		HTTPProxy: "cache.corp.example.com:1234",
	***REMOVED***,
	want: "http://cache.corp.example.com:1234",
***REMOVED***, ***REMOVED***
	cfg: httpproxy.Config***REMOVED***
		HTTPProxy: "cache.corp.example.com",
	***REMOVED***,
	want: "http://cache.corp.example.com",
***REMOVED***, ***REMOVED***
	cfg: httpproxy.Config***REMOVED***
		HTTPProxy: "https://cache.corp.example.com",
	***REMOVED***,
	want: "https://cache.corp.example.com",
***REMOVED***, ***REMOVED***
	cfg: httpproxy.Config***REMOVED***
		HTTPProxy: "http://127.0.0.1:8080",
	***REMOVED***,
	want: "http://127.0.0.1:8080",
***REMOVED***, ***REMOVED***
	cfg: httpproxy.Config***REMOVED***
		HTTPProxy: "https://127.0.0.1:8080",
	***REMOVED***,
	want: "https://127.0.0.1:8080",
***REMOVED***, ***REMOVED***
	cfg: httpproxy.Config***REMOVED***
		HTTPProxy: "socks5://127.0.0.1",
	***REMOVED***,
	want: "socks5://127.0.0.1",
***REMOVED***, ***REMOVED***
	// Don't use secure for http
	cfg: httpproxy.Config***REMOVED***
		HTTPProxy:  "http.proxy.tld",
		HTTPSProxy: "secure.proxy.tld",
	***REMOVED***,
	req:  "http://insecure.tld/",
	want: "http://http.proxy.tld",
***REMOVED***, ***REMOVED***
	// Use secure for https.
	cfg: httpproxy.Config***REMOVED***
		HTTPProxy:  "http.proxy.tld",
		HTTPSProxy: "secure.proxy.tld",
	***REMOVED***,
	req:  "https://secure.tld/",
	want: "http://secure.proxy.tld",
***REMOVED***, ***REMOVED***
	cfg: httpproxy.Config***REMOVED***
		HTTPProxy:  "http.proxy.tld",
		HTTPSProxy: "https://secure.proxy.tld",
	***REMOVED***,
	req:  "https://secure.tld/",
	want: "https://secure.proxy.tld",
***REMOVED***, ***REMOVED***
	// Issue 16405: don't use HTTP_PROXY in a CGI environment,
	// where HTTP_PROXY can be attacker-controlled.
	cfg: httpproxy.Config***REMOVED***
		HTTPProxy: "http://10.1.2.3:8080",
		CGI:       true,
	***REMOVED***,
	want:    "<nil>",
	wanterr: errors.New("refusing to use HTTP_PROXY value in CGI environment; see golang.org/s/cgihttpproxy"),
***REMOVED***, ***REMOVED***
	// HTTPS proxy is still used even in CGI environment.
	// (perhaps dubious but it's the historical behaviour).
	cfg: httpproxy.Config***REMOVED***
		HTTPSProxy: "https://secure.proxy.tld",
		CGI:        true,
	***REMOVED***,
	req:  "https://secure.tld/",
	want: "https://secure.proxy.tld",
***REMOVED***, ***REMOVED***
	want: "<nil>",
***REMOVED***, ***REMOVED***
	cfg: httpproxy.Config***REMOVED***
		NoProxy:   "example.com",
		HTTPProxy: "proxy",
	***REMOVED***,
	req:  "http://example.com/",
	want: "<nil>",
***REMOVED***, ***REMOVED***
	cfg: httpproxy.Config***REMOVED***
		NoProxy:   ".example.com",
		HTTPProxy: "proxy",
	***REMOVED***,
	req:  "http://example.com/",
	want: "<nil>",
***REMOVED***, ***REMOVED***
	cfg: httpproxy.Config***REMOVED***
		NoProxy:   "ample.com",
		HTTPProxy: "proxy",
	***REMOVED***,
	req:  "http://example.com/",
	want: "http://proxy",
***REMOVED***, ***REMOVED***
	cfg: httpproxy.Config***REMOVED***
		NoProxy:   "example.com",
		HTTPProxy: "proxy",
	***REMOVED***,
	req:  "http://foo.example.com/",
	want: "<nil>",
***REMOVED***, ***REMOVED***
	cfg: httpproxy.Config***REMOVED***
		NoProxy:   ".foo.com",
		HTTPProxy: "proxy",
	***REMOVED***,
	req:  "http://example.com/",
	want: "http://proxy",
***REMOVED******REMOVED***

func testProxyForURL(t *testing.T, tt proxyForURLTest) ***REMOVED***
	setHelper(t)
	reqURLStr := tt.req
	if reqURLStr == "" ***REMOVED***
		reqURLStr = "http://example.com"
	***REMOVED***
	reqURL, err := url.Parse(reqURLStr)
	if err != nil ***REMOVED***
		t.Errorf("invalid URL %q", reqURLStr)
		return
	***REMOVED***
	cfg := tt.cfg
	proxyForURL := cfg.ProxyFunc()
	url, err := proxyForURL(reqURL)
	if g, e := fmt.Sprintf("%v", err), fmt.Sprintf("%v", tt.wanterr); g != e ***REMOVED***
		t.Errorf("%v: got error = %q, want %q", tt, g, e)
		return
	***REMOVED***
	if got := fmt.Sprintf("%s", url); got != tt.want ***REMOVED***
		t.Errorf("%v: got URL = %q, want %q", tt, url, tt.want)
	***REMOVED***

	// Check that changing the Config doesn't change the results
	// of the functuon.
	cfg = httpproxy.Config***REMOVED******REMOVED***
	url, err = proxyForURL(reqURL)
	if g, e := fmt.Sprintf("%v", err), fmt.Sprintf("%v", tt.wanterr); g != e ***REMOVED***
		t.Errorf("(after mutating config) %v: got error = %q, want %q", tt, g, e)
		return
	***REMOVED***
	if got := fmt.Sprintf("%s", url); got != tt.want ***REMOVED***
		t.Errorf("(after mutating config) %v: got URL = %q, want %q", tt, url, tt.want)
	***REMOVED***
***REMOVED***

func TestProxyForURL(t *testing.T) ***REMOVED***
	for _, tt := range proxyForURLTests ***REMOVED***
		testProxyForURL(t, tt)
	***REMOVED***
***REMOVED***

func TestFromEnvironment(t *testing.T) ***REMOVED***
	os.Setenv("HTTP_PROXY", "httpproxy")
	os.Setenv("HTTPS_PROXY", "httpsproxy")
	os.Setenv("NO_PROXY", "noproxy")
	os.Setenv("REQUEST_METHOD", "")
	got := httpproxy.FromEnvironment()
	want := httpproxy.Config***REMOVED***
		HTTPProxy:  "httpproxy",
		HTTPSProxy: "httpsproxy",
		NoProxy:    "noproxy",
	***REMOVED***
	if *got != want ***REMOVED***
		t.Errorf("unexpected proxy config, got %#v want %#v", got, want)
	***REMOVED***
***REMOVED***

func TestFromEnvironmentWithRequestMethod(t *testing.T) ***REMOVED***
	os.Setenv("HTTP_PROXY", "httpproxy")
	os.Setenv("HTTPS_PROXY", "httpsproxy")
	os.Setenv("NO_PROXY", "noproxy")
	os.Setenv("REQUEST_METHOD", "PUT")
	got := httpproxy.FromEnvironment()
	want := httpproxy.Config***REMOVED***
		HTTPProxy:  "httpproxy",
		HTTPSProxy: "httpsproxy",
		NoProxy:    "noproxy",
		CGI:        true,
	***REMOVED***
	if *got != want ***REMOVED***
		t.Errorf("unexpected proxy config, got %#v want %#v", got, want)
	***REMOVED***
***REMOVED***

func TestFromEnvironmentLowerCase(t *testing.T) ***REMOVED***
	os.Setenv("http_proxy", "httpproxy")
	os.Setenv("https_proxy", "httpsproxy")
	os.Setenv("no_proxy", "noproxy")
	os.Setenv("REQUEST_METHOD", "")
	got := httpproxy.FromEnvironment()
	want := httpproxy.Config***REMOVED***
		HTTPProxy:  "httpproxy",
		HTTPSProxy: "httpsproxy",
		NoProxy:    "noproxy",
	***REMOVED***
	if *got != want ***REMOVED***
		t.Errorf("unexpected proxy config, got %#v want %#v", got, want)
	***REMOVED***
***REMOVED***

var UseProxyTests = []struct ***REMOVED***
	host  string
	match bool
***REMOVED******REMOVED***
	// Never proxy localhost:
	***REMOVED***"localhost", false***REMOVED***,
	***REMOVED***"127.0.0.1", false***REMOVED***,
	***REMOVED***"127.0.0.2", false***REMOVED***,
	***REMOVED***"[::1]", false***REMOVED***,
	***REMOVED***"[::2]", true***REMOVED***, // not a loopback address

	***REMOVED***"barbaz.net", false***REMOVED***,     // match as .barbaz.net
	***REMOVED***"foobar.com", false***REMOVED***,     // have a port but match
	***REMOVED***"foofoobar.com", true***REMOVED***,   // not match as a part of foobar.com
	***REMOVED***"baz.com", true***REMOVED***,         // not match as a part of barbaz.com
	***REMOVED***"localhost.net", true***REMOVED***,   // not match as suffix of address
	***REMOVED***"local.localhost", true***REMOVED***, // not match as prefix as address
	***REMOVED***"barbarbaz.net", true***REMOVED***,   // not match because NO_PROXY have a '.'
	***REMOVED***"www.foobar.com", false***REMOVED***, // match because NO_PROXY includes "foobar.com"
***REMOVED***

func TestUseProxy(t *testing.T) ***REMOVED***
	cfg := &httpproxy.Config***REMOVED***
		NoProxy: "foobar.com, .barbaz.net",
	***REMOVED***
	for _, test := range UseProxyTests ***REMOVED***
		if httpproxy.ExportUseProxy(cfg, test.host+":80") != test.match ***REMOVED***
			t.Errorf("useProxy(%v) = %v, want %v", test.host, !test.match, test.match)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestInvalidNoProxy(t *testing.T) ***REMOVED***
	cfg := &httpproxy.Config***REMOVED***
		NoProxy: ":1",
	***REMOVED***
	ok := httpproxy.ExportUseProxy(cfg, "example.com:80") // should not panic
	if !ok ***REMOVED***
		t.Errorf("useProxy unexpected return; got false; want true")
	***REMOVED***
***REMOVED***
