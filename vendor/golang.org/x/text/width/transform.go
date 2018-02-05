// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package width

import (
	"unicode/utf8"

	"golang.org/x/text/transform"
)

type foldTransform struct ***REMOVED***
	transform.NopResetter
***REMOVED***

func (foldTransform) Span(src []byte, atEOF bool) (n int, err error) ***REMOVED***
	for n < len(src) ***REMOVED***
		if src[n] < utf8.RuneSelf ***REMOVED***
			// ASCII fast path.
			for n++; n < len(src) && src[n] < utf8.RuneSelf; n++ ***REMOVED***
			***REMOVED***
			continue
		***REMOVED***
		v, size := trie.lookup(src[n:])
		if size == 0 ***REMOVED*** // incomplete UTF-8 encoding
			if !atEOF ***REMOVED***
				err = transform.ErrShortSrc
			***REMOVED*** else ***REMOVED***
				n = len(src)
			***REMOVED***
			break
		***REMOVED***
		if elem(v)&tagNeedsFold != 0 ***REMOVED***
			err = transform.ErrEndOfSpan
			break
		***REMOVED***
		n += size
	***REMOVED***
	return n, err
***REMOVED***

func (foldTransform) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) ***REMOVED***
	for nSrc < len(src) ***REMOVED***
		if src[nSrc] < utf8.RuneSelf ***REMOVED***
			// ASCII fast path.
			start, end := nSrc, len(src)
			if d := len(dst) - nDst; d < end-start ***REMOVED***
				end = nSrc + d
			***REMOVED***
			for nSrc++; nSrc < end && src[nSrc] < utf8.RuneSelf; nSrc++ ***REMOVED***
			***REMOVED***
			n := copy(dst[nDst:], src[start:nSrc])
			if nDst += n; nDst == len(dst) ***REMOVED***
				nSrc = start + n
				if nSrc == len(src) ***REMOVED***
					return nDst, nSrc, nil
				***REMOVED***
				if src[nSrc] < utf8.RuneSelf ***REMOVED***
					return nDst, nSrc, transform.ErrShortDst
				***REMOVED***
			***REMOVED***
			continue
		***REMOVED***
		v, size := trie.lookup(src[nSrc:])
		if size == 0 ***REMOVED*** // incomplete UTF-8 encoding
			if !atEOF ***REMOVED***
				return nDst, nSrc, transform.ErrShortSrc
			***REMOVED***
			size = 1 // gobble 1 byte
		***REMOVED***
		if elem(v)&tagNeedsFold == 0 ***REMOVED***
			if size != copy(dst[nDst:], src[nSrc:nSrc+size]) ***REMOVED***
				return nDst, nSrc, transform.ErrShortDst
			***REMOVED***
			nDst += size
		***REMOVED*** else ***REMOVED***
			data := inverseData[byte(v)]
			if len(dst)-nDst < int(data[0]) ***REMOVED***
				return nDst, nSrc, transform.ErrShortDst
			***REMOVED***
			i := 1
			for end := int(data[0]); i < end; i++ ***REMOVED***
				dst[nDst] = data[i]
				nDst++
			***REMOVED***
			dst[nDst] = data[i] ^ src[nSrc+size-1]
			nDst++
		***REMOVED***
		nSrc += size
	***REMOVED***
	return nDst, nSrc, nil
***REMOVED***

type narrowTransform struct ***REMOVED***
	transform.NopResetter
***REMOVED***

func (narrowTransform) Span(src []byte, atEOF bool) (n int, err error) ***REMOVED***
	for n < len(src) ***REMOVED***
		if src[n] < utf8.RuneSelf ***REMOVED***
			// ASCII fast path.
			for n++; n < len(src) && src[n] < utf8.RuneSelf; n++ ***REMOVED***
			***REMOVED***
			continue
		***REMOVED***
		v, size := trie.lookup(src[n:])
		if size == 0 ***REMOVED*** // incomplete UTF-8 encoding
			if !atEOF ***REMOVED***
				err = transform.ErrShortSrc
			***REMOVED*** else ***REMOVED***
				n = len(src)
			***REMOVED***
			break
		***REMOVED***
		if k := elem(v).kind(); byte(v) == 0 || k != EastAsianFullwidth && k != EastAsianWide && k != EastAsianAmbiguous ***REMOVED***
		***REMOVED*** else ***REMOVED***
			err = transform.ErrEndOfSpan
			break
		***REMOVED***
		n += size
	***REMOVED***
	return n, err
***REMOVED***

func (narrowTransform) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) ***REMOVED***
	for nSrc < len(src) ***REMOVED***
		if src[nSrc] < utf8.RuneSelf ***REMOVED***
			// ASCII fast path.
			start, end := nSrc, len(src)
			if d := len(dst) - nDst; d < end-start ***REMOVED***
				end = nSrc + d
			***REMOVED***
			for nSrc++; nSrc < end && src[nSrc] < utf8.RuneSelf; nSrc++ ***REMOVED***
			***REMOVED***
			n := copy(dst[nDst:], src[start:nSrc])
			if nDst += n; nDst == len(dst) ***REMOVED***
				nSrc = start + n
				if nSrc == len(src) ***REMOVED***
					return nDst, nSrc, nil
				***REMOVED***
				if src[nSrc] < utf8.RuneSelf ***REMOVED***
					return nDst, nSrc, transform.ErrShortDst
				***REMOVED***
			***REMOVED***
			continue
		***REMOVED***
		v, size := trie.lookup(src[nSrc:])
		if size == 0 ***REMOVED*** // incomplete UTF-8 encoding
			if !atEOF ***REMOVED***
				return nDst, nSrc, transform.ErrShortSrc
			***REMOVED***
			size = 1 // gobble 1 byte
		***REMOVED***
		if k := elem(v).kind(); byte(v) == 0 || k != EastAsianFullwidth && k != EastAsianWide && k != EastAsianAmbiguous ***REMOVED***
			if size != copy(dst[nDst:], src[nSrc:nSrc+size]) ***REMOVED***
				return nDst, nSrc, transform.ErrShortDst
			***REMOVED***
			nDst += size
		***REMOVED*** else ***REMOVED***
			data := inverseData[byte(v)]
			if len(dst)-nDst < int(data[0]) ***REMOVED***
				return nDst, nSrc, transform.ErrShortDst
			***REMOVED***
			i := 1
			for end := int(data[0]); i < end; i++ ***REMOVED***
				dst[nDst] = data[i]
				nDst++
			***REMOVED***
			dst[nDst] = data[i] ^ src[nSrc+size-1]
			nDst++
		***REMOVED***
		nSrc += size
	***REMOVED***
	return nDst, nSrc, nil
***REMOVED***

type wideTransform struct ***REMOVED***
	transform.NopResetter
***REMOVED***

func (wideTransform) Span(src []byte, atEOF bool) (n int, err error) ***REMOVED***
	for n < len(src) ***REMOVED***
		// TODO: Consider ASCII fast path. Special-casing ASCII handling can
		// reduce the ns/op of BenchmarkWideASCII by about 30%. This is probably
		// not enough to warrant the extra code and complexity.
		v, size := trie.lookup(src[n:])
		if size == 0 ***REMOVED*** // incomplete UTF-8 encoding
			if !atEOF ***REMOVED***
				err = transform.ErrShortSrc
			***REMOVED*** else ***REMOVED***
				n = len(src)
			***REMOVED***
			break
		***REMOVED***
		if k := elem(v).kind(); byte(v) == 0 || k != EastAsianHalfwidth && k != EastAsianNarrow ***REMOVED***
		***REMOVED*** else ***REMOVED***
			err = transform.ErrEndOfSpan
			break
		***REMOVED***
		n += size
	***REMOVED***
	return n, err
***REMOVED***

func (wideTransform) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) ***REMOVED***
	for nSrc < len(src) ***REMOVED***
		// TODO: Consider ASCII fast path. Special-casing ASCII handling can
		// reduce the ns/op of BenchmarkWideASCII by about 30%. This is probably
		// not enough to warrant the extra code and complexity.
		v, size := trie.lookup(src[nSrc:])
		if size == 0 ***REMOVED*** // incomplete UTF-8 encoding
			if !atEOF ***REMOVED***
				return nDst, nSrc, transform.ErrShortSrc
			***REMOVED***
			size = 1 // gobble 1 byte
		***REMOVED***
		if k := elem(v).kind(); byte(v) == 0 || k != EastAsianHalfwidth && k != EastAsianNarrow ***REMOVED***
			if size != copy(dst[nDst:], src[nSrc:nSrc+size]) ***REMOVED***
				return nDst, nSrc, transform.ErrShortDst
			***REMOVED***
			nDst += size
		***REMOVED*** else ***REMOVED***
			data := inverseData[byte(v)]
			if len(dst)-nDst < int(data[0]) ***REMOVED***
				return nDst, nSrc, transform.ErrShortDst
			***REMOVED***
			i := 1
			for end := int(data[0]); i < end; i++ ***REMOVED***
				dst[nDst] = data[i]
				nDst++
			***REMOVED***
			dst[nDst] = data[i] ^ src[nSrc+size-1]
			nDst++
		***REMOVED***
		nSrc += size
	***REMOVED***
	return nDst, nSrc, nil
***REMOVED***
