package gelf

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

const (
	DefaultMaxReconnect   = 3
	DefaultReconnectDelay = 1
)

type TCPWriter struct ***REMOVED***
	GelfWriter
	mu             sync.Mutex
	MaxReconnect   int
	ReconnectDelay time.Duration
***REMOVED***

func NewTCPWriter(addr string) (*TCPWriter, error) ***REMOVED***
	var err error
	w := new(TCPWriter)
	w.MaxReconnect = DefaultMaxReconnect
	w.ReconnectDelay = DefaultReconnectDelay
	w.proto = "tcp"
	w.addr = addr

	if w.conn, err = net.Dial("tcp", addr); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if w.hostname, err = os.Hostname(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return w, nil
***REMOVED***

// WriteMessage sends the specified message to the GELF server
// specified in the call to New().  It assumes all the fields are
// filled out appropriately.  In general, clients will want to use
// Write, rather than WriteMessage.
func (w *TCPWriter) WriteMessage(m *Message) (err error) ***REMOVED***
	messageBytes, err := m.toBytes()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	messageBytes = append(messageBytes, 0)

	n, err := w.writeToSocketWithReconnectAttempts(messageBytes)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if n != len(messageBytes) ***REMOVED***
		return fmt.Errorf("bad write (%d/%d)", n, len(messageBytes))
	***REMOVED***

	return nil
***REMOVED***

func (w *TCPWriter) Write(p []byte) (n int, err error) ***REMOVED***
	file, line := getCallerIgnoringLogMulti(1)

	m := constructMessage(p, w.hostname, w.Facility, file, line)

	if err = w.WriteMessage(m); err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	return len(p), nil
***REMOVED***

func (w *TCPWriter) writeToSocketWithReconnectAttempts(zBytes []byte) (n int, err error) ***REMOVED***
	var errConn error
	var i int

	w.mu.Lock()
	for i = 0; i <= w.MaxReconnect; i++ ***REMOVED***
		errConn = nil

		if w.conn != nil ***REMOVED***
			n, err = w.conn.Write(zBytes)
		***REMOVED*** else ***REMOVED***
			err = fmt.Errorf("Connection was nil, will attempt reconnect")
		***REMOVED***
		if err != nil ***REMOVED***
			time.Sleep(w.ReconnectDelay * time.Second)
			w.conn, errConn = net.Dial("tcp", w.addr)
		***REMOVED*** else ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	w.mu.Unlock()

	if i > w.MaxReconnect ***REMOVED***
		return 0, fmt.Errorf("Maximum reconnection attempts was reached; giving up")
	***REMOVED***
	if errConn != nil ***REMOVED***
		return 0, fmt.Errorf("Write Failed: %s\nReconnection failed: %s", err, errConn)
	***REMOVED***
	return n, nil
***REMOVED***
