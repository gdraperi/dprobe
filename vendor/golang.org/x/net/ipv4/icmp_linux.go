// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipv4

func (f *icmpFilter) accept(typ ICMPType) ***REMOVED***
	f.Data &^= 1 << (uint32(typ) & 31)
***REMOVED***

func (f *icmpFilter) block(typ ICMPType) ***REMOVED***
	f.Data |= 1 << (uint32(typ) & 31)
***REMOVED***

func (f *icmpFilter) setAll(block bool) ***REMOVED***
	if block ***REMOVED***
		f.Data = 1<<32 - 1
	***REMOVED*** else ***REMOVED***
		f.Data = 0
	***REMOVED***
***REMOVED***

func (f *icmpFilter) willBlock(typ ICMPType) bool ***REMOVED***
	return f.Data&(1<<(uint32(typ)&31)) != 0
***REMOVED***
