// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package trace

import (
	"net/http"
	"reflect"
	"testing"
)

type s struct***REMOVED******REMOVED***

func (s) String() string ***REMOVED*** return "lazy string" ***REMOVED***

// TestReset checks whether all the fields are zeroed after reset.
func TestReset(t *testing.T) ***REMOVED***
	tr := New("foo", "bar")
	tr.LazyLog(s***REMOVED******REMOVED***, false)
	tr.LazyPrintf("%d", 1)
	tr.SetRecycler(func(_ interface***REMOVED******REMOVED***) ***REMOVED******REMOVED***)
	tr.SetTraceInfo(3, 4)
	tr.SetMaxEvents(100)
	tr.SetError()
	tr.Finish()

	tr.(*trace).reset()

	if !reflect.DeepEqual(tr, new(trace)) ***REMOVED***
		t.Errorf("reset didn't clear all fields: %+v", tr)
	***REMOVED***
***REMOVED***

// TestResetLog checks whether all the fields are zeroed after reset.
func TestResetLog(t *testing.T) ***REMOVED***
	el := NewEventLog("foo", "bar")
	el.Printf("message")
	el.Errorf("error")
	el.Finish()

	el.(*eventLog).reset()

	if !reflect.DeepEqual(el, new(eventLog)) ***REMOVED***
		t.Errorf("reset didn't clear all fields: %+v", el)
	***REMOVED***
***REMOVED***

func TestAuthRequest(t *testing.T) ***REMOVED***
	testCases := []struct ***REMOVED***
		host string
		want bool
	***REMOVED******REMOVED***
		***REMOVED***host: "192.168.23.1", want: false***REMOVED***,
		***REMOVED***host: "192.168.23.1:8080", want: false***REMOVED***,
		***REMOVED***host: "malformed remote addr", want: false***REMOVED***,
		***REMOVED***host: "localhost", want: true***REMOVED***,
		***REMOVED***host: "localhost:8080", want: true***REMOVED***,
		***REMOVED***host: "127.0.0.1", want: true***REMOVED***,
		***REMOVED***host: "127.0.0.1:8080", want: true***REMOVED***,
		***REMOVED***host: "::1", want: true***REMOVED***,
		***REMOVED***host: "[::1]:8080", want: true***REMOVED***,
	***REMOVED***
	for _, tt := range testCases ***REMOVED***
		req := &http.Request***REMOVED***RemoteAddr: tt.host***REMOVED***
		any, sensitive := AuthRequest(req)
		if any != tt.want || sensitive != tt.want ***REMOVED***
			t.Errorf("AuthRequest(%q) = %t, %t; want %t, %t", tt.host, any, sensitive, tt.want, tt.want)
		***REMOVED***
	***REMOVED***
***REMOVED***

// TestParseTemplate checks that all templates used by this package are valid
// as they are parsed on first usage
func TestParseTemplate(t *testing.T) ***REMOVED***
	if tmpl := distTmpl(); tmpl == nil ***REMOVED***
		t.Error("invalid template returned from distTmpl()")
	***REMOVED***
	if tmpl := pageTmpl(); tmpl == nil ***REMOVED***
		t.Error("invalid template returned from pageTmpl()")
	***REMOVED***
	if tmpl := eventsTmpl(); tmpl == nil ***REMOVED***
		t.Error("invalid template returned from eventsTmpl()")
	***REMOVED***
***REMOVED***

func benchmarkTrace(b *testing.B, maxEvents, numEvents int) ***REMOVED***
	numSpans := (b.N + numEvents + 1) / numEvents

	for i := 0; i < numSpans; i++ ***REMOVED***
		tr := New("test", "test")
		tr.SetMaxEvents(maxEvents)
		for j := 0; j < numEvents; j++ ***REMOVED***
			tr.LazyPrintf("%d", j)
		***REMOVED***
		tr.Finish()
	***REMOVED***
***REMOVED***

func BenchmarkTrace_Default_2(b *testing.B) ***REMOVED***
	benchmarkTrace(b, 0, 2)
***REMOVED***

func BenchmarkTrace_Default_10(b *testing.B) ***REMOVED***
	benchmarkTrace(b, 0, 10)
***REMOVED***

func BenchmarkTrace_Default_100(b *testing.B) ***REMOVED***
	benchmarkTrace(b, 0, 100)
***REMOVED***

func BenchmarkTrace_Default_1000(b *testing.B) ***REMOVED***
	benchmarkTrace(b, 0, 1000)
***REMOVED***

func BenchmarkTrace_Default_10000(b *testing.B) ***REMOVED***
	benchmarkTrace(b, 0, 10000)
***REMOVED***

func BenchmarkTrace_10_2(b *testing.B) ***REMOVED***
	benchmarkTrace(b, 10, 2)
***REMOVED***

func BenchmarkTrace_10_10(b *testing.B) ***REMOVED***
	benchmarkTrace(b, 10, 10)
***REMOVED***

func BenchmarkTrace_10_100(b *testing.B) ***REMOVED***
	benchmarkTrace(b, 10, 100)
***REMOVED***

func BenchmarkTrace_10_1000(b *testing.B) ***REMOVED***
	benchmarkTrace(b, 10, 1000)
***REMOVED***

func BenchmarkTrace_10_10000(b *testing.B) ***REMOVED***
	benchmarkTrace(b, 10, 10000)
***REMOVED***

func BenchmarkTrace_100_2(b *testing.B) ***REMOVED***
	benchmarkTrace(b, 100, 2)
***REMOVED***

func BenchmarkTrace_100_10(b *testing.B) ***REMOVED***
	benchmarkTrace(b, 100, 10)
***REMOVED***

func BenchmarkTrace_100_100(b *testing.B) ***REMOVED***
	benchmarkTrace(b, 100, 100)
***REMOVED***

func BenchmarkTrace_100_1000(b *testing.B) ***REMOVED***
	benchmarkTrace(b, 100, 1000)
***REMOVED***

func BenchmarkTrace_100_10000(b *testing.B) ***REMOVED***
	benchmarkTrace(b, 100, 10000)
***REMOVED***

func BenchmarkTrace_1000_2(b *testing.B) ***REMOVED***
	benchmarkTrace(b, 1000, 2)
***REMOVED***

func BenchmarkTrace_1000_10(b *testing.B) ***REMOVED***
	benchmarkTrace(b, 1000, 10)
***REMOVED***

func BenchmarkTrace_1000_100(b *testing.B) ***REMOVED***
	benchmarkTrace(b, 1000, 100)
***REMOVED***

func BenchmarkTrace_1000_1000(b *testing.B) ***REMOVED***
	benchmarkTrace(b, 1000, 1000)
***REMOVED***

func BenchmarkTrace_1000_10000(b *testing.B) ***REMOVED***
	benchmarkTrace(b, 1000, 10000)
***REMOVED***
