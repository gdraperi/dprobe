// +build windows

package dns

import "net"

type SessionUDP struct ***REMOVED***
	raddr *net.UDPAddr
***REMOVED***

// ReadFromSessionUDP acts just like net.UDPConn.ReadFrom(), but returns a session object instead of a
// net.UDPAddr.
func ReadFromSessionUDP(conn *net.UDPConn, b []byte) (int, *SessionUDP, error) ***REMOVED***
	n, raddr, err := conn.ReadFrom(b)
	if err != nil ***REMOVED***
		return n, nil, err
	***REMOVED***
	session := &SessionUDP***REMOVED***raddr.(*net.UDPAddr)***REMOVED***
	return n, session, err
***REMOVED***

// WriteToSessionUDP acts just like net.UDPConn.WritetTo(), but uses a *SessionUDP instead of a net.Addr.
func WriteToSessionUDP(conn *net.UDPConn, b []byte, session *SessionUDP) (int, error) ***REMOVED***
	n, err := conn.WriteTo(b, session.raddr)
	return n, err
***REMOVED***

func (s *SessionUDP) RemoteAddr() net.Addr ***REMOVED*** return s.raddr ***REMOVED***

// setUDPSocketOptions sets the UDP socket options.
// This function is implemented on a per platform basis. See udp_*.go for more details
func setUDPSocketOptions(conn *net.UDPConn) error ***REMOVED***
	return nil
***REMOVED***
