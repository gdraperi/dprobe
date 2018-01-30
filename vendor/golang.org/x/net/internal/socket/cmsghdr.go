// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package socket

func (h *cmsghdr) len() int ***REMOVED*** return int(h.Len) ***REMOVED***
func (h *cmsghdr) lvl() int ***REMOVED*** return int(h.Level) ***REMOVED***
func (h *cmsghdr) typ() int ***REMOVED*** return int(h.Type) ***REMOVED***
