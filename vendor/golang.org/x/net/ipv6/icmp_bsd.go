// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd netbsd openbsd

package ipv6

func (f *icmpv6Filter) accept(typ ICMPType) ***REMOVED***
	f.Filt[typ>>5] |= 1 << (uint32(typ) & 31)
***REMOVED***

func (f *icmpv6Filter) block(typ ICMPType) ***REMOVED***
	f.Filt[typ>>5] &^= 1 << (uint32(typ) & 31)
***REMOVED***

func (f *icmpv6Filter) setAll(block bool) ***REMOVED***
	for i := range f.Filt ***REMOVED***
		if block ***REMOVED***
			f.Filt[i] = 0
		***REMOVED*** else ***REMOVED***
			f.Filt[i] = 1<<32 - 1
		***REMOVED***
	***REMOVED***
***REMOVED***

func (f *icmpv6Filter) willBlock(typ ICMPType) bool ***REMOVED***
	return f.Filt[typ>>5]&(1<<(uint32(typ)&31)) == 0
***REMOVED***
