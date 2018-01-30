// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package route

func (w *wireFormat) parseInterfaceMessage(typ RIBType, b []byte) (Message, error) ***REMOVED***
	var extOff, bodyOff int
	if typ == sysNET_RT_IFLISTL ***REMOVED***
		if len(b) < 20 ***REMOVED***
			return nil, errMessageTooShort
		***REMOVED***
		extOff = int(nativeEndian.Uint16(b[18:20]))
		bodyOff = int(nativeEndian.Uint16(b[16:18]))
	***REMOVED*** else ***REMOVED***
		extOff = w.extOff
		bodyOff = w.bodyOff
	***REMOVED***
	if len(b) < extOff || len(b) < bodyOff ***REMOVED***
		return nil, errInvalidMessage
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
		Flags:   int(nativeEndian.Uint32(b[8:12])),
		Index:   int(nativeEndian.Uint16(b[12:14])),
		Addrs:   make([]Addr, sysRTAX_MAX),
		extOff:  extOff,
		raw:     b[:l],
	***REMOVED***
	a, err := parseLinkAddr(b[bodyOff:])
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	m.Addrs[sysRTAX_IFP] = a
	m.Name = a.(*LinkAddr).Name
	return m, nil
***REMOVED***

func (w *wireFormat) parseInterfaceAddrMessage(typ RIBType, b []byte) (Message, error) ***REMOVED***
	var bodyOff int
	if typ == sysNET_RT_IFLISTL ***REMOVED***
		if len(b) < 24 ***REMOVED***
			return nil, errMessageTooShort
		***REMOVED***
		bodyOff = int(nativeEndian.Uint16(b[16:18]))
	***REMOVED*** else ***REMOVED***
		bodyOff = w.bodyOff
	***REMOVED***
	if len(b) < bodyOff ***REMOVED***
		return nil, errInvalidMessage
	***REMOVED***
	l := int(nativeEndian.Uint16(b[:2]))
	if len(b) < l ***REMOVED***
		return nil, errInvalidMessage
	***REMOVED***
	m := &InterfaceAddrMessage***REMOVED***
		Version: int(b[2]),
		Type:    int(b[3]),
		Flags:   int(nativeEndian.Uint32(b[8:12])),
		Index:   int(nativeEndian.Uint16(b[12:14])),
		raw:     b[:l],
	***REMOVED***
	var err error
	m.Addrs, err = parseAddrs(uint(nativeEndian.Uint32(b[4:8])), parseKernelInetAddr, b[bodyOff:])
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return m, nil
***REMOVED***
