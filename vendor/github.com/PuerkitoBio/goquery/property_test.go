package goquery

import (
	"regexp"
	"strings"
	"testing"
)

func TestAttrExists(t *testing.T) ***REMOVED***
	if val, ok := Doc().Find("a").Attr("href"); !ok ***REMOVED***
		t.Error("Expected a value for the href attribute.")
	***REMOVED*** else ***REMOVED***
		t.Logf("Href of first anchor: %v.", val)
	***REMOVED***
***REMOVED***

func TestAttrOr(t *testing.T) ***REMOVED***
	if val := Doc().Find("a").AttrOr("fake-attribute", "alternative"); val != "alternative" ***REMOVED***
		t.Error("Expected an alternative value for 'fake-attribute' attribute.")
	***REMOVED*** else ***REMOVED***
		t.Logf("Value returned for not existing attribute: %v.", val)
	***REMOVED***
	if val := Doc().Find("zz").AttrOr("fake-attribute", "alternative"); val != "alternative" ***REMOVED***
		t.Error("Expected an alternative value for 'fake-attribute' on an empty selection.")
	***REMOVED*** else ***REMOVED***
		t.Logf("Value returned for empty selection: %v.", val)
	***REMOVED***
***REMOVED***

func TestAttrNotExist(t *testing.T) ***REMOVED***
	if val, ok := Doc().Find("div.row-fluid").Attr("href"); ok ***REMOVED***
		t.Errorf("Expected no value for the href attribute, got %v.", val)
	***REMOVED***
***REMOVED***

func TestRemoveAttr(t *testing.T) ***REMOVED***
	sel := Doc2Clone().Find("div")

	sel.RemoveAttr("id")

	_, ok := sel.Attr("id")
	if ok ***REMOVED***
		t.Error("Expected there to be no id attributes set")
	***REMOVED***
***REMOVED***

func TestSetAttr(t *testing.T) ***REMOVED***
	sel := Doc2Clone().Find("#main")

	sel.SetAttr("id", "not-main")

	val, ok := sel.Attr("id")
	if !ok ***REMOVED***
		t.Error("Expected an id attribute on main")
	***REMOVED***

	if val != "not-main" ***REMOVED***
		t.Errorf("Expected an attribute id to be not-main, got %s", val)
	***REMOVED***
***REMOVED***

func TestSetAttr2(t *testing.T) ***REMOVED***
	sel := Doc2Clone().Find("#main")

	sel.SetAttr("foo", "bar")

	val, ok := sel.Attr("foo")
	if !ok ***REMOVED***
		t.Error("Expected an 'foo' attribute on main")
	***REMOVED***

	if val != "bar" ***REMOVED***
		t.Errorf("Expected an attribute 'foo' to be 'bar', got '%s'", val)
	***REMOVED***
***REMOVED***

func TestText(t *testing.T) ***REMOVED***
	txt := Doc().Find("h1").Text()
	if strings.Trim(txt, " \n\r\t") != "Provok.in" ***REMOVED***
		t.Errorf("Expected text to be Provok.in, found %s.", txt)
	***REMOVED***
***REMOVED***

func TestText2(t *testing.T) ***REMOVED***
	txt := Doc().Find(".hero-unit .container-fluid .row-fluid:nth-child(1)").Text()
	if ok, e := regexp.MatchString(`^\s+Provok\.in\s+Prove your point.\s+$`, txt); !ok || e != nil ***REMOVED***
		t.Errorf("Expected text to be Provok.in Prove your point., found %s.", txt)
		if e != nil ***REMOVED***
			t.Logf("Error: %s.", e.Error())
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestText3(t *testing.T) ***REMOVED***
	txt := Doc().Find(".pvk-gutter").First().Text()
	// There's an &nbsp; character in there...
	if ok, e := regexp.MatchString(`^[\s\x***REMOVED***00A0***REMOVED***]+$`, txt); !ok || e != nil ***REMOVED***
		t.Errorf("Expected spaces, found <%v>.", txt)
		if e != nil ***REMOVED***
			t.Logf("Error: %s.", e.Error())
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestHtml(t *testing.T) ***REMOVED***
	txt, e := Doc().Find("h1").Html()
	if e != nil ***REMOVED***
		t.Errorf("Error: %s.", e)
	***REMOVED***

	if ok, e := regexp.MatchString(`^\s*<a href="/">Provok<span class="green">\.</span><span class="red">i</span>n</a>\s*$`, txt); !ok || e != nil ***REMOVED***
		t.Errorf("Unexpected HTML content, found %s.", txt)
		if e != nil ***REMOVED***
			t.Logf("Error: %s.", e.Error())
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestNbsp(t *testing.T) ***REMOVED***
	src := `<p>Some&nbsp;text</p>`
	d, err := NewDocumentFromReader(strings.NewReader(src))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	txt := d.Find("p").Text()
	ix := strings.Index(txt, "\u00a0")
	if ix != 4 ***REMOVED***
		t.Errorf("Text: expected a non-breaking space at index 4, got %d", ix)
	***REMOVED***

	h, err := d.Find("p").Html()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	ix = strings.Index(h, "\u00a0")
	if ix != 4 ***REMOVED***
		t.Errorf("Html: expected a non-breaking space at index 4, got %d", ix)
	***REMOVED***
***REMOVED***

func TestAddClass(t *testing.T) ***REMOVED***
	sel := Doc2Clone().Find("#main")
	sel.AddClass("main main main")

	// Make sure that class was only added once
	if a, ok := sel.Attr("class"); !ok || a != "main" ***REMOVED***
		t.Error("Expected #main to have class main")
	***REMOVED***
***REMOVED***

func TestAddClassSimilar(t *testing.T) ***REMOVED***
	sel := Doc2Clone().Find("#nf5")
	sel.AddClass("odd")

	assertClass(t, sel, "odd")
	assertClass(t, sel, "odder")
	printSel(t, sel.Parent())
***REMOVED***

func TestAddEmptyClass(t *testing.T) ***REMOVED***
	sel := Doc2Clone().Find("#main")
	sel.AddClass("")

	// Make sure that class was only added once
	if a, ok := sel.Attr("class"); ok ***REMOVED***
		t.Errorf("Expected #main to not to have a class, have: %s", a)
	***REMOVED***
***REMOVED***

func TestAddClasses(t *testing.T) ***REMOVED***
	sel := Doc2Clone().Find("#main")
	sel.AddClass("a b")

	// Make sure that class was only added once
	if !sel.HasClass("a") || !sel.HasClass("b") ***REMOVED***
		t.Errorf("#main does not have classes")
	***REMOVED***
***REMOVED***

func TestHasClass(t *testing.T) ***REMOVED***
	sel := Doc().Find("div")
	if !sel.HasClass("span12") ***REMOVED***
		t.Error("Expected at least one div to have class span12.")
	***REMOVED***
***REMOVED***

func TestHasClassNone(t *testing.T) ***REMOVED***
	sel := Doc().Find("h2")
	if sel.HasClass("toto") ***REMOVED***
		t.Error("Expected h1 to have no class.")
	***REMOVED***
***REMOVED***

func TestHasClassNotFirst(t *testing.T) ***REMOVED***
	sel := Doc().Find(".alert")
	if !sel.HasClass("alert-error") ***REMOVED***
		t.Error("Expected .alert to also have class .alert-error.")
	***REMOVED***
***REMOVED***

func TestRemoveClass(t *testing.T) ***REMOVED***
	sel := Doc2Clone().Find("#nf1")
	sel.RemoveClass("one row")

	if !sel.HasClass("even") || sel.HasClass("one") || sel.HasClass("row") ***REMOVED***
		classes, _ := sel.Attr("class")
		t.Error("Expected #nf1 to have class even, has ", classes)
	***REMOVED***
***REMOVED***

func TestRemoveClassSimilar(t *testing.T) ***REMOVED***
	sel := Doc2Clone().Find("#nf5, #nf6")
	assertLength(t, sel.Nodes, 2)

	sel.RemoveClass("odd")
	assertClass(t, sel.Eq(0), "odder")
	printSel(t, sel)
***REMOVED***

func TestRemoveAllClasses(t *testing.T) ***REMOVED***
	sel := Doc2Clone().Find("#nf1")
	sel.RemoveClass()

	if a, ok := sel.Attr("class"); ok ***REMOVED***
		t.Error("All classes were not removed, has ", a)
	***REMOVED***

	sel = Doc2Clone().Find("#main")
	sel.RemoveClass()
	if a, ok := sel.Attr("class"); ok ***REMOVED***
		t.Error("All classes were not removed, has ", a)
	***REMOVED***
***REMOVED***

func TestToggleClass(t *testing.T) ***REMOVED***
	sel := Doc2Clone().Find("#nf1")

	sel.ToggleClass("one")
	if sel.HasClass("one") ***REMOVED***
		t.Error("Expected #nf1 to not have class one")
	***REMOVED***

	sel.ToggleClass("one")
	if !sel.HasClass("one") ***REMOVED***
		t.Error("Expected #nf1 to have class one")
	***REMOVED***

	sel.ToggleClass("one even row")
	if a, ok := sel.Attr("class"); ok ***REMOVED***
		t.Errorf("Expected #nf1 to have no classes, have %q", a)
	***REMOVED***
***REMOVED***
