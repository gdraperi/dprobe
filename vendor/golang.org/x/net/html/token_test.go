// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package html

import (
	"bytes"
	"io"
	"io/ioutil"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

type tokenTest struct ***REMOVED***
	// A short description of the test case.
	desc string
	// The HTML to parse.
	html string
	// The string representations of the expected tokens, joined by '$'.
	golden string
***REMOVED***

var tokenTests = []tokenTest***REMOVED***
	***REMOVED***
		"empty",
		"",
		"",
	***REMOVED***,
	// A single text node. The tokenizer should not break text nodes on whitespace,
	// nor should it normalize whitespace within a text node.
	***REMOVED***
		"text",
		"foo  bar",
		"foo  bar",
	***REMOVED***,
	// An entity.
	***REMOVED***
		"entity",
		"one &lt; two",
		"one &lt; two",
	***REMOVED***,
	// A start, self-closing and end tag. The tokenizer does not care if the start
	// and end tokens don't match; that is the job of the parser.
	***REMOVED***
		"tags",
		"<a>b<c/>d</e>",
		"<a>$b$<c/>$d$</e>",
	***REMOVED***,
	// Angle brackets that aren't a tag.
	***REMOVED***
		"not a tag #0",
		"<",
		"&lt;",
	***REMOVED***,
	***REMOVED***
		"not a tag #1",
		"</",
		"&lt;/",
	***REMOVED***,
	***REMOVED***
		"not a tag #2",
		"</>",
		"<!---->",
	***REMOVED***,
	***REMOVED***
		"not a tag #3",
		"a</>b",
		"a$<!---->$b",
	***REMOVED***,
	***REMOVED***
		"not a tag #4",
		"</ >",
		"<!-- -->",
	***REMOVED***,
	***REMOVED***
		"not a tag #5",
		"</.",
		"<!--.-->",
	***REMOVED***,
	***REMOVED***
		"not a tag #6",
		"</.>",
		"<!--.-->",
	***REMOVED***,
	***REMOVED***
		"not a tag #7",
		"a < b",
		"a &lt; b",
	***REMOVED***,
	***REMOVED***
		"not a tag #8",
		"<.>",
		"&lt;.&gt;",
	***REMOVED***,
	***REMOVED***
		"not a tag #9",
		"a<<<b>>>c",
		"a&lt;&lt;$<b>$&gt;&gt;c",
	***REMOVED***,
	***REMOVED***
		"not a tag #10",
		"if x<0 and y < 0 then x*y>0",
		"if x&lt;0 and y &lt; 0 then x*y&gt;0",
	***REMOVED***,
	***REMOVED***
		"not a tag #11",
		"<<p>",
		"&lt;$<p>",
	***REMOVED***,
	// EOF in a tag name.
	***REMOVED***
		"tag name eof #0",
		"<a",
		"",
	***REMOVED***,
	***REMOVED***
		"tag name eof #1",
		"<a ",
		"",
	***REMOVED***,
	***REMOVED***
		"tag name eof #2",
		"a<b",
		"a",
	***REMOVED***,
	***REMOVED***
		"tag name eof #3",
		"<a><b",
		"<a>",
	***REMOVED***,
	***REMOVED***
		"tag name eof #4",
		`<a x`,
		``,
	***REMOVED***,
	// Some malformed tags that are missing a '>'.
	***REMOVED***
		"malformed tag #0",
		`<p</p>`,
		`<p< p="">`,
	***REMOVED***,
	***REMOVED***
		"malformed tag #1",
		`<p </p>`,
		`<p <="" p="">`,
	***REMOVED***,
	***REMOVED***
		"malformed tag #2",
		`<p id`,
		``,
	***REMOVED***,
	***REMOVED***
		"malformed tag #3",
		`<p id=`,
		``,
	***REMOVED***,
	***REMOVED***
		"malformed tag #4",
		`<p id=>`,
		`<p id="">`,
	***REMOVED***,
	***REMOVED***
		"malformed tag #5",
		`<p id=0`,
		``,
	***REMOVED***,
	***REMOVED***
		"malformed tag #6",
		`<p id=0</p>`,
		`<p id="0&lt;/p">`,
	***REMOVED***,
	***REMOVED***
		"malformed tag #7",
		`<p id="0</p>`,
		``,
	***REMOVED***,
	***REMOVED***
		"malformed tag #8",
		`<p id="0"</p>`,
		`<p id="0" <="" p="">`,
	***REMOVED***,
	***REMOVED***
		"malformed tag #9",
		`<p></p id`,
		`<p>`,
	***REMOVED***,
	// Raw text and RCDATA.
	***REMOVED***
		"basic raw text",
		"<script><a></b></script>",
		"<script>$&lt;a&gt;&lt;/b&gt;$</script>",
	***REMOVED***,
	***REMOVED***
		"unfinished script end tag",
		"<SCRIPT>a</SCR",
		"<script>$a&lt;/SCR",
	***REMOVED***,
	***REMOVED***
		"broken script end tag",
		"<SCRIPT>a</SCR ipt>",
		"<script>$a&lt;/SCR ipt&gt;",
	***REMOVED***,
	***REMOVED***
		"EOF in script end tag",
		"<SCRIPT>a</SCRipt",
		"<script>$a&lt;/SCRipt",
	***REMOVED***,
	***REMOVED***
		"scriptx end tag",
		"<SCRIPT>a</SCRiptx",
		"<script>$a&lt;/SCRiptx",
	***REMOVED***,
	***REMOVED***
		"' ' completes script end tag",
		"<SCRIPT>a</SCRipt ",
		"<script>$a",
	***REMOVED***,
	***REMOVED***
		"'>' completes script end tag",
		"<SCRIPT>a</SCRipt>",
		"<script>$a$</script>",
	***REMOVED***,
	***REMOVED***
		"self-closing script end tag",
		"<SCRIPT>a</SCRipt/>",
		"<script>$a$</script>",
	***REMOVED***,
	***REMOVED***
		"nested script tag",
		"<SCRIPT>a</SCRipt<script>",
		"<script>$a&lt;/SCRipt&lt;script&gt;",
	***REMOVED***,
	***REMOVED***
		"script end tag after unfinished",
		"<SCRIPT>a</SCRipt</script>",
		"<script>$a&lt;/SCRipt$</script>",
	***REMOVED***,
	***REMOVED***
		"script/style mismatched tags",
		"<script>a</style>",
		"<script>$a&lt;/style&gt;",
	***REMOVED***,
	***REMOVED***
		"style element with entity",
		"<style>&apos;",
		"<style>$&amp;apos;",
	***REMOVED***,
	***REMOVED***
		"textarea with tag",
		"<textarea><div></textarea>",
		"<textarea>$&lt;div&gt;$</textarea>",
	***REMOVED***,
	***REMOVED***
		"title with tag and entity",
		"<title><b>K&amp;R C</b></title>",
		"<title>$&lt;b&gt;K&amp;R C&lt;/b&gt;$</title>",
	***REMOVED***,
	// DOCTYPE tests.
	***REMOVED***
		"Proper DOCTYPE",
		"<!DOCTYPE html>",
		"<!DOCTYPE html>",
	***REMOVED***,
	***REMOVED***
		"DOCTYPE with no space",
		"<!doctypehtml>",
		"<!DOCTYPE html>",
	***REMOVED***,
	***REMOVED***
		"DOCTYPE with two spaces",
		"<!doctype  html>",
		"<!DOCTYPE html>",
	***REMOVED***,
	***REMOVED***
		"looks like DOCTYPE but isn't",
		"<!DOCUMENT html>",
		"<!--DOCUMENT html-->",
	***REMOVED***,
	***REMOVED***
		"DOCTYPE at EOF",
		"<!DOCtype",
		"<!DOCTYPE >",
	***REMOVED***,
	// XML processing instructions.
	***REMOVED***
		"XML processing instruction",
		"<?xml?>",
		"<!--?xml?-->",
	***REMOVED***,
	// Comments.
	***REMOVED***
		"comment0",
		"abc<b><!-- skipme --></b>def",
		"abc$<b>$<!-- skipme -->$</b>$def",
	***REMOVED***,
	***REMOVED***
		"comment1",
		"a<!-->z",
		"a$<!---->$z",
	***REMOVED***,
	***REMOVED***
		"comment2",
		"a<!--->z",
		"a$<!---->$z",
	***REMOVED***,
	***REMOVED***
		"comment3",
		"a<!--x>-->z",
		"a$<!--x>-->$z",
	***REMOVED***,
	***REMOVED***
		"comment4",
		"a<!--x->-->z",
		"a$<!--x->-->$z",
	***REMOVED***,
	***REMOVED***
		"comment5",
		"a<!>z",
		"a$<!---->$z",
	***REMOVED***,
	***REMOVED***
		"comment6",
		"a<!->z",
		"a$<!----->$z",
	***REMOVED***,
	***REMOVED***
		"comment7",
		"a<!---<>z",
		"a$<!---<>z-->",
	***REMOVED***,
	***REMOVED***
		"comment8",
		"a<!--z",
		"a$<!--z-->",
	***REMOVED***,
	***REMOVED***
		"comment9",
		"a<!--z-",
		"a$<!--z-->",
	***REMOVED***,
	***REMOVED***
		"comment10",
		"a<!--z--",
		"a$<!--z-->",
	***REMOVED***,
	***REMOVED***
		"comment11",
		"a<!--z---",
		"a$<!--z--->",
	***REMOVED***,
	***REMOVED***
		"comment12",
		"a<!--z----",
		"a$<!--z---->",
	***REMOVED***,
	***REMOVED***
		"comment13",
		"a<!--x--!>z",
		"a$<!--x-->$z",
	***REMOVED***,
	// An attribute with a backslash.
	***REMOVED***
		"backslash",
		`<p id="a\"b">`,
		`<p id="a\" b"="">`,
	***REMOVED***,
	// Entities, tag name and attribute key lower-casing, and whitespace
	// normalization within a tag.
	***REMOVED***
		"tricky",
		"<p \t\n iD=\"a&quot;B\"  foo=\"bar\"><EM>te&lt;&amp;;xt</em></p>",
		`<p id="a&#34;B" foo="bar">$<em>$te&lt;&amp;;xt$</em>$</p>`,
	***REMOVED***,
	// A nonexistent entity. Tokenizing and converting back to a string should
	// escape the "&" to become "&amp;".
	***REMOVED***
		"noSuchEntity",
		`<a b="c&noSuchEntity;d">&lt;&alsoDoesntExist;&`,
		`<a b="c&amp;noSuchEntity;d">$&lt;&amp;alsoDoesntExist;&amp;`,
	***REMOVED***,
	***REMOVED***
		"entity without semicolon",
		`&notit;&notin;<a b="q=z&amp=5&notice=hello&not;=world">`,
		`¬it;∉$<a b="q=z&amp;amp=5&amp;notice=hello¬=world">`,
	***REMOVED***,
	***REMOVED***
		"entity with digits",
		"&frac12;",
		"½",
	***REMOVED***,
	// Attribute tests:
	// http://dev.w3.org/html5/pf-summary/Overview.html#attributes
	***REMOVED***
		"Empty attribute",
		`<input disabled FOO>`,
		`<input disabled="" foo="">`,
	***REMOVED***,
	***REMOVED***
		"Empty attribute, whitespace",
		`<input disabled FOO >`,
		`<input disabled="" foo="">`,
	***REMOVED***,
	***REMOVED***
		"Unquoted attribute value",
		`<input value=yes FOO=BAR>`,
		`<input value="yes" foo="BAR">`,
	***REMOVED***,
	***REMOVED***
		"Unquoted attribute value, spaces",
		`<input value = yes FOO = BAR>`,
		`<input value="yes" foo="BAR">`,
	***REMOVED***,
	***REMOVED***
		"Unquoted attribute value, trailing space",
		`<input value=yes FOO=BAR >`,
		`<input value="yes" foo="BAR">`,
	***REMOVED***,
	***REMOVED***
		"Single-quoted attribute value",
		`<input value='yes' FOO='BAR'>`,
		`<input value="yes" foo="BAR">`,
	***REMOVED***,
	***REMOVED***
		"Single-quoted attribute value, trailing space",
		`<input value='yes' FOO='BAR' >`,
		`<input value="yes" foo="BAR">`,
	***REMOVED***,
	***REMOVED***
		"Double-quoted attribute value",
		`<input value="I'm an attribute" FOO="BAR">`,
		`<input value="I&#39;m an attribute" foo="BAR">`,
	***REMOVED***,
	***REMOVED***
		"Attribute name characters",
		`<meta http-equiv="content-type">`,
		`<meta http-equiv="content-type">`,
	***REMOVED***,
	***REMOVED***
		"Mixed attributes",
		`a<P V="0 1" w='2' X=3 y>z`,
		`a$<p v="0 1" w="2" x="3" y="">$z`,
	***REMOVED***,
	***REMOVED***
		"Attributes with a solitary single quote",
		`<p id=can't><p id=won't>`,
		`<p id="can&#39;t">$<p id="won&#39;t">`,
	***REMOVED***,
***REMOVED***

func TestTokenizer(t *testing.T) ***REMOVED***
loop:
	for _, tt := range tokenTests ***REMOVED***
		z := NewTokenizer(strings.NewReader(tt.html))
		if tt.golden != "" ***REMOVED***
			for i, s := range strings.Split(tt.golden, "$") ***REMOVED***
				if z.Next() == ErrorToken ***REMOVED***
					t.Errorf("%s token %d: want %q got error %v", tt.desc, i, s, z.Err())
					continue loop
				***REMOVED***
				actual := z.Token().String()
				if s != actual ***REMOVED***
					t.Errorf("%s token %d: want %q got %q", tt.desc, i, s, actual)
					continue loop
				***REMOVED***
			***REMOVED***
		***REMOVED***
		z.Next()
		if z.Err() != io.EOF ***REMOVED***
			t.Errorf("%s: want EOF got %q", tt.desc, z.Err())
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestMaxBuffer(t *testing.T) ***REMOVED***
	// Exceeding the maximum buffer size generates ErrBufferExceeded.
	z := NewTokenizer(strings.NewReader("<" + strings.Repeat("t", 10)))
	z.SetMaxBuf(5)
	tt := z.Next()
	if got, want := tt, ErrorToken; got != want ***REMOVED***
		t.Fatalf("token type: got: %v want: %v", got, want)
	***REMOVED***
	if got, want := z.Err(), ErrBufferExceeded; got != want ***REMOVED***
		t.Errorf("error type: got: %v want: %v", got, want)
	***REMOVED***
	if got, want := string(z.Raw()), "<tttt"; got != want ***REMOVED***
		t.Fatalf("buffered before overflow: got: %q want: %q", got, want)
	***REMOVED***
***REMOVED***

func TestMaxBufferReconstruction(t *testing.T) ***REMOVED***
	// Exceeding the maximum buffer size at any point while tokenizing permits
	// reconstructing the original input.
tests:
	for _, test := range tokenTests ***REMOVED***
		for maxBuf := 1; ; maxBuf++ ***REMOVED***
			r := strings.NewReader(test.html)
			z := NewTokenizer(r)
			z.SetMaxBuf(maxBuf)
			var tokenized bytes.Buffer
			for ***REMOVED***
				tt := z.Next()
				tokenized.Write(z.Raw())
				if tt == ErrorToken ***REMOVED***
					if err := z.Err(); err != io.EOF && err != ErrBufferExceeded ***REMOVED***
						t.Errorf("%s: unexpected error: %v", test.desc, err)
					***REMOVED***
					break
				***REMOVED***
			***REMOVED***
			// Anything tokenized along with untokenized input or data left in the reader.
			assembled, err := ioutil.ReadAll(io.MultiReader(&tokenized, bytes.NewReader(z.Buffered()), r))
			if err != nil ***REMOVED***
				t.Errorf("%s: ReadAll: %v", test.desc, err)
				continue tests
			***REMOVED***
			if got, want := string(assembled), test.html; got != want ***REMOVED***
				t.Errorf("%s: reassembled html:\n got: %q\nwant: %q", test.desc, got, want)
				continue tests
			***REMOVED***
			// EOF indicates that we completed tokenization and hence found the max
			// maxBuf that generates ErrBufferExceeded, so continue to the next test.
			if z.Err() == io.EOF ***REMOVED***
				break
			***REMOVED***
		***REMOVED*** // buffer sizes
	***REMOVED*** // tests
***REMOVED***

func TestPassthrough(t *testing.T) ***REMOVED***
	// Accumulating the raw output for each parse event should reconstruct the
	// original input.
	for _, test := range tokenTests ***REMOVED***
		z := NewTokenizer(strings.NewReader(test.html))
		var parsed bytes.Buffer
		for ***REMOVED***
			tt := z.Next()
			parsed.Write(z.Raw())
			if tt == ErrorToken ***REMOVED***
				break
			***REMOVED***
		***REMOVED***
		if got, want := parsed.String(), test.html; got != want ***REMOVED***
			t.Errorf("%s: parsed output:\n got: %q\nwant: %q", test.desc, got, want)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestBufAPI(t *testing.T) ***REMOVED***
	s := "0<a>1</a>2<b>3<a>4<a>5</a>6</b>7</a>8<a/>9"
	z := NewTokenizer(bytes.NewBufferString(s))
	var result bytes.Buffer
	depth := 0
loop:
	for ***REMOVED***
		tt := z.Next()
		switch tt ***REMOVED***
		case ErrorToken:
			if z.Err() != io.EOF ***REMOVED***
				t.Error(z.Err())
			***REMOVED***
			break loop
		case TextToken:
			if depth > 0 ***REMOVED***
				result.Write(z.Text())
			***REMOVED***
		case StartTagToken, EndTagToken:
			tn, _ := z.TagName()
			if len(tn) == 1 && tn[0] == 'a' ***REMOVED***
				if tt == StartTagToken ***REMOVED***
					depth++
				***REMOVED*** else ***REMOVED***
					depth--
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	u := "14567"
	v := string(result.Bytes())
	if u != v ***REMOVED***
		t.Errorf("TestBufAPI: want %q got %q", u, v)
	***REMOVED***
***REMOVED***

func TestConvertNewlines(t *testing.T) ***REMOVED***
	testCases := map[string]string***REMOVED***
		"Mac\rDOS\r\nUnix\n":    "Mac\nDOS\nUnix\n",
		"Unix\nMac\rDOS\r\n":    "Unix\nMac\nDOS\n",
		"DOS\r\nDOS\r\nDOS\r\n": "DOS\nDOS\nDOS\n",
		"":         "",
		"\n":       "\n",
		"\n\r":     "\n\n",
		"\r":       "\n",
		"\r\n":     "\n",
		"\r\n\n":   "\n\n",
		"\r\n\r":   "\n\n",
		"\r\n\r\n": "\n\n",
		"\r\r":     "\n\n",
		"\r\r\n":   "\n\n",
		"\r\r\n\n": "\n\n\n",
		"\r\r\r\n": "\n\n\n",
		"\r \n":    "\n \n",
		"xyz":      "xyz",
	***REMOVED***
	for in, want := range testCases ***REMOVED***
		if got := string(convertNewlines([]byte(in))); got != want ***REMOVED***
			t.Errorf("input %q: got %q, want %q", in, got, want)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestReaderEdgeCases(t *testing.T) ***REMOVED***
	const s = "<p>An io.Reader can return (0, nil) or (n, io.EOF).</p>"
	testCases := []io.Reader***REMOVED***
		&zeroOneByteReader***REMOVED***s: s***REMOVED***,
		&eofStringsReader***REMOVED***s: s***REMOVED***,
		&stuckReader***REMOVED******REMOVED***,
	***REMOVED***
	for i, tc := range testCases ***REMOVED***
		got := []TokenType***REMOVED******REMOVED***
		z := NewTokenizer(tc)
		for ***REMOVED***
			tt := z.Next()
			if tt == ErrorToken ***REMOVED***
				break
			***REMOVED***
			got = append(got, tt)
		***REMOVED***
		if err := z.Err(); err != nil && err != io.EOF ***REMOVED***
			if err != io.ErrNoProgress ***REMOVED***
				t.Errorf("i=%d: %v", i, err)
			***REMOVED***
			continue
		***REMOVED***
		want := []TokenType***REMOVED***
			StartTagToken,
			TextToken,
			EndTagToken,
		***REMOVED***
		if !reflect.DeepEqual(got, want) ***REMOVED***
			t.Errorf("i=%d: got %v, want %v", i, got, want)
			continue
		***REMOVED***
	***REMOVED***
***REMOVED***

// zeroOneByteReader is like a strings.Reader that alternates between
// returning 0 bytes and 1 byte at a time.
type zeroOneByteReader struct ***REMOVED***
	s string
	n int
***REMOVED***

func (r *zeroOneByteReader) Read(p []byte) (int, error) ***REMOVED***
	if len(p) == 0 ***REMOVED***
		return 0, nil
	***REMOVED***
	if len(r.s) == 0 ***REMOVED***
		return 0, io.EOF
	***REMOVED***
	r.n++
	if r.n%2 != 0 ***REMOVED***
		return 0, nil
	***REMOVED***
	p[0], r.s = r.s[0], r.s[1:]
	return 1, nil
***REMOVED***

// eofStringsReader is like a strings.Reader but can return an (n, err) where
// n > 0 && err != nil.
type eofStringsReader struct ***REMOVED***
	s string
***REMOVED***

func (r *eofStringsReader) Read(p []byte) (int, error) ***REMOVED***
	n := copy(p, r.s)
	r.s = r.s[n:]
	if r.s != "" ***REMOVED***
		return n, nil
	***REMOVED***
	return n, io.EOF
***REMOVED***

// stuckReader is an io.Reader that always returns no data and no error.
type stuckReader struct***REMOVED******REMOVED***

func (*stuckReader) Read(p []byte) (int, error) ***REMOVED***
	return 0, nil
***REMOVED***

const (
	rawLevel = iota
	lowLevel
	highLevel
)

func benchmarkTokenizer(b *testing.B, level int) ***REMOVED***
	buf, err := ioutil.ReadFile("testdata/go1.html")
	if err != nil ***REMOVED***
		b.Fatalf("could not read testdata/go1.html: %v", err)
	***REMOVED***
	b.SetBytes(int64(len(buf)))
	runtime.GC()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		z := NewTokenizer(bytes.NewBuffer(buf))
		for ***REMOVED***
			tt := z.Next()
			if tt == ErrorToken ***REMOVED***
				if err := z.Err(); err != nil && err != io.EOF ***REMOVED***
					b.Fatalf("tokenizer error: %v", err)
				***REMOVED***
				break
			***REMOVED***
			switch level ***REMOVED***
			case rawLevel:
				// Calling z.Raw just returns the raw bytes of the token. It does
				// not unescape &lt; to <, or lower-case tag names and attribute keys.
				z.Raw()
			case lowLevel:
				// Caling z.Text, z.TagName and z.TagAttr returns []byte values
				// whose contents may change on the next call to z.Next.
				switch tt ***REMOVED***
				case TextToken, CommentToken, DoctypeToken:
					z.Text()
				case StartTagToken, SelfClosingTagToken:
					_, more := z.TagName()
					for more ***REMOVED***
						_, _, more = z.TagAttr()
					***REMOVED***
				case EndTagToken:
					z.TagName()
				***REMOVED***
			case highLevel:
				// Calling z.Token converts []byte values to strings whose validity
				// extend beyond the next call to z.Next.
				z.Token()
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func BenchmarkRawLevelTokenizer(b *testing.B)  ***REMOVED*** benchmarkTokenizer(b, rawLevel) ***REMOVED***
func BenchmarkLowLevelTokenizer(b *testing.B)  ***REMOVED*** benchmarkTokenizer(b, lowLevel) ***REMOVED***
func BenchmarkHighLevelTokenizer(b *testing.B) ***REMOVED*** benchmarkTokenizer(b, highLevel) ***REMOVED***
