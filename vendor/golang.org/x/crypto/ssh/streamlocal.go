package ssh

import (
	"errors"
	"io"
	"net"
)

// streamLocalChannelOpenDirectMsg is a struct used for SSH_MSG_CHANNEL_OPEN message
// with "direct-streamlocal@openssh.com" string.
//
// See openssh-portable/PROTOCOL, section 2.4. connection: Unix domain socket forwarding
// https://github.com/openssh/openssh-portable/blob/master/PROTOCOL#L235
type streamLocalChannelOpenDirectMsg struct ***REMOVED***
	socketPath string
	reserved0  string
	reserved1  uint32
***REMOVED***

// forwardedStreamLocalPayload is a struct used for SSH_MSG_CHANNEL_OPEN message
// with "forwarded-streamlocal@openssh.com" string.
type forwardedStreamLocalPayload struct ***REMOVED***
	SocketPath string
	Reserved0  string
***REMOVED***

// streamLocalChannelForwardMsg is a struct used for SSH2_MSG_GLOBAL_REQUEST message
// with "streamlocal-forward@openssh.com"/"cancel-streamlocal-forward@openssh.com" string.
type streamLocalChannelForwardMsg struct ***REMOVED***
	socketPath string
***REMOVED***

// ListenUnix is similar to ListenTCP but uses a Unix domain socket.
func (c *Client) ListenUnix(socketPath string) (net.Listener, error) ***REMOVED***
	m := streamLocalChannelForwardMsg***REMOVED***
		socketPath,
	***REMOVED***
	// send message
	ok, _, err := c.SendRequest("streamlocal-forward@openssh.com", true, Marshal(&m))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if !ok ***REMOVED***
		return nil, errors.New("ssh: streamlocal-forward@openssh.com request denied by peer")
	***REMOVED***
	ch := c.forwards.add(&net.UnixAddr***REMOVED***Name: socketPath, Net: "unix"***REMOVED***)

	return &unixListener***REMOVED***socketPath, c, ch***REMOVED***, nil
***REMOVED***

func (c *Client) dialStreamLocal(socketPath string) (Channel, error) ***REMOVED***
	msg := streamLocalChannelOpenDirectMsg***REMOVED***
		socketPath: socketPath,
	***REMOVED***
	ch, in, err := c.OpenChannel("direct-streamlocal@openssh.com", Marshal(&msg))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	go DiscardRequests(in)
	return ch, err
***REMOVED***

type unixListener struct ***REMOVED***
	socketPath string

	conn *Client
	in   <-chan forward
***REMOVED***

// Accept waits for and returns the next connection to the listener.
func (l *unixListener) Accept() (net.Conn, error) ***REMOVED***
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
		laddr: &net.UnixAddr***REMOVED***
			Name: l.socketPath,
			Net:  "unix",
		***REMOVED***,
		raddr: &net.UnixAddr***REMOVED***
			Name: "@",
			Net:  "unix",
		***REMOVED***,
	***REMOVED***, nil
***REMOVED***

// Close closes the listener.
func (l *unixListener) Close() error ***REMOVED***
	// this also closes the listener.
	l.conn.forwards.remove(&net.UnixAddr***REMOVED***Name: l.socketPath, Net: "unix"***REMOVED***)
	m := streamLocalChannelForwardMsg***REMOVED***
		l.socketPath,
	***REMOVED***
	ok, _, err := l.conn.SendRequest("cancel-streamlocal-forward@openssh.com", true, Marshal(&m))
	if err == nil && !ok ***REMOVED***
		err = errors.New("ssh: cancel-streamlocal-forward@openssh.com failed")
	***REMOVED***
	return err
***REMOVED***

// Addr returns the listener's network address.
func (l *unixListener) Addr() net.Addr ***REMOVED***
	return &net.UnixAddr***REMOVED***
		Name: l.socketPath,
		Net:  "unix",
	***REMOVED***
***REMOVED***
