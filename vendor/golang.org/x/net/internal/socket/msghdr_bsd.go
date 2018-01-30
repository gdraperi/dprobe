// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd netbsd openbsd

package socket

import "unsafe"

func (h *msghdr) pack(vs []iovec, bs [][]byte, oob []byte, sa []byte) ***REMOVED***
	for i := range vs ***REMOVED***
		vs[i].set(bs[i])
	***REMOVED***
	h.setIov(vs)
	if len(oob) > 0 ***REMOVED***
		h.Control = (*byte)(unsafe.Pointer(&oob[0]))
		h.Controllen = uint32(len(oob))
	***REMOVED***
	if sa != nil ***REMOVED***
		h.Name = (*byte)(unsafe.Pointer(&sa[0]))
		h.Namelen = uint32(len(sa))
	***REMOVED***
***REMOVED***

func (h *msghdr) name() []byte ***REMOVED***
	if h.Name != nil && h.Namelen > 0 ***REMOVED***
		return (*[sizeofSockaddrInet6]byte)(unsafe.Pointer(h.Name))[:h.Namelen]
	***REMOVED***
	return nil
***REMOVED***

func (h *msghdr) controllen() int ***REMOVED***
	return int(h.Controllen)
***REMOVED***

func (h *msghdr) flags() int ***REMOVED***
	return int(h.Flags)
***REMOVED***
