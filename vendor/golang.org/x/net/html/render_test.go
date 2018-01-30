// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package html

import (
	"bytes"
	"testing"
)

func TestRenderer(t *testing.T) ***REMOVED***
	nodes := [...]*Node***REMOVED***
		0: ***REMOVED***
			Type: ElementNode,
			Data: "html",
		***REMOVED***,
		1: ***REMOVED***
			Type: ElementNode,
			Data: "head",
		***REMOVED***,
		2: ***REMOVED***
			Type: ElementNode,
			Data: "body",
		***REMOVED***,
		3: ***REMOVED***
			Type: TextNode,
			Data: "0<1",
		***REMOVED***,
		4: ***REMOVED***
			Type: ElementNode,
			Data: "p",
			Attr: []Attribute***REMOVED***
				***REMOVED***
					Key: "id",
					Val: "A",
				***REMOVED***,
				***REMOVED***
					Key: "foo",
					Val: `abc"def`,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		5: ***REMOVED***
			Type: TextNode,
			Data: "2",
		***REMOVED***,
		6: ***REMOVED***
			Type: ElementNode,
			Data: "b",
			Attr: []Attribute***REMOVED***
				***REMOVED***
					Key: "empty",
					Val: "",
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		7: ***REMOVED***
			Type: TextNode,
			Data: "3",
		***REMOVED***,
		8: ***REMOVED***
			Type: ElementNode,
			Data: "i",
			Attr: []Attribute***REMOVED***
				***REMOVED***
					Key: "backslash",
					Val: `\`,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		9: ***REMOVED***
			Type: TextNode,
			Data: "&4",
		***REMOVED***,
		10: ***REMOVED***
			Type: TextNode,
			Data: "5",
		***REMOVED***,
		11: ***REMOVED***
			Type: ElementNode,
			Data: "blockquote",
		***REMOVED***,
		12: ***REMOVED***
			Type: ElementNode,
			Data: "br",
		***REMOVED***,
		13: ***REMOVED***
			Type: TextNode,
			Data: "6",
		***REMOVED***,
	***REMOVED***

	// Build a tree out of those nodes, based on a textual representation.
	// Only the ".\t"s are significant. The trailing HTML-like text is
	// just commentary. The "0:" prefixes are for easy cross-reference with
	// the nodes array.
	treeAsText := [...]string***REMOVED***
		0: `<html>`,
		1: `.	<head>`,
		2: `.	<body>`,
		3: `.	.	"0&lt;1"`,
		4: `.	.	<p id="A" foo="abc&#34;def">`,
		5: `.	.	.	"2"`,
		6: `.	.	.	<b empty="">`,
		7: `.	.	.	.	"3"`,
		8: `.	.	.	<i backslash="\">`,
		9: `.	.	.	.	"&amp;4"`,
		10: `.	.	"5"`,
		11: `.	.	<blockquote>`,
		12: `.	.	<br>`,
		13: `.	.	"6"`,
	***REMOVED***
	if len(nodes) != len(treeAsText) ***REMOVED***
		t.Fatal("len(nodes) != len(treeAsText)")
	***REMOVED***
	var stack [8]*Node
	for i, line := range treeAsText ***REMOVED***
		level := 0
		for line[0] == '.' ***REMOVED***
			// Strip a leading ".\t".
			line = line[2:]
			level++
		***REMOVED***
		n := nodes[i]
		if level == 0 ***REMOVED***
			if stack[0] != nil ***REMOVED***
				t.Fatal("multiple root nodes")
			***REMOVED***
			stack[0] = n
		***REMOVED*** else ***REMOVED***
			stack[level-1].AppendChild(n)
			stack[level] = n
			for i := level + 1; i < len(stack); i++ ***REMOVED***
				stack[i] = nil
			***REMOVED***
		***REMOVED***
		// At each stage of tree construction, we check all nodes for consistency.
		for j, m := range nodes ***REMOVED***
			if err := checkNodeConsistency(m); err != nil ***REMOVED***
				t.Fatalf("i=%d, j=%d: %v", i, j, err)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	want := `<html><head></head><body>0&lt;1<p id="A" foo="abc&#34;def">` +
		`2<b empty="">3</b><i backslash="\">&amp;4</i></p>` +
		`5<blockquote></blockquote><br/>6</body></html>`
	b := new(bytes.Buffer)
	if err := Render(b, nodes[0]); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if got := b.String(); got != want ***REMOVED***
		t.Errorf("got vs want:\n%s\n%s\n", got, want)
	***REMOVED***
***REMOVED***
