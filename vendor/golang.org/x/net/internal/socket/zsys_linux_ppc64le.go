// Created by cgo -godefs - DO NOT EDIT
// cgo -godefs defs_linux.go

package socket

const (
	sysAF_UNSPEC = 0x0
	sysAF_INET   = 0x2
	sysAF_INET6  = 0xa

	sysSOCK_RAW = 0x3
)

type iovec struct ***REMOVED***
	Base *byte
	Len  uint64
***REMOVED***

type msghdr struct ***REMOVED***
	Name       *byte
	Namelen    uint32
	Pad_cgo_0  [4]byte
	Iov        *iovec
	Iovlen     uint64
	Control    *byte
	Controllen uint64
	Flags      int32
	Pad_cgo_1  [4]byte
***REMOVED***

type mmsghdr struct ***REMOVED***
	Hdr       msghdr
	Len       uint32
	Pad_cgo_0 [4]byte
***REMOVED***

type cmsghdr struct ***REMOVED***
	Len   uint64
	Level int32
	Type  int32
***REMOVED***

type sockaddrInet struct ***REMOVED***
	Family uint16
	Port   uint16
	Addr   [4]byte /* in_addr */
	X__pad [8]uint8
***REMOVED***

type sockaddrInet6 struct ***REMOVED***
	Family   uint16
	Port     uint16
	Flowinfo uint32
	Addr     [16]byte /* in6_addr */
	Scope_id uint32
***REMOVED***

const (
	sizeofIovec   = 0x10
	sizeofMsghdr  = 0x38
	sizeofMmsghdr = 0x40
	sizeofCmsghdr = 0x10

	sizeofSockaddrInet  = 0x10
	sizeofSockaddrInet6 = 0x1c
)
