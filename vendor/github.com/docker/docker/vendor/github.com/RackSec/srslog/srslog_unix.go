package srslog

import (
	"errors"
	"io"
	"net"
)

// unixSyslog opens a connection to the syslog daemon running on the
// local machine using a Unix domain socket. This function exists because of
// Solaris support as implemented by gccgo.  On Solaris you can not
// simply open a TCP connection to the syslog daemon.  The gccgo
// sources have a syslog_solaris.go file that implements unixSyslog to
// return a type that satisfies the serverConn interface and simply calls the C
// library syslog function.
func unixSyslog() (conn serverConn, err error) ***REMOVED***
	logTypes := []string***REMOVED***"unixgram", "unix"***REMOVED***
	logPaths := []string***REMOVED***"/dev/log", "/var/run/syslog", "/var/run/log"***REMOVED***
	for _, network := range logTypes ***REMOVED***
		for _, path := range logPaths ***REMOVED***
			conn, err := net.Dial(network, path)
			if err != nil ***REMOVED***
				continue
			***REMOVED*** else ***REMOVED***
				return &localConn***REMOVED***conn: conn***REMOVED***, nil
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil, errors.New("Unix syslog delivery error")
***REMOVED***

// localConn adheres to the serverConn interface, allowing us to send syslog
// messages to the local syslog daemon over a Unix domain socket.
type localConn struct ***REMOVED***
	conn io.WriteCloser
***REMOVED***

// writeString formats syslog messages using time.Stamp instead of time.RFC3339,
// and omits the hostname (because it is expected to be used locally).
func (n *localConn) writeString(framer Framer, formatter Formatter, p Priority, hostname, tag, msg string) error ***REMOVED***
	if framer == nil ***REMOVED***
		framer = DefaultFramer
	***REMOVED***
	if formatter == nil ***REMOVED***
		formatter = UnixFormatter
	***REMOVED***
	_, err := n.conn.Write([]byte(framer(formatter(p, hostname, tag, msg))))
	return err
***REMOVED***

// close the (local) network connection
func (n *localConn) close() error ***REMOVED***
	return n.conn.Close()
***REMOVED***
