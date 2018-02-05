// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build icu

package main

/*
#cgo LDFLAGS: -licui18n -licuuc
#include <stdlib.h>
#include <unicode/ucol.h>
#include <unicode/uiter.h>
#include <unicode/utypes.h>
*/
import "C"
import (
	"fmt"
	"log"
	"unicode/utf16"
	"unicode/utf8"
	"unsafe"
)

func init() ***REMOVED***
	AddFactory(CollatorFactory***REMOVED***"icu", newUTF16,
		"Main ICU collator, using native strings."***REMOVED***)
	AddFactory(CollatorFactory***REMOVED***"icu8", newUTF8iter,
		"ICU collator using ICU iterators to process UTF8."***REMOVED***)
	AddFactory(CollatorFactory***REMOVED***"icu16", newUTF8conv,
		"ICU collation by first converting UTF8 to UTF16."***REMOVED***)
***REMOVED***

func icuCharP(s []byte) *C.char ***REMOVED***
	return (*C.char)(unsafe.Pointer(&s[0]))
***REMOVED***

func icuUInt8P(s []byte) *C.uint8_t ***REMOVED***
	return (*C.uint8_t)(unsafe.Pointer(&s[0]))
***REMOVED***

func icuUCharP(s []uint16) *C.UChar ***REMOVED***
	return (*C.UChar)(unsafe.Pointer(&s[0]))
***REMOVED***
func icuULen(s []uint16) C.int32_t ***REMOVED***
	return C.int32_t(len(s))
***REMOVED***
func icuSLen(s []byte) C.int32_t ***REMOVED***
	return C.int32_t(len(s))
***REMOVED***

// icuCollator implements a Collator based on ICU.
type icuCollator struct ***REMOVED***
	loc    *C.char
	col    *C.UCollator
	keyBuf []byte
***REMOVED***

const growBufSize = 10 * 1024 * 1024

func (c *icuCollator) init(locale string) error ***REMOVED***
	err := C.UErrorCode(0)
	c.loc = C.CString(locale)
	c.col = C.ucol_open(c.loc, &err)
	if err > 0 ***REMOVED***
		return fmt.Errorf("failed opening collator for %q", locale)
	***REMOVED*** else if err < 0 ***REMOVED***
		loc := C.ucol_getLocaleByType(c.col, 0, &err)
		fmt, ok := map[int]string***REMOVED***
			-127: "warning: using default collator: %s",
			-128: "warning: using fallback collator: %s",
		***REMOVED***[int(err)]
		if ok ***REMOVED***
			log.Printf(fmt, C.GoString(loc))
		***REMOVED***
	***REMOVED***
	c.keyBuf = make([]byte, 0, growBufSize)
	return nil
***REMOVED***

func (c *icuCollator) buf() (*C.uint8_t, C.int32_t) ***REMOVED***
	if len(c.keyBuf) == cap(c.keyBuf) ***REMOVED***
		c.keyBuf = make([]byte, 0, growBufSize)
	***REMOVED***
	b := c.keyBuf[len(c.keyBuf):cap(c.keyBuf)]
	return icuUInt8P(b), icuSLen(b)
***REMOVED***

func (c *icuCollator) extendBuf(n C.int32_t) []byte ***REMOVED***
	end := len(c.keyBuf) + int(n)
	if end > cap(c.keyBuf) ***REMOVED***
		if len(c.keyBuf) == 0 ***REMOVED***
			log.Fatalf("icuCollator: max string size exceeded: %v > %v", n, growBufSize)
		***REMOVED***
		c.keyBuf = make([]byte, 0, growBufSize)
		return nil
	***REMOVED***
	b := c.keyBuf[len(c.keyBuf):end]
	c.keyBuf = c.keyBuf[:end]
	return b
***REMOVED***

func (c *icuCollator) Close() error ***REMOVED***
	C.ucol_close(c.col)
	C.free(unsafe.Pointer(c.loc))
	return nil
***REMOVED***

// icuUTF16 implements the Collator interface.
type icuUTF16 struct ***REMOVED***
	icuCollator
***REMOVED***

func newUTF16(locale string) (Collator, error) ***REMOVED***
	c := &icuUTF16***REMOVED******REMOVED***
	return c, c.init(locale)
***REMOVED***

func (c *icuUTF16) Compare(a, b Input) int ***REMOVED***
	return int(C.ucol_strcoll(c.col, icuUCharP(a.UTF16), icuULen(a.UTF16), icuUCharP(b.UTF16), icuULen(b.UTF16)))
***REMOVED***

func (c *icuUTF16) Key(s Input) []byte ***REMOVED***
	bp, bn := c.buf()
	n := C.ucol_getSortKey(c.col, icuUCharP(s.UTF16), icuULen(s.UTF16), bp, bn)
	if b := c.extendBuf(n); b != nil ***REMOVED***
		return b
	***REMOVED***
	return c.Key(s)
***REMOVED***

// icuUTF8iter implements the Collator interface
// This implementation wraps the UTF8 string in an iterator
// which is passed to the collator.
type icuUTF8iter struct ***REMOVED***
	icuCollator
	a, b C.UCharIterator
***REMOVED***

func newUTF8iter(locale string) (Collator, error) ***REMOVED***
	c := &icuUTF8iter***REMOVED******REMOVED***
	return c, c.init(locale)
***REMOVED***

func (c *icuUTF8iter) Compare(a, b Input) int ***REMOVED***
	err := C.UErrorCode(0)
	C.uiter_setUTF8(&c.a, icuCharP(a.UTF8), icuSLen(a.UTF8))
	C.uiter_setUTF8(&c.b, icuCharP(b.UTF8), icuSLen(b.UTF8))
	return int(C.ucol_strcollIter(c.col, &c.a, &c.b, &err))
***REMOVED***

func (c *icuUTF8iter) Key(s Input) []byte ***REMOVED***
	err := C.UErrorCode(0)
	state := [2]C.uint32_t***REMOVED******REMOVED***
	C.uiter_setUTF8(&c.a, icuCharP(s.UTF8), icuSLen(s.UTF8))
	bp, bn := c.buf()
	n := C.ucol_nextSortKeyPart(c.col, &c.a, &(state[0]), bp, bn, &err)
	if n >= bn ***REMOVED***
		// Force failure.
		if c.extendBuf(n+1) != nil ***REMOVED***
			log.Fatal("expected extension to fail")
		***REMOVED***
		return c.Key(s)
	***REMOVED***
	return c.extendBuf(n)
***REMOVED***

// icuUTF8conv implements the Collator interface.
// This implementation first converts the give UTF8 string
// to UTF16 and then calls the main ICU collation function.
type icuUTF8conv struct ***REMOVED***
	icuCollator
***REMOVED***

func newUTF8conv(locale string) (Collator, error) ***REMOVED***
	c := &icuUTF8conv***REMOVED******REMOVED***
	return c, c.init(locale)
***REMOVED***

func (c *icuUTF8conv) Compare(sa, sb Input) int ***REMOVED***
	a := encodeUTF16(sa.UTF8)
	b := encodeUTF16(sb.UTF8)
	return int(C.ucol_strcoll(c.col, icuUCharP(a), icuULen(a), icuUCharP(b), icuULen(b)))
***REMOVED***

func (c *icuUTF8conv) Key(s Input) []byte ***REMOVED***
	a := encodeUTF16(s.UTF8)
	bp, bn := c.buf()
	n := C.ucol_getSortKey(c.col, icuUCharP(a), icuULen(a), bp, bn)
	if b := c.extendBuf(n); b != nil ***REMOVED***
		return b
	***REMOVED***
	return c.Key(s)
***REMOVED***

func encodeUTF16(b []byte) []uint16 ***REMOVED***
	a := []uint16***REMOVED******REMOVED***
	for len(b) > 0 ***REMOVED***
		r, sz := utf8.DecodeRune(b)
		b = b[sz:]
		r1, r2 := utf16.EncodeRune(r)
		if r1 != 0xFFFD ***REMOVED***
			a = append(a, uint16(r1), uint16(r2))
		***REMOVED*** else ***REMOVED***
			a = append(a, uint16(r))
		***REMOVED***
	***REMOVED***
	return a
***REMOVED***
