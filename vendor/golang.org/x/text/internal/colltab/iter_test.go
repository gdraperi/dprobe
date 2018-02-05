// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package colltab

import (
	"testing"
)

func TestDoNorm(t *testing.T) ***REMOVED***
	const div = -1 // The insertion point of the next block.
	tests := []struct ***REMOVED***
		in, out []int
	***REMOVED******REMOVED******REMOVED***
		in:  []int***REMOVED***4, div, 3***REMOVED***,
		out: []int***REMOVED***3, 4***REMOVED***,
	***REMOVED***, ***REMOVED***
		in:  []int***REMOVED***4, div, 3, 3, 3***REMOVED***,
		out: []int***REMOVED***3, 3, 3, 4***REMOVED***,
	***REMOVED***, ***REMOVED***
		in:  []int***REMOVED***0, 4, div, 3***REMOVED***,
		out: []int***REMOVED***0, 3, 4***REMOVED***,
	***REMOVED***, ***REMOVED***
		in:  []int***REMOVED***0, 0, 4, 5, div, 3, 3***REMOVED***,
		out: []int***REMOVED***0, 0, 3, 3, 4, 5***REMOVED***,
	***REMOVED***, ***REMOVED***
		in:  []int***REMOVED***0, 0, 1, 4, 5, div, 3, 3***REMOVED***,
		out: []int***REMOVED***0, 0, 1, 3, 3, 4, 5***REMOVED***,
	***REMOVED***, ***REMOVED***
		in:  []int***REMOVED***0, 0, 1, 4, 5, div, 4, 4***REMOVED***,
		out: []int***REMOVED***0, 0, 1, 4, 4, 4, 5***REMOVED***,
	***REMOVED***,
	***REMOVED***
	for j, tt := range tests ***REMOVED***
		i := Iter***REMOVED******REMOVED***
		var w, p int
		for k, cc := range tt.in ***REMOVED***

			if cc == div ***REMOVED***
				w = 100
				p = k
				continue
			***REMOVED***
			i.Elems = append(i.Elems, makeCE([]int***REMOVED***w, defaultSecondary, 2, cc***REMOVED***))
		***REMOVED***
		i.doNorm(p, i.Elems[p].CCC())
		if len(i.Elems) != len(tt.out) ***REMOVED***
			t.Errorf("%d: length was %d; want %d", j, len(i.Elems), len(tt.out))
		***REMOVED***
		prevCCC := uint8(0)
		for k, ce := range i.Elems ***REMOVED***
			if int(ce.CCC()) != tt.out[k] ***REMOVED***
				t.Errorf("%d:%d: unexpected CCC. Was %d; want %d", j, k, ce.CCC(), tt.out[k])
			***REMOVED***
			if k > 0 && ce.CCC() == prevCCC && i.Elems[k-1].Primary() > ce.Primary() ***REMOVED***
				t.Errorf("%d:%d: normalization crossed across CCC boundary.", j, k)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Combining rune overflow is tested in search/pattern_test.go.
***REMOVED***
