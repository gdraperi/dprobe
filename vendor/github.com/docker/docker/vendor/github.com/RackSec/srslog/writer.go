package srslog

import (
	"crypto/tls"
	"strings"
	"sync"
)

// A Writer is a connection to a syslog server.
type Writer struct ***REMOVED***
	priority  Priority
	tag       string
	hostname  string
	network   string
	raddr     string
	tlsConfig *tls.Config
	framer    Framer
	formatter Formatter

	mu   sync.RWMutex // guards conn
	conn serverConn
***REMOVED***

// getConn provides access to the internal conn, protected by a mutex. The
// conn is threadsafe, so it can be used while unlocked, but we want to avoid
// race conditions on grabbing a reference to it.
func (w *Writer) getConn() serverConn ***REMOVED***
	w.mu.RLock()
	conn := w.conn
	w.mu.RUnlock()
	return conn
***REMOVED***

// setConn updates the internal conn, protected by a mutex.
func (w *Writer) setConn(c serverConn) ***REMOVED***
	w.mu.Lock()
	w.conn = c
	w.mu.Unlock()
***REMOVED***

// connect makes a connection to the syslog server.
func (w *Writer) connect() (serverConn, error) ***REMOVED***
	conn := w.getConn()
	if conn != nil ***REMOVED***
		// ignore err from close, it makes sense to continue anyway
		conn.close()
		w.setConn(nil)
	***REMOVED***

	var hostname string
	var err error
	dialer := w.getDialer()
	conn, hostname, err = dialer.Call()
	if err == nil ***REMOVED***
		w.setConn(conn)
		w.hostname = hostname

		return conn, nil
	***REMOVED*** else ***REMOVED***
		return nil, err
	***REMOVED***
***REMOVED***

// SetFormatter changes the formatter function for subsequent messages.
func (w *Writer) SetFormatter(f Formatter) ***REMOVED***
	w.formatter = f
***REMOVED***

// SetFramer changes the framer function for subsequent messages.
func (w *Writer) SetFramer(f Framer) ***REMOVED***
	w.framer = f
***REMOVED***

// Write sends a log message to the syslog daemon using the default priority
// passed into `srslog.New` or the `srslog.Dial*` functions.
func (w *Writer) Write(b []byte) (int, error) ***REMOVED***
	return w.writeAndRetry(w.priority, string(b))
***REMOVED***

// WriteWithPriority sends a log message with a custom priority
func (w *Writer) WriteWithPriority(p Priority, b []byte) (int, error) ***REMOVED***
	return w.writeAndRetry(p, string(b))
***REMOVED***

// Close closes a connection to the syslog daemon.
func (w *Writer) Close() error ***REMOVED***
	conn := w.getConn()
	if conn != nil ***REMOVED***
		err := conn.close()
		w.setConn(nil)
		return err
	***REMOVED***
	return nil
***REMOVED***

// Emerg logs a message with severity LOG_EMERG; this overrides the default
// priority passed to `srslog.New` and the `srslog.Dial*` functions.
func (w *Writer) Emerg(m string) (err error) ***REMOVED***
	_, err = w.writeAndRetry(LOG_EMERG, m)
	return err
***REMOVED***

// Alert logs a message with severity LOG_ALERT; this overrides the default
// priority passed to `srslog.New` and the `srslog.Dial*` functions.
func (w *Writer) Alert(m string) (err error) ***REMOVED***
	_, err = w.writeAndRetry(LOG_ALERT, m)
	return err
***REMOVED***

// Crit logs a message with severity LOG_CRIT; this overrides the default
// priority passed to `srslog.New` and the `srslog.Dial*` functions.
func (w *Writer) Crit(m string) (err error) ***REMOVED***
	_, err = w.writeAndRetry(LOG_CRIT, m)
	return err
***REMOVED***

// Err logs a message with severity LOG_ERR; this overrides the default
// priority passed to `srslog.New` and the `srslog.Dial*` functions.
func (w *Writer) Err(m string) (err error) ***REMOVED***
	_, err = w.writeAndRetry(LOG_ERR, m)
	return err
***REMOVED***

// Warning logs a message with severity LOG_WARNING; this overrides the default
// priority passed to `srslog.New` and the `srslog.Dial*` functions.
func (w *Writer) Warning(m string) (err error) ***REMOVED***
	_, err = w.writeAndRetry(LOG_WARNING, m)
	return err
***REMOVED***

// Notice logs a message with severity LOG_NOTICE; this overrides the default
// priority passed to `srslog.New` and the `srslog.Dial*` functions.
func (w *Writer) Notice(m string) (err error) ***REMOVED***
	_, err = w.writeAndRetry(LOG_NOTICE, m)
	return err
***REMOVED***

// Info logs a message with severity LOG_INFO; this overrides the default
// priority passed to `srslog.New` and the `srslog.Dial*` functions.
func (w *Writer) Info(m string) (err error) ***REMOVED***
	_, err = w.writeAndRetry(LOG_INFO, m)
	return err
***REMOVED***

// Debug logs a message with severity LOG_DEBUG; this overrides the default
// priority passed to `srslog.New` and the `srslog.Dial*` functions.
func (w *Writer) Debug(m string) (err error) ***REMOVED***
	_, err = w.writeAndRetry(LOG_DEBUG, m)
	return err
***REMOVED***

func (w *Writer) writeAndRetry(p Priority, s string) (int, error) ***REMOVED***
	pr := (w.priority & facilityMask) | (p & severityMask)

	conn := w.getConn()
	if conn != nil ***REMOVED***
		if n, err := w.write(conn, pr, s); err == nil ***REMOVED***
			return n, err
		***REMOVED***
	***REMOVED***

	var err error
	if conn, err = w.connect(); err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	return w.write(conn, pr, s)
***REMOVED***

// write generates and writes a syslog formatted string. It formats the
// message based on the current Formatter and Framer.
func (w *Writer) write(conn serverConn, p Priority, msg string) (int, error) ***REMOVED***
	// ensure it ends in a \n
	if !strings.HasSuffix(msg, "\n") ***REMOVED***
		msg += "\n"
	***REMOVED***

	err := conn.writeString(w.framer, w.formatter, p, w.hostname, w.tag, msg)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	// Note: return the length of the input, not the number of
	// bytes printed by Fprintf, because this must behave like
	// an io.Writer.
	return len(msg), nil
***REMOVED***
