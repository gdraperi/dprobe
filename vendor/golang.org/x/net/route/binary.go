// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd netbsd openbsd

package route

// This file contains duplicates of encoding/binary package.
//
// This package is supposed to be used by the net package of standard
// library. Therefore the package set used in the package must be the
// same as net package.

var (
	littleEndian binaryLittleEndian
	bigEndian    binaryBigEndian
)

type binaryByteOrder interface ***REMOVED***
	Uint16([]byte) uint16
	Uint32([]byte) uint32
	PutUint16([]byte, uint16)
	PutUint32([]byte, uint32)
	Uint64([]byte) uint64
***REMOVED***

type binaryLittleEndian struct***REMOVED******REMOVED***

func (binaryLittleEndian) Uint16(b []byte) uint16 ***REMOVED***
	_ = b[1] // bounds check hint to compiler; see golang.org/issue/14808
	return uint16(b[0]) | uint16(b[1])<<8
***REMOVED***

func (binaryLittleEndian) PutUint16(b []byte, v uint16) ***REMOVED***
	_ = b[1] // early bounds check to guarantee safety of writes below
	b[0] = byte(v)
	b[1] = byte(v >> 8)
***REMOVED***

func (binaryLittleEndian) Uint32(b []byte) uint32 ***REMOVED***
	_ = b[3] // bounds check hint to compiler; see golang.org/issue/14808
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
***REMOVED***

func (binaryLittleEndian) PutUint32(b []byte, v uint32) ***REMOVED***
	_ = b[3] // early bounds check to guarantee safety of writes below
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
***REMOVED***

func (binaryLittleEndian) Uint64(b []byte) uint64 ***REMOVED***
	_ = b[7] // bounds check hint to compiler; see golang.org/issue/14808
	return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
		uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
***REMOVED***

type binaryBigEndian struct***REMOVED******REMOVED***

func (binaryBigEndian) Uint16(b []byte) uint16 ***REMOVED***
	_ = b[1] // bounds check hint to compiler; see golang.org/issue/14808
	return uint16(b[1]) | uint16(b[0])<<8
***REMOVED***

func (binaryBigEndian) PutUint16(b []byte, v uint16) ***REMOVED***
	_ = b[1] // early bounds check to guarantee safety of writes below
	b[0] = byte(v >> 8)
	b[1] = byte(v)
***REMOVED***

func (binaryBigEndian) Uint32(b []byte) uint32 ***REMOVED***
	_ = b[3] // bounds check hint to compiler; see golang.org/issue/14808
	return uint32(b[3]) | uint32(b[2])<<8 | uint32(b[1])<<16 | uint32(b[0])<<24
***REMOVED***

func (binaryBigEndian) PutUint32(b []byte, v uint32) ***REMOVED***
	_ = b[3] // early bounds check to guarantee safety of writes below
	b[0] = byte(v >> 24)
	b[1] = byte(v >> 16)
	b[2] = byte(v >> 8)
	b[3] = byte(v)
***REMOVED***

func (binaryBigEndian) Uint64(b []byte) uint64 ***REMOVED***
	_ = b[7] // bounds check hint to compiler; see golang.org/issue/14808
	return uint64(b[7]) | uint64(b[6])<<8 | uint64(b[5])<<16 | uint64(b[4])<<24 |
		uint64(b[3])<<32 | uint64(b[2])<<40 | uint64(b[1])<<48 | uint64(b[0])<<56
***REMOVED***
