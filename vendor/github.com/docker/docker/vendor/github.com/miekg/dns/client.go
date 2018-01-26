package dns

// A client implementation.

import (
	"bytes"
	"io"
	"net"
	"time"
)

const dnsTimeout time.Duration = 2 * time.Second
const tcpIdleTimeout time.Duration = 8 * time.Second

// A Conn represents a connection to a DNS server.
type Conn struct ***REMOVED***
	net.Conn                         // a net.Conn holding the connection
	UDPSize        uint16            // minimum receive buffer for UDP messages
	TsigSecret     map[string]string // secret(s) for Tsig map[<zonename>]<base64 secret>, zonename must be fully qualified
	rtt            time.Duration
	t              time.Time
	tsigRequestMAC string
***REMOVED***

// A Client defines parameters for a DNS client.
type Client struct ***REMOVED***
	Net            string            // if "tcp" a TCP query will be initiated, otherwise an UDP one (default is "" for UDP)
	UDPSize        uint16            // minimum receive buffer for UDP messages
	DialTimeout    time.Duration     // net.DialTimeout, defaults to 2 seconds
	ReadTimeout    time.Duration     // net.Conn.SetReadTimeout value for connections, defaults to 2 seconds
	WriteTimeout   time.Duration     // net.Conn.SetWriteTimeout value for connections, defaults to 2 seconds
	TsigSecret     map[string]string // secret(s) for Tsig map[<zonename>]<base64 secret>, zonename must be fully qualified
	SingleInflight bool              // if true suppress multiple outstanding queries for the same Qname, Qtype and Qclass
	group          singleflight
***REMOVED***

// Exchange performs a synchronous UDP query. It sends the message m to the address
// contained in a and waits for an reply. Exchange does not retry a failed query, nor
// will it fall back to TCP in case of truncation.
// If you need to send a DNS message on an already existing connection, you can use the
// following:
//
//	co := &dns.Conn***REMOVED***Conn: c***REMOVED*** // c is your net.Conn
//	co.WriteMsg(m)
//	in, err := co.ReadMsg()
//	co.Close()
//
func Exchange(m *Msg, a string) (r *Msg, err error) ***REMOVED***
	var co *Conn
	co, err = DialTimeout("udp", a, dnsTimeout)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	defer co.Close()

	opt := m.IsEdns0()
	// If EDNS0 is used use that for size.
	if opt != nil && opt.UDPSize() >= MinMsgSize ***REMOVED***
		co.UDPSize = opt.UDPSize()
	***REMOVED***

	co.SetWriteDeadline(time.Now().Add(dnsTimeout))
	if err = co.WriteMsg(m); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	co.SetReadDeadline(time.Now().Add(dnsTimeout))
	r, err = co.ReadMsg()
	if err == nil && r.Id != m.Id ***REMOVED***
		err = ErrId
	***REMOVED***
	return r, err
***REMOVED***

// ExchangeConn performs a synchronous query. It sends the message m via the connection
// c and waits for a reply. The connection c is not closed by ExchangeConn.
// This function is going away, but can easily be mimicked:
//
//	co := &dns.Conn***REMOVED***Conn: c***REMOVED*** // c is your net.Conn
//	co.WriteMsg(m)
//	in, _  := co.ReadMsg()
//	co.Close()
//
func ExchangeConn(c net.Conn, m *Msg) (r *Msg, err error) ***REMOVED***
	println("dns: this function is deprecated")
	co := new(Conn)
	co.Conn = c
	if err = co.WriteMsg(m); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	r, err = co.ReadMsg()
	if err == nil && r.Id != m.Id ***REMOVED***
		err = ErrId
	***REMOVED***
	return r, err
***REMOVED***

// Exchange performs an synchronous query. It sends the message m to the address
// contained in a and waits for an reply. Basic use pattern with a *dns.Client:
//
//	c := new(dns.Client)
//	in, rtt, err := c.Exchange(message, "127.0.0.1:53")
//
// Exchange does not retry a failed query, nor will it fall back to TCP in
// case of truncation.
func (c *Client) Exchange(m *Msg, a string) (r *Msg, rtt time.Duration, err error) ***REMOVED***
	if !c.SingleInflight ***REMOVED***
		return c.exchange(m, a)
	***REMOVED***
	// This adds a bunch of garbage, TODO(miek).
	t := "nop"
	if t1, ok := TypeToString[m.Question[0].Qtype]; ok ***REMOVED***
		t = t1
	***REMOVED***
	cl := "nop"
	if cl1, ok := ClassToString[m.Question[0].Qclass]; ok ***REMOVED***
		cl = cl1
	***REMOVED***
	r, rtt, err, shared := c.group.Do(m.Question[0].Name+t+cl, func() (*Msg, time.Duration, error) ***REMOVED***
		return c.exchange(m, a)
	***REMOVED***)
	if err != nil ***REMOVED***
		return r, rtt, err
	***REMOVED***
	if shared ***REMOVED***
		return r.Copy(), rtt, nil
	***REMOVED***
	return r, rtt, nil
***REMOVED***

func (c *Client) dialTimeout() time.Duration ***REMOVED***
	if c.DialTimeout != 0 ***REMOVED***
		return c.DialTimeout
	***REMOVED***
	return dnsTimeout
***REMOVED***

func (c *Client) readTimeout() time.Duration ***REMOVED***
	if c.ReadTimeout != 0 ***REMOVED***
		return c.ReadTimeout
	***REMOVED***
	return dnsTimeout
***REMOVED***

func (c *Client) writeTimeout() time.Duration ***REMOVED***
	if c.WriteTimeout != 0 ***REMOVED***
		return c.WriteTimeout
	***REMOVED***
	return dnsTimeout
***REMOVED***

func (c *Client) exchange(m *Msg, a string) (r *Msg, rtt time.Duration, err error) ***REMOVED***
	var co *Conn
	if c.Net == "" ***REMOVED***
		co, err = DialTimeout("udp", a, c.dialTimeout())
	***REMOVED*** else ***REMOVED***
		co, err = DialTimeout(c.Net, a, c.dialTimeout())
	***REMOVED***
	if err != nil ***REMOVED***
		return nil, 0, err
	***REMOVED***
	defer co.Close()

	opt := m.IsEdns0()
	// If EDNS0 is used use that for size.
	if opt != nil && opt.UDPSize() >= MinMsgSize ***REMOVED***
		co.UDPSize = opt.UDPSize()
	***REMOVED***
	// Otherwise use the client's configured UDP size.
	if opt == nil && c.UDPSize >= MinMsgSize ***REMOVED***
		co.UDPSize = c.UDPSize
	***REMOVED***

	co.TsigSecret = c.TsigSecret
	co.SetWriteDeadline(time.Now().Add(c.writeTimeout()))
	if err = co.WriteMsg(m); err != nil ***REMOVED***
		return nil, 0, err
	***REMOVED***

	co.SetReadDeadline(time.Now().Add(c.readTimeout()))
	r, err = co.ReadMsg()
	if err == nil && r.Id != m.Id ***REMOVED***
		err = ErrId
	***REMOVED***
	return r, co.rtt, err
***REMOVED***

// ReadMsg reads a message from the connection co.
// If the received message contains a TSIG record the transaction
// signature is verified.
func (co *Conn) ReadMsg() (*Msg, error) ***REMOVED***
	p, err := co.ReadMsgHeader(nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	m := new(Msg)
	if err := m.Unpack(p); err != nil ***REMOVED***
		// If ErrTruncated was returned, we still want to allow the user to use
		// the message, but naively they can just check err if they don't want
		// to use a truncated message
		if err == ErrTruncated ***REMOVED***
			return m, err
		***REMOVED***
		return nil, err
	***REMOVED***
	if t := m.IsTsig(); t != nil ***REMOVED***
		if _, ok := co.TsigSecret[t.Hdr.Name]; !ok ***REMOVED***
			return m, ErrSecret
		***REMOVED***
		// Need to work on the original message p, as that was used to calculate the tsig.
		err = TsigVerify(p, co.TsigSecret[t.Hdr.Name], co.tsigRequestMAC, false)
	***REMOVED***
	return m, err
***REMOVED***

// ReadMsgHeader reads a DNS message, parses and populates hdr (when hdr is not nil).
// Returns message as a byte slice to be parsed with Msg.Unpack later on.
// Note that error handling on the message body is not possible as only the header is parsed.
func (co *Conn) ReadMsgHeader(hdr *Header) ([]byte, error) ***REMOVED***
	var (
		p   []byte
		n   int
		err error
	)

	if t, ok := co.Conn.(*net.TCPConn); ok ***REMOVED***
		// First two bytes specify the length of the entire message.
		l, err := tcpMsgLen(t)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		p = make([]byte, l)
		n, err = tcpRead(t, p)
	***REMOVED*** else ***REMOVED***
		if co.UDPSize > MinMsgSize ***REMOVED***
			p = make([]byte, co.UDPSize)
		***REMOVED*** else ***REMOVED***
			p = make([]byte, MinMsgSize)
		***REMOVED***
		n, err = co.Read(p)
	***REMOVED***

	if err != nil ***REMOVED***
		return nil, err
	***REMOVED*** else if n < headerSize ***REMOVED***
		return nil, ErrShortRead
	***REMOVED***

	p = p[:n]
	if hdr != nil ***REMOVED***
		if _, err = UnpackStruct(hdr, p, 0); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	return p, err
***REMOVED***

// tcpMsgLen is a helper func to read first two bytes of stream as uint16 packet length.
func tcpMsgLen(t *net.TCPConn) (int, error) ***REMOVED***
	p := []byte***REMOVED***0, 0***REMOVED***
	n, err := t.Read(p)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	if n != 2 ***REMOVED***
		return 0, ErrShortRead
	***REMOVED***
	l, _ := unpackUint16(p, 0)
	if l == 0 ***REMOVED***
		return 0, ErrShortRead
	***REMOVED***
	return int(l), nil
***REMOVED***

// tcpRead calls TCPConn.Read enough times to fill allocated buffer.
func tcpRead(t *net.TCPConn, p []byte) (int, error) ***REMOVED***
	n, err := t.Read(p)
	if err != nil ***REMOVED***
		return n, err
	***REMOVED***
	for n < len(p) ***REMOVED***
		j, err := t.Read(p[n:])
		if err != nil ***REMOVED***
			return n, err
		***REMOVED***
		n += j
	***REMOVED***
	return n, err
***REMOVED***

// Read implements the net.Conn read method.
func (co *Conn) Read(p []byte) (n int, err error) ***REMOVED***
	if co.Conn == nil ***REMOVED***
		return 0, ErrConnEmpty
	***REMOVED***
	if len(p) < 2 ***REMOVED***
		return 0, io.ErrShortBuffer
	***REMOVED***
	if t, ok := co.Conn.(*net.TCPConn); ok ***REMOVED***
		l, err := tcpMsgLen(t)
		if err != nil ***REMOVED***
			return 0, err
		***REMOVED***
		if l > len(p) ***REMOVED***
			return int(l), io.ErrShortBuffer
		***REMOVED***
		return tcpRead(t, p[:l])
	***REMOVED***
	// UDP connection
	n, err = co.Conn.Read(p)
	if err != nil ***REMOVED***
		return n, err
	***REMOVED***

	co.rtt = time.Since(co.t)
	return n, err
***REMOVED***

// WriteMsg sends a message throught the connection co.
// If the message m contains a TSIG record the transaction
// signature is calculated.
func (co *Conn) WriteMsg(m *Msg) (err error) ***REMOVED***
	var out []byte
	if t := m.IsTsig(); t != nil ***REMOVED***
		mac := ""
		if _, ok := co.TsigSecret[t.Hdr.Name]; !ok ***REMOVED***
			return ErrSecret
		***REMOVED***
		out, mac, err = TsigGenerate(m, co.TsigSecret[t.Hdr.Name], co.tsigRequestMAC, false)
		// Set for the next read, allthough only used in zone transfers
		co.tsigRequestMAC = mac
	***REMOVED*** else ***REMOVED***
		out, err = m.Pack()
	***REMOVED***
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	co.t = time.Now()
	if _, err = co.Write(out); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

// Write implements the net.Conn Write method.
func (co *Conn) Write(p []byte) (n int, err error) ***REMOVED***
	if t, ok := co.Conn.(*net.TCPConn); ok ***REMOVED***
		lp := len(p)
		if lp < 2 ***REMOVED***
			return 0, io.ErrShortBuffer
		***REMOVED***
		if lp > MaxMsgSize ***REMOVED***
			return 0, &Error***REMOVED***err: "message too large"***REMOVED***
		***REMOVED***
		l := make([]byte, 2, lp+2)
		l[0], l[1] = packUint16(uint16(lp))
		p = append(l, p...)
		n, err := io.Copy(t, bytes.NewReader(p))
		return int(n), err
	***REMOVED***
	n, err = co.Conn.(*net.UDPConn).Write(p)
	return n, err
***REMOVED***

// Dial connects to the address on the named network.
func Dial(network, address string) (conn *Conn, err error) ***REMOVED***
	conn = new(Conn)
	conn.Conn, err = net.Dial(network, address)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return conn, nil
***REMOVED***

// DialTimeout acts like Dial but takes a timeout.
func DialTimeout(network, address string, timeout time.Duration) (conn *Conn, err error) ***REMOVED***
	conn = new(Conn)
	conn.Conn, err = net.DialTimeout(network, address, timeout)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return conn, nil
***REMOVED***
