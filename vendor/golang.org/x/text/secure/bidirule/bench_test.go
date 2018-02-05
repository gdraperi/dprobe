// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bidirule

import (
	"testing"

	"golang.org/x/text/internal/testtext"
)

var benchData = []struct***REMOVED*** name, data string ***REMOVED******REMOVED***
	***REMOVED***"ascii", "Scheveningen"***REMOVED***,
	***REMOVED***"arabic", "دبي"***REMOVED***,
	***REMOVED***"hangul", "다음과"***REMOVED***,
***REMOVED***

func doBench(b *testing.B, fn func(b *testing.B, data string)) ***REMOVED***
	for _, d := range benchData ***REMOVED***
		testtext.Bench(b, d.name, func(b *testing.B) ***REMOVED*** fn(b, d.data) ***REMOVED***)
	***REMOVED***
***REMOVED***

func BenchmarkSpan(b *testing.B) ***REMOVED***
	r := New()
	doBench(b, func(b *testing.B, str string) ***REMOVED***
		b.SetBytes(int64(len(str)))
		data := []byte(str)
		for i := 0; i < b.N; i++ ***REMOVED***
			r.Reset()
			r.Span(data, true)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func BenchmarkDirectionASCII(b *testing.B) ***REMOVED***
	doBench(b, func(b *testing.B, str string) ***REMOVED***
		b.SetBytes(int64(len(str)))
		data := []byte(str)
		for i := 0; i < b.N; i++ ***REMOVED***
			Direction(data)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func BenchmarkDirectionStringASCII(b *testing.B) ***REMOVED***
	doBench(b, func(b *testing.B, str string) ***REMOVED***
		b.SetBytes(int64(len(str)))
		for i := 0; i < b.N; i++ ***REMOVED***
			DirectionString(str)
		***REMOVED***
	***REMOVED***)
***REMOVED***
