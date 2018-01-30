package goquery

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

// Test helper functions and members
var doc *Document
var doc2 *Document
var doc3 *Document
var docB *Document
var docW *Document

func Doc() *Document ***REMOVED***
	if doc == nil ***REMOVED***
		doc = loadDoc("page.html")
	***REMOVED***
	return doc
***REMOVED***

func Doc2() *Document ***REMOVED***
	if doc2 == nil ***REMOVED***
		doc2 = loadDoc("page2.html")
	***REMOVED***
	return doc2
***REMOVED***

func Doc2Clone() *Document ***REMOVED***
	return CloneDocument(Doc2())
***REMOVED***

func Doc3() *Document ***REMOVED***
	if doc3 == nil ***REMOVED***
		doc3 = loadDoc("page3.html")
	***REMOVED***
	return doc3
***REMOVED***

func Doc3Clone() *Document ***REMOVED***
	return CloneDocument(Doc3())
***REMOVED***

func DocB() *Document ***REMOVED***
	if docB == nil ***REMOVED***
		docB = loadDoc("gotesting.html")
	***REMOVED***
	return docB
***REMOVED***

func DocW() *Document ***REMOVED***
	if docW == nil ***REMOVED***
		docW = loadDoc("gowiki.html")
	***REMOVED***
	return docW
***REMOVED***

func assertLength(t *testing.T, nodes []*html.Node, length int) ***REMOVED***
	if len(nodes) != length ***REMOVED***
		t.Errorf("Expected %d nodes, found %d.", length, len(nodes))
		for i, n := range nodes ***REMOVED***
			t.Logf("Node %d: %+v.", i, n)
		***REMOVED***
	***REMOVED***
***REMOVED***

func assertClass(t *testing.T, sel *Selection, class string) ***REMOVED***
	if !sel.HasClass(class) ***REMOVED***
		t.Errorf("Expected node to have class %s, found %+v.", class, sel.Get(0))
	***REMOVED***
***REMOVED***

func assertPanic(t *testing.T) ***REMOVED***
	if e := recover(); e == nil ***REMOVED***
		t.Error("Expected a panic.")
	***REMOVED***
***REMOVED***

func assertEqual(t *testing.T, s1 *Selection, s2 *Selection) ***REMOVED***
	if s1 != s2 ***REMOVED***
		t.Error("Expected selection objects to be the same.")
	***REMOVED***
***REMOVED***

func assertSelectionIs(t *testing.T, sel *Selection, is ...string) ***REMOVED***
	for i := 0; i < sel.Length(); i++ ***REMOVED***
		if !sel.Eq(i).Is(is[i]) ***REMOVED***
			t.Errorf("Expected node %d to be %s, found %+v", i, is[i], sel.Get(i))
		***REMOVED***
	***REMOVED***
***REMOVED***

func printSel(t *testing.T, sel *Selection) ***REMOVED***
	if testing.Verbose() ***REMOVED***
		h, err := sel.Html()
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		t.Log(h)
	***REMOVED***
***REMOVED***

func loadDoc(page string) *Document ***REMOVED***
	var f *os.File
	var e error

	if f, e = os.Open(fmt.Sprintf("./testdata/%s", page)); e != nil ***REMOVED***
		panic(e.Error())
	***REMOVED***
	defer f.Close()

	var node *html.Node
	if node, e = html.Parse(f); e != nil ***REMOVED***
		panic(e.Error())
	***REMOVED***
	return NewDocumentFromNode(node)
***REMOVED***

func TestNewDocument(t *testing.T) ***REMOVED***
	if f, e := os.Open("./testdata/page.html"); e != nil ***REMOVED***
		t.Error(e.Error())
	***REMOVED*** else ***REMOVED***
		defer f.Close()
		if node, e := html.Parse(f); e != nil ***REMOVED***
			t.Error(e.Error())
		***REMOVED*** else ***REMOVED***
			doc = NewDocumentFromNode(node)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestNewDocumentFromReader(t *testing.T) ***REMOVED***
	cases := []struct ***REMOVED***
		src string
		err bool
		sel string
		cnt int
	***REMOVED******REMOVED***
		0: ***REMOVED***
			src: `
<html>
<head>
<title>Test</title>
<body>
<h1>Hi</h1>
</body>
</html>`,
			sel: "h1",
			cnt: 1,
		***REMOVED***,
		1: ***REMOVED***
			// Actually pretty hard to make html.Parse return an error
			// based on content...
			src: `<html><body><aef<eqf>>>qq></body></ht>`,
		***REMOVED***,
	***REMOVED***
	buf := bytes.NewBuffer(nil)

	for i, c := range cases ***REMOVED***
		buf.Reset()
		buf.WriteString(c.src)

		d, e := NewDocumentFromReader(buf)
		if (e != nil) != c.err ***REMOVED***
			if c.err ***REMOVED***
				t.Errorf("[%d] - expected error, got none", i)
			***REMOVED*** else ***REMOVED***
				t.Errorf("[%d] - expected no error, got %s", i, e)
			***REMOVED***
		***REMOVED***
		if c.sel != "" ***REMOVED***
			s := d.Find(c.sel)
			if s.Length() != c.cnt ***REMOVED***
				t.Errorf("[%d] - expected %d nodes, found %d", i, c.cnt, s.Length())
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestNewDocumentFromResponseNil(t *testing.T) ***REMOVED***
	_, e := NewDocumentFromResponse(nil)
	if e == nil ***REMOVED***
		t.Error("Expected error, got none")
	***REMOVED***
***REMOVED***

func TestIssue103(t *testing.T) ***REMOVED***
	d, err := NewDocumentFromReader(strings.NewReader("<html><title>Scientists Stored These Images in DNA—Then Flawlessly Retrieved Them</title></html>"))
	if err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***
	text := d.Find("title").Text()
	for i, r := range text ***REMOVED***
		t.Logf("%d: %d - %q\n", i, r, string(r))
	***REMOVED***
	t.Log(text)
***REMOVED***
