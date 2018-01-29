// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssh

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

// Client implements a traditional SSH client that supports shells,
// subprocesses, TCP port/streamlocal forwarding and tunneled dialing.
type Client struct ***REMOVED***
	Conn

	forwards        forwardList // forwarded tcpip connections from the remote side
	mu              sync.Mutex
	channelHandlers map[string]chan NewChannel
***REMOVED***

// HandleChannelOpen returns a channel on which NewChannel requests
// for the given type are sent. If the type already is being handled,
// nil is returned. The channel is closed when the connection is closed.
func (c *Client) HandleChannelOpen(channelType string) <-chan NewChannel ***REMOVED***
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.channelHandlers == nil ***REMOVED***
		// The SSH channel has been closed.
		c := make(chan NewChannel)
		close(c)
		return c
	***REMOVED***

	ch := c.channelHandlers[channelType]
	if ch != nil ***REMOVED***
		return nil
	***REMOVED***

	ch = make(chan NewChannel, chanSize)
	c.channelHandlers[channelType] = ch
	return ch
***REMOVED***

// NewClient creates a Client on top of the given connection.
func NewClient(c Conn, chans <-chan NewChannel, reqs <-chan *Request) *Client ***REMOVED***
	conn := &Client***REMOVED***
		Conn:            c,
		channelHandlers: make(map[string]chan NewChannel, 1),
	***REMOVED***

	go conn.handleGlobalRequests(reqs)
	go conn.handleChannelOpens(chans)
	go func() ***REMOVED***
		conn.Wait()
		conn.forwards.closeAll()
	***REMOVED***()
	go conn.forwards.handleChannels(conn.HandleChannelOpen("forwarded-tcpip"))
	go conn.forwards.handleChannels(conn.HandleChannelOpen("forwarded-streamlocal@openssh.com"))
	return conn
***REMOVED***

// NewClientConn establishes an authenticated SSH connection using c
// as the underlying transport.  The Request and NewChannel channels
// must be serviced or the connection will hang.
func NewClientConn(c net.Conn, addr string, config *ClientConfig) (Conn, <-chan NewChannel, <-chan *Request, error) ***REMOVED***
	fullConf := *config
	fullConf.SetDefaults()
	if fullConf.HostKeyCallback == nil ***REMOVED***
		c.Close()
		return nil, nil, nil, errors.New("ssh: must specify HostKeyCallback")
	***REMOVED***

	conn := &connection***REMOVED***
		sshConn: sshConn***REMOVED***conn: c***REMOVED***,
	***REMOVED***

	if err := conn.clientHandshake(addr, &fullConf); err != nil ***REMOVED***
		c.Close()
		return nil, nil, nil, fmt.Errorf("ssh: handshake failed: %v", err)
	***REMOVED***
	conn.mux = newMux(conn.transport)
	return conn, conn.mux.incomingChannels, conn.mux.incomingRequests, nil
***REMOVED***

// clientHandshake performs the client side key exchange. See RFC 4253 Section
// 7.
func (c *connection) clientHandshake(dialAddress string, config *ClientConfig) error ***REMOVED***
	if config.ClientVersion != "" ***REMOVED***
		c.clientVersion = []byte(config.ClientVersion)
	***REMOVED*** else ***REMOVED***
		c.clientVersion = []byte(packageVersion)
	***REMOVED***
	var err error
	c.serverVersion, err = exchangeVersions(c.sshConn.conn, c.clientVersion)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	c.transport = newClientTransport(
		newTransport(c.sshConn.conn, config.Rand, true /* is client */),
		c.clientVersion, c.serverVersion, config, dialAddress, c.sshConn.RemoteAddr())
	if err := c.transport.waitSession(); err != nil ***REMOVED***
		return err
	***REMOVED***

	c.sessionID = c.transport.getSessionID()
	return c.clientAuthenticate(config)
***REMOVED***

// verifyHostKeySignature verifies the host key obtained in the key
// exchange.
func verifyHostKeySignature(hostKey PublicKey, result *kexResult) error ***REMOVED***
	sig, rest, ok := parseSignatureBody(result.Signature)
	if len(rest) > 0 || !ok ***REMOVED***
		return errors.New("ssh: signature parse error")
	***REMOVED***

	return hostKey.Verify(result.H, sig)
***REMOVED***

// NewSession opens a new Session for this client. (A session is a remote
// execution of a program.)
func (c *Client) NewSession() (*Session, error) ***REMOVED***
	ch, in, err := c.OpenChannel("session", nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return newSession(ch, in)
***REMOVED***

func (c *Client) handleGlobalRequests(incoming <-chan *Request) ***REMOVED***
	for r := range incoming ***REMOVED***
		// This handles keepalive messages and matches
		// the behaviour of OpenSSH.
		r.Reply(false, nil)
	***REMOVED***
***REMOVED***

// handleChannelOpens channel open messages from the remote side.
func (c *Client) handleChannelOpens(in <-chan NewChannel) ***REMOVED***
	for ch := range in ***REMOVED***
		c.mu.Lock()
		handler := c.channelHandlers[ch.ChannelType()]
		c.mu.Unlock()

		if handler != nil ***REMOVED***
			handler <- ch
		***REMOVED*** else ***REMOVED***
			ch.Reject(UnknownChannelType, fmt.Sprintf("unknown channel type: %v", ch.ChannelType()))
		***REMOVED***
	***REMOVED***

	c.mu.Lock()
	for _, ch := range c.channelHandlers ***REMOVED***
		close(ch)
	***REMOVED***
	c.channelHandlers = nil
	c.mu.Unlock()
***REMOVED***

// Dial starts a client connection to the given SSH server. It is a
// convenience function that connects to the given network address,
// initiates the SSH handshake, and then sets up a Client.  For access
// to incoming channels and requests, use net.Dial with NewClientConn
// instead.
func Dial(network, addr string, config *ClientConfig) (*Client, error) ***REMOVED***
	conn, err := net.DialTimeout(network, addr, config.Timeout)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	c, chans, reqs, err := NewClientConn(conn, addr, config)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return NewClient(c, chans, reqs), nil
***REMOVED***

// HostKeyCallback is the function type used for verifying server
// keys.  A HostKeyCallback must return nil if the host key is OK, or
// an error to reject it. It receives the hostname as passed to Dial
// or NewClientConn. The remote address is the RemoteAddr of the
// net.Conn underlying the the SSH connection.
type HostKeyCallback func(hostname string, remote net.Addr, key PublicKey) error

// BannerCallback is the function type used for treat the banner sent by
// the server. A BannerCallback receives the message sent by the remote server.
type BannerCallback func(message string) error

// A ClientConfig structure is used to configure a Client. It must not be
// modified after having been passed to an SSH function.
type ClientConfig struct ***REMOVED***
	// Config contains configuration that is shared between clients and
	// servers.
	Config

	// User contains the username to authenticate as.
	User string

	// Auth contains possible authentication methods to use with the
	// server. Only the first instance of a particular RFC 4252 method will
	// be used during authentication.
	Auth []AuthMethod

	// HostKeyCallback is called during the cryptographic
	// handshake to validate the server's host key. The client
	// configuration must supply this callback for the connection
	// to succeed. The functions InsecureIgnoreHostKey or
	// FixedHostKey can be used for simplistic host key checks.
	HostKeyCallback HostKeyCallback

	// BannerCallback is called during the SSH dance to display a custom
	// server's message. The client configuration can supply this callback to
	// handle it as wished. The function BannerDisplayStderr can be used for
	// simplistic display on Stderr.
	BannerCallback BannerCallback

	// ClientVersion contains the version identification string that will
	// be used for the connection. If empty, a reasonable default is used.
	ClientVersion string

	// HostKeyAlgorithms lists the key types that the client will
	// accept from the server as host key, in order of
	// preference. If empty, a reasonable default is used. Any
	// string returned from PublicKey.Type method may be used, or
	// any of the CertAlgoXxxx and KeyAlgoXxxx constants.
	HostKeyAlgorithms []string

	// Timeout is the maximum amount of time for the TCP connection to establish.
	//
	// A Timeout of zero means no timeout.
	Timeout time.Duration
***REMOVED***

// InsecureIgnoreHostKey returns a function that can be used for
// ClientConfig.HostKeyCallback to accept any host key. It should
// not be used for production code.
func InsecureIgnoreHostKey() HostKeyCallback ***REMOVED***
	return func(hostname string, remote net.Addr, key PublicKey) error ***REMOVED***
		return nil
	***REMOVED***
***REMOVED***

type fixedHostKey struct ***REMOVED***
	key PublicKey
***REMOVED***

func (f *fixedHostKey) check(hostname string, remote net.Addr, key PublicKey) error ***REMOVED***
	if f.key == nil ***REMOVED***
		return fmt.Errorf("ssh: required host key was nil")
	***REMOVED***
	if !bytes.Equal(key.Marshal(), f.key.Marshal()) ***REMOVED***
		return fmt.Errorf("ssh: host key mismatch")
	***REMOVED***
	return nil
***REMOVED***

// FixedHostKey returns a function for use in
// ClientConfig.HostKeyCallback to accept only a specific host key.
func FixedHostKey(key PublicKey) HostKeyCallback ***REMOVED***
	hk := &fixedHostKey***REMOVED***key***REMOVED***
	return hk.check
***REMOVED***

// BannerDisplayStderr returns a function that can be used for
// ClientConfig.BannerCallback to display banners on os.Stderr.
func BannerDisplayStderr() BannerCallback ***REMOVED***
	return func(banner string) error ***REMOVED***
		_, err := os.Stderr.WriteString(banner)

		return err
	***REMOVED***
***REMOVED***
