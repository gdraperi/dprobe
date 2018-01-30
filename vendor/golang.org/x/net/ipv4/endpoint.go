// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipv4

import (
	"net"
	"syscall"
	"time"

	"golang.org/x/net/internal/socket"
)

// BUG(mikio): On Windows, the JoinSourceSpecificGroup,
// LeaveSourceSpecificGroup, ExcludeSourceSpecificGroup and
// IncludeSourceSpecificGroup methods of PacketConn and RawConn are
// not implemented.

// A Conn represents a network endpoint that uses the IPv4 transport.
// It is used to control basic IP-level socket options such as TOS and
// TTL.
type Conn struct ***REMOVED***
	genericOpt
***REMOVED***

type genericOpt struct ***REMOVED***
	*socket.Conn
***REMOVED***

func (c *genericOpt) ok() bool ***REMOVED*** return c != nil && c.Conn != nil ***REMOVED***

// NewConn returns a new Conn.
func NewConn(c net.Conn) *Conn ***REMOVED***
	cc, _ := socket.NewConn(c)
	return &Conn***REMOVED***
		genericOpt: genericOpt***REMOVED***Conn: cc***REMOVED***,
	***REMOVED***
***REMOVED***

// A PacketConn represents a packet network endpoint that uses the
// IPv4 transport. It is used to control several IP-level socket
// options including multicasting. It also provides datagram based
// network I/O methods specific to the IPv4 and higher layer protocols
// such as UDP.
type PacketConn struct ***REMOVED***
	genericOpt
	dgramOpt
	payloadHandler
***REMOVED***

type dgramOpt struct ***REMOVED***
	*socket.Conn
***REMOVED***

func (c *dgramOpt) ok() bool ***REMOVED*** return c != nil && c.Conn != nil ***REMOVED***

// SetControlMessage sets the per packet IP-level socket options.
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
	return c.payloadHandler.PacketConn.SetDeadline(t)
***REMOVED***

// SetReadDeadline sets the read deadline associated with the
// endpoint.
func (c *PacketConn) SetReadDeadline(t time.Time) error ***REMOVED***
	if !c.payloadHandler.ok() ***REMOVED***
		return syscall.EINVAL
	***REMOVED***
	return c.payloadHandler.PacketConn.SetReadDeadline(t)
***REMOVED***

// SetWriteDeadline sets the write deadline associated with the
// endpoint.
func (c *PacketConn) SetWriteDeadline(t time.Time) error ***REMOVED***
	if !c.payloadHandler.ok() ***REMOVED***
		return syscall.EINVAL
	***REMOVED***
	return c.payloadHandler.PacketConn.SetWriteDeadline(t)
***REMOVED***

// Close closes the endpoint.
func (c *PacketConn) Close() error ***REMOVED***
	if !c.payloadHandler.ok() ***REMOVED***
		return syscall.EINVAL
	***REMOVED***
	return c.payloadHandler.PacketConn.Close()
***REMOVED***

// NewPacketConn returns a new PacketConn using c as its underlying
// transport.
func NewPacketConn(c net.PacketConn) *PacketConn ***REMOVED***
	cc, _ := socket.NewConn(c.(net.Conn))
	p := &PacketConn***REMOVED***
		genericOpt:     genericOpt***REMOVED***Conn: cc***REMOVED***,
		dgramOpt:       dgramOpt***REMOVED***Conn: cc***REMOVED***,
		payloadHandler: payloadHandler***REMOVED***PacketConn: c, Conn: cc***REMOVED***,
	***REMOVED***
	return p
***REMOVED***

// A RawConn represents a packet network endpoint that uses the IPv4
// transport. It is used to control several IP-level socket options
// including IPv4 header manipulation. It also provides datagram
// based network I/O methods specific to the IPv4 and higher layer
// protocols that handle IPv4 datagram directly such as OSPF, GRE.
type RawConn struct ***REMOVED***
	genericOpt
	dgramOpt
	packetHandler
***REMOVED***

// SetControlMessage sets the per packet IP-level socket options.
func (c *RawConn) SetControlMessage(cf ControlFlags, on bool) error ***REMOVED***
	if !c.packetHandler.ok() ***REMOVED***
		return syscall.EINVAL
	***REMOVED***
	return setControlMessage(c.dgramOpt.Conn, &c.packetHandler.rawOpt, cf, on)
***REMOVED***

// SetDeadline sets the read and write deadlines associated with the
// endpoint.
func (c *RawConn) SetDeadline(t time.Time) error ***REMOVED***
	if !c.packetHandler.ok() ***REMOVED***
		return syscall.EINVAL
	***REMOVED***
	return c.packetHandler.IPConn.SetDeadline(t)
***REMOVED***

// SetReadDeadline sets the read deadline associated with the
// endpoint.
func (c *RawConn) SetReadDeadline(t time.Time) error ***REMOVED***
	if !c.packetHandler.ok() ***REMOVED***
		return syscall.EINVAL
	***REMOVED***
	return c.packetHandler.IPConn.SetReadDeadline(t)
***REMOVED***

// SetWriteDeadline sets the write deadline associated with the
// endpoint.
func (c *RawConn) SetWriteDeadline(t time.Time) error ***REMOVED***
	if !c.packetHandler.ok() ***REMOVED***
		return syscall.EINVAL
	***REMOVED***
	return c.packetHandler.IPConn.SetWriteDeadline(t)
***REMOVED***

// Close closes the endpoint.
func (c *RawConn) Close() error ***REMOVED***
	if !c.packetHandler.ok() ***REMOVED***
		return syscall.EINVAL
	***REMOVED***
	return c.packetHandler.IPConn.Close()
***REMOVED***

// NewRawConn returns a new RawConn using c as its underlying
// transport.
func NewRawConn(c net.PacketConn) (*RawConn, error) ***REMOVED***
	cc, err := socket.NewConn(c.(net.Conn))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	r := &RawConn***REMOVED***
		genericOpt:    genericOpt***REMOVED***Conn: cc***REMOVED***,
		dgramOpt:      dgramOpt***REMOVED***Conn: cc***REMOVED***,
		packetHandler: packetHandler***REMOVED***IPConn: c.(*net.IPConn), Conn: cc***REMOVED***,
	***REMOVED***
	so, ok := sockOpts[ssoHeaderPrepend]
	if !ok ***REMOVED***
		return nil, errOpNoSupport
	***REMOVED***
	if err := so.SetInt(r.dgramOpt.Conn, boolint(true)); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return r, nil
***REMOVED***
