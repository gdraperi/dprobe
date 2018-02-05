package rangetable

import (
	"reflect"
	"testing"
	"unicode"
)

var (
	empty = &unicode.RangeTable***REMOVED******REMOVED***
	many  = &unicode.RangeTable***REMOVED***
		R16:         []unicode.Range16***REMOVED******REMOVED***0, 0xffff, 5***REMOVED******REMOVED***,
		R32:         []unicode.Range32***REMOVED******REMOVED***0x10004, 0x10009, 5***REMOVED******REMOVED***,
		LatinOffset: 0,
	***REMOVED***
)

func TestVisit(t *testing.T) ***REMOVED***
	Visit(empty, func(got rune) ***REMOVED***
		t.Error("call from empty RangeTable")
	***REMOVED***)

	var want rune
	Visit(many, func(got rune) ***REMOVED***
		if got != want ***REMOVED***
			t.Errorf("got %U; want %U", got, want)
		***REMOVED***
		want += 5
	***REMOVED***)
	if want -= 5; want != 0x10009 ***REMOVED***
		t.Errorf("last run was %U; want U+10009", want)
	***REMOVED***
***REMOVED***

func TestNew(t *testing.T) ***REMOVED***
	for i, rt := range []*unicode.RangeTable***REMOVED***
		empty,
		unicode.Co,
		unicode.Letter,
		unicode.ASCII_Hex_Digit,
		many,
		maxRuneTable,
	***REMOVED*** ***REMOVED***
		var got, want []rune
		Visit(rt, func(r rune) ***REMOVED***
			want = append(want, r)
		***REMOVED***)
		Visit(New(want...), func(r rune) ***REMOVED***
			got = append(got, r)
		***REMOVED***)
		if !reflect.DeepEqual(got, want) ***REMOVED***
			t.Errorf("%d:\ngot  %v;\nwant %v", i, got, want)
		***REMOVED***
	***REMOVED***
***REMOVED***
