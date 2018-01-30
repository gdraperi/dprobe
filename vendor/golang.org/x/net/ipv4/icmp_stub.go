// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !linux

package ipv4

const sizeofICMPFilter = 0x0

type icmpFilter struct ***REMOVED***
***REMOVED***

func (f *icmpFilter) accept(typ ICMPType) ***REMOVED***
***REMOVED***

func (f *icmpFilter) block(typ ICMPType) ***REMOVED***
***REMOVED***

func (f *icmpFilter) setAll(block bool) ***REMOVED***
***REMOVED***

func (f *icmpFilter) willBlock(typ ICMPType) bool ***REMOVED***
	return false
***REMOVED***
