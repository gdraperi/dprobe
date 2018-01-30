package goquery

import (
	"testing"
)

func TestIs(t *testing.T) ***REMOVED***
	sel := Doc().Find(".footer p:nth-child(1)")
	if !sel.Is("p") ***REMOVED***
		t.Error("Expected .footer p:nth-child(1) to be p.")
	***REMOVED***
***REMOVED***

func TestIsInvalid(t *testing.T) ***REMOVED***
	sel := Doc().Find(".footer p:nth-child(1)")
	if sel.Is("") ***REMOVED***
		t.Error("Is should not succeed with invalid selector string")
	***REMOVED***
***REMOVED***

func TestIsPositional(t *testing.T) ***REMOVED***
	sel := Doc().Find(".footer p:nth-child(2)")
	if !sel.Is("p:nth-child(2)") ***REMOVED***
		t.Error("Expected .footer p:nth-child(2) to be p:nth-child(2).")
	***REMOVED***
***REMOVED***

func TestIsPositionalNot(t *testing.T) ***REMOVED***
	sel := Doc().Find(".footer p:nth-child(1)")
	if sel.Is("p:nth-child(2)") ***REMOVED***
		t.Error("Expected .footer p:nth-child(1) NOT to be p:nth-child(2).")
	***REMOVED***
***REMOVED***

func TestIsFunction(t *testing.T) ***REMOVED***
	ok := Doc().Find("div").IsFunction(func(i int, s *Selection) bool ***REMOVED***
		return s.HasClass("container-fluid")
	***REMOVED***)

	if !ok ***REMOVED***
		t.Error("Expected some div to have a container-fluid class.")
	***REMOVED***
***REMOVED***

func TestIsFunctionRollback(t *testing.T) ***REMOVED***
	ok := Doc().Find("div").IsFunction(func(i int, s *Selection) bool ***REMOVED***
		return s.HasClass("container-fluid")
	***REMOVED***)

	if !ok ***REMOVED***
		t.Error("Expected some div to have a container-fluid class.")
	***REMOVED***
***REMOVED***

func TestIsSelection(t *testing.T) ***REMOVED***
	sel := Doc().Find("div")
	sel2 := Doc().Find(".pvk-gutter")

	if !sel.IsSelection(sel2) ***REMOVED***
		t.Error("Expected some div to have a pvk-gutter class.")
	***REMOVED***
***REMOVED***

func TestIsSelectionNot(t *testing.T) ***REMOVED***
	sel := Doc().Find("div")
	sel2 := Doc().Find("a")

	if sel.IsSelection(sel2) ***REMOVED***
		t.Error("Expected some div NOT to be an anchor.")
	***REMOVED***
***REMOVED***

func TestIsNodes(t *testing.T) ***REMOVED***
	sel := Doc().Find("div")
	sel2 := Doc().Find(".footer")

	if !sel.IsNodes(sel2.Nodes[0]) ***REMOVED***
		t.Error("Expected some div to have a footer class.")
	***REMOVED***
***REMOVED***

func TestDocContains(t *testing.T) ***REMOVED***
	sel := Doc().Find("h1")
	if !Doc().Contains(sel.Nodes[0]) ***REMOVED***
		t.Error("Expected document to contain H1 tag.")
	***REMOVED***
***REMOVED***

func TestSelContains(t *testing.T) ***REMOVED***
	sel := Doc().Find(".row-fluid")
	sel2 := Doc().Find("a[ng-click]")
	if !sel.Contains(sel2.Nodes[0]) ***REMOVED***
		t.Error("Expected .row-fluid to contain a[ng-click] tag.")
	***REMOVED***
***REMOVED***

func TestSelNotContains(t *testing.T) ***REMOVED***
	sel := Doc().Find("a.link")
	sel2 := Doc().Find("span")
	if sel.Contains(sel2.Nodes[0]) ***REMOVED***
		t.Error("Expected a.link to NOT contain span tag.")
	***REMOVED***
***REMOVED***
