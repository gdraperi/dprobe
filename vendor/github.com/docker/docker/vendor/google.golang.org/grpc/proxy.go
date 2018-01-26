/*
 *
 * Copyright 2017, Google Inc.
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are
 * met:
 *
 *     * Redistributions of source code must retain the above copyright
 * notice, this list of conditions and the following disclaimer.
 *     * Redistributions in binary form must reproduce the above
 * copyright notice, this list of conditions and the following disclaimer
 * in the documentation and/or other materials provided with the
 * distribution.
 *     * Neither the name of Google Inc. nor the names of its
 * contributors may be used to endorse or promote products derived from
 * this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 * "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
 * A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
 * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
 * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

package grpc

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"

	"golang.org/x/net/context"
)

var (
	// errDisabled indicates that proxy is disabled for the address.
	errDisabled = errors.New("proxy is disabled for the address")
	// The following variable will be overwritten in the tests.
	httpProxyFromEnvironment = http.ProxyFromEnvironment
)

func mapAddress(ctx context.Context, address string) (string, error) ***REMOVED***
	req := &http.Request***REMOVED***
		URL: &url.URL***REMOVED***
			Scheme: "https",
			Host:   address,
		***REMOVED***,
	***REMOVED***
	url, err := httpProxyFromEnvironment(req)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	if url == nil ***REMOVED***
		return "", errDisabled
	***REMOVED***
	return url.Host, nil
***REMOVED***

// To read a response from a net.Conn, http.ReadResponse() takes a bufio.Reader.
// It's possible that this reader reads more than what's need for the response and stores
// those bytes in the buffer.
// bufConn wraps the original net.Conn and the bufio.Reader to make sure we don't lose the
// bytes in the buffer.
type bufConn struct ***REMOVED***
	net.Conn
	r io.Reader
***REMOVED***

func (c *bufConn) Read(b []byte) (int, error) ***REMOVED***
	return c.r.Read(b)
***REMOVED***

func doHTTPConnectHandshake(ctx context.Context, conn net.Conn, addr string) (_ net.Conn, err error) ***REMOVED***
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			conn.Close()
		***REMOVED***
	***REMOVED***()

	req := (&http.Request***REMOVED***
		Method: http.MethodConnect,
		URL:    &url.URL***REMOVED***Host: addr***REMOVED***,
		Header: map[string][]string***REMOVED***"User-Agent": ***REMOVED***grpcUA***REMOVED******REMOVED***,
	***REMOVED***)

	if err := sendHTTPRequest(ctx, req, conn); err != nil ***REMOVED***
		return nil, fmt.Errorf("failed to write the HTTP request: %v", err)
	***REMOVED***

	r := bufio.NewReader(conn)
	resp, err := http.ReadResponse(r, req)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("reading server HTTP response: %v", err)
	***REMOVED***
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK ***REMOVED***
		dump, err := httputil.DumpResponse(resp, true)
		if err != nil ***REMOVED***
			return nil, fmt.Errorf("failed to do connect handshake, status code: %s", resp.Status)
		***REMOVED***
		return nil, fmt.Errorf("failed to do connect handshake, response: %q", dump)
	***REMOVED***

	return &bufConn***REMOVED***Conn: conn, r: r***REMOVED***, nil
***REMOVED***

// newProxyDialer returns a dialer that connects to proxy first if necessary.
// The returned dialer checks if a proxy is necessary, dial to the proxy with the
// provided dialer, does HTTP CONNECT handshake and returns the connection.
func newProxyDialer(dialer func(context.Context, string) (net.Conn, error)) func(context.Context, string) (net.Conn, error) ***REMOVED***
	return func(ctx context.Context, addr string) (conn net.Conn, err error) ***REMOVED***
		var skipHandshake bool
		newAddr, err := mapAddress(ctx, addr)
		if err != nil ***REMOVED***
			if err != errDisabled ***REMOVED***
				return nil, err
			***REMOVED***
			skipHandshake = true
			newAddr = addr
		***REMOVED***

		conn, err = dialer(ctx, newAddr)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		if !skipHandshake ***REMOVED***
			conn, err = doHTTPConnectHandshake(ctx, conn, addr)
		***REMOVED***
		return
	***REMOVED***
***REMOVED***
