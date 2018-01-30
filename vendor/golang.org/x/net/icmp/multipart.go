// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package icmp

import "golang.org/x/net/internal/iana"

// multipartMessageBodyDataLen takes b as an original datagram and
// exts as extensions, and returns a required length for message body
// and a required length for a padded original datagram in wire
// format.
func multipartMessageBodyDataLen(proto int, b []byte, exts []Extension) (bodyLen, dataLen int) ***REMOVED***
	for _, ext := range exts ***REMOVED***
		bodyLen += ext.Len(proto)
	***REMOVED***
	if bodyLen > 0 ***REMOVED***
		dataLen = multipartMessageOrigDatagramLen(proto, b)
		bodyLen += 4 // length of extension header
	***REMOVED*** else ***REMOVED***
		dataLen = len(b)
	***REMOVED***
	bodyLen += dataLen
	return bodyLen, dataLen
***REMOVED***

// multipartMessageOrigDatagramLen takes b as an original datagram,
// and returns a required length for a padded orignal datagram in wire
// format.
func multipartMessageOrigDatagramLen(proto int, b []byte) int ***REMOVED***
	roundup := func(b []byte, align int) int ***REMOVED***
		// According to RFC 4884, the padded original datagram
		// field must contain at least 128 octets.
		if len(b) < 128 ***REMOVED***
			return 128
		***REMOVED***
		r := len(b)
		return (r + align - 1) & ^(align - 1)
	***REMOVED***
	switch proto ***REMOVED***
	case iana.ProtocolICMP:
		return roundup(b, 4)
	case iana.ProtocolIPv6ICMP:
		return roundup(b, 8)
	default:
		return len(b)
	***REMOVED***
***REMOVED***

// marshalMultipartMessageBody takes data as an original datagram and
// exts as extesnsions, and returns a binary encoding of message body.
// It can be used for non-multipart message bodies when exts is nil.
func marshalMultipartMessageBody(proto int, data []byte, exts []Extension) ([]byte, error) ***REMOVED***
	bodyLen, dataLen := multipartMessageBodyDataLen(proto, data, exts)
	b := make([]byte, 4+bodyLen)
	copy(b[4:], data)
	off := dataLen + 4
	if len(exts) > 0 ***REMOVED***
		b[dataLen+4] = byte(extensionVersion << 4)
		off += 4 // length of object header
		for _, ext := range exts ***REMOVED***
			switch ext := ext.(type) ***REMOVED***
			case *MPLSLabelStack:
				if err := ext.marshal(proto, b[off:]); err != nil ***REMOVED***
					return nil, err
				***REMOVED***
				off += ext.Len(proto)
			case *InterfaceInfo:
				attrs, l := ext.attrsAndLen(proto)
				if err := ext.marshal(proto, b[off:], attrs, l); err != nil ***REMOVED***
					return nil, err
				***REMOVED***
				off += ext.Len(proto)
			***REMOVED***
		***REMOVED***
		s := checksum(b[dataLen+4:])
		b[dataLen+4+2] ^= byte(s)
		b[dataLen+4+3] ^= byte(s >> 8)
		switch proto ***REMOVED***
		case iana.ProtocolICMP:
			b[1] = byte(dataLen / 4)
		case iana.ProtocolIPv6ICMP:
			b[0] = byte(dataLen / 8)
		***REMOVED***
	***REMOVED***
	return b, nil
***REMOVED***

// parseMultipartMessageBody parses b as either a non-multipart
// message body or a multipart message body.
func parseMultipartMessageBody(proto int, b []byte) ([]byte, []Extension, error) ***REMOVED***
	var l int
	switch proto ***REMOVED***
	case iana.ProtocolICMP:
		l = 4 * int(b[1])
	case iana.ProtocolIPv6ICMP:
		l = 8 * int(b[0])
	***REMOVED***
	if len(b) == 4 ***REMOVED***
		return nil, nil, nil
	***REMOVED***
	exts, l, err := parseExtensions(b[4:], l)
	if err != nil ***REMOVED***
		l = len(b) - 4
	***REMOVED***
	data := make([]byte, l)
	copy(data, b[4:])
	return data, exts, nil
***REMOVED***
