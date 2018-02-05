// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin

package main

/*
#cgo LDFLAGS: -framework CoreFoundation
#include <CoreFoundation/CFBase.h>
#include <CoreFoundation/CoreFoundation.h>
*/
import "C"
import (
	"unsafe"
)

func init() ***REMOVED***
	AddFactory(CollatorFactory***REMOVED***"osx", newOSX16Collator,
		"OS X/Darwin collator, using native strings."***REMOVED***)
	AddFactory(CollatorFactory***REMOVED***"osx8", newOSX8Collator,
		"OS X/Darwin collator for UTF-8."***REMOVED***)
***REMOVED***

func osxUInt8P(s []byte) *C.UInt8 ***REMOVED***
	return (*C.UInt8)(unsafe.Pointer(&s[0]))
***REMOVED***

func osxCharP(s []uint16) *C.UniChar ***REMOVED***
	return (*C.UniChar)(unsafe.Pointer(&s[0]))
***REMOVED***

// osxCollator implements an Collator based on OS X's CoreFoundation.
type osxCollator struct ***REMOVED***
	loc C.CFLocaleRef
	opt C.CFStringCompareFlags
***REMOVED***

func (c *osxCollator) init(locale string) ***REMOVED***
	l := C.CFStringCreateWithBytes(
		C.kCFAllocatorDefault,
		osxUInt8P([]byte(locale)),
		C.CFIndex(len(locale)),
		C.kCFStringEncodingUTF8,
		C.Boolean(0),
	)
	c.loc = C.CFLocaleCreate(C.kCFAllocatorDefault, l)
***REMOVED***

func newOSX8Collator(locale string) (Collator, error) ***REMOVED***
	c := &osx8Collator***REMOVED******REMOVED***
	c.init(locale)
	return c, nil
***REMOVED***

func newOSX16Collator(locale string) (Collator, error) ***REMOVED***
	c := &osx16Collator***REMOVED******REMOVED***
	c.init(locale)
	return c, nil
***REMOVED***

func (c osxCollator) Key(s Input) []byte ***REMOVED***
	return nil // sort keys not supported by OS X CoreFoundation
***REMOVED***

type osx8Collator struct ***REMOVED***
	osxCollator
***REMOVED***

type osx16Collator struct ***REMOVED***
	osxCollator
***REMOVED***

func (c osx16Collator) Compare(a, b Input) int ***REMOVED***
	sa := C.CFStringCreateWithCharactersNoCopy(
		C.kCFAllocatorDefault,
		osxCharP(a.UTF16),
		C.CFIndex(len(a.UTF16)),
		C.kCFAllocatorDefault,
	)
	sb := C.CFStringCreateWithCharactersNoCopy(
		C.kCFAllocatorDefault,
		osxCharP(b.UTF16),
		C.CFIndex(len(b.UTF16)),
		C.kCFAllocatorDefault,
	)
	_range := C.CFRangeMake(0, C.CFStringGetLength(sa))
	return int(C.CFStringCompareWithOptionsAndLocale(sa, sb, _range, c.opt, c.loc))
***REMOVED***

func (c osx8Collator) Compare(a, b Input) int ***REMOVED***
	sa := C.CFStringCreateWithBytesNoCopy(
		C.kCFAllocatorDefault,
		osxUInt8P(a.UTF8),
		C.CFIndex(len(a.UTF8)),
		C.kCFStringEncodingUTF8,
		C.Boolean(0),
		C.kCFAllocatorDefault,
	)
	sb := C.CFStringCreateWithBytesNoCopy(
		C.kCFAllocatorDefault,
		osxUInt8P(b.UTF8),
		C.CFIndex(len(b.UTF8)),
		C.kCFStringEncodingUTF8,
		C.Boolean(0),
		C.kCFAllocatorDefault,
	)
	_range := C.CFRangeMake(0, C.CFStringGetLength(sa))
	return int(C.CFStringCompareWithOptionsAndLocale(sa, sb, _range, c.opt, c.loc))
***REMOVED***
