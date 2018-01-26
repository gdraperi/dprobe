/*
 *
 * Copyright 2014, Google Inc.
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

// Package credentials implements various credentials supported by gRPC library,
// which encapsulate all the state needed by a client to authenticate with a
// server and make various assertions, e.g., about the client's identity, role,
// or whether it is authorized to make a particular call.
package credentials // import "google.golang.org/grpc/credentials"

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"strings"

	"golang.org/x/net/context"
)

var (
	// alpnProtoStr are the specified application level protocols for gRPC.
	alpnProtoStr = []string***REMOVED***"h2"***REMOVED***
)

// PerRPCCredentials defines the common interface for the credentials which need to
// attach security information to every RPC (e.g., oauth2).
type PerRPCCredentials interface ***REMOVED***
	// GetRequestMetadata gets the current request metadata, refreshing
	// tokens if required. This should be called by the transport layer on
	// each request, and the data should be populated in headers or other
	// context. uri is the URI of the entry point for the request. When
	// supported by the underlying implementation, ctx can be used for
	// timeout and cancellation.
	// TODO(zhaoq): Define the set of the qualified keys instead of leaving
	// it as an arbitrary string.
	GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error)
	// RequireTransportSecurity indicates whether the credentials requires
	// transport security.
	RequireTransportSecurity() bool
***REMOVED***

// ProtocolInfo provides information regarding the gRPC wire protocol version,
// security protocol, security protocol version in use, server name, etc.
type ProtocolInfo struct ***REMOVED***
	// ProtocolVersion is the gRPC wire protocol version.
	ProtocolVersion string
	// SecurityProtocol is the security protocol in use.
	SecurityProtocol string
	// SecurityVersion is the security protocol version.
	SecurityVersion string
	// ServerName is the user-configured server name.
	ServerName string
***REMOVED***

// AuthInfo defines the common interface for the auth information the users are interested in.
type AuthInfo interface ***REMOVED***
	AuthType() string
***REMOVED***

var (
	// ErrConnDispatched indicates that rawConn has been dispatched out of gRPC
	// and the caller should not close rawConn.
	ErrConnDispatched = errors.New("credentials: rawConn is dispatched out of gRPC")
)

// TransportCredentials defines the common interface for all the live gRPC wire
// protocols and supported transport security protocols (e.g., TLS, SSL).
type TransportCredentials interface ***REMOVED***
	// ClientHandshake does the authentication handshake specified by the corresponding
	// authentication protocol on rawConn for clients. It returns the authenticated
	// connection and the corresponding auth information about the connection.
	// Implementations must use the provided context to implement timely cancellation.
	// gRPC will try to reconnect if the error returned is a temporary error
	// (io.EOF, context.DeadlineExceeded or err.Temporary() == true).
	// If the returned error is a wrapper error, implementations should make sure that
	// the error implements Temporary() to have the correct retry behaviors.
	ClientHandshake(context.Context, string, net.Conn) (net.Conn, AuthInfo, error)
	// ServerHandshake does the authentication handshake for servers. It returns
	// the authenticated connection and the corresponding auth information about
	// the connection.
	ServerHandshake(net.Conn) (net.Conn, AuthInfo, error)
	// Info provides the ProtocolInfo of this TransportCredentials.
	Info() ProtocolInfo
	// Clone makes a copy of this TransportCredentials.
	Clone() TransportCredentials
	// OverrideServerName overrides the server name used to verify the hostname on the returned certificates from the server.
	// gRPC internals also use it to override the virtual hosting name if it is set.
	// It must be called before dialing. Currently, this is only used by grpclb.
	OverrideServerName(string) error
***REMOVED***

// TLSInfo contains the auth information for a TLS authenticated connection.
// It implements the AuthInfo interface.
type TLSInfo struct ***REMOVED***
	State tls.ConnectionState
***REMOVED***

// AuthType returns the type of TLSInfo as a string.
func (t TLSInfo) AuthType() string ***REMOVED***
	return "tls"
***REMOVED***

// tlsCreds is the credentials required for authenticating a connection using TLS.
type tlsCreds struct ***REMOVED***
	// TLS configuration
	config *tls.Config
***REMOVED***

func (c tlsCreds) Info() ProtocolInfo ***REMOVED***
	return ProtocolInfo***REMOVED***
		SecurityProtocol: "tls",
		SecurityVersion:  "1.2",
		ServerName:       c.config.ServerName,
	***REMOVED***
***REMOVED***

func (c *tlsCreds) ClientHandshake(ctx context.Context, addr string, rawConn net.Conn) (_ net.Conn, _ AuthInfo, err error) ***REMOVED***
	// use local cfg to avoid clobbering ServerName if using multiple endpoints
	cfg := cloneTLSConfig(c.config)
	if cfg.ServerName == "" ***REMOVED***
		colonPos := strings.LastIndex(addr, ":")
		if colonPos == -1 ***REMOVED***
			colonPos = len(addr)
		***REMOVED***
		cfg.ServerName = addr[:colonPos]
	***REMOVED***
	conn := tls.Client(rawConn, cfg)
	errChannel := make(chan error, 1)
	go func() ***REMOVED***
		errChannel <- conn.Handshake()
	***REMOVED***()
	select ***REMOVED***
	case err := <-errChannel:
		if err != nil ***REMOVED***
			return nil, nil, err
		***REMOVED***
	case <-ctx.Done():
		return nil, nil, ctx.Err()
	***REMOVED***
	return conn, TLSInfo***REMOVED***conn.ConnectionState()***REMOVED***, nil
***REMOVED***

func (c *tlsCreds) ServerHandshake(rawConn net.Conn) (net.Conn, AuthInfo, error) ***REMOVED***
	conn := tls.Server(rawConn, c.config)
	if err := conn.Handshake(); err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	return conn, TLSInfo***REMOVED***conn.ConnectionState()***REMOVED***, nil
***REMOVED***

func (c *tlsCreds) Clone() TransportCredentials ***REMOVED***
	return NewTLS(c.config)
***REMOVED***

func (c *tlsCreds) OverrideServerName(serverNameOverride string) error ***REMOVED***
	c.config.ServerName = serverNameOverride
	return nil
***REMOVED***

// NewTLS uses c to construct a TransportCredentials based on TLS.
func NewTLS(c *tls.Config) TransportCredentials ***REMOVED***
	tc := &tlsCreds***REMOVED***cloneTLSConfig(c)***REMOVED***
	tc.config.NextProtos = alpnProtoStr
	return tc
***REMOVED***

// NewClientTLSFromCert constructs a TLS from the input certificate for client.
// serverNameOverride is for testing only. If set to a non empty string,
// it will override the virtual host name of authority (e.g. :authority header field) in requests.
func NewClientTLSFromCert(cp *x509.CertPool, serverNameOverride string) TransportCredentials ***REMOVED***
	return NewTLS(&tls.Config***REMOVED***ServerName: serverNameOverride, RootCAs: cp***REMOVED***)
***REMOVED***

// NewClientTLSFromFile constructs a TLS from the input certificate file for client.
// serverNameOverride is for testing only. If set to a non empty string,
// it will override the virtual host name of authority (e.g. :authority header field) in requests.
func NewClientTLSFromFile(certFile, serverNameOverride string) (TransportCredentials, error) ***REMOVED***
	b, err := ioutil.ReadFile(certFile)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(b) ***REMOVED***
		return nil, fmt.Errorf("credentials: failed to append certificates")
	***REMOVED***
	return NewTLS(&tls.Config***REMOVED***ServerName: serverNameOverride, RootCAs: cp***REMOVED***), nil
***REMOVED***

// NewServerTLSFromCert constructs a TLS from the input certificate for server.
func NewServerTLSFromCert(cert *tls.Certificate) TransportCredentials ***REMOVED***
	return NewTLS(&tls.Config***REMOVED***Certificates: []tls.Certificate***REMOVED****cert***REMOVED******REMOVED***)
***REMOVED***

// NewServerTLSFromFile constructs a TLS from the input certificate file and key
// file for server.
func NewServerTLSFromFile(certFile, keyFile string) (TransportCredentials, error) ***REMOVED***
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return NewTLS(&tls.Config***REMOVED***Certificates: []tls.Certificate***REMOVED***cert***REMOVED******REMOVED***), nil
***REMOVED***
