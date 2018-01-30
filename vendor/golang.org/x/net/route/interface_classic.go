// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly netbsd

package route

import "runtime"

func (w *wireFormat) parseInterfaceMessage(_ RIBType, b []byte) (Message, error) ***REMOVED***
	if len(b) < w.bodyOff ***REMOVED***
		return nil, errMessageTooShort
	***REMOVED***
	l := int(nativeEndian.Uint16(b[:2]))
	if len(b) < l ***REMOVED***
		return nil, errInvalidMessage
	***REMOVED***
	attrs := uint(nativeEndian.Uint32(b[4:8]))
	if attrs&sysRTA_IFP == 0 ***REMOVED***
		return nil, nil
	***REMOVED***
	m := &InterfaceMessage***REMOVED***
		Version: int(b[2]),
		Type:    int(b[3]),
		Addrs:   make([]Addr, sysRTAX_MAX),
		Flags:   int(nativeEndian.Uint32(b[8:12])),
		Index:   int(nativeEndian.Uint16(b[12:14])),
		extOff:  w.extOff,
		raw:     b[:l],
	***REMOVED***
	a, err := parseLinkAddr(b[w.bodyOff:])
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	m.Addrs[sysRTAX_IFP] = a
	m.Name = a.(*LinkAddr).Name
	return m, nil
***REMOVED***

func (w *wireFormat) parseInterfaceAddrMessage(_ RIBType, b []byte) (Message, error) ***REMOVED***
	if len(b) < w.bodyOff ***REMOVED***
		return nil, errMessageTooShort
	***REMOVED***
	l := int(nativeEndian.Uint16(b[:2]))
	if len(b) < l ***REMOVED***
		return nil, errInvalidMessage
	***REMOVED***
	m := &InterfaceAddrMessage***REMOVED***
		Version: int(b[2]),
		Type:    int(b[3]),
		Flags:   int(nativeEndian.Uint32(b[8:12])),
		raw:     b[:l],
	***REMOVED***
	if runtime.GOOS == "netbsd" ***REMOVED***
		m.Index = int(nativeEndian.Uint16(b[16:18]))
	***REMOVED*** else ***REMOVED***
		m.Index = int(nativeEndian.Uint16(b[12:14]))
	***REMOVED***
	var err error
	m.Addrs, err = parseAddrs(uint(nativeEndian.Uint32(b[4:8])), parseKernelInetAddr, b[w.bodyOff:])
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return m, nil
***REMOVED***
