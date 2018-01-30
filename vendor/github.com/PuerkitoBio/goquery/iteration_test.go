package goquery

import (
	"testing"

	"golang.org/x/net/html"
)

func TestEach(t *testing.T) ***REMOVED***
	var cnt int

	sel := Doc().Find(".hero-unit .row-fluid").Each(func(i int, n *Selection) ***REMOVED***
		cnt++
		t.Logf("At index %v, node %v", i, n.Nodes[0].Data)
	***REMOVED***).Find("a")

	if cnt != 4 ***REMOVED***
		t.Errorf("Expected Each() to call function 4 times, got %v times.", cnt)
	***REMOVED***
	assertLength(t, sel.Nodes, 6)
***REMOVED***

func TestEachWithBreak(t *testing.T) ***REMOVED***
	var cnt int

	sel := Doc().Find(".hero-unit .row-fluid").EachWithBreak(func(i int, n *Selection) bool ***REMOVED***
		cnt++
		t.Logf("At index %v, node %v", i, n.Nodes[0].Data)
		return false
	***REMOVED***).Find("a")

	if cnt != 1 ***REMOVED***
		t.Errorf("Expected Each() to call function 1 time, got %v times.", cnt)
	***REMOVED***
	assertLength(t, sel.Nodes, 6)
***REMOVED***

func TestEachEmptySelection(t *testing.T) ***REMOVED***
	var cnt int

	sel := Doc().Find("zzzz")
	sel.Each(func(i int, n *Selection) ***REMOVED***
		cnt++
	***REMOVED***)
	if cnt > 0 ***REMOVED***
		t.Error("Expected Each() to not be called on empty Selection.")
	***REMOVED***
	sel2 := sel.Find("div")
	assertLength(t, sel2.Nodes, 0)
***REMOVED***

func TestMap(t *testing.T) ***REMOVED***
	sel := Doc().Find(".pvk-content")
	vals := sel.Map(func(i int, s *Selection) string ***REMOVED***
		n := s.Get(0)
		if n.Type == html.ElementNode ***REMOVED***
			return n.Data
		***REMOVED***
		return ""
	***REMOVED***)
	for _, v := range vals ***REMOVED***
		if v != "div" ***REMOVED***
			t.Error("Expected Map array result to be all 'div's.")
		***REMOVED***
	***REMOVED***
	if len(vals) != 3 ***REMOVED***
		t.Errorf("Expected Map array result to have a length of 3, found %v.", len(vals))
	***REMOVED***
***REMOVED***

func TestForRange(t *testing.T) ***REMOVED***
	sel := Doc().Find(".pvk-content")
	initLen := sel.Length()
	for i := range sel.Nodes ***REMOVED***
		single := sel.Eq(i)
		//h, err := single.Html()
		//if err != nil ***REMOVED***
		//	t.Fatal(err)
		//***REMOVED***
		//fmt.Println(i, h)
		if single.Length() != 1 ***REMOVED***
			t.Errorf("%d: expected length of 1, got %d", i, single.Length())
		***REMOVED***
	***REMOVED***
	if sel.Length() != initLen ***REMOVED***
		t.Errorf("expected initial selection to still have length %d, got %d", initLen, sel.Length())
	***REMOVED***
***REMOVED***
