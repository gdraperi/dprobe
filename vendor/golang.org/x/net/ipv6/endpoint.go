// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipv6

import (
	"net"
	"syscall"
	"time"

	"golang.org/x/net/internal/socket"
)

// BUG(mikio): On Windows, the JoinSourceSpecificGroup,
// LeaveSourceSpecificGroup, ExcludeSourceSpecificGroup and
// IncludeSourceSpecificGroup methods of PacketConn are not
// implemented.

// A Conn represents a network endpoint that uses IPv6 transport.
// It allows to set basic IP-level socket options such as traffic
// class and hop limit.
type Conn struct ***REMOVED***
	genericOpt
***REMOVED***

type genericOpt struct ***REMOVED***
	*socket.Conn
***REMOVED***

func (c *genericOpt) ok() bool ***REMOVED*** return c != nil && c.Conn != nil ***REMOVED***

// PathMTU returns a path MTU value for the destination associated
// with the endpoint.
func (c *Conn) PathMTU() (int, error) ***REMOVED***
	if !c.ok() ***REMOVED***
		return 0, syscall.EINVAL
	***REMOVED***
	so, ok := sockOpts[ssoPathMTU]
	if !ok ***REMOVED***
		return 0, errOpNoSupport
	***REMOVED***
	_, mtu, err := so.getMTUInfo(c.Conn)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	return mtu, nil
***REMOVED***

// NewConn returns a new Conn.
func NewConn(c net.Conn) *Conn ***REMOVED***
	cc, _ := socket.NewConn(c)
	return &Conn***REMOVED***
		genericOpt: genericOpt***REMOVED***Conn: cc***REMOVED***,
	***REMOVED***
***REMOVED***

// A PacketConn represents a packet network endpoint that uses IPv6
// transport. It is used to control several IP-level socket options
// including IPv6 header manipulation. It also provides datagram
// based network I/O methods specific to the IPv6 and higher layer
// protocols such as OSPF, GRE, and UDP.
type PacketConn struct ***REMOVED***
	genericOpt
	dgramOpt
	payloadHandler
***REMOVED***

type dgramOpt struct ***REMOVED***
	*socket.Conn
***REMOVED***

func (c *dgramOpt) ok() bool ***REMOVED*** return c != nil && c.Conn != nil ***REMOVED***

// SetControlMessage allows to receive the per packet basis IP-level
// socket options.
func (c *PacketConn) SetControlMessage(cf ControlFlags, on bool) error ***REMOVED***
	if !c.payloadHandler.ok() ***REMOVED***
		return syscall.EINVAL
	***REMOVED***
	return setControlMessage(c.dgramOpt.Conn, &c.payloadHandler.rawOpt, cf, on)
***REMOVED***

// SetDeadline sets the read and write deadlines associated with the
// endpoint.
func (c *PacketConn) SetDeadline(t time.Time) error ***REMOVED***
	if !c.payloadHandler.ok() ***REMOVED***
		return syscall.EINVAL
	***REMOVED***
	return c.payloadHandler.SetDeadline(t)
***REMOVED***

// SetReadDeadline sets the read deadline associated with the
// endpoint.
func (c *PacketConn) SetReadDeadline(t time.Time) error ***REMOVED***
	if !c.payloadHandler.ok() ***REMOVED***
		return syscall.EINVAL
	***REMOVED***
	return c.payloadHandler.SetReadDeadline(t)
***REMOVED***

// SetWriteDeadline sets the write deadline associated with the
// endpoint.
func (c *PacketConn) SetWriteDeadline(t time.Time) error ***REMOVED***
	if !c.payloadHandler.ok() ***REMOVED***
		return syscall.EINVAL
	***REMOVED***
	return c.payloadHandler.SetWriteDeadline(t)
***REMOVED***

// Close closes the endpoint.
func (c *PacketConn) Close() error ***REMOVED***
	if !c.payloadHandler.ok() ***REMOVED***
		return syscall.EINVAL
	***REMOVED***
	return c.payloadHandler.Close()
***REMOVED***

// NewPacketConn returns a new PacketConn using c as its underlying
// transport.
func NewPacketConn(c net.PacketConn) *PacketConn ***REMOVED***
	cc, _ := socket.NewConn(c.(net.Conn))
	return &PacketConn***REMOVED***
		genericOpt:     genericOpt***REMOVED***Conn: cc***REMOVED***,
		dgramOpt:       dgramOpt***REMOVED***Conn: cc***REMOVED***,
		payloadHandler: payloadHandler***REMOVED***PacketConn: c, Conn: cc***REMOVED***,
	***REMOVED***
***REMOVED***
