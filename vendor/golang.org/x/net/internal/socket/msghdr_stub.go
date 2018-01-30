// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !darwin,!dragonfly,!freebsd,!linux,!netbsd,!openbsd,!solaris

package socket

type msghdr struct***REMOVED******REMOVED***

func (h *msghdr) pack(vs []iovec, bs [][]byte, oob []byte, sa []byte) ***REMOVED******REMOVED***
func (h *msghdr) name() []byte                                        ***REMOVED*** return nil ***REMOVED***
func (h *msghdr) controllen() int                                     ***REMOVED*** return 0 ***REMOVED***
func (h *msghdr) flags() int                                          ***REMOVED*** return 0 ***REMOVED***
