package srslog

import (
	"net"
)

// netConn has an internal net.Conn and adheres to the serverConn interface,
// allowing us to send syslog messages over the network.
type netConn struct ***REMOVED***
	conn net.Conn
***REMOVED***

// writeString formats syslog messages using time.RFC3339 and includes the
// hostname, and sends the message to the connection.
func (n *netConn) writeString(framer Framer, formatter Formatter, p Priority, hostname, tag, msg string) error ***REMOVED***
	if framer == nil ***REMOVED***
		framer = DefaultFramer
	***REMOVED***
	if formatter == nil ***REMOVED***
		formatter = DefaultFormatter
	***REMOVED***
	formattedMessage := framer(formatter(p, hostname, tag, msg))
	_, err := n.conn.Write([]byte(formattedMessage))
	return err
***REMOVED***

// close the network connection
func (n *netConn) close() error ***REMOVED***
	return n.conn.Close()
***REMOVED***
