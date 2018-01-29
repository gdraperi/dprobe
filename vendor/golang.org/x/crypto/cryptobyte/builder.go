// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cryptobyte

import (
	"errors"
	"fmt"
)

// A Builder builds byte strings from fixed-length and length-prefixed values.
// Builders either allocate space as needed, or are ‘fixed’, which means that
// they write into a given buffer and produce an error if it's exhausted.
//
// The zero value is a usable Builder that allocates space as needed.
//
// Simple values are marshaled and appended to a Builder using methods on the
// Builder. Length-prefixed values are marshaled by providing a
// BuilderContinuation, which is a function that writes the inner contents of
// the value to a given Builder. See the documentation for BuilderContinuation
// for details.
type Builder struct ***REMOVED***
	err            error
	result         []byte
	fixedSize      bool
	child          *Builder
	offset         int
	pendingLenLen  int
	pendingIsASN1  bool
	inContinuation *bool
***REMOVED***

// NewBuilder creates a Builder that appends its output to the given buffer.
// Like append(), the slice will be reallocated if its capacity is exceeded.
// Use Bytes to get the final buffer.
func NewBuilder(buffer []byte) *Builder ***REMOVED***
	return &Builder***REMOVED***
		result: buffer,
	***REMOVED***
***REMOVED***

// NewFixedBuilder creates a Builder that appends its output into the given
// buffer. This builder does not reallocate the output buffer. Writes that
// would exceed the buffer's capacity are treated as an error.
func NewFixedBuilder(buffer []byte) *Builder ***REMOVED***
	return &Builder***REMOVED***
		result:    buffer,
		fixedSize: true,
	***REMOVED***
***REMOVED***

// Bytes returns the bytes written by the builder or an error if one has
// occurred during during building.
func (b *Builder) Bytes() ([]byte, error) ***REMOVED***
	if b.err != nil ***REMOVED***
		return nil, b.err
	***REMOVED***
	return b.result[b.offset:], nil
***REMOVED***

// BytesOrPanic returns the bytes written by the builder or panics if an error
// has occurred during building.
func (b *Builder) BytesOrPanic() []byte ***REMOVED***
	if b.err != nil ***REMOVED***
		panic(b.err)
	***REMOVED***
	return b.result[b.offset:]
***REMOVED***

// AddUint8 appends an 8-bit value to the byte string.
func (b *Builder) AddUint8(v uint8) ***REMOVED***
	b.add(byte(v))
***REMOVED***

// AddUint16 appends a big-endian, 16-bit value to the byte string.
func (b *Builder) AddUint16(v uint16) ***REMOVED***
	b.add(byte(v>>8), byte(v))
***REMOVED***

// AddUint24 appends a big-endian, 24-bit value to the byte string. The highest
// byte of the 32-bit input value is silently truncated.
func (b *Builder) AddUint24(v uint32) ***REMOVED***
	b.add(byte(v>>16), byte(v>>8), byte(v))
***REMOVED***

// AddUint32 appends a big-endian, 32-bit value to the byte string.
func (b *Builder) AddUint32(v uint32) ***REMOVED***
	b.add(byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
***REMOVED***

// AddBytes appends a sequence of bytes to the byte string.
func (b *Builder) AddBytes(v []byte) ***REMOVED***
	b.add(v...)
***REMOVED***

// BuilderContinuation is continuation-passing interface for building
// length-prefixed byte sequences. Builder methods for length-prefixed
// sequences (AddUint8LengthPrefixed etc) will invoke the BuilderContinuation
// supplied to them. The child builder passed to the continuation can be used
// to build the content of the length-prefixed sequence. For example:
//
//   parent := cryptobyte.NewBuilder()
//   parent.AddUint8LengthPrefixed(func (child *Builder) ***REMOVED***
//     child.AddUint8(42)
//     child.AddUint8LengthPrefixed(func (grandchild *Builder) ***REMOVED***
//       grandchild.AddUint8(5)
// ***REMOVED***)
//   ***REMOVED***)
//
// It is an error to write more bytes to the child than allowed by the reserved
// length prefix. After the continuation returns, the child must be considered
// invalid, i.e. users must not store any copies or references of the child
// that outlive the continuation.
//
// If the continuation panics with a value of type BuildError then the inner
// error will be returned as the error from Bytes. If the child panics
// otherwise then Bytes will repanic with the same value.
type BuilderContinuation func(child *Builder)

// BuildError wraps an error. If a BuilderContinuation panics with this value,
// the panic will be recovered and the inner error will be returned from
// Builder.Bytes.
type BuildError struct ***REMOVED***
	Err error
***REMOVED***

// AddUint8LengthPrefixed adds a 8-bit length-prefixed byte sequence.
func (b *Builder) AddUint8LengthPrefixed(f BuilderContinuation) ***REMOVED***
	b.addLengthPrefixed(1, false, f)
***REMOVED***

// AddUint16LengthPrefixed adds a big-endian, 16-bit length-prefixed byte sequence.
func (b *Builder) AddUint16LengthPrefixed(f BuilderContinuation) ***REMOVED***
	b.addLengthPrefixed(2, false, f)
***REMOVED***

// AddUint24LengthPrefixed adds a big-endian, 24-bit length-prefixed byte sequence.
func (b *Builder) AddUint24LengthPrefixed(f BuilderContinuation) ***REMOVED***
	b.addLengthPrefixed(3, false, f)
***REMOVED***

// AddUint32LengthPrefixed adds a big-endian, 32-bit length-prefixed byte sequence.
func (b *Builder) AddUint32LengthPrefixed(f BuilderContinuation) ***REMOVED***
	b.addLengthPrefixed(4, false, f)
***REMOVED***

func (b *Builder) callContinuation(f BuilderContinuation, arg *Builder) ***REMOVED***
	if !*b.inContinuation ***REMOVED***
		*b.inContinuation = true

		defer func() ***REMOVED***
			*b.inContinuation = false

			r := recover()
			if r == nil ***REMOVED***
				return
			***REMOVED***

			if buildError, ok := r.(BuildError); ok ***REMOVED***
				b.err = buildError.Err
			***REMOVED*** else ***REMOVED***
				panic(r)
			***REMOVED***
		***REMOVED***()
	***REMOVED***

	f(arg)
***REMOVED***

func (b *Builder) addLengthPrefixed(lenLen int, isASN1 bool, f BuilderContinuation) ***REMOVED***
	// Subsequent writes can be ignored if the builder has encountered an error.
	if b.err != nil ***REMOVED***
		return
	***REMOVED***

	offset := len(b.result)
	b.add(make([]byte, lenLen)...)

	if b.inContinuation == nil ***REMOVED***
		b.inContinuation = new(bool)
	***REMOVED***

	b.child = &Builder***REMOVED***
		result:         b.result,
		fixedSize:      b.fixedSize,
		offset:         offset,
		pendingLenLen:  lenLen,
		pendingIsASN1:  isASN1,
		inContinuation: b.inContinuation,
	***REMOVED***

	b.callContinuation(f, b.child)
	b.flushChild()
	if b.child != nil ***REMOVED***
		panic("cryptobyte: internal error")
	***REMOVED***
***REMOVED***

func (b *Builder) flushChild() ***REMOVED***
	if b.child == nil ***REMOVED***
		return
	***REMOVED***
	b.child.flushChild()
	child := b.child
	b.child = nil

	if child.err != nil ***REMOVED***
		b.err = child.err
		return
	***REMOVED***

	length := len(child.result) - child.pendingLenLen - child.offset

	if length < 0 ***REMOVED***
		panic("cryptobyte: internal error") // result unexpectedly shrunk
	***REMOVED***

	if child.pendingIsASN1 ***REMOVED***
		// For ASN.1, we reserved a single byte for the length. If that turned out
		// to be incorrect, we have to move the contents along in order to make
		// space.
		if child.pendingLenLen != 1 ***REMOVED***
			panic("cryptobyte: internal error")
		***REMOVED***
		var lenLen, lenByte uint8
		if int64(length) > 0xfffffffe ***REMOVED***
			b.err = errors.New("pending ASN.1 child too long")
			return
		***REMOVED*** else if length > 0xffffff ***REMOVED***
			lenLen = 5
			lenByte = 0x80 | 4
		***REMOVED*** else if length > 0xffff ***REMOVED***
			lenLen = 4
			lenByte = 0x80 | 3
		***REMOVED*** else if length > 0xff ***REMOVED***
			lenLen = 3
			lenByte = 0x80 | 2
		***REMOVED*** else if length > 0x7f ***REMOVED***
			lenLen = 2
			lenByte = 0x80 | 1
		***REMOVED*** else ***REMOVED***
			lenLen = 1
			lenByte = uint8(length)
			length = 0
		***REMOVED***

		// Insert the initial length byte, make space for successive length bytes,
		// and adjust the offset.
		child.result[child.offset] = lenByte
		extraBytes := int(lenLen - 1)
		if extraBytes != 0 ***REMOVED***
			child.add(make([]byte, extraBytes)...)
			childStart := child.offset + child.pendingLenLen
			copy(child.result[childStart+extraBytes:], child.result[childStart:])
		***REMOVED***
		child.offset++
		child.pendingLenLen = extraBytes
	***REMOVED***

	l := length
	for i := child.pendingLenLen - 1; i >= 0; i-- ***REMOVED***
		child.result[child.offset+i] = uint8(l)
		l >>= 8
	***REMOVED***
	if l != 0 ***REMOVED***
		b.err = fmt.Errorf("cryptobyte: pending child length %d exceeds %d-byte length prefix", length, child.pendingLenLen)
		return
	***REMOVED***

	if !b.fixedSize ***REMOVED***
		b.result = child.result // In case child reallocated result.
	***REMOVED***
***REMOVED***

func (b *Builder) add(bytes ...byte) ***REMOVED***
	if b.err != nil ***REMOVED***
		return
	***REMOVED***
	if b.child != nil ***REMOVED***
		panic("attempted write while child is pending")
	***REMOVED***
	if len(b.result)+len(bytes) < len(bytes) ***REMOVED***
		b.err = errors.New("cryptobyte: length overflow")
	***REMOVED***
	if b.fixedSize && len(b.result)+len(bytes) > cap(b.result) ***REMOVED***
		b.err = errors.New("cryptobyte: Builder is exceeding its fixed-size buffer")
		return
	***REMOVED***
	b.result = append(b.result, bytes...)
***REMOVED***

// A MarshalingValue marshals itself into a Builder.
type MarshalingValue interface ***REMOVED***
	// Marshal is called by Builder.AddValue. It receives a pointer to a builder
	// to marshal itself into. It may return an error that occurred during
	// marshaling, such as unset or invalid values.
	Marshal(b *Builder) error
***REMOVED***

// AddValue calls Marshal on v, passing a pointer to the builder to append to.
// If Marshal returns an error, it is set on the Builder so that subsequent
// appends don't have an effect.
func (b *Builder) AddValue(v MarshalingValue) ***REMOVED***
	err := v.Marshal(b)
	if err != nil ***REMOVED***
		b.err = err
	***REMOVED***
***REMOVED***
