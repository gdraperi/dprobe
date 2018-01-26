package dns

import (
	"time"
)

// Envelope is used when doing a zone transfer with a remote server.
type Envelope struct ***REMOVED***
	RR    []RR  // The set of RRs in the answer section of the xfr reply message.
	Error error // If something went wrong, this contains the error.
***REMOVED***

// A Transfer defines parameters that are used during a zone transfer.
type Transfer struct ***REMOVED***
	*Conn
	DialTimeout    time.Duration     // net.DialTimeout, defaults to 2 seconds
	ReadTimeout    time.Duration     // net.Conn.SetReadTimeout value for connections, defaults to 2 seconds
	WriteTimeout   time.Duration     // net.Conn.SetWriteTimeout value for connections, defaults to 2 seconds
	TsigSecret     map[string]string // Secret(s) for Tsig map[<zonename>]<base64 secret>, zonename must be fully qualified
	tsigTimersOnly bool
***REMOVED***

// Think we need to away to stop the transfer

// In performs an incoming transfer with the server in a.
// If you would like to set the source IP, or some other attribute
// of a Dialer for a Transfer, you can do so by specifying the attributes
// in the Transfer.Conn:
//
//	d := net.Dialer***REMOVED***LocalAddr: transfer_source***REMOVED***
//	con, err := d.Dial("tcp", master)
//	dnscon := &dns.Conn***REMOVED***Conn:con***REMOVED***
//	transfer = &dns.Transfer***REMOVED***Conn: dnscon***REMOVED***
//	channel, err := transfer.In(message, master)
//
func (t *Transfer) In(q *Msg, a string) (env chan *Envelope, err error) ***REMOVED***
	timeout := dnsTimeout
	if t.DialTimeout != 0 ***REMOVED***
		timeout = t.DialTimeout
	***REMOVED***
	if t.Conn == nil ***REMOVED***
		t.Conn, err = DialTimeout("tcp", a, timeout)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	if err := t.WriteMsg(q); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	env = make(chan *Envelope)
	go func() ***REMOVED***
		if q.Question[0].Qtype == TypeAXFR ***REMOVED***
			go t.inAxfr(q.Id, env)
			return
		***REMOVED***
		if q.Question[0].Qtype == TypeIXFR ***REMOVED***
			go t.inIxfr(q.Id, env)
			return
		***REMOVED***
	***REMOVED***()
	return env, nil
***REMOVED***

func (t *Transfer) inAxfr(id uint16, c chan *Envelope) ***REMOVED***
	first := true
	defer t.Close()
	defer close(c)
	timeout := dnsTimeout
	if t.ReadTimeout != 0 ***REMOVED***
		timeout = t.ReadTimeout
	***REMOVED***
	for ***REMOVED***
		t.Conn.SetReadDeadline(time.Now().Add(timeout))
		in, err := t.ReadMsg()
		if err != nil ***REMOVED***
			c <- &Envelope***REMOVED***nil, err***REMOVED***
			return
		***REMOVED***
		if id != in.Id ***REMOVED***
			c <- &Envelope***REMOVED***in.Answer, ErrId***REMOVED***
			return
		***REMOVED***
		if first ***REMOVED***
			if !isSOAFirst(in) ***REMOVED***
				c <- &Envelope***REMOVED***in.Answer, ErrSoa***REMOVED***
				return
			***REMOVED***
			first = !first
			// only one answer that is SOA, receive more
			if len(in.Answer) == 1 ***REMOVED***
				t.tsigTimersOnly = true
				c <- &Envelope***REMOVED***in.Answer, nil***REMOVED***
				continue
			***REMOVED***
		***REMOVED***

		if !first ***REMOVED***
			t.tsigTimersOnly = true // Subsequent envelopes use this.
			if isSOALast(in) ***REMOVED***
				c <- &Envelope***REMOVED***in.Answer, nil***REMOVED***
				return
			***REMOVED***
			c <- &Envelope***REMOVED***in.Answer, nil***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func (t *Transfer) inIxfr(id uint16, c chan *Envelope) ***REMOVED***
	serial := uint32(0) // The first serial seen is the current server serial
	first := true
	defer t.Close()
	defer close(c)
	timeout := dnsTimeout
	if t.ReadTimeout != 0 ***REMOVED***
		timeout = t.ReadTimeout
	***REMOVED***
	for ***REMOVED***
		t.SetReadDeadline(time.Now().Add(timeout))
		in, err := t.ReadMsg()
		if err != nil ***REMOVED***
			c <- &Envelope***REMOVED***nil, err***REMOVED***
			return
		***REMOVED***
		if id != in.Id ***REMOVED***
			c <- &Envelope***REMOVED***in.Answer, ErrId***REMOVED***
			return
		***REMOVED***
		if first ***REMOVED***
			// A single SOA RR signals "no changes"
			if len(in.Answer) == 1 && isSOAFirst(in) ***REMOVED***
				c <- &Envelope***REMOVED***in.Answer, nil***REMOVED***
				return
			***REMOVED***

			// Check if the returned answer is ok
			if !isSOAFirst(in) ***REMOVED***
				c <- &Envelope***REMOVED***in.Answer, ErrSoa***REMOVED***
				return
			***REMOVED***
			// This serial is important
			serial = in.Answer[0].(*SOA).Serial
			first = !first
		***REMOVED***

		// Now we need to check each message for SOA records, to see what we need to do
		if !first ***REMOVED***
			t.tsigTimersOnly = true
			// If the last record in the IXFR contains the servers' SOA,  we should quit
			if v, ok := in.Answer[len(in.Answer)-1].(*SOA); ok ***REMOVED***
				if v.Serial == serial ***REMOVED***
					c <- &Envelope***REMOVED***in.Answer, nil***REMOVED***
					return
				***REMOVED***
			***REMOVED***
			c <- &Envelope***REMOVED***in.Answer, nil***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

// Out performs an outgoing transfer with the client connecting in w.
// Basic use pattern:
//
//	ch := make(chan *dns.Envelope)
//	tr := new(dns.Transfer)
//	tr.Out(w, r, ch)
//	c <- &dns.Envelope***REMOVED***RR: []dns.RR***REMOVED***soa, rr1, rr2, rr3, soa***REMOVED******REMOVED***
//	close(ch)
//	w.Hijack()
//	// w.Close() // Client closes connection
//
// The server is responsible for sending the correct sequence of RRs through the
// channel ch.
func (t *Transfer) Out(w ResponseWriter, q *Msg, ch chan *Envelope) error ***REMOVED***
	for x := range ch ***REMOVED***
		r := new(Msg)
		// Compress?
		r.SetReply(q)
		r.Authoritative = true
		// assume it fits TODO(miek): fix
		r.Answer = append(r.Answer, x.RR...)
		if err := w.WriteMsg(r); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	w.TsigTimersOnly(true)
	return nil
***REMOVED***

// ReadMsg reads a message from the transfer connection t.
func (t *Transfer) ReadMsg() (*Msg, error) ***REMOVED***
	m := new(Msg)
	p := make([]byte, MaxMsgSize)
	n, err := t.Read(p)
	if err != nil && n == 0 ***REMOVED***
		return nil, err
	***REMOVED***
	p = p[:n]
	if err := m.Unpack(p); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if ts := m.IsTsig(); ts != nil && t.TsigSecret != nil ***REMOVED***
		if _, ok := t.TsigSecret[ts.Hdr.Name]; !ok ***REMOVED***
			return m, ErrSecret
		***REMOVED***
		// Need to work on the original message p, as that was used to calculate the tsig.
		err = TsigVerify(p, t.TsigSecret[ts.Hdr.Name], t.tsigRequestMAC, t.tsigTimersOnly)
		t.tsigRequestMAC = ts.MAC
	***REMOVED***
	return m, err
***REMOVED***

// WriteMsg writes a message through the transfer connection t.
func (t *Transfer) WriteMsg(m *Msg) (err error) ***REMOVED***
	var out []byte
	if ts := m.IsTsig(); ts != nil && t.TsigSecret != nil ***REMOVED***
		if _, ok := t.TsigSecret[ts.Hdr.Name]; !ok ***REMOVED***
			return ErrSecret
		***REMOVED***
		out, t.tsigRequestMAC, err = TsigGenerate(m, t.TsigSecret[ts.Hdr.Name], t.tsigRequestMAC, t.tsigTimersOnly)
	***REMOVED*** else ***REMOVED***
		out, err = m.Pack()
	***REMOVED***
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if _, err = t.Write(out); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

func isSOAFirst(in *Msg) bool ***REMOVED***
	if len(in.Answer) > 0 ***REMOVED***
		return in.Answer[0].Header().Rrtype == TypeSOA
	***REMOVED***
	return false
***REMOVED***

func isSOALast(in *Msg) bool ***REMOVED***
	if len(in.Answer) > 0 ***REMOVED***
		return in.Answer[len(in.Answer)-1].Header().Rrtype == TypeSOA
	***REMOVED***
	return false
***REMOVED***
