// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package route

import "syscall"

func (m *RouteMessage) marshal() ([]byte, error) ***REMOVED***
	l := sizeofRtMsghdr + addrsSpace(m.Addrs)
	b := make([]byte, l)
	nativeEndian.PutUint16(b[:2], uint16(l))
	if m.Version == 0 ***REMOVED***
		b[2] = sysRTM_VERSION
	***REMOVED*** else ***REMOVED***
		b[2] = byte(m.Version)
	***REMOVED***
	b[3] = byte(m.Type)
	nativeEndian.PutUint16(b[4:6], uint16(sizeofRtMsghdr))
	nativeEndian.PutUint32(b[16:20], uint32(m.Flags))
	nativeEndian.PutUint16(b[6:8], uint16(m.Index))
	nativeEndian.PutUint32(b[24:28], uint32(m.ID))
	nativeEndian.PutUint32(b[28:32], uint32(m.Seq))
	attrs, err := marshalAddrs(b[sizeofRtMsghdr:], m.Addrs)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if attrs > 0 ***REMOVED***
		nativeEndian.PutUint32(b[12:16], uint32(attrs))
	***REMOVED***
	return b, nil
***REMOVED***

func (*wireFormat) parseRouteMessage(_ RIBType, b []byte) (Message, error) ***REMOVED***
	if len(b) < sizeofRtMsghdr ***REMOVED***
		return nil, errMessageTooShort
	***REMOVED***
	l := int(nativeEndian.Uint16(b[:2]))
	if len(b) < l ***REMOVED***
		return nil, errInvalidMessage
	***REMOVED***
	m := &RouteMessage***REMOVED***
		Version: int(b[2]),
		Type:    int(b[3]),
		Flags:   int(nativeEndian.Uint32(b[16:20])),
		Index:   int(nativeEndian.Uint16(b[6:8])),
		ID:      uintptr(nativeEndian.Uint32(b[24:28])),
		Seq:     int(nativeEndian.Uint32(b[28:32])),
		raw:     b[:l],
	***REMOVED***
	ll := int(nativeEndian.Uint16(b[4:6]))
	if len(b) < ll ***REMOVED***
		return nil, errInvalidMessage
	***REMOVED***
	errno := syscall.Errno(nativeEndian.Uint32(b[32:36]))
	if errno != 0 ***REMOVED***
		m.Err = errno
	***REMOVED***
	as, err := parseAddrs(uint(nativeEndian.Uint32(b[12:16])), parseKernelInetAddr, b[ll:])
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	m.Addrs = as
	return m, nil
***REMOVED***
