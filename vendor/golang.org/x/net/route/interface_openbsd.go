// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package route

func (*wireFormat) parseInterfaceMessage(_ RIBType, b []byte) (Message, error) ***REMOVED***
	if len(b) < 32 ***REMOVED***
		return nil, errMessageTooShort
	***REMOVED***
	l := int(nativeEndian.Uint16(b[:2]))
	if len(b) < l ***REMOVED***
		return nil, errInvalidMessage
	***REMOVED***
	attrs := uint(nativeEndian.Uint32(b[12:16]))
	if attrs&sysRTA_IFP == 0 ***REMOVED***
		return nil, nil
	***REMOVED***
	m := &InterfaceMessage***REMOVED***
		Version: int(b[2]),
		Type:    int(b[3]),
		Flags:   int(nativeEndian.Uint32(b[16:20])),
		Index:   int(nativeEndian.Uint16(b[6:8])),
		Addrs:   make([]Addr, sysRTAX_MAX),
		raw:     b[:l],
	***REMOVED***
	ll := int(nativeEndian.Uint16(b[4:6]))
	if len(b) < ll ***REMOVED***
		return nil, errInvalidMessage
	***REMOVED***
	a, err := parseLinkAddr(b[ll:])
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	m.Addrs[sysRTAX_IFP] = a
	m.Name = a.(*LinkAddr).Name
	return m, nil
***REMOVED***

func (*wireFormat) parseInterfaceAddrMessage(_ RIBType, b []byte) (Message, error) ***REMOVED***
	if len(b) < 24 ***REMOVED***
		return nil, errMessageTooShort
	***REMOVED***
	l := int(nativeEndian.Uint16(b[:2]))
	if len(b) < l ***REMOVED***
		return nil, errInvalidMessage
	***REMOVED***
	bodyOff := int(nativeEndian.Uint16(b[4:6]))
	if len(b) < bodyOff ***REMOVED***
		return nil, errInvalidMessage
	***REMOVED***
	m := &InterfaceAddrMessage***REMOVED***
		Version: int(b[2]),
		Type:    int(b[3]),
		Flags:   int(nativeEndian.Uint32(b[12:16])),
		Index:   int(nativeEndian.Uint16(b[6:8])),
		raw:     b[:l],
	***REMOVED***
	var err error
	m.Addrs, err = parseAddrs(uint(nativeEndian.Uint32(b[12:16])), parseKernelInetAddr, b[bodyOff:])
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return m, nil
***REMOVED***

func (*wireFormat) parseInterfaceAnnounceMessage(_ RIBType, b []byte) (Message, error) ***REMOVED***
	if len(b) < 26 ***REMOVED***
		return nil, errMessageTooShort
	***REMOVED***
	l := int(nativeEndian.Uint16(b[:2]))
	if len(b) < l ***REMOVED***
		return nil, errInvalidMessage
	***REMOVED***
	m := &InterfaceAnnounceMessage***REMOVED***
		Version: int(b[2]),
		Type:    int(b[3]),
		Index:   int(nativeEndian.Uint16(b[6:8])),
		What:    int(nativeEndian.Uint16(b[8:10])),
		raw:     b[:l],
	***REMOVED***
	for i := 0; i < 16; i++ ***REMOVED***
		if b[10+i] != 0 ***REMOVED***
			continue
		***REMOVED***
		m.Name = string(b[10 : 10+i])
		break
	***REMOVED***
	return m, nil
***REMOVED***
