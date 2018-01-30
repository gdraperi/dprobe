package goquery

import (
	"reflect"
	"sort"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

var allNodes = `<!doctype html>
<html>
	<head>
		<meta a="b">
	</head>
	<body>
		<p><!-- this is a comment -->
		This is some text.
		</p>
		<div></div>
		<h1 class="header"></h1>
		<h2 class="header"></h2>
	</body>
</html>`

func TestNodeName(t *testing.T) ***REMOVED***
	doc, err := NewDocumentFromReader(strings.NewReader(allNodes))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	n0 := doc.Nodes[0]
	nDT := n0.FirstChild
	sMeta := doc.Find("meta")
	nMeta := sMeta.Get(0)
	sP := doc.Find("p")
	nP := sP.Get(0)
	nComment := nP.FirstChild
	nText := nComment.NextSibling

	cases := []struct ***REMOVED***
		node *html.Node
		typ  html.NodeType
		want string
	***REMOVED******REMOVED***
		***REMOVED***n0, html.DocumentNode, nodeNames[html.DocumentNode]***REMOVED***,
		***REMOVED***nDT, html.DoctypeNode, "html"***REMOVED***,
		***REMOVED***nMeta, html.ElementNode, "meta"***REMOVED***,
		***REMOVED***nP, html.ElementNode, "p"***REMOVED***,
		***REMOVED***nComment, html.CommentNode, nodeNames[html.CommentNode]***REMOVED***,
		***REMOVED***nText, html.TextNode, nodeNames[html.TextNode]***REMOVED***,
	***REMOVED***
	for i, c := range cases ***REMOVED***
		got := NodeName(newSingleSelection(c.node, doc))
		if c.node.Type != c.typ ***REMOVED***
			t.Errorf("%d: want type %v, got %v", i, c.typ, c.node.Type)
		***REMOVED***
		if got != c.want ***REMOVED***
			t.Errorf("%d: want %q, got %q", i, c.want, got)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestNodeNameMultiSel(t *testing.T) ***REMOVED***
	doc, err := NewDocumentFromReader(strings.NewReader(allNodes))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	in := []string***REMOVED***"p", "h1", "div"***REMOVED***
	var out []string
	doc.Find(strings.Join(in, ", ")).Each(func(i int, s *Selection) ***REMOVED***
		got := NodeName(s)
		out = append(out, got)
	***REMOVED***)
	sort.Strings(in)
	sort.Strings(out)
	if !reflect.DeepEqual(in, out) ***REMOVED***
		t.Errorf("want %v, got %v", in, out)
	***REMOVED***
***REMOVED***

func TestOuterHtml(t *testing.T) ***REMOVED***
	doc, err := NewDocumentFromReader(strings.NewReader(allNodes))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	n0 := doc.Nodes[0]
	nDT := n0.FirstChild
	sMeta := doc.Find("meta")
	sP := doc.Find("p")
	nP := sP.Get(0)
	nComment := nP.FirstChild
	nText := nComment.NextSibling
	sHeaders := doc.Find(".header")

	cases := []struct ***REMOVED***
		node *html.Node
		sel  *Selection
		want string
	***REMOVED******REMOVED***
		***REMOVED***nDT, nil, "<!DOCTYPE html>"***REMOVED***, // render makes DOCTYPE all caps
		***REMOVED***nil, sMeta, `<meta a="b"/>`***REMOVED***, // and auto-closes the meta
		***REMOVED***nil, sP, `<p><!-- this is a comment -->
		This is some text.
		</p>`***REMOVED***,
		***REMOVED***nComment, nil, "<!-- this is a comment -->"***REMOVED***,
		***REMOVED***nText, nil, `
		This is some text.
		`***REMOVED***,
		***REMOVED***nil, sHeaders, `<h1 class="header"></h1>`***REMOVED***,
	***REMOVED***
	for i, c := range cases ***REMOVED***
		if c.sel == nil ***REMOVED***
			c.sel = newSingleSelection(c.node, doc)
		***REMOVED***
		got, err := OuterHtml(c.sel)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***

		if got != c.want ***REMOVED***
			t.Errorf("%d: want %q, got %q", i, c.want, got)
		***REMOVED***
	***REMOVED***
***REMOVED***
