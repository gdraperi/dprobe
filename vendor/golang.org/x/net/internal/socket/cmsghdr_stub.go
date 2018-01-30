// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !darwin,!dragonfly,!freebsd,!linux,!netbsd,!openbsd,!solaris

package socket

type cmsghdr struct***REMOVED******REMOVED***

const sizeofCmsghdr = 0

func (h *cmsghdr) len() int ***REMOVED*** return 0 ***REMOVED***
func (h *cmsghdr) lvl() int ***REMOVED*** return 0 ***REMOVED***
func (h *cmsghdr) typ() int ***REMOVED*** return 0 ***REMOVED***

func (h *cmsghdr) set(l, lvl, typ int) ***REMOVED******REMOVED***
