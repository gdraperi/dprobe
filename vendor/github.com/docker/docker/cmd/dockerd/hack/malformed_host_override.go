// +build !windows

package hack

import "net"

// MalformedHostHeaderOverride is a wrapper to be able
// to overcome the 400 Bad request coming from old docker
// clients that send an invalid Host header.
type MalformedHostHeaderOverride struct ***REMOVED***
	net.Listener
***REMOVED***

// MalformedHostHeaderOverrideConn wraps the underlying unix
// connection and keeps track of the first read from http.Server
// which just reads the headers.
type MalformedHostHeaderOverrideConn struct ***REMOVED***
	net.Conn
	first bool
***REMOVED***

var closeConnHeader = []byte("\r\nConnection: close\r")

// Read reads the first *read* request from http.Server to inspect
// the Host header. If the Host starts with / then we're talking to
// an old docker client which send an invalid Host header. To not
// error out in http.Server we rewrite the first bytes of the request
// to sanitize the Host header itself.
// In case we're not dealing with old docker clients the data is just passed
// to the server w/o modification.
func (l *MalformedHostHeaderOverrideConn) Read(b []byte) (n int, err error) ***REMOVED***
	// http.Server uses a 4k buffer
	if l.first && len(b) == 4096 ***REMOVED***
		// This keeps track of the first read from http.Server which just reads
		// the headers
		l.first = false
		// The first read of the connection by http.Server is done limited to
		// DefaultMaxHeaderBytes (usually 1 << 20) + 4096.
		// Here we do the first read which gets us all the http headers to
		// be inspected and modified below.
		c, err := l.Conn.Read(b)
		if err != nil ***REMOVED***
			return c, err
		***REMOVED***

		var (
			start, end    int
			firstLineFeed = -1
			buf           []byte
		)
		for i := 0; i <= c-1-7; i++ ***REMOVED***
			if b[i] == '\n' && firstLineFeed == -1 ***REMOVED***
				firstLineFeed = i
			***REMOVED***
			if b[i] != '\n' ***REMOVED***
				continue
			***REMOVED***

			if b[i+1] == '\r' && b[i+2] == '\n' ***REMOVED***
				return c, nil
			***REMOVED***

			if b[i+1] != 'H' ***REMOVED***
				continue
			***REMOVED***
			if b[i+2] != 'o' ***REMOVED***
				continue
			***REMOVED***
			if b[i+3] != 's' ***REMOVED***
				continue
			***REMOVED***
			if b[i+4] != 't' ***REMOVED***
				continue
			***REMOVED***
			if b[i+5] != ':' ***REMOVED***
				continue
			***REMOVED***
			if b[i+6] != ' ' ***REMOVED***
				continue
			***REMOVED***
			if b[i+7] != '/' ***REMOVED***
				continue
			***REMOVED***
			// ensure clients other than the docker clients do not get this hack
			if i != firstLineFeed ***REMOVED***
				return c, nil
			***REMOVED***
			start = i + 7
			// now find where the value ends
			for ii, bbb := range b[start:c] ***REMOVED***
				if bbb == '\n' ***REMOVED***
					end = start + ii
					break
				***REMOVED***
			***REMOVED***
			buf = make([]byte, 0, c+len(closeConnHeader)-(end-start))
			// strip the value of the host header and
			// inject `Connection: close` to ensure we don't reuse this connection
			buf = append(buf, b[:start]...)
			buf = append(buf, closeConnHeader...)
			buf = append(buf, b[end:c]...)
			copy(b, buf)
			break
		***REMOVED***
		if len(buf) == 0 ***REMOVED***
			return c, nil
		***REMOVED***
		return len(buf), nil
	***REMOVED***
	return l.Conn.Read(b)
***REMOVED***

// Accept makes the listener accepts connections and wraps the connection
// in a MalformedHostHeaderOverrideConn initializing first to true.
func (l *MalformedHostHeaderOverride) Accept() (net.Conn, error) ***REMOVED***
	c, err := l.Listener.Accept()
	if err != nil ***REMOVED***
		return c, err
	***REMOVED***
	return &MalformedHostHeaderOverrideConn***REMOVED***c, true***REMOVED***, nil
***REMOVED***
