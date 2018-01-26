// +build !windows

package dns

import (
	"net"
	"syscall"
)

// SessionUDP holds the remote address and the associated
// out-of-band data.
type SessionUDP struct ***REMOVED***
	raddr   *net.UDPAddr
	context []byte
***REMOVED***

// RemoteAddr returns the remote network address.
func (s *SessionUDP) RemoteAddr() net.Addr ***REMOVED*** return s.raddr ***REMOVED***

// setUDPSocketOptions sets the UDP socket options.
// This function is implemented on a per platform basis. See udp_*.go for more details
func setUDPSocketOptions(conn *net.UDPConn) error ***REMOVED***
	sa, err := getUDPSocketName(conn)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	switch sa.(type) ***REMOVED***
	case *syscall.SockaddrInet6:
		v6only, err := getUDPSocketOptions6Only(conn)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		setUDPSocketOptions6(conn)
		if !v6only ***REMOVED***
			setUDPSocketOptions4(conn)
		***REMOVED***
	case *syscall.SockaddrInet4:
		setUDPSocketOptions4(conn)
	***REMOVED***
	return nil
***REMOVED***

// ReadFromSessionUDP acts just like net.UDPConn.ReadFrom(), but returns a session object instead of a
// net.UDPAddr.
func ReadFromSessionUDP(conn *net.UDPConn, b []byte) (int, *SessionUDP, error) ***REMOVED***
	oob := make([]byte, 40)
	n, oobn, _, raddr, err := conn.ReadMsgUDP(b, oob)
	if err != nil ***REMOVED***
		return n, nil, err
	***REMOVED***
	return n, &SessionUDP***REMOVED***raddr, oob[:oobn]***REMOVED***, err
***REMOVED***

// WriteToSessionUDP acts just like net.UDPConn.WritetTo(), but uses a *SessionUDP instead of a net.Addr.
func WriteToSessionUDP(conn *net.UDPConn, b []byte, session *SessionUDP) (int, error) ***REMOVED***
	n, _, err := conn.WriteMsgUDP(b, session.context, session.raddr)
	return n, err
***REMOVED***
