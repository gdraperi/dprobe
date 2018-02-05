// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package idna

// This file implements the Punycode algorithm from RFC 3492.

import (
	"math"
	"strings"
	"unicode/utf8"
)

// These parameter values are specified in section 5.
//
// All computation is done with int32s, so that overflow behavior is identical
// regardless of whether int is 32-bit or 64-bit.
const (
	base        int32 = 36
	damp        int32 = 700
	initialBias int32 = 72
	initialN    int32 = 128
	skew        int32 = 38
	tmax        int32 = 26
	tmin        int32 = 1
)

func punyError(s string) error ***REMOVED*** return &labelError***REMOVED***s, "A3"***REMOVED*** ***REMOVED***

// decode decodes a string as specified in section 6.2.
func decode(encoded string) (string, error) ***REMOVED***
	if encoded == "" ***REMOVED***
		return "", nil
	***REMOVED***
	pos := 1 + strings.LastIndex(encoded, "-")
	if pos == 1 ***REMOVED***
		return "", punyError(encoded)
	***REMOVED***
	if pos == len(encoded) ***REMOVED***
		return encoded[:len(encoded)-1], nil
	***REMOVED***
	output := make([]rune, 0, len(encoded))
	if pos != 0 ***REMOVED***
		for _, r := range encoded[:pos-1] ***REMOVED***
			output = append(output, r)
		***REMOVED***
	***REMOVED***
	i, n, bias := int32(0), initialN, initialBias
	for pos < len(encoded) ***REMOVED***
		oldI, w := i, int32(1)
		for k := base; ; k += base ***REMOVED***
			if pos == len(encoded) ***REMOVED***
				return "", punyError(encoded)
			***REMOVED***
			digit, ok := decodeDigit(encoded[pos])
			if !ok ***REMOVED***
				return "", punyError(encoded)
			***REMOVED***
			pos++
			i += digit * w
			if i < 0 ***REMOVED***
				return "", punyError(encoded)
			***REMOVED***
			t := k - bias
			if t < tmin ***REMOVED***
				t = tmin
			***REMOVED*** else if t > tmax ***REMOVED***
				t = tmax
			***REMOVED***
			if digit < t ***REMOVED***
				break
			***REMOVED***
			w *= base - t
			if w >= math.MaxInt32/base ***REMOVED***
				return "", punyError(encoded)
			***REMOVED***
		***REMOVED***
		x := int32(len(output) + 1)
		bias = adapt(i-oldI, x, oldI == 0)
		n += i / x
		i %= x
		if n > utf8.MaxRune || len(output) >= 1024 ***REMOVED***
			return "", punyError(encoded)
		***REMOVED***
		output = append(output, 0)
		copy(output[i+1:], output[i:])
		output[i] = n
		i++
	***REMOVED***
	return string(output), nil
***REMOVED***

// encode encodes a string as specified in section 6.3 and prepends prefix to
// the result.
//
// The "while h < length(input)" line in the specification becomes "for
// remaining != 0" in the Go code, because len(s) in Go is in bytes, not runes.
func encode(prefix, s string) (string, error) ***REMOVED***
	output := make([]byte, len(prefix), len(prefix)+1+2*len(s))
	copy(output, prefix)
	delta, n, bias := int32(0), initialN, initialBias
	b, remaining := int32(0), int32(0)
	for _, r := range s ***REMOVED***
		if r < 0x80 ***REMOVED***
			b++
			output = append(output, byte(r))
		***REMOVED*** else ***REMOVED***
			remaining++
		***REMOVED***
	***REMOVED***
	h := b
	if b > 0 ***REMOVED***
		output = append(output, '-')
	***REMOVED***
	for remaining != 0 ***REMOVED***
		m := int32(0x7fffffff)
		for _, r := range s ***REMOVED***
			if m > r && r >= n ***REMOVED***
				m = r
			***REMOVED***
		***REMOVED***
		delta += (m - n) * (h + 1)
		if delta < 0 ***REMOVED***
			return "", punyError(s)
		***REMOVED***
		n = m
		for _, r := range s ***REMOVED***
			if r < n ***REMOVED***
				delta++
				if delta < 0 ***REMOVED***
					return "", punyError(s)
				***REMOVED***
				continue
			***REMOVED***
			if r > n ***REMOVED***
				continue
			***REMOVED***
			q := delta
			for k := base; ; k += base ***REMOVED***
				t := k - bias
				if t < tmin ***REMOVED***
					t = tmin
				***REMOVED*** else if t > tmax ***REMOVED***
					t = tmax
				***REMOVED***
				if q < t ***REMOVED***
					break
				***REMOVED***
				output = append(output, encodeDigit(t+(q-t)%(base-t)))
				q = (q - t) / (base - t)
			***REMOVED***
			output = append(output, encodeDigit(q))
			bias = adapt(delta, h+1, h == b)
			delta = 0
			h++
			remaining--
		***REMOVED***
		delta++
		n++
	***REMOVED***
	return string(output), nil
***REMOVED***

func decodeDigit(x byte) (digit int32, ok bool) ***REMOVED***
	switch ***REMOVED***
	case '0' <= x && x <= '9':
		return int32(x - ('0' - 26)), true
	case 'A' <= x && x <= 'Z':
		return int32(x - 'A'), true
	case 'a' <= x && x <= 'z':
		return int32(x - 'a'), true
	***REMOVED***
	return 0, false
***REMOVED***

func encodeDigit(digit int32) byte ***REMOVED***
	switch ***REMOVED***
	case 0 <= digit && digit < 26:
		return byte(digit + 'a')
	case 26 <= digit && digit < 36:
		return byte(digit + ('0' - 26))
	***REMOVED***
	panic("idna: internal error in punycode encoding")
***REMOVED***

// adapt is the bias adaptation function specified in section 6.1.
func adapt(delta, numPoints int32, firstTime bool) int32 ***REMOVED***
	if firstTime ***REMOVED***
		delta /= damp
	***REMOVED*** else ***REMOVED***
		delta /= 2
	***REMOVED***
	delta += delta / numPoints
	k := int32(0)
	for delta > ((base-tmin)*tmax)/2 ***REMOVED***
		delta /= base - tmin
		k += base
	***REMOVED***
	return k + (base-tmin+1)*delta/(delta+skew)
***REMOVED***
