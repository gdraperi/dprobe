package goquery

import (
	"testing"
)

const (
	wrapHtml = "<div id=\"ins\">test string<div><p><em><b></b></em></p></div></div>"
)

func TestAfter(t *testing.T) ***REMOVED***
	doc := Doc2Clone()
	doc.Find("#main").After("#nf6")

	assertLength(t, doc.Find("#main #nf6").Nodes, 0)
	assertLength(t, doc.Find("#foot #nf6").Nodes, 0)
	assertLength(t, doc.Find("#main + #nf6").Nodes, 1)
	printSel(t, doc.Selection)
***REMOVED***

func TestAfterMany(t *testing.T) ***REMOVED***
	doc := Doc2Clone()
	doc.Find(".one").After("#nf6")

	assertLength(t, doc.Find("#foot #nf6").Nodes, 1)
	assertLength(t, doc.Find("#main #nf6").Nodes, 1)
	assertLength(t, doc.Find(".one + #nf6").Nodes, 2)
	printSel(t, doc.Selection)
***REMOVED***

func TestAfterWithRemoved(t *testing.T) ***REMOVED***
	doc := Doc2Clone()
	s := doc.Find("#main").Remove()
	s.After("#nf6")

	assertLength(t, s.Find("#nf6").Nodes, 0)
	assertLength(t, doc.Find("#nf6").Nodes, 0)
	printSel(t, doc.Selection)
***REMOVED***

func TestAfterSelection(t *testing.T) ***REMOVED***
	doc := Doc2Clone()
	doc.Find("#main").AfterSelection(doc.Find("#nf1, #nf2"))

	assertLength(t, doc.Find("#main #nf1, #main #nf2").Nodes, 0)
	assertLength(t, doc.Find("#foot #nf1, #foot #nf2").Nodes, 0)
	assertLength(t, doc.Find("#main + #nf1, #nf1 + #nf2").Nodes, 2)
	printSel(t, doc.Selection)
***REMOVED***

func TestAfterHtml(t *testing.T) ***REMOVED***
	doc := Doc2Clone()
	doc.Find("#main").AfterHtml("<strong>new node</strong>")

	assertLength(t, doc.Find("#main + strong").Nodes, 1)
	printSel(t, doc.Selection)
***REMOVED***

func TestAppend(t *testing.T) ***REMOVED***
	doc := Doc2Clone()
	doc.Find("#main").Append("#nf6")

	assertLength(t, doc.Find("#foot #nf6").Nodes, 0)
	assertLength(t, doc.Find("#main #nf6").Nodes, 1)
	printSel(t, doc.Selection)
***REMOVED***

func TestAppendBody(t *testing.T) ***REMOVED***
	doc := Doc2Clone()
	doc.Find("body").Append("#nf6")

	assertLength(t, doc.Find("#foot #nf6").Nodes, 0)
	assertLength(t, doc.Find("#main #nf6").Nodes, 0)
	assertLength(t, doc.Find("body > #nf6").Nodes, 1)
	printSel(t, doc.Selection)
***REMOVED***

func TestAppendSelection(t *testing.T) ***REMOVED***
	doc := Doc2Clone()
	doc.Find("#main").AppendSelection(doc.Find("#nf1, #nf2"))

	assertLength(t, doc.Find("#foot #nf1").Nodes, 0)
	assertLength(t, doc.Find("#foot #nf2").Nodes, 0)
	assertLength(t, doc.Find("#main #nf1").Nodes, 1)
	assertLength(t, doc.Find("#main #nf2").Nodes, 1)
	printSel(t, doc.Selection)
***REMOVED***

func TestAppendSelectionExisting(t *testing.T) ***REMOVED***
	doc := Doc2Clone()
	doc.Find("#main").AppendSelection(doc.Find("#n1, #n2"))

	assertClass(t, doc.Find("#main :nth-child(1)"), "three")
	assertClass(t, doc.Find("#main :nth-child(5)"), "one")
	assertClass(t, doc.Find("#main :nth-child(6)"), "two")
	printSel(t, doc.Selection)
***REMOVED***

func TestAppendClone(t *testing.T) ***REMOVED***
	doc := Doc2Clone()
	doc.Find("#n1").AppendSelection(doc.Find("#nf1").Clone())

	assertLength(t, doc.Find("#foot #nf1").Nodes, 1)
	assertLength(t, doc.Find("#main #nf1").Nodes, 1)
	printSel(t, doc.Selection)
***REMOVED***

func TestAppendHtml(t *testing.T) ***REMOVED***
	doc := Doc2Clone()
	doc.Find("div").AppendHtml("<strong>new node</strong>")

	assertLength(t, doc.Find("strong").Nodes, 14)
	printSel(t, doc.Selection)
***REMOVED***

func TestBefore(t *testing.T) ***REMOVED***
	doc := Doc2Clone()
	doc.Find("#main").Before("#nf6")

	assertLength(t, doc.Find("#main #nf6").Nodes, 0)
	assertLength(t, doc.Find("#foot #nf6").Nodes, 0)
	assertLength(t, doc.Find("body > #nf6:first-child").Nodes, 1)
	printSel(t, doc.Selection)
***REMOVED***

func TestBeforeWithRemoved(t *testing.T) ***REMOVED***
	doc := Doc2Clone()
	s := doc.Find("#main").Remove()
	s.Before("#nf6")

	assertLength(t, s.Find("#nf6").Nodes, 0)
	assertLength(t, doc.Find("#nf6").Nodes, 0)
	printSel(t, doc.Selection)
***REMOVED***

func TestBeforeSelection(t *testing.T) ***REMOVED***
	doc := Doc2Clone()
	doc.Find("#main").BeforeSelection(doc.Find("#nf1, #nf2"))

	assertLength(t, doc.Find("#main #nf1, #main #nf2").Nodes, 0)
	assertLength(t, doc.Find("#foot #nf1, #foot #nf2").Nodes, 0)
	assertLength(t, doc.Find("body > #nf1:first-child, #nf1 + #nf2").Nodes, 2)
	printSel(t, doc.Selection)
***REMOVED***

func TestBeforeHtml(t *testing.T) ***REMOVED***
	doc := Doc2Clone()
	doc.Find("#main").BeforeHtml("<strong>new node</strong>")

	assertLength(t, doc.Find("body > strong:first-child").Nodes, 1)
	printSel(t, doc.Selection)
***REMOVED***

func TestEmpty(t *testing.T) ***REMOVED***
	doc := Doc2Clone()
	s := doc.Find("#main").Empty()

	assertLength(t, doc.Find("#main").Children().Nodes, 0)
	assertLength(t, s.Filter("div").Nodes, 6)
	printSel(t, doc.Selection)
***REMOVED***

func TestPrepend(t *testing.T) ***REMOVED***
	doc := Doc2Clone()
	doc.Find("#main").Prepend("#nf6")

	assertLength(t, doc.Find("#foot #nf6").Nodes, 0)
	assertLength(t, doc.Find("#main #nf6:first-child").Nodes, 1)
	printSel(t, doc.Selection)
***REMOVED***

func TestPrependBody(t *testing.T) ***REMOVED***
	doc := Doc2Clone()
	doc.Find("body").Prepend("#nf6")

	assertLength(t, doc.Find("#foot #nf6").Nodes, 0)
	assertLength(t, doc.Find("#main #nf6").Nodes, 0)
	assertLength(t, doc.Find("body > #nf6:first-child").Nodes, 1)
	printSel(t, doc.Selection)
***REMOVED***

func TestPrependSelection(t *testing.T) ***REMOVED***
	doc := Doc2Clone()
	doc.Find("#main").PrependSelection(doc.Find("#nf1, #nf2"))

	assertLength(t, doc.Find("#foot #nf1").Nodes, 0)
	assertLength(t, doc.Find("#foot #nf2").Nodes, 0)
	assertLength(t, doc.Find("#main #nf1:first-child").Nodes, 1)
	assertLength(t, doc.Find("#main #nf2:nth-child(2)").Nodes, 1)
	printSel(t, doc.Selection)
***REMOVED***

func TestPrependSelectionExisting(t *testing.T) ***REMOVED***
	doc := Doc2Clone()
	doc.Find("#main").PrependSelection(doc.Find("#n5, #n6"))

	assertClass(t, doc.Find("#main :nth-child(1)"), "five")
	assertClass(t, doc.Find("#main :nth-child(2)"), "six")
	assertClass(t, doc.Find("#main :nth-child(5)"), "three")
	assertClass(t, doc.Find("#main :nth-child(6)"), "four")
	printSel(t, doc.Selection)
***REMOVED***

func TestPrependClone(t *testing.T) ***REMOVED***
	doc := Doc2Clone()
	doc.Find("#n1").PrependSelection(doc.Find("#nf1").Clone())

	assertLength(t, doc.Find("#foot #nf1:first-child").Nodes, 1)
	assertLength(t, doc.Find("#main #nf1:first-child").Nodes, 1)
	printSel(t, doc.Selection)
***REMOVED***

func TestPrependHtml(t *testing.T) ***REMOVED***
	doc := Doc2Clone()
	doc.Find("div").PrependHtml("<strong>new node</strong>")

	assertLength(t, doc.Find("strong:first-child").Nodes, 14)
	printSel(t, doc.Selection)
***REMOVED***

func TestRemove(t *testing.T) ***REMOVED***
	doc := Doc2Clone()
	doc.Find("#nf1").Remove()

	assertLength(t, doc.Find("#foot #nf1").Nodes, 0)
	printSel(t, doc.Selection)
***REMOVED***

func TestRemoveAll(t *testing.T) ***REMOVED***
	doc := Doc2Clone()
	doc.Find("*").Remove()

	assertLength(t, doc.Find("*").Nodes, 0)
	printSel(t, doc.Selection)
***REMOVED***

func TestRemoveRoot(t *testing.T) ***REMOVED***
	doc := Doc2Clone()
	doc.Find("html").Remove()

	assertLength(t, doc.Find("html").Nodes, 0)
	printSel(t, doc.Selection)
***REMOVED***

func TestRemoveFiltered(t *testing.T) ***REMOVED***
	doc := Doc2Clone()
	nf6 := doc.Find("#nf6")
	s := doc.Find("div").RemoveFiltered("#nf6")

	assertLength(t, doc.Find("#nf6").Nodes, 0)
	assertLength(t, s.Nodes, 1)
	if nf6.Nodes[0] != s.Nodes[0] ***REMOVED***
		t.Error("Removed node does not match original")
	***REMOVED***
	printSel(t, doc.Selection)
***REMOVED***

func TestReplaceWith(t *testing.T) ***REMOVED***
	doc := Doc2Clone()

	doc.Find("#nf6").ReplaceWith("#main")
	assertLength(t, doc.Find("#foot #main:last-child").Nodes, 1)
	printSel(t, doc.Selection)

	doc.Find("#foot").ReplaceWith("#main")
	assertLength(t, doc.Find("#foot").Nodes, 0)
	assertLength(t, doc.Find("#main").Nodes, 1)

	printSel(t, doc.Selection)
***REMOVED***

func TestReplaceWithHtml(t *testing.T) ***REMOVED***
	doc := Doc2Clone()
	doc.Find("#main, #foot").ReplaceWithHtml("<div id=\"replace\"></div>")

	assertLength(t, doc.Find("#replace").Nodes, 2)

	printSel(t, doc.Selection)
***REMOVED***

func TestSetHtml(t *testing.T) ***REMOVED***
	doc := Doc2Clone()
	q := doc.Find("#main, #foot")
	q.SetHtml(`<div id="replace">test</div>`)

	assertLength(t, doc.Find("#replace").Nodes, 2)
	assertLength(t, doc.Find("#main, #foot").Nodes, 2)

	if q.Text() != "testtest" ***REMOVED***
		t.Errorf("Expected text to be %v, found %v", "testtest", q.Text())
	***REMOVED***

	printSel(t, doc.Selection)
***REMOVED***

func TestSetHtmlNoMatch(t *testing.T) ***REMOVED***
	doc := Doc2Clone()
	q := doc.Find("#notthere")
	q.SetHtml(`<div id="replace">test</div>`)

	assertLength(t, doc.Find("#replace").Nodes, 0)

	printSel(t, doc.Selection)
***REMOVED***

func TestSetHtmlEmpty(t *testing.T) ***REMOVED***
	doc := Doc2Clone()
	q := doc.Find("#main")
	q.SetHtml(``)

	assertLength(t, doc.Find("#main").Nodes, 1)
	assertLength(t, doc.Find("#main").Children().Nodes, 0)
	printSel(t, doc.Selection)
***REMOVED***

func TestSetText(t *testing.T) ***REMOVED***
	doc := Doc2Clone()
	q := doc.Find("#main, #foot")
	repl := "<div id=\"replace\">test</div>"
	q.SetText(repl)

	assertLength(t, doc.Find("#replace").Nodes, 0)
	assertLength(t, doc.Find("#main, #foot").Nodes, 2)

	if q.Text() != (repl + repl) ***REMOVED***
		t.Errorf("Expected text to be %v, found %v", (repl + repl), q.Text())
	***REMOVED***

	h, err := q.Html()
	if err != nil ***REMOVED***
		t.Errorf("Error: %v", err)
	***REMOVED***
	esc := "&lt;div id=&#34;replace&#34;&gt;test&lt;/div&gt;"
	if h != esc ***REMOVED***
		t.Errorf("Expected html to be %v, found %v", esc, h)
	***REMOVED***

	printSel(t, doc.Selection)
***REMOVED***

func TestReplaceWithSelection(t *testing.T) ***REMOVED***
	doc := Doc2Clone()
	sel := doc.Find("#nf6").ReplaceWithSelection(doc.Find("#nf5"))

	assertSelectionIs(t, sel, "#nf6")
	assertLength(t, doc.Find("#nf6").Nodes, 0)
	assertLength(t, doc.Find("#nf5").Nodes, 1)

	printSel(t, doc.Selection)
***REMOVED***

func TestUnwrap(t *testing.T) ***REMOVED***
	doc := Doc2Clone()

	doc.Find("#nf5").Unwrap()
	assertLength(t, doc.Find("#foot").Nodes, 0)
	assertLength(t, doc.Find("body > #nf1").Nodes, 1)
	assertLength(t, doc.Find("body > #nf5").Nodes, 1)

	printSel(t, doc.Selection)

	doc = Doc2Clone()

	doc.Find("#nf5, #n1").Unwrap()
	assertLength(t, doc.Find("#foot").Nodes, 0)
	assertLength(t, doc.Find("#main").Nodes, 0)
	assertLength(t, doc.Find("body > #n1").Nodes, 1)
	assertLength(t, doc.Find("body > #nf5").Nodes, 1)

	printSel(t, doc.Selection)
***REMOVED***

func TestUnwrapBody(t *testing.T) ***REMOVED***
	doc := Doc2Clone()

	doc.Find("#main").Unwrap()
	assertLength(t, doc.Find("body").Nodes, 1)
	assertLength(t, doc.Find("body > #main").Nodes, 1)

	printSel(t, doc.Selection)
***REMOVED***

func TestUnwrapHead(t *testing.T) ***REMOVED***
	doc := Doc2Clone()

	doc.Find("title").Unwrap()
	assertLength(t, doc.Find("head").Nodes, 0)
	assertLength(t, doc.Find("head > title").Nodes, 0)
	assertLength(t, doc.Find("title").Nodes, 1)

	printSel(t, doc.Selection)
***REMOVED***

func TestUnwrapHtml(t *testing.T) ***REMOVED***
	doc := Doc2Clone()

	doc.Find("head").Unwrap()
	assertLength(t, doc.Find("html").Nodes, 0)
	assertLength(t, doc.Find("html head").Nodes, 0)
	assertLength(t, doc.Find("head").Nodes, 1)

	printSel(t, doc.Selection)
***REMOVED***

func TestWrap(t *testing.T) ***REMOVED***
	doc := Doc2Clone()
	doc.Find("#nf1").Wrap("#nf2")
	nf1 := doc.Find("#foot #nf2 #nf1")
	assertLength(t, nf1.Nodes, 1)

	nf2 := doc.Find("#nf2")
	assertLength(t, nf2.Nodes, 2)

	printSel(t, doc.Selection)
***REMOVED***

func TestWrapEmpty(t *testing.T) ***REMOVED***
	doc := Doc2Clone()
	doc.Find("#nf1").Wrap("#doesnt-exist")

	origHtml, _ := Doc2().Html()
	newHtml, _ := doc.Html()

	if origHtml != newHtml ***REMOVED***
		t.Error("Expected the two documents to be identical.")
	***REMOVED***

	printSel(t, doc.Selection)
***REMOVED***

func TestWrapHtml(t *testing.T) ***REMOVED***
	doc := Doc2Clone()
	doc.Find(".odd").WrapHtml(wrapHtml)
	nf2 := doc.Find("#ins #nf2")
	assertLength(t, nf2.Nodes, 1)
	printSel(t, doc.Selection)
***REMOVED***

func TestWrapSelection(t *testing.T) ***REMOVED***
	doc := Doc2Clone()
	doc.Find("#nf1").WrapSelection(doc.Find("#nf2"))
	nf1 := doc.Find("#foot #nf2 #nf1")
	assertLength(t, nf1.Nodes, 1)

	nf2 := doc.Find("#nf2")
	assertLength(t, nf2.Nodes, 2)

	printSel(t, doc.Selection)
***REMOVED***

func TestWrapAll(t *testing.T) ***REMOVED***
	doc := Doc2Clone()
	doc.Find(".odd").WrapAll("#nf1")
	nf1 := doc.Find("#main #nf1")
	assertLength(t, nf1.Nodes, 1)

	sel := nf1.Find("#n2 ~ #n4 ~ #n6 ~ #nf2 ~ #nf4 ~ #nf6")
	assertLength(t, sel.Nodes, 1)

	printSel(t, doc.Selection)
***REMOVED***

func TestWrapAllHtml(t *testing.T) ***REMOVED***
	doc := Doc2Clone()
	doc.Find(".odd").WrapAllHtml(wrapHtml)
	nf1 := doc.Find("#main div#ins div p em b #n2 ~ #n4 ~ #n6 ~ #nf2 ~ #nf4 ~ #nf6")
	assertLength(t, nf1.Nodes, 1)
	printSel(t, doc.Selection)
***REMOVED***

func TestWrapInnerNoContent(t *testing.T) ***REMOVED***
	doc := Doc2Clone()
	doc.Find(".one").WrapInner(".two")

	twos := doc.Find(".two")
	assertLength(t, twos.Nodes, 4)
	assertLength(t, doc.Find(".one .two").Nodes, 2)

	printSel(t, doc.Selection)
***REMOVED***

func TestWrapInnerWithContent(t *testing.T) ***REMOVED***
	doc := Doc3Clone()
	doc.Find(".one").WrapInner(".two")

	twos := doc.Find(".two")
	assertLength(t, twos.Nodes, 4)
	assertLength(t, doc.Find(".one .two").Nodes, 2)

	printSel(t, doc.Selection)
***REMOVED***

func TestWrapInnerNoWrapper(t *testing.T) ***REMOVED***
	doc := Doc2Clone()
	doc.Find(".one").WrapInner(".not-exist")

	twos := doc.Find(".two")
	assertLength(t, twos.Nodes, 2)
	assertLength(t, doc.Find(".one").Nodes, 2)
	assertLength(t, doc.Find(".one .two").Nodes, 0)

	printSel(t, doc.Selection)
***REMOVED***

func TestWrapInnerHtml(t *testing.T) ***REMOVED***
	doc := Doc2Clone()
	doc.Find("#foot").WrapInnerHtml(wrapHtml)

	foot := doc.Find("#foot div#ins div p em b #nf1 ~ #nf2 ~ #nf3")
	assertLength(t, foot.Nodes, 1)

	printSel(t, doc.Selection)
***REMOVED***
