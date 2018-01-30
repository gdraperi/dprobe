package goquery

import (
	"testing"
)

func BenchmarkFirst(b *testing.B) ***REMOVED***
	b.StopTimer()
	sel := DocB().Find("dd")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		sel.First()
	***REMOVED***
***REMOVED***

func BenchmarkLast(b *testing.B) ***REMOVED***
	b.StopTimer()
	sel := DocB().Find("dd")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		sel.Last()
	***REMOVED***
***REMOVED***

func BenchmarkEq(b *testing.B) ***REMOVED***
	b.StopTimer()
	sel := DocB().Find("dd")
	j := 0
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		sel.Eq(j)
		if j++; j >= sel.Length() ***REMOVED***
			j = 0
		***REMOVED***
	***REMOVED***
***REMOVED***

func BenchmarkSlice(b *testing.B) ***REMOVED***
	b.StopTimer()
	sel := DocB().Find("dd")
	j := 0
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		sel.Slice(j, j+4)
		if j++; j >= (sel.Length() - 4) ***REMOVED***
			j = 0
		***REMOVED***
	***REMOVED***
***REMOVED***

func BenchmarkGet(b *testing.B) ***REMOVED***
	b.StopTimer()
	sel := DocB().Find("dd")
	j := 0
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		sel.Get(j)
		if j++; j >= sel.Length() ***REMOVED***
			j = 0
		***REMOVED***
	***REMOVED***
***REMOVED***

func BenchmarkIndex(b *testing.B) ***REMOVED***
	var j int

	b.StopTimer()
	sel := DocB().Find("#Main")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		j = sel.Index()
	***REMOVED***
	if j != 3 ***REMOVED***
		b.Fatalf("want 3, got %d", j)
	***REMOVED***
***REMOVED***

func BenchmarkIndexSelector(b *testing.B) ***REMOVED***
	var j int

	b.StopTimer()
	sel := DocB().Find("#manual-nav dl dd:nth-child(1)")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		j = sel.IndexSelector("dd")
	***REMOVED***
	if j != 4 ***REMOVED***
		b.Fatalf("want 4, got %d", j)
	***REMOVED***
***REMOVED***

func BenchmarkIndexOfNode(b *testing.B) ***REMOVED***
	var j int

	b.StopTimer()
	sel := DocB().Find("span a")
	sel2 := DocB().Find("span a:nth-child(3)")
	n := sel2.Get(0)
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		j = sel.IndexOfNode(n)
	***REMOVED***
	if j != 2 ***REMOVED***
		b.Fatalf("want 2, got %d", j)
	***REMOVED***
***REMOVED***

func BenchmarkIndexOfSelection(b *testing.B) ***REMOVED***
	var j int
	b.StopTimer()
	sel := DocB().Find("span a")
	sel2 := DocB().Find("span a:nth-child(3)")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		j = sel.IndexOfSelection(sel2)
	***REMOVED***
	if j != 2 ***REMOVED***
		b.Fatalf("want 2, got %d", j)
	***REMOVED***
***REMOVED***
