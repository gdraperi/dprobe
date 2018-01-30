package goquery

import (
	"testing"
)

func BenchmarkAttr(b *testing.B) ***REMOVED***
	var s string

	b.StopTimer()
	sel := DocW().Find("h1")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		s, _ = sel.Attr("id")
	***REMOVED***
	if s != "firstHeading" ***REMOVED***
		b.Fatalf("want firstHeading, got %q", s)
	***REMOVED***
***REMOVED***

func BenchmarkText(b *testing.B) ***REMOVED***
	b.StopTimer()
	sel := DocW().Find("h2")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		sel.Text()
	***REMOVED***
***REMOVED***

func BenchmarkLength(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocW().Find("h2")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		n = sel.Length()
	***REMOVED***
	if n != 14 ***REMOVED***
		b.Fatalf("want 14, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkHtml(b *testing.B) ***REMOVED***
	b.StopTimer()
	sel := DocW().Find("h2")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		sel.Html()
	***REMOVED***
***REMOVED***
