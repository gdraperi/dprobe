// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package ipv6

import "golang.org/x/net/internal/socket"

func setControlMessage(c *socket.Conn, opt *rawOpt, cf ControlFlags, on bool) error ***REMOVED***
	opt.Lock()
	defer opt.Unlock()
	if so, ok := sockOpts[ssoReceiveTrafficClass]; ok && cf&FlagTrafficClass != 0 ***REMOVED***
		if err := so.SetInt(c, boolint(on)); err != nil ***REMOVED***
			return err
		***REMOVED***
		if on ***REMOVED***
			opt.set(FlagTrafficClass)
		***REMOVED*** else ***REMOVED***
			opt.clear(FlagTrafficClass)
		***REMOVED***
	***REMOVED***
	if so, ok := sockOpts[ssoReceiveHopLimit]; ok && cf&FlagHopLimit != 0 ***REMOVED***
		if err := so.SetInt(c, boolint(on)); err != nil ***REMOVED***
			return err
		***REMOVED***
		if on ***REMOVED***
			opt.set(FlagHopLimit)
		***REMOVED*** else ***REMOVED***
			opt.clear(FlagHopLimit)
		***REMOVED***
	***REMOVED***
	if so, ok := sockOpts[ssoReceivePacketInfo]; ok && cf&flagPacketInfo != 0 ***REMOVED***
		if err := so.SetInt(c, boolint(on)); err != nil ***REMOVED***
			return err
		***REMOVED***
		if on ***REMOVED***
			opt.set(cf & flagPacketInfo)
		***REMOVED*** else ***REMOVED***
			opt.clear(cf & flagPacketInfo)
		***REMOVED***
	***REMOVED***
	if so, ok := sockOpts[ssoReceivePathMTU]; ok && cf&FlagPathMTU != 0 ***REMOVED***
		if err := so.SetInt(c, boolint(on)); err != nil ***REMOVED***
			return err
		***REMOVED***
		if on ***REMOVED***
			opt.set(FlagPathMTU)
		***REMOVED*** else ***REMOVED***
			opt.clear(FlagPathMTU)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***
