package goquery

import (
	"testing"
)

func BenchmarkAdd(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocB().Find("dd")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.Add("h2[title]").Length()
		***REMOVED*** else ***REMOVED***
			sel.Add("h2[title]")
		***REMOVED***
	***REMOVED***
	if n != 43 ***REMOVED***
		b.Fatalf("want 43, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkAddSelection(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocB().Find("dd")
	sel2 := DocB().Find("h2[title]")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.AddSelection(sel2).Length()
		***REMOVED*** else ***REMOVED***
			sel.AddSelection(sel2)
		***REMOVED***
	***REMOVED***
	if n != 43 ***REMOVED***
		b.Fatalf("want 43, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkAddNodes(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocB().Find("dd")
	sel2 := DocB().Find("h2[title]")
	nodes := sel2.Nodes
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.AddNodes(nodes...).Length()
		***REMOVED*** else ***REMOVED***
			sel.AddNodes(nodes...)
		***REMOVED***
	***REMOVED***
	if n != 43 ***REMOVED***
		b.Fatalf("want 43, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkAddNodesBig(b *testing.B) ***REMOVED***
	var n int

	doc := DocW()
	sel := doc.Find("li")
	// make nodes > 1000
	nodes := sel.Nodes
	nodes = append(nodes, nodes...)
	nodes = append(nodes, nodes...)
	sel = doc.Find("xyz")
	b.ResetTimer()

	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.AddNodes(nodes...).Length()
		***REMOVED*** else ***REMOVED***
			sel.AddNodes(nodes...)
		***REMOVED***
	***REMOVED***
	if n != 373 ***REMOVED***
		b.Fatalf("want 373, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkAndSelf(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocB().Find("dd").Parent()
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.AndSelf().Length()
		***REMOVED*** else ***REMOVED***
			sel.AndSelf()
		***REMOVED***
	***REMOVED***
	if n != 44 ***REMOVED***
		b.Fatalf("want 44, got %d", n)
	***REMOVED***
***REMOVED***
