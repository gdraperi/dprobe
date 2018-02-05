// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package norm

import (
	"fmt"
	"testing"

	"golang.org/x/text/transform"
)

func TestTransform(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		f       Form
		in, out string
		eof     bool
		dstSize int
		err     error
	***REMOVED******REMOVED***
		***REMOVED***NFC, "ab", "ab", true, 2, nil***REMOVED***,
		***REMOVED***NFC, "qx", "qx", true, 2, nil***REMOVED***,
		***REMOVED***NFD, "qx", "qx", true, 2, nil***REMOVED***,
		***REMOVED***NFC, "", "", true, 1, nil***REMOVED***,
		***REMOVED***NFD, "", "", true, 1, nil***REMOVED***,
		***REMOVED***NFC, "", "", false, 1, nil***REMOVED***,
		***REMOVED***NFD, "", "", false, 1, nil***REMOVED***,

		// Normalized segment does not fit in destination.
		***REMOVED***NFD, "ö", "", true, 1, transform.ErrShortDst***REMOVED***,
		***REMOVED***NFD, "ö", "", true, 2, transform.ErrShortDst***REMOVED***,

		// As an artifact of the algorithm, only full segments are written.
		// This is not strictly required, and some bytes could be written.
		// In practice, for Transform to not block, the destination buffer
		// should be at least MaxSegmentSize to work anyway and these edge
		// conditions will be relatively rare.
		***REMOVED***NFC, "ab", "", true, 1, transform.ErrShortDst***REMOVED***,
		// This is even true for inert runes.
		***REMOVED***NFC, "qx", "", true, 1, transform.ErrShortDst***REMOVED***,
		***REMOVED***NFC, "a\u0300abc", "\u00e0a", true, 4, transform.ErrShortDst***REMOVED***,

		// We cannot write a segment if successive runes could still change the result.
		***REMOVED***NFD, "ö", "", false, 3, transform.ErrShortSrc***REMOVED***,
		***REMOVED***NFC, "a\u0300", "", false, 4, transform.ErrShortSrc***REMOVED***,
		***REMOVED***NFD, "a\u0300", "", false, 4, transform.ErrShortSrc***REMOVED***,
		***REMOVED***NFC, "ö", "", false, 3, transform.ErrShortSrc***REMOVED***,

		***REMOVED***NFC, "a\u0300", "", true, 1, transform.ErrShortDst***REMOVED***,
		// Theoretically could fit, but won't due to simplified checks.
		***REMOVED***NFC, "a\u0300", "", true, 2, transform.ErrShortDst***REMOVED***,
		***REMOVED***NFC, "a\u0300", "", true, 3, transform.ErrShortDst***REMOVED***,
		***REMOVED***NFC, "a\u0300", "\u00e0", true, 4, nil***REMOVED***,

		***REMOVED***NFD, "öa\u0300", "o\u0308", false, 8, transform.ErrShortSrc***REMOVED***,
		***REMOVED***NFD, "öa\u0300ö", "o\u0308a\u0300", true, 8, transform.ErrShortDst***REMOVED***,
		***REMOVED***NFD, "öa\u0300ö", "o\u0308a\u0300", false, 12, transform.ErrShortSrc***REMOVED***,

		// Illegal input is copied verbatim.
		***REMOVED***NFD, "\xbd\xb2=\xbc ", "\xbd\xb2=\xbc ", true, 8, nil***REMOVED***,
	***REMOVED***
	b := make([]byte, 100)
	for i, tt := range tests ***REMOVED***
		nDst, _, err := tt.f.Transform(b[:tt.dstSize], []byte(tt.in), tt.eof)
		out := string(b[:nDst])
		if out != tt.out || err != tt.err ***REMOVED***
			t.Errorf("%d: was %+q (%v); want %+q (%v)", i, out, err, tt.out, tt.err)
		***REMOVED***
		if want := tt.f.String(tt.in)[:nDst]; want != out ***REMOVED***
			t.Errorf("%d: incorrect normalization: was %+q; want %+q", i, out, want)
		***REMOVED***
	***REMOVED***
***REMOVED***

var transBufSizes = []int***REMOVED***
	MaxTransformChunkSize,
	3 * MaxTransformChunkSize / 2,
	2 * MaxTransformChunkSize,
	3 * MaxTransformChunkSize,
	100 * MaxTransformChunkSize,
***REMOVED***

func doTransNorm(f Form, buf []byte, b []byte) []byte ***REMOVED***
	acc := []byte***REMOVED******REMOVED***
	for p := 0; p < len(b); ***REMOVED***
		nd, ns, _ := f.Transform(buf[:], b[p:], true)
		p += ns
		acc = append(acc, buf[:nd]...)
	***REMOVED***
	return acc
***REMOVED***

func TestTransformNorm(t *testing.T) ***REMOVED***
	for _, sz := range transBufSizes ***REMOVED***
		buf := make([]byte, sz)
		runNormTests(t, fmt.Sprintf("Transform:%d", sz), func(f Form, out []byte, s string) []byte ***REMOVED***
			return doTransNorm(f, buf, append(out, s...))
		***REMOVED***)
	***REMOVED***
***REMOVED***
