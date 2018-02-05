// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package norm

import (
	"strings"
	"testing"
)

func doIterNorm(f Form, s string) []byte ***REMOVED***
	acc := []byte***REMOVED******REMOVED***
	i := Iter***REMOVED******REMOVED***
	i.InitString(f, s)
	for !i.Done() ***REMOVED***
		acc = append(acc, i.Next()...)
	***REMOVED***
	return acc
***REMOVED***

func TestIterNext(t *testing.T) ***REMOVED***
	runNormTests(t, "IterNext", func(f Form, out []byte, s string) []byte ***REMOVED***
		return doIterNorm(f, string(append(out, s...)))
	***REMOVED***)
***REMOVED***

type SegmentTest struct ***REMOVED***
	in  string
	out []string
***REMOVED***

var segmentTests = []SegmentTest***REMOVED***
	***REMOVED***"\u1E0A\u0323a", []string***REMOVED***"\x44\u0323\u0307", "a", ""***REMOVED******REMOVED***,
	***REMOVED***rep('a', segSize), append(strings.Split(rep('a', segSize), ""), "")***REMOVED***,
	***REMOVED***rep('a', segSize+2), append(strings.Split(rep('a', segSize+2), ""), "")***REMOVED***,
	***REMOVED***rep('a', segSize) + "\u0300aa",
		append(strings.Split(rep('a', segSize-1), ""), "a\u0300", "a", "a", "")***REMOVED***,

	// U+0f73 is NOT treated as a starter as it is a modifier
	***REMOVED***"a" + grave(29) + "\u0f73", []string***REMOVED***"a" + grave(29), cgj + "\u0f73"***REMOVED******REMOVED***,
	***REMOVED***"a\u0f73", []string***REMOVED***"a\u0f73"***REMOVED******REMOVED***,

	// U+ff9e is treated as a non-starter.
	// TODO: should we? Note that this will only affect iteration, as whether
	// or not we do so does not affect the normalization output and will either
	// way result in consistent iteration output.
	***REMOVED***"a" + grave(30) + "\uff9e", []string***REMOVED***"a" + grave(30), cgj + "\uff9e"***REMOVED******REMOVED***,
	***REMOVED***"a\uff9e", []string***REMOVED***"a\uff9e"***REMOVED******REMOVED***,
***REMOVED***

var segmentTestsK = []SegmentTest***REMOVED***
	***REMOVED***"\u3332", []string***REMOVED***"\u30D5", "\u30A1", "\u30E9", "\u30C3", "\u30C8\u3099", ""***REMOVED******REMOVED***,
	// last segment of multi-segment decomposition needs normalization
	***REMOVED***"\u3332\u093C", []string***REMOVED***"\u30D5", "\u30A1", "\u30E9", "\u30C3", "\u30C8\u093C\u3099", ""***REMOVED******REMOVED***,
	***REMOVED***"\u320E", []string***REMOVED***"\x28", "\uAC00", "\x29"***REMOVED******REMOVED***,

	// last segment should be copied to start of buffer.
	***REMOVED***"\ufdfa", []string***REMOVED***"\u0635", "\u0644", "\u0649", " ", "\u0627", "\u0644", "\u0644", "\u0647", " ", "\u0639", "\u0644", "\u064a", "\u0647", " ", "\u0648", "\u0633", "\u0644", "\u0645", ""***REMOVED******REMOVED***,
	***REMOVED***"\ufdfa" + grave(30), []string***REMOVED***"\u0635", "\u0644", "\u0649", " ", "\u0627", "\u0644", "\u0644", "\u0647", " ", "\u0639", "\u0644", "\u064a", "\u0647", " ", "\u0648", "\u0633", "\u0644", "\u0645" + grave(30), ""***REMOVED******REMOVED***,
	***REMOVED***"\uFDFA" + grave(64), []string***REMOVED***"\u0635", "\u0644", "\u0649", " ", "\u0627", "\u0644", "\u0644", "\u0647", " ", "\u0639", "\u0644", "\u064a", "\u0647", " ", "\u0648", "\u0633", "\u0644", "\u0645" + grave(30), cgj + grave(30), cgj + grave(4), ""***REMOVED******REMOVED***,

	// Hangul and Jamo are grouped together.
	***REMOVED***"\uAC00", []string***REMOVED***"\u1100\u1161", ""***REMOVED******REMOVED***,
	***REMOVED***"\uAC01", []string***REMOVED***"\u1100\u1161\u11A8", ""***REMOVED******REMOVED***,
	***REMOVED***"\u1100\u1161", []string***REMOVED***"\u1100\u1161", ""***REMOVED******REMOVED***,
***REMOVED***

// Note that, by design, segmentation is equal for composing and decomposing forms.
func TestIterSegmentation(t *testing.T) ***REMOVED***
	segmentTest(t, "SegmentTestD", NFD, segmentTests)
	segmentTest(t, "SegmentTestC", NFC, segmentTests)
	segmentTest(t, "SegmentTestKD", NFKD, segmentTestsK)
	segmentTest(t, "SegmentTestKC", NFKC, segmentTestsK)
***REMOVED***

func segmentTest(t *testing.T, name string, f Form, tests []SegmentTest) ***REMOVED***
	iter := Iter***REMOVED******REMOVED***
	for i, tt := range tests ***REMOVED***
		iter.InitString(f, tt.in)
		for j, seg := range tt.out ***REMOVED***
			if seg == "" ***REMOVED***
				if !iter.Done() ***REMOVED***
					res := string(iter.Next())
					t.Errorf(`%s:%d:%d: expected Done()==true, found segment %+q`, name, i, j, res)
				***REMOVED***
				continue
			***REMOVED***
			if iter.Done() ***REMOVED***
				t.Errorf("%s:%d:%d: Done()==true, want false", name, i, j)
			***REMOVED***
			seg = f.String(seg)
			if res := string(iter.Next()); res != seg ***REMOVED***
				t.Errorf(`%s:%d:%d" segment was %+q (%d); want %+q (%d)`, name, i, j, pc(res), len(res), pc(seg), len(seg))
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
