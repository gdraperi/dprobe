// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package icmp

import "encoding/binary"

// A PacketTooBig represents an ICMP packet too big message body.
type PacketTooBig struct ***REMOVED***
	MTU  int    // maximum transmission unit of the nexthop link
	Data []byte // data, known as original datagram field
***REMOVED***

// Len implements the Len method of MessageBody interface.
func (p *PacketTooBig) Len(proto int) int ***REMOVED***
	if p == nil ***REMOVED***
		return 0
	***REMOVED***
	return 4 + len(p.Data)
***REMOVED***

// Marshal implements the Marshal method of MessageBody interface.
func (p *PacketTooBig) Marshal(proto int) ([]byte, error) ***REMOVED***
	b := make([]byte, 4+len(p.Data))
	binary.BigEndian.PutUint32(b[:4], uint32(p.MTU))
	copy(b[4:], p.Data)
	return b, nil
***REMOVED***

// parsePacketTooBig parses b as an ICMP packet too big message body.
func parsePacketTooBig(proto int, b []byte) (MessageBody, error) ***REMOVED***
	bodyLen := len(b)
	if bodyLen < 4 ***REMOVED***
		return nil, errMessageTooShort
	***REMOVED***
	p := &PacketTooBig***REMOVED***MTU: int(binary.BigEndian.Uint32(b[:4]))***REMOVED***
	if bodyLen > 4 ***REMOVED***
		p.Data = make([]byte, bodyLen-4)
		copy(p.Data, b[4:])
	***REMOVED***
	return p, nil
***REMOVED***
