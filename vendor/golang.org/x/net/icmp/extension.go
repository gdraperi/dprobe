// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package icmp

import "encoding/binary"

// An Extension represents an ICMP extension.
type Extension interface ***REMOVED***
	// Len returns the length of ICMP extension.
	// Proto must be either the ICMPv4 or ICMPv6 protocol number.
	Len(proto int) int

	// Marshal returns the binary encoding of ICMP extension.
	// Proto must be either the ICMPv4 or ICMPv6 protocol number.
	Marshal(proto int) ([]byte, error)
***REMOVED***

const extensionVersion = 2

func validExtensionHeader(b []byte) bool ***REMOVED***
	v := int(b[0]&0xf0) >> 4
	s := binary.BigEndian.Uint16(b[2:4])
	if s != 0 ***REMOVED***
		s = checksum(b)
	***REMOVED***
	if v != extensionVersion || s != 0 ***REMOVED***
		return false
	***REMOVED***
	return true
***REMOVED***

// parseExtensions parses b as a list of ICMP extensions.
// The length attribute l must be the length attribute field in
// received icmp messages.
//
// It will return a list of ICMP extensions and an adjusted length
// attribute that represents the length of the padded original
// datagram field. Otherwise, it returns an error.
func parseExtensions(b []byte, l int) ([]Extension, int, error) ***REMOVED***
	// Still a lot of non-RFC 4884 compliant implementations are
	// out there. Set the length attribute l to 128 when it looks
	// inappropriate for backwards compatibility.
	//
	// A minimal extension at least requires 8 octets; 4 octets
	// for an extension header, and 4 octets for a single object
	// header.
	//
	// See RFC 4884 for further information.
	if 128 > l || l+8 > len(b) ***REMOVED***
		l = 128
	***REMOVED***
	if l+8 > len(b) ***REMOVED***
		return nil, -1, errNoExtension
	***REMOVED***
	if !validExtensionHeader(b[l:]) ***REMOVED***
		if l == 128 ***REMOVED***
			return nil, -1, errNoExtension
		***REMOVED***
		l = 128
		if !validExtensionHeader(b[l:]) ***REMOVED***
			return nil, -1, errNoExtension
		***REMOVED***
	***REMOVED***
	var exts []Extension
	for b = b[l+4:]; len(b) >= 4; ***REMOVED***
		ol := int(binary.BigEndian.Uint16(b[:2]))
		if 4 > ol || ol > len(b) ***REMOVED***
			break
		***REMOVED***
		switch b[2] ***REMOVED***
		case classMPLSLabelStack:
			ext, err := parseMPLSLabelStack(b[:ol])
			if err != nil ***REMOVED***
				return nil, -1, err
			***REMOVED***
			exts = append(exts, ext)
		case classInterfaceInfo:
			ext, err := parseInterfaceInfo(b[:ol])
			if err != nil ***REMOVED***
				return nil, -1, err
			***REMOVED***
			exts = append(exts, ext)
		***REMOVED***
		b = b[ol:]
	***REMOVED***
	return exts, l, nil
***REMOVED***
