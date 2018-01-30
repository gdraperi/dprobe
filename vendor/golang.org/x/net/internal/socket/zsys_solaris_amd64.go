// Created by cgo -godefs - DO NOT EDIT
// cgo -godefs defs_solaris.go

package socket

const (
	sysAF_UNSPEC = 0x0
	sysAF_INET   = 0x2
	sysAF_INET6  = 0x1a

	sysSOCK_RAW = 0x4
)

type iovec struct ***REMOVED***
	Base *int8
	Len  uint64
***REMOVED***

type msghdr struct ***REMOVED***
	Name         *byte
	Namelen      uint32
	Pad_cgo_0    [4]byte
	Iov          *iovec
	Iovlen       int32
	Pad_cgo_1    [4]byte
	Accrights    *int8
	Accrightslen int32
	Pad_cgo_2    [4]byte
***REMOVED***

type cmsghdr struct ***REMOVED***
	Len   uint32
	Level int32
	Type  int32
***REMOVED***

type sockaddrInet struct ***REMOVED***
	Family uint16
	Port   uint16
	Addr   [4]byte /* in_addr */
	Zero   [8]int8
***REMOVED***

type sockaddrInet6 struct ***REMOVED***
	Family         uint16
	Port           uint16
	Flowinfo       uint32
	Addr           [16]byte /* in6_addr */
	Scope_id       uint32
	X__sin6_src_id uint32
***REMOVED***

const (
	sizeofIovec   = 0x10
	sizeofMsghdr  = 0x30
	sizeofCmsghdr = 0xc

	sizeofSockaddrInet  = 0x10
	sizeofSockaddrInet6 = 0x20
)
