// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipv4

import (
	"encoding/binary"
	"fmt"
	"net"
	"runtime"
	"syscall"

	"golang.org/x/net/internal/socket"
)

const (
	Version      = 4  // protocol version
	HeaderLen    = 20 // header length without extension headers
	maxHeaderLen = 60 // sensible default, revisit if later RFCs define new usage of version and header length fields
)

type HeaderFlags int

const (
	MoreFragments HeaderFlags = 1 << iota // more fragments flag
	DontFragment                          // don't fragment flag
)

// A Header represents an IPv4 header.
type Header struct ***REMOVED***
	Version  int         // protocol version
	Len      int         // header length
	TOS      int         // type-of-service
	TotalLen int         // packet total length
	ID       int         // identification
	Flags    HeaderFlags // flags
	FragOff  int         // fragment offset
	TTL      int         // time-to-live
	Protocol int         // next protocol
	Checksum int         // checksum
	Src      net.IP      // source address
	Dst      net.IP      // destination address
	Options  []byte      // options, extension headers
***REMOVED***

func (h *Header) String() string ***REMOVED***
	if h == nil ***REMOVED***
		return "<nil>"
	***REMOVED***
	return fmt.Sprintf("ver=%d hdrlen=%d tos=%#x totallen=%d id=%#x flags=%#x fragoff=%#x ttl=%d proto=%d cksum=%#x src=%v dst=%v", h.Version, h.Len, h.TOS, h.TotalLen, h.ID, h.Flags, h.FragOff, h.TTL, h.Protocol, h.Checksum, h.Src, h.Dst)
***REMOVED***

// Marshal returns the binary encoding of h.
func (h *Header) Marshal() ([]byte, error) ***REMOVED***
	if h == nil ***REMOVED***
		return nil, syscall.EINVAL
	***REMOVED***
	if h.Len < HeaderLen ***REMOVED***
		return nil, errHeaderTooShort
	***REMOVED***
	hdrlen := HeaderLen + len(h.Options)
	b := make([]byte, hdrlen)
	b[0] = byte(Version<<4 | (hdrlen >> 2 & 0x0f))
	b[1] = byte(h.TOS)
	flagsAndFragOff := (h.FragOff & 0x1fff) | int(h.Flags<<13)
	switch runtime.GOOS ***REMOVED***
	case "darwin", "dragonfly", "netbsd":
		socket.NativeEndian.PutUint16(b[2:4], uint16(h.TotalLen))
		socket.NativeEndian.PutUint16(b[6:8], uint16(flagsAndFragOff))
	case "freebsd":
		if freebsdVersion < 1100000 ***REMOVED***
			socket.NativeEndian.PutUint16(b[2:4], uint16(h.TotalLen))
			socket.NativeEndian.PutUint16(b[6:8], uint16(flagsAndFragOff))
		***REMOVED*** else ***REMOVED***
			binary.BigEndian.PutUint16(b[2:4], uint16(h.TotalLen))
			binary.BigEndian.PutUint16(b[6:8], uint16(flagsAndFragOff))
		***REMOVED***
	default:
		binary.BigEndian.PutUint16(b[2:4], uint16(h.TotalLen))
		binary.BigEndian.PutUint16(b[6:8], uint16(flagsAndFragOff))
	***REMOVED***
	binary.BigEndian.PutUint16(b[4:6], uint16(h.ID))
	b[8] = byte(h.TTL)
	b[9] = byte(h.Protocol)
	binary.BigEndian.PutUint16(b[10:12], uint16(h.Checksum))
	if ip := h.Src.To4(); ip != nil ***REMOVED***
		copy(b[12:16], ip[:net.IPv4len])
	***REMOVED***
	if ip := h.Dst.To4(); ip != nil ***REMOVED***
		copy(b[16:20], ip[:net.IPv4len])
	***REMOVED*** else ***REMOVED***
		return nil, errMissingAddress
	***REMOVED***
	if len(h.Options) > 0 ***REMOVED***
		copy(b[HeaderLen:], h.Options)
	***REMOVED***
	return b, nil
***REMOVED***

// Parse parses b as an IPv4 header and sotres the result in h.
func (h *Header) Parse(b []byte) error ***REMOVED***
	if h == nil || len(b) < HeaderLen ***REMOVED***
		return errHeaderTooShort
	***REMOVED***
	hdrlen := int(b[0]&0x0f) << 2
	if hdrlen > len(b) ***REMOVED***
		return errBufferTooShort
	***REMOVED***
	h.Version = int(b[0] >> 4)
	h.Len = hdrlen
	h.TOS = int(b[1])
	h.ID = int(binary.BigEndian.Uint16(b[4:6]))
	h.TTL = int(b[8])
	h.Protocol = int(b[9])
	h.Checksum = int(binary.BigEndian.Uint16(b[10:12]))
	h.Src = net.IPv4(b[12], b[13], b[14], b[15])
	h.Dst = net.IPv4(b[16], b[17], b[18], b[19])
	switch runtime.GOOS ***REMOVED***
	case "darwin", "dragonfly", "netbsd":
		h.TotalLen = int(socket.NativeEndian.Uint16(b[2:4])) + hdrlen
		h.FragOff = int(socket.NativeEndian.Uint16(b[6:8]))
	case "freebsd":
		if freebsdVersion < 1100000 ***REMOVED***
			h.TotalLen = int(socket.NativeEndian.Uint16(b[2:4]))
			if freebsdVersion < 1000000 ***REMOVED***
				h.TotalLen += hdrlen
			***REMOVED***
			h.FragOff = int(socket.NativeEndian.Uint16(b[6:8]))
		***REMOVED*** else ***REMOVED***
			h.TotalLen = int(binary.BigEndian.Uint16(b[2:4]))
			h.FragOff = int(binary.BigEndian.Uint16(b[6:8]))
		***REMOVED***
	default:
		h.TotalLen = int(binary.BigEndian.Uint16(b[2:4]))
		h.FragOff = int(binary.BigEndian.Uint16(b[6:8]))
	***REMOVED***
	h.Flags = HeaderFlags(h.FragOff&0xe000) >> 13
	h.FragOff = h.FragOff & 0x1fff
	optlen := hdrlen - HeaderLen
	if optlen > 0 && len(b) >= hdrlen ***REMOVED***
		if cap(h.Options) < optlen ***REMOVED***
			h.Options = make([]byte, optlen)
		***REMOVED*** else ***REMOVED***
			h.Options = h.Options[:optlen]
		***REMOVED***
		copy(h.Options, b[HeaderLen:hdrlen])
	***REMOVED***
	return nil
***REMOVED***

// ParseHeader parses b as an IPv4 header.
func ParseHeader(b []byte) (*Header, error) ***REMOVED***
	h := new(Header)
	if err := h.Parse(b); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return h, nil
***REMOVED***
