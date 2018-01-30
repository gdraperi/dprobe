// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package icmp

// A TimeExceeded represents an ICMP time exceeded message body.
type TimeExceeded struct ***REMOVED***
	Data       []byte      // data, known as original datagram field
	Extensions []Extension // extensions
***REMOVED***

// Len implements the Len method of MessageBody interface.
func (p *TimeExceeded) Len(proto int) int ***REMOVED***
	if p == nil ***REMOVED***
		return 0
	***REMOVED***
	l, _ := multipartMessageBodyDataLen(proto, p.Data, p.Extensions)
	return 4 + l
***REMOVED***

// Marshal implements the Marshal method of MessageBody interface.
func (p *TimeExceeded) Marshal(proto int) ([]byte, error) ***REMOVED***
	return marshalMultipartMessageBody(proto, p.Data, p.Extensions)
***REMOVED***

// parseTimeExceeded parses b as an ICMP time exceeded message body.
func parseTimeExceeded(proto int, b []byte) (MessageBody, error) ***REMOVED***
	if len(b) < 4 ***REMOVED***
		return nil, errMessageTooShort
	***REMOVED***
	p := &TimeExceeded***REMOVED******REMOVED***
	var err error
	p.Data, p.Extensions, err = parseMultipartMessageBody(proto, b)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return p, nil
***REMOVED***
