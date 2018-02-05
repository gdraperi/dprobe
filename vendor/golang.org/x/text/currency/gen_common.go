// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"time"

	"golang.org/x/text/language"
)

// This file contains code common to gen.go and the package code.

const (
	cashShift = 3
	roundMask = 0x7

	nonTenderBit = 0x8000
)

// currencyInfo contains information about a currency.
// bits 0..2: index into roundings for standard rounding
// bits 3..5: index into roundings for cash rounding
type currencyInfo byte

// roundingType defines the scale (number of fractional decimals) and increments
// in terms of units of size 10^-scale. For example, for scale == 2 and
// increment == 1, the currency is rounded to units of 0.01.
type roundingType struct ***REMOVED***
	scale, increment uint8
***REMOVED***

// roundings contains rounding data for currencies. This struct is
// created by hand as it is very unlikely to change much.
var roundings = [...]roundingType***REMOVED***
	***REMOVED***2, 1***REMOVED***, // default
	***REMOVED***0, 1***REMOVED***,
	***REMOVED***1, 1***REMOVED***,
	***REMOVED***3, 1***REMOVED***,
	***REMOVED***4, 1***REMOVED***,
	***REMOVED***2, 5***REMOVED***, // cash rounding alternative
	***REMOVED***2, 50***REMOVED***,
***REMOVED***

// regionToCode returns a 16-bit region code. Only two-letter codes are
// supported. (Three-letter codes are not needed.)
func regionToCode(r language.Region) uint16 ***REMOVED***
	if s := r.String(); len(s) == 2 ***REMOVED***
		return uint16(s[0])<<8 | uint16(s[1])
	***REMOVED***
	return 0
***REMOVED***

func toDate(t time.Time) uint32 ***REMOVED***
	y := t.Year()
	if y == 1 ***REMOVED***
		return 0
	***REMOVED***
	date := uint32(y) << 4
	date |= uint32(t.Month())
	date <<= 5
	date |= uint32(t.Day())
	return date
***REMOVED***

func fromDate(date uint32) time.Time ***REMOVED***
	return time.Date(int(date>>9), time.Month((date>>5)&0xf), int(date&0x1f), 0, 0, 0, 0, time.UTC)
***REMOVED***
