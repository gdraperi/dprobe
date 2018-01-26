package memberlist

import (
	"fmt"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/armon/go-metrics"
	sockaddr "github.com/hashicorp/go-sockaddr"
)

const (
	// udpPacketBufSize is used to buffer incoming packets during read
	// operations.
	udpPacketBufSize = 65536

	// udpRecvBufSize is a large buffer size that we attempt to set UDP
	// sockets to in order to handle a large volume of messages.
	udpRecvBufSize = 2 * 1024 * 1024
)

// NetTransportConfig is used to configure a net transport.
type NetTransportConfig struct ***REMOVED***
	// BindAddrs is a list of addresses to bind to for both TCP and UDP
	// communications.
	BindAddrs []string

	// BindPort is the port to listen on, for each address above.
	BindPort int

	// Logger is a logger for operator messages.
	Logger *log.Logger
***REMOVED***

// NetTransport is a Transport implementation that uses connectionless UDP for
// packet operations, and ad-hoc TCP connections for stream operations.
type NetTransport struct ***REMOVED***
	config       *NetTransportConfig
	packetCh     chan *Packet
	streamCh     chan net.Conn
	logger       *log.Logger
	wg           sync.WaitGroup
	tcpListeners []*net.TCPListener
	udpListeners []*net.UDPConn
	shutdown     int32
***REMOVED***

// NewNetTransport returns a net transport with the given configuration. On
// success all the network listeners will be created and listening.
func NewNetTransport(config *NetTransportConfig) (*NetTransport, error) ***REMOVED***
	// If we reject the empty list outright we can assume that there's at
	// least one listener of each type later during operation.
	if len(config.BindAddrs) == 0 ***REMOVED***
		return nil, fmt.Errorf("At least one bind address is required")
	***REMOVED***

	// Build out the new transport.
	var ok bool
	t := NetTransport***REMOVED***
		config:   config,
		packetCh: make(chan *Packet),
		streamCh: make(chan net.Conn),
		logger:   config.Logger,
	***REMOVED***

	// Clean up listeners if there's an error.
	defer func() ***REMOVED***
		if !ok ***REMOVED***
			t.Shutdown()
		***REMOVED***
	***REMOVED***()

	// Build all the TCP and UDP listeners.
	port := config.BindPort
	for _, addr := range config.BindAddrs ***REMOVED***
		ip := net.ParseIP(addr)

		tcpAddr := &net.TCPAddr***REMOVED***IP: ip, Port: port***REMOVED***
		tcpLn, err := net.ListenTCP("tcp", tcpAddr)
		if err != nil ***REMOVED***
			return nil, fmt.Errorf("Failed to start TCP listener on %q port %d: %v", addr, port, err)
		***REMOVED***
		t.tcpListeners = append(t.tcpListeners, tcpLn)

		// If the config port given was zero, use the first TCP listener
		// to pick an available port and then apply that to everything
		// else.
		if port == 0 ***REMOVED***
			port = tcpLn.Addr().(*net.TCPAddr).Port
		***REMOVED***

		udpAddr := &net.UDPAddr***REMOVED***IP: ip, Port: port***REMOVED***
		udpLn, err := net.ListenUDP("udp", udpAddr)
		if err != nil ***REMOVED***
			return nil, fmt.Errorf("Failed to start UDP listener on %q port %d: %v", addr, port, err)
		***REMOVED***
		if err := setUDPRecvBuf(udpLn); err != nil ***REMOVED***
			return nil, fmt.Errorf("Failed to resize UDP buffer: %v", err)
		***REMOVED***
		t.udpListeners = append(t.udpListeners, udpLn)
	***REMOVED***

	// Fire them up now that we've been able to create them all.
	for i := 0; i < len(config.BindAddrs); i++ ***REMOVED***
		t.wg.Add(2)
		go t.tcpListen(t.tcpListeners[i])
		go t.udpListen(t.udpListeners[i])
	***REMOVED***

	ok = true
	return &t, nil
***REMOVED***

// GetAutoBindPort returns the bind port that was automatically given by the
// kernel, if a bind port of 0 was given.
func (t *NetTransport) GetAutoBindPort() int ***REMOVED***
	// We made sure there's at least one TCP listener, and that one's
	// port was applied to all the others for the dynamic bind case.
	return t.tcpListeners[0].Addr().(*net.TCPAddr).Port
***REMOVED***

// See Transport.
func (t *NetTransport) FinalAdvertiseAddr(ip string, port int) (net.IP, int, error) ***REMOVED***
	var advertiseAddr net.IP
	var advertisePort int
	if ip != "" ***REMOVED***
		// If they've supplied an address, use that.
		advertiseAddr = net.ParseIP(ip)
		if advertiseAddr == nil ***REMOVED***
			return nil, 0, fmt.Errorf("Failed to parse advertise address %q", ip)
		***REMOVED***

		// Ensure IPv4 conversion if necessary.
		if ip4 := advertiseAddr.To4(); ip4 != nil ***REMOVED***
			advertiseAddr = ip4
		***REMOVED***
		advertisePort = port
	***REMOVED*** else ***REMOVED***
		if t.config.BindAddrs[0] == "0.0.0.0" ***REMOVED***
			// Otherwise, if we're not bound to a specific IP, let's
			// use a suitable private IP address.
			var err error
			ip, err = sockaddr.GetPrivateIP()
			if err != nil ***REMOVED***
				return nil, 0, fmt.Errorf("Failed to get interface addresses: %v", err)
			***REMOVED***
			if ip == "" ***REMOVED***
				return nil, 0, fmt.Errorf("No private IP address found, and explicit IP not provided")
			***REMOVED***

			advertiseAddr = net.ParseIP(ip)
			if advertiseAddr == nil ***REMOVED***
				return nil, 0, fmt.Errorf("Failed to parse advertise address: %q", ip)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			// Use the IP that we're bound to, based on the first
			// TCP listener, which we already ensure is there.
			advertiseAddr = t.tcpListeners[0].Addr().(*net.TCPAddr).IP
		***REMOVED***

		// Use the port we are bound to.
		advertisePort = t.GetAutoBindPort()
	***REMOVED***

	return advertiseAddr, advertisePort, nil
***REMOVED***

// See Transport.
func (t *NetTransport) WriteTo(b []byte, addr string) (time.Time, error) ***REMOVED***
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil ***REMOVED***
		return time.Time***REMOVED******REMOVED***, err
	***REMOVED***

	// We made sure there's at least one UDP listener, so just use the
	// packet sending interface on the first one. Take the time after the
	// write call comes back, which will underestimate the time a little,
	// but help account for any delays before the write occurs.
	_, err = t.udpListeners[0].WriteTo(b, udpAddr)
	return time.Now(), err
***REMOVED***

// See Transport.
func (t *NetTransport) PacketCh() <-chan *Packet ***REMOVED***
	return t.packetCh
***REMOVED***

// See Transport.
func (t *NetTransport) DialTimeout(addr string, timeout time.Duration) (net.Conn, error) ***REMOVED***
	dialer := net.Dialer***REMOVED***Timeout: timeout***REMOVED***
	return dialer.Dial("tcp", addr)
***REMOVED***

// See Transport.
func (t *NetTransport) StreamCh() <-chan net.Conn ***REMOVED***
	return t.streamCh
***REMOVED***

// See Transport.
func (t *NetTransport) Shutdown() error ***REMOVED***
	// This will avoid log spam about errors when we shut down.
	atomic.StoreInt32(&t.shutdown, 1)

	// Rip through all the connections and shut them down.
	for _, conn := range t.tcpListeners ***REMOVED***
		conn.Close()
	***REMOVED***
	for _, conn := range t.udpListeners ***REMOVED***
		conn.Close()
	***REMOVED***

	// Block until all the listener threads have died.
	t.wg.Wait()
	return nil
***REMOVED***

// tcpListen is a long running goroutine that accepts incoming TCP connections
// and hands them off to the stream channel.
func (t *NetTransport) tcpListen(tcpLn *net.TCPListener) ***REMOVED***
	defer t.wg.Done()
	for ***REMOVED***
		conn, err := tcpLn.AcceptTCP()
		if err != nil ***REMOVED***
			if s := atomic.LoadInt32(&t.shutdown); s == 1 ***REMOVED***
				break
			***REMOVED***

			t.logger.Printf("[ERR] memberlist: Error accepting TCP connection: %v", err)
			continue
		***REMOVED***

		t.streamCh <- conn
	***REMOVED***
***REMOVED***

// udpListen is a long running goroutine that accepts incoming UDP packets and
// hands them off to the packet channel.
func (t *NetTransport) udpListen(udpLn *net.UDPConn) ***REMOVED***
	defer t.wg.Done()
	for ***REMOVED***
		// Do a blocking read into a fresh buffer. Grab a time stamp as
		// close as possible to the I/O.
		buf := make([]byte, udpPacketBufSize)
		n, addr, err := udpLn.ReadFrom(buf)
		ts := time.Now()
		if err != nil ***REMOVED***
			if s := atomic.LoadInt32(&t.shutdown); s == 1 ***REMOVED***
				break
			***REMOVED***

			t.logger.Printf("[ERR] memberlist: Error reading UDP packet: %v", err)
			continue
		***REMOVED***

		// Check the length - it needs to have at least one byte to be a
		// proper message.
		if n < 1 ***REMOVED***
			t.logger.Printf("[ERR] memberlist: UDP packet too short (%d bytes) %s",
				len(buf), LogAddress(addr))
			continue
		***REMOVED***

		// Ingest the packet.
		metrics.IncrCounter([]string***REMOVED***"memberlist", "udp", "received"***REMOVED***, float32(n))
		t.packetCh <- &Packet***REMOVED***
			Buf:       buf[:n],
			From:      addr,
			Timestamp: ts,
		***REMOVED***
	***REMOVED***
***REMOVED***

// setUDPRecvBuf is used to resize the UDP receive window. The function
// attempts to set the read buffer to `udpRecvBuf` but backs off until
// the read buffer can be set.
func setUDPRecvBuf(c *net.UDPConn) error ***REMOVED***
	size := udpRecvBufSize
	var err error
	for size > 0 ***REMOVED***
		if err = c.SetReadBuffer(size); err == nil ***REMOVED***
			return nil
		***REMOVED***
		size = size / 2
	***REMOVED***
	return err
***REMOVED***
