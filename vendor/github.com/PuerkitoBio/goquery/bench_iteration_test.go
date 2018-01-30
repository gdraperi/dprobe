package goquery

import (
	"testing"
)

func BenchmarkEach(b *testing.B) ***REMOVED***
	var tmp, n int

	b.StopTimer()
	sel := DocW().Find("td")
	f := func(i int, s *Selection) ***REMOVED***
		tmp++
	***REMOVED***
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		sel.Each(f)
		if n == 0 ***REMOVED***
			n = tmp
		***REMOVED***
	***REMOVED***
	if n != 59 ***REMOVED***
		b.Fatalf("want 59, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkMap(b *testing.B) ***REMOVED***
	var tmp, n int

	b.StopTimer()
	sel := DocW().Find("td")
	f := func(i int, s *Selection) string ***REMOVED***
		tmp++
		return string(tmp)
	***REMOVED***
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		sel.Map(f)
		if n == 0 ***REMOVED***
			n = tmp
		***REMOVED***
	***REMOVED***
	if n != 59 ***REMOVED***
		b.Fatalf("want 59, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkEachWithBreak(b *testing.B) ***REMOVED***
	var tmp, n int

	b.StopTimer()
	sel := DocW().Find("td")
	f := func(i int, s *Selection) bool ***REMOVED***
		tmp++
		return tmp < 10
	***REMOVED***
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		tmp = 0
		sel.EachWithBreak(f)
		if n == 0 ***REMOVED***
			n = tmp
		***REMOVED***
	***REMOVED***
	if n != 10 ***REMOVED***
		b.Fatalf("want 10, got %d", n)
	***REMOVED***
***REMOVED***
