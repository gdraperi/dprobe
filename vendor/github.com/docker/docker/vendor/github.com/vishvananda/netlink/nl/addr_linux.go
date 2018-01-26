package nl

import (
	"syscall"
	"unsafe"
)

type IfAddrmsg struct ***REMOVED***
	syscall.IfAddrmsg
***REMOVED***

func NewIfAddrmsg(family int) *IfAddrmsg ***REMOVED***
	return &IfAddrmsg***REMOVED***
		IfAddrmsg: syscall.IfAddrmsg***REMOVED***
			Family: uint8(family),
		***REMOVED***,
	***REMOVED***
***REMOVED***

// struct ifaddrmsg ***REMOVED***
//   __u8    ifa_family;
//   __u8    ifa_prefixlen;  /* The prefix length    */
//   __u8    ifa_flags;  /* Flags      */
//   __u8    ifa_scope;  /* Address scope    */
//   __u32   ifa_index;  /* Link index     */
// ***REMOVED***;

// type IfAddrmsg struct ***REMOVED***
// 	Family    uint8
// 	Prefixlen uint8
// 	Flags     uint8
// 	Scope     uint8
// 	Index     uint32
// ***REMOVED***
// SizeofIfAddrmsg     = 0x8

func DeserializeIfAddrmsg(b []byte) *IfAddrmsg ***REMOVED***
	return (*IfAddrmsg)(unsafe.Pointer(&b[0:syscall.SizeofIfAddrmsg][0]))
***REMOVED***

func (msg *IfAddrmsg) Serialize() []byte ***REMOVED***
	return (*(*[syscall.SizeofIfAddrmsg]byte)(unsafe.Pointer(msg)))[:]
***REMOVED***

func (msg *IfAddrmsg) Len() int ***REMOVED***
	return syscall.SizeofIfAddrmsg
***REMOVED***

// struct ifa_cacheinfo ***REMOVED***
// 	__u32	ifa_prefered;
// 	__u32	ifa_valid;
// 	__u32	cstamp; /* created timestamp, hundredths of seconds */
// 	__u32	tstamp; /* updated timestamp, hundredths of seconds */
// ***REMOVED***;

const IFA_CACHEINFO = 6
const SizeofIfaCacheInfo = 0x10

type IfaCacheInfo struct ***REMOVED***
	IfaPrefered uint32
	IfaValid    uint32
	Cstamp      uint32
	Tstamp      uint32
***REMOVED***

func (msg *IfaCacheInfo) Len() int ***REMOVED***
	return SizeofIfaCacheInfo
***REMOVED***

func DeserializeIfaCacheInfo(b []byte) *IfaCacheInfo ***REMOVED***
	return (*IfaCacheInfo)(unsafe.Pointer(&b[0:SizeofIfaCacheInfo][0]))
***REMOVED***

func (msg *IfaCacheInfo) Serialize() []byte ***REMOVED***
	return (*(*[SizeofIfaCacheInfo]byte)(unsafe.Pointer(msg)))[:]
***REMOVED***
