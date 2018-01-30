package goquery

import (
	"testing"
)

func BenchmarkFind(b *testing.B) ***REMOVED***
	var n int

	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = DocB().Find("dd").Length()

		***REMOVED*** else ***REMOVED***
			DocB().Find("dd")
		***REMOVED***
	***REMOVED***
	if n != 41 ***REMOVED***
		b.Fatalf("want 41, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkFindWithinSelection(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocW().Find("ul")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.Find("a[class]").Length()
		***REMOVED*** else ***REMOVED***
			sel.Find("a[class]")
		***REMOVED***
	***REMOVED***
	if n != 39 ***REMOVED***
		b.Fatalf("want 39, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkFindSelection(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocW().Find("ul")
	sel2 := DocW().Find("span")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.FindSelection(sel2).Length()
		***REMOVED*** else ***REMOVED***
			sel.FindSelection(sel2)
		***REMOVED***
	***REMOVED***
	if n != 73 ***REMOVED***
		b.Fatalf("want 73, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkFindNodes(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocW().Find("ul")
	sel2 := DocW().Find("span")
	nodes := sel2.Nodes
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.FindNodes(nodes...).Length()
		***REMOVED*** else ***REMOVED***
			sel.FindNodes(nodes...)
		***REMOVED***
	***REMOVED***
	if n != 73 ***REMOVED***
		b.Fatalf("want 73, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkContents(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocW().Find(".toclevel-1")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.Contents().Length()
		***REMOVED*** else ***REMOVED***
			sel.Contents()
		***REMOVED***
	***REMOVED***
	if n != 16 ***REMOVED***
		b.Fatalf("want 16, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkContentsFiltered(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocW().Find(".toclevel-1")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.ContentsFiltered("a[href=\"#Examples\"]").Length()
		***REMOVED*** else ***REMOVED***
			sel.ContentsFiltered("a[href=\"#Examples\"]")
		***REMOVED***
	***REMOVED***
	if n != 1 ***REMOVED***
		b.Fatalf("want 1, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkChildren(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocW().Find(".toclevel-2")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.Children().Length()
		***REMOVED*** else ***REMOVED***
			sel.Children()
		***REMOVED***
	***REMOVED***
	if n != 2 ***REMOVED***
		b.Fatalf("want 2, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkChildrenFiltered(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocW().Find("h3")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.ChildrenFiltered(".editsection").Length()
		***REMOVED*** else ***REMOVED***
			sel.ChildrenFiltered(".editsection")
		***REMOVED***
	***REMOVED***
	if n != 2 ***REMOVED***
		b.Fatalf("want 2, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkParent(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocW().Find("li")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.Parent().Length()
		***REMOVED*** else ***REMOVED***
			sel.Parent()
		***REMOVED***
	***REMOVED***
	if n != 55 ***REMOVED***
		b.Fatalf("want 55, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkParentFiltered(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocW().Find("li")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.ParentFiltered("ul[id]").Length()
		***REMOVED*** else ***REMOVED***
			sel.ParentFiltered("ul[id]")
		***REMOVED***
	***REMOVED***
	if n != 4 ***REMOVED***
		b.Fatalf("want 4, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkParents(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocW().Find("th a")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.Parents().Length()
		***REMOVED*** else ***REMOVED***
			sel.Parents()
		***REMOVED***
	***REMOVED***
	if n != 73 ***REMOVED***
		b.Fatalf("want 73, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkParentsFiltered(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocW().Find("th a")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.ParentsFiltered("tr").Length()
		***REMOVED*** else ***REMOVED***
			sel.ParentsFiltered("tr")
		***REMOVED***
	***REMOVED***
	if n != 18 ***REMOVED***
		b.Fatalf("want 18, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkParentsUntil(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocW().Find("th a")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.ParentsUntil("table").Length()
		***REMOVED*** else ***REMOVED***
			sel.ParentsUntil("table")
		***REMOVED***
	***REMOVED***
	if n != 52 ***REMOVED***
		b.Fatalf("want 52, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkParentsUntilSelection(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocW().Find("th a")
	sel2 := DocW().Find("#content")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.ParentsUntilSelection(sel2).Length()
		***REMOVED*** else ***REMOVED***
			sel.ParentsUntilSelection(sel2)
		***REMOVED***
	***REMOVED***
	if n != 70 ***REMOVED***
		b.Fatalf("want 70, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkParentsUntilNodes(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocW().Find("th a")
	sel2 := DocW().Find("#content")
	nodes := sel2.Nodes
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.ParentsUntilNodes(nodes...).Length()
		***REMOVED*** else ***REMOVED***
			sel.ParentsUntilNodes(nodes...)
		***REMOVED***
	***REMOVED***
	if n != 70 ***REMOVED***
		b.Fatalf("want 70, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkParentsFilteredUntil(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocW().Find(".toclevel-1 a")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.ParentsFilteredUntil(":nth-child(1)", "ul").Length()
		***REMOVED*** else ***REMOVED***
			sel.ParentsFilteredUntil(":nth-child(1)", "ul")
		***REMOVED***
	***REMOVED***
	if n != 2 ***REMOVED***
		b.Fatalf("want 2, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkParentsFilteredUntilSelection(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocW().Find(".toclevel-1 a")
	sel2 := DocW().Find("ul")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.ParentsFilteredUntilSelection(":nth-child(1)", sel2).Length()
		***REMOVED*** else ***REMOVED***
			sel.ParentsFilteredUntilSelection(":nth-child(1)", sel2)
		***REMOVED***
	***REMOVED***
	if n != 2 ***REMOVED***
		b.Fatalf("want 2, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkParentsFilteredUntilNodes(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocW().Find(".toclevel-1 a")
	sel2 := DocW().Find("ul")
	nodes := sel2.Nodes
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.ParentsFilteredUntilNodes(":nth-child(1)", nodes...).Length()
		***REMOVED*** else ***REMOVED***
			sel.ParentsFilteredUntilNodes(":nth-child(1)", nodes...)
		***REMOVED***
	***REMOVED***
	if n != 2 ***REMOVED***
		b.Fatalf("want 2, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkSiblings(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocW().Find("ul li:nth-child(1)")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.Siblings().Length()
		***REMOVED*** else ***REMOVED***
			sel.Siblings()
		***REMOVED***
	***REMOVED***
	if n != 293 ***REMOVED***
		b.Fatalf("want 293, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkSiblingsFiltered(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocW().Find("ul li:nth-child(1)")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.SiblingsFiltered("[class]").Length()
		***REMOVED*** else ***REMOVED***
			sel.SiblingsFiltered("[class]")
		***REMOVED***
	***REMOVED***
	if n != 46 ***REMOVED***
		b.Fatalf("want 46, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkNext(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocW().Find("li:nth-child(1)")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.Next().Length()
		***REMOVED*** else ***REMOVED***
			sel.Next()
		***REMOVED***
	***REMOVED***
	if n != 49 ***REMOVED***
		b.Fatalf("want 49, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkNextFiltered(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocW().Find("li:nth-child(1)")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.NextFiltered("[class]").Length()
		***REMOVED*** else ***REMOVED***
			sel.NextFiltered("[class]")
		***REMOVED***
	***REMOVED***
	if n != 6 ***REMOVED***
		b.Fatalf("want 6, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkNextAll(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocW().Find("li:nth-child(3)")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.NextAll().Length()
		***REMOVED*** else ***REMOVED***
			sel.NextAll()
		***REMOVED***
	***REMOVED***
	if n != 234 ***REMOVED***
		b.Fatalf("want 234, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkNextAllFiltered(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocW().Find("li:nth-child(3)")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.NextAllFiltered("[class]").Length()
		***REMOVED*** else ***REMOVED***
			sel.NextAllFiltered("[class]")
		***REMOVED***
	***REMOVED***
	if n != 33 ***REMOVED***
		b.Fatalf("want 33, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkPrev(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocW().Find("li:last-child")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.Prev().Length()
		***REMOVED*** else ***REMOVED***
			sel.Prev()
		***REMOVED***
	***REMOVED***
	if n != 49 ***REMOVED***
		b.Fatalf("want 49, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkPrevFiltered(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocW().Find("li:last-child")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.PrevFiltered("[class]").Length()
		***REMOVED*** else ***REMOVED***
			sel.PrevFiltered("[class]")
		***REMOVED***
	***REMOVED***
	// There is one more Prev li with a class, compared to Next li with a class
	// (confirmed by looking at the HTML, this is ok)
	if n != 7 ***REMOVED***
		b.Fatalf("want 7, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkPrevAll(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocW().Find("li:nth-child(4)")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.PrevAll().Length()
		***REMOVED*** else ***REMOVED***
			sel.PrevAll()
		***REMOVED***
	***REMOVED***
	if n != 78 ***REMOVED***
		b.Fatalf("want 78, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkPrevAllFiltered(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocW().Find("li:nth-child(4)")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.PrevAllFiltered("[class]").Length()
		***REMOVED*** else ***REMOVED***
			sel.PrevAllFiltered("[class]")
		***REMOVED***
	***REMOVED***
	if n != 6 ***REMOVED***
		b.Fatalf("want 6, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkNextUntil(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocW().Find("li:first-child")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.NextUntil(":nth-child(4)").Length()
		***REMOVED*** else ***REMOVED***
			sel.NextUntil(":nth-child(4)")
		***REMOVED***
	***REMOVED***
	if n != 84 ***REMOVED***
		b.Fatalf("want 84, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkNextUntilSelection(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocW().Find("h2")
	sel2 := DocW().Find("ul")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.NextUntilSelection(sel2).Length()
		***REMOVED*** else ***REMOVED***
			sel.NextUntilSelection(sel2)
		***REMOVED***
	***REMOVED***
	if n != 42 ***REMOVED***
		b.Fatalf("want 42, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkNextUntilNodes(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocW().Find("h2")
	sel2 := DocW().Find("p")
	nodes := sel2.Nodes
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.NextUntilNodes(nodes...).Length()
		***REMOVED*** else ***REMOVED***
			sel.NextUntilNodes(nodes...)
		***REMOVED***
	***REMOVED***
	if n != 12 ***REMOVED***
		b.Fatalf("want 12, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkPrevUntil(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocW().Find("li:last-child")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.PrevUntil(":nth-child(4)").Length()
		***REMOVED*** else ***REMOVED***
			sel.PrevUntil(":nth-child(4)")
		***REMOVED***
	***REMOVED***
	if n != 238 ***REMOVED***
		b.Fatalf("want 238, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkPrevUntilSelection(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocW().Find("h2")
	sel2 := DocW().Find("ul")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.PrevUntilSelection(sel2).Length()
		***REMOVED*** else ***REMOVED***
			sel.PrevUntilSelection(sel2)
		***REMOVED***
	***REMOVED***
	if n != 49 ***REMOVED***
		b.Fatalf("want 49, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkPrevUntilNodes(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocW().Find("h2")
	sel2 := DocW().Find("p")
	nodes := sel2.Nodes
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.PrevUntilNodes(nodes...).Length()
		***REMOVED*** else ***REMOVED***
			sel.PrevUntilNodes(nodes...)
		***REMOVED***
	***REMOVED***
	if n != 11 ***REMOVED***
		b.Fatalf("want 11, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkNextFilteredUntil(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocW().Find("h2")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.NextFilteredUntil("p", "div").Length()
		***REMOVED*** else ***REMOVED***
			sel.NextFilteredUntil("p", "div")
		***REMOVED***
	***REMOVED***
	if n != 22 ***REMOVED***
		b.Fatalf("want 22, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkNextFilteredUntilSelection(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocW().Find("h2")
	sel2 := DocW().Find("div")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.NextFilteredUntilSelection("p", sel2).Length()
		***REMOVED*** else ***REMOVED***
			sel.NextFilteredUntilSelection("p", sel2)
		***REMOVED***
	***REMOVED***
	if n != 22 ***REMOVED***
		b.Fatalf("want 22, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkNextFilteredUntilNodes(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocW().Find("h2")
	sel2 := DocW().Find("div")
	nodes := sel2.Nodes
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.NextFilteredUntilNodes("p", nodes...).Length()
		***REMOVED*** else ***REMOVED***
			sel.NextFilteredUntilNodes("p", nodes...)
		***REMOVED***
	***REMOVED***
	if n != 22 ***REMOVED***
		b.Fatalf("want 22, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkPrevFilteredUntil(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocW().Find("h2")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.PrevFilteredUntil("p", "div").Length()
		***REMOVED*** else ***REMOVED***
			sel.PrevFilteredUntil("p", "div")
		***REMOVED***
	***REMOVED***
	if n != 20 ***REMOVED***
		b.Fatalf("want 20, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkPrevFilteredUntilSelection(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocW().Find("h2")
	sel2 := DocW().Find("div")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.PrevFilteredUntilSelection("p", sel2).Length()
		***REMOVED*** else ***REMOVED***
			sel.PrevFilteredUntilSelection("p", sel2)
		***REMOVED***
	***REMOVED***
	if n != 20 ***REMOVED***
		b.Fatalf("want 20, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkPrevFilteredUntilNodes(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := DocW().Find("h2")
	sel2 := DocW().Find("div")
	nodes := sel2.Nodes
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.PrevFilteredUntilNodes("p", nodes...).Length()
		***REMOVED*** else ***REMOVED***
			sel.PrevFilteredUntilNodes("p", nodes...)
		***REMOVED***
	***REMOVED***
	if n != 20 ***REMOVED***
		b.Fatalf("want 20, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkClosest(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := Doc().Find(".container-fluid")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.Closest(".pvk-content").Length()
		***REMOVED*** else ***REMOVED***
			sel.Closest(".pvk-content")
		***REMOVED***
	***REMOVED***
	if n != 2 ***REMOVED***
		b.Fatalf("want 2, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkClosestSelection(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := Doc().Find(".container-fluid")
	sel2 := Doc().Find(".pvk-content")
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.ClosestSelection(sel2).Length()
		***REMOVED*** else ***REMOVED***
			sel.ClosestSelection(sel2)
		***REMOVED***
	***REMOVED***
	if n != 2 ***REMOVED***
		b.Fatalf("want 2, got %d", n)
	***REMOVED***
***REMOVED***

func BenchmarkClosestNodes(b *testing.B) ***REMOVED***
	var n int

	b.StopTimer()
	sel := Doc().Find(".container-fluid")
	nodes := Doc().Find(".pvk-content").Nodes
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if n == 0 ***REMOVED***
			n = sel.ClosestNodes(nodes...).Length()
		***REMOVED*** else ***REMOVED***
			sel.ClosestNodes(nodes...)
		***REMOVED***
	***REMOVED***
	if n != 2 ***REMOVED***
		b.Fatalf("want 2, got %d", n)
	***REMOVED***
***REMOVED***
