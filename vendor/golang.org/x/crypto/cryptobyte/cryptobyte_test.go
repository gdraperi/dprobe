// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cryptobyte

import (
	"bytes"
	"errors"
	"fmt"
	"testing"
)

func builderBytesEq(b *Builder, want ...byte) error ***REMOVED***
	got := b.BytesOrPanic()
	if !bytes.Equal(got, want) ***REMOVED***
		return fmt.Errorf("Bytes() = %v, want %v", got, want)
	***REMOVED***
	return nil
***REMOVED***

func TestContinuationError(t *testing.T) ***REMOVED***
	const errorStr = "TestContinuationError"
	var b Builder
	b.AddUint8LengthPrefixed(func(b *Builder) ***REMOVED***
		b.AddUint8(1)
		panic(BuildError***REMOVED***Err: errors.New(errorStr)***REMOVED***)
	***REMOVED***)

	ret, err := b.Bytes()
	if ret != nil ***REMOVED***
		t.Error("expected nil result")
	***REMOVED***
	if err == nil ***REMOVED***
		t.Fatal("unexpected nil error")
	***REMOVED***
	if s := err.Error(); s != errorStr ***REMOVED***
		t.Errorf("expected error %q, got %v", errorStr, s)
	***REMOVED***
***REMOVED***

func TestContinuationNonError(t *testing.T) ***REMOVED***
	defer func() ***REMOVED***
		recover()
	***REMOVED***()

	var b Builder
	b.AddUint8LengthPrefixed(func(b *Builder) ***REMOVED***
		b.AddUint8(1)
		panic(1)
	***REMOVED***)

	t.Error("Builder did not panic")
***REMOVED***

func TestGeneratedPanic(t *testing.T) ***REMOVED***
	defer func() ***REMOVED***
		recover()
	***REMOVED***()

	var b Builder
	b.AddUint8LengthPrefixed(func(b *Builder) ***REMOVED***
		var p *byte
		*p = 0
	***REMOVED***)

	t.Error("Builder did not panic")
***REMOVED***

func TestBytes(t *testing.T) ***REMOVED***
	var b Builder
	v := []byte("foobarbaz")
	b.AddBytes(v[0:3])
	b.AddBytes(v[3:4])
	b.AddBytes(v[4:9])
	if err := builderBytesEq(&b, v...); err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***
	s := String(b.BytesOrPanic())
	for _, w := range []string***REMOVED***"foo", "bar", "baz"***REMOVED*** ***REMOVED***
		var got []byte
		if !s.ReadBytes(&got, 3) ***REMOVED***
			t.Errorf("ReadBytes() = false, want true (w = %v)", w)
		***REMOVED***
		want := []byte(w)
		if !bytes.Equal(got, want) ***REMOVED***
			t.Errorf("ReadBytes(): got = %v, want %v", got, want)
		***REMOVED***
	***REMOVED***
	if len(s) != 0 ***REMOVED***
		t.Errorf("len(s) = %d, want 0", len(s))
	***REMOVED***
***REMOVED***

func TestUint8(t *testing.T) ***REMOVED***
	var b Builder
	b.AddUint8(42)
	if err := builderBytesEq(&b, 42); err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***

	var s String = b.BytesOrPanic()
	var v uint8
	if !s.ReadUint8(&v) ***REMOVED***
		t.Error("ReadUint8() = false, want true")
	***REMOVED***
	if v != 42 ***REMOVED***
		t.Errorf("v = %d, want 42", v)
	***REMOVED***
	if len(s) != 0 ***REMOVED***
		t.Errorf("len(s) = %d, want 0", len(s))
	***REMOVED***
***REMOVED***

func TestUint16(t *testing.T) ***REMOVED***
	var b Builder
	b.AddUint16(65534)
	if err := builderBytesEq(&b, 255, 254); err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***
	var s String = b.BytesOrPanic()
	var v uint16
	if !s.ReadUint16(&v) ***REMOVED***
		t.Error("ReadUint16() == false, want true")
	***REMOVED***
	if v != 65534 ***REMOVED***
		t.Errorf("v = %d, want 65534", v)
	***REMOVED***
	if len(s) != 0 ***REMOVED***
		t.Errorf("len(s) = %d, want 0", len(s))
	***REMOVED***
***REMOVED***

func TestUint24(t *testing.T) ***REMOVED***
	var b Builder
	b.AddUint24(0xfffefd)
	if err := builderBytesEq(&b, 255, 254, 253); err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***

	var s String = b.BytesOrPanic()
	var v uint32
	if !s.ReadUint24(&v) ***REMOVED***
		t.Error("ReadUint8() = false, want true")
	***REMOVED***
	if v != 0xfffefd ***REMOVED***
		t.Errorf("v = %d, want fffefd", v)
	***REMOVED***
	if len(s) != 0 ***REMOVED***
		t.Errorf("len(s) = %d, want 0", len(s))
	***REMOVED***
***REMOVED***

func TestUint24Truncation(t *testing.T) ***REMOVED***
	var b Builder
	b.AddUint24(0x10111213)
	if err := builderBytesEq(&b, 0x11, 0x12, 0x13); err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***
***REMOVED***

func TestUint32(t *testing.T) ***REMOVED***
	var b Builder
	b.AddUint32(0xfffefdfc)
	if err := builderBytesEq(&b, 255, 254, 253, 252); err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***

	var s String = b.BytesOrPanic()
	var v uint32
	if !s.ReadUint32(&v) ***REMOVED***
		t.Error("ReadUint8() = false, want true")
	***REMOVED***
	if v != 0xfffefdfc ***REMOVED***
		t.Errorf("v = %x, want fffefdfc", v)
	***REMOVED***
	if len(s) != 0 ***REMOVED***
		t.Errorf("len(s) = %d, want 0", len(s))
	***REMOVED***
***REMOVED***

func TestUMultiple(t *testing.T) ***REMOVED***
	var b Builder
	b.AddUint8(23)
	b.AddUint32(0xfffefdfc)
	b.AddUint16(42)
	if err := builderBytesEq(&b, 23, 255, 254, 253, 252, 0, 42); err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***

	var s String = b.BytesOrPanic()
	var (
		x uint8
		y uint32
		z uint16
	)
	if !s.ReadUint8(&x) || !s.ReadUint32(&y) || !s.ReadUint16(&z) ***REMOVED***
		t.Error("ReadUint8() = false, want true")
	***REMOVED***
	if x != 23 || y != 0xfffefdfc || z != 42 ***REMOVED***
		t.Errorf("x, y, z = %d, %d, %d; want 23, 4294901244, 5", x, y, z)
	***REMOVED***
	if len(s) != 0 ***REMOVED***
		t.Errorf("len(s) = %d, want 0", len(s))
	***REMOVED***
***REMOVED***

func TestUint8LengthPrefixedSimple(t *testing.T) ***REMOVED***
	var b Builder
	b.AddUint8LengthPrefixed(func(c *Builder) ***REMOVED***
		c.AddUint8(23)
		c.AddUint8(42)
	***REMOVED***)
	if err := builderBytesEq(&b, 2, 23, 42); err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***

	var base, child String = b.BytesOrPanic(), nil
	var x, y uint8
	if !base.ReadUint8LengthPrefixed(&child) || !child.ReadUint8(&x) ||
		!child.ReadUint8(&y) ***REMOVED***
		t.Error("parsing failed")
	***REMOVED***
	if x != 23 || y != 42 ***REMOVED***
		t.Errorf("want x, y == 23, 42; got %d, %d", x, y)
	***REMOVED***
	if len(base) != 0 ***REMOVED***
		t.Errorf("len(base) = %d, want 0", len(base))
	***REMOVED***
	if len(child) != 0 ***REMOVED***
		t.Errorf("len(child) = %d, want 0", len(child))
	***REMOVED***
***REMOVED***

func TestUint8LengthPrefixedMulti(t *testing.T) ***REMOVED***
	var b Builder
	b.AddUint8LengthPrefixed(func(c *Builder) ***REMOVED***
		c.AddUint8(23)
		c.AddUint8(42)
	***REMOVED***)
	b.AddUint8(5)
	b.AddUint8LengthPrefixed(func(c *Builder) ***REMOVED***
		c.AddUint8(123)
		c.AddUint8(234)
	***REMOVED***)
	if err := builderBytesEq(&b, 2, 23, 42, 5, 2, 123, 234); err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***

	var s, child String = b.BytesOrPanic(), nil
	var u, v, w, x, y uint8
	if !s.ReadUint8LengthPrefixed(&child) || !child.ReadUint8(&u) || !child.ReadUint8(&v) ||
		!s.ReadUint8(&w) || !s.ReadUint8LengthPrefixed(&child) || !child.ReadUint8(&x) || !child.ReadUint8(&y) ***REMOVED***
		t.Error("parsing failed")
	***REMOVED***
	if u != 23 || v != 42 || w != 5 || x != 123 || y != 234 ***REMOVED***
		t.Errorf("u, v, w, x, y = %d, %d, %d, %d, %d; want 23, 42, 5, 123, 234",
			u, v, w, x, y)
	***REMOVED***
	if len(s) != 0 ***REMOVED***
		t.Errorf("len(s) = %d, want 0", len(s))
	***REMOVED***
	if len(child) != 0 ***REMOVED***
		t.Errorf("len(child) = %d, want 0", len(child))
	***REMOVED***
***REMOVED***

func TestUint8LengthPrefixedNested(t *testing.T) ***REMOVED***
	var b Builder
	b.AddUint8LengthPrefixed(func(c *Builder) ***REMOVED***
		c.AddUint8(5)
		c.AddUint8LengthPrefixed(func(d *Builder) ***REMOVED***
			d.AddUint8(23)
			d.AddUint8(42)
		***REMOVED***)
		c.AddUint8(123)
	***REMOVED***)
	if err := builderBytesEq(&b, 5, 5, 2, 23, 42, 123); err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***

	var base, child1, child2 String = b.BytesOrPanic(), nil, nil
	var u, v, w, x uint8
	if !base.ReadUint8LengthPrefixed(&child1) ***REMOVED***
		t.Error("parsing base failed")
	***REMOVED***
	if !child1.ReadUint8(&u) || !child1.ReadUint8LengthPrefixed(&child2) || !child1.ReadUint8(&x) ***REMOVED***
		t.Error("parsing child1 failed")
	***REMOVED***
	if !child2.ReadUint8(&v) || !child2.ReadUint8(&w) ***REMOVED***
		t.Error("parsing child2 failed")
	***REMOVED***
	if u != 5 || v != 23 || w != 42 || x != 123 ***REMOVED***
		t.Errorf("u, v, w, x = %d, %d, %d, %d, want 5, 23, 42, 123",
			u, v, w, x)
	***REMOVED***
	if len(base) != 0 ***REMOVED***
		t.Errorf("len(base) = %d, want 0", len(base))
	***REMOVED***
	if len(child1) != 0 ***REMOVED***
		t.Errorf("len(child1) = %d, want 0", len(child1))
	***REMOVED***
	if len(base) != 0 ***REMOVED***
		t.Errorf("len(child2) = %d, want 0", len(child2))
	***REMOVED***
***REMOVED***

func TestPreallocatedBuffer(t *testing.T) ***REMOVED***
	var buf [5]byte
	b := NewBuilder(buf[0:0])
	b.AddUint8(1)
	b.AddUint8LengthPrefixed(func(c *Builder) ***REMOVED***
		c.AddUint8(3)
		c.AddUint8(4)
	***REMOVED***)
	b.AddUint16(1286) // Outgrow buf by one byte.
	want := []byte***REMOVED***1, 2, 3, 4, 0***REMOVED***
	if !bytes.Equal(buf[:], want) ***REMOVED***
		t.Errorf("buf = %v want %v", buf, want)
	***REMOVED***
	if err := builderBytesEq(b, 1, 2, 3, 4, 5, 6); err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***
***REMOVED***

func TestWriteWithPendingChild(t *testing.T) ***REMOVED***
	var b Builder
	b.AddUint8LengthPrefixed(func(c *Builder) ***REMOVED***
		c.AddUint8LengthPrefixed(func(d *Builder) ***REMOVED***
			defer func() ***REMOVED***
				if recover() == nil ***REMOVED***
					t.Errorf("recover() = nil, want error; c.AddUint8() did not panic")
				***REMOVED***
			***REMOVED***()
			c.AddUint8(2) // panics

			defer func() ***REMOVED***
				if recover() == nil ***REMOVED***
					t.Errorf("recover() = nil, want error; b.AddUint8() did not panic")
				***REMOVED***
			***REMOVED***()
			b.AddUint8(2) // panics
		***REMOVED***)

		defer func() ***REMOVED***
			if recover() == nil ***REMOVED***
				t.Errorf("recover() = nil, want error; b.AddUint8() did not panic")
			***REMOVED***
		***REMOVED***()
		b.AddUint8(2) // panics
	***REMOVED***)
***REMOVED***

// ASN.1

func TestASN1Int64(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		in   int64
		want []byte
	***REMOVED******REMOVED***
		***REMOVED***-0x800000, []byte***REMOVED***2, 3, 128, 0, 0***REMOVED******REMOVED***,
		***REMOVED***-256, []byte***REMOVED***2, 2, 255, 0***REMOVED******REMOVED***,
		***REMOVED***-129, []byte***REMOVED***2, 2, 255, 127***REMOVED******REMOVED***,
		***REMOVED***-128, []byte***REMOVED***2, 1, 128***REMOVED******REMOVED***,
		***REMOVED***-1, []byte***REMOVED***2, 1, 255***REMOVED******REMOVED***,
		***REMOVED***0, []byte***REMOVED***2, 1, 0***REMOVED******REMOVED***,
		***REMOVED***1, []byte***REMOVED***2, 1, 1***REMOVED******REMOVED***,
		***REMOVED***2, []byte***REMOVED***2, 1, 2***REMOVED******REMOVED***,
		***REMOVED***127, []byte***REMOVED***2, 1, 127***REMOVED******REMOVED***,
		***REMOVED***128, []byte***REMOVED***2, 2, 0, 128***REMOVED******REMOVED***,
		***REMOVED***256, []byte***REMOVED***2, 2, 1, 0***REMOVED******REMOVED***,
		***REMOVED***0x800000, []byte***REMOVED***2, 4, 0, 128, 0, 0***REMOVED******REMOVED***,
	***REMOVED***
	for i, tt := range tests ***REMOVED***
		var b Builder
		b.AddASN1Int64(tt.in)
		if err := builderBytesEq(&b, tt.want...); err != nil ***REMOVED***
			t.Errorf("%v, (i = %d; in = %v)", err, i, tt.in)
		***REMOVED***

		var n int64
		s := String(b.BytesOrPanic())
		ok := s.ReadASN1Integer(&n)
		if !ok || n != tt.in ***REMOVED***
			t.Errorf("s.ReadASN1Integer(&n) = %v, n = %d; want true, n = %d (i = %d)",
				ok, n, tt.in, i)
		***REMOVED***
		if len(s) != 0 ***REMOVED***
			t.Errorf("len(s) = %d, want 0", len(s))
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestASN1Uint64(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		in   uint64
		want []byte
	***REMOVED******REMOVED***
		***REMOVED***0, []byte***REMOVED***2, 1, 0***REMOVED******REMOVED***,
		***REMOVED***1, []byte***REMOVED***2, 1, 1***REMOVED******REMOVED***,
		***REMOVED***2, []byte***REMOVED***2, 1, 2***REMOVED******REMOVED***,
		***REMOVED***127, []byte***REMOVED***2, 1, 127***REMOVED******REMOVED***,
		***REMOVED***128, []byte***REMOVED***2, 2, 0, 128***REMOVED******REMOVED***,
		***REMOVED***256, []byte***REMOVED***2, 2, 1, 0***REMOVED******REMOVED***,
		***REMOVED***0x800000, []byte***REMOVED***2, 4, 0, 128, 0, 0***REMOVED******REMOVED***,
		***REMOVED***0x7fffffffffffffff, []byte***REMOVED***2, 8, 127, 255, 255, 255, 255, 255, 255, 255***REMOVED******REMOVED***,
		***REMOVED***0x8000000000000000, []byte***REMOVED***2, 9, 0, 128, 0, 0, 0, 0, 0, 0, 0***REMOVED******REMOVED***,
		***REMOVED***0xffffffffffffffff, []byte***REMOVED***2, 9, 0, 255, 255, 255, 255, 255, 255, 255, 255***REMOVED******REMOVED***,
	***REMOVED***
	for i, tt := range tests ***REMOVED***
		var b Builder
		b.AddASN1Uint64(tt.in)
		if err := builderBytesEq(&b, tt.want...); err != nil ***REMOVED***
			t.Errorf("%v, (i = %d; in = %v)", err, i, tt.in)
		***REMOVED***

		var n uint64
		s := String(b.BytesOrPanic())
		ok := s.ReadASN1Integer(&n)
		if !ok || n != tt.in ***REMOVED***
			t.Errorf("s.ReadASN1Integer(&n) = %v, n = %d; want true, n = %d (i = %d)",
				ok, n, tt.in, i)
		***REMOVED***
		if len(s) != 0 ***REMOVED***
			t.Errorf("len(s) = %d, want 0", len(s))
		***REMOVED***
	***REMOVED***
***REMOVED***
