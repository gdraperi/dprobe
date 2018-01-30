// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build amd64
// +build solaris

package socket

import "unsafe"

func (h *msghdr) pack(vs []iovec, bs [][]byte, oob []byte, sa []byte) ***REMOVED***
	for i := range vs ***REMOVED***
		vs[i].set(bs[i])
	***REMOVED***
	if len(vs) > 0 ***REMOVED***
		h.Iov = &vs[0]
		h.Iovlen = int32(len(vs))
	***REMOVED***
	if len(oob) > 0 ***REMOVED***
		h.Accrights = (*int8)(unsafe.Pointer(&oob[0]))
		h.Accrightslen = int32(len(oob))
	***REMOVED***
	if sa != nil ***REMOVED***
		h.Name = (*byte)(unsafe.Pointer(&sa[0]))
		h.Namelen = uint32(len(sa))
	***REMOVED***
***REMOVED***

func (h *msghdr) controllen() int ***REMOVED***
	return int(h.Accrightslen)
***REMOVED***

func (h *msghdr) flags() int ***REMOVED***
	return int(NativeEndian.Uint32(h.Pad_cgo_2[:]))
***REMOVED***
