// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package blake2b

import (
	"encoding/binary"
	"errors"
	"io"
)

// XOF defines the interface to hash functions that
// support arbitrary-length output.
type XOF interface ***REMOVED***
	// Write absorbs more data into the hash's state. It panics if called
	// after Read.
	io.Writer

	// Read reads more output from the hash. It returns io.EOF if the limit
	// has been reached.
	io.Reader

	// Clone returns a copy of the XOF in its current state.
	Clone() XOF

	// Reset resets the XOF to its initial state.
	Reset()
***REMOVED***

// OutputLengthUnknown can be used as the size argument to NewXOF to indicate
// the the length of the output is not known in advance.
const OutputLengthUnknown = 0

// magicUnknownOutputLength is a magic value for the output size that indicates
// an unknown number of output bytes.
const magicUnknownOutputLength = (1 << 32) - 1

// maxOutputLength is the absolute maximum number of bytes to produce when the
// number of output bytes is unknown.
const maxOutputLength = (1 << 32) * 64

// NewXOF creates a new variable-output-length hash. The hash either produce a
// known number of bytes (1 <= size < 2**32-1), or an unknown number of bytes
// (size == OutputLengthUnknown). In the latter case, an absolute limit of
// 256GiB applies.
//
// A non-nil key turns the hash into a MAC. The key must between
// zero and 32 bytes long.
func NewXOF(size uint32, key []byte) (XOF, error) ***REMOVED***
	if len(key) > Size ***REMOVED***
		return nil, errKeySize
	***REMOVED***
	if size == magicUnknownOutputLength ***REMOVED***
		// 2^32-1 indicates an unknown number of bytes and thus isn't a
		// valid length.
		return nil, errors.New("blake2b: XOF length too large")
	***REMOVED***
	if size == OutputLengthUnknown ***REMOVED***
		size = magicUnknownOutputLength
	***REMOVED***
	x := &xof***REMOVED***
		d: digest***REMOVED***
			size:   Size,
			keyLen: len(key),
		***REMOVED***,
		length: size,
	***REMOVED***
	copy(x.d.key[:], key)
	x.Reset()
	return x, nil
***REMOVED***

type xof struct ***REMOVED***
	d                digest
	length           uint32
	remaining        uint64
	cfg, root, block [Size]byte
	offset           int
	nodeOffset       uint32
	readMode         bool
***REMOVED***

func (x *xof) Write(p []byte) (n int, err error) ***REMOVED***
	if x.readMode ***REMOVED***
		panic("blake2b: write to XOF after read")
	***REMOVED***
	return x.d.Write(p)
***REMOVED***

func (x *xof) Clone() XOF ***REMOVED***
	clone := *x
	return &clone
***REMOVED***

func (x *xof) Reset() ***REMOVED***
	x.cfg[0] = byte(Size)
	binary.LittleEndian.PutUint32(x.cfg[4:], uint32(Size)) // leaf length
	binary.LittleEndian.PutUint32(x.cfg[12:], x.length)    // XOF length
	x.cfg[17] = byte(Size)                                 // inner hash size

	x.d.Reset()
	x.d.h[1] ^= uint64(x.length) << 32

	x.remaining = uint64(x.length)
	if x.remaining == magicUnknownOutputLength ***REMOVED***
		x.remaining = maxOutputLength
	***REMOVED***
	x.offset, x.nodeOffset = 0, 0
	x.readMode = false
***REMOVED***

func (x *xof) Read(p []byte) (n int, err error) ***REMOVED***
	if !x.readMode ***REMOVED***
		x.d.finalize(&x.root)
		x.readMode = true
	***REMOVED***

	if x.remaining == 0 ***REMOVED***
		return 0, io.EOF
	***REMOVED***

	n = len(p)
	if uint64(n) > x.remaining ***REMOVED***
		n = int(x.remaining)
		p = p[:n]
	***REMOVED***

	if x.offset > 0 ***REMOVED***
		blockRemaining := Size - x.offset
		if n < blockRemaining ***REMOVED***
			x.offset += copy(p, x.block[x.offset:])
			x.remaining -= uint64(n)
			return
		***REMOVED***
		copy(p, x.block[x.offset:])
		p = p[blockRemaining:]
		x.offset = 0
		x.remaining -= uint64(blockRemaining)
	***REMOVED***

	for len(p) >= Size ***REMOVED***
		binary.LittleEndian.PutUint32(x.cfg[8:], x.nodeOffset)
		x.nodeOffset++

		x.d.initConfig(&x.cfg)
		x.d.Write(x.root[:])
		x.d.finalize(&x.block)

		copy(p, x.block[:])
		p = p[Size:]
		x.remaining -= uint64(Size)
	***REMOVED***

	if todo := len(p); todo > 0 ***REMOVED***
		if x.remaining < uint64(Size) ***REMOVED***
			x.cfg[0] = byte(x.remaining)
		***REMOVED***
		binary.LittleEndian.PutUint32(x.cfg[8:], x.nodeOffset)
		x.nodeOffset++

		x.d.initConfig(&x.cfg)
		x.d.Write(x.root[:])
		x.d.finalize(&x.block)

		x.offset = copy(p, x.block[:todo])
		x.remaining -= uint64(todo)
	***REMOVED***
	return
***REMOVED***

func (d *digest) initConfig(cfg *[Size]byte) ***REMOVED***
	d.offset, d.c[0], d.c[1] = 0, 0, 0
	for i := range d.h ***REMOVED***
		d.h[i] = iv[i] ^ binary.LittleEndian.Uint64(cfg[i*8:])
	***REMOVED***
***REMOVED***
