// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.7

package precis

import (
	"testing"

	"golang.org/x/text/internal/testtext"
)

var benchData = []struct***REMOVED*** name, str string ***REMOVED******REMOVED***
	***REMOVED***"ASCII", "Malvolio"***REMOVED***,
	***REMOVED***"NotNormalized", "abcdefg\u0301\u031f"***REMOVED***,
	***REMOVED***"Arabic", "دبي"***REMOVED***,
	***REMOVED***"Hangul", "동일조건변경허락"***REMOVED***,
***REMOVED***

var benchProfiles = []struct ***REMOVED***
	name string
	p    *Profile
***REMOVED******REMOVED***
	***REMOVED***"FreeForm", NewFreeform()***REMOVED***,
	***REMOVED***"Nickname", Nickname***REMOVED***,
	***REMOVED***"OpaqueString", OpaqueString***REMOVED***,
	***REMOVED***"UsernameCaseMapped", UsernameCaseMapped***REMOVED***,
	***REMOVED***"UsernameCasePreserved", UsernameCasePreserved***REMOVED***,
***REMOVED***

func doBench(b *testing.B, f func(b *testing.B, p *Profile, s string)) ***REMOVED***
	for _, bp := range benchProfiles ***REMOVED***
		for _, d := range benchData ***REMOVED***
			testtext.Bench(b, bp.name+"/"+d.name, func(b *testing.B) ***REMOVED***
				f(b, bp.p, d.str)
			***REMOVED***)
		***REMOVED***
	***REMOVED***
***REMOVED***

func BenchmarkString(b *testing.B) ***REMOVED***
	doBench(b, func(b *testing.B, p *Profile, s string) ***REMOVED***
		for i := 0; i < b.N; i++ ***REMOVED***
			p.String(s)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func BenchmarkBytes(b *testing.B) ***REMOVED***
	doBench(b, func(b *testing.B, p *Profile, s string) ***REMOVED***
		src := []byte(s)
		b.ResetTimer()
		for i := 0; i < b.N; i++ ***REMOVED***
			p.Bytes(src)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func BenchmarkAppend(b *testing.B) ***REMOVED***
	doBench(b, func(b *testing.B, p *Profile, s string) ***REMOVED***
		src := []byte(s)
		dst := make([]byte, 0, 4096)
		b.ResetTimer()
		for i := 0; i < b.N; i++ ***REMOVED***
			p.Append(dst, src)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func BenchmarkTransform(b *testing.B) ***REMOVED***
	doBench(b, func(b *testing.B, p *Profile, s string) ***REMOVED***
		src := []byte(s)
		dst := make([]byte, 2*len(s))
		t := p.NewTransformer()
		b.ResetTimer()
		for i := 0; i < b.N; i++ ***REMOVED***
			t.Transform(dst, src, true)
		***REMOVED***
	***REMOVED***)
***REMOVED***
