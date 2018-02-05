// Copyright 2017 Frank Schroeder. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package properties

import (
	"reflect"
	"testing"
	"time"
)

func TestDecodeValues(t *testing.T) ***REMOVED***
	type S struct ***REMOVED***
		S   string
		BT  bool
		BF  bool
		I   int
		I8  int8
		I16 int16
		I32 int32
		I64 int64
		U   uint
		U8  uint8
		U16 uint16
		U32 uint32
		U64 uint64
		F32 float32
		F64 float64
		D   time.Duration
		TM  time.Time
	***REMOVED***
	in := `
	S=abc
	BT=true
	BF=false
	I=-1
	I8=-8
	I16=-16
	I32=-32
	I64=-64
	U=1
	U8=8
	U16=16
	U32=32
	U64=64
	F32=3.2
	F64=6.4
	D=5s
	TM=2015-01-02T12:34:56Z
	`
	out := &S***REMOVED***
		S:   "abc",
		BT:  true,
		BF:  false,
		I:   -1,
		I8:  -8,
		I16: -16,
		I32: -32,
		I64: -64,
		U:   1,
		U8:  8,
		U16: 16,
		U32: 32,
		U64: 64,
		F32: 3.2,
		F64: 6.4,
		D:   5 * time.Second,
		TM:  tm(t, time.RFC3339, "2015-01-02T12:34:56Z"),
	***REMOVED***
	testDecode(t, in, &S***REMOVED******REMOVED***, out)
***REMOVED***

func TestDecodeValueDefaults(t *testing.T) ***REMOVED***
	type S struct ***REMOVED***
		S   string        `properties:",default=abc"`
		BT  bool          `properties:",default=true"`
		BF  bool          `properties:",default=false"`
		I   int           `properties:",default=-1"`
		I8  int8          `properties:",default=-8"`
		I16 int16         `properties:",default=-16"`
		I32 int32         `properties:",default=-32"`
		I64 int64         `properties:",default=-64"`
		U   uint          `properties:",default=1"`
		U8  uint8         `properties:",default=8"`
		U16 uint16        `properties:",default=16"`
		U32 uint32        `properties:",default=32"`
		U64 uint64        `properties:",default=64"`
		F32 float32       `properties:",default=3.2"`
		F64 float64       `properties:",default=6.4"`
		D   time.Duration `properties:",default=5s"`
		TM  time.Time     `properties:",default=2015-01-02T12:34:56Z"`
	***REMOVED***
	out := &S***REMOVED***
		S:   "abc",
		BT:  true,
		BF:  false,
		I:   -1,
		I8:  -8,
		I16: -16,
		I32: -32,
		I64: -64,
		U:   1,
		U8:  8,
		U16: 16,
		U32: 32,
		U64: 64,
		F32: 3.2,
		F64: 6.4,
		D:   5 * time.Second,
		TM:  tm(t, time.RFC3339, "2015-01-02T12:34:56Z"),
	***REMOVED***
	testDecode(t, "", &S***REMOVED******REMOVED***, out)
***REMOVED***

func TestDecodeArrays(t *testing.T) ***REMOVED***
	type S struct ***REMOVED***
		S   []string
		B   []bool
		I   []int
		I8  []int8
		I16 []int16
		I32 []int32
		I64 []int64
		U   []uint
		U8  []uint8
		U16 []uint16
		U32 []uint32
		U64 []uint64
		F32 []float32
		F64 []float64
		D   []time.Duration
		TM  []time.Time
	***REMOVED***
	in := `
	S=a;b
	B=true;false
	I=-1;-2
	I8=-8;-9
	I16=-16;-17
	I32=-32;-33
	I64=-64;-65
	U=1;2
	U8=8;9
	U16=16;17
	U32=32;33
	U64=64;65
	F32=3.2;3.3
	F64=6.4;6.5
	D=4s;5s
	TM=2015-01-01T00:00:00Z;2016-01-01T00:00:00Z
	`
	out := &S***REMOVED***
		S:   []string***REMOVED***"a", "b"***REMOVED***,
		B:   []bool***REMOVED***true, false***REMOVED***,
		I:   []int***REMOVED***-1, -2***REMOVED***,
		I8:  []int8***REMOVED***-8, -9***REMOVED***,
		I16: []int16***REMOVED***-16, -17***REMOVED***,
		I32: []int32***REMOVED***-32, -33***REMOVED***,
		I64: []int64***REMOVED***-64, -65***REMOVED***,
		U:   []uint***REMOVED***1, 2***REMOVED***,
		U8:  []uint8***REMOVED***8, 9***REMOVED***,
		U16: []uint16***REMOVED***16, 17***REMOVED***,
		U32: []uint32***REMOVED***32, 33***REMOVED***,
		U64: []uint64***REMOVED***64, 65***REMOVED***,
		F32: []float32***REMOVED***3.2, 3.3***REMOVED***,
		F64: []float64***REMOVED***6.4, 6.5***REMOVED***,
		D:   []time.Duration***REMOVED***4 * time.Second, 5 * time.Second***REMOVED***,
		TM:  []time.Time***REMOVED***tm(t, time.RFC3339, "2015-01-01T00:00:00Z"), tm(t, time.RFC3339, "2016-01-01T00:00:00Z")***REMOVED***,
	***REMOVED***
	testDecode(t, in, &S***REMOVED******REMOVED***, out)
***REMOVED***

func TestDecodeArrayDefaults(t *testing.T) ***REMOVED***
	type S struct ***REMOVED***
		S   []string        `properties:",default=a;b"`
		B   []bool          `properties:",default=true;false"`
		I   []int           `properties:",default=-1;-2"`
		I8  []int8          `properties:",default=-8;-9"`
		I16 []int16         `properties:",default=-16;-17"`
		I32 []int32         `properties:",default=-32;-33"`
		I64 []int64         `properties:",default=-64;-65"`
		U   []uint          `properties:",default=1;2"`
		U8  []uint8         `properties:",default=8;9"`
		U16 []uint16        `properties:",default=16;17"`
		U32 []uint32        `properties:",default=32;33"`
		U64 []uint64        `properties:",default=64;65"`
		F32 []float32       `properties:",default=3.2;3.3"`
		F64 []float64       `properties:",default=6.4;6.5"`
		D   []time.Duration `properties:",default=4s;5s"`
		TM  []time.Time     `properties:",default=2015-01-01T00:00:00Z;2016-01-01T00:00:00Z"`
	***REMOVED***
	out := &S***REMOVED***
		S:   []string***REMOVED***"a", "b"***REMOVED***,
		B:   []bool***REMOVED***true, false***REMOVED***,
		I:   []int***REMOVED***-1, -2***REMOVED***,
		I8:  []int8***REMOVED***-8, -9***REMOVED***,
		I16: []int16***REMOVED***-16, -17***REMOVED***,
		I32: []int32***REMOVED***-32, -33***REMOVED***,
		I64: []int64***REMOVED***-64, -65***REMOVED***,
		U:   []uint***REMOVED***1, 2***REMOVED***,
		U8:  []uint8***REMOVED***8, 9***REMOVED***,
		U16: []uint16***REMOVED***16, 17***REMOVED***,
		U32: []uint32***REMOVED***32, 33***REMOVED***,
		U64: []uint64***REMOVED***64, 65***REMOVED***,
		F32: []float32***REMOVED***3.2, 3.3***REMOVED***,
		F64: []float64***REMOVED***6.4, 6.5***REMOVED***,
		D:   []time.Duration***REMOVED***4 * time.Second, 5 * time.Second***REMOVED***,
		TM:  []time.Time***REMOVED***tm(t, time.RFC3339, "2015-01-01T00:00:00Z"), tm(t, time.RFC3339, "2016-01-01T00:00:00Z")***REMOVED***,
	***REMOVED***
	testDecode(t, "", &S***REMOVED******REMOVED***, out)
***REMOVED***

func TestDecodeSkipUndef(t *testing.T) ***REMOVED***
	type S struct ***REMOVED***
		X     string `properties:"-"`
		Undef string `properties:",default=some value"`
	***REMOVED***
	in := `X=ignore`
	out := &S***REMOVED***"", "some value"***REMOVED***
	testDecode(t, in, &S***REMOVED******REMOVED***, out)
***REMOVED***

func TestDecodeStruct(t *testing.T) ***REMOVED***
	type A struct ***REMOVED***
		S string
		T string `properties:"t"`
		U string `properties:"u,default=uuu"`
	***REMOVED***
	type S struct ***REMOVED***
		A A
		B A `properties:"b"`
	***REMOVED***
	in := `
	A.S=sss
	A.t=ttt
	b.S=SSS
	b.t=TTT
	`
	out := &S***REMOVED***
		A***REMOVED***S: "sss", T: "ttt", U: "uuu"***REMOVED***,
		A***REMOVED***S: "SSS", T: "TTT", U: "uuu"***REMOVED***,
	***REMOVED***
	testDecode(t, in, &S***REMOVED******REMOVED***, out)
***REMOVED***

func TestDecodeMap(t *testing.T) ***REMOVED***
	type S struct ***REMOVED***
		A string `properties:"a"`
	***REMOVED***
	type X struct ***REMOVED***
		A map[string]string
		B map[string][]string
		C map[string]map[string]string
		D map[string]S
		E map[string]int
		F map[string]int `properties:"-"`
	***REMOVED***
	in := `
	A.foo=bar
	A.bar=bang
	B.foo=a;b;c
	B.bar=1;2;3
	C.foo.one=1
	C.foo.two=2
	C.bar.three=3
	C.bar.four=4
	D.foo.a=bar
	`
	out := &X***REMOVED***
		A: map[string]string***REMOVED***"foo": "bar", "bar": "bang"***REMOVED***,
		B: map[string][]string***REMOVED***"foo": []string***REMOVED***"a", "b", "c"***REMOVED***, "bar": []string***REMOVED***"1", "2", "3"***REMOVED******REMOVED***,
		C: map[string]map[string]string***REMOVED***"foo": map[string]string***REMOVED***"one": "1", "two": "2"***REMOVED***, "bar": map[string]string***REMOVED***"three": "3", "four": "4"***REMOVED******REMOVED***,
		D: map[string]S***REMOVED***"foo": S***REMOVED***"bar"***REMOVED******REMOVED***,
		E: map[string]int***REMOVED******REMOVED***,
	***REMOVED***
	testDecode(t, in, &X***REMOVED******REMOVED***, out)
***REMOVED***

func testDecode(t *testing.T, in string, v, out interface***REMOVED******REMOVED***) ***REMOVED***
	p, err := parse(in)
	if err != nil ***REMOVED***
		t.Fatalf("got %v want nil", err)
	***REMOVED***
	if err := p.Decode(v); err != nil ***REMOVED***
		t.Fatalf("got %v want nil", err)
	***REMOVED***
	if got, want := v, out; !reflect.DeepEqual(got, want) ***REMOVED***
		t.Fatalf("\ngot  %+v\nwant %+v", got, want)
	***REMOVED***
***REMOVED***

func tm(t *testing.T, layout, s string) time.Time ***REMOVED***
	tm, err := time.Parse(layout, s)
	if err != nil ***REMOVED***
		t.Fatalf("got %v want nil", err)
	***REMOVED***
	return tm
***REMOVED***
