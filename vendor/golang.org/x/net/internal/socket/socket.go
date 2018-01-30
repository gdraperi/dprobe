// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package socket provides a portable interface for socket system
// calls.
package socket // import "golang.org/x/net/internal/socket"

import (
	"errors"
	"net"
	"unsafe"
)

// An Option represents a sticky socket option.
type Option struct ***REMOVED***
	Level int // level
	Name  int // name; must be equal or greater than 1
	Len   int // length of value in bytes; must be equal or greater than 1
***REMOVED***

// Get reads a value for the option from the kernel.
// It returns the number of bytes written into b.
func (o *Option) Get(c *Conn, b []byte) (int, error) ***REMOVED***
	if o.Name < 1 || o.Len < 1 ***REMOVED***
		return 0, errors.New("invalid option")
	***REMOVED***
	if len(b) < o.Len ***REMOVED***
		return 0, errors.New("short buffer")
	***REMOVED***
	return o.get(c, b)
***REMOVED***

// GetInt returns an integer value for the option.
//
// The Len field of Option must be either 1 or 4.
func (o *Option) GetInt(c *Conn) (int, error) ***REMOVED***
	if o.Len != 1 && o.Len != 4 ***REMOVED***
		return 0, errors.New("invalid option")
	***REMOVED***
	var b []byte
	var bb [4]byte
	if o.Len == 1 ***REMOVED***
		b = bb[:1]
	***REMOVED*** else ***REMOVED***
		b = bb[:4]
	***REMOVED***
	n, err := o.get(c, b)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	if n != o.Len ***REMOVED***
		return 0, errors.New("invalid option length")
	***REMOVED***
	if o.Len == 1 ***REMOVED***
		return int(b[0]), nil
	***REMOVED***
	return int(NativeEndian.Uint32(b[:4])), nil
***REMOVED***

// Set writes the option and value to the kernel.
func (o *Option) Set(c *Conn, b []byte) error ***REMOVED***
	if o.Name < 1 || o.Len < 1 ***REMOVED***
		return errors.New("invalid option")
	***REMOVED***
	if len(b) < o.Len ***REMOVED***
		return errors.New("short buffer")
	***REMOVED***
	return o.set(c, b)
***REMOVED***

// SetInt writes the option and value to the kernel.
//
// The Len field of Option must be either 1 or 4.
func (o *Option) SetInt(c *Conn, v int) error ***REMOVED***
	if o.Len != 1 && o.Len != 4 ***REMOVED***
		return errors.New("invalid option")
	***REMOVED***
	var b []byte
	if o.Len == 1 ***REMOVED***
		b = []byte***REMOVED***byte(v)***REMOVED***
	***REMOVED*** else ***REMOVED***
		var bb [4]byte
		NativeEndian.PutUint32(bb[:o.Len], uint32(v))
		b = bb[:4]
	***REMOVED***
	return o.set(c, b)
***REMOVED***

func controlHeaderLen() int ***REMOVED***
	return roundup(sizeofCmsghdr)
***REMOVED***

func controlMessageLen(dataLen int) int ***REMOVED***
	return roundup(sizeofCmsghdr) + dataLen
***REMOVED***

// ControlMessageSpace returns the whole length of control message.
func ControlMessageSpace(dataLen int) int ***REMOVED***
	return roundup(sizeofCmsghdr) + roundup(dataLen)
***REMOVED***

// A ControlMessage represents the head message in a stream of control
// messages.
//
// A control message comprises of a header, data and a few padding
// fields to conform to the interface to the kernel.
//
// See RFC 3542 for further information.
type ControlMessage []byte

// Data returns the data field of the control message at the head on
// m.
func (m ControlMessage) Data(dataLen int) []byte ***REMOVED***
	l := controlHeaderLen()
	if len(m) < l || len(m) < l+dataLen ***REMOVED***
		return nil
	***REMOVED***
	return m[l : l+dataLen]
***REMOVED***

// Next returns the control message at the next on m.
//
// Next works only for standard control messages.
func (m ControlMessage) Next(dataLen int) ControlMessage ***REMOVED***
	l := ControlMessageSpace(dataLen)
	if len(m) < l ***REMOVED***
		return nil
	***REMOVED***
	return m[l:]
***REMOVED***

// MarshalHeader marshals the header fields of the control message at
// the head on m.
func (m ControlMessage) MarshalHeader(lvl, typ, dataLen int) error ***REMOVED***
	if len(m) < controlHeaderLen() ***REMOVED***
		return errors.New("short message")
	***REMOVED***
	h := (*cmsghdr)(unsafe.Pointer(&m[0]))
	h.set(controlMessageLen(dataLen), lvl, typ)
	return nil
***REMOVED***

// ParseHeader parses and returns the header fields of the control
// message at the head on m.
func (m ControlMessage) ParseHeader() (lvl, typ, dataLen int, err error) ***REMOVED***
	l := controlHeaderLen()
	if len(m) < l ***REMOVED***
		return 0, 0, 0, errors.New("short message")
	***REMOVED***
	h := (*cmsghdr)(unsafe.Pointer(&m[0]))
	return h.lvl(), h.typ(), int(uint64(h.len()) - uint64(l)), nil
***REMOVED***

// Marshal marshals the control message at the head on m, and returns
// the next control message.
func (m ControlMessage) Marshal(lvl, typ int, data []byte) (ControlMessage, error) ***REMOVED***
	l := len(data)
	if len(m) < ControlMessageSpace(l) ***REMOVED***
		return nil, errors.New("short message")
	***REMOVED***
	h := (*cmsghdr)(unsafe.Pointer(&m[0]))
	h.set(controlMessageLen(l), lvl, typ)
	if l > 0 ***REMOVED***
		copy(m.Data(l), data)
	***REMOVED***
	return m.Next(l), nil
***REMOVED***

// Parse parses m as a single or multiple control messages.
//
// Parse works for both standard and compatible messages.
func (m ControlMessage) Parse() ([]ControlMessage, error) ***REMOVED***
	var ms []ControlMessage
	for len(m) >= controlHeaderLen() ***REMOVED***
		h := (*cmsghdr)(unsafe.Pointer(&m[0]))
		l := h.len()
		if l <= 0 ***REMOVED***
			return nil, errors.New("invalid header length")
		***REMOVED***
		if uint64(l) < uint64(controlHeaderLen()) ***REMOVED***
			return nil, errors.New("invalid message length")
		***REMOVED***
		if uint64(l) > uint64(len(m)) ***REMOVED***
			return nil, errors.New("short buffer")
		***REMOVED***
		// On message reception:
		//
		// |<- ControlMessageSpace --------------->|
		// |<- controlMessageLen ---------->|      |
		// |<- controlHeaderLen ->|         |      |
		// +---------------+------+---------+------+
		// |    Header     | PadH |  Data   | PadD |
		// +---------------+------+---------+------+
		//
		// On compatible message reception:
		//
		// | ... |<- controlMessageLen ----------->|
		// | ... |<- controlHeaderLen ->|          |
		// +-----+---------------+------+----------+
		// | ... |    Header     | PadH |   Data   |
		// +-----+---------------+------+----------+
		ms = append(ms, ControlMessage(m[:l]))
		ll := l - controlHeaderLen()
		if len(m) >= ControlMessageSpace(ll) ***REMOVED***
			m = m[ControlMessageSpace(ll):]
		***REMOVED*** else ***REMOVED***
			m = m[controlMessageLen(ll):]
		***REMOVED***
	***REMOVED***
	return ms, nil
***REMOVED***

// NewControlMessage returns a new stream of control messages.
func NewControlMessage(dataLen []int) ControlMessage ***REMOVED***
	var l int
	for i := range dataLen ***REMOVED***
		l += ControlMessageSpace(dataLen[i])
	***REMOVED***
	return make([]byte, l)
***REMOVED***

// A Message represents an IO message.
type Message struct ***REMOVED***
	// When writing, the Buffers field must contain at least one
	// byte to write.
	// When reading, the Buffers field will always contain a byte
	// to read.
	Buffers [][]byte

	// OOB contains protocol-specific control or miscellaneous
	// ancillary data known as out-of-band data.
	OOB []byte

	// Addr specifies a destination address when writing.
	// It can be nil when the underlying protocol of the raw
	// connection uses connection-oriented communication.
	// After a successful read, it may contain the source address
	// on the received packet.
	Addr net.Addr

	N     int // # of bytes read or written from/to Buffers
	NN    int // # of bytes read or written from/to OOB
	Flags int // protocol-specific information on the received message
***REMOVED***

// RecvMsg wraps recvmsg system call.
//
// The provided flags is a set of platform-dependent flags, such as
// syscall.MSG_PEEK.
func (c *Conn) RecvMsg(m *Message, flags int) error ***REMOVED***
	return c.recvMsg(m, flags)
***REMOVED***

// SendMsg wraps sendmsg system call.
//
// The provided flags is a set of platform-dependent flags, such as
// syscall.MSG_DONTROUTE.
func (c *Conn) SendMsg(m *Message, flags int) error ***REMOVED***
	return c.sendMsg(m, flags)
***REMOVED***

// RecvMsgs wraps recvmmsg system call.
//
// It returns the number of processed messages.
//
// The provided flags is a set of platform-dependent flags, such as
// syscall.MSG_PEEK.
//
// Only Linux supports this.
func (c *Conn) RecvMsgs(ms []Message, flags int) (int, error) ***REMOVED***
	return c.recvMsgs(ms, flags)
***REMOVED***

// SendMsgs wraps sendmmsg system call.
//
// It returns the number of processed messages.
//
// The provided flags is a set of platform-dependent flags, such as
// syscall.MSG_DONTROUTE.
//
// Only Linux supports this.
func (c *Conn) SendMsgs(ms []Message, flags int) (int, error) ***REMOVED***
	return c.sendMsgs(ms, flags)
***REMOVED***
