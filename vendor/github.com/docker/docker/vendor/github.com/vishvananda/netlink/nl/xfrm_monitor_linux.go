package nl

import (
	"unsafe"
)

const (
	SizeofXfrmUserExpire = 0xe8
)

// struct xfrm_user_expire ***REMOVED***
// 	struct xfrm_usersa_info		state;
// 	__u8				hard;
// ***REMOVED***;

type XfrmUserExpire struct ***REMOVED***
	XfrmUsersaInfo XfrmUsersaInfo
	Hard           uint8
	Pad            [7]byte
***REMOVED***

func (msg *XfrmUserExpire) Len() int ***REMOVED***
	return SizeofXfrmUserExpire
***REMOVED***

func DeserializeXfrmUserExpire(b []byte) *XfrmUserExpire ***REMOVED***
	return (*XfrmUserExpire)(unsafe.Pointer(&b[0:SizeofXfrmUserExpire][0]))
***REMOVED***

func (msg *XfrmUserExpire) Serialize() []byte ***REMOVED***
	return (*(*[SizeofXfrmUserExpire]byte)(unsafe.Pointer(msg)))[:]
***REMOVED***
