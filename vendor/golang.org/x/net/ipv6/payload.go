// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipv6

import (
	"net"

	"golang.org/x/net/internal/socket"
)

// BUG(mikio): On Windows, the ControlMessage for ReadFrom and WriteTo
// methods of PacketConn is not implemented.

// A payloadHandler represents the IPv6 datagram payload handler.
type payloadHandler struct ***REMOVED***
	net.PacketConn
	*socket.Conn
	rawOpt
***REMOVED***

func (c *payloadHandler) ok() bool ***REMOVED*** return c != nil && c.PacketConn != nil && c.Conn != nil ***REMOVED***
