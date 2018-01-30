// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package html

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"

	"golang.org/x/net/html/atom"
)

// readParseTest reads a single test case from r.
func readParseTest(r *bufio.Reader) (text, want, context string, err error) ***REMOVED***
	line, err := r.ReadSlice('\n')
	if err != nil ***REMOVED***
		return "", "", "", err
	***REMOVED***
	var b []byte

	// Read the HTML.
	if string(line) != "#data\n" ***REMOVED***
		return "", "", "", fmt.Errorf(`got %q want "#data\n"`, line)
	***REMOVED***
	for ***REMOVED***
		line, err = r.ReadSlice('\n')
		if err != nil ***REMOVED***
			return "", "", "", err
		***REMOVED***
		if line[0] == '#' ***REMOVED***
			break
		***REMOVED***
		b = append(b, line...)
	***REMOVED***
	text = strings.TrimSuffix(string(b), "\n")
	b = b[:0]

	// Skip the error list.
	if string(line) != "#errors\n" ***REMOVED***
		return "", "", "", fmt.Errorf(`got %q want "#errors\n"`, line)
	***REMOVED***
	for ***REMOVED***
		line, err = r.ReadSlice('\n')
		if err != nil ***REMOVED***
			return "", "", "", err
		***REMOVED***
		if line[0] == '#' ***REMOVED***
			break
		***REMOVED***
	***REMOVED***

	if string(line) == "#document-fragment\n" ***REMOVED***
		line, err = r.ReadSlice('\n')
		if err != nil ***REMOVED***
			return "", "", "", err
		***REMOVED***
		context = strings.TrimSpace(string(line))
		line, err = r.ReadSlice('\n')
		if err != nil ***REMOVED***
			return "", "", "", err
		***REMOVED***
	***REMOVED***

	// Read the dump of what the parse tree should be.
	if string(line) != "#document\n" ***REMOVED***
		return "", "", "", fmt.Errorf(`got %q want "#document\n"`, line)
	***REMOVED***
	inQuote := false
	for ***REMOVED***
		line, err = r.ReadSlice('\n')
		if err != nil && err != io.EOF ***REMOVED***
			return "", "", "", err
		***REMOVED***
		trimmed := bytes.Trim(line, "| \n")
		if len(trimmed) > 0 ***REMOVED***
			if line[0] == '|' && trimmed[0] == '"' ***REMOVED***
				inQuote = true
			***REMOVED***
			if trimmed[len(trimmed)-1] == '"' && !(line[0] == '|' && len(trimmed) == 1) ***REMOVED***
				inQuote = false
			***REMOVED***
		***REMOVED***
		if len(line) == 0 || len(line) == 1 && line[0] == '\n' && !inQuote ***REMOVED***
			break
		***REMOVED***
		b = append(b, line...)
	***REMOVED***
	return text, string(b), context, nil
***REMOVED***

func dumpIndent(w io.Writer, level int) ***REMOVED***
	io.WriteString(w, "| ")
	for i := 0; i < level; i++ ***REMOVED***
		io.WriteString(w, "  ")
	***REMOVED***
***REMOVED***

type sortedAttributes []Attribute

func (a sortedAttributes) Len() int ***REMOVED***
	return len(a)
***REMOVED***

func (a sortedAttributes) Less(i, j int) bool ***REMOVED***
	if a[i].Namespace != a[j].Namespace ***REMOVED***
		return a[i].Namespace < a[j].Namespace
	***REMOVED***
	return a[i].Key < a[j].Key
***REMOVED***

func (a sortedAttributes) Swap(i, j int) ***REMOVED***
	a[i], a[j] = a[j], a[i]
***REMOVED***

func dumpLevel(w io.Writer, n *Node, level int) error ***REMOVED***
	dumpIndent(w, level)
	switch n.Type ***REMOVED***
	case ErrorNode:
		return errors.New("unexpected ErrorNode")
	case DocumentNode:
		return errors.New("unexpected DocumentNode")
	case ElementNode:
		if n.Namespace != "" ***REMOVED***
			fmt.Fprintf(w, "<%s %s>", n.Namespace, n.Data)
		***REMOVED*** else ***REMOVED***
			fmt.Fprintf(w, "<%s>", n.Data)
		***REMOVED***
		attr := sortedAttributes(n.Attr)
		sort.Sort(attr)
		for _, a := range attr ***REMOVED***
			io.WriteString(w, "\n")
			dumpIndent(w, level+1)
			if a.Namespace != "" ***REMOVED***
				fmt.Fprintf(w, `%s %s="%s"`, a.Namespace, a.Key, a.Val)
			***REMOVED*** else ***REMOVED***
				fmt.Fprintf(w, `%s="%s"`, a.Key, a.Val)
			***REMOVED***
		***REMOVED***
	case TextNode:
		fmt.Fprintf(w, `"%s"`, n.Data)
	case CommentNode:
		fmt.Fprintf(w, "<!-- %s -->", n.Data)
	case DoctypeNode:
		fmt.Fprintf(w, "<!DOCTYPE %s", n.Data)
		if n.Attr != nil ***REMOVED***
			var p, s string
			for _, a := range n.Attr ***REMOVED***
				switch a.Key ***REMOVED***
				case "public":
					p = a.Val
				case "system":
					s = a.Val
				***REMOVED***
			***REMOVED***
			if p != "" || s != "" ***REMOVED***
				fmt.Fprintf(w, ` "%s"`, p)
				fmt.Fprintf(w, ` "%s"`, s)
			***REMOVED***
		***REMOVED***
		io.WriteString(w, ">")
	case scopeMarkerNode:
		return errors.New("unexpected scopeMarkerNode")
	default:
		return errors.New("unknown node type")
	***REMOVED***
	io.WriteString(w, "\n")
	for c := n.FirstChild; c != nil; c = c.NextSibling ***REMOVED***
		if err := dumpLevel(w, c, level+1); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func dump(n *Node) (string, error) ***REMOVED***
	if n == nil || n.FirstChild == nil ***REMOVED***
		return "", nil
	***REMOVED***
	var b bytes.Buffer
	for c := n.FirstChild; c != nil; c = c.NextSibling ***REMOVED***
		if err := dumpLevel(&b, c, 0); err != nil ***REMOVED***
			return "", err
		***REMOVED***
	***REMOVED***
	return b.String(), nil
***REMOVED***

const testDataDir = "testdata/webkit/"

func TestParser(t *testing.T) ***REMOVED***
	testFiles, err := filepath.Glob(testDataDir + "*.dat")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	for _, tf := range testFiles ***REMOVED***
		f, err := os.Open(tf)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		defer f.Close()
		r := bufio.NewReader(f)

		for i := 0; ; i++ ***REMOVED***
			text, want, context, err := readParseTest(r)
			if err == io.EOF ***REMOVED***
				break
			***REMOVED***
			if err != nil ***REMOVED***
				t.Fatal(err)
			***REMOVED***

			err = testParseCase(text, want, context)

			if err != nil ***REMOVED***
				t.Errorf("%s test #%d %q, %s", tf, i, text, err)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

// testParseCase tests one test case from the test files. If the test does not
// pass, it returns an error that explains the failure.
// text is the HTML to be parsed, want is a dump of the correct parse tree,
// and context is the name of the context node, if any.
func testParseCase(text, want, context string) (err error) ***REMOVED***
	defer func() ***REMOVED***
		if x := recover(); x != nil ***REMOVED***
			switch e := x.(type) ***REMOVED***
			case error:
				err = e
			default:
				err = fmt.Errorf("%v", e)
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	var doc *Node
	if context == "" ***REMOVED***
		doc, err = Parse(strings.NewReader(text))
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		contextNode := &Node***REMOVED***
			Type:     ElementNode,
			DataAtom: atom.Lookup([]byte(context)),
			Data:     context,
		***REMOVED***
		nodes, err := ParseFragment(strings.NewReader(text), contextNode)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		doc = &Node***REMOVED***
			Type: DocumentNode,
		***REMOVED***
		for _, n := range nodes ***REMOVED***
			doc.AppendChild(n)
		***REMOVED***
	***REMOVED***

	if err := checkTreeConsistency(doc); err != nil ***REMOVED***
		return err
	***REMOVED***

	got, err := dump(doc)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	// Compare the parsed tree to the #document section.
	if got != want ***REMOVED***
		return fmt.Errorf("got vs want:\n----\n%s----\n%s----", got, want)
	***REMOVED***

	if renderTestBlacklist[text] || context != "" ***REMOVED***
		return nil
	***REMOVED***

	// Check that rendering and re-parsing results in an identical tree.
	pr, pw := io.Pipe()
	go func() ***REMOVED***
		pw.CloseWithError(Render(pw, doc))
	***REMOVED***()
	doc1, err := Parse(pr)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	got1, err := dump(doc1)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if got != got1 ***REMOVED***
		return fmt.Errorf("got vs got1:\n----\n%s----\n%s----", got, got1)
	***REMOVED***

	return nil
***REMOVED***

// Some test input result in parse trees are not 'well-formed' despite
// following the HTML5 recovery algorithms. Rendering and re-parsing such a
// tree will not result in an exact clone of that tree. We blacklist such
// inputs from the render test.
var renderTestBlacklist = map[string]bool***REMOVED***
	// The second <a> will be reparented to the first <table>'s parent. This
	// results in an <a> whose parent is an <a>, which is not 'well-formed'.
	`<a><table><td><a><table></table><a></tr><a></table><b>X</b>C<a>Y`: true,
	// The same thing with a <p>:
	`<p><table></p>`: true,
	// More cases of <a> being reparented:
	`<a href="blah">aba<table><a href="foo">br<tr><td></td></tr>x</table>aoe`: true,
	`<a><table><a></table><p><a><div><a>`:                                     true,
	`<a><table><td><a><table></table><a></tr><a></table><a>`:                  true,
	// A similar reparenting situation involving <nobr>:
	`<!DOCTYPE html><body><b><nobr>1<table><nobr></b><i><nobr>2<nobr></i>3`: true,
	// A <plaintext> element is reparented, putting it before a table.
	// A <plaintext> element can't have anything after it in HTML.
	`<table><plaintext><td>`:                                   true,
	`<!doctype html><table><plaintext></plaintext>`:            true,
	`<!doctype html><table><tbody><plaintext></plaintext>`:     true,
	`<!doctype html><table><tbody><tr><plaintext></plaintext>`: true,
	// A form inside a table inside a form doesn't work either.
	`<!doctype html><form><table></form><form></table></form>`: true,
	// A script that ends at EOF may escape its own closing tag when rendered.
	`<!doctype html><script><!--<script `:          true,
	`<!doctype html><script><!--<script <`:         true,
	`<!doctype html><script><!--<script <a`:        true,
	`<!doctype html><script><!--<script </`:        true,
	`<!doctype html><script><!--<script </s`:       true,
	`<!doctype html><script><!--<script </script`:  true,
	`<!doctype html><script><!--<script </scripta`: true,
	`<!doctype html><script><!--<script -`:         true,
	`<!doctype html><script><!--<script -a`:        true,
	`<!doctype html><script><!--<script -<`:        true,
	`<!doctype html><script><!--<script --`:        true,
	`<!doctype html><script><!--<script --a`:       true,
	`<!doctype html><script><!--<script --<`:       true,
	`<script><!--<script `:                         true,
	`<script><!--<script <a`:                       true,
	`<script><!--<script </script`:                 true,
	`<script><!--<script </scripta`:                true,
	`<script><!--<script -`:                        true,
	`<script><!--<script -a`:                       true,
	`<script><!--<script --`:                       true,
	`<script><!--<script --a`:                      true,
	`<script><!--<script <`:                        true,
	`<script><!--<script </`:                       true,
	`<script><!--<script </s`:                      true,
	// Reconstructing the active formatting elements results in a <plaintext>
	// element that contains an <a> element.
	`<!doctype html><p><a><plaintext>b`: true,
***REMOVED***

func TestNodeConsistency(t *testing.T) ***REMOVED***
	// inconsistentNode is a Node whose DataAtom and Data do not agree.
	inconsistentNode := &Node***REMOVED***
		Type:     ElementNode,
		DataAtom: atom.Frameset,
		Data:     "table",
	***REMOVED***
	_, err := ParseFragment(strings.NewReader("<p>hello</p>"), inconsistentNode)
	if err == nil ***REMOVED***
		t.Errorf("got nil error, want non-nil")
	***REMOVED***
***REMOVED***

func BenchmarkParser(b *testing.B) ***REMOVED***
	buf, err := ioutil.ReadFile("testdata/go1.html")
	if err != nil ***REMOVED***
		b.Fatalf("could not read testdata/go1.html: %v", err)
	***REMOVED***
	b.SetBytes(int64(len(buf)))
	runtime.GC()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		Parse(bytes.NewBuffer(buf))
	***REMOVED***
***REMOVED***
