// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssh

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Listen requests the remote peer open a listening socket on
// addr. Incoming connections will be available by calling Accept on
// the returned net.Listener. The listener must be serviced, or the
// SSH connection may hang.
// N must be "tcp", "tcp4", "tcp6", or "unix".
func (c *Client) Listen(n, addr string) (net.Listener, error) ***REMOVED***
	switch n ***REMOVED***
	case "tcp", "tcp4", "tcp6":
		laddr, err := net.ResolveTCPAddr(n, addr)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return c.ListenTCP(laddr)
	case "unix":
		return c.ListenUnix(addr)
	default:
		return nil, fmt.Errorf("ssh: unsupported protocol: %s", n)
	***REMOVED***
***REMOVED***

// Automatic port allocation is broken with OpenSSH before 6.0. See
// also https://bugzilla.mindrot.org/show_bug.cgi?id=2017.  In
// particular, OpenSSH 5.9 sends a channelOpenMsg with port number 0,
// rather than the actual port number. This means you can never open
// two different listeners with auto allocated ports. We work around
// this by trying explicit ports until we succeed.

const openSSHPrefix = "OpenSSH_"

var portRandomizer = rand.New(rand.NewSource(time.Now().UnixNano()))

// isBrokenOpenSSHVersion returns true if the given version string
// specifies a version of OpenSSH that is known to have a bug in port
// forwarding.
func isBrokenOpenSSHVersion(versionStr string) bool ***REMOVED***
	i := strings.Index(versionStr, openSSHPrefix)
	if i < 0 ***REMOVED***
		return false
	***REMOVED***
	i += len(openSSHPrefix)
	j := i
	for ; j < len(versionStr); j++ ***REMOVED***
		if versionStr[j] < '0' || versionStr[j] > '9' ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	version, _ := strconv.Atoi(versionStr[i:j])
	return version < 6
***REMOVED***

// autoPortListenWorkaround simulates automatic port allocation by
// trying random ports repeatedly.
func (c *Client) autoPortListenWorkaround(laddr *net.TCPAddr) (net.Listener, error) ***REMOVED***
	var sshListener net.Listener
	var err error
	const tries = 10
	for i := 0; i < tries; i++ ***REMOVED***
		addr := *laddr
		addr.Port = 1024 + portRandomizer.Intn(60000)
		sshListener, err = c.ListenTCP(&addr)
		if err == nil ***REMOVED***
			laddr.Port = addr.Port
			return sshListener, err
		***REMOVED***
	***REMOVED***
	return nil, fmt.Errorf("ssh: listen on random port failed after %d tries: %v", tries, err)
***REMOVED***

// RFC 4254 7.1
type channelForwardMsg struct ***REMOVED***
	addr  string
	rport uint32
***REMOVED***

// ListenTCP requests the remote peer open a listening socket
// on laddr. Incoming connections will be available by calling
// Accept on the returned net.Listener.
func (c *Client) ListenTCP(laddr *net.TCPAddr) (net.Listener, error) ***REMOVED***
	if laddr.Port == 0 && isBrokenOpenSSHVersion(string(c.ServerVersion())) ***REMOVED***
		return c.autoPortListenWorkaround(laddr)
	***REMOVED***

	m := channelForwardMsg***REMOVED***
		laddr.IP.String(),
		uint32(laddr.Port),
	***REMOVED***
	// send message
	ok, resp, err := c.SendRequest("tcpip-forward", true, Marshal(&m))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if !ok ***REMOVED***
		return nil, errors.New("ssh: tcpip-forward request denied by peer")
	***REMOVED***

	// If the original port was 0, then the remote side will
	// supply a real port number in the response.
	if laddr.Port == 0 ***REMOVED***
		var p struct ***REMOVED***
			Port uint32
		***REMOVED***
		if err := Unmarshal(resp, &p); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		laddr.Port = int(p.Port)
	***REMOVED***

	// Register this forward, using the port number we obtained.
	ch := c.forwards.add(laddr)

	return &tcpListener***REMOVED***laddr, c, ch***REMOVED***, nil
***REMOVED***

// forwardList stores a mapping between remote
// forward requests and the tcpListeners.
type forwardList struct ***REMOVED***
	sync.Mutex
	entries []forwardEntry
***REMOVED***

// forwardEntry represents an established mapping of a laddr on a
// remote ssh server to a channel connected to a tcpListener.
type forwardEntry struct ***REMOVED***
	laddr net.Addr
	c     chan forward
***REMOVED***

// forward represents an incoming forwarded tcpip connection. The
// arguments to add/remove/lookup should be address as specified in
// the original forward-request.
type forward struct ***REMOVED***
	newCh NewChannel // the ssh client channel underlying this forward
	raddr net.Addr   // the raddr of the incoming connection
***REMOVED***

func (l *forwardList) add(addr net.Addr) chan forward ***REMOVED***
	l.Lock()
	defer l.Unlock()
	f := forwardEntry***REMOVED***
		laddr: addr,
		c:     make(chan forward, 1),
	***REMOVED***
	l.entries = append(l.entries, f)
	return f.c
***REMOVED***

// See RFC 4254, section 7.2
type forwardedTCPPayload struct ***REMOVED***
	Addr       string
	Port       uint32
	OriginAddr string
	OriginPort uint32
***REMOVED***

// parseTCPAddr parses the originating address from the remote into a *net.TCPAddr.
func parseTCPAddr(addr string, port uint32) (*net.TCPAddr, error) ***REMOVED***
	if port == 0 || port > 65535 ***REMOVED***
		return nil, fmt.Errorf("ssh: port number out of range: %d", port)
	***REMOVED***
	ip := net.ParseIP(string(addr))
	if ip == nil ***REMOVED***
		return nil, fmt.Errorf("ssh: cannot parse IP address %q", addr)
	***REMOVED***
	return &net.TCPAddr***REMOVED***IP: ip, Port: int(port)***REMOVED***, nil
***REMOVED***

func (l *forwardList) handleChannels(in <-chan NewChannel) ***REMOVED***
	for ch := range in ***REMOVED***
		var (
			laddr net.Addr
			raddr net.Addr
			err   error
		)
		switch channelType := ch.ChannelType(); channelType ***REMOVED***
		case "forwarded-tcpip":
			var payload forwardedTCPPayload
			if err = Unmarshal(ch.ExtraData(), &payload); err != nil ***REMOVED***
				ch.Reject(ConnectionFailed, "could not parse forwarded-tcpip payload: "+err.Error())
				continue
			***REMOVED***

			// RFC 4254 section 7.2 specifies that incoming
			// addresses should list the address, in string
			// format. It is implied that this should be an IP
			// address, as it would be impossible to connect to it
			// otherwise.
			laddr, err = parseTCPAddr(payload.Addr, payload.Port)
			if err != nil ***REMOVED***
				ch.Reject(ConnectionFailed, err.Error())
				continue
			***REMOVED***
			raddr, err = parseTCPAddr(payload.OriginAddr, payload.OriginPort)
			if err != nil ***REMOVED***
				ch.Reject(ConnectionFailed, err.Error())
				continue
			***REMOVED***

		case "forwarded-streamlocal@openssh.com":
			var payload forwardedStreamLocalPayload
			if err = Unmarshal(ch.ExtraData(), &payload); err != nil ***REMOVED***
				ch.Reject(ConnectionFailed, "could not parse forwarded-streamlocal@openssh.com payload: "+err.Error())
				continue
			***REMOVED***
			laddr = &net.UnixAddr***REMOVED***
				Name: payload.SocketPath,
				Net:  "unix",
			***REMOVED***
			raddr = &net.UnixAddr***REMOVED***
				Name: "@",
				Net:  "unix",
			***REMOVED***
		default:
			panic(fmt.Errorf("ssh: unknown channel type %s", channelType))
		***REMOVED***
		if ok := l.forward(laddr, raddr, ch); !ok ***REMOVED***
			// Section 7.2, implementations MUST reject spurious incoming
			// connections.
			ch.Reject(Prohibited, "no forward for address")
			continue
		***REMOVED***

	***REMOVED***
***REMOVED***

// remove removes the forward entry, and the channel feeding its
// listener.
func (l *forwardList) remove(addr net.Addr) ***REMOVED***
	l.Lock()
	defer l.Unlock()
	for i, f := range l.entries ***REMOVED***
		if addr.Network() == f.laddr.Network() && addr.String() == f.laddr.String() ***REMOVED***
			l.entries = append(l.entries[:i], l.entries[i+1:]...)
			close(f.c)
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

// closeAll closes and clears all forwards.
func (l *forwardList) closeAll() ***REMOVED***
	l.Lock()
	defer l.Unlock()
	for _, f := range l.entries ***REMOVED***
		close(f.c)
	***REMOVED***
	l.entries = nil
***REMOVED***

func (l *forwardList) forward(laddr, raddr net.Addr, ch NewChannel) bool ***REMOVED***
	l.Lock()
	defer l.Unlock()
	for _, f := range l.entries ***REMOVED***
		if laddr.Network() == f.laddr.Network() && laddr.String() == f.laddr.String() ***REMOVED***
			f.c <- forward***REMOVED***newCh: ch, raddr: raddr***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

type tcpListener struct ***REMOVED***
	laddr *net.TCPAddr

	conn *Client
	in   <-chan forward
***REMOVED***

// Accept waits for and returns the next connection to the listener.
func (l *tcpListener) Accept() (net.Conn, error) ***REMOVED***
	s, ok := <-l.in
	if !ok ***REMOVED***
		return nil, io.EOF
	***REMOVED***
	ch, incoming, err := s.newCh.Accept()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	go DiscardRequests(incoming)

	return &chanConn***REMOVED***
		Channel: ch,
		laddr:   l.laddr,
		raddr:   s.raddr,
	***REMOVED***, nil
***REMOVED***

// Close closes the listener.
func (l *tcpListener) Close() error ***REMOVED***
	m := channelForwardMsg***REMOVED***
		l.laddr.IP.String(),
		uint32(l.laddr.Port),
	***REMOVED***

	// this also closes the listener.
	l.conn.forwards.remove(l.laddr)
	ok, _, err := l.conn.SendRequest("cancel-tcpip-forward", true, Marshal(&m))
	if err == nil && !ok ***REMOVED***
		err = errors.New("ssh: cancel-tcpip-forward failed")
	***REMOVED***
	return err
***REMOVED***

// Addr returns the listener's network address.
func (l *tcpListener) Addr() net.Addr ***REMOVED***
	return l.laddr
***REMOVED***

// Dial initiates a connection to the addr from the remote host.
// The resulting connection has a zero LocalAddr() and RemoteAddr().
func (c *Client) Dial(n, addr string) (net.Conn, error) ***REMOVED***
	var ch Channel
	switch n ***REMOVED***
	case "tcp", "tcp4", "tcp6":
		// Parse the address into host and numeric port.
		host, portString, err := net.SplitHostPort(addr)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		port, err := strconv.ParseUint(portString, 10, 16)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		ch, err = c.dial(net.IPv4zero.String(), 0, host, int(port))
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		// Use a zero address for local and remote address.
		zeroAddr := &net.TCPAddr***REMOVED***
			IP:   net.IPv4zero,
			Port: 0,
		***REMOVED***
		return &chanConn***REMOVED***
			Channel: ch,
			laddr:   zeroAddr,
			raddr:   zeroAddr,
		***REMOVED***, nil
	case "unix":
		var err error
		ch, err = c.dialStreamLocal(addr)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return &chanConn***REMOVED***
			Channel: ch,
			laddr: &net.UnixAddr***REMOVED***
				Name: "@",
				Net:  "unix",
			***REMOVED***,
			raddr: &net.UnixAddr***REMOVED***
				Name: addr,
				Net:  "unix",
			***REMOVED***,
		***REMOVED***, nil
	default:
		return nil, fmt.Errorf("ssh: unsupported protocol: %s", n)
	***REMOVED***
***REMOVED***

// DialTCP connects to the remote address raddr on the network net,
// which must be "tcp", "tcp4", or "tcp6".  If laddr is not nil, it is used
// as the local address for the connection.
func (c *Client) DialTCP(n string, laddr, raddr *net.TCPAddr) (net.Conn, error) ***REMOVED***
	if laddr == nil ***REMOVED***
		laddr = &net.TCPAddr***REMOVED***
			IP:   net.IPv4zero,
			Port: 0,
		***REMOVED***
	***REMOVED***
	ch, err := c.dial(laddr.IP.String(), laddr.Port, raddr.IP.String(), raddr.Port)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &chanConn***REMOVED***
		Channel: ch,
		laddr:   laddr,
		raddr:   raddr,
	***REMOVED***, nil
***REMOVED***

// RFC 4254 7.2
type channelOpenDirectMsg struct ***REMOVED***
	raddr string
	rport uint32
	laddr string
	lport uint32
***REMOVED***

func (c *Client) dial(laddr string, lport int, raddr string, rport int) (Channel, error) ***REMOVED***
	msg := channelOpenDirectMsg***REMOVED***
		raddr: raddr,
		rport: uint32(rport),
		laddr: laddr,
		lport: uint32(lport),
	***REMOVED***
	ch, in, err := c.OpenChannel("direct-tcpip", Marshal(&msg))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	go DiscardRequests(in)
	return ch, err
***REMOVED***

type tcpChan struct ***REMOVED***
	Channel // the backing channel
***REMOVED***

// chanConn fulfills the net.Conn interface without
// the tcpChan having to hold laddr or raddr directly.
type chanConn struct ***REMOVED***
	Channel
	laddr, raddr net.Addr
***REMOVED***

// LocalAddr returns the local network address.
func (t *chanConn) LocalAddr() net.Addr ***REMOVED***
	return t.laddr
***REMOVED***

// RemoteAddr returns the remote network address.
func (t *chanConn) RemoteAddr() net.Addr ***REMOVED***
	return t.raddr
***REMOVED***

// SetDeadline sets the read and write deadlines associated
// with the connection.
func (t *chanConn) SetDeadline(deadline time.Time) error ***REMOVED***
	if err := t.SetReadDeadline(deadline); err != nil ***REMOVED***
		return err
	***REMOVED***
	return t.SetWriteDeadline(deadline)
***REMOVED***

// SetReadDeadline sets the read deadline.
// A zero value for t means Read will not time out.
// After the deadline, the error from Read will implement net.Error
// with Timeout() == true.
func (t *chanConn) SetReadDeadline(deadline time.Time) error ***REMOVED***
	// for compatibility with previous version,
	// the error message contains "tcpChan"
	return errors.New("ssh: tcpChan: deadline not supported")
***REMOVED***

// SetWriteDeadline exists to satisfy the net.Conn interface
// but is not implemented by this type.  It always returns an error.
func (t *chanConn) SetWriteDeadline(deadline time.Time) error ***REMOVED***
	return errors.New("ssh: tcpChan: deadline not supported")
***REMOVED***
