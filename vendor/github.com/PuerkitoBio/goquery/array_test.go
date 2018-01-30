package goquery

import (
	"testing"
)

func TestFirst(t *testing.T) ***REMOVED***
	sel := Doc().Find(".pvk-content").First()
	assertLength(t, sel.Nodes, 1)
***REMOVED***

func TestFirstEmpty(t *testing.T) ***REMOVED***
	sel := Doc().Find(".pvk-zzcontentzz").First()
	assertLength(t, sel.Nodes, 0)
***REMOVED***

func TestFirstInvalid(t *testing.T) ***REMOVED***
	sel := Doc().Find("").First()
	assertLength(t, sel.Nodes, 0)
***REMOVED***

func TestFirstRollback(t *testing.T) ***REMOVED***
	sel := Doc().Find(".pvk-content")
	sel2 := sel.First().End()
	assertEqual(t, sel, sel2)
***REMOVED***

func TestLast(t *testing.T) ***REMOVED***
	sel := Doc().Find(".pvk-content").Last()
	assertLength(t, sel.Nodes, 1)

	// Should contain Footer
	foot := Doc().Find(".footer")
	if !sel.Contains(foot.Nodes[0]) ***REMOVED***
		t.Error("Last .pvk-content should contain .footer.")
	***REMOVED***
***REMOVED***

func TestLastEmpty(t *testing.T) ***REMOVED***
	sel := Doc().Find(".pvk-zzcontentzz").Last()
	assertLength(t, sel.Nodes, 0)
***REMOVED***

func TestLastInvalid(t *testing.T) ***REMOVED***
	sel := Doc().Find("").Last()
	assertLength(t, sel.Nodes, 0)
***REMOVED***

func TestLastRollback(t *testing.T) ***REMOVED***
	sel := Doc().Find(".pvk-content")
	sel2 := sel.Last().End()
	assertEqual(t, sel, sel2)
***REMOVED***

func TestEq(t *testing.T) ***REMOVED***
	sel := Doc().Find(".pvk-content").Eq(1)
	assertLength(t, sel.Nodes, 1)
***REMOVED***

func TestEqNegative(t *testing.T) ***REMOVED***
	sel := Doc().Find(".pvk-content").Eq(-1)
	assertLength(t, sel.Nodes, 1)

	// Should contain Footer
	foot := Doc().Find(".footer")
	if !sel.Contains(foot.Nodes[0]) ***REMOVED***
		t.Error("Index -1 of .pvk-content should contain .footer.")
	***REMOVED***
***REMOVED***

func TestEqEmpty(t *testing.T) ***REMOVED***
	sel := Doc().Find("something_random_that_does_not_exists").Eq(0)
	assertLength(t, sel.Nodes, 0)
***REMOVED***

func TestEqInvalid(t *testing.T) ***REMOVED***
	sel := Doc().Find("").Eq(0)
	assertLength(t, sel.Nodes, 0)
***REMOVED***

func TestEqInvalidPositive(t *testing.T) ***REMOVED***
	sel := Doc().Find(".pvk-content").Eq(3)
	assertLength(t, sel.Nodes, 0)
***REMOVED***

func TestEqInvalidNegative(t *testing.T) ***REMOVED***
	sel := Doc().Find(".pvk-content").Eq(-4)
	assertLength(t, sel.Nodes, 0)
***REMOVED***

func TestEqRollback(t *testing.T) ***REMOVED***
	sel := Doc().Find(".pvk-content")
	sel2 := sel.Eq(1).End()
	assertEqual(t, sel, sel2)
***REMOVED***

func TestSlice(t *testing.T) ***REMOVED***
	sel := Doc().Find(".pvk-content").Slice(0, 2)

	assertLength(t, sel.Nodes, 2)
	assertSelectionIs(t, sel, "#pc1", "#pc2")
***REMOVED***

func TestSliceToEnd(t *testing.T) ***REMOVED***
	sel := Doc().Find(".pvk-content").Slice(1, ToEnd)

	assertLength(t, sel.Nodes, 2)
	assertSelectionIs(t, sel.Eq(0), "#pc2")
	if _, ok := sel.Eq(1).Attr("id"); ok ***REMOVED***
		t.Error("Want no attribute ID, got one")
	***REMOVED***
***REMOVED***

func TestSliceEmpty(t *testing.T) ***REMOVED***
	defer assertPanic(t)
	Doc().Find("x").Slice(0, 2)
***REMOVED***

func TestSliceInvalid(t *testing.T) ***REMOVED***
	defer assertPanic(t)
	Doc().Find("").Slice(0, 2)
***REMOVED***

func TestSliceInvalidToEnd(t *testing.T) ***REMOVED***
	defer assertPanic(t)
	Doc().Find("").Slice(2, ToEnd)
***REMOVED***

func TestSliceOutOfBounds(t *testing.T) ***REMOVED***
	defer assertPanic(t)
	Doc().Find(".pvk-content").Slice(2, 12)
***REMOVED***

func TestNegativeSliceStart(t *testing.T) ***REMOVED***
	sel := Doc().Find(".container-fluid").Slice(-2, 3)
	assertLength(t, sel.Nodes, 1)
	assertSelectionIs(t, sel.Eq(0), "#cf3")
***REMOVED***

func TestNegativeSliceEnd(t *testing.T) ***REMOVED***
	sel := Doc().Find(".container-fluid").Slice(1, -1)
	assertLength(t, sel.Nodes, 2)
	assertSelectionIs(t, sel.Eq(0), "#cf2")
	assertSelectionIs(t, sel.Eq(1), "#cf3")
***REMOVED***

func TestNegativeSliceBoth(t *testing.T) ***REMOVED***
	sel := Doc().Find(".container-fluid").Slice(-3, -1)
	assertLength(t, sel.Nodes, 2)
	assertSelectionIs(t, sel.Eq(0), "#cf2")
	assertSelectionIs(t, sel.Eq(1), "#cf3")
***REMOVED***

func TestNegativeSliceToEnd(t *testing.T) ***REMOVED***
	sel := Doc().Find(".container-fluid").Slice(-3, ToEnd)
	assertLength(t, sel.Nodes, 3)
	assertSelectionIs(t, sel, "#cf2", "#cf3", "#cf4")
***REMOVED***

func TestNegativeSliceOutOfBounds(t *testing.T) ***REMOVED***
	defer assertPanic(t)
	Doc().Find(".container-fluid").Slice(-12, -7)
***REMOVED***

func TestSliceRollback(t *testing.T) ***REMOVED***
	sel := Doc().Find(".pvk-content")
	sel2 := sel.Slice(0, 2).End()
	assertEqual(t, sel, sel2)
***REMOVED***

func TestGet(t *testing.T) ***REMOVED***
	sel := Doc().Find(".pvk-content")
	node := sel.Get(1)
	if sel.Nodes[1] != node ***REMOVED***
		t.Errorf("Expected node %v to be %v.", node, sel.Nodes[1])
	***REMOVED***
***REMOVED***

func TestGetNegative(t *testing.T) ***REMOVED***
	sel := Doc().Find(".pvk-content")
	node := sel.Get(-3)
	if sel.Nodes[0] != node ***REMOVED***
		t.Errorf("Expected node %v to be %v.", node, sel.Nodes[0])
	***REMOVED***
***REMOVED***

func TestGetInvalid(t *testing.T) ***REMOVED***
	defer assertPanic(t)
	sel := Doc().Find(".pvk-content")
	sel.Get(129)
***REMOVED***

func TestIndex(t *testing.T) ***REMOVED***
	sel := Doc().Find(".pvk-content")
	if i := sel.Index(); i != 1 ***REMOVED***
		t.Errorf("Expected index of 1, got %v.", i)
	***REMOVED***
***REMOVED***

func TestIndexSelector(t *testing.T) ***REMOVED***
	sel := Doc().Find(".hero-unit")
	if i := sel.IndexSelector("div"); i != 4 ***REMOVED***
		t.Errorf("Expected index of 4, got %v.", i)
	***REMOVED***
***REMOVED***

func TestIndexSelectorInvalid(t *testing.T) ***REMOVED***
	sel := Doc().Find(".hero-unit")
	if i := sel.IndexSelector(""); i != -1 ***REMOVED***
		t.Errorf("Expected index of -1, got %v.", i)
	***REMOVED***
***REMOVED***

func TestIndexOfNode(t *testing.T) ***REMOVED***
	sel := Doc().Find("div.pvk-gutter")
	if i := sel.IndexOfNode(sel.Nodes[1]); i != 1 ***REMOVED***
		t.Errorf("Expected index of 1, got %v.", i)
	***REMOVED***
***REMOVED***

func TestIndexOfNilNode(t *testing.T) ***REMOVED***
	sel := Doc().Find("div.pvk-gutter")
	if i := sel.IndexOfNode(nil); i != -1 ***REMOVED***
		t.Errorf("Expected index of -1, got %v.", i)
	***REMOVED***
***REMOVED***

func TestIndexOfSelection(t *testing.T) ***REMOVED***
	sel := Doc().Find("div")
	sel2 := Doc().Find(".hero-unit")
	if i := sel.IndexOfSelection(sel2); i != 4 ***REMOVED***
		t.Errorf("Expected index of 4, got %v.", i)
	***REMOVED***
***REMOVED***
