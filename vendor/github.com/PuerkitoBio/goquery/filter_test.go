package goquery

import (
	"testing"
)

func TestFilter(t *testing.T) ***REMOVED***
	sel := Doc().Find(".span12").Filter(".alert")
	assertLength(t, sel.Nodes, 1)
***REMOVED***

func TestFilterNone(t *testing.T) ***REMOVED***
	sel := Doc().Find(".span12").Filter(".zzalert")
	assertLength(t, sel.Nodes, 0)
***REMOVED***

func TestFilterInvalid(t *testing.T) ***REMOVED***
	sel := Doc().Find(".span12").Filter("")
	assertLength(t, sel.Nodes, 0)
***REMOVED***

func TestFilterRollback(t *testing.T) ***REMOVED***
	sel := Doc().Find(".pvk-content")
	sel2 := sel.Filter(".alert").End()
	assertEqual(t, sel, sel2)
***REMOVED***

func TestFilterFunction(t *testing.T) ***REMOVED***
	sel := Doc().Find(".pvk-content").FilterFunction(func(i int, s *Selection) bool ***REMOVED***
		return i > 0
	***REMOVED***)
	assertLength(t, sel.Nodes, 2)
***REMOVED***

func TestFilterFunctionRollback(t *testing.T) ***REMOVED***
	sel := Doc().Find(".pvk-content")
	sel2 := sel.FilterFunction(func(i int, s *Selection) bool ***REMOVED***
		return i > 0
	***REMOVED***).End()
	assertEqual(t, sel, sel2)
***REMOVED***

func TestFilterNode(t *testing.T) ***REMOVED***
	sel := Doc().Find(".pvk-content")
	sel2 := sel.FilterNodes(sel.Nodes[2])
	assertLength(t, sel2.Nodes, 1)
***REMOVED***

func TestFilterNodeRollback(t *testing.T) ***REMOVED***
	sel := Doc().Find(".pvk-content")
	sel2 := sel.FilterNodes(sel.Nodes[2]).End()
	assertEqual(t, sel, sel2)
***REMOVED***

func TestFilterSelection(t *testing.T) ***REMOVED***
	sel := Doc().Find(".link")
	sel2 := Doc().Find("a[ng-click]")
	sel3 := sel.FilterSelection(sel2)
	assertLength(t, sel3.Nodes, 1)
***REMOVED***

func TestFilterSelectionRollback(t *testing.T) ***REMOVED***
	sel := Doc().Find(".link")
	sel2 := Doc().Find("a[ng-click]")
	sel2 = sel.FilterSelection(sel2).End()
	assertEqual(t, sel, sel2)
***REMOVED***

func TestFilterSelectionNil(t *testing.T) ***REMOVED***
	var sel2 *Selection

	sel := Doc().Find(".link")
	sel3 := sel.FilterSelection(sel2)
	assertLength(t, sel3.Nodes, 0)
***REMOVED***

func TestNot(t *testing.T) ***REMOVED***
	sel := Doc().Find(".span12").Not(".alert")
	assertLength(t, sel.Nodes, 1)
***REMOVED***

func TestNotInvalid(t *testing.T) ***REMOVED***
	sel := Doc().Find(".span12").Not("")
	assertLength(t, sel.Nodes, 2)
***REMOVED***

func TestNotRollback(t *testing.T) ***REMOVED***
	sel := Doc().Find(".span12")
	sel2 := sel.Not(".alert").End()
	assertEqual(t, sel, sel2)
***REMOVED***

func TestNotNone(t *testing.T) ***REMOVED***
	sel := Doc().Find(".span12").Not(".zzalert")
	assertLength(t, sel.Nodes, 2)
***REMOVED***

func TestNotFunction(t *testing.T) ***REMOVED***
	sel := Doc().Find(".pvk-content").NotFunction(func(i int, s *Selection) bool ***REMOVED***
		return i > 0
	***REMOVED***)
	assertLength(t, sel.Nodes, 1)
***REMOVED***

func TestNotFunctionRollback(t *testing.T) ***REMOVED***
	sel := Doc().Find(".pvk-content")
	sel2 := sel.NotFunction(func(i int, s *Selection) bool ***REMOVED***
		return i > 0
	***REMOVED***).End()
	assertEqual(t, sel, sel2)
***REMOVED***

func TestNotNode(t *testing.T) ***REMOVED***
	sel := Doc().Find(".pvk-content")
	sel2 := sel.NotNodes(sel.Nodes[2])
	assertLength(t, sel2.Nodes, 2)
***REMOVED***

func TestNotNodeRollback(t *testing.T) ***REMOVED***
	sel := Doc().Find(".pvk-content")
	sel2 := sel.NotNodes(sel.Nodes[2]).End()
	assertEqual(t, sel, sel2)
***REMOVED***

func TestNotSelection(t *testing.T) ***REMOVED***
	sel := Doc().Find(".link")
	sel2 := Doc().Find("a[ng-click]")
	sel3 := sel.NotSelection(sel2)
	assertLength(t, sel3.Nodes, 6)
***REMOVED***

func TestNotSelectionRollback(t *testing.T) ***REMOVED***
	sel := Doc().Find(".link")
	sel2 := Doc().Find("a[ng-click]")
	sel2 = sel.NotSelection(sel2).End()
	assertEqual(t, sel, sel2)
***REMOVED***

func TestIntersection(t *testing.T) ***REMOVED***
	sel := Doc().Find(".pvk-gutter")
	sel2 := Doc().Find("div").Intersection(sel)
	assertLength(t, sel2.Nodes, 6)
***REMOVED***

func TestIntersectionRollback(t *testing.T) ***REMOVED***
	sel := Doc().Find(".pvk-gutter")
	sel2 := Doc().Find("div")
	sel2 = sel.Intersection(sel2).End()
	assertEqual(t, sel, sel2)
***REMOVED***

func TestHas(t *testing.T) ***REMOVED***
	sel := Doc().Find(".container-fluid").Has(".center-content")
	assertLength(t, sel.Nodes, 2)
	// Has() returns the high-level .container-fluid div, and the one that is the immediate parent of center-content
***REMOVED***

func TestHasInvalid(t *testing.T) ***REMOVED***
	sel := Doc().Find(".container-fluid").Has("")
	assertLength(t, sel.Nodes, 0)
***REMOVED***

func TestHasRollback(t *testing.T) ***REMOVED***
	sel := Doc().Find(".container-fluid")
	sel2 := sel.Has(".center-content").End()
	assertEqual(t, sel, sel2)
***REMOVED***

func TestHasNodes(t *testing.T) ***REMOVED***
	sel := Doc().Find(".container-fluid")
	sel2 := Doc().Find(".center-content")
	sel = sel.HasNodes(sel2.Nodes...)
	assertLength(t, sel.Nodes, 2)
	// Has() returns the high-level .container-fluid div, and the one that is the immediate parent of center-content
***REMOVED***

func TestHasNodesRollback(t *testing.T) ***REMOVED***
	sel := Doc().Find(".container-fluid")
	sel2 := Doc().Find(".center-content")
	sel2 = sel.HasNodes(sel2.Nodes...).End()
	assertEqual(t, sel, sel2)
***REMOVED***

func TestHasSelection(t *testing.T) ***REMOVED***
	sel := Doc().Find("p")
	sel2 := Doc().Find("small")
	sel = sel.HasSelection(sel2)
	assertLength(t, sel.Nodes, 1)
***REMOVED***

func TestHasSelectionRollback(t *testing.T) ***REMOVED***
	sel := Doc().Find("p")
	sel2 := Doc().Find("small")
	sel2 = sel.HasSelection(sel2).End()
	assertEqual(t, sel, sel2)
***REMOVED***

func TestEnd(t *testing.T) ***REMOVED***
	sel := Doc().Find("p").Has("small").End()
	assertLength(t, sel.Nodes, 4)
***REMOVED***

func TestEndToTop(t *testing.T) ***REMOVED***
	sel := Doc().Find("p").Has("small").End().End().End()
	assertLength(t, sel.Nodes, 0)
***REMOVED***
