// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd netbsd

package route

import (
	"runtime"
	"syscall"
)

func (m *RouteMessage) marshal() ([]byte, error) ***REMOVED***
	w, ok := wireFormats[m.Type]
	if !ok ***REMOVED***
		return nil, errUnsupportedMessage
	***REMOVED***
	l := w.bodyOff + addrsSpace(m.Addrs)
	if runtime.GOOS == "darwin" ***REMOVED***
		// Fix stray pointer writes on macOS.
		// See golang.org/issue/22456.
		l += 1024
	***REMOVED***
	b := make([]byte, l)
	nativeEndian.PutUint16(b[:2], uint16(l))
	if m.Version == 0 ***REMOVED***
		b[2] = sysRTM_VERSION
	***REMOVED*** else ***REMOVED***
		b[2] = byte(m.Version)
	***REMOVED***
	b[3] = byte(m.Type)
	nativeEndian.PutUint32(b[8:12], uint32(m.Flags))
	nativeEndian.PutUint16(b[4:6], uint16(m.Index))
	nativeEndian.PutUint32(b[16:20], uint32(m.ID))
	nativeEndian.PutUint32(b[20:24], uint32(m.Seq))
	attrs, err := marshalAddrs(b[w.bodyOff:], m.Addrs)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if attrs > 0 ***REMOVED***
		nativeEndian.PutUint32(b[12:16], uint32(attrs))
	***REMOVED***
	return b, nil
***REMOVED***

func (w *wireFormat) parseRouteMessage(typ RIBType, b []byte) (Message, error) ***REMOVED***
	if len(b) < w.bodyOff ***REMOVED***
		return nil, errMessageTooShort
	***REMOVED***
	l := int(nativeEndian.Uint16(b[:2]))
	if len(b) < l ***REMOVED***
		return nil, errInvalidMessage
	***REMOVED***
	m := &RouteMessage***REMOVED***
		Version: int(b[2]),
		Type:    int(b[3]),
		Flags:   int(nativeEndian.Uint32(b[8:12])),
		Index:   int(nativeEndian.Uint16(b[4:6])),
		ID:      uintptr(nativeEndian.Uint32(b[16:20])),
		Seq:     int(nativeEndian.Uint32(b[20:24])),
		extOff:  w.extOff,
		raw:     b[:l],
	***REMOVED***
	errno := syscall.Errno(nativeEndian.Uint32(b[28:32]))
	if errno != 0 ***REMOVED***
		m.Err = errno
	***REMOVED***
	var err error
	m.Addrs, err = parseAddrs(uint(nativeEndian.Uint32(b[12:16])), parseKernelInetAddr, b[w.bodyOff:])
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return m, nil
***REMOVED***
