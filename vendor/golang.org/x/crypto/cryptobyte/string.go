// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package cryptobyte contains types that help with parsing and constructing
// length-prefixed, binary messages, including ASN.1 DER. (The asn1 subpackage
// contains useful ASN.1 constants.)
//
// The String type is for parsing. It wraps a []byte slice and provides helper
// functions for consuming structures, value by value.
//
// The Builder type is for constructing messages. It providers helper functions
// for appending values and also for appending length-prefixed submessages –
// without having to worry about calculating the length prefix ahead of time.
//
// See the documentation and examples for the Builder and String types to get
// started.
package cryptobyte // import "golang.org/x/crypto/cryptobyte"

// String represents a string of bytes. It provides methods for parsing
// fixed-length and length-prefixed values from it.
type String []byte

// read advances a String by n bytes and returns them. If less than n bytes
// remain, it returns nil.
func (s *String) read(n int) []byte ***REMOVED***
	if len(*s) < n ***REMOVED***
		return nil
	***REMOVED***
	v := (*s)[:n]
	*s = (*s)[n:]
	return v
***REMOVED***

// Skip advances the String by n byte and reports whether it was successful.
func (s *String) Skip(n int) bool ***REMOVED***
	return s.read(n) != nil
***REMOVED***

// ReadUint8 decodes an 8-bit value into out and advances over it. It
// returns true on success and false on error.
func (s *String) ReadUint8(out *uint8) bool ***REMOVED***
	v := s.read(1)
	if v == nil ***REMOVED***
		return false
	***REMOVED***
	*out = uint8(v[0])
	return true
***REMOVED***

// ReadUint16 decodes a big-endian, 16-bit value into out and advances over it.
// It returns true on success and false on error.
func (s *String) ReadUint16(out *uint16) bool ***REMOVED***
	v := s.read(2)
	if v == nil ***REMOVED***
		return false
	***REMOVED***
	*out = uint16(v[0])<<8 | uint16(v[1])
	return true
***REMOVED***

// ReadUint24 decodes a big-endian, 24-bit value into out and advances over it.
// It returns true on success and false on error.
func (s *String) ReadUint24(out *uint32) bool ***REMOVED***
	v := s.read(3)
	if v == nil ***REMOVED***
		return false
	***REMOVED***
	*out = uint32(v[0])<<16 | uint32(v[1])<<8 | uint32(v[2])
	return true
***REMOVED***

// ReadUint32 decodes a big-endian, 32-bit value into out and advances over it.
// It returns true on success and false on error.
func (s *String) ReadUint32(out *uint32) bool ***REMOVED***
	v := s.read(4)
	if v == nil ***REMOVED***
		return false
	***REMOVED***
	*out = uint32(v[0])<<24 | uint32(v[1])<<16 | uint32(v[2])<<8 | uint32(v[3])
	return true
***REMOVED***

func (s *String) readUnsigned(out *uint32, length int) bool ***REMOVED***
	v := s.read(length)
	if v == nil ***REMOVED***
		return false
	***REMOVED***
	var result uint32
	for i := 0; i < length; i++ ***REMOVED***
		result <<= 8
		result |= uint32(v[i])
	***REMOVED***
	*out = result
	return true
***REMOVED***

func (s *String) readLengthPrefixed(lenLen int, outChild *String) bool ***REMOVED***
	lenBytes := s.read(lenLen)
	if lenBytes == nil ***REMOVED***
		return false
	***REMOVED***
	var length uint32
	for _, b := range lenBytes ***REMOVED***
		length = length << 8
		length = length | uint32(b)
	***REMOVED***
	if int(length) < 0 ***REMOVED***
		// This currently cannot overflow because we read uint24 at most, but check
		// anyway in case that changes in the future.
		return false
	***REMOVED***
	v := s.read(int(length))
	if v == nil ***REMOVED***
		return false
	***REMOVED***
	*outChild = v
	return true
***REMOVED***

// ReadUint8LengthPrefixed reads the content of an 8-bit length-prefixed value
// into out and advances over it. It returns true on success and false on
// error.
func (s *String) ReadUint8LengthPrefixed(out *String) bool ***REMOVED***
	return s.readLengthPrefixed(1, out)
***REMOVED***

// ReadUint16LengthPrefixed reads the content of a big-endian, 16-bit
// length-prefixed value into out and advances over it. It returns true on
// success and false on error.
func (s *String) ReadUint16LengthPrefixed(out *String) bool ***REMOVED***
	return s.readLengthPrefixed(2, out)
***REMOVED***

// ReadUint24LengthPrefixed reads the content of a big-endian, 24-bit
// length-prefixed value into out and advances over it. It returns true on
// success and false on error.
func (s *String) ReadUint24LengthPrefixed(out *String) bool ***REMOVED***
	return s.readLengthPrefixed(3, out)
***REMOVED***

// ReadBytes reads n bytes into out and advances over them. It returns true on
// success and false and error.
func (s *String) ReadBytes(out *[]byte, n int) bool ***REMOVED***
	v := s.read(n)
	if v == nil ***REMOVED***
		return false
	***REMOVED***
	*out = v
	return true
***REMOVED***

// CopyBytes copies len(out) bytes into out and advances over them. It returns
// true on success and false on error.
func (s *String) CopyBytes(out []byte) bool ***REMOVED***
	n := len(out)
	v := s.read(n)
	if v == nil ***REMOVED***
		return false
	***REMOVED***
	return copy(out, v) == n
***REMOVED***

// Empty reports whether the string does not contain any bytes.
func (s String) Empty() bool ***REMOVED***
	return len(s) == 0
***REMOVED***
