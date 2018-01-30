// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package ipv4

import (
	"unsafe"

	"golang.org/x/net/internal/iana"
	"golang.org/x/net/internal/socket"
)

func setControlMessage(c *socket.Conn, opt *rawOpt, cf ControlFlags, on bool) error ***REMOVED***
	opt.Lock()
	defer opt.Unlock()
	if so, ok := sockOpts[ssoReceiveTTL]; ok && cf&FlagTTL != 0 ***REMOVED***
		if err := so.SetInt(c, boolint(on)); err != nil ***REMOVED***
			return err
		***REMOVED***
		if on ***REMOVED***
			opt.set(FlagTTL)
		***REMOVED*** else ***REMOVED***
			opt.clear(FlagTTL)
		***REMOVED***
	***REMOVED***
	if so, ok := sockOpts[ssoPacketInfo]; ok ***REMOVED***
		if cf&(FlagSrc|FlagDst|FlagInterface) != 0 ***REMOVED***
			if err := so.SetInt(c, boolint(on)); err != nil ***REMOVED***
				return err
			***REMOVED***
			if on ***REMOVED***
				opt.set(cf & (FlagSrc | FlagDst | FlagInterface))
			***REMOVED*** else ***REMOVED***
				opt.clear(cf & (FlagSrc | FlagDst | FlagInterface))
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if so, ok := sockOpts[ssoReceiveDst]; ok && cf&FlagDst != 0 ***REMOVED***
			if err := so.SetInt(c, boolint(on)); err != nil ***REMOVED***
				return err
			***REMOVED***
			if on ***REMOVED***
				opt.set(FlagDst)
			***REMOVED*** else ***REMOVED***
				opt.clear(FlagDst)
			***REMOVED***
		***REMOVED***
		if so, ok := sockOpts[ssoReceiveInterface]; ok && cf&FlagInterface != 0 ***REMOVED***
			if err := so.SetInt(c, boolint(on)); err != nil ***REMOVED***
				return err
			***REMOVED***
			if on ***REMOVED***
				opt.set(FlagInterface)
			***REMOVED*** else ***REMOVED***
				opt.clear(FlagInterface)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func marshalTTL(b []byte, cm *ControlMessage) []byte ***REMOVED***
	m := socket.ControlMessage(b)
	m.MarshalHeader(iana.ProtocolIP, sysIP_RECVTTL, 1)
	return m.Next(1)
***REMOVED***

func parseTTL(cm *ControlMessage, b []byte) ***REMOVED***
	cm.TTL = int(*(*byte)(unsafe.Pointer(&b[:1][0])))
***REMOVED***
