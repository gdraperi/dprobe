package nl

import (
	"syscall"
	"unsafe"
)

type RtMsg struct ***REMOVED***
	syscall.RtMsg
***REMOVED***

func NewRtMsg() *RtMsg ***REMOVED***
	return &RtMsg***REMOVED***
		RtMsg: syscall.RtMsg***REMOVED***
			Table:    syscall.RT_TABLE_MAIN,
			Scope:    syscall.RT_SCOPE_UNIVERSE,
			Protocol: syscall.RTPROT_BOOT,
			Type:     syscall.RTN_UNICAST,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func NewRtDelMsg() *RtMsg ***REMOVED***
	return &RtMsg***REMOVED***
		RtMsg: syscall.RtMsg***REMOVED***
			Table: syscall.RT_TABLE_MAIN,
			Scope: syscall.RT_SCOPE_NOWHERE,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func (msg *RtMsg) Len() int ***REMOVED***
	return syscall.SizeofRtMsg
***REMOVED***

func DeserializeRtMsg(b []byte) *RtMsg ***REMOVED***
	return (*RtMsg)(unsafe.Pointer(&b[0:syscall.SizeofRtMsg][0]))
***REMOVED***

func (msg *RtMsg) Serialize() []byte ***REMOVED***
	return (*(*[syscall.SizeofRtMsg]byte)(unsafe.Pointer(msg)))[:]
***REMOVED***

type RtNexthop struct ***REMOVED***
	syscall.RtNexthop
	Children []NetlinkRequestData
***REMOVED***

func DeserializeRtNexthop(b []byte) *RtNexthop ***REMOVED***
	return (*RtNexthop)(unsafe.Pointer(&b[0:syscall.SizeofRtNexthop][0]))
***REMOVED***

func (msg *RtNexthop) Len() int ***REMOVED***
	if len(msg.Children) == 0 ***REMOVED***
		return syscall.SizeofRtNexthop
	***REMOVED***

	l := 0
	for _, child := range msg.Children ***REMOVED***
		l += rtaAlignOf(child.Len())
	***REMOVED***
	l += syscall.SizeofRtNexthop
	return rtaAlignOf(l)
***REMOVED***

func (msg *RtNexthop) Serialize() []byte ***REMOVED***
	length := msg.Len()
	msg.RtNexthop.Len = uint16(length)
	buf := make([]byte, length)
	copy(buf, (*(*[syscall.SizeofRtNexthop]byte)(unsafe.Pointer(msg)))[:])
	next := rtaAlignOf(syscall.SizeofRtNexthop)
	if len(msg.Children) > 0 ***REMOVED***
		for _, child := range msg.Children ***REMOVED***
			childBuf := child.Serialize()
			copy(buf[next:], childBuf)
			next += rtaAlignOf(len(childBuf))
		***REMOVED***
	***REMOVED***
	return buf
***REMOVED***
