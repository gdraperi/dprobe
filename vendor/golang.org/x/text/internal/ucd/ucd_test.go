package ucd

import (
	"strings"
	"testing"
)

const file = `
# Comments should be skipped
# rune;  bool;  uint; int; float; runes; # Y
0..0005; Y;     0;    2;      -5.25 ;  0 1 2 3 4 5;
6..0007; Yes  ; 6;    1;     -4.25  ;  0006 0007;
8;       T ;    8 ;   0 ;-3.25  ;;# T
9;       True  ;9  ;  -1;-2.25  ;  0009;

# more comments to be ignored
@Part0  

A;       N;   10  ;   -2;  -1.25; ;# N
B;       No;   11 ;   -3;  -0.25; 
C;        False;12;   -4;   0.75;
D;        ;13;-5;1.75;

@Part1   # Another part. 
# We test part comments get removed by not commenting the the next line.
E..10FFFF; F;   14  ; -6;   2.75;
`

var want = []struct ***REMOVED***
	start, end rune
***REMOVED******REMOVED***
	***REMOVED***0x00, 0x05***REMOVED***,
	***REMOVED***0x06, 0x07***REMOVED***,
	***REMOVED***0x08, 0x08***REMOVED***,
	***REMOVED***0x09, 0x09***REMOVED***,
	***REMOVED***0x0A, 0x0A***REMOVED***,
	***REMOVED***0x0B, 0x0B***REMOVED***,
	***REMOVED***0x0C, 0x0C***REMOVED***,
	***REMOVED***0x0D, 0x0D***REMOVED***,
	***REMOVED***0x0E, 0x10FFFF***REMOVED***,
***REMOVED***

func TestGetters(t *testing.T) ***REMOVED***
	parts := [][2]string***REMOVED***
		***REMOVED***"Part0", ""***REMOVED***,
		***REMOVED***"Part1", "Another part."***REMOVED***,
	***REMOVED***
	handler := func(p *Parser) ***REMOVED***
		if len(parts) == 0 ***REMOVED***
			t.Error("Part handler invoked too many times.")
			return
		***REMOVED***
		want := parts[0]
		parts = parts[1:]
		if got0, got1 := p.String(0), p.Comment(); got0 != want[0] || got1 != want[1] ***REMOVED***
			t.Errorf(`part: got %q, %q; want %q"`, got0, got1, want)
		***REMOVED***
	***REMOVED***

	p := New(strings.NewReader(file), KeepRanges, Part(handler))
	for i := 0; p.Next(); i++ ***REMOVED***
		start, end := p.Range(0)
		w := want[i]
		if start != w.start || end != w.end ***REMOVED***
			t.Fatalf("%d:Range(0); got %#x..%#x; want %#x..%#x", i, start, end, w.start, w.end)
		***REMOVED***
		if w.start == w.end && p.Rune(0) != w.start ***REMOVED***
			t.Errorf("%d:Range(0).start: got %U; want %U", i, p.Rune(0), w.start)
		***REMOVED***
		if got, want := p.Bool(1), w.start <= 9; got != want ***REMOVED***
			t.Errorf("%d:Bool(1): got %v; want %v", i, got, want)
		***REMOVED***
		if got := p.Rune(4); got != 0 || p.Err() == nil ***REMOVED***
			t.Errorf("%d:Rune(%q): got no error; want error", i, p.String(1))
		***REMOVED***
		p.err = nil
		if got := p.Uint(2); rune(got) != start ***REMOVED***
			t.Errorf("%d:Uint(2): got %v; want %v", i, got, start)
		***REMOVED***
		if got, want := p.Int(3), 2-i; got != want ***REMOVED***
			t.Errorf("%d:Int(3): got %v; want %v", i, got, want)
		***REMOVED***
		if got, want := p.Float(4), -5.25+float64(i); got != want ***REMOVED***
			t.Errorf("%d:Int(3): got %v; want %v", i, got, want)
		***REMOVED***
		if got := p.Runes(5); got == nil ***REMOVED***
			if p.String(5) != "" ***REMOVED***
				t.Errorf("%d:Runes(5): expected non-empty list", i)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if got[0] != start || got[len(got)-1] != end ***REMOVED***
				t.Errorf("%d:Runes(5): got %#x; want %#x..%#x", i, got, start, end)
			***REMOVED***
		***REMOVED***
		if got := p.Comment(); got != "" && got != p.String(1) ***REMOVED***
			t.Errorf("%d:Comment(): got %v; want %v", i, got, p.String(1))
		***REMOVED***
	***REMOVED***
	if err := p.Err(); err != nil ***REMOVED***
		t.Errorf("Parser error: %v", err)
	***REMOVED***
	if len(parts) != 0 ***REMOVED***
		t.Errorf("expected %d more invocations of part handler", len(parts))
	***REMOVED***
***REMOVED***
