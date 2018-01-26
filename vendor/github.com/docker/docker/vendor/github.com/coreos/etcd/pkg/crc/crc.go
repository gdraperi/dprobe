// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package crc provides utility function for cyclic redundancy check
// algorithms.
package crc

import (
	"hash"
	"hash/crc32"
)

// The size of a CRC-32 checksum in bytes.
const Size = 4

type digest struct ***REMOVED***
	crc uint32
	tab *crc32.Table
***REMOVED***

// New creates a new hash.Hash32 computing the CRC-32 checksum
// using the polynomial represented by the Table.
// Modified by xiangli to take a prevcrc.
func New(prev uint32, tab *crc32.Table) hash.Hash32 ***REMOVED*** return &digest***REMOVED***prev, tab***REMOVED*** ***REMOVED***

func (d *digest) Size() int ***REMOVED*** return Size ***REMOVED***

func (d *digest) BlockSize() int ***REMOVED*** return 1 ***REMOVED***

func (d *digest) Reset() ***REMOVED*** d.crc = 0 ***REMOVED***

func (d *digest) Write(p []byte) (n int, err error) ***REMOVED***
	d.crc = crc32.Update(d.crc, d.tab, p)
	return len(p), nil
***REMOVED***

func (d *digest) Sum32() uint32 ***REMOVED*** return d.crc ***REMOVED***

func (d *digest) Sum(in []byte) []byte ***REMOVED***
	s := d.Sum32()
	return append(in, byte(s>>24), byte(s>>16), byte(s>>8), byte(s))
***REMOVED***
