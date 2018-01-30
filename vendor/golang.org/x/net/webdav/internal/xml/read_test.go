// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xml

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"
	"time"
)

// Stripped down Atom feed data structures.

func TestUnmarshalFeed(t *testing.T) ***REMOVED***
	var f Feed
	if err := Unmarshal([]byte(atomFeedString), &f); err != nil ***REMOVED***
		t.Fatalf("Unmarshal: %s", err)
	***REMOVED***
	if !reflect.DeepEqual(f, atomFeed) ***REMOVED***
		t.Fatalf("have %#v\nwant %#v", f, atomFeed)
	***REMOVED***
***REMOVED***

// hget http://codereview.appspot.com/rss/mine/rsc
const atomFeedString = `
<?xml version="1.0" encoding="utf-8"?>
<feed xmlns="http://www.w3.org/2005/Atom" xml:lang="en-us" updated="2009-10-04T01:35:58+00:00"><title>Code Review - My issues</title><link href="http://codereview.appspot.com/" rel="alternate"></link><link href="http://codereview.appspot.com/rss/mine/rsc" rel="self"></link><id>http://codereview.appspot.com/</id><author><name>rietveld&lt;&gt;</name></author><entry><title>rietveld: an attempt at pubsubhubbub
</title><link href="http://codereview.appspot.com/126085" rel="alternate"></link><updated>2009-10-04T01:35:58+00:00</updated><author><name>email-address-removed</name></author><id>urn:md5:134d9179c41f806be79b3a5f7877d19a</id><summary type="html">
  An attempt at adding pubsubhubbub support to Rietveld.
http://code.google.com/p/pubsubhubbub
http://code.google.com/p/rietveld/issues/detail?id=155

The server side of the protocol is trivial:
  1. add a &amp;lt;link rel=&amp;quot;hub&amp;quot; href=&amp;quot;hub-server&amp;quot;&amp;gt; tag to all
     feeds that will be pubsubhubbubbed.
  2. every time one of those feeds changes, tell the hub
     with a simple POST request.

I have tested this by adding debug prints to a local hub
server and checking that the server got the right publish
requests.

I can&amp;#39;t quite get the server to work, but I think the bug
is not in my code.  I think that the server expects to be
able to grab the feed and see the feed&amp;#39;s actual URL in
the link rel=&amp;quot;self&amp;quot;, but the default value for that drops
the :port from the URL, and I cannot for the life of me
figure out how to get the Atom generator deep inside
django not to do that, or even where it is doing that,
or even what code is running to generate the Atom feed.
(I thought I knew but I added some assert False statements
and it kept running!)

Ignoring that particular problem, I would appreciate
feedback on the right way to get the two values at
the top of feeds.py marked NOTE(rsc).


</summary></entry><entry><title>rietveld: correct tab handling
</title><link href="http://codereview.appspot.com/124106" rel="alternate"></link><updated>2009-10-03T23:02:17+00:00</updated><author><name>email-address-removed</name></author><id>urn:md5:0a2a4f19bb815101f0ba2904aed7c35a</id><summary type="html">
  This fixes the buggy tab rendering that can be seen at
http://codereview.appspot.com/116075/diff/1/2

The fundamental problem was that the tab code was
not being told what column the text began in, so it
didn&amp;#39;t know where to put the tab stops.  Another problem
was that some of the code assumed that string byte
offsets were the same as column offsets, which is only
true if there are no tabs.

In the process of fixing this, I cleaned up the arguments
to Fold and ExpandTabs and renamed them Break and
_ExpandTabs so that I could be sure that I found all the
call sites.  I also wanted to verify that ExpandTabs was
not being used from outside intra_region_diff.py.


</summary></entry></feed> 	   `

type Feed struct ***REMOVED***
	XMLName Name      `xml:"http://www.w3.org/2005/Atom feed"`
	Title   string    `xml:"title"`
	Id      string    `xml:"id"`
	Link    []Link    `xml:"link"`
	Updated time.Time `xml:"updated,attr"`
	Author  Person    `xml:"author"`
	Entry   []Entry   `xml:"entry"`
***REMOVED***

type Entry struct ***REMOVED***
	Title   string    `xml:"title"`
	Id      string    `xml:"id"`
	Link    []Link    `xml:"link"`
	Updated time.Time `xml:"updated"`
	Author  Person    `xml:"author"`
	Summary Text      `xml:"summary"`
***REMOVED***

type Link struct ***REMOVED***
	Rel  string `xml:"rel,attr,omitempty"`
	Href string `xml:"href,attr"`
***REMOVED***

type Person struct ***REMOVED***
	Name     string `xml:"name"`
	URI      string `xml:"uri"`
	Email    string `xml:"email"`
	InnerXML string `xml:",innerxml"`
***REMOVED***

type Text struct ***REMOVED***
	Type string `xml:"type,attr,omitempty"`
	Body string `xml:",chardata"`
***REMOVED***

var atomFeed = Feed***REMOVED***
	XMLName: Name***REMOVED***"http://www.w3.org/2005/Atom", "feed"***REMOVED***,
	Title:   "Code Review - My issues",
	Link: []Link***REMOVED***
		***REMOVED***Rel: "alternate", Href: "http://codereview.appspot.com/"***REMOVED***,
		***REMOVED***Rel: "self", Href: "http://codereview.appspot.com/rss/mine/rsc"***REMOVED***,
	***REMOVED***,
	Id:      "http://codereview.appspot.com/",
	Updated: ParseTime("2009-10-04T01:35:58+00:00"),
	Author: Person***REMOVED***
		Name:     "rietveld<>",
		InnerXML: "<name>rietveld&lt;&gt;</name>",
	***REMOVED***,
	Entry: []Entry***REMOVED***
		***REMOVED***
			Title: "rietveld: an attempt at pubsubhubbub\n",
			Link: []Link***REMOVED***
				***REMOVED***Rel: "alternate", Href: "http://codereview.appspot.com/126085"***REMOVED***,
			***REMOVED***,
			Updated: ParseTime("2009-10-04T01:35:58+00:00"),
			Author: Person***REMOVED***
				Name:     "email-address-removed",
				InnerXML: "<name>email-address-removed</name>",
			***REMOVED***,
			Id: "urn:md5:134d9179c41f806be79b3a5f7877d19a",
			Summary: Text***REMOVED***
				Type: "html",
				Body: `
  An attempt at adding pubsubhubbub support to Rietveld.
http://code.google.com/p/pubsubhubbub
http://code.google.com/p/rietveld/issues/detail?id=155

The server side of the protocol is trivial:
  1. add a &lt;link rel=&quot;hub&quot; href=&quot;hub-server&quot;&gt; tag to all
     feeds that will be pubsubhubbubbed.
  2. every time one of those feeds changes, tell the hub
     with a simple POST request.

I have tested this by adding debug prints to a local hub
server and checking that the server got the right publish
requests.

I can&#39;t quite get the server to work, but I think the bug
is not in my code.  I think that the server expects to be
able to grab the feed and see the feed&#39;s actual URL in
the link rel=&quot;self&quot;, but the default value for that drops
the :port from the URL, and I cannot for the life of me
figure out how to get the Atom generator deep inside
django not to do that, or even where it is doing that,
or even what code is running to generate the Atom feed.
(I thought I knew but I added some assert False statements
and it kept running!)

Ignoring that particular problem, I would appreciate
feedback on the right way to get the two values at
the top of feeds.py marked NOTE(rsc).


`,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			Title: "rietveld: correct tab handling\n",
			Link: []Link***REMOVED***
				***REMOVED***Rel: "alternate", Href: "http://codereview.appspot.com/124106"***REMOVED***,
			***REMOVED***,
			Updated: ParseTime("2009-10-03T23:02:17+00:00"),
			Author: Person***REMOVED***
				Name:     "email-address-removed",
				InnerXML: "<name>email-address-removed</name>",
			***REMOVED***,
			Id: "urn:md5:0a2a4f19bb815101f0ba2904aed7c35a",
			Summary: Text***REMOVED***
				Type: "html",
				Body: `
  This fixes the buggy tab rendering that can be seen at
http://codereview.appspot.com/116075/diff/1/2

The fundamental problem was that the tab code was
not being told what column the text began in, so it
didn&#39;t know where to put the tab stops.  Another problem
was that some of the code assumed that string byte
offsets were the same as column offsets, which is only
true if there are no tabs.

In the process of fixing this, I cleaned up the arguments
to Fold and ExpandTabs and renamed them Break and
_ExpandTabs so that I could be sure that I found all the
call sites.  I also wanted to verify that ExpandTabs was
not being used from outside intra_region_diff.py.


`,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***,
***REMOVED***

const pathTestString = `
<Result>
    <Before>1</Before>
    <Items>
        <Item1>
            <Value>A</Value>
        </Item1>
        <Item2>
            <Value>B</Value>
        </Item2>
        <Item1>
            <Value>C</Value>
            <Value>D</Value>
        </Item1>
        <_>
            <Value>E</Value>
        </_>
    </Items>
    <After>2</After>
</Result>
`

type PathTestItem struct ***REMOVED***
	Value string
***REMOVED***

type PathTestA struct ***REMOVED***
	Items         []PathTestItem `xml:">Item1"`
	Before, After string
***REMOVED***

type PathTestB struct ***REMOVED***
	Other         []PathTestItem `xml:"Items>Item1"`
	Before, After string
***REMOVED***

type PathTestC struct ***REMOVED***
	Values1       []string `xml:"Items>Item1>Value"`
	Values2       []string `xml:"Items>Item2>Value"`
	Before, After string
***REMOVED***

type PathTestSet struct ***REMOVED***
	Item1 []PathTestItem
***REMOVED***

type PathTestD struct ***REMOVED***
	Other         PathTestSet `xml:"Items"`
	Before, After string
***REMOVED***

type PathTestE struct ***REMOVED***
	Underline     string `xml:"Items>_>Value"`
	Before, After string
***REMOVED***

var pathTests = []interface***REMOVED******REMOVED******REMOVED***
	&PathTestA***REMOVED***Items: []PathTestItem***REMOVED******REMOVED***"A"***REMOVED***, ***REMOVED***"D"***REMOVED******REMOVED***, Before: "1", After: "2"***REMOVED***,
	&PathTestB***REMOVED***Other: []PathTestItem***REMOVED******REMOVED***"A"***REMOVED***, ***REMOVED***"D"***REMOVED******REMOVED***, Before: "1", After: "2"***REMOVED***,
	&PathTestC***REMOVED***Values1: []string***REMOVED***"A", "C", "D"***REMOVED***, Values2: []string***REMOVED***"B"***REMOVED***, Before: "1", After: "2"***REMOVED***,
	&PathTestD***REMOVED***Other: PathTestSet***REMOVED***Item1: []PathTestItem***REMOVED******REMOVED***"A"***REMOVED***, ***REMOVED***"D"***REMOVED******REMOVED******REMOVED***, Before: "1", After: "2"***REMOVED***,
	&PathTestE***REMOVED***Underline: "E", Before: "1", After: "2"***REMOVED***,
***REMOVED***

func TestUnmarshalPaths(t *testing.T) ***REMOVED***
	for _, pt := range pathTests ***REMOVED***
		v := reflect.New(reflect.TypeOf(pt).Elem()).Interface()
		if err := Unmarshal([]byte(pathTestString), v); err != nil ***REMOVED***
			t.Fatalf("Unmarshal: %s", err)
		***REMOVED***
		if !reflect.DeepEqual(v, pt) ***REMOVED***
			t.Fatalf("have %#v\nwant %#v", v, pt)
		***REMOVED***
	***REMOVED***
***REMOVED***

type BadPathTestA struct ***REMOVED***
	First  string `xml:"items>item1"`
	Other  string `xml:"items>item2"`
	Second string `xml:"items"`
***REMOVED***

type BadPathTestB struct ***REMOVED***
	Other  string `xml:"items>item2>value"`
	First  string `xml:"items>item1"`
	Second string `xml:"items>item1>value"`
***REMOVED***

type BadPathTestC struct ***REMOVED***
	First  string
	Second string `xml:"First"`
***REMOVED***

type BadPathTestD struct ***REMOVED***
	BadPathEmbeddedA
	BadPathEmbeddedB
***REMOVED***

type BadPathEmbeddedA struct ***REMOVED***
	First string
***REMOVED***

type BadPathEmbeddedB struct ***REMOVED***
	Second string `xml:"First"`
***REMOVED***

var badPathTests = []struct ***REMOVED***
	v, e interface***REMOVED******REMOVED***
***REMOVED******REMOVED***
	***REMOVED***&BadPathTestA***REMOVED******REMOVED***, &TagPathError***REMOVED***reflect.TypeOf(BadPathTestA***REMOVED******REMOVED***), "First", "items>item1", "Second", "items"***REMOVED******REMOVED***,
	***REMOVED***&BadPathTestB***REMOVED******REMOVED***, &TagPathError***REMOVED***reflect.TypeOf(BadPathTestB***REMOVED******REMOVED***), "First", "items>item1", "Second", "items>item1>value"***REMOVED******REMOVED***,
	***REMOVED***&BadPathTestC***REMOVED******REMOVED***, &TagPathError***REMOVED***reflect.TypeOf(BadPathTestC***REMOVED******REMOVED***), "First", "", "Second", "First"***REMOVED******REMOVED***,
	***REMOVED***&BadPathTestD***REMOVED******REMOVED***, &TagPathError***REMOVED***reflect.TypeOf(BadPathTestD***REMOVED******REMOVED***), "First", "", "Second", "First"***REMOVED******REMOVED***,
***REMOVED***

func TestUnmarshalBadPaths(t *testing.T) ***REMOVED***
	for _, tt := range badPathTests ***REMOVED***
		err := Unmarshal([]byte(pathTestString), tt.v)
		if !reflect.DeepEqual(err, tt.e) ***REMOVED***
			t.Fatalf("Unmarshal with %#v didn't fail properly:\nhave %#v,\nwant %#v", tt.v, err, tt.e)
		***REMOVED***
	***REMOVED***
***REMOVED***

const OK = "OK"
const withoutNameTypeData = `
<?xml version="1.0" charset="utf-8"?>
<Test3 Attr="OK" />`

type TestThree struct ***REMOVED***
	XMLName Name   `xml:"Test3"`
	Attr    string `xml:",attr"`
***REMOVED***

func TestUnmarshalWithoutNameType(t *testing.T) ***REMOVED***
	var x TestThree
	if err := Unmarshal([]byte(withoutNameTypeData), &x); err != nil ***REMOVED***
		t.Fatalf("Unmarshal: %s", err)
	***REMOVED***
	if x.Attr != OK ***REMOVED***
		t.Fatalf("have %v\nwant %v", x.Attr, OK)
	***REMOVED***
***REMOVED***

func TestUnmarshalAttr(t *testing.T) ***REMOVED***
	type ParamVal struct ***REMOVED***
		Int int `xml:"int,attr"`
	***REMOVED***

	type ParamPtr struct ***REMOVED***
		Int *int `xml:"int,attr"`
	***REMOVED***

	type ParamStringPtr struct ***REMOVED***
		Int *string `xml:"int,attr"`
	***REMOVED***

	x := []byte(`<Param int="1" />`)

	p1 := &ParamPtr***REMOVED******REMOVED***
	if err := Unmarshal(x, p1); err != nil ***REMOVED***
		t.Fatalf("Unmarshal: %s", err)
	***REMOVED***
	if p1.Int == nil ***REMOVED***
		t.Fatalf("Unmarshal failed in to *int field")
	***REMOVED*** else if *p1.Int != 1 ***REMOVED***
		t.Fatalf("Unmarshal with %s failed:\nhave %#v,\n want %#v", x, p1.Int, 1)
	***REMOVED***

	p2 := &ParamVal***REMOVED******REMOVED***
	if err := Unmarshal(x, p2); err != nil ***REMOVED***
		t.Fatalf("Unmarshal: %s", err)
	***REMOVED***
	if p2.Int != 1 ***REMOVED***
		t.Fatalf("Unmarshal with %s failed:\nhave %#v,\n want %#v", x, p2.Int, 1)
	***REMOVED***

	p3 := &ParamStringPtr***REMOVED******REMOVED***
	if err := Unmarshal(x, p3); err != nil ***REMOVED***
		t.Fatalf("Unmarshal: %s", err)
	***REMOVED***
	if p3.Int == nil ***REMOVED***
		t.Fatalf("Unmarshal failed in to *string field")
	***REMOVED*** else if *p3.Int != "1" ***REMOVED***
		t.Fatalf("Unmarshal with %s failed:\nhave %#v,\n want %#v", x, p3.Int, 1)
	***REMOVED***
***REMOVED***

type Tables struct ***REMOVED***
	HTable string `xml:"http://www.w3.org/TR/html4/ table"`
	FTable string `xml:"http://www.w3schools.com/furniture table"`
***REMOVED***

var tables = []struct ***REMOVED***
	xml string
	tab Tables
	ns  string
***REMOVED******REMOVED***
	***REMOVED***
		xml: `<Tables>` +
			`<table xmlns="http://www.w3.org/TR/html4/">hello</table>` +
			`<table xmlns="http://www.w3schools.com/furniture">world</table>` +
			`</Tables>`,
		tab: Tables***REMOVED***"hello", "world"***REMOVED***,
	***REMOVED***,
	***REMOVED***
		xml: `<Tables>` +
			`<table xmlns="http://www.w3schools.com/furniture">world</table>` +
			`<table xmlns="http://www.w3.org/TR/html4/">hello</table>` +
			`</Tables>`,
		tab: Tables***REMOVED***"hello", "world"***REMOVED***,
	***REMOVED***,
	***REMOVED***
		xml: `<Tables xmlns:f="http://www.w3schools.com/furniture" xmlns:h="http://www.w3.org/TR/html4/">` +
			`<f:table>world</f:table>` +
			`<h:table>hello</h:table>` +
			`</Tables>`,
		tab: Tables***REMOVED***"hello", "world"***REMOVED***,
	***REMOVED***,
	***REMOVED***
		xml: `<Tables>` +
			`<table>bogus</table>` +
			`</Tables>`,
		tab: Tables***REMOVED******REMOVED***,
	***REMOVED***,
	***REMOVED***
		xml: `<Tables>` +
			`<table>only</table>` +
			`</Tables>`,
		tab: Tables***REMOVED***HTable: "only"***REMOVED***,
		ns:  "http://www.w3.org/TR/html4/",
	***REMOVED***,
	***REMOVED***
		xml: `<Tables>` +
			`<table>only</table>` +
			`</Tables>`,
		tab: Tables***REMOVED***FTable: "only"***REMOVED***,
		ns:  "http://www.w3schools.com/furniture",
	***REMOVED***,
	***REMOVED***
		xml: `<Tables>` +
			`<table>only</table>` +
			`</Tables>`,
		tab: Tables***REMOVED******REMOVED***,
		ns:  "something else entirely",
	***REMOVED***,
***REMOVED***

func TestUnmarshalNS(t *testing.T) ***REMOVED***
	for i, tt := range tables ***REMOVED***
		var dst Tables
		var err error
		if tt.ns != "" ***REMOVED***
			d := NewDecoder(strings.NewReader(tt.xml))
			d.DefaultSpace = tt.ns
			err = d.Decode(&dst)
		***REMOVED*** else ***REMOVED***
			err = Unmarshal([]byte(tt.xml), &dst)
		***REMOVED***
		if err != nil ***REMOVED***
			t.Errorf("#%d: Unmarshal: %v", i, err)
			continue
		***REMOVED***
		want := tt.tab
		if dst != want ***REMOVED***
			t.Errorf("#%d: dst=%+v, want %+v", i, dst, want)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestRoundTrip(t *testing.T) ***REMOVED***
	// From issue 7535
	const s = `<ex:element xmlns:ex="http://example.com/schema"></ex:element>`
	in := bytes.NewBufferString(s)
	for i := 0; i < 10; i++ ***REMOVED***
		out := &bytes.Buffer***REMOVED******REMOVED***
		d := NewDecoder(in)
		e := NewEncoder(out)

		for ***REMOVED***
			t, err := d.Token()
			if err == io.EOF ***REMOVED***
				break
			***REMOVED***
			if err != nil ***REMOVED***
				fmt.Println("failed:", err)
				return
			***REMOVED***
			e.EncodeToken(t)
		***REMOVED***
		e.Flush()
		in = out
	***REMOVED***
	if got := in.String(); got != s ***REMOVED***
		t.Errorf("have: %q\nwant: %q\n", got, s)
	***REMOVED***
***REMOVED***

func TestMarshalNS(t *testing.T) ***REMOVED***
	dst := Tables***REMOVED***"hello", "world"***REMOVED***
	data, err := Marshal(&dst)
	if err != nil ***REMOVED***
		t.Fatalf("Marshal: %v", err)
	***REMOVED***
	want := `<Tables><table xmlns="http://www.w3.org/TR/html4/">hello</table><table xmlns="http://www.w3schools.com/furniture">world</table></Tables>`
	str := string(data)
	if str != want ***REMOVED***
		t.Errorf("have: %q\nwant: %q\n", str, want)
	***REMOVED***
***REMOVED***

type TableAttrs struct ***REMOVED***
	TAttr TAttr
***REMOVED***

type TAttr struct ***REMOVED***
	HTable string `xml:"http://www.w3.org/TR/html4/ table,attr"`
	FTable string `xml:"http://www.w3schools.com/furniture table,attr"`
	Lang   string `xml:"http://www.w3.org/XML/1998/namespace lang,attr,omitempty"`
	Other1 string `xml:"http://golang.org/xml/ other,attr,omitempty"`
	Other2 string `xml:"http://golang.org/xmlfoo/ other,attr,omitempty"`
	Other3 string `xml:"http://golang.org/json/ other,attr,omitempty"`
	Other4 string `xml:"http://golang.org/2/json/ other,attr,omitempty"`
***REMOVED***

var tableAttrs = []struct ***REMOVED***
	xml string
	tab TableAttrs
	ns  string
***REMOVED******REMOVED***
	***REMOVED***
		xml: `<TableAttrs xmlns:f="http://www.w3schools.com/furniture" xmlns:h="http://www.w3.org/TR/html4/"><TAttr ` +
			`h:table="hello" f:table="world" ` +
			`/></TableAttrs>`,
		tab: TableAttrs***REMOVED***TAttr***REMOVED***HTable: "hello", FTable: "world"***REMOVED******REMOVED***,
	***REMOVED***,
	***REMOVED***
		xml: `<TableAttrs><TAttr xmlns:f="http://www.w3schools.com/furniture" xmlns:h="http://www.w3.org/TR/html4/" ` +
			`h:table="hello" f:table="world" ` +
			`/></TableAttrs>`,
		tab: TableAttrs***REMOVED***TAttr***REMOVED***HTable: "hello", FTable: "world"***REMOVED******REMOVED***,
	***REMOVED***,
	***REMOVED***
		xml: `<TableAttrs><TAttr ` +
			`h:table="hello" f:table="world" xmlns:f="http://www.w3schools.com/furniture" xmlns:h="http://www.w3.org/TR/html4/" ` +
			`/></TableAttrs>`,
		tab: TableAttrs***REMOVED***TAttr***REMOVED***HTable: "hello", FTable: "world"***REMOVED******REMOVED***,
	***REMOVED***,
	***REMOVED***
		// Default space does not apply to attribute names.
		xml: `<TableAttrs xmlns="http://www.w3schools.com/furniture" xmlns:h="http://www.w3.org/TR/html4/"><TAttr ` +
			`h:table="hello" table="world" ` +
			`/></TableAttrs>`,
		tab: TableAttrs***REMOVED***TAttr***REMOVED***HTable: "hello", FTable: ""***REMOVED******REMOVED***,
	***REMOVED***,
	***REMOVED***
		// Default space does not apply to attribute names.
		xml: `<TableAttrs xmlns:f="http://www.w3schools.com/furniture"><TAttr xmlns="http://www.w3.org/TR/html4/" ` +
			`table="hello" f:table="world" ` +
			`/></TableAttrs>`,
		tab: TableAttrs***REMOVED***TAttr***REMOVED***HTable: "", FTable: "world"***REMOVED******REMOVED***,
	***REMOVED***,
	***REMOVED***
		xml: `<TableAttrs><TAttr ` +
			`table="bogus" ` +
			`/></TableAttrs>`,
		tab: TableAttrs***REMOVED******REMOVED***,
	***REMOVED***,
	***REMOVED***
		// Default space does not apply to attribute names.
		xml: `<TableAttrs xmlns:h="http://www.w3.org/TR/html4/"><TAttr ` +
			`h:table="hello" table="world" ` +
			`/></TableAttrs>`,
		tab: TableAttrs***REMOVED***TAttr***REMOVED***HTable: "hello", FTable: ""***REMOVED******REMOVED***,
		ns:  "http://www.w3schools.com/furniture",
	***REMOVED***,
	***REMOVED***
		// Default space does not apply to attribute names.
		xml: `<TableAttrs xmlns:f="http://www.w3schools.com/furniture"><TAttr ` +
			`table="hello" f:table="world" ` +
			`/></TableAttrs>`,
		tab: TableAttrs***REMOVED***TAttr***REMOVED***HTable: "", FTable: "world"***REMOVED******REMOVED***,
		ns:  "http://www.w3.org/TR/html4/",
	***REMOVED***,
	***REMOVED***
		xml: `<TableAttrs><TAttr ` +
			`table="bogus" ` +
			`/></TableAttrs>`,
		tab: TableAttrs***REMOVED******REMOVED***,
		ns:  "something else entirely",
	***REMOVED***,
***REMOVED***

func TestUnmarshalNSAttr(t *testing.T) ***REMOVED***
	for i, tt := range tableAttrs ***REMOVED***
		var dst TableAttrs
		var err error
		if tt.ns != "" ***REMOVED***
			d := NewDecoder(strings.NewReader(tt.xml))
			d.DefaultSpace = tt.ns
			err = d.Decode(&dst)
		***REMOVED*** else ***REMOVED***
			err = Unmarshal([]byte(tt.xml), &dst)
		***REMOVED***
		if err != nil ***REMOVED***
			t.Errorf("#%d: Unmarshal: %v", i, err)
			continue
		***REMOVED***
		want := tt.tab
		if dst != want ***REMOVED***
			t.Errorf("#%d: dst=%+v, want %+v", i, dst, want)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestMarshalNSAttr(t *testing.T) ***REMOVED***
	src := TableAttrs***REMOVED***TAttr***REMOVED***"hello", "world", "en_US", "other1", "other2", "other3", "other4"***REMOVED******REMOVED***
	data, err := Marshal(&src)
	if err != nil ***REMOVED***
		t.Fatalf("Marshal: %v", err)
	***REMOVED***
	want := `<TableAttrs><TAttr xmlns:json_1="http://golang.org/2/json/" xmlns:json="http://golang.org/json/" xmlns:_xmlfoo="http://golang.org/xmlfoo/" xmlns:_xml="http://golang.org/xml/" xmlns:furniture="http://www.w3schools.com/furniture" xmlns:html4="http://www.w3.org/TR/html4/" html4:table="hello" furniture:table="world" xml:lang="en_US" _xml:other="other1" _xmlfoo:other="other2" json:other="other3" json_1:other="other4"></TAttr></TableAttrs>`
	str := string(data)
	if str != want ***REMOVED***
		t.Errorf("Marshal:\nhave: %#q\nwant: %#q\n", str, want)
	***REMOVED***

	var dst TableAttrs
	if err := Unmarshal(data, &dst); err != nil ***REMOVED***
		t.Errorf("Unmarshal: %v", err)
	***REMOVED***

	if dst != src ***REMOVED***
		t.Errorf("Unmarshal = %q, want %q", dst, src)
	***REMOVED***
***REMOVED***

type MyCharData struct ***REMOVED***
	body string
***REMOVED***

func (m *MyCharData) UnmarshalXML(d *Decoder, start StartElement) error ***REMOVED***
	for ***REMOVED***
		t, err := d.Token()
		if err == io.EOF ***REMOVED*** // found end of element
			break
		***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if char, ok := t.(CharData); ok ***REMOVED***
			m.body += string(char)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

var _ Unmarshaler = (*MyCharData)(nil)

func (m *MyCharData) UnmarshalXMLAttr(attr Attr) error ***REMOVED***
	panic("must not call")
***REMOVED***

type MyAttr struct ***REMOVED***
	attr string
***REMOVED***

func (m *MyAttr) UnmarshalXMLAttr(attr Attr) error ***REMOVED***
	m.attr = attr.Value
	return nil
***REMOVED***

var _ UnmarshalerAttr = (*MyAttr)(nil)

type MyStruct struct ***REMOVED***
	Data *MyCharData
	Attr *MyAttr `xml:",attr"`

	Data2 MyCharData
	Attr2 MyAttr `xml:",attr"`
***REMOVED***

func TestUnmarshaler(t *testing.T) ***REMOVED***
	xml := `<?xml version="1.0" encoding="utf-8"?>
		<MyStruct Attr="attr1" Attr2="attr2">
		<Data>hello <!-- comment -->world</Data>
		<Data2>howdy <!-- comment -->world</Data2>
		</MyStruct>
	`

	var m MyStruct
	if err := Unmarshal([]byte(xml), &m); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if m.Data == nil || m.Attr == nil || m.Data.body != "hello world" || m.Attr.attr != "attr1" || m.Data2.body != "howdy world" || m.Attr2.attr != "attr2" ***REMOVED***
		t.Errorf("m=%#+v\n", m)
	***REMOVED***
***REMOVED***

type Pea struct ***REMOVED***
	Cotelydon string
***REMOVED***

type Pod struct ***REMOVED***
	Pea interface***REMOVED******REMOVED*** `xml:"Pea"`
***REMOVED***

// https://golang.org/issue/6836
func TestUnmarshalIntoInterface(t *testing.T) ***REMOVED***
	pod := new(Pod)
	pod.Pea = new(Pea)
	xml := `<Pod><Pea><Cotelydon>Green stuff</Cotelydon></Pea></Pod>`
	err := Unmarshal([]byte(xml), pod)
	if err != nil ***REMOVED***
		t.Fatalf("failed to unmarshal %q: %v", xml, err)
	***REMOVED***
	pea, ok := pod.Pea.(*Pea)
	if !ok ***REMOVED***
		t.Fatalf("unmarshalled into wrong type: have %T want *Pea", pod.Pea)
	***REMOVED***
	have, want := pea.Cotelydon, "Green stuff"
	if have != want ***REMOVED***
		t.Errorf("failed to unmarshal into interface, have %q want %q", have, want)
	***REMOVED***
***REMOVED***
