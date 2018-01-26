package memberlist

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

// MockNetwork is used as a factory that produces MockTransport instances which
// are uniquely addressed and wired up to talk to each other.
type MockNetwork struct ***REMOVED***
	transports map[string]*MockTransport
	port       int
***REMOVED***

// NewTransport returns a new MockTransport with a unique address, wired up to
// talk to the other transports in the MockNetwork.
func (n *MockNetwork) NewTransport() *MockTransport ***REMOVED***
	n.port += 1
	addr := fmt.Sprintf("127.0.0.1:%d", n.port)
	transport := &MockTransport***REMOVED***
		net:      n,
		addr:     &MockAddress***REMOVED***addr***REMOVED***,
		packetCh: make(chan *Packet),
		streamCh: make(chan net.Conn),
	***REMOVED***

	if n.transports == nil ***REMOVED***
		n.transports = make(map[string]*MockTransport)
	***REMOVED***
	n.transports[addr] = transport
	return transport
***REMOVED***

// MockAddress is a wrapper which adds the net.Addr interface to our mock
// address scheme.
type MockAddress struct ***REMOVED***
	addr string
***REMOVED***

// See net.Addr.
func (a *MockAddress) Network() string ***REMOVED***
	return "mock"
***REMOVED***

// See net.Addr.
func (a *MockAddress) String() string ***REMOVED***
	return a.addr
***REMOVED***

// MockTransport directly plumbs messages to other transports its MockNetwork.
type MockTransport struct ***REMOVED***
	net      *MockNetwork
	addr     *MockAddress
	packetCh chan *Packet
	streamCh chan net.Conn
***REMOVED***

// See Transport.
func (t *MockTransport) FinalAdvertiseAddr(string, int) (net.IP, int, error) ***REMOVED***
	host, portStr, err := net.SplitHostPort(t.addr.String())
	if err != nil ***REMOVED***
		return nil, 0, err
	***REMOVED***

	ip := net.ParseIP(host)
	if ip == nil ***REMOVED***
		return nil, 0, fmt.Errorf("Failed to parse IP %q", host)
	***REMOVED***

	port, err := strconv.ParseInt(portStr, 10, 16)
	if err != nil ***REMOVED***
		return nil, 0, err
	***REMOVED***

	return ip, int(port), nil
***REMOVED***

// See Transport.
func (t *MockTransport) WriteTo(b []byte, addr string) (time.Time, error) ***REMOVED***
	dest, ok := t.net.transports[addr]
	if !ok ***REMOVED***
		return time.Time***REMOVED******REMOVED***, fmt.Errorf("No route to %q", addr)
	***REMOVED***

	now := time.Now()
	dest.packetCh <- &Packet***REMOVED***
		Buf:       b,
		From:      t.addr,
		Timestamp: now,
	***REMOVED***
	return now, nil
***REMOVED***

// See Transport.
func (t *MockTransport) PacketCh() <-chan *Packet ***REMOVED***
	return t.packetCh
***REMOVED***

// See Transport.
func (t *MockTransport) DialTimeout(addr string, timeout time.Duration) (net.Conn, error) ***REMOVED***
	dest, ok := t.net.transports[addr]
	if !ok ***REMOVED***
		return nil, fmt.Errorf("No route to %q", addr)
	***REMOVED***

	p1, p2 := net.Pipe()
	dest.streamCh <- p1
	return p2, nil
***REMOVED***

// See Transport.
func (t *MockTransport) StreamCh() <-chan net.Conn ***REMOVED***
	return t.streamCh
***REMOVED***

// See Transport.
func (t *MockTransport) Shutdown() error ***REMOVED***
	return nil
***REMOVED***
