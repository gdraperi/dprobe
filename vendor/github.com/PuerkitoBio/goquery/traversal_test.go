package goquery

import (
	"strings"
	"testing"
)

func TestFind(t *testing.T) ***REMOVED***
	sel := Doc().Find("div.row-fluid")
	assertLength(t, sel.Nodes, 9)
***REMOVED***

func TestFindRollback(t *testing.T) ***REMOVED***
	sel := Doc().Find("div.row-fluid")
	sel2 := sel.Find("a").End()
	assertEqual(t, sel, sel2)
***REMOVED***

func TestFindNotSelf(t *testing.T) ***REMOVED***
	sel := Doc().Find("h1").Find("h1")
	assertLength(t, sel.Nodes, 0)
***REMOVED***

func TestFindInvalid(t *testing.T) ***REMOVED***
	sel := Doc().Find(":+ ^")
	assertLength(t, sel.Nodes, 0)
***REMOVED***

func TestFindBig(t *testing.T) ***REMOVED***
	doc := DocW()
	sel := doc.Find("li")
	assertLength(t, sel.Nodes, 373)
	sel2 := doc.Find("span")
	assertLength(t, sel2.Nodes, 448)
	sel3 := sel.FindSelection(sel2)
	assertLength(t, sel3.Nodes, 248)
***REMOVED***

func TestChainedFind(t *testing.T) ***REMOVED***
	sel := Doc().Find("div.hero-unit").Find(".row-fluid")
	assertLength(t, sel.Nodes, 4)
***REMOVED***

func TestChainedFindInvalid(t *testing.T) ***REMOVED***
	sel := Doc().Find("div.hero-unit").Find("")
	assertLength(t, sel.Nodes, 0)
***REMOVED***

func TestChildren(t *testing.T) ***REMOVED***
	sel := Doc().Find(".pvk-content").Children()
	assertLength(t, sel.Nodes, 5)
***REMOVED***

func TestChildrenRollback(t *testing.T) ***REMOVED***
	sel := Doc().Find(".pvk-content")
	sel2 := sel.Children().End()
	assertEqual(t, sel, sel2)
***REMOVED***

func TestContents(t *testing.T) ***REMOVED***
	sel := Doc().Find(".pvk-content").Contents()
	assertLength(t, sel.Nodes, 13)
***REMOVED***

func TestContentsRollback(t *testing.T) ***REMOVED***
	sel := Doc().Find(".pvk-content")
	sel2 := sel.Contents().End()
	assertEqual(t, sel, sel2)
***REMOVED***

func TestChildrenFiltered(t *testing.T) ***REMOVED***
	sel := Doc().Find(".pvk-content").ChildrenFiltered(".hero-unit")
	assertLength(t, sel.Nodes, 1)
***REMOVED***

func TestChildrenFilteredInvalid(t *testing.T) ***REMOVED***
	sel := Doc().Find(".pvk-content").ChildrenFiltered("")
	assertLength(t, sel.Nodes, 0)
***REMOVED***

func TestChildrenFilteredRollback(t *testing.T) ***REMOVED***
	sel := Doc().Find(".pvk-content")
	sel2 := sel.ChildrenFiltered(".hero-unit").End()
	assertEqual(t, sel, sel2)
***REMOVED***

func TestContentsFiltered(t *testing.T) ***REMOVED***
	sel := Doc().Find(".pvk-content").ContentsFiltered(".hero-unit")
	assertLength(t, sel.Nodes, 1)
***REMOVED***

func TestContentsFilteredInvalid(t *testing.T) ***REMOVED***
	sel := Doc().Find(".pvk-content").ContentsFiltered("~")
	assertLength(t, sel.Nodes, 0)
***REMOVED***

func TestContentsFilteredRollback(t *testing.T) ***REMOVED***
	sel := Doc().Find(".pvk-content")
	sel2 := sel.ContentsFiltered(".hero-unit").End()
	assertEqual(t, sel, sel2)
***REMOVED***

func TestChildrenFilteredNone(t *testing.T) ***REMOVED***
	sel := Doc().Find(".pvk-content").ChildrenFiltered("a.btn")
	assertLength(t, sel.Nodes, 0)
***REMOVED***

func TestParent(t *testing.T) ***REMOVED***
	sel := Doc().Find(".container-fluid").Parent()
	assertLength(t, sel.Nodes, 3)
***REMOVED***

func TestParentRollback(t *testing.T) ***REMOVED***
	sel := Doc().Find(".container-fluid")
	sel2 := sel.Parent().End()
	assertEqual(t, sel, sel2)
***REMOVED***

func TestParentBody(t *testing.T) ***REMOVED***
	sel := Doc().Find("body").Parent()
	assertLength(t, sel.Nodes, 1)
***REMOVED***

func TestParentFiltered(t *testing.T) ***REMOVED***
	sel := Doc().Find(".container-fluid").ParentFiltered(".hero-unit")
	assertLength(t, sel.Nodes, 1)
	assertClass(t, sel, "hero-unit")
***REMOVED***

func TestParentFilteredInvalid(t *testing.T) ***REMOVED***
	sel := Doc().Find(".container-fluid").ParentFiltered("")
	assertLength(t, sel.Nodes, 0)
***REMOVED***

func TestParentFilteredRollback(t *testing.T) ***REMOVED***
	sel := Doc().Find(".container-fluid")
	sel2 := sel.ParentFiltered(".hero-unit").End()
	assertEqual(t, sel, sel2)
***REMOVED***

func TestParents(t *testing.T) ***REMOVED***
	sel := Doc().Find(".container-fluid").Parents()
	assertLength(t, sel.Nodes, 8)
***REMOVED***

func TestParentsOrder(t *testing.T) ***REMOVED***
	sel := Doc().Find("#cf2").Parents()
	assertLength(t, sel.Nodes, 6)
	assertSelectionIs(t, sel, ".hero-unit", ".pvk-content", "div.row-fluid", "#cf1", "body", "html")
***REMOVED***

func TestParentsRollback(t *testing.T) ***REMOVED***
	sel := Doc().Find(".container-fluid")
	sel2 := sel.Parents().End()
	assertEqual(t, sel, sel2)
***REMOVED***

func TestParentsFiltered(t *testing.T) ***REMOVED***
	sel := Doc().Find(".container-fluid").ParentsFiltered("body")
	assertLength(t, sel.Nodes, 1)
***REMOVED***

func TestParentsFilteredInvalid(t *testing.T) ***REMOVED***
	sel := Doc().Find(".container-fluid").ParentsFiltered("")
	assertLength(t, sel.Nodes, 0)
***REMOVED***

func TestParentsFilteredRollback(t *testing.T) ***REMOVED***
	sel := Doc().Find(".container-fluid")
	sel2 := sel.ParentsFiltered("body").End()
	assertEqual(t, sel, sel2)
***REMOVED***

func TestParentsUntil(t *testing.T) ***REMOVED***
	sel := Doc().Find(".container-fluid").ParentsUntil("body")
	assertLength(t, sel.Nodes, 6)
***REMOVED***

func TestParentsUntilInvalid(t *testing.T) ***REMOVED***
	sel := Doc().Find(".container-fluid").ParentsUntil("")
	assertLength(t, sel.Nodes, 8)
***REMOVED***

func TestParentsUntilRollback(t *testing.T) ***REMOVED***
	sel := Doc().Find(".container-fluid")
	sel2 := sel.ParentsUntil("body").End()
	assertEqual(t, sel, sel2)
***REMOVED***

func TestParentsUntilSelection(t *testing.T) ***REMOVED***
	sel := Doc().Find(".container-fluid")
	sel2 := Doc().Find(".pvk-content")
	sel = sel.ParentsUntilSelection(sel2)
	assertLength(t, sel.Nodes, 3)
***REMOVED***

func TestParentsUntilSelectionRollback(t *testing.T) ***REMOVED***
	sel := Doc().Find(".container-fluid")
	sel2 := Doc().Find(".pvk-content")
	sel2 = sel.ParentsUntilSelection(sel2).End()
	assertEqual(t, sel, sel2)
***REMOVED***

func TestParentsUntilNodes(t *testing.T) ***REMOVED***
	sel := Doc().Find(".container-fluid")
	sel2 := Doc().Find(".pvk-content, .hero-unit")
	sel = sel.ParentsUntilNodes(sel2.Nodes...)
	assertLength(t, sel.Nodes, 2)
***REMOVED***

func TestParentsUntilNodesRollback(t *testing.T) ***REMOVED***
	sel := Doc().Find(".container-fluid")
	sel2 := Doc().Find(".pvk-content, .hero-unit")
	sel2 = sel.ParentsUntilNodes(sel2.Nodes...).End()
	assertEqual(t, sel, sel2)
***REMOVED***

func TestParentsFilteredUntil(t *testing.T) ***REMOVED***
	sel := Doc().Find(".container-fluid").ParentsFilteredUntil(".pvk-content", "body")
	assertLength(t, sel.Nodes, 2)
***REMOVED***

func TestParentsFilteredUntilInvalid(t *testing.T) ***REMOVED***
	sel := Doc().Find(".container-fluid").ParentsFilteredUntil("", "")
	assertLength(t, sel.Nodes, 0)
***REMOVED***

func TestParentsFilteredUntilRollback(t *testing.T) ***REMOVED***
	sel := Doc().Find(".container-fluid")
	sel2 := sel.ParentsFilteredUntil(".pvk-content", "body").End()
	assertEqual(t, sel, sel2)
***REMOVED***

func TestParentsFilteredUntilSelection(t *testing.T) ***REMOVED***
	sel := Doc().Find(".container-fluid")
	sel2 := Doc().Find(".row-fluid")
	sel = sel.ParentsFilteredUntilSelection("div", sel2)
	assertLength(t, sel.Nodes, 3)
***REMOVED***

func TestParentsFilteredUntilSelectionRollback(t *testing.T) ***REMOVED***
	sel := Doc().Find(".container-fluid")
	sel2 := Doc().Find(".row-fluid")
	sel2 = sel.ParentsFilteredUntilSelection("div", sel2).End()
	assertEqual(t, sel, sel2)
***REMOVED***

func TestParentsFilteredUntilNodes(t *testing.T) ***REMOVED***
	sel := Doc().Find(".container-fluid")
	sel2 := Doc().Find(".row-fluid")
	sel = sel.ParentsFilteredUntilNodes("body", sel2.Nodes...)
	assertLength(t, sel.Nodes, 1)
***REMOVED***

func TestParentsFilteredUntilNodesRollback(t *testing.T) ***REMOVED***
	sel := Doc().Find(".container-fluid")
	sel2 := Doc().Find(".row-fluid")
	sel2 = sel.ParentsFilteredUntilNodes("body", sel2.Nodes...).End()
	assertEqual(t, sel, sel2)
***REMOVED***

func TestSiblings(t *testing.T) ***REMOVED***
	sel := Doc().Find("h1").Siblings()
	assertLength(t, sel.Nodes, 1)
***REMOVED***

func TestSiblingsRollback(t *testing.T) ***REMOVED***
	sel := Doc().Find("h1")
	sel2 := sel.Siblings().End()
	assertEqual(t, sel, sel2)
***REMOVED***

func TestSiblings2(t *testing.T) ***REMOVED***
	sel := Doc().Find(".pvk-gutter").Siblings()
	assertLength(t, sel.Nodes, 9)
***REMOVED***

func TestSiblings3(t *testing.T) ***REMOVED***
	sel := Doc().Find("body>.container-fluid").Siblings()
	assertLength(t, sel.Nodes, 0)
***REMOVED***

func TestSiblingsFiltered(t *testing.T) ***REMOVED***
	sel := Doc().Find(".pvk-gutter").SiblingsFiltered(".pvk-content")
	assertLength(t, sel.Nodes, 3)
***REMOVED***

func TestSiblingsFilteredInvalid(t *testing.T) ***REMOVED***
	sel := Doc().Find(".pvk-gutter").SiblingsFiltered("")
	assertLength(t, sel.Nodes, 0)
***REMOVED***

func TestSiblingsFilteredRollback(t *testing.T) ***REMOVED***
	sel := Doc().Find(".pvk-gutter")
	sel2 := sel.SiblingsFiltered(".pvk-content").End()
	assertEqual(t, sel, sel2)
***REMOVED***

func TestNext(t *testing.T) ***REMOVED***
	sel := Doc().Find("h1").Next()
	assertLength(t, sel.Nodes, 1)
***REMOVED***

func TestNextRollback(t *testing.T) ***REMOVED***
	sel := Doc().Find("h1")
	sel2 := sel.Next().End()
	assertEqual(t, sel, sel2)
***REMOVED***

func TestNext2(t *testing.T) ***REMOVED***
	sel := Doc().Find(".close").Next()
	assertLength(t, sel.Nodes, 1)
***REMOVED***

func TestNextNone(t *testing.T) ***REMOVED***
	sel := Doc().Find("small").Next()
	assertLength(t, sel.Nodes, 0)
***REMOVED***

func TestNextFiltered(t *testing.T) ***REMOVED***
	sel := Doc().Find(".container-fluid").NextFiltered("div")
	assertLength(t, sel.Nodes, 2)
***REMOVED***

func TestNextFilteredInvalid(t *testing.T) ***REMOVED***
	sel := Doc().Find(".container-fluid").NextFiltered("")
	assertLength(t, sel.Nodes, 0)
***REMOVED***

func TestNextFilteredRollback(t *testing.T) ***REMOVED***
	sel := Doc().Find(".container-fluid")
	sel2 := sel.NextFiltered("div").End()
	assertEqual(t, sel, sel2)
***REMOVED***

func TestNextFiltered2(t *testing.T) ***REMOVED***
	sel := Doc().Find(".container-fluid").NextFiltered("[ng-view]")
	assertLength(t, sel.Nodes, 1)
***REMOVED***

func TestPrev(t *testing.T) ***REMOVED***
	sel := Doc().Find(".red").Prev()
	assertLength(t, sel.Nodes, 1)
	assertClass(t, sel, "green")
***REMOVED***

func TestPrevRollback(t *testing.T) ***REMOVED***
	sel := Doc().Find(".red")
	sel2 := sel.Prev().End()
	assertEqual(t, sel, sel2)
***REMOVED***

func TestPrev2(t *testing.T) ***REMOVED***
	sel := Doc().Find(".row-fluid").Prev()
	assertLength(t, sel.Nodes, 5)
***REMOVED***

func TestPrevNone(t *testing.T) ***REMOVED***
	sel := Doc().Find("h2").Prev()
	assertLength(t, sel.Nodes, 0)
***REMOVED***

func TestPrevFiltered(t *testing.T) ***REMOVED***
	sel := Doc().Find(".row-fluid").PrevFiltered(".row-fluid")
	assertLength(t, sel.Nodes, 5)
***REMOVED***

func TestPrevFilteredInvalid(t *testing.T) ***REMOVED***
	sel := Doc().Find(".row-fluid").PrevFiltered("")
	assertLength(t, sel.Nodes, 0)
***REMOVED***

func TestPrevFilteredRollback(t *testing.T) ***REMOVED***
	sel := Doc().Find(".row-fluid")
	sel2 := sel.PrevFiltered(".row-fluid").End()
	assertEqual(t, sel, sel2)
***REMOVED***

func TestNextAll(t *testing.T) ***REMOVED***
	sel := Doc().Find("#cf2 div:nth-child(1)").NextAll()
	assertLength(t, sel.Nodes, 3)
***REMOVED***

func TestNextAllRollback(t *testing.T) ***REMOVED***
	sel := Doc().Find("#cf2 div:nth-child(1)")
	sel2 := sel.NextAll().End()
	assertEqual(t, sel, sel2)
***REMOVED***

func TestNextAll2(t *testing.T) ***REMOVED***
	sel := Doc().Find("div[ng-cloak]").NextAll()
	assertLength(t, sel.Nodes, 1)
***REMOVED***

func TestNextAllNone(t *testing.T) ***REMOVED***
	sel := Doc().Find(".footer").NextAll()
	assertLength(t, sel.Nodes, 0)
***REMOVED***

func TestNextAllFiltered(t *testing.T) ***REMOVED***
	sel := Doc().Find("#cf2 .row-fluid").NextAllFiltered("[ng-cloak]")
	assertLength(t, sel.Nodes, 2)
***REMOVED***

func TestNextAllFilteredInvalid(t *testing.T) ***REMOVED***
	sel := Doc().Find("#cf2 .row-fluid").NextAllFiltered("")
	assertLength(t, sel.Nodes, 0)
***REMOVED***

func TestNextAllFilteredRollback(t *testing.T) ***REMOVED***
	sel := Doc().Find("#cf2 .row-fluid")
	sel2 := sel.NextAllFiltered("[ng-cloak]").End()
	assertEqual(t, sel, sel2)
***REMOVED***

func TestNextAllFiltered2(t *testing.T) ***REMOVED***
	sel := Doc().Find(".close").NextAllFiltered("h4")
	assertLength(t, sel.Nodes, 1)
***REMOVED***

func TestPrevAll(t *testing.T) ***REMOVED***
	sel := Doc().Find("[ng-view]").PrevAll()
	assertLength(t, sel.Nodes, 2)
***REMOVED***

func TestPrevAllOrder(t *testing.T) ***REMOVED***
	sel := Doc().Find("[ng-view]").PrevAll()
	assertLength(t, sel.Nodes, 2)
	assertSelectionIs(t, sel, "#cf4", "#cf3")
***REMOVED***

func TestPrevAllRollback(t *testing.T) ***REMOVED***
	sel := Doc().Find("[ng-view]")
	sel2 := sel.PrevAll().End()
	assertEqual(t, sel, sel2)
***REMOVED***

func TestPrevAll2(t *testing.T) ***REMOVED***
	sel := Doc().Find(".pvk-gutter").PrevAll()
	assertLength(t, sel.Nodes, 6)
***REMOVED***

func TestPrevAllFiltered(t *testing.T) ***REMOVED***
	sel := Doc().Find(".pvk-gutter").PrevAllFiltered(".pvk-content")
	assertLength(t, sel.Nodes, 3)
***REMOVED***

func TestPrevAllFilteredInvalid(t *testing.T) ***REMOVED***
	sel := Doc().Find(".pvk-gutter").PrevAllFiltered("")
	assertLength(t, sel.Nodes, 0)
***REMOVED***

func TestPrevAllFilteredRollback(t *testing.T) ***REMOVED***
	sel := Doc().Find(".pvk-gutter")
	sel2 := sel.PrevAllFiltered(".pvk-content").End()
	assertEqual(t, sel, sel2)
***REMOVED***

func TestNextUntil(t *testing.T) ***REMOVED***
	sel := Doc().Find(".alert a").NextUntil("p")
	assertLength(t, sel.Nodes, 1)
	assertSelectionIs(t, sel, "h4")
***REMOVED***

func TestNextUntilInvalid(t *testing.T) ***REMOVED***
	sel := Doc().Find(".alert a").NextUntil("")
	assertLength(t, sel.Nodes, 2)
***REMOVED***

func TestNextUntil2(t *testing.T) ***REMOVED***
	sel := Doc().Find("#cf2-1").NextUntil("[ng-cloak]")
	assertLength(t, sel.Nodes, 1)
	assertSelectionIs(t, sel, "#cf2-2")
***REMOVED***

func TestNextUntilOrder(t *testing.T) ***REMOVED***
	sel := Doc().Find("#cf2-1").NextUntil("#cf2-4")
	assertLength(t, sel.Nodes, 2)
	assertSelectionIs(t, sel, "#cf2-2", "#cf2-3")
***REMOVED***

func TestNextUntilRollback(t *testing.T) ***REMOVED***
	sel := Doc().Find("#cf2-1")
	sel2 := sel.PrevUntil("#cf2-4").End()
	assertEqual(t, sel, sel2)
***REMOVED***

func TestNextUntilSelection(t *testing.T) ***REMOVED***
	sel := Doc2().Find("#n2")
	sel2 := Doc2().Find("#n4")
	sel2 = sel.NextUntilSelection(sel2)
	assertLength(t, sel2.Nodes, 1)
	assertSelectionIs(t, sel2, "#n3")
***REMOVED***

func TestNextUntilSelectionRollback(t *testing.T) ***REMOVED***
	sel := Doc2().Find("#n2")
	sel2 := Doc2().Find("#n4")
	sel2 = sel.NextUntilSelection(sel2).End()
	assertEqual(t, sel, sel2)
***REMOVED***

func TestNextUntilNodes(t *testing.T) ***REMOVED***
	sel := Doc2().Find("#n2")
	sel2 := Doc2().Find("#n5")
	sel2 = sel.NextUntilNodes(sel2.Nodes...)
	assertLength(t, sel2.Nodes, 2)
	assertSelectionIs(t, sel2, "#n3", "#n4")
***REMOVED***

func TestNextUntilNodesRollback(t *testing.T) ***REMOVED***
	sel := Doc2().Find("#n2")
	sel2 := Doc2().Find("#n5")
	sel2 = sel.NextUntilNodes(sel2.Nodes...).End()
	assertEqual(t, sel, sel2)
***REMOVED***

func TestPrevUntil(t *testing.T) ***REMOVED***
	sel := Doc().Find(".alert p").PrevUntil("a")
	assertLength(t, sel.Nodes, 1)
	assertSelectionIs(t, sel, "h4")
***REMOVED***

func TestPrevUntilInvalid(t *testing.T) ***REMOVED***
	sel := Doc().Find(".alert p").PrevUntil("")
	assertLength(t, sel.Nodes, 2)
***REMOVED***

func TestPrevUntil2(t *testing.T) ***REMOVED***
	sel := Doc().Find("[ng-cloak]").PrevUntil(":not([ng-cloak])")
	assertLength(t, sel.Nodes, 1)
	assertSelectionIs(t, sel, "[ng-cloak]")
***REMOVED***

func TestPrevUntilOrder(t *testing.T) ***REMOVED***
	sel := Doc().Find("#cf2-4").PrevUntil("#cf2-1")
	assertLength(t, sel.Nodes, 2)
	assertSelectionIs(t, sel, "#cf2-3", "#cf2-2")
***REMOVED***

func TestPrevUntilRollback(t *testing.T) ***REMOVED***
	sel := Doc().Find("#cf2-4")
	sel2 := sel.PrevUntil("#cf2-1").End()
	assertEqual(t, sel, sel2)
***REMOVED***

func TestPrevUntilSelection(t *testing.T) ***REMOVED***
	sel := Doc2().Find("#n4")
	sel2 := Doc2().Find("#n2")
	sel2 = sel.PrevUntilSelection(sel2)
	assertLength(t, sel2.Nodes, 1)
	assertSelectionIs(t, sel2, "#n3")
***REMOVED***

func TestPrevUntilSelectionRollback(t *testing.T) ***REMOVED***
	sel := Doc2().Find("#n4")
	sel2 := Doc2().Find("#n2")
	sel2 = sel.PrevUntilSelection(sel2).End()
	assertEqual(t, sel, sel2)
***REMOVED***

func TestPrevUntilNodes(t *testing.T) ***REMOVED***
	sel := Doc2().Find("#n5")
	sel2 := Doc2().Find("#n2")
	sel2 = sel.PrevUntilNodes(sel2.Nodes...)
	assertLength(t, sel2.Nodes, 2)
	assertSelectionIs(t, sel2, "#n4", "#n3")
***REMOVED***

func TestPrevUntilNodesRollback(t *testing.T) ***REMOVED***
	sel := Doc2().Find("#n5")
	sel2 := Doc2().Find("#n2")
	sel2 = sel.PrevUntilNodes(sel2.Nodes...).End()
	assertEqual(t, sel, sel2)
***REMOVED***

func TestNextFilteredUntil(t *testing.T) ***REMOVED***
	sel := Doc2().Find(".two").NextFilteredUntil(".even", ".six")
	assertLength(t, sel.Nodes, 4)
	assertSelectionIs(t, sel, "#n3", "#n5", "#nf3", "#nf5")
***REMOVED***

func TestNextFilteredUntilInvalid(t *testing.T) ***REMOVED***
	sel := Doc2().Find(".two").NextFilteredUntil("", "")
	assertLength(t, sel.Nodes, 0)
***REMOVED***

func TestNextFilteredUntilRollback(t *testing.T) ***REMOVED***
	sel := Doc2().Find(".two")
	sel2 := sel.NextFilteredUntil(".even", ".six").End()
	assertEqual(t, sel, sel2)
***REMOVED***

func TestNextFilteredUntilSelection(t *testing.T) ***REMOVED***
	sel := Doc2().Find(".even")
	sel2 := Doc2().Find(".five")
	sel = sel.NextFilteredUntilSelection(".even", sel2)
	assertLength(t, sel.Nodes, 2)
	assertSelectionIs(t, sel, "#n3", "#nf3")
***REMOVED***

func TestNextFilteredUntilSelectionRollback(t *testing.T) ***REMOVED***
	sel := Doc2().Find(".even")
	sel2 := Doc2().Find(".five")
	sel3 := sel.NextFilteredUntilSelection(".even", sel2).End()
	assertEqual(t, sel, sel3)
***REMOVED***

func TestNextFilteredUntilNodes(t *testing.T) ***REMOVED***
	sel := Doc2().Find(".even")
	sel2 := Doc2().Find(".four")
	sel = sel.NextFilteredUntilNodes(".odd", sel2.Nodes...)
	assertLength(t, sel.Nodes, 4)
	assertSelectionIs(t, sel, "#n2", "#n6", "#nf2", "#nf6")
***REMOVED***

func TestNextFilteredUntilNodesRollback(t *testing.T) ***REMOVED***
	sel := Doc2().Find(".even")
	sel2 := Doc2().Find(".four")
	sel3 := sel.NextFilteredUntilNodes(".odd", sel2.Nodes...).End()
	assertEqual(t, sel, sel3)
***REMOVED***

func TestPrevFilteredUntil(t *testing.T) ***REMOVED***
	sel := Doc2().Find(".five").PrevFilteredUntil(".odd", ".one")
	assertLength(t, sel.Nodes, 4)
	assertSelectionIs(t, sel, "#n4", "#n2", "#nf4", "#nf2")
***REMOVED***

func TestPrevFilteredUntilInvalid(t *testing.T) ***REMOVED***
	sel := Doc2().Find(".five").PrevFilteredUntil("", "")
	assertLength(t, sel.Nodes, 0)
***REMOVED***

func TestPrevFilteredUntilRollback(t *testing.T) ***REMOVED***
	sel := Doc2().Find(".four")
	sel2 := sel.PrevFilteredUntil(".odd", ".one").End()
	assertEqual(t, sel, sel2)
***REMOVED***

func TestPrevFilteredUntilSelection(t *testing.T) ***REMOVED***
	sel := Doc2().Find(".odd")
	sel2 := Doc2().Find(".two")
	sel = sel.PrevFilteredUntilSelection(".odd", sel2)
	assertLength(t, sel.Nodes, 2)
	assertSelectionIs(t, sel, "#n4", "#nf4")
***REMOVED***

func TestPrevFilteredUntilSelectionRollback(t *testing.T) ***REMOVED***
	sel := Doc2().Find(".even")
	sel2 := Doc2().Find(".five")
	sel3 := sel.PrevFilteredUntilSelection(".even", sel2).End()
	assertEqual(t, sel, sel3)
***REMOVED***

func TestPrevFilteredUntilNodes(t *testing.T) ***REMOVED***
	sel := Doc2().Find(".even")
	sel2 := Doc2().Find(".four")
	sel = sel.PrevFilteredUntilNodes(".odd", sel2.Nodes...)
	assertLength(t, sel.Nodes, 2)
	assertSelectionIs(t, sel, "#n2", "#nf2")
***REMOVED***

func TestPrevFilteredUntilNodesRollback(t *testing.T) ***REMOVED***
	sel := Doc2().Find(".even")
	sel2 := Doc2().Find(".four")
	sel3 := sel.PrevFilteredUntilNodes(".odd", sel2.Nodes...).End()
	assertEqual(t, sel, sel3)
***REMOVED***

func TestClosestItself(t *testing.T) ***REMOVED***
	sel := Doc2().Find(".three")
	sel2 := sel.Closest(".row")
	assertLength(t, sel2.Nodes, sel.Length())
	assertSelectionIs(t, sel2, "#n3", "#nf3")
***REMOVED***

func TestClosestNoDupes(t *testing.T) ***REMOVED***
	sel := Doc().Find(".span12")
	sel2 := sel.Closest(".pvk-content")
	assertLength(t, sel2.Nodes, 1)
	assertClass(t, sel2, "pvk-content")
***REMOVED***

func TestClosestNone(t *testing.T) ***REMOVED***
	sel := Doc().Find("h4")
	sel2 := sel.Closest("a")
	assertLength(t, sel2.Nodes, 0)
***REMOVED***

func TestClosestInvalid(t *testing.T) ***REMOVED***
	sel := Doc().Find("h4")
	sel2 := sel.Closest("")
	assertLength(t, sel2.Nodes, 0)
***REMOVED***

func TestClosestMany(t *testing.T) ***REMOVED***
	sel := Doc().Find(".container-fluid")
	sel2 := sel.Closest(".pvk-content")
	assertLength(t, sel2.Nodes, 2)
	assertSelectionIs(t, sel2, "#pc1", "#pc2")
***REMOVED***

func TestClosestRollback(t *testing.T) ***REMOVED***
	sel := Doc().Find(".container-fluid")
	sel2 := sel.Closest(".pvk-content").End()
	assertEqual(t, sel, sel2)
***REMOVED***

func TestClosestSelectionItself(t *testing.T) ***REMOVED***
	sel := Doc2().Find(".three")
	sel2 := sel.ClosestSelection(Doc2().Find(".row"))
	assertLength(t, sel2.Nodes, sel.Length())
***REMOVED***

func TestClosestSelectionNoDupes(t *testing.T) ***REMOVED***
	sel := Doc().Find(".span12")
	sel2 := sel.ClosestSelection(Doc().Find(".pvk-content"))
	assertLength(t, sel2.Nodes, 1)
	assertClass(t, sel2, "pvk-content")
***REMOVED***

func TestClosestSelectionNone(t *testing.T) ***REMOVED***
	sel := Doc().Find("h4")
	sel2 := sel.ClosestSelection(Doc().Find("a"))
	assertLength(t, sel2.Nodes, 0)
***REMOVED***

func TestClosestSelectionMany(t *testing.T) ***REMOVED***
	sel := Doc().Find(".container-fluid")
	sel2 := sel.ClosestSelection(Doc().Find(".pvk-content"))
	assertLength(t, sel2.Nodes, 2)
	assertSelectionIs(t, sel2, "#pc1", "#pc2")
***REMOVED***

func TestClosestSelectionRollback(t *testing.T) ***REMOVED***
	sel := Doc().Find(".container-fluid")
	sel2 := sel.ClosestSelection(Doc().Find(".pvk-content")).End()
	assertEqual(t, sel, sel2)
***REMOVED***

func TestClosestNodesItself(t *testing.T) ***REMOVED***
	sel := Doc2().Find(".three")
	sel2 := sel.ClosestNodes(Doc2().Find(".row").Nodes...)
	assertLength(t, sel2.Nodes, sel.Length())
***REMOVED***

func TestClosestNodesNoDupes(t *testing.T) ***REMOVED***
	sel := Doc().Find(".span12")
	sel2 := sel.ClosestNodes(Doc().Find(".pvk-content").Nodes...)
	assertLength(t, sel2.Nodes, 1)
	assertClass(t, sel2, "pvk-content")
***REMOVED***

func TestClosestNodesNone(t *testing.T) ***REMOVED***
	sel := Doc().Find("h4")
	sel2 := sel.ClosestNodes(Doc().Find("a").Nodes...)
	assertLength(t, sel2.Nodes, 0)
***REMOVED***

func TestClosestNodesMany(t *testing.T) ***REMOVED***
	sel := Doc().Find(".container-fluid")
	sel2 := sel.ClosestNodes(Doc().Find(".pvk-content").Nodes...)
	assertLength(t, sel2.Nodes, 2)
	assertSelectionIs(t, sel2, "#pc1", "#pc2")
***REMOVED***

func TestClosestNodesRollback(t *testing.T) ***REMOVED***
	sel := Doc().Find(".container-fluid")
	sel2 := sel.ClosestNodes(Doc().Find(".pvk-content").Nodes...).End()
	assertEqual(t, sel, sel2)
***REMOVED***

func TestIssue26(t *testing.T) ***REMOVED***
	img1 := `<img src="assets/images/gallery/thumb-1.jpg" alt="150x150" />`
	img2 := `<img alt="150x150" src="assets/images/gallery/thumb-1.jpg" />`
	cases := []struct ***REMOVED***
		s string
		l int
	***REMOVED******REMOVED***
		***REMOVED***s: img1 + img2, l: 2***REMOVED***,
		***REMOVED***s: img1, l: 1***REMOVED***,
		***REMOVED***s: img2, l: 1***REMOVED***,
	***REMOVED***
	for _, c := range cases ***REMOVED***
		doc, err := NewDocumentFromReader(strings.NewReader(c.s))
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		sel := doc.Find("img[src]")
		assertLength(t, sel.Nodes, c.l)
	***REMOVED***
***REMOVED***
