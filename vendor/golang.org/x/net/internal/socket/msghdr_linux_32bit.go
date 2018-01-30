// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build arm mips mipsle 386
// +build linux

package socket

import "unsafe"

func (h *msghdr) setIov(vs []iovec) ***REMOVED***
	l := len(vs)
	if l == 0 ***REMOVED***
		return
	***REMOVED***
	h.Iov = &vs[0]
	h.Iovlen = uint32(l)
***REMOVED***

func (h *msghdr) setControl(b []byte) ***REMOVED***
	h.Control = (*byte)(unsafe.Pointer(&b[0]))
	h.Controllen = uint32(len(b))
***REMOVED***