// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package icmp

import "encoding/binary"

// An Echo represents an ICMP echo request or reply message body.
type Echo struct ***REMOVED***
	ID   int    // identifier
	Seq  int    // sequence number
	Data []byte // data
***REMOVED***

// Len implements the Len method of MessageBody interface.
func (p *Echo) Len(proto int) int ***REMOVED***
	if p == nil ***REMOVED***
		return 0
	***REMOVED***
	return 4 + len(p.Data)
***REMOVED***

// Marshal implements the Marshal method of MessageBody interface.
func (p *Echo) Marshal(proto int) ([]byte, error) ***REMOVED***
	b := make([]byte, 4+len(p.Data))
	binary.BigEndian.PutUint16(b[:2], uint16(p.ID))
	binary.BigEndian.PutUint16(b[2:4], uint16(p.Seq))
	copy(b[4:], p.Data)
	return b, nil
***REMOVED***

// parseEcho parses b as an ICMP echo request or reply message body.
func parseEcho(proto int, b []byte) (MessageBody, error) ***REMOVED***
	bodyLen := len(b)
	if bodyLen < 4 ***REMOVED***
		return nil, errMessageTooShort
	***REMOVED***
	p := &Echo***REMOVED***ID: int(binary.BigEndian.Uint16(b[:2])), Seq: int(binary.BigEndian.Uint16(b[2:4]))***REMOVED***
	if bodyLen > 4 ***REMOVED***
		p.Data = make([]byte, bodyLen-4)
		copy(p.Data, b[4:])
	***REMOVED***
	return p, nil
***REMOVED***
