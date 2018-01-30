// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd netbsd openbsd

package route

// A Message represents a routing message.
type Message interface ***REMOVED***
	// Sys returns operating system-specific information.
	Sys() []Sys
***REMOVED***

// A Sys reprensents operating system-specific information.
type Sys interface ***REMOVED***
	// SysType returns a type of operating system-specific
	// information.
	SysType() SysType
***REMOVED***

// A SysType represents a type of operating system-specific
// information.
type SysType int

const (
	SysMetrics SysType = iota
	SysStats
)

// ParseRIB parses b as a routing information base and returns a list
// of routing messages.
func ParseRIB(typ RIBType, b []byte) ([]Message, error) ***REMOVED***
	if !typ.parseable() ***REMOVED***
		return nil, errUnsupportedMessage
	***REMOVED***
	var msgs []Message
	nmsgs, nskips := 0, 0
	for len(b) > 4 ***REMOVED***
		nmsgs++
		l := int(nativeEndian.Uint16(b[:2]))
		if l == 0 ***REMOVED***
			return nil, errInvalidMessage
		***REMOVED***
		if len(b) < l ***REMOVED***
			return nil, errMessageTooShort
		***REMOVED***
		if b[2] != sysRTM_VERSION ***REMOVED***
			b = b[l:]
			continue
		***REMOVED***
		if w, ok := wireFormats[int(b[3])]; !ok ***REMOVED***
			nskips++
		***REMOVED*** else ***REMOVED***
			m, err := w.parse(typ, b)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			if m == nil ***REMOVED***
				nskips++
			***REMOVED*** else ***REMOVED***
				msgs = append(msgs, m)
			***REMOVED***
		***REMOVED***
		b = b[l:]
	***REMOVED***
	// We failed to parse any of the messages - version mismatch?
	if nmsgs != len(msgs)+nskips ***REMOVED***
		return nil, errMessageMismatch
	***REMOVED***
	return msgs, nil
***REMOVED***
