// DNS server implementation.

package dns

import (
	"bytes"
	"io"
	"net"
	"sync"
	"time"
)

// Maximum number of TCP queries before we close the socket.
const maxTCPQueries = 128

// Handler is implemented by any value that implements ServeDNS.
type Handler interface ***REMOVED***
	ServeDNS(w ResponseWriter, r *Msg)
***REMOVED***

// A ResponseWriter interface is used by an DNS handler to
// construct an DNS response.
type ResponseWriter interface ***REMOVED***
	// LocalAddr returns the net.Addr of the server
	LocalAddr() net.Addr
	// RemoteAddr returns the net.Addr of the client that sent the current request.
	RemoteAddr() net.Addr
	// WriteMsg writes a reply back to the client.
	WriteMsg(*Msg) error
	// Write writes a raw buffer back to the client.
	Write([]byte) (int, error)
	// Close closes the connection.
	Close() error
	// TsigStatus returns the status of the Tsig.
	TsigStatus() error
	// TsigTimersOnly sets the tsig timers only boolean.
	TsigTimersOnly(bool)
	// Hijack lets the caller take over the connection.
	// After a call to Hijack(), the DNS package will not do anything with the connection.
	Hijack()
***REMOVED***

type response struct ***REMOVED***
	hijacked       bool // connection has been hijacked by handler
	tsigStatus     error
	tsigTimersOnly bool
	tsigRequestMAC string
	tsigSecret     map[string]string // the tsig secrets
	udp            *net.UDPConn      // i/o connection if UDP was used
	tcp            *net.TCPConn      // i/o connection if TCP was used
	udpSession     *SessionUDP       // oob data to get egress interface right
	remoteAddr     net.Addr          // address of the client
	writer         Writer            // writer to output the raw DNS bits
***REMOVED***

// ServeMux is an DNS request multiplexer. It matches the
// zone name of each incoming request against a list of
// registered patterns add calls the handler for the pattern
// that most closely matches the zone name. ServeMux is DNSSEC aware, meaning
// that queries for the DS record are redirected to the parent zone (if that
// is also registered), otherwise the child gets the query.
// ServeMux is also safe for concurrent access from multiple goroutines.
type ServeMux struct ***REMOVED***
	z map[string]Handler
	m *sync.RWMutex
***REMOVED***

// NewServeMux allocates and returns a new ServeMux.
func NewServeMux() *ServeMux ***REMOVED*** return &ServeMux***REMOVED***z: make(map[string]Handler), m: new(sync.RWMutex)***REMOVED*** ***REMOVED***

// DefaultServeMux is the default ServeMux used by Serve.
var DefaultServeMux = NewServeMux()

// The HandlerFunc type is an adapter to allow the use of
// ordinary functions as DNS handlers.  If f is a function
// with the appropriate signature, HandlerFunc(f) is a
// Handler object that calls f.
type HandlerFunc func(ResponseWriter, *Msg)

// ServeDNS calls f(w, r).
func (f HandlerFunc) ServeDNS(w ResponseWriter, r *Msg) ***REMOVED***
	f(w, r)
***REMOVED***

// HandleFailed returns a HandlerFunc that returns SERVFAIL for every request it gets.
func HandleFailed(w ResponseWriter, r *Msg) ***REMOVED***
	m := new(Msg)
	m.SetRcode(r, RcodeServerFailure)
	// does not matter if this write fails
	w.WriteMsg(m)
***REMOVED***

func failedHandler() Handler ***REMOVED*** return HandlerFunc(HandleFailed) ***REMOVED***

// ListenAndServe Starts a server on addresss and network speficied. Invoke handler
// for incoming queries.
func ListenAndServe(addr string, network string, handler Handler) error ***REMOVED***
	server := &Server***REMOVED***Addr: addr, Net: network, Handler: handler***REMOVED***
	return server.ListenAndServe()
***REMOVED***

// ActivateAndServe activates a server with a listener from systemd,
// l and p should not both be non-nil.
// If both l and p are not nil only p will be used.
// Invoke handler for incoming queries.
func ActivateAndServe(l net.Listener, p net.PacketConn, handler Handler) error ***REMOVED***
	server := &Server***REMOVED***Listener: l, PacketConn: p, Handler: handler***REMOVED***
	return server.ActivateAndServe()
***REMOVED***

func (mux *ServeMux) match(q string, t uint16) Handler ***REMOVED***
	mux.m.RLock()
	defer mux.m.RUnlock()
	var handler Handler
	b := make([]byte, len(q)) // worst case, one label of length q
	off := 0
	end := false
	for ***REMOVED***
		l := len(q[off:])
		for i := 0; i < l; i++ ***REMOVED***
			b[i] = q[off+i]
			if b[i] >= 'A' && b[i] <= 'Z' ***REMOVED***
				b[i] |= ('a' - 'A')
			***REMOVED***
		***REMOVED***
		if h, ok := mux.z[string(b[:l])]; ok ***REMOVED*** // 'causes garbage, might want to change the map key
			if t != TypeDS ***REMOVED***
				return h
			***REMOVED***
			// Continue for DS to see if we have a parent too, if so delegeate to the parent
			handler = h
		***REMOVED***
		off, end = NextLabel(q, off)
		if end ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	// Wildcard match, if we have found nothing try the root zone as a last resort.
	if h, ok := mux.z["."]; ok ***REMOVED***
		return h
	***REMOVED***
	return handler
***REMOVED***

// Handle adds a handler to the ServeMux for pattern.
func (mux *ServeMux) Handle(pattern string, handler Handler) ***REMOVED***
	if pattern == "" ***REMOVED***
		panic("dns: invalid pattern " + pattern)
	***REMOVED***
	mux.m.Lock()
	mux.z[Fqdn(pattern)] = handler
	mux.m.Unlock()
***REMOVED***

// HandleFunc adds a handler function to the ServeMux for pattern.
func (mux *ServeMux) HandleFunc(pattern string, handler func(ResponseWriter, *Msg)) ***REMOVED***
	mux.Handle(pattern, HandlerFunc(handler))
***REMOVED***

// HandleRemove deregistrars the handler specific for pattern from the ServeMux.
func (mux *ServeMux) HandleRemove(pattern string) ***REMOVED***
	if pattern == "" ***REMOVED***
		panic("dns: invalid pattern " + pattern)
	***REMOVED***
	mux.m.Lock()
	delete(mux.z, Fqdn(pattern))
	mux.m.Unlock()
***REMOVED***

// ServeDNS dispatches the request to the handler whose
// pattern most closely matches the request message. If DefaultServeMux
// is used the correct thing for DS queries is done: a possible parent
// is sought.
// If no handler is found a standard SERVFAIL message is returned
// If the request message does not have exactly one question in the
// question section a SERVFAIL is returned, unlesss Unsafe is true.
func (mux *ServeMux) ServeDNS(w ResponseWriter, request *Msg) ***REMOVED***
	var h Handler
	if len(request.Question) < 1 ***REMOVED*** // allow more than one question
		h = failedHandler()
	***REMOVED*** else ***REMOVED***
		if h = mux.match(request.Question[0].Name, request.Question[0].Qtype); h == nil ***REMOVED***
			h = failedHandler()
		***REMOVED***
	***REMOVED***
	h.ServeDNS(w, request)
***REMOVED***

// Handle registers the handler with the given pattern
// in the DefaultServeMux. The documentation for
// ServeMux explains how patterns are matched.
func Handle(pattern string, handler Handler) ***REMOVED*** DefaultServeMux.Handle(pattern, handler) ***REMOVED***

// HandleRemove deregisters the handle with the given pattern
// in the DefaultServeMux.
func HandleRemove(pattern string) ***REMOVED*** DefaultServeMux.HandleRemove(pattern) ***REMOVED***

// HandleFunc registers the handler function with the given pattern
// in the DefaultServeMux.
func HandleFunc(pattern string, handler func(ResponseWriter, *Msg)) ***REMOVED***
	DefaultServeMux.HandleFunc(pattern, handler)
***REMOVED***

// Writer writes raw DNS messages; each call to Write should send an entire message.
type Writer interface ***REMOVED***
	io.Writer
***REMOVED***

// Reader reads raw DNS messages; each call to ReadTCP or ReadUDP should return an entire message.
type Reader interface ***REMOVED***
	// ReadTCP reads a raw message from a TCP connection. Implementations may alter
	// connection properties, for example the read-deadline.
	ReadTCP(conn *net.TCPConn, timeout time.Duration) ([]byte, error)
	// ReadUDP reads a raw message from a UDP connection. Implementations may alter
	// connection properties, for example the read-deadline.
	ReadUDP(conn *net.UDPConn, timeout time.Duration) ([]byte, *SessionUDP, error)
***REMOVED***

// defaultReader is an adapter for the Server struct that implements the Reader interface
// using the readTCP and readUDP func of the embedded Server.
type defaultReader struct ***REMOVED***
	*Server
***REMOVED***

func (dr *defaultReader) ReadTCP(conn *net.TCPConn, timeout time.Duration) ([]byte, error) ***REMOVED***
	return dr.readTCP(conn, timeout)
***REMOVED***

func (dr *defaultReader) ReadUDP(conn *net.UDPConn, timeout time.Duration) ([]byte, *SessionUDP, error) ***REMOVED***
	return dr.readUDP(conn, timeout)
***REMOVED***

// DecorateReader is a decorator hook for extending or supplanting the functionality of a Reader.
// Implementations should never return a nil Reader.
type DecorateReader func(Reader) Reader

// DecorateWriter is a decorator hook for extending or supplanting the functionality of a Writer.
// Implementations should never return a nil Writer.
type DecorateWriter func(Writer) Writer

// A Server defines parameters for running an DNS server.
type Server struct ***REMOVED***
	// Address to listen on, ":dns" if empty.
	Addr string
	// if "tcp" it will invoke a TCP listener, otherwise an UDP one.
	Net string
	// TCP Listener to use, this is to aid in systemd's socket activation.
	Listener net.Listener
	// UDP "Listener" to use, this is to aid in systemd's socket activation.
	PacketConn net.PacketConn
	// Handler to invoke, dns.DefaultServeMux if nil.
	Handler Handler
	// Default buffer size to use to read incoming UDP messages. If not set
	// it defaults to MinMsgSize (512 B).
	UDPSize int
	// The net.Conn.SetReadTimeout value for new connections, defaults to 2 * time.Second.
	ReadTimeout time.Duration
	// The net.Conn.SetWriteTimeout value for new connections, defaults to 2 * time.Second.
	WriteTimeout time.Duration
	// TCP idle timeout for multiple queries, if nil, defaults to 8 * time.Second (RFC 5966).
	IdleTimeout func() time.Duration
	// Secret(s) for Tsig map[<zonename>]<base64 secret>.
	TsigSecret map[string]string
	// Unsafe instructs the server to disregard any sanity checks and directly hand the message to
	// the handler. It will specfically not check if the query has the QR bit not set.
	Unsafe bool
	// If NotifyStartedFunc is set it is called once the server has started listening.
	NotifyStartedFunc func()
	// DecorateReader is optional, allows customization of the process that reads raw DNS messages.
	DecorateReader DecorateReader
	// DecorateWriter is optional, allows customization of the process that writes raw DNS messages.
	DecorateWriter DecorateWriter

	// Graceful shutdown handling

	inFlight sync.WaitGroup

	lock    sync.RWMutex
	started bool
***REMOVED***

// ListenAndServe starts a nameserver on the configured address in *Server.
func (srv *Server) ListenAndServe() error ***REMOVED***
	srv.lock.Lock()
	defer srv.lock.Unlock()
	if srv.started ***REMOVED***
		return &Error***REMOVED***err: "server already started"***REMOVED***
	***REMOVED***
	addr := srv.Addr
	if addr == "" ***REMOVED***
		addr = ":domain"
	***REMOVED***
	if srv.UDPSize == 0 ***REMOVED***
		srv.UDPSize = MinMsgSize
	***REMOVED***
	switch srv.Net ***REMOVED***
	case "tcp", "tcp4", "tcp6":
		a, e := net.ResolveTCPAddr(srv.Net, addr)
		if e != nil ***REMOVED***
			return e
		***REMOVED***
		l, e := net.ListenTCP(srv.Net, a)
		if e != nil ***REMOVED***
			return e
		***REMOVED***
		srv.Listener = l
		srv.started = true
		srv.lock.Unlock()
		e = srv.serveTCP(l)
		srv.lock.Lock() // to satisfy the defer at the top
		return e
	case "udp", "udp4", "udp6":
		a, e := net.ResolveUDPAddr(srv.Net, addr)
		if e != nil ***REMOVED***
			return e
		***REMOVED***
		l, e := net.ListenUDP(srv.Net, a)
		if e != nil ***REMOVED***
			return e
		***REMOVED***
		if e := setUDPSocketOptions(l); e != nil ***REMOVED***
			return e
		***REMOVED***
		srv.PacketConn = l
		srv.started = true
		srv.lock.Unlock()
		e = srv.serveUDP(l)
		srv.lock.Lock() // to satisfy the defer at the top
		return e
	***REMOVED***
	return &Error***REMOVED***err: "bad network"***REMOVED***
***REMOVED***

// ActivateAndServe starts a nameserver with the PacketConn or Listener
// configured in *Server. Its main use is to start a server from systemd.
func (srv *Server) ActivateAndServe() error ***REMOVED***
	srv.lock.Lock()
	defer srv.lock.Unlock()
	if srv.started ***REMOVED***
		return &Error***REMOVED***err: "server already started"***REMOVED***
	***REMOVED***
	pConn := srv.PacketConn
	l := srv.Listener
	if pConn != nil ***REMOVED***
		if srv.UDPSize == 0 ***REMOVED***
			srv.UDPSize = MinMsgSize
		***REMOVED***
		if t, ok := pConn.(*net.UDPConn); ok ***REMOVED***
			if e := setUDPSocketOptions(t); e != nil ***REMOVED***
				return e
			***REMOVED***
			srv.started = true
			srv.lock.Unlock()
			e := srv.serveUDP(t)
			srv.lock.Lock() // to satisfy the defer at the top
			return e
		***REMOVED***
	***REMOVED***
	if l != nil ***REMOVED***
		if t, ok := l.(*net.TCPListener); ok ***REMOVED***
			srv.started = true
			srv.lock.Unlock()
			e := srv.serveTCP(t)
			srv.lock.Lock() // to satisfy the defer at the top
			return e
		***REMOVED***
	***REMOVED***
	return &Error***REMOVED***err: "bad listeners"***REMOVED***
***REMOVED***

// Shutdown gracefully shuts down a server. After a call to Shutdown, ListenAndServe and
// ActivateAndServe will return. All in progress queries are completed before the server
// is taken down. If the Shutdown is taking longer than the reading timeout an error
// is returned.
func (srv *Server) Shutdown() error ***REMOVED***
	srv.lock.Lock()
	if !srv.started ***REMOVED***
		srv.lock.Unlock()
		return &Error***REMOVED***err: "server not started"***REMOVED***
	***REMOVED***
	srv.started = false
	srv.lock.Unlock()

	if srv.PacketConn != nil ***REMOVED***
		srv.PacketConn.Close()
	***REMOVED***
	if srv.Listener != nil ***REMOVED***
		srv.Listener.Close()
	***REMOVED***

	fin := make(chan bool)
	go func() ***REMOVED***
		srv.inFlight.Wait()
		fin <- true
	***REMOVED***()

	select ***REMOVED***
	case <-time.After(srv.getReadTimeout()):
		return &Error***REMOVED***err: "server shutdown is pending"***REMOVED***
	case <-fin:
		return nil
	***REMOVED***
***REMOVED***

// getReadTimeout is a helper func to use system timeout if server did not intend to change it.
func (srv *Server) getReadTimeout() time.Duration ***REMOVED***
	rtimeout := dnsTimeout
	if srv.ReadTimeout != 0 ***REMOVED***
		rtimeout = srv.ReadTimeout
	***REMOVED***
	return rtimeout
***REMOVED***

// serveTCP starts a TCP listener for the server.
// Each request is handled in a separate goroutine.
func (srv *Server) serveTCP(l *net.TCPListener) error ***REMOVED***
	defer l.Close()

	if srv.NotifyStartedFunc != nil ***REMOVED***
		srv.NotifyStartedFunc()
	***REMOVED***

	reader := Reader(&defaultReader***REMOVED***srv***REMOVED***)
	if srv.DecorateReader != nil ***REMOVED***
		reader = srv.DecorateReader(reader)
	***REMOVED***

	handler := srv.Handler
	if handler == nil ***REMOVED***
		handler = DefaultServeMux
	***REMOVED***
	rtimeout := srv.getReadTimeout()
	// deadline is not used here
	for ***REMOVED***
		rw, e := l.AcceptTCP()
		if e != nil ***REMOVED***
			if neterr, ok := e.(net.Error); ok && neterr.Temporary() ***REMOVED***
				continue
			***REMOVED***
			return e
		***REMOVED***
		m, e := reader.ReadTCP(rw, rtimeout)
		srv.lock.RLock()
		if !srv.started ***REMOVED***
			srv.lock.RUnlock()
			return nil
		***REMOVED***
		srv.lock.RUnlock()
		if e != nil ***REMOVED***
			continue
		***REMOVED***
		srv.inFlight.Add(1)
		go srv.serve(rw.RemoteAddr(), handler, m, nil, nil, rw)
	***REMOVED***
***REMOVED***

// serveUDP starts a UDP listener for the server.
// Each request is handled in a separate goroutine.
func (srv *Server) serveUDP(l *net.UDPConn) error ***REMOVED***
	defer l.Close()

	if srv.NotifyStartedFunc != nil ***REMOVED***
		srv.NotifyStartedFunc()
	***REMOVED***

	reader := Reader(&defaultReader***REMOVED***srv***REMOVED***)
	if srv.DecorateReader != nil ***REMOVED***
		reader = srv.DecorateReader(reader)
	***REMOVED***

	handler := srv.Handler
	if handler == nil ***REMOVED***
		handler = DefaultServeMux
	***REMOVED***
	rtimeout := srv.getReadTimeout()
	// deadline is not used here
	for ***REMOVED***
		m, s, e := reader.ReadUDP(l, rtimeout)
		srv.lock.RLock()
		if !srv.started ***REMOVED***
			srv.lock.RUnlock()
			return nil
		***REMOVED***
		srv.lock.RUnlock()
		if e != nil ***REMOVED***
			continue
		***REMOVED***
		srv.inFlight.Add(1)
		go srv.serve(s.RemoteAddr(), handler, m, l, s, nil)
	***REMOVED***
***REMOVED***

// Serve a new connection.
func (srv *Server) serve(a net.Addr, h Handler, m []byte, u *net.UDPConn, s *SessionUDP, t *net.TCPConn) ***REMOVED***
	defer srv.inFlight.Done()

	w := &response***REMOVED***tsigSecret: srv.TsigSecret, udp: u, tcp: t, remoteAddr: a, udpSession: s***REMOVED***
	if srv.DecorateWriter != nil ***REMOVED***
		w.writer = srv.DecorateWriter(w)
	***REMOVED*** else ***REMOVED***
		w.writer = w
	***REMOVED***

	q := 0 // counter for the amount of TCP queries we get

	reader := Reader(&defaultReader***REMOVED***srv***REMOVED***)
	if srv.DecorateReader != nil ***REMOVED***
		reader = srv.DecorateReader(reader)
	***REMOVED***
Redo:
	req := new(Msg)
	err := req.Unpack(m)
	if err != nil ***REMOVED*** // Send a FormatError back
		x := new(Msg)
		x.SetRcodeFormatError(req)
		w.WriteMsg(x)
		goto Exit
	***REMOVED***
	if !srv.Unsafe && req.Response ***REMOVED***
		goto Exit
	***REMOVED***

	w.tsigStatus = nil
	if w.tsigSecret != nil ***REMOVED***
		if t := req.IsTsig(); t != nil ***REMOVED***
			secret := t.Hdr.Name
			if _, ok := w.tsigSecret[secret]; !ok ***REMOVED***
				w.tsigStatus = ErrKeyAlg
			***REMOVED***
			w.tsigStatus = TsigVerify(m, w.tsigSecret[secret], "", false)
			w.tsigTimersOnly = false
			w.tsigRequestMAC = req.Extra[len(req.Extra)-1].(*TSIG).MAC
		***REMOVED***
	***REMOVED***
	h.ServeDNS(w, req) // Writes back to the client

Exit:
	if w.tcp == nil ***REMOVED***
		return
	***REMOVED***
	// TODO(miek): make this number configurable?
	if q > maxTCPQueries ***REMOVED*** // close socket after this many queries
		w.Close()
		return
	***REMOVED***

	if w.hijacked ***REMOVED***
		return // client calls Close()
	***REMOVED***
	if u != nil ***REMOVED*** // UDP, "close" and return
		w.Close()
		return
	***REMOVED***
	idleTimeout := tcpIdleTimeout
	if srv.IdleTimeout != nil ***REMOVED***
		idleTimeout = srv.IdleTimeout()
	***REMOVED***
	m, e := reader.ReadTCP(w.tcp, idleTimeout)
	if e == nil ***REMOVED***
		q++
		goto Redo
	***REMOVED***
	w.Close()
	return
***REMOVED***

func (srv *Server) readTCP(conn *net.TCPConn, timeout time.Duration) ([]byte, error) ***REMOVED***
	conn.SetReadDeadline(time.Now().Add(timeout))
	l := make([]byte, 2)
	n, err := conn.Read(l)
	if err != nil || n != 2 ***REMOVED***
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return nil, ErrShortRead
	***REMOVED***
	length, _ := unpackUint16(l, 0)
	if length == 0 ***REMOVED***
		return nil, ErrShortRead
	***REMOVED***
	m := make([]byte, int(length))
	n, err = conn.Read(m[:int(length)])
	if err != nil || n == 0 ***REMOVED***
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return nil, ErrShortRead
	***REMOVED***
	i := n
	for i < int(length) ***REMOVED***
		j, err := conn.Read(m[i:int(length)])
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		i += j
	***REMOVED***
	n = i
	m = m[:n]
	return m, nil
***REMOVED***

func (srv *Server) readUDP(conn *net.UDPConn, timeout time.Duration) ([]byte, *SessionUDP, error) ***REMOVED***
	conn.SetReadDeadline(time.Now().Add(timeout))
	m := make([]byte, srv.UDPSize)
	n, s, e := ReadFromSessionUDP(conn, m)
	if e != nil || n == 0 ***REMOVED***
		if e != nil ***REMOVED***
			return nil, nil, e
		***REMOVED***
		return nil, nil, ErrShortRead
	***REMOVED***
	m = m[:n]
	return m, s, nil
***REMOVED***

// WriteMsg implements the ResponseWriter.WriteMsg method.
func (w *response) WriteMsg(m *Msg) (err error) ***REMOVED***
	var data []byte
	if w.tsigSecret != nil ***REMOVED*** // if no secrets, dont check for the tsig (which is a longer check)
		if t := m.IsTsig(); t != nil ***REMOVED***
			data, w.tsigRequestMAC, err = TsigGenerate(m, w.tsigSecret[t.Hdr.Name], w.tsigRequestMAC, w.tsigTimersOnly)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			_, err = w.writer.Write(data)
			return err
		***REMOVED***
	***REMOVED***
	data, err = m.Pack()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	_, err = w.writer.Write(data)
	return err
***REMOVED***

// Write implements the ResponseWriter.Write method.
func (w *response) Write(m []byte) (int, error) ***REMOVED***
	switch ***REMOVED***
	case w.udp != nil:
		n, err := WriteToSessionUDP(w.udp, m, w.udpSession)
		return n, err
	case w.tcp != nil:
		lm := len(m)
		if lm < 2 ***REMOVED***
			return 0, io.ErrShortBuffer
		***REMOVED***
		if lm > MaxMsgSize ***REMOVED***
			return 0, &Error***REMOVED***err: "message too large"***REMOVED***
		***REMOVED***
		l := make([]byte, 2, 2+lm)
		l[0], l[1] = packUint16(uint16(lm))
		m = append(l, m...)

		n, err := io.Copy(w.tcp, bytes.NewReader(m))
		return int(n), err
	***REMOVED***
	panic("not reached")
***REMOVED***

// LocalAddr implements the ResponseWriter.LocalAddr method.
func (w *response) LocalAddr() net.Addr ***REMOVED***
	if w.tcp != nil ***REMOVED***
		return w.tcp.LocalAddr()
	***REMOVED***
	return w.udp.LocalAddr()
***REMOVED***

// RemoteAddr implements the ResponseWriter.RemoteAddr method.
func (w *response) RemoteAddr() net.Addr ***REMOVED*** return w.remoteAddr ***REMOVED***

// TsigStatus implements the ResponseWriter.TsigStatus method.
func (w *response) TsigStatus() error ***REMOVED*** return w.tsigStatus ***REMOVED***

// TsigTimersOnly implements the ResponseWriter.TsigTimersOnly method.
func (w *response) TsigTimersOnly(b bool) ***REMOVED*** w.tsigTimersOnly = b ***REMOVED***

// Hijack implements the ResponseWriter.Hijack method.
func (w *response) Hijack() ***REMOVED*** w.hijacked = true ***REMOVED***

// Close implements the ResponseWriter.Close method
func (w *response) Close() error ***REMOVED***
	// Can't close the udp conn, as that is actually the listener.
	if w.tcp != nil ***REMOVED***
		e := w.tcp.Close()
		w.tcp = nil
		return e
	***REMOVED***
	return nil
***REMOVED***
