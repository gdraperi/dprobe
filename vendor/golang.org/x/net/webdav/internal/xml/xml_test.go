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
	"unicode/utf8"
)

const testInput = `
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN"
  "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<body xmlns:foo="ns1" xmlns="ns2" xmlns:tag="ns3" ` +
	"\r\n\t" + `  >
  <hello lang="en">World &lt;&gt;&apos;&quot; &#x767d;&#40300;翔</hello>
  <query>&何; &is-it;</query>
  <goodbye />
  <outer foo:attr="value" xmlns:tag="ns4">
    <inner/>
  </outer>
  <tag:name>
    <![CDATA[Some text here.]]>
  </tag:name>
</body><!-- missing final newline -->`

var testEntity = map[string]string***REMOVED***"何": "What", "is-it": "is it?"***REMOVED***

var rawTokens = []Token***REMOVED***
	CharData("\n"),
	ProcInst***REMOVED***"xml", []byte(`version="1.0" encoding="UTF-8"`)***REMOVED***,
	CharData("\n"),
	Directive(`DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN"
  "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd"`),
	CharData("\n"),
	StartElement***REMOVED***Name***REMOVED***"", "body"***REMOVED***, []Attr***REMOVED******REMOVED***Name***REMOVED***"xmlns", "foo"***REMOVED***, "ns1"***REMOVED***, ***REMOVED***Name***REMOVED***"", "xmlns"***REMOVED***, "ns2"***REMOVED***, ***REMOVED***Name***REMOVED***"xmlns", "tag"***REMOVED***, "ns3"***REMOVED******REMOVED******REMOVED***,
	CharData("\n  "),
	StartElement***REMOVED***Name***REMOVED***"", "hello"***REMOVED***, []Attr***REMOVED******REMOVED***Name***REMOVED***"", "lang"***REMOVED***, "en"***REMOVED******REMOVED******REMOVED***,
	CharData("World <>'\" 白鵬翔"),
	EndElement***REMOVED***Name***REMOVED***"", "hello"***REMOVED******REMOVED***,
	CharData("\n  "),
	StartElement***REMOVED***Name***REMOVED***"", "query"***REMOVED***, []Attr***REMOVED******REMOVED******REMOVED***,
	CharData("What is it?"),
	EndElement***REMOVED***Name***REMOVED***"", "query"***REMOVED******REMOVED***,
	CharData("\n  "),
	StartElement***REMOVED***Name***REMOVED***"", "goodbye"***REMOVED***, []Attr***REMOVED******REMOVED******REMOVED***,
	EndElement***REMOVED***Name***REMOVED***"", "goodbye"***REMOVED******REMOVED***,
	CharData("\n  "),
	StartElement***REMOVED***Name***REMOVED***"", "outer"***REMOVED***, []Attr***REMOVED******REMOVED***Name***REMOVED***"foo", "attr"***REMOVED***, "value"***REMOVED***, ***REMOVED***Name***REMOVED***"xmlns", "tag"***REMOVED***, "ns4"***REMOVED******REMOVED******REMOVED***,
	CharData("\n    "),
	StartElement***REMOVED***Name***REMOVED***"", "inner"***REMOVED***, []Attr***REMOVED******REMOVED******REMOVED***,
	EndElement***REMOVED***Name***REMOVED***"", "inner"***REMOVED******REMOVED***,
	CharData("\n  "),
	EndElement***REMOVED***Name***REMOVED***"", "outer"***REMOVED******REMOVED***,
	CharData("\n  "),
	StartElement***REMOVED***Name***REMOVED***"tag", "name"***REMOVED***, []Attr***REMOVED******REMOVED******REMOVED***,
	CharData("\n    "),
	CharData("Some text here."),
	CharData("\n  "),
	EndElement***REMOVED***Name***REMOVED***"tag", "name"***REMOVED******REMOVED***,
	CharData("\n"),
	EndElement***REMOVED***Name***REMOVED***"", "body"***REMOVED******REMOVED***,
	Comment(" missing final newline "),
***REMOVED***

var cookedTokens = []Token***REMOVED***
	CharData("\n"),
	ProcInst***REMOVED***"xml", []byte(`version="1.0" encoding="UTF-8"`)***REMOVED***,
	CharData("\n"),
	Directive(`DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN"
  "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd"`),
	CharData("\n"),
	StartElement***REMOVED***Name***REMOVED***"ns2", "body"***REMOVED***, []Attr***REMOVED******REMOVED***Name***REMOVED***"xmlns", "foo"***REMOVED***, "ns1"***REMOVED***, ***REMOVED***Name***REMOVED***"", "xmlns"***REMOVED***, "ns2"***REMOVED***, ***REMOVED***Name***REMOVED***"xmlns", "tag"***REMOVED***, "ns3"***REMOVED******REMOVED******REMOVED***,
	CharData("\n  "),
	StartElement***REMOVED***Name***REMOVED***"ns2", "hello"***REMOVED***, []Attr***REMOVED******REMOVED***Name***REMOVED***"", "lang"***REMOVED***, "en"***REMOVED******REMOVED******REMOVED***,
	CharData("World <>'\" 白鵬翔"),
	EndElement***REMOVED***Name***REMOVED***"ns2", "hello"***REMOVED******REMOVED***,
	CharData("\n  "),
	StartElement***REMOVED***Name***REMOVED***"ns2", "query"***REMOVED***, []Attr***REMOVED******REMOVED******REMOVED***,
	CharData("What is it?"),
	EndElement***REMOVED***Name***REMOVED***"ns2", "query"***REMOVED******REMOVED***,
	CharData("\n  "),
	StartElement***REMOVED***Name***REMOVED***"ns2", "goodbye"***REMOVED***, []Attr***REMOVED******REMOVED******REMOVED***,
	EndElement***REMOVED***Name***REMOVED***"ns2", "goodbye"***REMOVED******REMOVED***,
	CharData("\n  "),
	StartElement***REMOVED***Name***REMOVED***"ns2", "outer"***REMOVED***, []Attr***REMOVED******REMOVED***Name***REMOVED***"ns1", "attr"***REMOVED***, "value"***REMOVED***, ***REMOVED***Name***REMOVED***"xmlns", "tag"***REMOVED***, "ns4"***REMOVED******REMOVED******REMOVED***,
	CharData("\n    "),
	StartElement***REMOVED***Name***REMOVED***"ns2", "inner"***REMOVED***, []Attr***REMOVED******REMOVED******REMOVED***,
	EndElement***REMOVED***Name***REMOVED***"ns2", "inner"***REMOVED******REMOVED***,
	CharData("\n  "),
	EndElement***REMOVED***Name***REMOVED***"ns2", "outer"***REMOVED******REMOVED***,
	CharData("\n  "),
	StartElement***REMOVED***Name***REMOVED***"ns3", "name"***REMOVED***, []Attr***REMOVED******REMOVED******REMOVED***,
	CharData("\n    "),
	CharData("Some text here."),
	CharData("\n  "),
	EndElement***REMOVED***Name***REMOVED***"ns3", "name"***REMOVED******REMOVED***,
	CharData("\n"),
	EndElement***REMOVED***Name***REMOVED***"ns2", "body"***REMOVED******REMOVED***,
	Comment(" missing final newline "),
***REMOVED***

const testInputAltEncoding = `
<?xml version="1.0" encoding="x-testing-uppercase"?>
<TAG>VALUE</TAG>`

var rawTokensAltEncoding = []Token***REMOVED***
	CharData("\n"),
	ProcInst***REMOVED***"xml", []byte(`version="1.0" encoding="x-testing-uppercase"`)***REMOVED***,
	CharData("\n"),
	StartElement***REMOVED***Name***REMOVED***"", "tag"***REMOVED***, []Attr***REMOVED******REMOVED******REMOVED***,
	CharData("value"),
	EndElement***REMOVED***Name***REMOVED***"", "tag"***REMOVED******REMOVED***,
***REMOVED***

var xmlInput = []string***REMOVED***
	// unexpected EOF cases
	"<",
	"<t",
	"<t ",
	"<t/",
	"<!",
	"<!-",
	"<!--",
	"<!--c-",
	"<!--c--",
	"<!d",
	"<t></",
	"<t></t",
	"<?",
	"<?p",
	"<t a",
	"<t a=",
	"<t a='",
	"<t a=''",
	"<t/><![",
	"<t/><![C",
	"<t/><![CDATA[d",
	"<t/><![CDATA[d]",
	"<t/><![CDATA[d]]",

	// other Syntax errors
	"<>",
	"<t/a",
	"<0 />",
	"<?0 >",
	//	"<!0 >",	// let the Token() caller handle
	"</0>",
	"<t 0=''>",
	"<t a='&'>",
	"<t a='<'>",
	"<t>&nbspc;</t>",
	"<t a>",
	"<t a=>",
	"<t a=v>",
	//	"<![CDATA[d]]>",	// let the Token() caller handle
	"<t></e>",
	"<t></>",
	"<t></t!",
	"<t>cdata]]></t>",
***REMOVED***

func TestRawToken(t *testing.T) ***REMOVED***
	d := NewDecoder(strings.NewReader(testInput))
	d.Entity = testEntity
	testRawToken(t, d, testInput, rawTokens)
***REMOVED***

const nonStrictInput = `
<tag>non&entity</tag>
<tag>&unknown;entity</tag>
<tag>&#123</tag>
<tag>&#zzz;</tag>
<tag>&なまえ3;</tag>
<tag>&lt-gt;</tag>
<tag>&;</tag>
<tag>&0a;</tag>
`

var nonStringEntity = map[string]string***REMOVED***"": "oops!", "0a": "oops!"***REMOVED***

var nonStrictTokens = []Token***REMOVED***
	CharData("\n"),
	StartElement***REMOVED***Name***REMOVED***"", "tag"***REMOVED***, []Attr***REMOVED******REMOVED******REMOVED***,
	CharData("non&entity"),
	EndElement***REMOVED***Name***REMOVED***"", "tag"***REMOVED******REMOVED***,
	CharData("\n"),
	StartElement***REMOVED***Name***REMOVED***"", "tag"***REMOVED***, []Attr***REMOVED******REMOVED******REMOVED***,
	CharData("&unknown;entity"),
	EndElement***REMOVED***Name***REMOVED***"", "tag"***REMOVED******REMOVED***,
	CharData("\n"),
	StartElement***REMOVED***Name***REMOVED***"", "tag"***REMOVED***, []Attr***REMOVED******REMOVED******REMOVED***,
	CharData("&#123"),
	EndElement***REMOVED***Name***REMOVED***"", "tag"***REMOVED******REMOVED***,
	CharData("\n"),
	StartElement***REMOVED***Name***REMOVED***"", "tag"***REMOVED***, []Attr***REMOVED******REMOVED******REMOVED***,
	CharData("&#zzz;"),
	EndElement***REMOVED***Name***REMOVED***"", "tag"***REMOVED******REMOVED***,
	CharData("\n"),
	StartElement***REMOVED***Name***REMOVED***"", "tag"***REMOVED***, []Attr***REMOVED******REMOVED******REMOVED***,
	CharData("&なまえ3;"),
	EndElement***REMOVED***Name***REMOVED***"", "tag"***REMOVED******REMOVED***,
	CharData("\n"),
	StartElement***REMOVED***Name***REMOVED***"", "tag"***REMOVED***, []Attr***REMOVED******REMOVED******REMOVED***,
	CharData("&lt-gt;"),
	EndElement***REMOVED***Name***REMOVED***"", "tag"***REMOVED******REMOVED***,
	CharData("\n"),
	StartElement***REMOVED***Name***REMOVED***"", "tag"***REMOVED***, []Attr***REMOVED******REMOVED******REMOVED***,
	CharData("&;"),
	EndElement***REMOVED***Name***REMOVED***"", "tag"***REMOVED******REMOVED***,
	CharData("\n"),
	StartElement***REMOVED***Name***REMOVED***"", "tag"***REMOVED***, []Attr***REMOVED******REMOVED******REMOVED***,
	CharData("&0a;"),
	EndElement***REMOVED***Name***REMOVED***"", "tag"***REMOVED******REMOVED***,
	CharData("\n"),
***REMOVED***

func TestNonStrictRawToken(t *testing.T) ***REMOVED***
	d := NewDecoder(strings.NewReader(nonStrictInput))
	d.Strict = false
	testRawToken(t, d, nonStrictInput, nonStrictTokens)
***REMOVED***

type downCaser struct ***REMOVED***
	t *testing.T
	r io.ByteReader
***REMOVED***

func (d *downCaser) ReadByte() (c byte, err error) ***REMOVED***
	c, err = d.r.ReadByte()
	if c >= 'A' && c <= 'Z' ***REMOVED***
		c += 'a' - 'A'
	***REMOVED***
	return
***REMOVED***

func (d *downCaser) Read(p []byte) (int, error) ***REMOVED***
	d.t.Fatalf("unexpected Read call on downCaser reader")
	panic("unreachable")
***REMOVED***

func TestRawTokenAltEncoding(t *testing.T) ***REMOVED***
	d := NewDecoder(strings.NewReader(testInputAltEncoding))
	d.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) ***REMOVED***
		if charset != "x-testing-uppercase" ***REMOVED***
			t.Fatalf("unexpected charset %q", charset)
		***REMOVED***
		return &downCaser***REMOVED***t, input.(io.ByteReader)***REMOVED***, nil
	***REMOVED***
	testRawToken(t, d, testInputAltEncoding, rawTokensAltEncoding)
***REMOVED***

func TestRawTokenAltEncodingNoConverter(t *testing.T) ***REMOVED***
	d := NewDecoder(strings.NewReader(testInputAltEncoding))
	token, err := d.RawToken()
	if token == nil ***REMOVED***
		t.Fatalf("expected a token on first RawToken call")
	***REMOVED***
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	token, err = d.RawToken()
	if token != nil ***REMOVED***
		t.Errorf("expected a nil token; got %#v", token)
	***REMOVED***
	if err == nil ***REMOVED***
		t.Fatalf("expected an error on second RawToken call")
	***REMOVED***
	const encoding = "x-testing-uppercase"
	if !strings.Contains(err.Error(), encoding) ***REMOVED***
		t.Errorf("expected error to contain %q; got error: %v",
			encoding, err)
	***REMOVED***
***REMOVED***

func testRawToken(t *testing.T, d *Decoder, raw string, rawTokens []Token) ***REMOVED***
	lastEnd := int64(0)
	for i, want := range rawTokens ***REMOVED***
		start := d.InputOffset()
		have, err := d.RawToken()
		end := d.InputOffset()
		if err != nil ***REMOVED***
			t.Fatalf("token %d: unexpected error: %s", i, err)
		***REMOVED***
		if !reflect.DeepEqual(have, want) ***REMOVED***
			var shave, swant string
			if _, ok := have.(CharData); ok ***REMOVED***
				shave = fmt.Sprintf("CharData(%q)", have)
			***REMOVED*** else ***REMOVED***
				shave = fmt.Sprintf("%#v", have)
			***REMOVED***
			if _, ok := want.(CharData); ok ***REMOVED***
				swant = fmt.Sprintf("CharData(%q)", want)
			***REMOVED*** else ***REMOVED***
				swant = fmt.Sprintf("%#v", want)
			***REMOVED***
			t.Errorf("token %d = %s, want %s", i, shave, swant)
		***REMOVED***

		// Check that InputOffset returned actual token.
		switch ***REMOVED***
		case start < lastEnd:
			t.Errorf("token %d: position [%d,%d) for %T is before previous token", i, start, end, have)
		case start >= end:
			// Special case: EndElement can be synthesized.
			if start == end && end == lastEnd ***REMOVED***
				break
			***REMOVED***
			t.Errorf("token %d: position [%d,%d) for %T is empty", i, start, end, have)
		case end > int64(len(raw)):
			t.Errorf("token %d: position [%d,%d) for %T extends beyond input", i, start, end, have)
		default:
			text := raw[start:end]
			if strings.ContainsAny(text, "<>") && (!strings.HasPrefix(text, "<") || !strings.HasSuffix(text, ">")) ***REMOVED***
				t.Errorf("token %d: misaligned raw token %#q for %T", i, text, have)
			***REMOVED***
		***REMOVED***
		lastEnd = end
	***REMOVED***
***REMOVED***

// Ensure that directives (specifically !DOCTYPE) include the complete
// text of any nested directives, noting that < and > do not change
// nesting depth if they are in single or double quotes.

var nestedDirectivesInput = `
<!DOCTYPE [<!ENTITY rdf "http://www.w3.org/1999/02/22-rdf-syntax-ns#">]>
<!DOCTYPE [<!ENTITY xlt ">">]>
<!DOCTYPE [<!ENTITY xlt "<">]>
<!DOCTYPE [<!ENTITY xlt '>'>]>
<!DOCTYPE [<!ENTITY xlt '<'>]>
<!DOCTYPE [<!ENTITY xlt '">'>]>
<!DOCTYPE [<!ENTITY xlt "'<">]>
`

var nestedDirectivesTokens = []Token***REMOVED***
	CharData("\n"),
	Directive(`DOCTYPE [<!ENTITY rdf "http://www.w3.org/1999/02/22-rdf-syntax-ns#">]`),
	CharData("\n"),
	Directive(`DOCTYPE [<!ENTITY xlt ">">]`),
	CharData("\n"),
	Directive(`DOCTYPE [<!ENTITY xlt "<">]`),
	CharData("\n"),
	Directive(`DOCTYPE [<!ENTITY xlt '>'>]`),
	CharData("\n"),
	Directive(`DOCTYPE [<!ENTITY xlt '<'>]`),
	CharData("\n"),
	Directive(`DOCTYPE [<!ENTITY xlt '">'>]`),
	CharData("\n"),
	Directive(`DOCTYPE [<!ENTITY xlt "'<">]`),
	CharData("\n"),
***REMOVED***

func TestNestedDirectives(t *testing.T) ***REMOVED***
	d := NewDecoder(strings.NewReader(nestedDirectivesInput))

	for i, want := range nestedDirectivesTokens ***REMOVED***
		have, err := d.Token()
		if err != nil ***REMOVED***
			t.Fatalf("token %d: unexpected error: %s", i, err)
		***REMOVED***
		if !reflect.DeepEqual(have, want) ***REMOVED***
			t.Errorf("token %d = %#v want %#v", i, have, want)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestToken(t *testing.T) ***REMOVED***
	d := NewDecoder(strings.NewReader(testInput))
	d.Entity = testEntity

	for i, want := range cookedTokens ***REMOVED***
		have, err := d.Token()
		if err != nil ***REMOVED***
			t.Fatalf("token %d: unexpected error: %s", i, err)
		***REMOVED***
		if !reflect.DeepEqual(have, want) ***REMOVED***
			t.Errorf("token %d = %#v want %#v", i, have, want)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestSyntax(t *testing.T) ***REMOVED***
	for i := range xmlInput ***REMOVED***
		d := NewDecoder(strings.NewReader(xmlInput[i]))
		var err error
		for _, err = d.Token(); err == nil; _, err = d.Token() ***REMOVED***
		***REMOVED***
		if _, ok := err.(*SyntaxError); !ok ***REMOVED***
			t.Fatalf(`xmlInput "%s": expected SyntaxError not received`, xmlInput[i])
		***REMOVED***
	***REMOVED***
***REMOVED***

type allScalars struct ***REMOVED***
	True1     bool
	True2     bool
	False1    bool
	False2    bool
	Int       int
	Int8      int8
	Int16     int16
	Int32     int32
	Int64     int64
	Uint      int
	Uint8     uint8
	Uint16    uint16
	Uint32    uint32
	Uint64    uint64
	Uintptr   uintptr
	Float32   float32
	Float64   float64
	String    string
	PtrString *string
***REMOVED***

var all = allScalars***REMOVED***
	True1:     true,
	True2:     true,
	False1:    false,
	False2:    false,
	Int:       1,
	Int8:      -2,
	Int16:     3,
	Int32:     -4,
	Int64:     5,
	Uint:      6,
	Uint8:     7,
	Uint16:    8,
	Uint32:    9,
	Uint64:    10,
	Uintptr:   11,
	Float32:   13.0,
	Float64:   14.0,
	String:    "15",
	PtrString: &sixteen,
***REMOVED***

var sixteen = "16"

const testScalarsInput = `<allscalars>
	<True1>true</True1>
	<True2>1</True2>
	<False1>false</False1>
	<False2>0</False2>
	<Int>1</Int>
	<Int8>-2</Int8>
	<Int16>3</Int16>
	<Int32>-4</Int32>
	<Int64>5</Int64>
	<Uint>6</Uint>
	<Uint8>7</Uint8>
	<Uint16>8</Uint16>
	<Uint32>9</Uint32>
	<Uint64>10</Uint64>
	<Uintptr>11</Uintptr>
	<Float>12.0</Float>
	<Float32>13.0</Float32>
	<Float64>14.0</Float64>
	<String>15</String>
	<PtrString>16</PtrString>
</allscalars>`

func TestAllScalars(t *testing.T) ***REMOVED***
	var a allScalars
	err := Unmarshal([]byte(testScalarsInput), &a)

	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if !reflect.DeepEqual(a, all) ***REMOVED***
		t.Errorf("have %+v want %+v", a, all)
	***REMOVED***
***REMOVED***

type item struct ***REMOVED***
	Field_a string
***REMOVED***

func TestIssue569(t *testing.T) ***REMOVED***
	data := `<item><Field_a>abcd</Field_a></item>`
	var i item
	err := Unmarshal([]byte(data), &i)

	if err != nil || i.Field_a != "abcd" ***REMOVED***
		t.Fatal("Expecting abcd")
	***REMOVED***
***REMOVED***

func TestUnquotedAttrs(t *testing.T) ***REMOVED***
	data := "<tag attr=azAZ09:-_\t>"
	d := NewDecoder(strings.NewReader(data))
	d.Strict = false
	token, err := d.Token()
	if _, ok := err.(*SyntaxError); ok ***REMOVED***
		t.Errorf("Unexpected error: %v", err)
	***REMOVED***
	if token.(StartElement).Name.Local != "tag" ***REMOVED***
		t.Errorf("Unexpected tag name: %v", token.(StartElement).Name.Local)
	***REMOVED***
	attr := token.(StartElement).Attr[0]
	if attr.Value != "azAZ09:-_" ***REMOVED***
		t.Errorf("Unexpected attribute value: %v", attr.Value)
	***REMOVED***
	if attr.Name.Local != "attr" ***REMOVED***
		t.Errorf("Unexpected attribute name: %v", attr.Name.Local)
	***REMOVED***
***REMOVED***

func TestValuelessAttrs(t *testing.T) ***REMOVED***
	tests := [][3]string***REMOVED***
		***REMOVED***"<p nowrap>", "p", "nowrap"***REMOVED***,
		***REMOVED***"<p nowrap >", "p", "nowrap"***REMOVED***,
		***REMOVED***"<input checked/>", "input", "checked"***REMOVED***,
		***REMOVED***"<input checked />", "input", "checked"***REMOVED***,
	***REMOVED***
	for _, test := range tests ***REMOVED***
		d := NewDecoder(strings.NewReader(test[0]))
		d.Strict = false
		token, err := d.Token()
		if _, ok := err.(*SyntaxError); ok ***REMOVED***
			t.Errorf("Unexpected error: %v", err)
		***REMOVED***
		if token.(StartElement).Name.Local != test[1] ***REMOVED***
			t.Errorf("Unexpected tag name: %v", token.(StartElement).Name.Local)
		***REMOVED***
		attr := token.(StartElement).Attr[0]
		if attr.Value != test[2] ***REMOVED***
			t.Errorf("Unexpected attribute value: %v", attr.Value)
		***REMOVED***
		if attr.Name.Local != test[2] ***REMOVED***
			t.Errorf("Unexpected attribute name: %v", attr.Name.Local)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestCopyTokenCharData(t *testing.T) ***REMOVED***
	data := []byte("same data")
	var tok1 Token = CharData(data)
	tok2 := CopyToken(tok1)
	if !reflect.DeepEqual(tok1, tok2) ***REMOVED***
		t.Error("CopyToken(CharData) != CharData")
	***REMOVED***
	data[1] = 'o'
	if reflect.DeepEqual(tok1, tok2) ***REMOVED***
		t.Error("CopyToken(CharData) uses same buffer.")
	***REMOVED***
***REMOVED***

func TestCopyTokenStartElement(t *testing.T) ***REMOVED***
	elt := StartElement***REMOVED***Name***REMOVED***"", "hello"***REMOVED***, []Attr***REMOVED******REMOVED***Name***REMOVED***"", "lang"***REMOVED***, "en"***REMOVED******REMOVED******REMOVED***
	var tok1 Token = elt
	tok2 := CopyToken(tok1)
	if tok1.(StartElement).Attr[0].Value != "en" ***REMOVED***
		t.Error("CopyToken overwrote Attr[0]")
	***REMOVED***
	if !reflect.DeepEqual(tok1, tok2) ***REMOVED***
		t.Error("CopyToken(StartElement) != StartElement")
	***REMOVED***
	tok1.(StartElement).Attr[0] = Attr***REMOVED***Name***REMOVED***"", "lang"***REMOVED***, "de"***REMOVED***
	if reflect.DeepEqual(tok1, tok2) ***REMOVED***
		t.Error("CopyToken(CharData) uses same buffer.")
	***REMOVED***
***REMOVED***

func TestSyntaxErrorLineNum(t *testing.T) ***REMOVED***
	testInput := "<P>Foo<P>\n\n<P>Bar</>\n"
	d := NewDecoder(strings.NewReader(testInput))
	var err error
	for _, err = d.Token(); err == nil; _, err = d.Token() ***REMOVED***
	***REMOVED***
	synerr, ok := err.(*SyntaxError)
	if !ok ***REMOVED***
		t.Error("Expected SyntaxError.")
	***REMOVED***
	if synerr.Line != 3 ***REMOVED***
		t.Error("SyntaxError didn't have correct line number.")
	***REMOVED***
***REMOVED***

func TestTrailingRawToken(t *testing.T) ***REMOVED***
	input := `<FOO></FOO>  `
	d := NewDecoder(strings.NewReader(input))
	var err error
	for _, err = d.RawToken(); err == nil; _, err = d.RawToken() ***REMOVED***
	***REMOVED***
	if err != io.EOF ***REMOVED***
		t.Fatalf("d.RawToken() = _, %v, want _, io.EOF", err)
	***REMOVED***
***REMOVED***

func TestTrailingToken(t *testing.T) ***REMOVED***
	input := `<FOO></FOO>  `
	d := NewDecoder(strings.NewReader(input))
	var err error
	for _, err = d.Token(); err == nil; _, err = d.Token() ***REMOVED***
	***REMOVED***
	if err != io.EOF ***REMOVED***
		t.Fatalf("d.Token() = _, %v, want _, io.EOF", err)
	***REMOVED***
***REMOVED***

func TestEntityInsideCDATA(t *testing.T) ***REMOVED***
	input := `<test><![CDATA[ &val=foo ]]></test>`
	d := NewDecoder(strings.NewReader(input))
	var err error
	for _, err = d.Token(); err == nil; _, err = d.Token() ***REMOVED***
	***REMOVED***
	if err != io.EOF ***REMOVED***
		t.Fatalf("d.Token() = _, %v, want _, io.EOF", err)
	***REMOVED***
***REMOVED***

var characterTests = []struct ***REMOVED***
	in  string
	err string
***REMOVED******REMOVED***
	***REMOVED***"\x12<doc/>", "illegal character code U+0012"***REMOVED***,
	***REMOVED***"<?xml version=\"1.0\"?>\x0b<doc/>", "illegal character code U+000B"***REMOVED***,
	***REMOVED***"\xef\xbf\xbe<doc/>", "illegal character code U+FFFE"***REMOVED***,
	***REMOVED***"<?xml version=\"1.0\"?><doc>\r\n<hiya/>\x07<toots/></doc>", "illegal character code U+0007"***REMOVED***,
	***REMOVED***"<?xml version=\"1.0\"?><doc \x12='value'>what's up</doc>", "expected attribute name in element"***REMOVED***,
	***REMOVED***"<doc>&abc\x01;</doc>", "invalid character entity &abc (no semicolon)"***REMOVED***,
	***REMOVED***"<doc>&\x01;</doc>", "invalid character entity & (no semicolon)"***REMOVED***,
	***REMOVED***"<doc>&\xef\xbf\xbe;</doc>", "invalid character entity &\uFFFE;"***REMOVED***,
	***REMOVED***"<doc>&hello;</doc>", "invalid character entity &hello;"***REMOVED***,
***REMOVED***

func TestDisallowedCharacters(t *testing.T) ***REMOVED***

	for i, tt := range characterTests ***REMOVED***
		d := NewDecoder(strings.NewReader(tt.in))
		var err error

		for err == nil ***REMOVED***
			_, err = d.Token()
		***REMOVED***
		synerr, ok := err.(*SyntaxError)
		if !ok ***REMOVED***
			t.Fatalf("input %d d.Token() = _, %v, want _, *SyntaxError", i, err)
		***REMOVED***
		if synerr.Msg != tt.err ***REMOVED***
			t.Fatalf("input %d synerr.Msg wrong: want %q, got %q", i, tt.err, synerr.Msg)
		***REMOVED***
	***REMOVED***
***REMOVED***

type procInstEncodingTest struct ***REMOVED***
	expect, got string
***REMOVED***

var procInstTests = []struct ***REMOVED***
	input  string
	expect [2]string
***REMOVED******REMOVED***
	***REMOVED***`version="1.0" encoding="utf-8"`, [2]string***REMOVED***"1.0", "utf-8"***REMOVED******REMOVED***,
	***REMOVED***`version="1.0" encoding='utf-8'`, [2]string***REMOVED***"1.0", "utf-8"***REMOVED******REMOVED***,
	***REMOVED***`version="1.0" encoding='utf-8' `, [2]string***REMOVED***"1.0", "utf-8"***REMOVED******REMOVED***,
	***REMOVED***`version="1.0" encoding=utf-8`, [2]string***REMOVED***"1.0", ""***REMOVED******REMOVED***,
	***REMOVED***`encoding="FOO" `, [2]string***REMOVED***"", "FOO"***REMOVED******REMOVED***,
***REMOVED***

func TestProcInstEncoding(t *testing.T) ***REMOVED***
	for _, test := range procInstTests ***REMOVED***
		if got := procInst("version", test.input); got != test.expect[0] ***REMOVED***
			t.Errorf("procInst(version, %q) = %q; want %q", test.input, got, test.expect[0])
		***REMOVED***
		if got := procInst("encoding", test.input); got != test.expect[1] ***REMOVED***
			t.Errorf("procInst(encoding, %q) = %q; want %q", test.input, got, test.expect[1])
		***REMOVED***
	***REMOVED***
***REMOVED***

// Ensure that directives with comments include the complete
// text of any nested directives.

var directivesWithCommentsInput = `
<!DOCTYPE [<!-- a comment --><!ENTITY rdf "http://www.w3.org/1999/02/22-rdf-syntax-ns#">]>
<!DOCTYPE [<!ENTITY go "Golang"><!-- a comment-->]>
<!DOCTYPE <!-> <!> <!----> <!-->--> <!--->--> [<!ENTITY go "Golang"><!-- a comment-->]>
`

var directivesWithCommentsTokens = []Token***REMOVED***
	CharData("\n"),
	Directive(`DOCTYPE [<!ENTITY rdf "http://www.w3.org/1999/02/22-rdf-syntax-ns#">]`),
	CharData("\n"),
	Directive(`DOCTYPE [<!ENTITY go "Golang">]`),
	CharData("\n"),
	Directive(`DOCTYPE <!-> <!>    [<!ENTITY go "Golang">]`),
	CharData("\n"),
***REMOVED***

func TestDirectivesWithComments(t *testing.T) ***REMOVED***
	d := NewDecoder(strings.NewReader(directivesWithCommentsInput))

	for i, want := range directivesWithCommentsTokens ***REMOVED***
		have, err := d.Token()
		if err != nil ***REMOVED***
			t.Fatalf("token %d: unexpected error: %s", i, err)
		***REMOVED***
		if !reflect.DeepEqual(have, want) ***REMOVED***
			t.Errorf("token %d = %#v want %#v", i, have, want)
		***REMOVED***
	***REMOVED***
***REMOVED***

// Writer whose Write method always returns an error.
type errWriter struct***REMOVED******REMOVED***

func (errWriter) Write(p []byte) (n int, err error) ***REMOVED*** return 0, fmt.Errorf("unwritable") ***REMOVED***

func TestEscapeTextIOErrors(t *testing.T) ***REMOVED***
	expectErr := "unwritable"
	err := EscapeText(errWriter***REMOVED******REMOVED***, []byte***REMOVED***'A'***REMOVED***)

	if err == nil || err.Error() != expectErr ***REMOVED***
		t.Errorf("have %v, want %v", err, expectErr)
	***REMOVED***
***REMOVED***

func TestEscapeTextInvalidChar(t *testing.T) ***REMOVED***
	input := []byte("A \x00 terminated string.")
	expected := "A \uFFFD terminated string."

	buff := new(bytes.Buffer)
	if err := EscapeText(buff, input); err != nil ***REMOVED***
		t.Fatalf("have %v, want nil", err)
	***REMOVED***
	text := buff.String()

	if text != expected ***REMOVED***
		t.Errorf("have %v, want %v", text, expected)
	***REMOVED***
***REMOVED***

func TestIssue5880(t *testing.T) ***REMOVED***
	type T []byte
	data, err := Marshal(T***REMOVED***192, 168, 0, 1***REMOVED***)
	if err != nil ***REMOVED***
		t.Errorf("Marshal error: %v", err)
	***REMOVED***
	if !utf8.Valid(data) ***REMOVED***
		t.Errorf("Marshal generated invalid UTF-8: %x", data)
	***REMOVED***
***REMOVED***
