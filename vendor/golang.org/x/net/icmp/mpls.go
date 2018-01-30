// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package icmp

import "encoding/binary"

// A MPLSLabel represents a MPLS label stack entry.
type MPLSLabel struct ***REMOVED***
	Label int  // label value
	TC    int  // traffic class; formerly experimental use
	S     bool // bottom of stack
	TTL   int  // time to live
***REMOVED***

const (
	classMPLSLabelStack        = 1
	typeIncomingMPLSLabelStack = 1
)

// A MPLSLabelStack represents a MPLS label stack.
type MPLSLabelStack struct ***REMOVED***
	Class  int // extension object class number
	Type   int // extension object sub-type
	Labels []MPLSLabel
***REMOVED***

// Len implements the Len method of Extension interface.
func (ls *MPLSLabelStack) Len(proto int) int ***REMOVED***
	return 4 + (4 * len(ls.Labels))
***REMOVED***

// Marshal implements the Marshal method of Extension interface.
func (ls *MPLSLabelStack) Marshal(proto int) ([]byte, error) ***REMOVED***
	b := make([]byte, ls.Len(proto))
	if err := ls.marshal(proto, b); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return b, nil
***REMOVED***

func (ls *MPLSLabelStack) marshal(proto int, b []byte) error ***REMOVED***
	l := ls.Len(proto)
	binary.BigEndian.PutUint16(b[:2], uint16(l))
	b[2], b[3] = classMPLSLabelStack, typeIncomingMPLSLabelStack
	off := 4
	for _, ll := range ls.Labels ***REMOVED***
		b[off], b[off+1], b[off+2] = byte(ll.Label>>12), byte(ll.Label>>4&0xff), byte(ll.Label<<4&0xf0)
		b[off+2] |= byte(ll.TC << 1 & 0x0e)
		if ll.S ***REMOVED***
			b[off+2] |= 0x1
		***REMOVED***
		b[off+3] = byte(ll.TTL)
		off += 4
	***REMOVED***
	return nil
***REMOVED***

func parseMPLSLabelStack(b []byte) (Extension, error) ***REMOVED***
	ls := &MPLSLabelStack***REMOVED***
		Class: int(b[2]),
		Type:  int(b[3]),
	***REMOVED***
	for b = b[4:]; len(b) >= 4; b = b[4:] ***REMOVED***
		ll := MPLSLabel***REMOVED***
			Label: int(b[0])<<12 | int(b[1])<<4 | int(b[2])>>4,
			TC:    int(b[2]&0x0e) >> 1,
			TTL:   int(b[3]),
		***REMOVED***
		if b[2]&0x1 != 0 ***REMOVED***
			ll.S = true
		***REMOVED***
		ls.Labels = append(ls.Labels, ll)
	***REMOVED***
	return ls, nil
***REMOVED***
