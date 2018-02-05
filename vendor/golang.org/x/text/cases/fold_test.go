// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cases

import (
	"testing"

	"golang.org/x/text/internal/testtext"
)

var foldTestCases = []string***REMOVED***
	"βß\u13f8",        // "βssᏰ"
	"ab\u13fc\uab7aꭰ", // abᏴᎪᎠ
	"aﬃﬄaﬆ",           // affifflast
	"Iİiı\u0345",      // ii̇iıι
	"µµΜΜςσΣΣ",        // μμμμσσσσ
***REMOVED***

func TestFold(t *testing.T) ***REMOVED***
	for _, tc := range foldTestCases ***REMOVED***
		testEntry := func(name string, c Caser, m func(r rune) string) ***REMOVED***
			want := ""
			for _, r := range tc ***REMOVED***
				want += m(r)
			***REMOVED***
			if got := c.String(tc); got != want ***REMOVED***
				t.Errorf("%s(%s) = %+q; want %+q", name, tc, got, want)
			***REMOVED***
			dst := make([]byte, 256) // big enough to hold any result
			src := []byte(tc)
			v := testtext.AllocsPerRun(20, func() ***REMOVED***
				c.Transform(dst, src, true)
			***REMOVED***)
			if v > 0 ***REMOVED***
				t.Errorf("%s(%s): number of allocs was %f; want 0", name, tc, v)
			***REMOVED***
		***REMOVED***
		testEntry("FullFold", Fold(), func(r rune) string ***REMOVED***
			return runeFoldData(r).full
		***REMOVED***)
		// TODO:
		// testEntry("SimpleFold", Fold(Compact), func(r rune) string ***REMOVED***
		// 	return runeFoldData(r).simple
		// ***REMOVED***)
		// testEntry("SpecialFold", Fold(Turkic), func(r rune) string ***REMOVED***
		// 	return runeFoldData(r).special
		// ***REMOVED***)
	***REMOVED***
***REMOVED***
