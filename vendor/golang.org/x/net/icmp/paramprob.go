// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package icmp

import (
	"encoding/binary"
	"golang.org/x/net/internal/iana"
)

// A ParamProb represents an ICMP parameter problem message body.
type ParamProb struct ***REMOVED***
	Pointer    uintptr     // offset within the data where the error was detected
	Data       []byte      // data, known as original datagram field
	Extensions []Extension // extensions
***REMOVED***

// Len implements the Len method of MessageBody interface.
func (p *ParamProb) Len(proto int) int ***REMOVED***
	if p == nil ***REMOVED***
		return 0
	***REMOVED***
	l, _ := multipartMessageBodyDataLen(proto, p.Data, p.Extensions)
	return 4 + l
***REMOVED***

// Marshal implements the Marshal method of MessageBody interface.
func (p *ParamProb) Marshal(proto int) ([]byte, error) ***REMOVED***
	if proto == iana.ProtocolIPv6ICMP ***REMOVED***
		b := make([]byte, p.Len(proto))
		binary.BigEndian.PutUint32(b[:4], uint32(p.Pointer))
		copy(b[4:], p.Data)
		return b, nil
	***REMOVED***
	b, err := marshalMultipartMessageBody(proto, p.Data, p.Extensions)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	b[0] = byte(p.Pointer)
	return b, nil
***REMOVED***

// parseParamProb parses b as an ICMP parameter problem message body.
func parseParamProb(proto int, b []byte) (MessageBody, error) ***REMOVED***
	if len(b) < 4 ***REMOVED***
		return nil, errMessageTooShort
	***REMOVED***
	p := &ParamProb***REMOVED******REMOVED***
	if proto == iana.ProtocolIPv6ICMP ***REMOVED***
		p.Pointer = uintptr(binary.BigEndian.Uint32(b[:4]))
		p.Data = make([]byte, len(b)-4)
		copy(p.Data, b[4:])
		return p, nil
	***REMOVED***
	p.Pointer = uintptr(b[0])
	var err error
	p.Data, p.Extensions, err = parseMultipartMessageBody(proto, b)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return p, nil
***REMOVED***
