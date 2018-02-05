// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package collate

import (
	"bytes"
	"testing"

	"golang.org/x/text/internal/colltab"
	"golang.org/x/text/language"
)

type weightsTest struct ***REMOVED***
	opt     opts
	in, out ColElems
***REMOVED***

type opts struct ***REMOVED***
	lev int
	alt alternateHandling
	top int

	backwards bool
	caseLevel bool
***REMOVED***

// ignore returns an initialized boolean array based on the given Level.
// A negative value means using the default setting of quaternary.
func ignore(level colltab.Level) (ignore [colltab.NumLevels]bool) ***REMOVED***
	if level < 0 ***REMOVED***
		level = colltab.Quaternary
	***REMOVED***
	for i := range ignore ***REMOVED***
		ignore[i] = level < colltab.Level(i)
	***REMOVED***
	return ignore
***REMOVED***

func makeCE(w []int) colltab.Elem ***REMOVED***
	ce, err := colltab.MakeElem(w[0], w[1], w[2], uint8(w[3]))
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return ce
***REMOVED***

func (o opts) collator() *Collator ***REMOVED***
	c := &Collator***REMOVED***
		options: options***REMOVED***
			ignore:      ignore(colltab.Level(o.lev - 1)),
			alternate:   o.alt,
			backwards:   o.backwards,
			caseLevel:   o.caseLevel,
			variableTop: uint32(o.top),
		***REMOVED***,
	***REMOVED***
	return c
***REMOVED***

const (
	maxQ = 0x1FFFFF
)

func wpq(p, q int) Weights ***REMOVED***
	return W(p, defaults.Secondary, defaults.Tertiary, q)
***REMOVED***

func wsq(s, q int) Weights ***REMOVED***
	return W(0, s, defaults.Tertiary, q)
***REMOVED***

func wq(q int) Weights ***REMOVED***
	return W(0, 0, 0, q)
***REMOVED***

var zero = W(0, 0, 0, 0)

var processTests = []weightsTest***REMOVED***
	// Shifted
	***REMOVED*** // simple sequence of non-variables
		opt: opts***REMOVED***alt: altShifted, top: 100***REMOVED***,
		in:  ColElems***REMOVED***W(200), W(300), W(400)***REMOVED***,
		out: ColElems***REMOVED***wpq(200, maxQ), wpq(300, maxQ), wpq(400, maxQ)***REMOVED***,
	***REMOVED***,
	***REMOVED*** // first is a variable
		opt: opts***REMOVED***alt: altShifted, top: 250***REMOVED***,
		in:  ColElems***REMOVED***W(200), W(300), W(400)***REMOVED***,
		out: ColElems***REMOVED***wq(200), wpq(300, maxQ), wpq(400, maxQ)***REMOVED***,
	***REMOVED***,
	***REMOVED*** // all but first are variable
		opt: opts***REMOVED***alt: altShifted, top: 999***REMOVED***,
		in:  ColElems***REMOVED***W(1000), W(200), W(300), W(400)***REMOVED***,
		out: ColElems***REMOVED***wpq(1000, maxQ), wq(200), wq(300), wq(400)***REMOVED***,
	***REMOVED***,
	***REMOVED*** // first is a modifier
		opt: opts***REMOVED***alt: altShifted, top: 999***REMOVED***,
		in:  ColElems***REMOVED***W(0, 10), W(1000)***REMOVED***,
		out: ColElems***REMOVED***wsq(10, maxQ), wpq(1000, maxQ)***REMOVED***,
	***REMOVED***,
	***REMOVED*** // primary ignorables
		opt: opts***REMOVED***alt: altShifted, top: 250***REMOVED***,
		in:  ColElems***REMOVED***W(200), W(0, 10), W(300), W(0, 15), W(400)***REMOVED***,
		out: ColElems***REMOVED***wq(200), zero, wpq(300, maxQ), wsq(15, maxQ), wpq(400, maxQ)***REMOVED***,
	***REMOVED***,
	***REMOVED*** // secondary ignorables
		opt: opts***REMOVED***alt: altShifted, top: 250***REMOVED***,
		in:  ColElems***REMOVED***W(200), W(0, 0, 10), W(300), W(0, 0, 15), W(400)***REMOVED***,
		out: ColElems***REMOVED***wq(200), zero, wpq(300, maxQ), W(0, 0, 15, maxQ), wpq(400, maxQ)***REMOVED***,
	***REMOVED***,
	***REMOVED*** // tertiary ignorables, no change
		opt: opts***REMOVED***alt: altShifted, top: 250***REMOVED***,
		in:  ColElems***REMOVED***W(200), zero, W(300), zero, W(400)***REMOVED***,
		out: ColElems***REMOVED***wq(200), zero, wpq(300, maxQ), zero, wpq(400, maxQ)***REMOVED***,
	***REMOVED***,

	// ShiftTrimmed (same as Shifted)
	***REMOVED*** // simple sequence of non-variables
		opt: opts***REMOVED***alt: altShiftTrimmed, top: 100***REMOVED***,
		in:  ColElems***REMOVED***W(200), W(300), W(400)***REMOVED***,
		out: ColElems***REMOVED***wpq(200, maxQ), wpq(300, maxQ), wpq(400, maxQ)***REMOVED***,
	***REMOVED***,
	***REMOVED*** // first is a variable
		opt: opts***REMOVED***alt: altShiftTrimmed, top: 250***REMOVED***,
		in:  ColElems***REMOVED***W(200), W(300), W(400)***REMOVED***,
		out: ColElems***REMOVED***wq(200), wpq(300, maxQ), wpq(400, maxQ)***REMOVED***,
	***REMOVED***,
	***REMOVED*** // all but first are variable
		opt: opts***REMOVED***alt: altShiftTrimmed, top: 999***REMOVED***,
		in:  ColElems***REMOVED***W(1000), W(200), W(300), W(400)***REMOVED***,
		out: ColElems***REMOVED***wpq(1000, maxQ), wq(200), wq(300), wq(400)***REMOVED***,
	***REMOVED***,
	***REMOVED*** // first is a modifier
		opt: opts***REMOVED***alt: altShiftTrimmed, top: 999***REMOVED***,
		in:  ColElems***REMOVED***W(0, 10), W(1000)***REMOVED***,
		out: ColElems***REMOVED***wsq(10, maxQ), wpq(1000, maxQ)***REMOVED***,
	***REMOVED***,
	***REMOVED*** // primary ignorables
		opt: opts***REMOVED***alt: altShiftTrimmed, top: 250***REMOVED***,
		in:  ColElems***REMOVED***W(200), W(0, 10), W(300), W(0, 15), W(400)***REMOVED***,
		out: ColElems***REMOVED***wq(200), zero, wpq(300, maxQ), wsq(15, maxQ), wpq(400, maxQ)***REMOVED***,
	***REMOVED***,
	***REMOVED*** // secondary ignorables
		opt: opts***REMOVED***alt: altShiftTrimmed, top: 250***REMOVED***,
		in:  ColElems***REMOVED***W(200), W(0, 0, 10), W(300), W(0, 0, 15), W(400)***REMOVED***,
		out: ColElems***REMOVED***wq(200), zero, wpq(300, maxQ), W(0, 0, 15, maxQ), wpq(400, maxQ)***REMOVED***,
	***REMOVED***,
	***REMOVED*** // tertiary ignorables, no change
		opt: opts***REMOVED***alt: altShiftTrimmed, top: 250***REMOVED***,
		in:  ColElems***REMOVED***W(200), zero, W(300), zero, W(400)***REMOVED***,
		out: ColElems***REMOVED***wq(200), zero, wpq(300, maxQ), zero, wpq(400, maxQ)***REMOVED***,
	***REMOVED***,

	// Blanked
	***REMOVED*** // simple sequence of non-variables
		opt: opts***REMOVED***alt: altBlanked, top: 100***REMOVED***,
		in:  ColElems***REMOVED***W(200), W(300), W(400)***REMOVED***,
		out: ColElems***REMOVED***W(200), W(300), W(400)***REMOVED***,
	***REMOVED***,
	***REMOVED*** // first is a variable
		opt: opts***REMOVED***alt: altBlanked, top: 250***REMOVED***,
		in:  ColElems***REMOVED***W(200), W(300), W(400)***REMOVED***,
		out: ColElems***REMOVED***zero, W(300), W(400)***REMOVED***,
	***REMOVED***,
	***REMOVED*** // all but first are variable
		opt: opts***REMOVED***alt: altBlanked, top: 999***REMOVED***,
		in:  ColElems***REMOVED***W(1000), W(200), W(300), W(400)***REMOVED***,
		out: ColElems***REMOVED***W(1000), zero, zero, zero***REMOVED***,
	***REMOVED***,
	***REMOVED*** // first is a modifier
		opt: opts***REMOVED***alt: altBlanked, top: 999***REMOVED***,
		in:  ColElems***REMOVED***W(0, 10), W(1000)***REMOVED***,
		out: ColElems***REMOVED***W(0, 10), W(1000)***REMOVED***,
	***REMOVED***,
	***REMOVED*** // primary ignorables
		opt: opts***REMOVED***alt: altBlanked, top: 250***REMOVED***,
		in:  ColElems***REMOVED***W(200), W(0, 10), W(300), W(0, 15), W(400)***REMOVED***,
		out: ColElems***REMOVED***zero, zero, W(300), W(0, 15), W(400)***REMOVED***,
	***REMOVED***,
	***REMOVED*** // secondary ignorables
		opt: opts***REMOVED***alt: altBlanked, top: 250***REMOVED***,
		in:  ColElems***REMOVED***W(200), W(0, 0, 10), W(300), W(0, 0, 15), W(400)***REMOVED***,
		out: ColElems***REMOVED***zero, zero, W(300), W(0, 0, 15), W(400)***REMOVED***,
	***REMOVED***,
	***REMOVED*** // tertiary ignorables, no change
		opt: opts***REMOVED***alt: altBlanked, top: 250***REMOVED***,
		in:  ColElems***REMOVED***W(200), zero, W(300), zero, W(400)***REMOVED***,
		out: ColElems***REMOVED***zero, zero, W(300), zero, W(400)***REMOVED***,
	***REMOVED***,

	// Non-ignorable: input is always equal to output.
	***REMOVED*** // all but first are variable
		opt: opts***REMOVED***alt: altNonIgnorable, top: 999***REMOVED***,
		in:  ColElems***REMOVED***W(1000), W(200), W(300), W(400)***REMOVED***,
		out: ColElems***REMOVED***W(1000), W(200), W(300), W(400)***REMOVED***,
	***REMOVED***,
	***REMOVED*** // primary ignorables
		opt: opts***REMOVED***alt: altNonIgnorable, top: 250***REMOVED***,
		in:  ColElems***REMOVED***W(200), W(0, 10), W(300), W(0, 15), W(400)***REMOVED***,
		out: ColElems***REMOVED***W(200), W(0, 10), W(300), W(0, 15), W(400)***REMOVED***,
	***REMOVED***,
	***REMOVED*** // secondary ignorables
		opt: opts***REMOVED***alt: altNonIgnorable, top: 250***REMOVED***,
		in:  ColElems***REMOVED***W(200), W(0, 0, 10), W(300), W(0, 0, 15), W(400)***REMOVED***,
		out: ColElems***REMOVED***W(200), W(0, 0, 10), W(300), W(0, 0, 15), W(400)***REMOVED***,
	***REMOVED***,
	***REMOVED*** // tertiary ignorables, no change
		opt: opts***REMOVED***alt: altNonIgnorable, top: 250***REMOVED***,
		in:  ColElems***REMOVED***W(200), zero, W(300), zero, W(400)***REMOVED***,
		out: ColElems***REMOVED***W(200), zero, W(300), zero, W(400)***REMOVED***,
	***REMOVED***,
***REMOVED***

func TestProcessWeights(t *testing.T) ***REMOVED***
	for i, tt := range processTests ***REMOVED***
		in := convertFromWeights(tt.in)
		out := convertFromWeights(tt.out)
		processWeights(tt.opt.alt, uint32(tt.opt.top), in)
		for j, w := range in ***REMOVED***
			if w != out[j] ***REMOVED***
				t.Errorf("%d: Weights %d was %v; want %v", i, j, w, out[j])
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

type keyFromElemTest struct ***REMOVED***
	opt opts
	in  ColElems
	out []byte
***REMOVED***

var defS = byte(defaults.Secondary)
var defT = byte(defaults.Tertiary)

const sep = 0 // separator byte

var keyFromElemTests = []keyFromElemTest***REMOVED***
	***REMOVED*** // simple primary and secondary weights.
		opts***REMOVED***alt: altShifted***REMOVED***,
		ColElems***REMOVED***W(0x200), W(0x7FFF), W(0, 0x30), W(0x100)***REMOVED***,
		[]byte***REMOVED***0x2, 0, 0x7F, 0xFF, 0x1, 0x00, // primary
			sep, sep, 0, defS, 0, defS, 0, 0x30, 0, defS, // secondary
			sep, sep, defT, defT, defT, defT, // tertiary
			sep, 0xFF, 0xFF, 0xFF, 0xFF, // quaternary
		***REMOVED***,
	***REMOVED***,
	***REMOVED*** // same as first, but with zero element that need to be removed
		opts***REMOVED***alt: altShifted***REMOVED***,
		ColElems***REMOVED***W(0x200), zero, W(0x7FFF), W(0, 0x30), zero, W(0x100)***REMOVED***,
		[]byte***REMOVED***0x2, 0, 0x7F, 0xFF, 0x1, 0x00, // primary
			sep, sep, 0, defS, 0, defS, 0, 0x30, 0, defS, // secondary
			sep, sep, defT, defT, defT, defT, // tertiary
			sep, 0xFF, 0xFF, 0xFF, 0xFF, // quaternary
		***REMOVED***,
	***REMOVED***,
	***REMOVED*** // same as first, with large primary values
		opts***REMOVED***alt: altShifted***REMOVED***,
		ColElems***REMOVED***W(0x200), W(0x8000), W(0, 0x30), W(0x12345)***REMOVED***,
		[]byte***REMOVED***0x2, 0, 0x80, 0x80, 0x00, 0x81, 0x23, 0x45, // primary
			sep, sep, 0, defS, 0, defS, 0, 0x30, 0, defS, // secondary
			sep, sep, defT, defT, defT, defT, // tertiary
			sep, 0xFF, 0xFF, 0xFF, 0xFF, // quaternary
		***REMOVED***,
	***REMOVED***,
	***REMOVED*** // same as first, but with the secondary level backwards
		opts***REMOVED***alt: altShifted, backwards: true***REMOVED***,
		ColElems***REMOVED***W(0x200), W(0x7FFF), W(0, 0x30), W(0x100)***REMOVED***,
		[]byte***REMOVED***0x2, 0, 0x7F, 0xFF, 0x1, 0x00, // primary
			sep, sep, 0, defS, 0, 0x30, 0, defS, 0, defS, // secondary
			sep, sep, defT, defT, defT, defT, // tertiary
			sep, 0xFF, 0xFF, 0xFF, 0xFF, // quaternary
		***REMOVED***,
	***REMOVED***,
	***REMOVED*** // same as first, ignoring quaternary level
		opts***REMOVED***alt: altShifted, lev: 3***REMOVED***,
		ColElems***REMOVED***W(0x200), zero, W(0x7FFF), W(0, 0x30), zero, W(0x100)***REMOVED***,
		[]byte***REMOVED***0x2, 0, 0x7F, 0xFF, 0x1, 0x00, // primary
			sep, sep, 0, defS, 0, defS, 0, 0x30, 0, defS, // secondary
			sep, sep, defT, defT, defT, defT, // tertiary
		***REMOVED***,
	***REMOVED***,
	***REMOVED*** // same as first, ignoring tertiary level
		opts***REMOVED***alt: altShifted, lev: 2***REMOVED***,
		ColElems***REMOVED***W(0x200), zero, W(0x7FFF), W(0, 0x30), zero, W(0x100)***REMOVED***,
		[]byte***REMOVED***0x2, 0, 0x7F, 0xFF, 0x1, 0x00, // primary
			sep, sep, 0, defS, 0, defS, 0, 0x30, 0, defS, // secondary
		***REMOVED***,
	***REMOVED***,
	***REMOVED*** // same as first, ignoring secondary level
		opts***REMOVED***alt: altShifted, lev: 1***REMOVED***,
		ColElems***REMOVED***W(0x200), zero, W(0x7FFF), W(0, 0x30), zero, W(0x100)***REMOVED***,
		[]byte***REMOVED***0x2, 0, 0x7F, 0xFF, 0x1, 0x00***REMOVED***,
	***REMOVED***,
	***REMOVED*** // simple primary and secondary weights.
		opts***REMOVED***alt: altShiftTrimmed, top: 0x250***REMOVED***,
		ColElems***REMOVED***W(0x300), W(0x200), W(0x7FFF), W(0, 0x30), W(0x800)***REMOVED***,
		[]byte***REMOVED***0x3, 0, 0x7F, 0xFF, 0x8, 0x00, // primary
			sep, sep, 0, defS, 0, defS, 0, 0x30, 0, defS, // secondary
			sep, sep, defT, defT, defT, defT, // tertiary
			sep, 0xFF, 0x2, 0, // quaternary
		***REMOVED***,
	***REMOVED***,
	***REMOVED*** // as first, primary with case level enabled
		opts***REMOVED***alt: altShifted, lev: 1, caseLevel: true***REMOVED***,
		ColElems***REMOVED***W(0x200), W(0x7FFF), W(0, 0x30), W(0x100)***REMOVED***,
		[]byte***REMOVED***0x2, 0, 0x7F, 0xFF, 0x1, 0x00, // primary
			sep, sep, // secondary
			sep, sep, defT, defT, defT, defT, // tertiary
		***REMOVED***,
	***REMOVED***,
***REMOVED***

func TestKeyFromElems(t *testing.T) ***REMOVED***
	buf := Buffer***REMOVED******REMOVED***
	for i, tt := range keyFromElemTests ***REMOVED***
		buf.Reset()
		in := convertFromWeights(tt.in)
		processWeights(tt.opt.alt, uint32(tt.opt.top), in)
		tt.opt.collator().keyFromElems(&buf, in)
		res := buf.key
		if len(res) != len(tt.out) ***REMOVED***
			t.Errorf("%d: len(ws) was %d; want %d (%X should be %X)", i, len(res), len(tt.out), res, tt.out)
		***REMOVED***
		n := len(res)
		if len(tt.out) < n ***REMOVED***
			n = len(tt.out)
		***REMOVED***
		for j, c := range res[:n] ***REMOVED***
			if c != tt.out[j] ***REMOVED***
				t.Errorf("%d: byte %d was %X; want %X", i, j, c, tt.out[j])
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestGetColElems(t *testing.T) ***REMOVED***
	for i, tt := range appendNextTests ***REMOVED***
		c, err := makeTable(tt.in)
		if err != nil ***REMOVED***
			// error is reported in TestAppendNext
			continue
		***REMOVED***
		// Create one large test per table
		str := make([]byte, 0, 4000)
		out := ColElems***REMOVED******REMOVED***
		for len(str) < 3000 ***REMOVED***
			for _, chk := range tt.chk ***REMOVED***
				str = append(str, chk.in[:chk.n]...)
				out = append(out, chk.out...)
			***REMOVED***
		***REMOVED***
		for j, chk := range append(tt.chk, check***REMOVED***string(str), len(str), out***REMOVED***) ***REMOVED***
			out := convertFromWeights(chk.out)
			ce := c.getColElems([]byte(chk.in)[:chk.n])
			if len(ce) != len(out) ***REMOVED***
				t.Errorf("%d:%d: len(ws) was %d; want %d", i, j, len(ce), len(out))
				continue
			***REMOVED***
			cnt := 0
			for k, w := range ce ***REMOVED***
				w, _ = colltab.MakeElem(w.Primary(), w.Secondary(), int(w.Tertiary()), 0)
				if w != out[k] ***REMOVED***
					t.Errorf("%d:%d: Weights %d was %X; want %X", i, j, k, w, out[k])
					cnt++
				***REMOVED***
				if cnt > 10 ***REMOVED***
					break
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

type keyTest struct ***REMOVED***
	in  string
	out []byte
***REMOVED***

var keyTests = []keyTest***REMOVED***
	***REMOVED***"abc",
		[]byte***REMOVED***0, 100, 0, 200, 1, 44, 0, 0, 0, 32, 0, 32, 0, 32, 0, 0, 2, 2, 2, 0, 255, 255, 255***REMOVED***,
	***REMOVED***,
	***REMOVED***"a\u0301",
		[]byte***REMOVED***0, 102, 0, 0, 0, 32, 0, 0, 2, 0, 255***REMOVED***,
	***REMOVED***,
	***REMOVED***"aaaaa",
		[]byte***REMOVED***0, 100, 0, 100, 0, 100, 0, 100, 0, 100, 0, 0,
			0, 32, 0, 32, 0, 32, 0, 32, 0, 32, 0, 0,
			2, 2, 2, 2, 2, 0,
			255, 255, 255, 255, 255,
		***REMOVED***,
	***REMOVED***,
	// Issue 16391: incomplete rune at end of UTF-8 sequence.
	***REMOVED***"\xc2", []byte***REMOVED***133, 255, 253, 0, 0, 0, 32, 0, 0, 2, 0, 255***REMOVED******REMOVED***,
	***REMOVED***"\xc2a", []byte***REMOVED***133, 255, 253, 0, 100, 0, 0, 0, 32, 0, 32, 0, 0, 2, 2, 0, 255, 255***REMOVED******REMOVED***,
***REMOVED***

func TestKey(t *testing.T) ***REMOVED***
	c, _ := makeTable(appendNextTests[4].in)
	c.alternate = altShifted
	c.ignore = ignore(colltab.Quaternary)
	buf := Buffer***REMOVED******REMOVED***
	keys1 := [][]byte***REMOVED******REMOVED***
	keys2 := [][]byte***REMOVED******REMOVED***
	for _, tt := range keyTests ***REMOVED***
		keys1 = append(keys1, c.Key(&buf, []byte(tt.in)))
		keys2 = append(keys2, c.KeyFromString(&buf, tt.in))
	***REMOVED***
	// Separate generation from testing to ensure buffers are not overwritten.
	for i, tt := range keyTests ***REMOVED***
		if !bytes.Equal(keys1[i], tt.out) ***REMOVED***
			t.Errorf("%d: Key(%q) = %d; want %d", i, tt.in, keys1[i], tt.out)
		***REMOVED***
		if !bytes.Equal(keys2[i], tt.out) ***REMOVED***
			t.Errorf("%d: KeyFromString(%q) = %d; want %d", i, tt.in, keys2[i], tt.out)
		***REMOVED***
	***REMOVED***
***REMOVED***

type compareTest struct ***REMOVED***
	a, b string
	res  int // comparison result
***REMOVED***

var compareTests = []compareTest***REMOVED***
	***REMOVED***"a\u0301", "a", 1***REMOVED***,
	***REMOVED***"a\u0301b", "ab", 1***REMOVED***,
	***REMOVED***"a", "a\u0301", -1***REMOVED***,
	***REMOVED***"ab", "a\u0301b", -1***REMOVED***,
	***REMOVED***"bc", "a\u0301c", 1***REMOVED***,
	***REMOVED***"ab", "aB", -1***REMOVED***,
	***REMOVED***"a\u0301", "a\u0301", 0***REMOVED***,
	***REMOVED***"a", "a", 0***REMOVED***,
	// Only clip prefixes of whole runes.
	***REMOVED***"\u302E", "\u302F", 1***REMOVED***,
	// Don't clip prefixes when last rune of prefix may be part of contraction.
	***REMOVED***"a\u035E", "a\u0301\u035F", -1***REMOVED***,
	***REMOVED***"a\u0301\u035Fb", "a\u0301\u035F", -1***REMOVED***,
***REMOVED***

func TestCompare(t *testing.T) ***REMOVED***
	c, _ := makeTable(appendNextTests[4].in)
	for i, tt := range compareTests ***REMOVED***
		if res := c.Compare([]byte(tt.a), []byte(tt.b)); res != tt.res ***REMOVED***
			t.Errorf("%d: Compare(%q, %q) == %d; want %d", i, tt.a, tt.b, res, tt.res)
		***REMOVED***
		if res := c.CompareString(tt.a, tt.b); res != tt.res ***REMOVED***
			t.Errorf("%d: CompareString(%q, %q) == %d; want %d", i, tt.a, tt.b, res, tt.res)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestNumeric(t *testing.T) ***REMOVED***
	c := New(language.English, Loose, Numeric)

	for i, tt := range []struct ***REMOVED***
		a, b string
		want int
	***REMOVED******REMOVED***
		***REMOVED***"1", "2", -1***REMOVED***,
		***REMOVED***"2", "12", -1***REMOVED***,
		***REMOVED***"２", "１２", -1***REMOVED***, // Fullwidth is sorted as usual.
		***REMOVED***"₂", "₁₂", 1***REMOVED***,  // Subscript is not sorted as numbers.
		***REMOVED***"②", "①②", 1***REMOVED***,  // Circled is not sorted as numbers.
		***REMOVED*** // Imperial Aramaic, is not sorted as number.
			"\U00010859",
			"\U00010858\U00010859",
			1,
		***REMOVED***,
		***REMOVED***"12", "2", 1***REMOVED***,
		***REMOVED***"A-1", "A-2", -1***REMOVED***,
		***REMOVED***"A-2", "A-12", -1***REMOVED***,
		***REMOVED***"A-12", "A-2", 1***REMOVED***,
		***REMOVED***"A-0001", "A-1", 0***REMOVED***,
	***REMOVED*** ***REMOVED***
		if got := c.CompareString(tt.a, tt.b); got != tt.want ***REMOVED***
			t.Errorf("%d: CompareString(%s, %s) = %d; want %d", i, tt.a, tt.b, got, tt.want)
		***REMOVED***
	***REMOVED***
***REMOVED***
