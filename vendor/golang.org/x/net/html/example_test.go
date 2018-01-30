// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This example demonstrates parsing HTML data and walking the resulting tree.
package html_test

import (
	"fmt"
	"log"
	"strings"

	"golang.org/x/net/html"
)

func ExampleParse() ***REMOVED***
	s := `<p>Links:</p><ul><li><a href="foo">Foo</a><li><a href="/bar/baz">BarBaz</a></ul>`
	doc, err := html.Parse(strings.NewReader(s))
	if err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***
	var f func(*html.Node)
	f = func(n *html.Node) ***REMOVED***
		if n.Type == html.ElementNode && n.Data == "a" ***REMOVED***
			for _, a := range n.Attr ***REMOVED***
				if a.Key == "href" ***REMOVED***
					fmt.Println(a.Val)
					break
				***REMOVED***
			***REMOVED***
		***REMOVED***
		for c := n.FirstChild; c != nil; c = c.NextSibling ***REMOVED***
			f(c)
		***REMOVED***
	***REMOVED***
	f(doc)
	// Output:
	// foo
	// /bar/baz
***REMOVED***
