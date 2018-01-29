// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssh

import (
	"fmt"
	"net"
)

// OpenChannelError is returned if the other side rejects an
// OpenChannel request.
type OpenChannelError struct ***REMOVED***
	Reason  RejectionReason
	Message string
***REMOVED***

func (e *OpenChannelError) Error() string ***REMOVED***
	return fmt.Sprintf("ssh: rejected: %s (%s)", e.Reason, e.Message)
***REMOVED***

// ConnMetadata holds metadata for the connection.
type ConnMetadata interface ***REMOVED***
	// User returns the user ID for this connection.
	User() string

	// SessionID returns the session hash, also denoted by H.
	SessionID() []byte

	// ClientVersion returns the client's version string as hashed
	// into the session ID.
	ClientVersion() []byte

	// ServerVersion returns the server's version string as hashed
	// into the session ID.
	ServerVersion() []byte

	// RemoteAddr returns the remote address for this connection.
	RemoteAddr() net.Addr

	// LocalAddr returns the local address for this connection.
	LocalAddr() net.Addr
***REMOVED***

// Conn represents an SSH connection for both server and client roles.
// Conn is the basis for implementing an application layer, such
// as ClientConn, which implements the traditional shell access for
// clients.
type Conn interface ***REMOVED***
	ConnMetadata

	// SendRequest sends a global request, and returns the
	// reply. If wantReply is true, it returns the response status
	// and payload. See also RFC4254, section 4.
	SendRequest(name string, wantReply bool, payload []byte) (bool, []byte, error)

	// OpenChannel tries to open an channel. If the request is
	// rejected, it returns *OpenChannelError. On success it returns
	// the SSH Channel and a Go channel for incoming, out-of-band
	// requests. The Go channel must be serviced, or the
	// connection will hang.
	OpenChannel(name string, data []byte) (Channel, <-chan *Request, error)

	// Close closes the underlying network connection
	Close() error

	// Wait blocks until the connection has shut down, and returns the
	// error causing the shutdown.
	Wait() error

	// TODO(hanwen): consider exposing:
	//   RequestKeyChange
	//   Disconnect
***REMOVED***

// DiscardRequests consumes and rejects all requests from the
// passed-in channel.
func DiscardRequests(in <-chan *Request) ***REMOVED***
	for req := range in ***REMOVED***
		if req.WantReply ***REMOVED***
			req.Reply(false, nil)
		***REMOVED***
	***REMOVED***
***REMOVED***

// A connection represents an incoming connection.
type connection struct ***REMOVED***
	transport *handshakeTransport
	sshConn

	// The connection protocol.
	*mux
***REMOVED***

func (c *connection) Close() error ***REMOVED***
	return c.sshConn.conn.Close()
***REMOVED***

// sshconn provides net.Conn metadata, but disallows direct reads and
// writes.
type sshConn struct ***REMOVED***
	conn net.Conn

	user          string
	sessionID     []byte
	clientVersion []byte
	serverVersion []byte
***REMOVED***

func dup(src []byte) []byte ***REMOVED***
	dst := make([]byte, len(src))
	copy(dst, src)
	return dst
***REMOVED***

func (c *sshConn) User() string ***REMOVED***
	return c.user
***REMOVED***

func (c *sshConn) RemoteAddr() net.Addr ***REMOVED***
	return c.conn.RemoteAddr()
***REMOVED***

func (c *sshConn) Close() error ***REMOVED***
	return c.conn.Close()
***REMOVED***

func (c *sshConn) LocalAddr() net.Addr ***REMOVED***
	return c.conn.LocalAddr()
***REMOVED***

func (c *sshConn) SessionID() []byte ***REMOVED***
	return dup(c.sessionID)
***REMOVED***

func (c *sshConn) ClientVersion() []byte ***REMOVED***
	return dup(c.clientVersion)
***REMOVED***

func (c *sshConn) ServerVersion() []byte ***REMOVED***
	return dup(c.serverVersion)
***REMOVED***
