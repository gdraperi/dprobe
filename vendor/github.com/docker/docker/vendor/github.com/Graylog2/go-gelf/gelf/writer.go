// Copyright 2012 SocialCode. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package gelf

import (
	"net"
)

type Writer interface ***REMOVED***
	Close() error
	Write([]byte) (int, error)
	WriteMessage(*Message) error
***REMOVED***

// Writer implements io.Writer and is used to send both discrete
// messages to a graylog2 server, or data from a stream-oriented
// interface (like the functions in log).
type GelfWriter struct ***REMOVED***
	addr     string
	conn     net.Conn
	hostname string
	Facility string // defaults to current process name
	proto    string
***REMOVED***

// Close connection and interrupt blocked Read or Write operations
func (w *GelfWriter) Close() error ***REMOVED***
	if w.conn == nil ***REMOVED***
		return nil
	***REMOVED***
	return w.conn.Close()
***REMOVED***
