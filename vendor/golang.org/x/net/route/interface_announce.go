// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build dragonfly freebsd netbsd

package route

func (w *wireFormat) parseInterfaceAnnounceMessage(_ RIBType, b []byte) (Message, error) ***REMOVED***
	if len(b) < w.bodyOff ***REMOVED***
		return nil, errMessageTooShort
	***REMOVED***
	l := int(nativeEndian.Uint16(b[:2]))
	if len(b) < l ***REMOVED***
		return nil, errInvalidMessage
	***REMOVED***
	m := &InterfaceAnnounceMessage***REMOVED***
		Version: int(b[2]),
		Type:    int(b[3]),
		Index:   int(nativeEndian.Uint16(b[4:6])),
		What:    int(nativeEndian.Uint16(b[22:24])),
		raw:     b[:l],
	***REMOVED***
	for i := 0; i < 16; i++ ***REMOVED***
		if b[6+i] != 0 ***REMOVED***
			continue
		***REMOVED***
		m.Name = string(b[6 : 6+i])
		break
	***REMOVED***
	return m, nil
***REMOVED***
