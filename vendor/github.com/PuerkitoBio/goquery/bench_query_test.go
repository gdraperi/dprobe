package goquery

import (
	"testing"
)

func BenchmarkIs(b *testing.B) ***REMOVED***
	var y bool

	b.StopTimer()
	sel := DocW().Find("li")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		y = sel.Is(".toclevel-2")
	***REMOVED***
	if !y ***REMOVED***
		b.Fatal("want true")
	***REMOVED***
***REMOVED***

func BenchmarkIsPositional(b *testing.B) ***REMOVED***
	var y bool

	b.StopTimer()
	sel := DocW().Find("li")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		y = sel.Is("li:nth-child(2)")
	***REMOVED***
	if !y ***REMOVED***
		b.Fatal("want true")
	***REMOVED***
***REMOVED***

func BenchmarkIsFunction(b *testing.B) ***REMOVED***
	var y bool

	b.StopTimer()
	sel := DocW().Find(".toclevel-1")
	f := func(i int, s *Selection) bool ***REMOVED***
		return i == 8
	***REMOVED***
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		y = sel.IsFunction(f)
	***REMOVED***
	if !y ***REMOVED***
		b.Fatal("want true")
	***REMOVED***
***REMOVED***

func BenchmarkIsSelection(b *testing.B) ***REMOVED***
	var y bool

	b.StopTimer()
	sel := DocW().Find("li")
	sel2 := DocW().Find(".toclevel-2")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		y = sel.IsSelection(sel2)
	***REMOVED***
	if !y ***REMOVED***
		b.Fatal("want true")
	***REMOVED***
***REMOVED***

func BenchmarkIsNodes(b *testing.B) ***REMOVED***
	var y bool

	b.StopTimer()
	sel := DocW().Find("li")
	sel2 := DocW().Find(".toclevel-2")
	nodes := sel2.Nodes
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		y = sel.IsNodes(nodes...)
	***REMOVED***
	if !y ***REMOVED***
		b.Fatal("want true")
	***REMOVED***
***REMOVED***

func BenchmarkHasClass(b *testing.B) ***REMOVED***
	var y bool

	b.StopTimer()
	sel := DocW().Find("span")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		y = sel.HasClass("official")
	***REMOVED***
	if !y ***REMOVED***
		b.Fatal("want true")
	***REMOVED***
***REMOVED***

func BenchmarkContains(b *testing.B) ***REMOVED***
	var y bool

	b.StopTimer()
	sel := DocW().Find("span.url")
	sel2 := DocW().Find("a[rel=\"nofollow\"]")
	node := sel2.Nodes[0]
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		y = sel.Contains(node)
	***REMOVED***
	if !y ***REMOVED***
		b.Fatal("want true")
	***REMOVED***
***REMOVED***
