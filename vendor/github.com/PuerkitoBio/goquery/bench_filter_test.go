package goquery

import (
	"testing"
)

func BenchmarkFilter(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocW().Find("li")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.Filter(".toclevel-1").Length()
		***REMOVED*** else ***REMOVED***
			sel.Filter(".toclevel-1")
		***REMOVED***
	***REMOVED***
	if n != 13 ***REMOVED***
		b.Fatalf("want 13, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkNot(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocW().Find("li")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.Not(".toclevel-2").Length()
		***REMOVED*** else ***REMOVED***
			sel.Filter(".toclevel-2")
		***REMOVED***
	***REMOVED***
	if n != 371 ***REMOVED***
		b.Fatalf("want 371, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkFilterFunction(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocW().Find("li")
	f := func(i int, s *Selection) bool ***REMOVED***
		return len(s.Get(0).Attr) > 0
	***REMOVED***
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.FilterFunction(f).Length()
		***REMOVED*** else ***REMOVED***
			sel.FilterFunction(f)
		***REMOVED***
	***REMOVED***
	if n != 112 ***REMOVED***
		b.Fatalf("want 112, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkNotFunction(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocW().Find("li")
	f := func(i int, s *Selection) bool ***REMOVED***
		return len(s.Get(0).Attr) > 0
	***REMOVED***
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.NotFunction(f).Length()
		***REMOVED*** else ***REMOVED***
			sel.NotFunction(f)
		***REMOVED***
	***REMOVED***
	if n != 261 ***REMOVED***
		b.Fatalf("want 261, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkFilterNodes(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocW().Find("li")
	sel2 := DocW().Find(".toclevel-2")
	nodes := sel2.Nodes
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.FilterNodes(nodes...).Length()
		***REMOVED*** else ***REMOVED***
			sel.FilterNodes(nodes...)
		***REMOVED***
	***REMOVED***
	if n != 2 ***REMOVED***
		b.Fatalf("want 2, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkNotNodes(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocW().Find("li")
	sel2 := DocW().Find(".toclevel-1")
	nodes := sel2.Nodes
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.NotNodes(nodes...).Length()
		***REMOVED*** else ***REMOVED***
			sel.NotNodes(nodes...)
		***REMOVED***
	***REMOVED***
	if n != 360 ***REMOVED***
		b.Fatalf("want 360, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkFilterSelection(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocW().Find("li")
	sel2 := DocW().Find(".toclevel-2")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.FilterSelection(sel2).Length()
		***REMOVED*** else ***REMOVED***
			sel.FilterSelection(sel2)
		***REMOVED***
	***REMOVED***
	if n != 2 ***REMOVED***
		b.Fatalf("want 2, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkNotSelection(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocW().Find("li")
	sel2 := DocW().Find(".toclevel-1")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.NotSelection(sel2).Length()
		***REMOVED*** else ***REMOVED***
			sel.NotSelection(sel2)
		***REMOVED***
	***REMOVED***
	if n != 360 ***REMOVED***
		b.Fatalf("want 360, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkHas(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocW().Find("h2")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.Has(".editsection").Length()
		***REMOVED*** else ***REMOVED***
			sel.Has(".editsection")
		***REMOVED***
	***REMOVED***
	if n != 13 ***REMOVED***
		b.Fatalf("want 13, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkHasNodes(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocW().Find("li")
	sel2 := DocW().Find(".tocnumber")
	nodes := sel2.Nodes
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.HasNodes(nodes...).Length()
		***REMOVED*** else ***REMOVED***
			sel.HasNodes(nodes...)
		***REMOVED***
	***REMOVED***
	if n != 15 ***REMOVED***
		b.Fatalf("want 15, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkHasSelection(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocW().Find("li")
	sel2 := DocW().Find(".tocnumber")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.HasSelection(sel2).Length()
		***REMOVED*** else ***REMOVED***
			sel.HasSelection(sel2)
		***REMOVED***
	***REMOVED***
	if n != 15 ***REMOVED***
		b.Fatalf("want 15, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkEnd(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocW().Find("li").Has(".tocnumber")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.End().Length()
		***REMOVED*** else ***REMOVED***
			sel.End()
		***REMOVED***
	***REMOVED***
	if n != 373 ***REMOVED***
		b.Fatalf("want 373, got %d", n)
	***REMOVED***
***REMOVED***
