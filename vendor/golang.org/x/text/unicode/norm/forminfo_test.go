// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build test

package norm

import "testing"

func TestProperties(t *testing.T) ***REMOVED***
	var d runeData
	CK := [2]string***REMOVED***"C", "K"***REMOVED***
	for k, r := 1, rune(0); r < 0x2ffff; r++ ***REMOVED***
		if k < len(testData) && r == testData[k].r ***REMOVED***
			d = testData[k]
			k++
		***REMOVED***
		s := string(r)
		for j, p := range []Properties***REMOVED***NFC.PropertiesString(s), NFKC.PropertiesString(s)***REMOVED*** ***REMOVED***
			f := d.f[j]
			if p.CCC() != d.ccc ***REMOVED***
				t.Errorf("%U: ccc(%s): was %d; want %d %X", r, CK[j], p.CCC(), d.ccc, p.index)
			***REMOVED***
			if p.isYesC() != (f.qc == Yes) ***REMOVED***
				t.Errorf("%U: YesC(%s): was %v; want %v", r, CK[j], p.isYesC(), f.qc == Yes)
			***REMOVED***
			if p.combinesBackward() != (f.qc == Maybe) ***REMOVED***
				t.Errorf("%U: combines backwards(%s): was %v; want %v", r, CK[j], p.combinesBackward(), f.qc == Maybe)
			***REMOVED***
			if p.nLeadingNonStarters() != d.nLead ***REMOVED***
				t.Errorf("%U: nLead(%s): was %d; want %d %#v %#v", r, CK[j], p.nLeadingNonStarters(), d.nLead, p, d)
			***REMOVED***
			if p.nTrailingNonStarters() != d.nTrail ***REMOVED***
				t.Errorf("%U: nTrail(%s): was %d; want %d %#v %#v", r, CK[j], p.nTrailingNonStarters(), d.nTrail, p, d)
			***REMOVED***
			if p.combinesForward() != f.combinesForward ***REMOVED***
				t.Errorf("%U: combines forward(%s): was %v; want %v %#v", r, CK[j], p.combinesForward(), f.combinesForward, p)
			***REMOVED***
			// Skip Hangul as it is algorithmically computed.
			if r >= hangulBase && r < hangulEnd ***REMOVED***
				continue
			***REMOVED***
			if p.hasDecomposition() ***REMOVED***
				if has := f.decomposition != ""; !has ***REMOVED***
					t.Errorf("%U: hasDecomposition(%s): was %v; want %v", r, CK[j], p.hasDecomposition(), has)
				***REMOVED***
				if string(p.Decomposition()) != f.decomposition ***REMOVED***
					t.Errorf("%U: decomp(%s): was %+q; want %+q", r, CK[j], p.Decomposition(), f.decomposition)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
