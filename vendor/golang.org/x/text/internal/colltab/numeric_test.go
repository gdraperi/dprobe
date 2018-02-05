// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package colltab

import (
	"reflect"
	"strings"
	"testing"

	"golang.org/x/text/internal/testtext"
)

const (
	digSec  = defaultSecondary
	digTert = defaultTertiary
)

var tPlus3 = e(0, 50, digTert+3)

// numWeighter is a testWeighter used for testing numericWeighter.
var numWeighter = testWeighter***REMOVED***
	"0": p(100),
	"０": []Elem***REMOVED***e(100, digSec, digTert+1)***REMOVED***, // U+FF10 FULLWIDTH DIGIT ZERO
	"₀": []Elem***REMOVED***e(100, digSec, digTert+5)***REMOVED***, // U+2080 SUBSCRIPT ZERO

	"1": p(101),
	// Allow non-primary collation elements to be inserted.
	"١": append(p(101), tPlus3), // U+0661 ARABIC-INDIC DIGIT ONE
	// Allow varying tertiary weight if the number is Nd.
	"１": []Elem***REMOVED***e(101, digSec, digTert+1)***REMOVED***, // U+FF11 FULLWIDTH DIGIT ONE
	"2": p(102),
	// Allow non-primary collation elements to be inserted.
	"٢": append(p(102), tPlus3), // U+0662 ARABIC-INDIC DIGIT TWO
	// Varying tertiary weights should be ignored.
	"２": []Elem***REMOVED***e(102, digSec, digTert+3)***REMOVED***, // U+FF12 FULLWIDTH DIGIT TWO
	"3": p(103),
	"4": p(104),
	"5": p(105),
	"6": p(106),
	"7": p(107),
	// Weights must be strictly monotonically increasing, but do not need to be
	// consecutive.
	"8": p(118),
	"9": p(119),
	// Allow non-primary collation elements to be inserted.
	"٩": append(p(119), tPlus3), // U+0669 ARABIC-INDIC DIGIT NINE
	// Varying tertiary weights should be ignored.
	"９": []Elem***REMOVED***e(119, digSec, digTert+1)***REMOVED***, // U+FF19 FULLWIDTH DIGIT NINE
	"₉": []Elem***REMOVED***e(119, digSec, digTert+5)***REMOVED***, // U+2089 SUBSCRIPT NINE

	"a": p(5),
	"b": p(6),
	"c": p(8, 2),

	"klm": p(99),

	"nop": p(121),

	"x": p(200),
	"y": p(201),
***REMOVED***

func p(w ...int) (elems []Elem) ***REMOVED***
	for _, x := range w ***REMOVED***
		e, _ := MakeElem(x, digSec, digTert, 0)
		elems = append(elems, e)
	***REMOVED***
	return elems
***REMOVED***

func TestNumericAppendNext(t *testing.T) ***REMOVED***
	for _, tt := range []struct ***REMOVED***
		in string
		w  []Elem
	***REMOVED******REMOVED***
		***REMOVED***"a", p(5)***REMOVED***,
		***REMOVED***"klm", p(99)***REMOVED***,
		***REMOVED***"aa", p(5, 5)***REMOVED***,
		***REMOVED***"1", p(120, 1, 101)***REMOVED***,
		***REMOVED***"0", p(120, 0)***REMOVED***,
		***REMOVED***"01", p(120, 1, 101)***REMOVED***,
		***REMOVED***"0001", p(120, 1, 101)***REMOVED***,
		***REMOVED***"10", p(120, 2, 101, 100)***REMOVED***,
		***REMOVED***"99", p(120, 2, 119, 119)***REMOVED***,
		***REMOVED***"9999", p(120, 4, 119, 119, 119, 119)***REMOVED***,
		***REMOVED***"1a", p(120, 1, 101, 5)***REMOVED***,
		***REMOVED***"0b", p(120, 0, 6)***REMOVED***,
		***REMOVED***"01c", p(120, 1, 101, 8, 2)***REMOVED***,
		***REMOVED***"10x", p(120, 2, 101, 100, 200)***REMOVED***,
		***REMOVED***"99y", p(120, 2, 119, 119, 201)***REMOVED***,
		***REMOVED***"9999nop", p(120, 4, 119, 119, 119, 119, 121)***REMOVED***,

		// Allow follow-up collation elements if they have a zero non-primary.
		***REMOVED***"١٢٩", []Elem***REMOVED***e(120), e(3), e(101), tPlus3, e(102), tPlus3, e(119), tPlus3***REMOVED******REMOVED***,
		***REMOVED***
			"１２９",
			[]Elem***REMOVED***
				e(120), e(3),
				e(101, digSec, digTert+1),
				e(102, digSec, digTert+3),
				e(119, digSec, digTert+1),
			***REMOVED***,
		***REMOVED***,

		// Ensure AppendNext* adds to the given buffer.
		***REMOVED***"a10", p(5, 120, 2, 101, 100)***REMOVED***,
	***REMOVED*** ***REMOVED***
		nw := NewNumericWeighter(numWeighter)

		b := []byte(tt.in)
		got := []Elem(nil)
		for n, sz := 0, 0; n < len(b); ***REMOVED***
			got, sz = nw.AppendNext(got, b[n:])
			n += sz
		***REMOVED***
		if !reflect.DeepEqual(got, tt.w) ***REMOVED***
			t.Errorf("AppendNext(%q) =\n%v; want\n%v", tt.in, got, tt.w)
		***REMOVED***

		got = nil
		for n, sz := 0, 0; n < len(tt.in); ***REMOVED***
			got, sz = nw.AppendNextString(got, tt.in[n:])
			n += sz
		***REMOVED***
		if !reflect.DeepEqual(got, tt.w) ***REMOVED***
			t.Errorf("AppendNextString(%q) =\n%v; want\n%v", tt.in, got, tt.w)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestNumericOverflow(t *testing.T) ***REMOVED***
	manyDigits := strings.Repeat("9", maxDigits+1) + "a"

	nw := NewNumericWeighter(numWeighter)

	got, n := nw.AppendNextString(nil, manyDigits)

	if n != maxDigits ***REMOVED***
		t.Errorf("n: got %d; want %d", n, maxDigits)
	***REMOVED***

	if got[1].Primary() != maxDigits ***REMOVED***
		t.Errorf("primary(e[1]): got %d; want %d", n, maxDigits)
	***REMOVED***
***REMOVED***

func TestNumericWeighterAlloc(t *testing.T) ***REMOVED***
	buf := make([]Elem, 100)
	w := NewNumericWeighter(numWeighter)
	s := "1234567890a"

	nNormal := testtext.AllocsPerRun(3, func() ***REMOVED*** numWeighter.AppendNextString(buf, s) ***REMOVED***)
	nNumeric := testtext.AllocsPerRun(3, func() ***REMOVED*** w.AppendNextString(buf, s) ***REMOVED***)
	if n := nNumeric - nNormal; n > 0 ***REMOVED***
		t.Errorf("got %f; want 0", n)
	***REMOVED***
***REMOVED***
